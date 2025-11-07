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
  WHERE similarity(b.b_id::text, sqlc.arg('searchQuery')) > 0.1

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
  WHERE similarity(t.tx_id::text, sqlc.arg('searchQuery')) > 0.1 AND i.out_index > -1
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

-- name: GetRecentBlocksForNetworkInfo :many
SELECT height, nbits, timestamp
FROM blocks
ORDER BY height DESC
LIMIT $1;

-- name: CountFuzzy :one
SELECT COUNT(*) AS total_count
FROM (
  SELECT 1
  FROM blocks b
  WHERE similarity(b.b_id::TEXT, sqlc.arg('searchQuery')) > 0.1
  
  UNION ALL
  
  SELECT 1
  FROM transactions t 
  JOIN tx_inputs i on i.tx_id = t.tx_id
  WHERE similarity(t.tx_id::text, sqlc.arg('searchQuery')) > 0.1 AND i.out_index > -1
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
SELECT
  b.*,
  miner.*
FROM blocks b
LEFT JOIN LATERAL (
  SELECT o.pub_key_hash, o.value
  FROM tx_inputs i
  JOIN tx_outputs o ON i.tx_id = o.tx_id AND i.b_id = b.b_id
  WHERE i.out_index = -1
  LIMIT 1
) AS miner ON TRUE
ORDER BY b.height DESC
OFFSET $1 
LIMIT $2;

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

-- name: DeleteBlockByHeight :exec
DELETE FROM blocks
WHERE height = $1;

-- name: GetMiners :many
WITH miner_stats AS (
	SELECT 
		o.pub_key_hash AS miner_pubkey,                 
		MIN(b.timestamp) AS first_mined_at,             
	  	MAX(b.timestamp) AS last_mined_at,              
		COUNT(DISTINCT b.b_id) AS mined_blocks          
	FROM blocks b
	JOIN transactions tx ON tx.b_id = b.b_id
	JOIN tx_inputs i ON i.tx_id = tx.tx_id
	JOIN tx_outputs o ON o.tx_id = tx.tx_id
	WHERE i.out_index = -1
	  AND b.height > 1
	GROUP BY o.pub_key_hash
)
SELECT 
	miner_pubkey,                                         
	mined_blocks,                                         
	first_mined_at,                                       
	last_mined_at,                                        
	SUM(mined_blocks) OVER () AS total_blocks_network,    
	ROUND(
	    mined_blocks * 100.0 / SUM(mined_blocks) OVER (),
	    2
  	) AS network_share_percent                            
FROM miner_stats
ORDER BY mined_blocks DESC
OFFSET $1
LIMIT $2;

-- name: CountMiners :one
SELECT 
	COUNT(DISTINCT o.pub_key_hash)              
FROM blocks b
JOIN transactions tx ON tx.b_id = b.b_id
JOIN tx_inputs i ON i.tx_id = tx.tx_id
JOIN tx_outputs o ON o.tx_id = tx.tx_id
WHERE i.out_index = -1
  AND b.height > 1;



-- name: CreateTransaction :one
insert into transactions (tx_id, b_id, fromHash, toHash, amount, fee, create_at)
values ($1, $2, $3, $4, $5, $6, $7) returning *;

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

-- name: GetCountTransaction :one
SELECT COUNT(*) FROM transactions
WHERE fromhash != '' AND tohash != '';

-- name: CountTodayTransactions :one
select COUNT(*) from transactions
where create_at >= EXTRACT(EPOCH FROM date_trunc('day', now()))
and create_at < EXTRACT(EPOCH FROM date_trunc('day', now()) + INTERVAL '1 day')
and fromHash != '' and toHash != '';

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

-- name: GetTxSummaryByPubKeyHash :one
WITH all_activity AS (
    SELECT 
        fromhash     AS pub_key_hash,
        COUNT(*)     AS total_tx_sent,
        SUM(amount::numeric)  AS total_sent,
        0::numeric            AS total_received
    FROM transactions
    WHERE fromhash = sqlc.arg('pub_key_hash')::TEXT
    GROUP BY fromhash

    UNION ALL

    SELECT 
        tohash       AS pub_key_hash,
        0                 AS total_tx_sent,
        0::numeric        AS total_sent,
        SUM(amount::numeric) AS total_received
    FROM transactions
    WHERE tohash = sqlc.arg('pub_key_hash')::TEXT
    GROUP BY tohash
),
aggregated AS (
    SELECT
        pub_key_hash,
        SUM(total_tx_sent)     AS total_tx,
        SUM(total_sent)        AS total_sent,
        SUM(total_received)    AS total_received
    FROM all_activity
    GROUP BY pub_key_hash
)
SELECT 
    pub_key_hash,
    COALESCE(total_tx, 0)             AS total_tx,
    COALESCE(total_sent, 0)::TEXT   AS total_sent,
    COALESCE(total_received, 0)::TEXT AS total_received
FROM aggregated;

-- name: GetRecentTransaction :many
WITH recent AS (
  SELECT 
    'sent'     AS type,
    tx.id,
    tx.tx_id,
    tx.b_id,
    tx.create_at,
    tx.amount,
    tx.fee,
    tx.fromhash,
    tx.tohash
  FROM transactions tx
  WHERE tx.fromhash = sqlc.arg('pub_key_hash')::TEXT

  UNION ALL

  SELECT 
    'received' AS type,
    tx.id,
    tx.tx_id,
    tx.b_id,
    tx.create_at,
    tx.amount,
    tx.fee,
    tx.fromhash,
    tx.tohash  
  FROM transactions tx
  WHERE tx.tohash = sqlc.arg('pub_key_hash')::TEXT
)
SELECT type, id, tx_id, b_id, create_at, amount, fee, fromhash, tohash
FROM recent
ORDER BY create_at DESC
OFFSET $1
LIMIT $2;

-- name: CountRecentTransaction :one
SELECT COUNT(*) FROM transactions
WHERE fromhash = sqlc.arg('pub_key_hash')::TEXT OR
tohash = sqlc.arg('pub_key_hash')::TEXT;

-- name: GetDetailTx :one
WITH last_block AS (
  SELECT MAX(height) AS last_block FROM blocks
),
miner_info AS (
  SELECT 
    i.b_id,
    o.pub_key_hash AS miner
  FROM tx_inputs i
  JOIN tx_outputs o ON o.tx_id = i.tx_id
  WHERE i.out_index = -1
)
SELECT 
  tx.tx_id,
  b.height,
  b.b_id,
  b.timestamp,
  tx.fromhash,
  tx.tohash,
  tx.amount,
  tx.fee,
  b.nbits,
  b.nonce,
  m.miner,
  lb.last_block
FROM transactions tx
JOIN blocks b ON tx.b_id = b.b_id
LEFT JOIN miner_info m ON m.b_id = b.b_id
CROSS JOIN last_block lb
WHERE tx.tx_id = sqlc.arg('tx_hash')::TEXT
LIMIT 1;