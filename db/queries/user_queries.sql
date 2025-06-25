-- name: GetAllUsers :many
select *, count(*) over() as total_count from users
limit $1
offset $2;

-- name: GetUserByEmail :one
select * from users
where email = $1
and deleted_at is null
limit 1;

-- name: CreateUser :exec
insert into users (id, email, password, role, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6);

-- name: SetUserDeletedStatus :exec
update users
set deleted_at = $2
where id = $1;
