package tab

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		Foreground(lipgloss.Color("#666")).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#c3f", Dark: "#c3f"}).
		Padding(0, 1)

	activeTab = tab.Copy().
			Foreground(lipgloss.NoColor{}).
			Bold(true).
			Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

type Model struct {
	activeTab int
	tabs      []string
	// width          int
}

func New(tabs []string) Model {
	return Model{
		tabs: tabs,
	}
}

// func (m *Model) SetWidth(w int) {
// 	m.width = w
// }

func (m *Model) NextTab() {
	nextTab := m.activeTab + 1
	if nextTab > len(m.tabs)-1 {
		m.activeTab = 0
	} else {
		m.activeTab = nextTab
	}
}

func (m *Model) ActiveTab() int {
	return m.activeTab
}

func (m Model) View(width int) string {
	tabCount := len(m.tabs)
	out := make([]string, tabCount)

	for i := 0; i < tabCount; i++ {
		if i == m.activeTab {
			out[i] = activeTab.Render(m.tabs[i])
		} else {
			out[i] = tab.Render(m.tabs[i])
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, out...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	return fmt.Sprintf("%s\n\n", row)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
