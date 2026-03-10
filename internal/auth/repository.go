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
func (r *repo) CreateUser(ctx context.Context, user *User) error {
    // تم إضافة علامات تنصيص حول "role" لأنها كلمة محجوزة
    query := `
        INSERT INTO users (user_identifier, phone_number, email, password_hash, full_name, "role", parent_id, is_active)
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

    return err
}

// GetUserByPhone fetches a user by phone number
func (r *repo) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
    // تم إضافة علامات تنصيص حول "role"
    query := `
        SELECT id, user_identifier, phone_number, email, password_hash, full_name, "role", parent_id, is_active, created_at
        FROM users
        WHERE phone_number = $1
    `

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
        return nil, err
    }

    return user, nil
}