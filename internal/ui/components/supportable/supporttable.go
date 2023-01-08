package supportable

import (
	"caniuse/internal/theme"
	"caniuse/internal/ui/components/tab"
	"caniuse/internal/ui/context"
	"caniuse/pkg/superscript"
	"fmt"
	"io"
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
	apiId string
	help  help.Model
	ctx   *context.AppContext
	tab   tab.Model
}

// used to parse notes
var reNotes = regexp.MustCompile(`#(\d)`)

func New(ctx *context.AppContext) Model {
	h := help.New()
	h.Styles.ShortKey = theme.HelpKey

	return Model{
		help: h,
		ctx:  ctx,
		tab:  tab.New([]string{"Notes", "Resources"}),
	}
}

func (m *Model) SetApiId(id string) {
	m.apiId = id
}

func (m *Model) IsActive() bool {
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
		case key.Matches(msg, keys.switchTab):
			m.tab.NextTab()
		}
	}

	return m, nil
}

const (
	noteTab = iota
	resourceTab
)

func (m Model) View() string {
	var (
		s              strings.Builder
		desktopColumns []string
		mobileColumns  []string
	)
	api := m.ctx.Db.Data.Api[m.apiId]

	fmt.Fprintf(&s, "%s %s\n\n%s\n\n",
		style.title.Render(api.Title),
		theme.RenderBrowserSupport(api.Usage),
		style.description.Render(wrap.String(api.Description, m.ctx.Screen.Width-style.body.GetHorizontalPadding())),
	)

	// building support table
	for _, id := range m.ctx.Db.BrowserIds() {
		var column strings.Builder
		isDesktop := m.ctx.Db.IsDesktopBrowser(id)
		browserName := style.browserName.Render(m.ctx.Db.Data.Browser[id].Name)

		fmt.Fprintf(&column, "%s\n", browserName)

		if !isDesktop && lipgloss.Height(browserName) == 1 {
			column.WriteString("\n")
		}

		for _, v := range api.Support[id] {
			version, supported := v[0], v[1]

			if supported == "y" {
				column.WriteString(style.supported.Render(version))
			} else if supported == "n" {
				column.WriteString(style.notSupported.Render(version))
			} else {
				// convert notes to superscript equivalent
				var notes strings.Builder
				for _, match := range reNotes.FindAllStringSubmatch(supported, -1) {
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

	fmt.Fprintf(&s, "%s\n\n%s\n\n\n%s\n\n%s\n\n",
		theme.BadgeInfo.Render("Desktop"),
		lipgloss.JoinHorizontal(lipgloss.Top, desktopColumns...),
		theme.BadgeInfo.Render("Mobile"),
		lipgloss.JoinHorizontal(lipgloss.Top, mobileColumns...),
	)

	s.WriteString(m.tab.View(m.ctx.Screen.Width))

	switch m.tab.ActiveTab() {
	case noteTab:
		for i := 1; i < len(api.Notes); i++ {
			k := strconv.Itoa(i)
			v := api.Notes[k]
			fmt.Fprintf(&s, "%s %s\n", theme.BadgeNeutral.Render(k), v)
		}
	case resourceTab:
		renderLink(&s, "Specification", api.Spec)
		for _, v := range api.Links {
			renderLink(&s, v.Title, v.URL)
		}
	}
	s.WriteString("\n")
	s.WriteString(m.help.View(keys))

	return style.body.Render(s.String())
}

func renderLink(w io.Writer, title string, link string) {
	fmt.Fprintf(w, "• %s %s\n",
		title,
		lipgloss.NewStyle().Foreground(theme.ColorNeutral).Render(link),
	)
}

type stylesheet struct {
	body         lipgloss.Style
	browserName  lipgloss.Style
	title        lipgloss.Style
	description  lipgloss.Style
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
	s.title = theme.Bold.Copy()
	s.description = theme.Text.Copy().Foreground(theme.ColorNeutral)
	s.supported = baseBadge.Copy().Inherit(theme.BadgeSuccess)
	s.partial = baseBadge.Copy().Inherit(theme.BadgeWarning)
	s.notSupported = baseBadge.Copy().Inherit(theme.BadgeError)

	s.browserName = theme.Text.Copy().
		Width(minWidth).
		Align(lipgloss.Center)

	return s
}()
