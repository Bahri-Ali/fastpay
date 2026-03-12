package auth

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    CreateUser(ctx context.Context, user *User) error
    GetUserByPhone(ctx context.Context, phone string) (*User, error)

    CreateSession(ctx context.Context, session *Session) error
    GetSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error)
    UpdateSessionExpiry(ctx context.Context, tokenHash string, newExpiry time.Time) error
}

type repo struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repo{db: db}
}


func (r *repo) CreateUser(ctx context.Context, user *User) error {
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

    return err
}

func (r *repo) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
    query := `
        SELECT id, user_identifier, phone_number, email, password_hash, full_name, role, parent_id, is_active, created_at
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


func (r *repo) CreateSession(ctx context.Context, session *Session) error {
    query := `
        INSERT INTO sessions (user_id, token_hash, expires_at, last_activity)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
    err := r.db.QueryRow(ctx, query,
        session.UserID,
        session.TokenHash,
        session.ExpiresAt,
        session.LastActivity,
    ).Scan(&session.ID, &session.CreatedAt)

    return err
}

func (r *repo) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
    query := `
        SELECT id, user_id, token_hash, expires_at, last_activity, created_at
        FROM sessions 
        WHERE token_hash = $1
    `

    session := &Session{}
    err := r.db.QueryRow(ctx, query, tokenHash).Scan(
        &session.ID,
        &session.UserID,
        &session.TokenHash,
        &session.ExpiresAt,
        &session.LastActivity,
        &session.CreatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil 
        }
        return nil, err
    }

    return session, nil
}

func (r *repo) UpdateSessionExpiry(ctx context.Context, tokenHash string, newExpiry time.Time) error {
    query := `
        UPDATE sessions 
        SET expires_at = $1, last_activity = NOW() 
        WHERE token_hash = $2
    `
    _, err := r.db.Exec(ctx, query, newExpiry, tokenHash)
    return err
}