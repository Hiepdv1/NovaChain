-- name: CreateBlock :one
insert into blocks (
    b_id, prev_hash, nonce, height,
    merkle_root, nbits, tx_count, nchain_work, size, timestamp
) values (
    $1, $2, $3, $4,
    $5, $6, $7, $8,
    $9, $10
) returning *;

-- name: SearchExact :many
SELECT type, keyword, value
FROM (
  SELECT 
    'block' AS type, 
    blocks.b_id AS keyword, 
    blocks.height AS value
  FROM blocks
  WHERE blocks.b_id = $1

  UNION ALL

  SELECT 
    'transaction' AS type, 
    transactions.tx_id AS keyword, 
    transactions.amount AS value
  FROM transactions
  WHERE transactions.tx_id = $1 
     OR transactions.fromHash = $1 
     OR transactions.toHash = $1

  UNION ALL

  SELECT 
    'reward' AS type, 
    tx_outputs.pub_key_hash AS keyword, 
    tx_outputs.value AS value
  FROM tx_outputs
  WHERE tx_outputs.pub_key_hash = $1
) AS unified;

-- name: SearchFuzzy :many
SELECT 
  type::TEXT,
  keyword::TEXT,
  data::JSONB
FROM (
  SELECT 
    'block' AS type,
    b.b_id AS keyword,
    jsonb_build_object(
        'height', b.height,
        'timestamp', b.timestamp,
        'size', b.size,
        'tx_count', b.tx_count,
        'miner', COALESCE(coinbase.pub_key_hash, 'unknown')
    ) AS data,
    similarity(b.b_id::text, sqlc.arg('searchQuery')) AS score
  FROM blocks b
  LEFT JOIN LATERAL (
    SELECT o.pub_key_hash
    FROM tx_inputs i
    JOIN tx_outputs o ON o.b_id = i.b_id AND o.b_id = b.b_id
    WHERE i.out_index = -1 AND i.b_id = b.b_id
    LIMIT 1
  ) coinbase ON true
  WHERE similarity(b.b_id::text, sqlc.arg('searchQuery')) > 0

  UNION ALL

  SELECT 
    'transaction' AS type,
    t.tx_id AS keyword,
    jsonb_build_object(
        'from', t.fromHash,
        'to', t.toHash,
        'amount', t.amount,
        'fee', t.fee,
        'timestamp', t.create_at
    ) AS data,
    similarity(t.tx_id::text, sqlc.arg('searchQuery')) AS score  
  FROM transactions t
  JOIN tx_inputs i on i.tx_id = t.tx_id
  WHERE similarity(t.tx_id::text, sqlc.arg('searchQuery')) > 0 AND i.out_index > -1
) AS unified
ORDER BY score DESC
OFFSET $1
LIMIT $2;

-- name: GetBlockDetailWithTransactions :one
SELECT 
  b.*,
  (
    SELECT jsonb_agg(tx_row)
    FROM (
      SELECT jsonb_build_object(
        'ID', tx.id,
        'BID', tx.b_id,
        'TxID', tx.tx_id,
        'Fromhash', tx.fromHash,
        'Tohash', tx.toHash,
        'Amount', tx.amount,
        'Fee', tx.fee,
        'CreateAt', tx.create_at
      ) AS tx_row
      FROM transactions tx
      WHERE tx.b_id = b.b_id
      ORDER BY tx.create_at DESC
      OFFSET sqlc.arg('offsetTx')
      LIMIT sqlc.arg('limitTx')
    ) sub
  ) AS transactions,

  (
    SELECT COALESCE(SUM(tx.fee), 0)
    FROM transactions tx
    WHERE tx.b_id = b.b_id
  ) AS total_fee,

  (
    SELECT COALESCE(o.pub_key_hash, 'Unknown') 
    FROM tx_inputs i
    JOIN tx_outputs o 
      ON o.b_id = b.b_id 
     AND o.b_id = i.b_id
    WHERE i.out_index = -1
    LIMIT 1
  ) AS miner

FROM blocks b
WHERE b.b_id = $1
LIMIT 1;


-- name: CountFuzzy :one
SELECT COUNT(*) AS total_count
FROM (
  SELECT 1
  FROM blocks b
  WHERE similarity(b.b_id::TEXT, sqlc.arg('searchQuery')) > 0
  
  UNION ALL
  
  SELECT 1
  FROM transactions t 
  JOIN tx_inputs i on i.tx_id = t.tx_id
  WHERE similarity(t.tx_id::text, sqlc.arg('searchQuery')) > 0 AND i.out_index > -1
) AS unified;

-- name: CountFuzzyByType :many
SELECT type, COUNT(*)::BIGINT AS total
FROM (
  SELECT 'block' AS type
  FROM blocks
  WHERE blocks.b_id % $1

  UNION ALL

  SELECT 'transaction' AS type
  FROM transactions
  WHERE transactions.tx_id % $1 
     OR transactions.fromHash % $1 
     OR transactions.toHash % $1
) AS unified
GROUP BY type;

-- name: GetBlockByHeight :one
select * from blocks where height = $1 limit 1;

-- name: GetListBlocksByHeight :many
select * from blocks where height = $1;

-- name: GetBlockByBID :one
select * from blocks where b_id = $1 limit 1;

-- name: GetLastBlock :one
select * from blocks order by height desc limit 1;

