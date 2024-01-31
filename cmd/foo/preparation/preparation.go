package preparation

import (
	"context"
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/rcy/megaworking/cmd/foo/create"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/db"
)

type Model struct {
	q        *db.Queries
	form     *huh.Form
	formData *db.Session
	create   create.Model
	state    string
}

func New(q *db.Queries) Model {
	return Model{
		q:      q,
		state:  "loading",
		create: create.New(q),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case messages.SessionNotFound:
		m.state = "notfound"

	case messages.SessionLoaded:
		m.state = msg.Session.State
		m.formData = msg.Session
		m.form = makeForm(m.formData)
		cmds = append(cmds, m.form.Init())
	}

	if m.state == "notfound" {
		model, cmd := m.create.Update(msg)
		m.create = model.(create.Model)
		cmds = append(cmds, cmd)
	}

	if m.state == "init" {
		if m.form != nil {
			model, cmd := m.form.Update(msg)
			m.form = model.(*huh.Form)
			cmds = append(cmds, cmd)
			if m.form.State == huh.StateCompleted {
				cmds = append(cmds, m.prepareSessionCmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) prepareSessionCmd() tea.Msg {
	s, err := m.q.PrepareSession(context.Background(), db.PrepareSessionParams{
		Accomplish:   m.formData.Accomplish,
		Important:    m.formData.Important,
		Complete:     m.formData.Complete,
		Distractions: m.formData.Distractions,
		Measurable:   m.formData.Measurable,
		Noteworthy:   m.formData.Noteworthy,
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}

	return messages.SessionLoaded{Session: &s}
}

func (m Model) View() string {
	switch m.state {
	case "notfound":
		return m.create.View()
	case "init":
		return m.form.View()
	case "prepared":
		return m.sessionPrepView()
	default:
		return "state==" + m.state
	}
}

func required(str string) error {
	if str == "" {
		return errors.New("This field is required")
	}
	return nil
}

func makeForm(data *db.Session) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What am I trying to accomplish?").
				Value(&data.Accomplish).
				Validate(required),
			huh.NewInput().
				Title("Why is this important and valuable?").
				Value(&data.Important),
			huh.NewInput().
				Title("How will I know when this is complete?").
				Value(&data.Complete),
			huh.NewInput().
				Title("Any risks / hazards? Potential distractions, procrastination, etc.").
				Value(&data.Distractions),
			huh.NewInput().
				Title("Is this concrete / measurable or subjective / ambiguous?").
				Value(&data.Measurable),
			huh.NewInput().
				Title("Anything else noteworthy?").
				Value(&data.Noteworthy),
		),
	)
}

func (m Model) sessionPrepView() string {
	d := m.formData
	strs := []string{
		lipgloss.NewStyle().Bold(true).Render("Session objective: " + d.Accomplish),
	}
	if d.Important != "" {
		strs = append(strs, "Why: "+d.Important)
	}
	if d.Complete != "" {
		strs = append(strs, "Completed: "+d.Complete)
	}
	if d.Distractions != "" {
		strs = append(strs, "Distractions: "+d.Distractions)
	}
	if d.Measurable != "" {
		strs = append(strs, "Measurable: "+d.Measurable)
	}
	if d.Noteworthy != "" {
		strs = append(strs, "Notes: "+d.Noteworthy)
	}

	str := strings.Join(strs, "\n")
	str += "\n\n"
	return lipgloss.NewStyle().SetString(str).Render()
}
