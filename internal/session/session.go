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
	state       state
	q           *db.Queries
	prepareForm huh.Form
	planForm    huh.Form
	reviewForm  huh.Form
}

func New(q *db.Queries) model {
	return model{
		state:       prepare,
		q:           q,
		prepareForm: *makePrepareForm(),
		planForm:    *makePlanForm(),
		reviewForm:  *makeReviewForm(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.prepareForm.Init(),
		m.planForm.Init(),
		m.reviewForm.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case prepare:
		newModel, cmd := m.prepareForm.Update(msg)
		if f, ok := newModel.(*huh.Form); ok {
			m.prepareForm = *f
		}

		if m.prepareForm.State == huh.StateCompleted {
			m.state = plan
			m.planForm = *makePlanForm()
			m.planForm.Init()
		}

		return m, cmd
	case plan:
		newModel, cmd := m.planForm.Update(msg)
		if f, ok := newModel.(*huh.Form); ok {
			m.planForm = *f
		}

		if m.planForm.State == huh.StateCompleted {
			m.state = prepare
			m.prepareForm = *makePrepareForm()
			m.prepareForm.Init()
		}

		return m, cmd
	case review:
		newModel, cmd := m.reviewForm.Update(msg)
		if f, ok := newModel.(*huh.Form); ok {
			m.reviewForm = *f
		}

		if m.reviewForm.State == huh.StateCompleted {
			m.state = work
		}

		return m, cmd
	case work:
		return m, nil
	default:
		panic("unhandled state")
	}
}

func (m model) View() string {
	switch m.state {
	case prepare:
		return m.prepareForm.View()
	case plan:
		return m.planForm.View()
	case review:
		return m.reviewForm.View()
	default:
		return fmt.Sprintf("state=%d", m.state)
	}
}

var prepParams = db.CreatePreparationParams{}

func makePrepareForm() *huh.Form {
	return huh.NewForm(
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
