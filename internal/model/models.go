// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package model

import (
	"time"
)

type Cycle struct {
	ID         int64
	CreatedAt  time.Time
	SessionID  int64
	Accomplish string
	Started    string
	Hazards    string
	Energy     int64
	Morale     int64
}

type Session struct {
	ID           int64
	CreatedAt    time.Time
	Accomplish   string
	Important    string
	Complete     string
	Distractions string
	Measurable   string
	Noteworthy   string
}
