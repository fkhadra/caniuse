package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"
)

type BrowserId string

type BrowserSupport [][2]string

func (b *BrowserSupport) UnmarshalJSON(data []byte) error {
	var (
		startVersion, prevSupport, prevVersion string
		out                                    [][2]string
	)

	input := strings.Split(
		strings.ReplaceAll(
			// remove curly braces from original string
			string(data[1:len(data)-1]),
			"\"", ""),
		",")

	for i, v := range input {
		var versionRange string
		isLastEntry := i == len(input)-1
		current := strings.Split(v, ":")
		version, support := strings.TrimSpace(current[0]), strings.TrimSpace(current[1])

		// first entry
		if startVersion == "" {
			startVersion = version
		} else if isLastEntry {
			if prevSupport != support {
				versionRange = fmt.Sprintf("%s-%s", startVersion, prevVersion)
				out = append(out,
					[2]string{versionRange, prevSupport},
					[2]string{version, support},
				)
			} else {
				versionRange = fmt.Sprintf("%s-%s", startVersion, version)
				out = append(out, [2]string{versionRange, support})
			}
		} else if prevSupport != "" && prevSupport != support {
			if startVersion == prevVersion {
				versionRange = startVersion
			} else {
				versionRange = fmt.Sprintf("%s-%s", startVersion, prevVersion)
			}
			out = append(out, [2]string{versionRange, prevSupport})
			startVersion = version
		}
		prevVersion = version
		prevSupport = support
	}

	*b = out

	return nil
}

type Api struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Spec        string `json:"spec"`
	Status      string `json:"status"`
	Links       []struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	} `json:"links"`
	Categories []string                     `json:"categories"`
	Support    map[BrowserId]BrowserSupport `json:"stats"`
	Usage      float64                      `json:"usage_perc_y"`
}

type BrowserInfo struct {
	Browser     string             `json:"browser"`
	LongName    string             `json:"long_name"`
	Abbr        string             `json:"abbr"`
	Prefix      string             `json:"prefix"`
	Type        string             `json:"type"`
	UsageGlobal map[string]float64 `json:"usage_global"`
}

type Db struct {
	Browser map[BrowserId]BrowserInfo `json:"agents"`
	Api     map[string]Api            `json:"data"`
	Updated int                       `json:"updated"`
}

const Url = "https://raw.githubusercontent.com/Fyrd/caniuse/main/data.json"

func DefaultConfigDir() string {
	homedir, _ := os.UserHomeDir()
	return path.Join(homedir, ".config/caniuse")
}

func DefaultDbPathname() string {
	return path.Join(DefaultConfigDir(), "data.json")
}

func Init() (db Db, err error) {
	if err = os.MkdirAll(DefaultConfigDir, 0755); err != nil {
		return
	}

	db, err = load(DefaultDbPathname)

	if err != nil && errors.Is(err, syscall.ENOENT) {
		// db does not exist, let's get it
		if err = download(DefaultDbPathname); err != nil {
			return
		}
		// load db again
		if db, err = load(DefaultDbPathname); err != nil {
			//something is really wrong
			return
		}
	}

	return db, nil
}

func (d Db) CheckForUpdate() (updated bool, err error) {
	dbPathname := path.Join(os.TempDir(), "caniuse.json")

	if err = download(dbPathname); err != nil {
		return
	}

	maybeNewDb, err := load(dbPathname)
	if err != nil {
		return
	}

	if maybeNewDb.Updated > d.Updated {
		if err = os.Rename(dbPathname, DefaultDbPathname()); err != nil {
			return
		}
		updated = true
	} else {
		os.Remove(dbPathname)
	}

	return
}

func load(dbPathname string) (db Db, err error) {
	f, err := os.Open(dbPathname)
	if err != nil {
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return
	}

	return
}

func download(dbPathname string) error {
	resp, err := http.Get(Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dbPathname)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
