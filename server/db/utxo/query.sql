-- name: CreateUTXO :one
INSERT INTO utxos (
    tx_id, output_index, value, pub_key_hash, block_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteUTXO :exec
DELETE FROM utxos
WHERE tx_id = $1 AND output_index = $2;

-- name: DeleteUTXOsByBlock :exec
DELETE FROM utxos
WHERE block_id = $1;

-- name: GetUTXOByID :one
SELECT * FROM utxos
WHERE id = $1 LIMIT 1;

-- name: GetUTXOByTxOut :one
SELECT * FROM utxos
WHERE tx_id = $1 AND output_index = $2 LIMIT 1;

-- name: GetUTXOsByTxID :many
SELECT * FROM utxos
WHERE tx_id = $1
ORDER BY output_index;

-- name: GetUTXOByPubKeyHash :one
SELECT * FROM utxos
WHERE pub_key_hash = $1
ORDER BY value DESC
LIMIT 1;

-- name: GetUTXOsByPubKeyHash :many
SELECT * FROM utxos
WHERE pub_key_hash = $1
ORDER BY value DESC;

