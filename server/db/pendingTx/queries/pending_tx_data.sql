-- name: InsertPendingTxData :one
INSERT INTO pending_tx_data (raw_tx, pub_key_hash, tx_ref)
VALUES ($1, $2, $3)
RETURNING *;

-- name: SelectPendingTxByPubKeyHash :one
SELECT ptd.*
FROM pending_tx_data ptd
JOIN pending_transactions pt ON pt.tx_id = ptd.tx_id
WHERE ptd.pub_key_hash = $1;

-- name: SelectPendingDataTxByTxRef :one
SELECT * FROM pending_tx_data
WHERE tx_ref = $1 LIMIT 1;
