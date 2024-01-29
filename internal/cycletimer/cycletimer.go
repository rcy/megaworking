package cycles

import "time"

type CycleTimer struct {
	workDuration  time.Duration
	restDuration  time.Duration
	restEndAt     time.Duration
	workEndAt     time.Duration
	cycleDuration time.Duration
	origin        time.Time
}

type state uint8

const (
	Rest state = iota
	Work
)

type Cycle struct {
	ID        int64
	State     state
	Remaining time.Duration
}

var origin = time.Date(2022, time.June, 0, 1, 0, 0, 0, time.UTC)

func New() CycleTimer {
	return NewCustom(30, 10, origin)
}

func NewCustom(workMinutes, restMinutes int, origin time.Time) CycleTimer {
	workDuration := time.Minute * time.Duration(workMinutes)
	restDuration := time.Minute * time.Duration(restMinutes)

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
	var cycle Cycle

	elapsed := when.Sub(cs.origin)
	pos := elapsed % cs.cycleDuration
	cycle.ID = int64(elapsed / cs.cycleDuration)

	if pos < cs.restEndAt {
		cycle.State = Rest
		cycle.Remaining = cs.restEndAt - pos
	} else {
		cycle.State = Work
		cycle.Remaining = cs.workEndAt - pos
	}

	return cycle
}
