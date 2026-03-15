package transaction

import (
    "context"
    "errors"
    // "fmt"
    // "math/rand"
    "time"

    // "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    ExecuteTransfer(ctx context.Context, senderID, receiverID string, amount float64) (*Transaction, error)
}

type repo struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repo{db: db}
}

// ExecuteTransfer performs the DB transaction (Lock -> Deduct -> Add -> Save)
func (r *repo) ExecuteTransfer(ctx context.Context, senderID, receiverID string, amount float64) (*Transaction, error) {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

    // 1. Lock Sender & Check Balance
    var senderBalance float64
    err = tx.QueryRow(ctx, 
        "SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE", senderID,
    ).Scan(&senderBalance)

    if err != nil {
        return nil, errors.New("sender wallet not found")
    }
    if senderBalance < amount {
        return nil, errors.New("insufficient balance")
    }

    // 2. Deduct from Sender
    _, err = tx.Exec(ctx, "UPDATE wallets SET balance = balance - $1 WHERE user_id = $2", amount, senderID)
    if err != nil {
        return nil, err
    }

    // 3. Add to Receiver
    res, err := tx.Exec(ctx, "UPDATE wallets SET balance = balance + $1 WHERE user_id = $2", amount, receiverID)
    if err != nil {
        return nil, err
    }
    if res.RowsAffected() == 0 {
        return nil, errors.New("receiver not found")
    }

    // 4. Save Transaction
    txn := &Transaction{
        SenderID:   senderID,
        ReceiverID: receiverID,
        Amount:     amount,
        Status:     "completed",
        CreatedAt:  time.Now(),
    }

    err = tx.QueryRow(ctx, `
        INSERT INTO transactions (sender_id, receiver_id, amount, status, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, senderID, receiverID, amount, "completed", txn.CreatedAt).Scan(&txn.ID)

    if err != nil {
        return nil, err
    }

    return txn, tx.Commit(ctx)
}