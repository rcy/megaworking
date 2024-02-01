package messages

import "github.com/rcy/megaworking/internal/db"

type QueryError struct {
	Err error
}
type SessionLoaded struct {
	Session *db.Session
}

type SessionNotFound struct{}

type Tick struct{}
