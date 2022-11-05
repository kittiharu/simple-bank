-- name: CreateAccount :one
INSERT INTO account (
  owner, amount, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE ID = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM account
WHERE ID = $1 LIMIT 1
FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM account
ORDER BY ID 
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account
SET amount = $2
WHERE ID = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account 
WHERE ID = $1;

-- name: AddAccountAmount :one
UPDATE account
SET amount = amount + sqlc.arg(amount)
WHERE ID = sqlc.arg(id)
RETURNING *;