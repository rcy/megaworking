-- name: Sessions :many
select * from sessions order by created_at desc;

-- name: CurrentSession :one
select * from sessions order by created_at desc limit 1;

-- name: PrepareSession :one
update sessions
set accomplish = ?,
    important = ?,
    complete = ?,
    distractions = ?,
    measurable = ?,
    noteworthy = ?,
    state = 'prepared'
where id = ?
returning *;

-- name: CreateSession :one
insert into sessions(
       num_cycles,
       start_at,
       start_cycle_timer_id
) values(?, ?, ?) returning *;

-- name: CreateCycle :one
insert into cycles(session_id, cycle_timer_id, accomplish, started, hazards, energy, morale) values(?,?,?,?,?,?,?) returning *;

-- name: SessionCycles :many
select * from cycles where session_id = ? order by cycle_timer_id;

-- name: SessionCycleByCycleTimerID :one
select * from cycles where session_id = ? and cycle_timer_id = ?;
