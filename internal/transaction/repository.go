package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    ExecuteTransfer(ctx context.Context, senderID, receiverID string, amount float64) (*Transaction, error)
    GetTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error)}

type repo struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repo{db: db}
}

func (r *repo) ExecuteTransfer(ctx context.Context, senderID, receiverID string, amount float64) (*Transaction, error) {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

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

    _, err = tx.Exec(ctx, "UPDATE wallets SET balance = balance - $1 WHERE user_id = $2", amount, senderID)
    if err != nil {
        return nil, err
    }

    res, err := tx.Exec(ctx, "UPDATE wallets SET balance = balance + $1 WHERE user_id = $2", amount, receiverID)
    if err != nil {
        return nil, err
    }
    if res.RowsAffected() == 0 {
        return nil, errors.New("receiver not found")
    }
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

func (r *repo) GetTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error) {
    query := `
        SELECT id, sender_id, receiver_id, amount, status, created_at
        FROM transactions
        WHERE sender_id = $1 OR receiver_id = $1
        ORDER BY created_at DESC
    `

    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []Transaction
    for rows.Next() {
        var t Transaction
        err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.Status, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }

    return transactions, nil
}