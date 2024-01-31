package planning

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/db"
)

type Model struct {
	q         *db.Queries
	numCycles int64
	cycles    []db.Cycle
}

func New(q *db.Queries) Model {
	return Model{
		q: q,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type cyclesLoaded struct {
	cycles []db.Cycle
}

func (m Model) fetchCycles(sessionID int64) func() tea.Msg {
	return func() tea.Msg {
		cycles, err := m.q.SessionCycles(context.TODO(), sessionID)
		if err != nil {
			return messages.QueryError{Err: err}
		}
		return cyclesLoaded{cycles: cycles}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case messages.SessionLoaded:
		m.numCycles = msg.Session.NumCycles
		cmds = append(cmds, m.fetchCycles(msg.Session.ID))

	case cyclesLoaded:
		m.cycles = msg.cycles
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%d %v load the form for the cycles", m.numCycles, m.cycles)
}
