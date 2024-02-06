package messages

import (
	"github.com/rcy/megaworking/internal/cycletimer"
	"github.com/rcy/megaworking/internal/db"
)

type QueryError struct {
	Err error
}
type SessionLoaded struct {
	Session *db.Session
	Cycles  []db.Cycle
}

type SessionCompleted struct {
	Session *db.Session
}

type SessionNotFound struct{}

type Tick struct{}

type CycleTimerUpdated struct {
	CycleTimer cycletimer.CycleTimer
}

type PhaseChanged struct {
	OldPhase cycletimer.Phase
	NewPhase cycletimer.Phase
}

type FinalCycleReviewCompleted struct{}
