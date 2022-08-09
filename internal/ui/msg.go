package ui

import (
	"caniuse/pkg/db"

	tea "github.com/charmbracelet/bubbletea"
)

type dbLoadedMsg struct {
	err error
	db  *db.Db
}

func loadDb() tea.Msg {
	d, err := db.Init()

	return dbLoadedMsg{
		db:  d,
		err: err,
	}
}
