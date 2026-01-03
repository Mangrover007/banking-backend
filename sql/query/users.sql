-- name: RegisterUser :one
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

-- name: FindUserByPhoneOrEmail :one
SELECT * FROM users
WHERE phone_number = $1 OR email = $2;

-- name: CreateSession :one
INSERT INTO sessions (fk_user_id)
VALUES (
    fk_user_id = $1
)
RETURNING id;

-- name: DeleteSession :execrows
DELETE FROM sessions
WHERE id = $1;

