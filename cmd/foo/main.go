package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/cmd/foo/planning"
	"github.com/rcy/megaworking/cmd/foo/preparation"
	"github.com/rcy/megaworking/internal/cycletimer"
	"github.com/rcy/megaworking/internal/db"
	_ "modernc.org/sqlite"
)

type model struct {
	q           *db.Queries
	model       tea.Model
	session     *db.Session
	progress    progress.Model
	preparation preparation.Model
	planning    planning.Model
}

func New(q *db.Queries) model {
	return model{
		q: q,
		progress: progress.New(
			progress.WithoutPercentage(),
			progress.WithWidth(80),
		),
		preparation: preparation.New(q),
		planning:    planning.New(q),
	}
}

type newModelMsg struct {
	model tea.Model
}

func newModelCmd(model tea.Model) func() tea.Msg {
	return func() tea.Msg {
		return newModelMsg{model: model}
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

	return messages.SessionLoaded{Session: &session}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchCurrentSession,
		m.preparation.Init(),
		m.progress.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 13
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case messages.QueryError:
		panic(msg.Err)

	case messages.SessionLoaded:
		m.session = msg.Session
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	if m.model != nil {
		m.model, cmd = m.model.Update(msg)
		cmds = append(cmds, cmd)
	}

	model, cmd := m.preparation.Update(msg)
	m.preparation = model.(preparation.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.planning.Update(msg)
	m.planning = model.(planning.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	cyc := cycletimer.New().CurrentCycle()

	str := ""
	str += m.progress.ViewAs(cyc.PhasePercentComplete())
	str += " " + cyc.Phase.String()
	str += " " + cyc.PhaseRemaining.Round(time.Second).String()
	str += `

┏┳┓┏━╸┏━╸┏━┓╻ ╻┏━┓┏━┓╻┏
┃┃┃┣╸ ┃╺┓┣━┫┃╻┃┃ ┃┣┳┛┣┻┓
╹ ╹┗━╸┗━┛╹ ╹┗┻┛┗━┛╹┗╸╹ ╹


`
	str += m.preparation.View()
	str += "\n\n"

	str += m.planning.View()
	str += "\n\n"

	if m.model != nil {
		str += m.model.View()
	}
	return str
}

func main() {
	sqldb, err := sql.Open("sqlite", os.Getenv("SQLITE_DB"))
	if err != nil {
		panic(err)
	}

	q := db.New(sqldb)
	m := New(q)
	tea.NewProgram(m, tea.WithAltScreen()).Run()
}
