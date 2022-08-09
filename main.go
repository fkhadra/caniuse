package main

import (
	"caniuse/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.New())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// var (
	// 	s              strings.Builder
	// 	desktopColumns []string
	// 	mobileColumns  []string
	// )
	// reNotes := regexp.MustCompile(`#(\d)`)

	// fmt.Fprintf(&s, "%s\n\n%s\n\n", theme.Bold.Render(api.Title), api.Description)

	// for id, v := range api.Support {
	// 	var column strings.Builder
	// 	isDesktop := d.Browser[id].Type == "desktop"

	// 	browserName := style.browserName.Render(d.Browser[id].Browser)

	// 	fmt.Fprintf(&column, "%s\n", browserName)

	// 	if !isDesktop && lipgloss.Height(browserName) == 1 {
	// 		column.WriteString("\n")
	// 	}

	// 	for _, supp := range v {
	// 		version, support := supp[0], supp[1]

	// 		if support == "y" {
	// 			column.WriteString(style.supported.Render(version))
	// 		} else if support == "n" {
	// 			column.WriteString(style.notSupported.Render(version))
	// 		} else {
	// 			var notes strings.Builder
	// 			for _, match := range reNotes.FindAllStringSubmatch(support, -1) {
	// 				n, _ := strconv.Atoi(match[1])
	// 				notes.WriteString(superscript.Itoa(n))
	// 			}
	// 			column.WriteString(style.partial.Render(version + notes.String()))
	// 		}

	// 		column.WriteString(" \n")
	// 	}

	// 	if isDesktop {
	// 		desktopColumns = append(desktopColumns, column.String())
	// 	} else {
	// 		mobileColumns = append(mobileColumns, column.String())
	// 	}

	// }

	// s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, desktopColumns...))
	// s.WriteString("\n\n")

	// s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, mobileColumns...))

	// fmt.Println(s.String())
}

// type tableStyles struct {
// 	browserName  lipgloss.Style
// 	supported    lipgloss.Style
// 	partial      lipgloss.Style
// 	notSupported lipgloss.Style
// }

// var style = func() (s tableStyles) {
// 	minWidth := 12
// 	baseBadge := theme.Badge.Copy().
// 		Width(minWidth).
// 		Align(lipgloss.Center)

// 	s.supported = baseBadge.Copy().Inherit(theme.BadgeSuccess)
// 	s.partial = baseBadge.Copy().Inherit(theme.BadgeWarning)
// 	s.notSupported = baseBadge.Copy().Inherit(theme.BadgeError)

// 	s.browserName = theme.Text.Copy().
// 		Width(minWidth).
// 		Align(lipgloss.Center)

// 	return s
// }()

// func f() (s style) {
// 	minWidth := 10
// 	return s
// }

// var style = struct {
// 	browserName lipgloss.Style
// }{
// 	browserName: lipgloss.NewStyle().Bold(true).Underline(true).Align(lipgloss.Center).wi,
// }
