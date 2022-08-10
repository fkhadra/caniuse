package db_test

import (
	"caniuse/pkg/db"
	"errors"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func removeConfigDir() {
	os.RemoveAll(db.DefaultConfigDir())
}

func TestDownloadDatabaseIfNotExist(t *testing.T) {
	removeConfigDir()
	_, err := os.Stat(db.DefaultConfigDir())

	// ensure config is removed
	assert.True(t, errors.Is(err, syscall.ENOENT))
	d, err := db.Init()

	assert.Nil(t, err)
	assert.NotEmpty(t, d)
}

func TestCheckForUpdate(t *testing.T) {
	d, _ := db.Init()

	// make current db old
	d.Updated = 0

	updated, err := d.CheckForUpdate()
	assert.Nil(t, err)
	assert.True(t, updated)

	// load updated database
	d, _ = db.Init()

	updated, err = d.CheckForUpdate()
	assert.Nil(t, err)
	assert.False(t, updated)
}

func TestBrowserIds(t *testing.T) {
	d, _ := db.Init()
	expected := []string{
		"and_chr",
		"and_ff",
		"and_qq",
		"and_uc",
		"android",
		"baidu",
		"bb",
		"chrome",
		"edge",
		"firefox",
		"ie",
		"ie_mob",
		"ios_saf",
		"kaios",
		"op_mini",
		"op_mob",
		"opera",
		"safari",
		"samsung",
	}

	assert.Equal(t, expected, d.BrowserIds())
	// cover caching
	assert.Equal(t, expected, d.BrowserIds())
}

func TestIsDesktopBrowser(t *testing.T) {
	d, _ := db.Init()

	assert.True(t, d.IsDesktopBrowser("chrome"))
	assert.False(t, d.IsDesktopBrowser("ios_saf"))
}
