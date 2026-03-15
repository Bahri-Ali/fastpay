package wallet

import (
    "context"
    "database/sql"
    "errors"
    "github.com/jackc/pgx/v5/pgxpool"

)


type Repository interface {
	CreateWallet(ctx context.Context , Wallet *Wallet) error
	GetWalletByUserID(ctx context.Context , userID string) (*Wallet , error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db : db}
}

func (r *repo) CreateWallet(ctx context.Context , Wallet *Wallet) error{
	query := `
        INSERT INTO wallets (user_id, balance, currency, is_frozen)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
	return r.db.QueryRow(ctx , query ,
	 Wallet.UserID,
    Wallet.Balance,
    Wallet.Currency,
    Wallet.IsFrozen,
    ).Scan(&Wallet.ID , &Wallet.CreatedAt)   
}

func (r *repo) GetWalletByUserID(ctx context.Context, userID string) (*Wallet, error) {
    query := `
        SELECT id, user_id, balance, currency, is_frozen, created_at
        FROM wallets WHERE user_id = $1
    `
    w := &Wallet{}
    err := r.db.QueryRow(ctx, query, userID).Scan(
        &w.ID, &w.UserID, &w.Balance, &w.Currency, &w.IsFrozen, &w.CreatedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return w, nil
}
