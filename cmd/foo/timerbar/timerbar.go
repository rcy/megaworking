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
		m.progress.Width = msg.Width - 15 // "Work 3/6 23m22s"

	case messages.CycleTimerUpdated:
		m.cycleTimer = &msg.CycleTimer
	}

	var cmds []tea.Cmd
	if m.cycleTimer != nil {
		phase := m.cycleTimer.CurrentCycle().Phase
		if phase != m.phase {
			cmds = append(cmds, phaseChangedCmd(m.phase, phase))
			m.phase = phase
		}
	}

	return m, tea.Batch(cmds...)
}

func phaseChangedCmd(old, new cycletimer.Phase) func() tea.Msg {
	return func() tea.Msg {
		return messages.PhaseChanged{
			OldPhase: old,
			NewPhase: new,
		}
	}
}

func (m Model) View() string {
	if m.cycleTimer == nil {
		return "\n"
	}

	if m.phase == cycletimer.Done {
		return "\n"
	}

	cyc := m.cycleTimer.CurrentCycle()

	cycleIndex := m.cycleTimer.CurrentCycle().ID - m.cycleTimer.FirstCycle().ID

	str := "\n"
	str += cyc.Phase.String() + fmt.Sprintf(" %d/%d %s ",
		cycleIndex+1,
		m.cycleTimer.NumCycles(),
		cyc.PhaseRemaining.Round(time.Second).String())
	str += m.progress.ViewAs(cyc.PhasePercentComplete())
	return str
}
