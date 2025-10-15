CREATE TABLE IF NOT EXISTS utxos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 

    tx_id CHAR(64) NOT NULL REFERENCES transactions(tx_id) ON DELETE CASCADE,  

    output_index BIGINT CHECK (output_index >= 0) NOT NULL,  

    value NUMERIC(20, 8) NOT NULL CHECK (value >= 0),  

    pub_key_hash CHAR(40) NOT NULL,                

    block_id CHAR(64) NOT NULL REFERENCES blocks(b_id) ON DELETE CASCADE,          

    CONSTRAINT uniq_utxo UNIQUE (tx_id, output_index)
);

CREATE INDEX idx_utxos_tx_out ON utxos(tx_id, output_index);

CREATE INDEX idx_utxos_pubkey ON utxos(pub_key_hash);

CREATE INDEX idx_utxos_block_id ON utxos(block_id);