package ui

import (
	"caniuse/pkg/db"

	tea "github.com/charmbracelet/bubbletea"
)

type dbLoadedMsg struct {
	err error
	db  *db.Db
}

func (m *Model) loadDb() tea.Msg {
	d, err := db.Init()

	return dbLoadedMsg{
		db:  d,
		err: err,
	}
}

type dbUpdatedMsg struct {
	err error
}

func (m *Model) updateDb() tea.Msg {
	return dbUpdatedMsg{
		err: m.ctx.Db.CheckForUpdate(),
	}
}
