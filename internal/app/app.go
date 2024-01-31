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
	q       *db.Queries
	session tea.Model
}

type NewSessionMsg struct {
	session tea.Model
}

func New(q *db.Queries) model {
	return model{
		q: q,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func NewSession(m model) func() tea.Msg {
	return func() tea.Msg {
		return NewSessionMsg{
			session: session.New(m.q, &db.Session{}),
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		// if m.session == nil {
		// 	if msg.Type == tea.KeyEnter {
		// 		return m, NewSession(m)
		// 	}
		// }

	case NewSessionMsg:
		m.session = msg.session
		return m, msg.session.Init()
	}

	if m.session != nil {
		newSession, newCmd := m.session.Update(msg)
		m.session = newSession
		return m, newCmd
	}

	return m, nil
}

func (m model) View() string {
	if m.session != nil {
		return m.session.View()
	}
	return "Welcome!"
}
