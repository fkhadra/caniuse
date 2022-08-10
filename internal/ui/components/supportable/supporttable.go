package supportable

import (
	"caniuse/internal/theme"
	"caniuse/internal/ui/context"
	"caniuse/pkg/superscript"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

type Model struct {
	apiId   string
	help    help.Model
	context *context.AppContext
}

// used to parse notes
var reNotes = regexp.MustCompile(`#(\d)`)

func New(ctx *context.AppContext) Model {
	h := help.New()
	h.Styles.ShortKey = theme.HelpKey

	return Model{
		help:    h,
		context: ctx,
	}
}

func (m *Model) SetApiId(id string) {
	m.apiId = id
}

func (m Model) IsActive() bool {
	return m.apiId != ""
}

func (m *Model) clearApiId() {
	m.apiId = ""
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.goBack):
			m.clearApiId()
		}
	}
	return m, nil
}

func (m Model) View() string {
	var (
		s              strings.Builder
		desktopColumns []string
		mobileColumns  []string
	)
	apiData := m.context.Db.Api[m.apiId]

	fmt.Fprintf(&s, "%s\n\n%s\n\n",
		style.title.Render(apiData.Title),
		wrap.String(apiData.Description, m.context.Screen.Width-style.body.GetHorizontalPadding()),
	)

	for _, id := range m.context.Db.BrowserIds() {
		var column strings.Builder
		isDesktop := m.context.Db.IsDesktopBrowser(id)

		browserName := style.browserName.Render(m.context.Db.Browser[id].Browser)

		fmt.Fprintf(&column, "%s\n", browserName)

		if !isDesktop && lipgloss.Height(browserName) == 1 {
			column.WriteString("\n")
		}

		for _, supp := range apiData.Support[id] {
			version, support := supp[0], supp[1]

			if support == "y" {
				column.WriteString(style.supported.Render(version))
			} else if support == "n" {
				column.WriteString(style.notSupported.Render(version))
			} else {
				var notes strings.Builder
				for _, match := range reNotes.FindAllStringSubmatch(support, -1) {
					n, _ := strconv.Atoi(match[1])
					notes.WriteString(superscript.Itoa(n))
				}
				column.WriteString(style.partial.Render(version + notes.String()))
			}

			column.WriteString(" \n")
		}

		if isDesktop {
			desktopColumns = append(desktopColumns, column.String())
		} else {
			mobileColumns = append(mobileColumns, column.String())
		}
	}

	fmt.Fprintf(&s, "%s\n\n%s\n\n\n%s\n\n%s\n\n%s",
		theme.BadgeInfo.Render("Desktop"),
		lipgloss.JoinHorizontal(lipgloss.Top, desktopColumns...),
		theme.BadgeInfo.Render("Mobile"),
		lipgloss.JoinHorizontal(lipgloss.Top, mobileColumns...),
		m.help.View(keys),
	)

	return style.body.Render(s.String())
}

type stylesheet struct {
	body         lipgloss.Style
	browserName  lipgloss.Style
	title        lipgloss.Style
	supported    lipgloss.Style
	partial      lipgloss.Style
	notSupported lipgloss.Style
}

var style = func() (s stylesheet) {
	minWidth := 12
	baseBadge := theme.Badge.Copy().
		Width(minWidth).
		Align(lipgloss.Center)

	s.body = lipgloss.NewStyle().Padding(0, 2, 0, 2)
	s.title = theme.Bold.Copy().Underline(true)
	s.supported = baseBadge.Copy().Inherit(theme.BadgeSuccess)
	s.partial = baseBadge.Copy().Inherit(theme.BadgeWarning)
	s.notSupported = baseBadge.Copy().Inherit(theme.BadgeError)

	s.browserName = theme.Text.Copy().
		Width(minWidth).
		Align(lipgloss.Center)

	return s
}()
