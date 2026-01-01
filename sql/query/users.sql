-- name: FindAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, phone_number, address, password)
VALUES (
    sqlc.arg(first_name),
    sqlc.arg(last_name),
    sqlc.arg(email),
    sqlc.arg(phone_number),
    sqlc.arg(address),
    sqlc.arg(password)
)
RETURNING *;

-- name: FindUserByPhone :one
SELECT * FROM users
WHERE phone_number = $1;

