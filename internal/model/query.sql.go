// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package model

import (
	"context"
)

const createPreparation = `-- name: CreatePreparation :one
insert into preparations(accomplish, important, complete, distractions, measurable, noteworthy) values(?,?,?,?,?,?) returning created_at, accomplish, important, complete, distractions, measurable, noteworthy
`

type CreatePreparationParams struct {
	Accomplish   string
	Important    string
	Complete     string
	Distractions string
	Measurable   string
	Noteworthy   string
}

func (q *Queries) CreatePreparation(ctx context.Context, arg CreatePreparationParams) (Preparation, error) {
	row := q.db.QueryRowContext(ctx, createPreparation,
		arg.Accomplish,
		arg.Important,
		arg.Complete,
		arg.Distractions,
		arg.Measurable,
		arg.Noteworthy,
	)
	var i Preparation
	err := row.Scan(
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
