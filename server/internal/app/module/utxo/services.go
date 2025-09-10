package utxo

type UTXOService struct {
	rpcRepo RPCUtxoRepository
	dbRepo  DbUTXORepository
}

func NewUTXOService(rpcRepo RPCUtxoRepository, dbRepo DbUTXORepository) *UTXOService {
	return &UTXOService{
		rpcRepo: rpcRepo,
		dbRepo:  dbRepo,
	}
}
