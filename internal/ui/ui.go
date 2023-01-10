package ui

import (
	"caniuse/internal/theme"
	"caniuse/internal/ui/components/apilist"
	"caniuse/internal/ui/context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Status int

const (
	StatusNotStarted = iota
	StatusInProgress
	StatusUpdating
	StatusOk
	StatusError
)

type Model struct {
	error      error
	ctx        *context.AppContext
	readyState Status
	apilist    apilist.Model
	spinner    spinner.Model
}

func New() Model {
	ctx := context.New()
	return Model{
		ctx:        ctx,
		apilist:    apilist.New(ctx),
		readyState: StatusInProgress,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(theme.SpinnerColor)),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadDb, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.SetSize(msg.Width, msg.Height)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case dbLoadedMsg:
		if msg.err != nil {
			m.readyState = StatusError
			m.error = msg.err
			return m, tea.Quit
		}

		m.ctx.SetDb(msg.db)

		if m.ctx.Db.ShouldCheckForUpdate() {
			m.readyState = StatusUpdating
			cmds = append(cmds, m.updateDb)
		} else {
			cmds = append(cmds, m.displayApiList())
		}
	// keep using previous database till update fails
	case dbUpdatedMsg:
		cmds = append(cmds, m.displayApiList())
	}

	m.apilist, cmd = m.apilist.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) displayApiList() tea.Cmd {
	m.readyState = StatusOk
	return m.apilist.Init()
}

func (m Model) View() string {
	var s strings.Builder

	if m.readyState == StatusError {
		fmt.Fprintf(&s, "%s %s", theme.TextError.Render("error:"), m.renderError())
	} else if m.readyState == StatusInProgress || m.readyState == StatusUpdating {
		msg := "Loading database..."
		if m.readyState == StatusUpdating {
			msg = "Updating database, please wait..."
		}

		fmt.Fprintf(&s, "%s\n%s %s", m.apilist.Placeholder(), m.spinner.View(), msg)

		// fill remaining space so we take up all the screen
		s.WriteString(lipgloss.NewStyle().Height(m.ctx.Screen.Height - lipgloss.Height(s.String())).Render(""))
	} else {
		s.WriteString(m.apilist.View())
	}

	return s.String()
}

func (m Model) renderError() string {
	var e *net.DNSError
	if errors.As(m.error, &e) {
		return theme.TextError.Render("No database found locally and we were unable to download it. Are you connected to the internet? 🤔\n")
	}

	return m.error.Error()
}
