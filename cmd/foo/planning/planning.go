package planning

import (
	"context"
	"database/sql"
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/cycletimer"
	"github.com/rcy/megaworking/internal/db"
)

type Model struct {
	q            *db.Queries
	session      *db.Session
	cycles       []db.Cycle
	currentCycle *db.Cycle
	cycleTimer   *cycletimer.CycleTimer
	formState    formState
	planData     *planData
	reviewData   *reviewData
	form         *huh.Form
	//phase      cycletimer.Phase
}

type formState int

const (
	none formState = iota
	plan
	review
)

type planData struct {
	db.Cycle
}

type reviewData struct {
	db.Cycle
}

func New(q *db.Queries) Model {
	return Model{
		q:         q,
		formState: none,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type currentCycleNotFound struct{}

type cycleLoaded struct {
	cycle *db.Cycle
}

func (m Model) fetchCurrentCycle() tea.Msg {
	log.Printf("fetchCurrentCycle %+v", m)

	if m.cycleTimer == nil {
		return nil
	}
	cycle, err := m.q.SessionCycleByCycleTimerID(context.TODO(), db.SessionCycleByCycleTimerIDParams{
		SessionID:    m.session.ID,
		CycleTimerID: m.cycleTimer.CurrentCycle().ID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return currentCycleNotFound{}
	}
	if err != nil {
		return messages.QueryError{Err: err}
	}
	return cycleLoaded{cycle: &cycle}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case messages.SessionLoaded:
		m.session = msg.Session
		m.cycles = msg.Cycles
		cmds = append(cmds, m.fetchCurrentCycle)
	case messages.CycleTimerUpdated:
		m.cycleTimer = &msg.CycleTimer
		cmds = append(cmds, m.fetchCurrentCycle)
	case currentCycleNotFound:
		if m.cycleTimer.CurrentCycle().Phase != cycletimer.Done {
			m.formState = plan
			m.planData = &planData{}
			m.form = makePlanForm(m.planData)
			cmds = append(cmds, m.form.Init())
		} else {
			cmds = append(cmds, m.finalCycleReviewCompleted)
		}
	case cycleLoaded:
		m.currentCycle = msg.cycle
		if m.formState == review {
			// after reviewing the last cycle, plan the current cycle
			cmds = append(cmds, m.fetchCurrentCycle)
		}
	case messages.PhaseChanged:
		log.Printf("planning: PhaseChanged: old=%s new=%s", msg.OldPhase, msg.NewPhase)
		if msg.OldPhase == cycletimer.Work {
			// we are coming out of a work phase, we need to review
			m.formState = review
			m.reviewData = &reviewData{}
			m.form = makeReviewForm(m.reviewData)
			cmds = append(cmds, m.form.Init())
		} else if msg.OldPhase == cycletimer.Void {
			// we are entering the first cycle
			cmds = append(cmds, m.fetchCurrentCycle)
			// m.formState = plan
			// m.planData = &planData{}
			// m.form = makePlanForm(m.planData)
			// cmds = append(cmds, m.form.Init())
		}
	}

	if m.form != nil {
		model, cmd := m.form.Update(msg)
		m.form = model.(*huh.Form)
		cmds = append(cmds, cmd)

		if m.form.State == huh.StateCompleted {
			m.form = nil
			if m.formState == plan {
				cmds = append(cmds, m.savePlanCmd)
				m.formState = none
			}
			if m.formState == review {
				cmds = append(cmds, m.saveReviewCmd)
				// if m.cycleTimer.CurrentCycle().Number <= m.cycleTimer.NumCycles() {
				// 	m.state = plan
				// 	m.planData = &planData{}
				// 	m.form = makePlanForm(m.planData)
				// 	cmds = append(cmds, m.form.Init())
				// } else {
				// 	m.state = idle
				// }
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) finalCycleReviewCompleted() tea.Msg {
	return messages.FinalCycleReviewCompleted{}
}

type planMsg struct{}

func (m Model) planCmd() tea.Msg {
	return planMsg{}
}

func (m Model) savePlanCmd() tea.Msg {
	log.Print("planCycleCmd")
	cycle, err := m.q.CreateCycle(context.Background(), db.CreateCycleParams{
		SessionID:    m.session.ID,
		CycleTimerID: m.cycleTimer.CurrentCycle().ID,
		Accomplish:   m.planData.Accomplish,
		Started:      m.planData.Started,
		Hazards:      m.planData.Hazards,
		Energy:       m.planData.Energy,
		Morale:       m.planData.Morale,
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}

	return cycleLoaded{cycle: &cycle}
}

func (m Model) saveReviewCmd() tea.Msg {
	cycle, err := m.q.UpdateCycle(context.Background(), db.UpdateCycleParams{
		ID:           m.cycleTimer.CurrentCycle().ID,
		Target:       m.reviewData.Target,
		Noteworthy:   m.reviewData.Noteworthy,
		Distractions: m.reviewData.Distractions,
		Improve:      m.reviewData.Improve,
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}

	return cycleLoaded{cycle: &cycle}
}

func (m Model) View() string {
	if m.form != nil {
		return "\n" + m.form.View()
	} else {
		return "\n" + m.CyclePlanView()
	}
}

func (m Model) CyclePlanView() string {
	if m.currentCycle != nil {
		return m.currentCycle.Accomplish
	}
	return ""
}

func makePlanForm(data *planData) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What am I trying to accomplish this cycle?").
				Value(&data.Accomplish),
			huh.NewInput().
				Title("How will I get started?").
				Value(&data.Started),
			huh.NewInput().
				Title("Any hazards present?").
				Value(&data.Hazards),
		),
		huh.NewGroup(
			huh.NewSelect[int64]().
				Title("Energy").
				Options(
					huh.NewOption("High", int64(1)),
					huh.NewOption("Medium", int64(0)),
					huh.NewOption("Low", int64(-1)),
				).
				Value(&data.Energy),
			huh.NewSelect[int64]().
				Title("Morale").
				Options(
					huh.NewOption("High", int64(1)),
					huh.NewOption("Medium", int64(0)),
					huh.NewOption("Low", int64(-1)),
				).
				Value(&data.Morale),
		),
	)
}

func makeReviewForm(data *reviewData) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int64]().
				Title("Completed cycle's target?").
				Options(
					huh.NewOption("Yes", int64(100)),
					huh.NewOption("Half", int64(50)),
					huh.NewOption("No", int64(0)),
				).
				Value(&data.Target),
			huh.NewInput().
				Title("Anything noteworthy?").
				Value(&data.Noteworthy),
			huh.NewInput().
				Title("Any distractions?").
				Value(&data.Distractions),
			huh.NewInput().
				Title("Things to improve for next cycle?").
				Value(&data.Improve),
		))
}
