// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"
)

const createCycle = `-- name: CreateCycle :one
insert into cycles(session_id, accomplish, started, hazards, energy, morale) values(?,?,?,?,?,?) returning id, created_at, session_id, accomplish, started, hazards, energy, morale
`

type CreateCycleParams struct {
	SessionID  int64
	Accomplish string
	Started    string
	Hazards    string
	Energy     int64
	Morale     int64
}

func (q *Queries) CreateCycle(ctx context.Context, arg CreateCycleParams) (Cycle, error) {
	row := q.db.QueryRowContext(ctx, createCycle,
		arg.SessionID,
		arg.Accomplish,
		arg.Started,
		arg.Hazards,
		arg.Energy,
		arg.Morale,
	)
	var i Cycle
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.SessionID,
		&i.Accomplish,
		&i.Started,
		&i.Hazards,
		&i.Energy,
		&i.Morale,
	)
	return i, err
}

const createPreparation = `-- name: CreatePreparation :one
insert into sessions(accomplish, important, complete, distractions, measurable, noteworthy) values(?,?,?,?,?,?) returning id, created_at, accomplish, important, complete, distractions, measurable, noteworthy
`

type CreatePreparationParams struct {
	Accomplish   string
	Important    string
	Complete     string
	Distractions string
	Measurable   string
	Noteworthy   string
}

func (q *Queries) CreatePreparation(ctx context.Context, arg CreatePreparationParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createPreparation,
		arg.Accomplish,
		arg.Important,
		arg.Complete,
		arg.Distractions,
		arg.Measurable,
		arg.Noteworthy,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Accomplish,
		&i.Important,
		&i.Complete,
		&i.Distractions,
		&i.Measurable,
		&i.Noteworthy,
	)
	return i, err
}