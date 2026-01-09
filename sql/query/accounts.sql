-- name: OpenAccount :one
INSERT INTO accounts (balance, type, fk_user_id)
VALUES (
    sqlc.arg(balance),
    sqlc.arg(type),
    sqlc.arg(fk_user_id)
)
RETURNING * ;

-- name: OpenSavingsAccount :one
INSERT INTO savings_accounts (account_id, interest_rate, min_balance, withdrawal_limit)
VALUES (
    sqlc.arg(account_id),
    sqlc.arg(interest_rate),
    sqlc.arg(min_balance),
    sqlc.arg(withdrawal_limit)
)
RETURNING * ;

-- name: OpenCheckingAccount :one
INSERT INTO checking_accounts (account_id, overdraft_limit, maintenance_fee)
VALUES (
    sqlc.arg(account_id),
    sqlc.arg(overdraft_limit),
    sqlc.arg(maintenance_fee)
)
RETURNING * ;

-- name: FindAccount :one
SELECT * FROM accounts
WHERE
    fk_user_id = $1 AND
    type = $2 ;

-- name: FindAccountByID :one
SELECT * FROM accounts
WHERE account_number = $1 ;

-- name: UpdateBalanceDeposit :execrows
UPDATE accounts
SET
    balance = balance + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE account_number = $1 ;

-- name: UpdateBalanceWithdraw :execrows
UPDATE accounts
SET
    balance = balance - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE account_number = $1 AND balance >= $2 ;

-- name: LockTransferAccount :execrows
SELECT * FROM accounts
WHERE account_number = $1
FOR UPDATE ;

-- name: FindAllAccountsByID :many
SELECT account_number, balance, type FROM
accounts INNER JOIN users
ON accounts.fk_user_id = users.id
WHERE id = $1 ;

