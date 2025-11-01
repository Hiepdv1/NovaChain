-- Pending TX Data
DROP INDEX IF EXISTS idx_pending_tx_data_txid;
DROP INDEX IF EXISTS idx_pending_tx_pubkeyhash;
DROP TABLE IF EXISTS pending_tx_data CASCADE;

-- Pending Transactions
DROP INDEX IF EXISTS idx_pending_txid;
DROP INDEX IF EXISTS idx_pending_tx_status_priority;
DROP INDEX IF EXISTS idx_pending_tx_address;
DROP TABLE IF EXISTS pending_transactions CASCADE;

