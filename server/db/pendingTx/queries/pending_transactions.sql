-- name: InsertPendingTransaction :one
INSERT INTO pending_transactions (tx_id, address, receiver_address, status, priority, message, amount, fee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdatePendingTransactionStatus :one
UPDATE pending_transactions
SET status = $2, updated_at = now()
WHERE tx_id = $1
RETURNING *;

-- name: UpdatePendingTransactionPriority :one
UPDATE pending_transactions
SET priority = $2, updated_at = now()
WHERE tx_id = $1
RETURNING *;

-- name: GetListPendingTxs :many
SELECT 
  id,
  tx_id,
  address,
  receiver_address,
  status,
  amount,
  fee,
  priority,
  created_at,
  updated_at
FROM pending_transactions
ORDER BY created_at DESC
OFFSET $1
LIMIT $2;

-- name: GetPendingTxsByStatus :many
SELECT 
  id,
  tx_id,
  address,
  receiver_address,
  status,
  amount,
  fee,
  priority,
  created_at,
  updated_at
FROM pending_transactions
WHERE status = ANY(sqlc.arg(statuses)::text[])
ORDER BY created_at DESC
OFFSET $1
LIMIT $2;

-- name: GetCountPendingTxs :one
SELECT COUNT(*) FROM pending_transactions;

-- name: GetCountPendingTxsByStatus :one
SELECT COUNT(*) FROM pending_transactions
WHERE status = ANY(sqlc.arg(statuses)::text[]);

-- name: SelectPendingTransactions :many
SELECT p.*, pd.raw_tx, pd.pub_key_hash
FROM pending_transactions p
JOIN pending_tx_data pd ON p.id = pd.tx_ref
WHERE p.status = $1
ORDER BY p.fee DESC, p.created_at ASC
LIMIT $2 OFFSET $3;

-- name: UpdatePendingTxsStatus :execrows
UPDATE pending_transactions
SET status = sqlc.arg(new_status),
    updated_at = NOW()
WHERE tx_id = ANY(sqlc.arg(tx_ids)::text[])
  AND status = ANY(sqlc.arg(old_status)::text[]);

-- name: SelectTxPendingByTxID :one
SELECT * FROM pending_transactions
WHERE tx_id = $1 LIMIT 1;
  
-- name: PendingTxsByAddress :many
SELECT p.*, pd.raw_tx, pd.pub_key_hash
FROM pending_transactions p
JOIN pending_tx_data pd ON p.id = pd.tx_ref
WHERE p.address = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: PendingTxsByAddressAndStatus :many
SELECT p.*, pd.raw_tx, pd.pub_key_hash
FROM pending_transactions p
JOIN pending_tx_data pd ON p.id = pd.tx_ref
WHERE p.address = $1
  AND status = ANY(sqlc.narg('status')::text[])
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPendingTransactionsByAddr :one
SELECT COUNT(*) FROM pending_transactions
WHERE address = $1;

-- name: PendingTxExists :one
SELECT EXISTS (
    SELECT 1
    FROM pending_transactions
    WHERE tx_id = $1
      AND status = ANY(sqlc.narg('status')::text[])
);

-- name: CountPendingTxs :one
select COUNT(*) FROM pending_transactions
where status in ('mining', 'pending');

-- name: CountTodayPendingTxs :one
select COUNT(*) FROM pending_transactions
where created_at >= date_trunc('day', now())
and created_at < date_trunc('day', now()) + INTERVAL '1 day' AND status in ('mining', 'pending');
