CREATE TABLE IF NOT EXISTS pending_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_id CHAR(64) NOT NULL,
    address CHAR(34) NOT NULL REFERENCES wallets(address) ON DELETE CASCADE,
    receiver_address CHAR(34) NOT NULL,
    amount NUMERIC(20,8) NOT NULL,
    fee NUMERIC(20,8) NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    priority INT DEFAULT 0,
    message TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_pending_tx_address
    ON pending_transactions(address);
CREATE INDEX idx_pending_tx_receiver_address
    ON pending_transactions(address);
CREATE INDEX idx_pending_tx_status_priority 
    ON pending_transactions(status, priority DESC, created_at ASC);
CREATE INDEX idx_pending_txid 
    ON pending_transactions(tx_id);

CREATE TABLE IF NOT EXISTS pending_tx_data (
    pending_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_ref UUID NOT NULL REFERENCES pending_transactions(id) ON DELETE CASCADE,
    raw_tx JSONB NOT NULL,
    pub_key_hash CHAR(64) NOT NULL
);

CREATE INDEX idx_pending_tx_pubkeyhash 
    ON pending_tx_data(pub_key_hash);