-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE id IN ($1);

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: ListActiveUsers :many
SELECT * FROM users
WHERE active
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  id, name, birth, email, location, created_at, updated_at, active
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
  name = CASE WHEN @name_do_update::boolean
  THEN @name::text ELSE name END,

  birth = CASE WHEN @birth_do_update::boolean
  THEN @birth::date ELSE birth END,

  email = CASE WHEN @email_do_update::boolean
  THEN @email::text ELSE email END,

  location = CASE WHEN @location_do_update::boolean
  THEN @location::text ELSE location END,

  active = CASE WHEN @active_do_update::boolean
  THEN @active::boolean ELSE active END,

  updated_at = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
RETURNING *;
