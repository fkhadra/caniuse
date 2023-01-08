package supportable

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	goBack    key.Binding
	switchTab key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.goBack,
	}
}

// unused but needed
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}, {}}
}

var keys = keyMap{
	goBack: key.NewBinding(
		key.WithKeys("esc", "ctrl+c", "/"),
		key.WithHelp("esc", "go back"),
	),
	switchTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch tab"),
	),
}
