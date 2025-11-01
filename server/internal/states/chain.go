package states

type SyncState struct {
	SyncStatus string
}

var (
	ChainSyncState = SyncState{
		SyncStatus: "N/A",
	}
)
