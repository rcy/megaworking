package preparation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/rcy/megaworking/cmd/foo/create"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/db"
)

type state int

const (
	start state = iota
	loading
	processingPrep
	prepared
	processingDebrief
	debrief
	completed
)

type Model struct {
	q        *db.Queries
	form     *huh.Form
	formData *formData
	session  *db.Session
	create   create.Model
	state    state
}

type formData struct {
	db.Session
	startAtString   string
	numCyclesString string
}

func New(q *db.Queries) Model {
	return Model{
		q:        q,
		state:    loading,
		create:   create.New(q),
		formData: &formData{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	//log.Printf("Update %v", m.formData)

	switch msg := msg.(type) {
	case messages.SessionNotFound:
		log.Print("SessionNotFound")
		m.state = start
		m.session = &db.Session{}
		m.form = makePrepareForm(m.formData)
		cmds = append(cmds, m.form.Init())

	case messages.SessionLoaded:
		switch msg.Session.Status {
		case "prepared":
			m.state = prepared
		}
		m.session = msg.Session

	case messages.FinalCycleReviewCompleted:
		m.state = debrief
		m.form = makeDebriefForm(m.formData)
		cmds = append(cmds, m.form.Init())

	default:
		//log.Print("default: ", msg)
	}

	if m.form != nil {
		model, cmd := m.form.Update(msg)
		m.form = model.(*huh.Form)
		cmds = append(cmds, cmd)

		switch m.state {
		case start:
			if m.form.State == huh.StateCompleted {
				m.state = processingPrep
				cmds = append(cmds, withContext(context.TODO(), m.savePrepCmd))
			}

		case debrief:
			if m.form.State == huh.StateCompleted {
				m.state = processingDebrief
				cmds = append(cmds, withContext(context.TODO(), m.saveDebriefCmd))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func withContext(ctx context.Context, fn func(context.Context) tea.Msg) func() tea.Msg {
	return func() tea.Msg {
		return fn(ctx)
	}
}

func (m Model) savePrepCmd(ctx context.Context) tea.Msg {
	log.Print("savePrepCmd")

	numCycles, _ := strconv.Atoi(m.formData.numCyclesString)
	log.Printf("prepareSessionCmd %v", m.formData)
	s, err := m.q.CreateSession(ctx, db.CreateSessionParams{
		NumCycles:         int64(numCycles),
		StartAt:           time.Now(), // TODO
		StartCycleTimerID: 0,
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}

	s, err = m.q.PrepareSession(context.Background(), db.PrepareSessionParams{
		ID:           s.ID,
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

func (m Model) saveDebriefCmd(ctx context.Context) tea.Msg {
	log.Print("saveDebriefCmd")
	s, err := m.q.DebriefSession(ctx, db.DebriefSessionParams{
		ID:        m.session.ID,
		Target:    m.formData.Target,
		Done:      m.formData.Done,
		Compare:   m.formData.Compare,
		Bogged:    m.formData.Bogged,
		Replicate: m.formData.Replicate,
		Takeaways: m.formData.Takeaways,
		Nextsteps: m.formData.Nextsteps,
	})
	if err != nil {
		return messages.QueryError{Err: err}
	}

	return messages.SessionCompleted{Session: &s}
}

func (m Model) View() string {
	switch m.state {
	case start:
		return "Take a few minutes to prepare, so that your next 4 hours are effective.\n\n" +
			m.form.View()
	case processingPrep, prepared:
		return m.sessionPrepView()
	case debrief:
		return "Take a few minutes to debrief, so that you can identify and lock-in lessons.\n\n" +
			m.form.View()
	default:
		return fmt.Sprintf("state==%d", m.state)
	}
}

func required(str string) error {
	if str == "" {
		return errors.New("This field is required")
	}
	return nil
}

func makePrepareForm(data *formData) *huh.Form {
	log.Println("makeForm", data)
	if data.numCyclesString == "" {
		data.numCyclesString = "2"
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("startAt").
				Title("When do you want to start working").
				Options(
					huh.NewOption("Start a new cycle right away", "now"),
					huh.NewOption("Join the next group cycle", "group"),
				).
				Value(&data.startAtString),
			huh.NewInput().
				Key("numCycles").
				Title("Number of cycles").
				Value(&data.numCyclesString),
		),
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
				Title("Potential distractions, procrastination? How am I going to deal with them?").
				Value(&data.Distractions),
			huh.NewInput().
				Title("Is this concrete / measurable or subjective / ambiguous?").
				Value(&data.Measurable),
			huh.NewInput().
				Title("Anything else noteworthy?").
				Value(&data.Noteworthy),
		),
	).WithHeight(20)
}

func makeDebriefForm(data *formData) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int64]().
				Title("Did you complete your session's targets?").
				Options(
					huh.NewOption("Yes", int64(100)),
					huh.NewOption("Half", int64(50)),
					huh.NewOption("No", int64(0)),
				).
				Value(&data.Target),
			huh.NewText().
				Title("What did I get done this session?").
				Value(&data.Done),
			huh.NewText().
				Title("What are the next steps?").
				Value(&data.Nextsteps),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("How did this compare to my normal work output?").
				Value(&data.Compare),
			huh.NewInput().
				Title("Did I get bogged down? Where?").
				Value(&data.Bogged),
			huh.NewInput().
				Title("What went well? How can I replicate this in the future?").
				Value(&data.Replicate),
			huh.NewInput().
				Title("Any other takeaways? Lessons to share with others?").
				Value(&data.Takeaways),
		),
	)
}

func (m Model) sessionPrepView() string {
	d := m.session
	if d == nil {
		return ""
	}
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
}
