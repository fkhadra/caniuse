package apilist

import (
	"caniuse/internal/theme"
	"caniuse/internal/ui/components/supportable"
	"caniuse/internal/ui/context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	list         list.Model
	searchInput  textinput.Model
	supportTable supportable.Model
	ctx          *context.AppContext
}

func New(ctx *context.AppContext) Model {
	l := list.New(make([]list.Item, 0), item{}, 0, 0)

	l.SetStatusBarItemName("entry", "entries")
	l.SetShowFilter(false)
	l.SetShowTitle(false)
	l.KeyMap.CursorUp = keys.up
	l.KeyMap.CursorDown = keys.down
	l.KeyMap.CancelWhileFiltering = keys.clearFilter
	l.KeyMap.AcceptWhileFiltering = keys.acceptWhileFiltering
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.PrevPage.SetKeys("left")
	l.KeyMap.NextPage.SetKeys("right")
	l.AdditionalShortHelpKeys = func() []key.Binding { return []key.Binding{keys.viewDetails} }
	l.Help.Styles.ShortKey = theme.HelpKey

	ti := textinput.New()
	ti.Placeholder = "..."
	ti.Prompt = ""
	ti.CharLimit = 20
	ti.CursorStyle = lipgloss.NewStyle().Foreground(theme.ColorHighlight)

	return Model{
		list:         l,
		searchInput:  ti,
		supportTable: supportable.New(ctx),
		ctx:          ctx,
	}
}

func (m *Model) Init() tea.Cmd {
	items := make([]list.Item, len(m.ctx.Db.Api))

	i := 0
	for id, v := range m.ctx.Db.Api {
		items[i] = item{
			id:          id,
			title:       v.Title,
			description: v.Description,
			usage:       v.Usage,
			categories:  v.Categories,
		}
		i++
	}

	m.list.SetSize(m.ctx.Screen.Width, m.ctx.Screen.Height-lipgloss.Height(m.renderSearchInput()))

	return tea.Batch(
		m.list.SetItems(items),
		textinput.Blink,
		m.searchInput.Focus(),
		list.EnableLiveFiltering,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.searchInput.Focused() {
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds = append(cmds, m.handleKeyboard(msg))
	}

	if !m.supportTable.IsActive() {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.supportTable, cmd = m.supportTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeyboard(msg tea.KeyMsg) tea.Cmd {
	var (
		cmd         tea.Cmd
		isFiltering = m.list.FilterState() == list.Filtering ||
			m.list.FilterState() == list.FilterApplied
	)

	switch {
	case key.Matches(msg, keys.up, keys.down):
		if isFiltering {
			m.searchInput.Blur()
			m.list.KeyMap.AcceptWhileFiltering.SetEnabled(true)
		}
	case key.Matches(msg, keys.filter):
		return tea.Batch(m.searchInput.Focus(), list.EnableLiveFiltering)
	case key.Matches(msg, keys.clearFilter):
		if !m.supportTable.IsActive() {
			m.searchInput.SetValue("")
			m.searchInput.Blur()
		} else {
			cmd = tea.Batch(m.searchInput.Focus(), list.EnableLiveFiltering)
		}
	case key.Matches(msg, keys.viewDetails):
		if m.list.FilterState() == list.Filtering {
			m.searchInput.Blur()
		} else if v, ok := m.list.SelectedItem().(item); ok {
			m.supportTable.SetApiId(v.id)
		}
	}

	return cmd
}

func (m Model) View() string {
	var (
		s    strings.Builder
		body string
	)

	if m.supportTable.IsActive() {
		body = m.supportTable.View()
	} else {
		body = m.list.View()
	}

	fmt.Fprintf(&s, "%s%s",
		m.renderSearchInput(),
		body,
	)

	return s.String()
}

func (m Model) Placeholder() string {
	return m.renderSearchInput()
}

func (m Model) renderSearchInput() string {
	input := style.searchInput.
		MaxWidth(m.ctx.Screen.Width).
		Render(fmt.Sprintf("\nCan I use %s ?", theme.Bold.Render(m.searchInput.View())))

	return fmt.Sprintf("%s\n\n", style.center.Width(m.ctx.Screen.Width).Render(input))
}

var style = func() (s struct {
	searchInput lipgloss.Style
	center      lipgloss.Style
}) {
	s.searchInput = lipgloss.NewStyle().
		BorderForeground(theme.ColorMagenta).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(80).
		Height(3).
		Margin(1).
		Padding(0, 2, 0, 2).
		Align(lipgloss.Center)

	s.center = lipgloss.NewStyle().Align(lipgloss.Center)

	return s
}()
