package supportable

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	goBack key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.goBack,
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
	goBack: key.NewBinding(
		key.WithKeys("esc", "ctrl+c", "/"),
		key.WithHelp("esc", "go back"),
	),
}
