package session

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/rcy/megaworking/internal/db"
)

type state int

const (
	prepare state = iota
	plan
	work
	review
	rest
	debrief
)

type model struct {
	state         state
	q             *db.Queries
	sessionParams *db.CreatePreparationParams
	form          *huh.Form
}

func New(q *db.Queries) model {
	sessionParams := &db.CreatePreparationParams{}

	return model{
		state:         prepare,
		q:             q,
		sessionParams: sessionParams,
		form:          makePrepareForm(sessionParams),
	}
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.form != nil {
		newModel, cmd := m.form.Update(msg)
		if f, ok := newModel.(*huh.Form); ok {
			m.form = f
		}
		cmds = append(cmds, cmd)
	}

	switch m.state {
	case prepare:
		if m.form != nil && m.form.State == huh.StateCompleted {
			m.state = plan
			m.form = makePlanForm()
			cmds = append(cmds, m.form.Init())
		}
	case plan:
		if m.form != nil && m.form.State == huh.StateCompleted {
			m.state = work
			m.form = makeReviewForm()
			cmds = append(cmds, m.form.Init())
		}
	case work:
	case review:
		if m.form != nil && m.form.State == huh.StateCompleted {
			m.state = work
			m.form = nil
		}
	case rest:
	case debrief:
	default:
		panic("unhandled state")
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case prepare:
		if m.form != nil {
			return m.form.View()
		}
		return "..."
	case plan:
		return m.form.View()
	case work:
		return "WORK"
	case review:
		return m.form.View()
	case rest:
		return "REST"
	case debrief:
		return "DEBRIEF"
	default:
		return fmt.Sprintf("state=%d", m.state)
	}
}

func makePrepareForm(values *db.CreatePreparationParams) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What am I trying to accomplish?").
				Value(&values.Accomplish),
			huh.NewInput().
				Title("Why is this important and valuable?").
				Value(&values.Important),
			huh.NewInput().
				Title("How will I know when this is complete?").
				Value(&values.Complete),
			huh.NewInput().
				Title("Any risks / hazards? Potential distractions, procrastination, etc.").
				Value(&values.Distractions),
			huh.NewInput().
				Title("Is this concrete / measurable or subjective / ambiguous?").
				Value(&values.Measurable),
			huh.NewInput().
				Title("Anything else noteworthy?").
				Value(&values.Noteworthy),
		),
	)
}

func makePlanForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What am I trying to accomplish this cycle?"),
			huh.NewInput().
				Title("How will I get started?"),
			huh.NewInput().
				Title("Any hazards present?"),
			huh.NewSelect[int64]().
				Title("Energy").
				Options(
					huh.NewOption("High", int64(1)),
					huh.NewOption("Medium", int64(0)),
					huh.NewOption("Low", int64(-1)),
				),
			huh.NewSelect[int64]().
				Title("Morale").
				Options(
					huh.NewOption("High", int64(1)),
					huh.NewOption("Medium", int64(0)),
					huh.NewOption("Low", int64(-1)),
				),
		),
	)
}

func makeReviewForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Completed cycle's target?"),
			huh.NewInput().
				Title("Anything noteworthy?"),
			huh.NewInput().
				Title("Any distractions?"),
			huh.NewInput().
				Title("Things to improve for next cycle?"),
		))
}
