package timerbar

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/rcy/megaworking/cmd/foo/messages"
	"github.com/rcy/megaworking/internal/cycletimer"
)

func TestOutput(t *testing.T) {
	m := New()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(300, 100),
	)
	// out, err := io.ReadAll(tm.Output())
	// if err != nil {
	// 	t.Error(err)
	// }
	//teatest.RequireEqualOutput(t, out)

	tm.Send(messages.CycleTimerUpdated{
		CycleTimer: cycletimer.New(time.Now(), 3),
	})

	tm.Send(tea.QuitMsg{})
	fm := tm.FinalModel(t, teatest.WithFinalTimeout(0))
	m, ok := fm.(Model)
	if !ok {
		t.Fatalf("wrong type %T", fm)
	}
	fmt.Printf("%+v", m.cycleTimer)
	if m.phase != cycletimer.Work {
		t.Errorf("wanted Done got %v", m.phase)
	}
}
