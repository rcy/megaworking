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

type SessionNotFound struct{}

type Tick struct{}

type CycleTimerUpdated struct {
	CycleTimer cycletimer.CycleTimer
}

type PhaseChanged struct{}
