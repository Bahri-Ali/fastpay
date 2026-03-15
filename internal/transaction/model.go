package transaction

import "time"

type Transaction struct {
    ID          string    `json:"id"`
    SenderID    string    `json:"sender_id"`
    ReceiverID  string    `json:"receiver_id"`
    Amount      float64   `json:"amount"`
    Status      string    `json:"status"` // "completed", "pending"
    CreatedAt   time.Time `json:"created_at"`
}