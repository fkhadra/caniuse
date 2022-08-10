package apilist

import (
	"caniuse/internal/theme"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
)

type item struct {
	id          string
	title       string
	description string
	usage       float64
	categories  []string
}

func (i item) FilterValue() string                       { return i.title }
func (i item) Height() int                               { return 2 }
func (i item) Spacing() int                              { return 1 }
func (i item) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (t item) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	el, ok := listItem.(item)
	if !ok {
		return
	}

	itemMaxWidth := m.Width() - itemStyle.title.GetPaddingLeft() - itemStyle.title.GetPaddingRight()
	desc := wordwrap.String(el.description, itemMaxWidth)
	lines := strings.Split(desc, "\n")
	descMaxLen := t.Height() - 1
	descLastLine := descMaxLen - 1

	if len(lines) > descMaxLen {
		lines[descLastLine] = truncate.StringWithTail(lines[descLastLine], uint(itemMaxWidth-10), "...")
		lines = lines[0:descMaxLen]
	}

	desc = strings.Join(lines, "\n")

	var (
		isSelected = index == m.Index()
		// emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied

		matchedRunes = m.MatchesForItem(index)
		title        = el.title
	)

	if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			unmatched := itemStyle.activeTitle.Inline(true)
			matched := unmatched.Copy().Inherit(itemStyle.filterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = itemStyle.activeTitle.Render(title)
		desc = itemStyle.activeDescription.Render(desc)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := itemStyle.title.Inline(true)
			matched := unmatched.Copy().Inherit(itemStyle.filterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = itemStyle.title.Render(title)
		desc = itemStyle.description.Render(desc)
	}

	fmt.Fprintf(w, "%s%s\n%s",
		itemStyle.title.Render(title),
		renderUsage(el.usage),
		itemStyle.description.Render(desc),
	)
}

func renderUsage(u float64) string {
	var color lipgloss.Style
	if u >= 0 && u <= 25 {
		color = theme.BadgeError
	} else if u >= 26 && u <= 75 {
		color = theme.BadgeWarning
	} else {
		color = theme.BadgeSuccess
	}

	return fmt.Sprintf("%s%s",
		theme.BadgeNeutral.Render("Browser support"),
		color.Render(fmt.Sprintf("%.2f%%", u)))
}

var itemStyle = func() (s struct {
	title             lipgloss.Style
	description       lipgloss.Style
	filterMatch       lipgloss.Style
	activeTitle       lipgloss.Style
	activeDescription lipgloss.Style
}) {
	s.title = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 2, 0, 2)

	s.description = s.title.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.filterMatch = lipgloss.NewStyle().
		Background(theme.ColorHighlight)

	s.activeTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)

	s.activeDescription = s.activeTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	return s
}()
