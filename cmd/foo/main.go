package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/cmd/foo/planning"
	"github.com/rcy/megaworking/cmd/foo/preparation"
	"github.com/rcy/megaworking/cmd/foo/timerbar"
	"github.com/rcy/megaworking/internal/cycletimer"
	"github.com/rcy/megaworking/internal/db"
	_ "modernc.org/sqlite"
)

// initializing
// preparing
// cycling.planning
// cycling.working
// cycling.reviewing
// cycling.resting
// debriefing

type model struct {
	q           *db.Queries
	session     *db.Session
	cycles      []db.Cycle
	bar         timerbar.Model
	preparation preparation.Model
	planning    planning.Model
}

func New(q *db.Queries) model {
	return model{
		q:           q,
		bar:         timerbar.New(),
		preparation: preparation.New(q),
		planning:    planning.New(q),
	}
}

func (m model) fetchCurrentSession() tea.Msg {
	session, err := m.q.CurrentSession(context.TODO())
	if errors.Is(err, sql.ErrNoRows) {
		return messages.SessionNotFound{}
	}
	if err != nil {
		return messages.QueryError{Err: err}
	}

	cycles, err := m.q.SessionCycles(context.TODO(), session.ID)
	if err != nil {
		return messages.QueryError{Err: err}
	}

	return messages.SessionLoaded{Session: &session, Cycles: cycles}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchCurrentSession,
		m.preparation.Init(),
		m.bar.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case messages.QueryError:
		panic(msg.Err)

	case messages.SessionLoaded:
		m.session = msg.Session
		m.cycles = msg.Cycles

		cycleTimer := cycletimer.NewCustom(
			30*time.Minute,
			10*time.Minute,
			msg.Session.StartAt,
			msg.Session.StartAt,
			msg.Session.NumCycles,
		)
		cmds = append(cmds, func() tea.Msg {
			return messages.CycleTimerUpdated{CycleTimer: cycleTimer}
		})

	case messages.SessionCompleted:
		m.session = msg.Session
		return m, tea.Quit
	}

	model, cmd := m.preparation.Update(msg)
	m.preparation = model.(preparation.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.planning.Update(msg)
	m.planning = model.(planning.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.bar.Update(msg)
	m.bar = model.(timerbar.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	str := ""
	str += `MEGA ` + m.bar.View()
	str += m.preparation.View()
	str += m.planning.View()
	return str
}

func main() {
	file, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sqldb, err := sql.Open("sqlite", os.Getenv("SQLITE_DB"))
	if err != nil {
		panic(err)
	}

	queries := db.New(sqldb)
	program := tea.NewProgram(New(queries), tea.WithAltScreen())
	go func() {
		for {
			program.Send(messages.Tick{})
			time.Sleep(time.Second)
		}
	}()
	finalModel, err := program.Run()
	if err != nil {
		panic(err)
	}
	if m, ok := finalModel.(model); ok {
		if m.session != nil {
			fmt.Println("session.Status==", m.session.Status)
			switch m.session.Status {
			case "completed":
				fmt.Println("next steps:", m.session.Nextsteps)
			}
		}
	}

}
