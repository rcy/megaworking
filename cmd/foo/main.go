package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/cmd/foo/planning"
	"github.com/rcy/megaworking/cmd/foo/preparation"
	"github.com/rcy/megaworking/cmd/foo/timer"
	"github.com/rcy/megaworking/internal/db"
	_ "modernc.org/sqlite"
)

type model struct {
	q           *db.Queries
	model       tea.Model
	session     *db.Session
	timer       timer.Model
	preparation preparation.Model
	planning    planning.Model
}

func New(q *db.Queries) model {
	return model{
		q:           q,
		timer:       timer.New(),
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
		m.timer.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

	model, cmd = m.timer.Update(msg)
	m.timer = model.(timer.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	str := ""
	str += m.timer.View() + "\n"
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
	p := tea.NewProgram(m, tea.WithAltScreen())
	go func() {
		for {
			p.Send(messages.Tick{})
			time.Sleep(time.Second)
		}
	}()
	p.Run()
}
