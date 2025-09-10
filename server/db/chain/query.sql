-- name: CreateBlock :one
insert into blocks (
    b_id, prev_hash, nonce, height,
    merkle_root, difficulty, tx_count, nchain_work, timestamp
) values (
    $1, $2, $3, $4,
    $5, $6, $7, $8,
    $9
) returning *;

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

-- name: DeleteBlockByBID :exec
delete from blocks
where b_id = $1;

-- name: CreateTransaction :one
insert into transactions (tx_id, b_id)
values ($1, $2) returning *;

-- name: GetCountTransaction :one
SELECT COUNT(*) FROM transactions;

-- name: GetTransactionByTxID :one
select * from transactions where tx_id = $1 limit 1;

-- name: GetListTransactionByBID :many
select * from transactions where b_id = $1;

-- name: GetListTransactions :many
select * from transactions offset $1 limit $2;

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
