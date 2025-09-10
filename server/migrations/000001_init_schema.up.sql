CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    b_id CHAR(64) NOT NULL UNIQUE,
    prev_hash CHAR(64),
    nonce BIGINT NOT NULL CHECK (nonce >= 0),
    height BIGINT NOT NULL CHECK (height > 0),
    merkle_root CHAR(64) NOT NULL,
    difficulty BIGINT NOT NULL CHECK (difficulty >= 0),
    tx_count BIGINT NOT NULL CHECK (tx_count >= 0),
    nchain_work TEXT NOT NULL,
    timestamp BIGINT NOT NULL CHECK (timestamp > 0)
);

CREATE INDEX idx_blocks_height ON blocks(height DESC);
CREATE INDEX idx_blocks_bid_prevhash ON blocks(b_id, prev_hash);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_id CHAR(64) NOT NULL UNIQUE,
    b_id CHAR(64) NOT NULL REFERENCES blocks(b_id) ON DELETE CASCADE,
    create_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_transactions_block ON transactions(b_id);
CREATE INDEX idx_transactions_txid ON transactions(tx_id);

CREATE TABLE IF NOT EXISTS tx_inputs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_id CHAR(64) NOT NULL REFERENCES transactions(tx_id) ON DELETE CASCADE,
    input_tx_id CHAR(64) REFERENCES transactions(tx_id) ON DELETE CASCADE,
    out_index BIGINT CHECK (out_index >= -1) NOT NULL,
    sig TEXT,
    pub_key TEXT,
    b_id CHAR(64) NOT NULL REFERENCES blocks(b_id) ON DELETE CASCADE
);

CREATE INDEX idx_txinputs_block_id ON tx_inputs(b_id);
CREATE INDEX idx_txinputs_txid ON tx_inputs(tx_id);
CREATE UNIQUE INDEX uniq_input_ref ON tx_inputs(input_tx_id, out_index);

CREATE TABLE IF NOT EXISTS tx_outputs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_id CHAR(64) NOT NULL REFERENCES transactions(tx_id) ON DELETE CASCADE,
    index BIGINT NOT NULL CHECK (index >= -1),
    value NUMERIC(20, 8) NOT NULL CHECK (value >= 0),
    pub_key_hash CHAR(64) NOT NULL,
    b_id CHAR(64) NOT NULL REFERENCES blocks(b_id) ON DELETE CASCADE

);

CREATE INDEX idx_txoutputs_block_id ON tx_outputs(b_id);
CREATE INDEX idx_txoutputs_pubkeyhash ON tx_outputs(pub_key_hash);
CREATE INDEX idx_txoutputs_txid_index ON tx_outputs(tx_id, index);

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address CHAR(64) NOT NULL UNIQUE,
    public_key TEXT NOT NULL UNIQUE,
    public_key_hash TEXT NOT NULL UNIQUE,
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

CREATE TABLE IF NOT EXISTS utxos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 

    tx_id CHAR(64) REFERENCES transactions(tx_id) ON DELETE CASCADE,  

    output_index BIGINT CHECK (output_index >= -1) NOT NULL,  

    value NUMERIC(20, 8) NOT NULL CHECK (value >= 0),  

    pub_key_hash CHAR(64) NOT NULL,                

    block_id CHAR(64) NOT NULL REFERENCES blocks(b_id) ON DELETE CASCADE,          

    CONSTRAINT uniq_utxo UNIQUE (tx_id, output_index)
);

CREATE INDEX idx_utxos_tx_out ON utxos(tx_id, output_index);

CREATE INDEX idx_utxos_pubkey ON utxos(pub_key_hash);

CREATE INDEX idx_utxos_block ON utxos(block_id);