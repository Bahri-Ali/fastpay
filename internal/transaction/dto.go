package transaction

// Request for initiating transfer
type TransferRequest struct {
    ReceiverID string  `json:"receiver_id" binding:"required"`
    Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// Response for client
type TransferResponse struct {
    Status           string  `json:"status"` // "completed" or "verification_required"
    TransactionID    string  `json:"transaction_id,omitempty"`
    VerificationToken string `json:"verification_token,omitempty"` // Used to verify PIN
    Message          string  `json:"message,omitempty"`
}

// Request for verifying PIN
type VerifyRequest struct {
    VerificationToken string `json:"verification_token" binding:"required"`
    PIN               string `json:"pin" binding:"required"`
}