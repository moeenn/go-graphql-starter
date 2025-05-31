-- name: GetUsers :many
select * from user
limit ?
offset ?;

-- name: GetUserByEmail :one
select * from user
where email = ?
and deleted_at is null
limit 1;

-- name: CreateUser :exec
insert into user (id, email, password, role, created_at, updated_at)
values (?, ?, ?, ?, ?, ?);

-- name: SetUserDeletedStatus :exec
update user
set deleted_at = ?
where id = ?;
