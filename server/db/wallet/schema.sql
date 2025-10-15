CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address CHAR(34) UNIQUE,
    public_key TEXT UNIQUE,
    public_key_hash VARCHAR(40) NOT NULL UNIQUE,
    balance NUMERIC(20, 8) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    create_at TIMESTAMP DEFAULT now(),
    last_login TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_wallets_address_pubkey ON wallets(address, public_key);
CREATE INDEX idx_wallets_address_pubkeyhash ON wallets(public_key_hash);

CREATE TABLE IF NOT EXISTS wallet_access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID REFERENCES wallets(id) ON DELETE CASCADE,
    access_time TIMESTAMP DEFAULT now(),
    ip TEXT,
    user_agent TEXT,
    access_type TEXT
);

CREATE INDEX idx_walletaccesslogs_walletid ON wallet_access_logs(wallet_id);
