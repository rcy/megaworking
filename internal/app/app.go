package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/internal/db"
	"github.com/rcy/megaworking/internal/session"
)

type appState int

const (
	welcome appState = iota
	inSession
)

type model struct {
	state   appState
	q       *db.Queries
	session tea.Model
}

type NewSessionMsg struct{}

func New(q *db.Queries) model {
	return model{
		state: welcome,
		q:     q,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func NewSession() tea.Msg {
	return NewSessionMsg{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if m.state == welcome {
			if msg.Type == tea.KeyEnter {
				return m, NewSession
			}
		}

	case NewSessionMsg:
		m.state = inSession
		m.session = session.New(m.q)
		cmd := m.session.Init()
		return m, cmd
	}

	switch m.state {
	case inSession:
		newSession, newCmd := m.session.Update(msg)
		m.session = newSession
		return m, newCmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case welcome:
		return "Welcome!"
	case inSession:
		return "*** SESSION ***\n\n" + m.session.View()
	default:
		return "???"
	}
}