-- name: GetListBlocks :many
select * from blocks order by height desc offset $1 limit $2;

-- name: IsExistingBlock :one
select exists (
    select 1 from blocks where b_id = $1
);

-- name: GetBestHeight :one
select height from blocks order by height desc limit 1;

-- name: GetBlockCountByHours :one
SELECT COUNT(*) 
FROM blocks
WHERE timestamp >= EXTRACT(EPOCH FROM NOW())::bigint - (sqlc.arg(hours)::bigint * 3600);

-- name: GetListBlockByHours :many
SELECT * FROM blocks 
WHERE timestamp >= EXTRACT(EPOCH FROM NOW())::bigint - (sqlc.arg(hours)::bigint * 3600);

-- name: DeleteBlockByBID :exec
delete from blocks
where b_id = $1;



-- name: CreateTransaction :one
insert into transactions (tx_id, b_id, fromHash, toHash, amount, fee, create_at)
values ($1, $2, $3, $4, $5, $6, $7) returning *;

-- name: GetCountTransaction :one
SELECT COUNT(*) FROM transactions;

-- name: SearchFuzzyTransactionsByBlock :many
SELECT 
  *
FROM transactions
WHERE 
  b_id = sqlc.arg('b_hash') AND
  similarity(tx_id::text, sqlc.arg('searchQuery')) > 0
ORDER BY 
  similarity(tx_id::text, sqlc.arg('searchQuery')) DESC
OFFSET $1
LIMIT $2;

-- name: CountFuzzyTransactionsByBlock :one
SELECT COUNT(*) AS total_count
FROM transactions
WHERE 
  b_id = sqlc.arg('b_hash') AND
  similarity(tx_id::text, sqlc.arg('searchQuery')) > 0;

-- name: GetTransactionByTxID :one
select * from transactions where tx_id = $1 limit 1;

-- name: GetListTransactionByBID :many
select * from transactions 
where b_id = $1
OFFSET $2
LIMIT $3;

-- name: GetFullTransactionByBID :many
select * FROM transactions where b_id = $1;

-- name: CountTransactionByBID :one
select COUNT(*) from transactions where b_id = $1;

-- name: GetListTransactions :many
select * from transactions
WHERE fromhash != '' AND tohash != ''
order by create_at desc offset $1 limit $2;

-- name: CountTransactions :one
select COUNT(*) from transactions;

-- name: CountTodayTransactions :one
select COUNT(*) from transactions
where create_at >= EXTRACT(EPOCH FROM date_trunc('day', now()))
and create_at < EXTRACT(EPOCH FROM date_trunc('day', now()) + INTERVAL '1 day');

-- name: GetListFullTransaction :many
SELECT tx.*, i.id as inID, i.input_tx_id, i.out_index, i.sig, i.pub_key, o.index, o.value, o.pub_key_hash
FROM transactions tx
JOIN tx_inputs i on i.tx_id = tx.tx_id
JOIN tx_outputs o on o.tx_id = tx.tx_id
where i.out_index > -1
order by create_at desc
offset $1 limit $2;



-- name: CreateTxInput :one
insert into tx_inputs (tx_id, input_tx_id, out_index, sig, b_id, pub_key)
values ($1, $2, $3, $4, $5, $6) returning *;

-- name: GetListTxInputByTxID :many
select * from tx_inputs where tx_id = $1;

-- name: GetTxInputByTxID :one
select * from tx_inputs where tx_id = $1 limit 1;

-- name: FindTxInputByBlockID :many
select * from tx_inputs where b_id = $1;



-- name: CreateTxOutput :one
insert into tx_outputs (tx_id, value, pub_key_hash, b_id, index)
values ($1, $2, $3, $4, $5) returning *;

-- name: FindListTxOutputByBlockID :many
select * from tx_outputs where b_id = $1;

-- name: GetListTxOutputByTxId :many
select * from tx_outputs where tx_id = $1;

-- name: GetTxOutputByTxID :one
select * from tx_outputs where tx_id = $1 limit 1;

-- name: GetTxOutputByTxIDAndIndex :one
select * from tx_outputs where tx_id = $1 and index = $2 limit 1;

-- name: CountDistinctMiners :one
SELECT COUNT(DISTINCT o.pub_key_hash) AS total_miners
FROM blocks b
JOIN transactions tx ON tx.b_id = b.b_id
JOIN tx_inputs i ON i.tx_id = tx.tx_id
JOIN tx_outputs o ON o.tx_id = tx.tx_id
WHERE i.out_index = -1
  AND b.height > 1
  AND o.pub_key_hash IS NOT NULL;

-- name: GetCountTodayWorkerMiners :one
SELECT COUNT(DISTINCT o.pub_key_hash) AS total_miners
FROM blocks b
JOIN transactions tx ON tx.b_id = b.b_id
JOIN tx_inputs i ON i.tx_id = tx.tx_id
JOIN tx_outputs o ON o.tx_id = tx.tx_id
WHERE i.out_index = -1
  AND b.height > 1
  AND o.pub_key_hash IS NOT NULL
  AND b.timestamp >= EXTRACT(EPOCH FROM date_trunc('day', now()))
  AND b.timestamp < EXTRACT(EPOCH FROM date_trunc('day', now()) + INTERVAL '1 day');