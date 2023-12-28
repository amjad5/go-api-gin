-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  name, phone_number
) VALUES (
  $1, $2
)
ON CONFLICT (phone_number) DO NOTHING
RETURNING *;

-- name: GetPhonenumber :one
SELECT phone_number from users
where phone_number = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  phone_number = $3
WHERE id = $1;

-- name: UpdateOtp :exec
UPDATE users
  set otp = $2,
  otp_expiration_time = $3
WHERE phone_number = $1;

-- name: VerifyOtp :one
SELECT phone_number, otp_expiration_time from users
where phone_number = $1 and otp = $2 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
