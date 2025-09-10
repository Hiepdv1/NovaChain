package types

type JWTWalletAuthPayload struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Pubkey  string `json:"pubkey"`
}
