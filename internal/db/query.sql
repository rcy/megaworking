-- name: CreateSession :one
insert into sessions(accomplish, important, complete, distractions, measurable, noteworthy) values(?,?,?,?,?,?) returning *;

-- name: CreateCycle :one
insert into cycles(session_id, accomplish, started, hazards, energy, morale) values(?,?,?,?,?,?) returning *;
