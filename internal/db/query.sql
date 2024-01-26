-- name: CreatePreparation :one
insert into preparations(accomplish) values(?) returning *;
