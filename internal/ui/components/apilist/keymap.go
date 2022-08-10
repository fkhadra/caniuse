package apilist

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	up                   key.Binding
	down                 key.Binding
	viewDetails          key.Binding
	clearFilter          key.Binding
	acceptWhileFiltering key.Binding
	filter               key.Binding
	quit                 key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.up,
		k.down,
		k.viewDetails,
		k.filter,
		k.quit,
	}
}

// unused for now
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{}, // first column
		{}, // second column
	}
}

var keys = keyMap{
	up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	viewDetails: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "view details"),
	),
	filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	clearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	acceptWhileFiltering: key.NewBinding(
		key.WithKeys("enter", "tab", "up", "down"),
		key.WithHelp("enter", "apply filter"),
	),
	quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}
