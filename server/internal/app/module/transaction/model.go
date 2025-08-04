package transaction

type Transaction struct {
	ID      string
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID        string
	Out       int64
	Signature string
	PubKey    string
}

type TxOutput struct {
	Value      float64
	PubKeyHash string
}
