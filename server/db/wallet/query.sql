-- name: CreateWallet :one
INSERT INTO wallets (
    address, public_key, public_key_hash, 
    balance, create_at, last_login
) VALUES (
    $1, $2, $3,
    $4, $5, $6
) RETURNING *;

-- name: UpdateWalletInfoByWalletID :one
UPDATE wallets
SET
    public_key = COALESCE(NULLIF(sqlc.narg('public_key')::text, ''), public_key),
    public_key_hash = COALESCE(NULLIF(sqlc.narg('public_key_hash')::text, ''), public_key_hash),
    balance = COALESCE(NULLIF(sqlc.narg('balance')::text, ''), balance),
    address = COALESCE(NULLIF(sqlc.narg('address')::text, ''), address),
    last_login = now()
WHERE id =  @id
RETURNING *;

-- name: ExistsWalletByAddrAndPubkey :one
SELECT EXISTS(SELECT 1 FROM wallets WHERE address = $1 AND public_key = $2);

-- name: GetWalletByAddress :one
SELECT * FROM wallets WHERE address = $1 LIMIT 1;

-- name: GetWalletByAddrAndPubkey :one
SELECT * FROM wallets WHERE address = $1 AND public_key = $2 LIMIT 1;

-- name: GetWalletByPubKeyHash :one
SELECT * FROM wallets WHERE public_key_hash = $1 LIMIT 1;

-- name: IncreaseWalletBalance :exec
UPDATE wallets
SET balance = balance + $1
WHERE address = $2 AND public_key = $3;

-- name: IncreaseWalletBalanceByPubKeyHash :exec
UPDATE wallets
SET balance = balance + $1
WHERE public_key_hash = $2;

-- name: DecreaseWalletBalance :exec
UPDATE wallets
SET balance = balance - $1
WHERE address = $2 AND public_key = $3 AND balance >= $1;

-- name: UpdateWalletLastLogin :exec
UPDATE wallets
SET last_login = now()
WHERE address = $1 AND public_key = $2;

-- name: CreateWalletAccessLog :exec
INSERT INTO wallet_access_logs (
    wallet_id, access_time, ip,
    user_agent, access_type
) VALUES (
    $1, $2, $3,
    $4, $5
);

-- name: GetListAccessLogByWalletID :many
SELECT * FROM wallet_access_logs WHERE wallet_id = $1 OFFSET $2 LIMIT $3;
