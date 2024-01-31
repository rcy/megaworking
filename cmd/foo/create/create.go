package create

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/db"
)

func New(q *db.Queries) Model {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("startAt").
				Title("When do you want to start working").
				Options(
					huh.NewOption("Join the next group cycle", "group"),
					huh.NewOption("Join the currently in progress group cycle", "group-in-progress"),
					huh.NewOption("Start a new independent cycle right away", "now"),
				),
			huh.NewSelect[int]().
				Key("numCycles").
				Title("Number of cycles").
				Options(
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
					huh.NewOption("6", 6).Selected(true),
					huh.NewOption("7", 7),
					huh.NewOption("8", 8),
					huh.NewOption("9", 9),
					huh.NewOption("10", 10),
				),
		),
	)

	return Model{
		q:    q,
		form: form,
	}
}

type Model struct {
	q    *db.Queries
	form *huh.Form
}

func (m Model) createSessionCmd() tea.Msg {
	s, err := m.q.CreateSession(context.TODO(), db.CreateSessionParams{
		NumCycles: int64(m.form.GetInt("numCycles")),
		StartAt:   time.Now(),
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}
	return messages.SessionLoaded{Session: &s}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	model, cmd := m.form.Update(msg)
	m.form = model.(*huh.Form)
	cmds = append(cmds, cmd)
	if m.form.State == huh.StateCompleted {
		cmds = append(cmds, m.createSessionCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.form != nil {
		return m.form.View()
	}
	return "?form"
}
