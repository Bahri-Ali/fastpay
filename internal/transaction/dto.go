package transaction

import "time"


type TransferRequest struct {
    ReceiverID string  `json:"receiver_id" binding:"required"`
    Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// Response for client
type TransferResponse struct {
    Status           string  `json:"status"` 
    TransactionID    string  `json:"transaction_id,omitempty"`
    VerificationToken string `json:"verification_token,omitempty"` 
    Message          string  `json:"message,omitempty"`
}


type VerifyRequest struct {
    VerificationToken string `json:"verification_token" binding:"required"`
    PIN               string `json:"pin" binding:"required"`
}

type TransactionItemDTO struct {
    ID           string    `json:"id"`
    Amount       float64   `json:"amount"`
    Currency     string    `json:"currency"`    
    Type         string    `json:"type"`         
    Status       string    `json:"status"`       
    Counterparty string    `json:"counterparty"` 
    Date         time.Time `json:"date"`
}

type TransactionCach struct {
    ID           string    `json:"id"`
    Amount       float64   `json:"amount"`
    Currency     string    `json:"currency"`    
    Type         string    `json:"type"`         
    Date         time.Time `json:"date"`
}

type TransactionListResponse struct {
    Transactions []TransactionItemDTO `json:"transactions"`
}

type TransactionCachSaved struct {
    Transactions []TransactionCach `json:"transactions"`
}
type WSEvent struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}