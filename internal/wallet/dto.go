package wallet

type WalletResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	IsFrozen bool    `json:"is_frozen"`
}
