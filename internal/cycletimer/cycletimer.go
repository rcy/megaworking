package cycletimer

import "time"

type CycleTimer struct {
	workDuration  time.Duration
	restDuration  time.Duration
	restEndAt     time.Duration
	workEndAt     time.Duration
	cycleDuration time.Duration
	origin        time.Time
}

type Phase uint8

func (p Phase) String() string {
	if p == Rest {
		return "Rest"
	}
	return "Work"
}

const (
	Void Phase = iota
	Rest
	Work
)

type Cycle struct {
	ID             int64
	Phase          Phase
	PhaseDuration  time.Duration
	PhaseRemaining time.Duration
}

func (c Cycle) PhasePercentComplete() float64 {
	return 1 - float64(c.PhaseRemaining)/float64(c.PhaseDuration)
}

var origin = time.Date(2022, time.June, 0, 1, 0, 0, 0, time.UTC)

func New() CycleTimer {
	return NewCustom(30*time.Minute, 10*time.Minute, origin)
}

func NewCustom(workDuration, restDuration time.Duration, origin time.Time) CycleTimer {
	return CycleTimer{
		origin:        origin,
		workDuration:  workDuration,
		restDuration:  restDuration,
		restEndAt:     restDuration,
		workEndAt:     restDuration + workDuration,
		cycleDuration: restDuration + workDuration,
	}
}

func (cs CycleTimer) CurrentCycle() Cycle {
	return cs.CycleAt(time.Now())
}

type Ticker struct {
	C       chan Cycle
	running bool
}

func (ti *Ticker) Stop() {
	ti.running = false
}

func (cs CycleTimer) NewTicker(interval time.Duration) *Ticker {
	t := &Ticker{
		C:       make(chan Cycle),
		running: true,
	}

	go func() {
		for t.running {
			t.C <- cs.CurrentCycle()
			time.Sleep(interval)
		}
	}()

	return t
}

func (cs CycleTimer) CycleAt(when time.Time) Cycle {
	cycle := Cycle{}

	if cs.cycleDuration == 0 {
		return cycle
	}

	elapsed := when.Sub(cs.origin)
	pos := elapsed % cs.cycleDuration
	cycle.ID = int64(elapsed / cs.cycleDuration)

	if pos < cs.restEndAt {
		cycle.Phase = Rest
		cycle.PhaseRemaining = cs.restEndAt - pos
		cycle.PhaseDuration = cs.restDuration
	} else {
		cycle.Phase = Work
		cycle.PhaseRemaining = cs.workEndAt - pos
		cycle.PhaseDuration = cs.workDuration
	}

	return cycle
}
