package dto

type TxInput struct {
	ID        string `json:"ID"`
	Out       int64  `json:"Out"`
	Signature string `json:"Signature"`
	PubKey    string `json:"PubKey"`
}

type TxOutput struct {
	Value      float64 `json:"Value"`
	PubKeyHash string  `json:"PubKeyHash"`
}

type Transaction struct {
	ID      string     `json:"ID"`
	Inputs  []TxInput  `json:"Inputs"`
	Outputs []TxOutput `json:"Outputs"`
}
