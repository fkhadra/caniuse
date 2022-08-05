package db

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"syscall"
)

type ApiLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Api struct {
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Spec        string                       `json:"spec"`
	Status      string                       `json:"status"`
	Links       []ApiLink                    `json:"links"`
	Categories  []string                     `json:"categories"`
	Stats       map[string]map[string]string `json:"stats"`
	Usage       float64                      `json:"usage_perc_y"`
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
	Browser map[string]BrowserInfo `json:"agents"`
	Api     map[string]Api         `json:"data"`
	Updated int                    `json:"updated"`
}

const Url = "https://raw.githubusercontent.com/Fyrd/caniuse/main/data.json"

var (
	DefaultConfigDir = func() string {
		homedir, _ := os.UserHomeDir()
		return path.Join(homedir, ".config/caniuse")
	}()

	DefaultDbPathname = path.Join(DefaultConfigDir, "data.json")
)

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
		if err = os.Rename(dbPathname, DefaultDbPathname); err != nil {
			return
		}
		updated = true
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
