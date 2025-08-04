package wallet

type Balance struct {
	Balance   float64 `json:"balance"`
	Address   string  `json:"address"`
	Timestamp int64   `json:"timestamp"`
	Error     *Error  `json:"error,omitempty"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
