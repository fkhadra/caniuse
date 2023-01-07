package context

import (
	"caniuse/pkg/db"
)

type AppContext struct {
	Db     *db.Db
	Screen struct {
		Width  int
		Height int
	}
}

func New() *AppContext { return &AppContext{} }

func (c *AppContext) SetDb(db *db.Db) {
	c.Db = db
}

func (c *AppContext) SetSize(w int, h int) {
	c.Screen.Width, c.Screen.Height = w, h
}
