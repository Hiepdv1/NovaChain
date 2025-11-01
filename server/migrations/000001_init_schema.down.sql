
-- UTXOs
DROP INDEX IF EXISTS idx_utxos_block;
DROP INDEX IF EXISTS idx_utxos_pubkey;
DROP INDEX IF EXISTS idx_utxos_tx_out;
DROP TABLE IF EXISTS utxos CASCADE;

-- Wallet Access Logs
DROP INDEX IF EXISTS idx_walletaccesslogs_walletid;
DROP TABLE IF EXISTS wallet_access_logs CASCADE;

-- Wallets
DROP INDEX IF EXISTS idx_wallets_address_pubkeyhash;
DROP INDEX IF EXISTS idx_wallets_address_pubkey;
DROP TABLE IF EXISTS wallets CASCADE;

-- TX Outputs
DROP INDEX IF EXISTS idx_txoutputs_txid_index;
DROP INDEX IF EXISTS idx_txoutputs_pubkeyhash;
DROP INDEX IF EXISTS idx_txoutputs_block_id;
DROP TABLE IF EXISTS tx_outputs CASCADE;

-- TX Inputs
DROP INDEX IF EXISTS uniq_input_ref;
DROP INDEX IF EXISTS idx_txinputs_txid;
DROP INDEX IF EXISTS idx_txinputs_block_id;
DROP TABLE IF EXISTS tx_inputs CASCADE;

-- Transactions
DROP INDEX IF EXISTS idx_transactions_txid;
DROP INDEX IF EXISTS idx_transactions_block;
DROP TABLE IF EXISTS transactions CASCADE;

-- Blocks
DROP INDEX IF EXISTS idx_blocks_bid_prevhash;
DROP INDEX IF EXISTS idx_blocks_height;
DROP TABLE IF EXISTS blocks CASCADE;
