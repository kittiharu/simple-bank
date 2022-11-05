-- name: CreateEntry :one
INSERT INTO entry (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entry
WHERE ID = $1 LIMIT 1;
