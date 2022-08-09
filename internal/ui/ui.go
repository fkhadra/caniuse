package ui

import (
	"caniuse/internal/constant"
	"caniuse/internal/ui/components/apilist"
	"caniuse/internal/ui/context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	error      error
	ctx        *context.AppContext
	readyState constant.Status
	apilist    apilist.Model
}

func New() Model {
	ctx := context.New()
	return Model{
		ctx:        ctx,
		apilist:    apilist.New(ctx),
		readyState: constant.StatusInProgress,
	}

}

func (m Model) Init() tea.Cmd {
	return loadDb
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case dbLoadedMsg:
		if msg.err != nil {
			m.readyState = constant.StatusError
			m.error = msg.err
			return m, tea.Quit
		}
		m.ctx.SetDb(msg.db)
		m.readyState = constant.StatusOk
		cmds = append(cmds, m.apilist.Init())
	case tea.WindowSizeMsg:
		m.ctx.SetSize(msg.Width, msg.Height)
	}

	m.apilist, cmd = m.apilist.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s strings.Builder

	if m.readyState == constant.StatusError {
		s.WriteString("Something is wrong")
	} else if m.readyState == constant.StatusInProgress {
		s.WriteString(m.apilist.Placeholder())
		// m.list.SetHeight(m.screen.height - lipgloss.Height(s.String()))
		// s.WriteString(m.list.View())
	} else {
		s.WriteString(m.apilist.View())
	}

	return s.String()
}
