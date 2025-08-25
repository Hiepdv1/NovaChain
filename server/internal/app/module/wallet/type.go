package wallet

import "ChainServer/internal/common/apperror"

type Balance struct {
	Balance   float64 `json:"balance"`
	Address   string  `json:"address"`
	Timestamp int64   `json:"timestamp"`
	Error     *Error  `json:"error,omitempty"`
}

type Error struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type BalanceResult struct {
	Balance *Balance
	Err     *apperror.AppError
}

type JWTWalletAuthPayload struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Pubkey  string `json:"pubkey"`
}
