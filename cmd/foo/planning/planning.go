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
	state        string
	formData     *formData
	form         *huh.Form
}

type formData struct {
	db.Cycle
}

func New(q *db.Queries) Model {
	return Model{
		q:     q,
		state: "init",
		//formData: &formData{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type currentCycleNotFound struct{}

type currentCycleLoaded struct {
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
	return currentCycleLoaded{cycle: &cycle}
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
		m.formData = &formData{}
		m.form = makePlanForm(m.formData)
		cmds = append(cmds, m.form.Init())
		m.state = "plan"
	case currentCycleLoaded:
		m.currentCycle = msg.cycle
	case messages.PhaseChanged:
		switch m.cycleTimer.CurrentCycle().Phase {
		case cycletimer.Rest:
			// do review of previous cycle
			cmds = append(cmds, m.fetchCurrentCycle)
		case cycletimer.Work:
		}
	}

	if m.form != nil {
		model, cmd := m.form.Update(msg)
		m.form = model.(*huh.Form)
		cmds = append(cmds, cmd)

		log.Printf("m.form: %+v", m.form)

		if m.form.State == huh.StateCompleted {
			cmds = append(cmds, m.savePlanCmd)
			m.form = nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) savePlanCmd() tea.Msg {
	log.Print("planCycleCmd")
	cycle, err := m.q.CreateCycle(context.Background(), db.CreateCycleParams{
		SessionID:    m.session.ID,
		CycleTimerID: m.cycleTimer.CurrentCycle().ID,
		Accomplish:   m.formData.Accomplish,
		Started:      m.formData.Started,
		Hazards:      m.formData.Hazards,
		Energy:       m.formData.Energy,
		Morale:       m.formData.Morale,
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}

	return currentCycleLoaded{cycle: &cycle}
}

func (m Model) View() string {
	if m.cycleTimer == nil {
		return "no timer yet"
	}
	cyc := m.cycleTimer.CurrentCycle()

	log.Print("planning")
	log.Printf("  state:%s", m.state)
	log.Printf("  Phase:%s", cyc.Phase)
	log.Printf("  cyc.ID:%d", cyc.ID)
	if m.currentCycle != nil {
		log.Printf("  current:%d", m.currentCycle.CycleTimerID)
	} else {
		log.Print("  current: -")
	}

	if m.state == "plan" {
		if m.form != nil {
			return m.form.View()
		}
	}

	return ""
}

func makePlanForm(data *formData) *huh.Form {
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

func makeReviewForm(data *db.Cycle) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int64]().
				Title("Completed cycle's target?").
				Options(
					huh.NewOption("Yes", int64(100)),
					huh.NewOption("Half", int64(50)),
					huh.NewOption("No", int64(0)),
				),
			huh.NewInput().
				Title("Anything noteworthy?"),
			huh.NewInput().
				Title("Any distractions?"),
			huh.NewInput().
				Title("Things to improve for next cycle?"),
		))
}
