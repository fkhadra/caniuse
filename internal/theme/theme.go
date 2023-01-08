package theme

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	ColorError     = lipgloss.AdaptiveColor{Light: "#F23D5C", Dark: "#F23D5C"}
	ColorSuccess   = lipgloss.AdaptiveColor{Light: "#00A300", Dark: "#00A300"}
	ColorWarning   = lipgloss.AdaptiveColor{Light: "#f1c40f", Dark: "#f1c40f"}
	ColorInfo      = lipgloss.AdaptiveColor{Light: "#3498db", Dark: "#3498db"}
	ColorNeutral   = lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}
	ColorHighlight = lipgloss.AdaptiveColor{Light: "#87D", Dark: "#87D"}
	ColorMagenta   = lipgloss.AdaptiveColor{Light: "#ba007a", Dark: "#ba007a"}
	SpinnerColor   = ColorMagenta

	Bold = lipgloss.NewStyle().Bold(true)

	Text           = lipgloss.NewStyle()
	TextInfo       = Text.Copy().Foreground(ColorInfo)
	TextWarning    = Text.Copy().Foreground(ColorWarning)
	TextError      = Text.Copy().Foreground(ColorError)
	TextSuccess    = Text.Copy().Foreground(ColorSuccess)
	TextHightlight = Text.Copy().Foreground(ColorHighlight)

	Badge = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fffff1")).
		Bold(true).
		Padding(0, 1, 0, 1)
	BadgeSuccess = Badge.Copy().Background(ColorSuccess)
	BadgeWarning = Badge.Copy().Background(ColorWarning)
	BadgeError   = Badge.Copy().Background(ColorError)
	BadgeInfo    = Badge.Copy().Background(ColorInfo)
	BadgeNeutral = Badge.Copy().Background(ColorNeutral).Foreground(lipgloss.Color("#fff"))

	HelpKey = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorInfo)
)

func RenderBrowserSupport(u float64) string {
	var color lipgloss.Style
	if u >= 0 && u <= 25 {
		color = BadgeError
	} else if u >= 26 && u <= 75 {
		color = BadgeWarning
	} else {
		color = BadgeSuccess
	}

	return fmt.Sprintf("%s%s",
		BadgeNeutral.Render("Browser support"),
		color.Render(fmt.Sprintf("%.2f%%", u)))
}
