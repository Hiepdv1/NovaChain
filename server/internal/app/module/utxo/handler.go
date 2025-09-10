package utxo

type UTXOHandler struct {
	service *UTXOService
}

func NewUTXOHandler(service *UTXOService) *UTXOHandler {
	return &UTXOHandler{
		service: service,
	}
}
