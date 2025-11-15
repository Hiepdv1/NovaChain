package p2p

const (
	MAX_HEADERS_PER_MSG = 500

	// Prefix Sync block
	PREFIX_BLOCK          = "block"
	PREFIX_BLOCK_SYNC     = "block_sync"
	PREFIX_HEADER         = "block_header"
	PREFIX_HEADER_SYNC    = "block_header_sync"
	PREFIX_HEADER_LOCATOR = "block_header_locator"
	PREFIX_TX_MINING      = "tx-mining"
	PREFIX_GET_DATA       = "getdata"
	PREFIX_GET_DATA_SYNC  = "get_data_sync"

	// Prefix Sync transaction
	PREFIX_TX_FROM_POOL = "gettxfrompool"
	PREFIX_TX           = "tx"
	PREFIX_INVENTORY    = "inv"
	PREFIX_DATA_TX      = "tx_Data"
	PREFIX_TXS_Data     = "txs"

	PREFIX_REQUEST_GOSSIP_PEERS = "request_gossip_peers"
	PREFOX_GOSSIP_PEERS         = "gossip_peers"
)
