-- name: CreatePreparation :one
insert into preparations(accomplish, important, complete, distractions, measurable, noteworthy) values(?,?,?,?,?,?) returning *;
