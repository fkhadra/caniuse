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
}
