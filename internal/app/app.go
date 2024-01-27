package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/rcy/megaworking/internal/db"
	"github.com/rcy/megaworking/internal/session"
)

type appState int

const (
	welcome appState = iota
	newSession
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
		m.state = newSession
		m.session = session.New(m.q)
		cmd := m.session.Init()
		return m, cmd
	}

	switch m.state {
	case newSession:
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
	case newSession:
		return "*** SESSION ***\n\n" + m.session.View()
	default:
		return "???"
	}
}

var prepParams = db.CreatePreparationParams{}

var sessionPrepForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("What am I trying to accomplish?").
			Value(&prepParams.Accomplish),
		huh.NewInput().
			Title("Why is this important and valuable?").
			Value(&prepParams.Important),
		huh.NewInput().
			Title("How will I know when this is complete?").
			Value(&prepParams.Complete),
		huh.NewInput().
			Title("Any risks / hazards? Potential distractions, procrastination, etc.").
			Value(&prepParams.Distractions),
		huh.NewInput().
			Title("Is this concrete / measurable or subjective / ambiguous?").
			Value(&prepParams.Measurable),
		huh.NewInput().
			Title("Anything else noteworthy?").
			Value(&prepParams.Noteworthy),
	),
)

var cycleParams = db.CreateCycleParams{}

var cyclePrepForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("What am I trying to accomplish this cycle?").
			Value(&cycleParams.Accomplish),
		huh.NewInput().
			Title("How will I get started?").
			Value(&cycleParams.Started),
		huh.NewInput().
			Title("Any hazards present?").
			Value(&cycleParams.Hazards),
		huh.NewSelect[int64]().
			Title("Energy").
			Options(
				huh.NewOption("High", int64(1)),
				huh.NewOption("Medium", int64(0)),
				huh.NewOption("Low", int64(-1)),
			).
			Value(&cycleParams.Energy),
		huh.NewSelect[int64]().
			Title("Morale").
			Options(
				huh.NewOption("High", int64(1)),
				huh.NewOption("Medium", int64(0)),
				huh.NewOption("Low", int64(-1)),
			).
			Value(&cycleParams.Morale),
	),
)
