package supportable

import (
	"caniuse/internal/theme"
	"caniuse/internal/ui/context"
	"caniuse/pkg/superscript"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	context *context.AppContext
}

// used to parse notes
var reNotes = regexp.MustCompile(`#(\d)`)

func New(ctx *context.AppContext) Model {
	return Model{
		context: ctx,
	}
}

func (m Model) View(api string) string {
	var (
		s              strings.Builder
		desktopColumns []string
		mobileColumns  []string
	)
	apiData := m.context.Db.Api[api]

	fmt.Fprintf(&s, "%s\n\n%s\n\n", theme.Bold.Render(apiData.Title), apiData.Description)

	for id, v := range apiData.Support {
		var column strings.Builder
		isDesktop := m.context.Db.IsDesktopBrowser(id)

		browserName := style.browserName.Render(m.context.Db.Browser[id].Browser)

		fmt.Fprintf(&column, "%s\n", browserName)

		if !isDesktop && lipgloss.Height(browserName) == 1 {
			column.WriteString("\n")
		}

		for _, supp := range v {
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

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, desktopColumns...))
	s.WriteString("\n\n")

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, mobileColumns...))

	return s.String()
}

type stylesheet struct {
	browserName  lipgloss.Style
	supported    lipgloss.Style
	partial      lipgloss.Style
	notSupported lipgloss.Style
}

var style = func() (s stylesheet) {
	minWidth := 12
	baseBadge := theme.Badge.Copy().
		Width(minWidth).
		Align(lipgloss.Center)

	s.supported = baseBadge.Copy().Inherit(theme.BadgeSuccess)
	s.partial = baseBadge.Copy().Inherit(theme.BadgeWarning)
	s.notSupported = baseBadge.Copy().Inherit(theme.BadgeError)

	s.browserName = theme.Text.Copy().
		Width(minWidth).
		Align(lipgloss.Center)

	return s
}()
