package timerbar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/cycletimer"
)

type Model struct {
	tick       int
	cycleTimer *cycletimer.CycleTimer
	progress   progress.Model
	phase      cycletimer.Phase
}

func New() Model {
	return Model{
		progress: progress.New(
			progress.WithoutPercentage(),
			progress.WithWidth(80),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return m.progress.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 13

	case messages.CycleTimerUpdated:
		m.cycleTimer = &msg.CycleTimer
	}

	var cmds []tea.Cmd
	if m.cycleTimer != nil {
		if m.cycleTimer.CurrentCycle().Phase != m.phase {
			cmds = append(cmds, phaseChangedCmd)
		}
		m.phase = m.cycleTimer.CurrentCycle().Phase
	}

	return m, tea.Batch(cmds...)
}

func phaseChangedCmd() tea.Msg {
	return messages.PhaseChanged{}
}

func (m Model) View() string {
	if m.cycleTimer == nil {
		return "no timer yet"
	}
	cyc := m.cycleTimer.CurrentCycle()

	str := ""
	str += fmt.Sprint(cyc.ID) + "\n"
	str += m.progress.ViewAs(cyc.PhasePercentComplete())
	str += " " + cyc.Phase.String()
	str += " " + cyc.PhaseRemaining.Round(time.Second).String()
	return str
}
