package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/rcy/megaworking/internal/cycletimer"
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
	error       error
	q           *db.Queries
	sessionData *db.Session
	cycleData   *db.Cycle
	form        *huh.Form
	spinner     spinner.Model
	timer       *cycletimer.CycleTimer
	cycle       cycletimer.Cycle
}

func New(q *db.Queries, sessionData *db.Session) model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot

	return model{
		state:       prepare,
		q:           q,
		sessionData: sessionData,
		//form:        makePrepareForm(sessionData),
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	//fmt.Println(reflect.TypeOf(msg).String())

	if s, ok := msg.(spinner.TickMsg); ok {
		spinner, cmd := m.spinner.Update(s)
		m.spinner = spinner
		cmds = append(cmds, cmd)
		// TODO return here rather than do more work below
	}

	if m.form != nil {
		newModel, cmd := m.form.Update(msg)
		if f, ok := newModel.(*huh.Form); ok {
			m.form = f
		}
		cmds = append(cmds, cmd)
	}

	if m.timer != nil {
		m.cycle = m.timer.CurrentCycle()
	}

	switch m.state {
	case prepare:
		if m.form != nil && m.form.State == huh.StateCompleted {
			m.state = plan
			//m.timer = cycletimer.NewCustom(10*time.Second, 5*time.Second, time.Now())
			session, err := m.q.PrepareSession(context.Background(), db.PrepareSessionParams{
				Accomplish:   m.sessionData.Accomplish,
				Important:    m.sessionData.Important,
				Complete:     m.sessionData.Complete,
				Distractions: m.sessionData.Distractions,
				Measurable:   m.sessionData.Measurable,
				Noteworthy:   m.sessionData.Noteworthy,
			})
			if err != nil {
				m.error = err
				break
			}
			m.sessionData = &session
			m.timer = cycletimer.New()
			m.cycleData = &db.Cycle{}
			m.form = makePlanForm(m.cycleData)
			cmds = append(cmds, m.form.Init())
		}
	case plan:
		if m.form != nil && m.form.State == huh.StateCompleted {
			m.form = nil
			if m.cycle.Phase == cycletimer.Work {
				m.state = work
			} else {
				m.state = rest
			}
		}
	case work:
		if m.cycle.Phase == cycletimer.Rest {
			m.state = review
			m.form = makeReviewForm(m.cycleData)
			cmds = append(cmds, m.form.Init())
		}
	case rest:
		if m.cycle.Phase == cycletimer.Work {
			m.state = work
		}
	case review:
		if m.form != nil && m.form.State == huh.StateCompleted {
			// TODO: check if completed the predetermined amount of cycles and route to either plan next cycle, or debrief session
			m.state = plan
			m.cycleData = &db.Cycle{}
			m.form = makePlanForm(m.cycleData)
			cmds = append(cmds, m.form.Init())
		}
	case debrief:
	default:
		panic("unhandled state")
	}

	return m, tea.Batch(cmds...)
}

var (
	sessionStyle = lipgloss.NewStyle()
)

func (m model) View() string {
	switch m.state {
	case prepare:
		if m.form != nil {
			return m.form.View()
		}
		return "..."

	case plan:
		str := m.sessionPrepView()
		str += "\n"
		str += m.form.View()
		return str
	case work:
		str := ""
		str += m.sessionPrepView()
		str += "\n\n"
		str += m.cyclePlanView()
		str += "\n"
		str += "\n"
		return str
	case rest:
		str := m.sessionPrepView()
		str += "\n"
		str += m.cyclePlanView()
		return str
	case review:
		return m.form.View()
	case debrief:
		return "DEBRIEF"
	default:
		return fmt.Sprintf("state=%d", m.state)
	}
}

func (m model) sessionPrepView0() string {
	str := "Session: " + m.sessionData.Accomplish + "\n"
	str += "Because: " + m.sessionData.Important + "\n"
	str += "Done when: " + m.sessionData.Complete + "\n"
	str += "Distractions: " + m.sessionData.Distractions + "\n"
	str += "Measurable: " + m.sessionData.Measurable + "\n"
	str += "Notes: " + m.sessionData.Noteworthy
	return sessionStyle.Foreground(lipgloss.Color("#bbbbbb")).Render(str)
}

func (m model) sessionPrepView() string {
	d := m.sessionData
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
	return lipgloss.NewStyle().SetString(str).Render()
	//return sessionStyle.Foreground(lipgloss.Color("#ffbbbb")).Render(str)
}

func (m model) cyclePlanView() string {
	str := "Cycle objective: " + m.cycleData.Accomplish + "\n"
	str += "First step: " + m.cycleData.Started + "\n"
	str += "Hazards: " + m.cycleData.Hazards + "\n"
	//str += "Energy: " + fmt.Sprint(m.cycleData.Energy) + " Morale: " + fmt.Sprint(m.cycleData.Morale)
	return sessionStyle.Foreground(lipgloss.Color("#ffffff")).Render(str)
}

func makePlanForm(data *db.Cycle) *huh.Form {
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
