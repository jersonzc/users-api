-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE id = ANY($1::text[]);

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: ListActiveUsers :many
SELECT * FROM users
WHERE active
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  id, name, birth, email, location, active
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  name = CASE WHEN @name_do_update::boolean
  THEN @name ELSE name END,

  birth = CASE WHEN @birth_do_update::boolean
  THEN @birth ELSE birth END,

  email = CASE WHEN @email_do_update::boolean
  THEN @email ELSE email END,

  location = CASE WHEN @location_do_update::boolean
  THEN @location ELSE location END,

  active = CASE WHEN @active_do_update::boolean
  THEN @active ELSE active END
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
