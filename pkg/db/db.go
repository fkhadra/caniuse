package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"syscall"
)

type Api struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Spec        string `json:"spec"`
	Status      string `json:"status"`
	Links       []struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	} `json:"links"`
	Categories []string                  `json:"categories"`
	Support    map[string]BrowserSupport `json:"stats"`
	Notes      map[string]string         `json:"notes_by_num"`
	Usage      float64                   `json:"usage_perc_y"`
}

type BrowserInfo struct {
	Name        string             `json:"browser"`
	LongName    string             `json:"long_name"`
	Abbr        string             `json:"abbr"`
	Prefix      string             `json:"prefix"`
	Type        string             `json:"type"`
	UsageGlobal map[string]float64 `json:"usage_global"`
}

type Data struct {
	Browser map[string]BrowserInfo `json:"agents"`
	Api     map[string]Api         `json:"data"`
	Updated int                    `json:"updated"`
}

type Db struct {
	Data *Data

	browserIds []string
}

type BrowserSupport [][2]string

// aggregate data while unmarshalling
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

		if startVersion == "" && isLastEntry {
			out = append(out,
				[2]string{version, support},
			)
		} else if startVersion == "" {
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

const Url = "https://raw.githubusercontent.com/Fyrd/caniuse/main/data.json"

var Dir = func() string {
	dir, _ := os.UserHomeDir()
	return path.Join(dir, ".config/caniuse")
}()

var Pathname = func() string {
	return path.Join(Dir, "db.json")
}()

func Init() (db *Db, err error) {
	db, err = load(Pathname)

	if err != nil && errors.Is(err, syscall.ENOENT) {

		if err = os.MkdirAll(Dir, 0755); err != nil {
			return
		}

		// db does not exist, let's get it
		if err = Download(Pathname); err != nil {
			return
		}
		// load db again
		if db, err = load(Pathname); err != nil {
			//something is really wrong
			return
		}
	}

	return
}

// used to iterate while keeping the order
func (d *Db) BrowserIds() []string {
	if d.browserIds != nil {
		return d.browserIds
	}

	keys := make([]string, len(d.Data.Browser))
	i := 0

	for id := range d.Data.Browser {
		keys[i] = id
		i++
	}

	sort.Strings(keys)
	d.browserIds = keys

	return keys
}

func (d *Db) CheckForUpdate() (updated bool, err error) {
	dbPathname := path.Join(os.TempDir(), "caniuse.json")

	// best effort to remove dangling crap, hence no error handling 🙈
	defer os.Remove(dbPathname)

	if err = Download(dbPathname); err != nil {
		return
	}

	latestDb, err := load(dbPathname)
	if err != nil {
		return
	}

	if latestDb.Data.Updated > d.Data.Updated {
		if err = os.Rename(dbPathname, Pathname); err != nil {
			return
		}

		d.Data = latestDb.Data
		updated = true
	}

	return
}

func (d *Db) IsDesktopBrowser(id string) bool {
	return d.Data.Browser[id].Type == "desktop"
}

func load(dbPathname string) (db *Db, err error) {
	f, err := os.Open(dbPathname)
	if err != nil {
		return
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return
	}

	var data Data
	err = json.Unmarshal(b, &data)

	if err != nil {
		return
	}

	db = &Db{
		Data: &data,
	}

	return
}

func Download(dbPathname string) (err error) {
	resp, err := http.Get(Url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(dbPathname)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return
}
