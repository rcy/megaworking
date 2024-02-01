package timer

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcy/megaworking/internal/cycletimer"
)

type Model struct {
	tick     int
	timer    cycletimer.CycleTimer
	progress progress.Model
}

func New() Model {
	return Model{
		timer: cycletimer.New(),
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

		// case messages.Tick:
		// 	m.tick++
	}
	return m, nil
}

func (m Model) View() string {
	cyc := m.timer.CurrentCycle()

	str := ""
	str += m.progress.ViewAs(cyc.PhasePercentComplete())
	str += " " + cyc.Phase.String()
	str += " " + cyc.PhaseRemaining.Round(time.Second).String()
	return str
}
