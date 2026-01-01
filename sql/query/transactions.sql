-- name: DepositTransaction :one
INSERT INTO transactions (type, status, amount, fk_sender)
VALUES (
    'DEPOSIT',
    sqlc.arg(status),
    sqlc.arg(amount),
    sqlc.arg(fk_sender)
)
RETURNING * ;

-- name: WithdrawTransaction :one
INSERT INTO transactions (type, status, amount, fk_sender)
VALUES (
    'WITHDRAW',
    sqlc.arg(status),
    sqlc.arg(amount),
    sqlc.arg(fk_sender)
)
RETURNING * ;

-- name: TransferTransaction :one
INSERT INTO transactions (type, status, amount, fk_sender, fk_recipient)
VALUES (
    'TRANSFER',
    sqlc.arg(status),
    sqlc.arg(amount),
    sqlc.arg(fk_sender),
    sqlc.arg(fk_recipient)
)
RETURNING * ;

-- name: UpdateTransactionStatus :execrows
UPDATE transactions
SET
    status = $2
WHERE id = $1 ;

