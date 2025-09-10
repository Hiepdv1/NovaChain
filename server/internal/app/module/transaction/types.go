package transaction

import "time"

type TransactionDataParsed struct {
	Fee       float64
	Amount    float64
	To        []byte // base58
	Timestamp time.Time
	Message   string
}

type NewTransactionParsed struct {
	Data   TransactionDataParsed
	Sig    []byte
	PubKey []byte
}
