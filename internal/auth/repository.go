package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context , user *User) error 
	GetUserByPhone(ctx context.Context, phone string) (*User, error)

}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repo{db: db}
}

func(r *repo) CreateUser(ctx context.Context , user *User) error {
	query := `
        INSERT INTO users (user_identifier, phone_number, email, password_hash, full_name, role, parent_id, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at
    `
	err := r.db.QueryRow(ctx, query,
        user.UserIdentifier,
        user.PhoneNumber,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.Role,
        user.ParentID,
        user.IsActive,
    ).Scan(&user.ID, &user.CreatedAt)

	if err != nil {return err}
	return nil

}

func ( r *repo) GetUserByPhone(ctx context.Context , phone string) (*User , error){

	query := `SELECT FROM users WHERE phone_number = $phone`
	user := &User{}
    err := r.db.QueryRow(ctx, query, phone).Scan(
        &user.ID,
        &user.UserIdentifier,
        &user.PhoneNumber,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.Role,
        &user.ParentID,
        &user.IsActive,
        &user.CreatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil 
		}
	}
	return user , err
}