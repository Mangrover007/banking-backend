-- name: OpenAccount :one
INSERT INTO accounts (account_number, balance, type, fk_user_id)
VALUES (
    sqlc.arg(account_number),
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

-- name: UpdateBalanceDeposit :execrows
UPDATE accounts
SET
    balance = balance + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 ;

-- name: UpdateBalanceWithdraw :execrows
UPDATE accounts
SET
    balance = balance - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND balance >= $2 ;

-- name: LockTransferAccount :execrows
SELECT * FROM accounts
WHERE id = $1
FOR UPDATE ;

