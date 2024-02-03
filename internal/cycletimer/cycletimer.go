package cycletimer

import (
	"log"
	"time"
)

type CycleTimer struct {
	workDuration  time.Duration
	restDuration  time.Duration
	restEndAt     time.Duration
	workEndAt     time.Duration
	cycleDuration time.Duration
	origin        time.Time
	startAt       time.Time
	numCycles     int64
}

type Phase uint8

func (p Phase) String() string {
	switch p {
	case Void:
		return "Void"
	case Rest:
		return "Rest"
	case Work:
		return "Work"
	case Done:
		return "Done"
	default:
		panic("unknown case")
	}
}

const (
	Void Phase = iota
	Rest
	Work
	Done
)

type Cycle struct {
	ID             int64
	Number         int64
	Phase          Phase
	PhaseDuration  time.Duration
	PhaseRemaining time.Duration
}

func (c Cycle) PhasePercentComplete() float64 {
	return 1 - float64(c.PhaseRemaining)/float64(c.PhaseDuration)
}

var origin = time.Date(2022, time.June, 0, 1, 0, 0, 0, time.UTC)

func New(startAt time.Time, numCycles int64) CycleTimer {
	panic("not implemented") // this breaks tests, work this out later
	return NewCustom(30*time.Minute, 10*time.Minute, origin, startAt, numCycles)
}

func NewCustom(workDuration, restDuration time.Duration, origin time.Time, startAt time.Time, numCycles int64) CycleTimer {
	return CycleTimer{
		origin:        origin,
		workDuration:  workDuration,
		restDuration:  restDuration,
		restEndAt:     restDuration,
		workEndAt:     restDuration + workDuration,
		cycleDuration: restDuration + workDuration,
		startAt:       startAt,
		numCycles:     numCycles,
	}
}

func (cs CycleTimer) CurrentCycle() Cycle {
	return cs.CycleAt(time.Now())
}

func (cs CycleTimer) FirstCycle() Cycle {
	return cs.CycleAt(cs.startAt)
}

func (cs CycleTimer) NumCycles() int64 {
	return cs.numCycles
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

func (cs CycleTimer) cycleIDAt(when time.Time) int64 {
	if cs.cycleDuration == 0 {
		return 0
	}

	elapsed := when.Sub(cs.origin)

	return int64(elapsed / cs.cycleDuration)
}

func (cs CycleTimer) CycleAt(when time.Time) Cycle {
	cycle := Cycle{}

	if cs.cycleDuration == 0 {
		return cycle
	}

	elapsed := when.Sub(cs.origin)
	pos := elapsed % cs.cycleDuration
	cycle.ID = cs.cycleIDAt(when)
	startingID := cs.cycleIDAt(cs.startAt)
	cycleNumber := cycle.ID - startingID + 1

	cycle.Number = cycleNumber
	log.Printf("cycleNumber %d\n", cycleNumber)

	if cycleNumber > cs.numCycles {
		cycle.Phase = Done
		cycle.Number = cs.numCycles
		return cycle
	}

	if pos < cs.restEndAt {
		cycle.Phase = Rest
		cycle.PhaseRemaining = cs.restEndAt - pos
		cycle.PhaseDuration = cs.restDuration
		return cycle
	}

	cycle.Phase = Work
	cycle.PhaseRemaining = cs.workEndAt - pos
	cycle.PhaseDuration = cs.workDuration
	return cycle
}
