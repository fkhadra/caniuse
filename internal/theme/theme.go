package theme

import "github.com/charmbracelet/lipgloss"

var (
	ErrorColor     = lipgloss.AdaptiveColor{Light: "#F23D5C", Dark: "#F23D5C"}
	SuccessColor   = lipgloss.AdaptiveColor{Light: "#00A300", Dark: "#00A300"}
	WarningColor   = lipgloss.AdaptiveColor{Light: "#f1c40f", Dark: "#f1c40f"}
	InfoColor      = lipgloss.AdaptiveColor{Light: "#3498db", Dark: "#3498db"}
	NeutralColor   = lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}
	HighlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#874BFD"}
	MagentaColor   = lipgloss.AdaptiveColor{Light: "#ba007a", Dark: "#ba007a"}
	SpinnerColor   = MagentaColor

	Bold = lipgloss.NewStyle().Bold(true)

	Text           = lipgloss.NewStyle().Bold(true)
	TextInfo       = Text.Copy().Foreground(InfoColor)
	TextWarning    = Text.Copy().Foreground(WarningColor)
	TextError      = Text.Copy().Foreground(ErrorColor)
	TextSuccess    = Text.Copy().Foreground(SuccessColor)
	TextHightlight = Text.Copy().Foreground(HighlightColor)

	Badge = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fffff1")).
		Bold(true).
		Padding(0, 1, 0, 1)
	BadgeSuccess = Badge.Copy().Background(SuccessColor)
	BadgeWarning = Badge.Copy().Background(WarningColor)
	BadgeError   = Badge.Copy().Background(ErrorColor)
	BadgeInfo    = Badge.Copy().Background(InfoColor)
	BadgeNeutral = Badge.Copy().Background(NeutralColor).Foreground(lipgloss.Color("#fff"))

	HelpKey = lipgloss.NewStyle().
		Bold(true).
		Foreground(InfoColor)
)
