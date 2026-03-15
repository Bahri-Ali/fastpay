package user

import (
    "context"
    "database/sql"
    "errors"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    GetUserByID(ctx context.Context, id string) (*user, error)
    UpdatePassword(ctx context.Context, userID, newHashedPassword string) error
}

type repo struct {
    db *pgxpool.Pool
}


type user struct {
    ID             string
    UserIdentifier string
    PhoneNumber    string
    Email          *string
    FullName       string
    PasswordHash   string
    Role           string
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repo{db: db}
}

func (r *repo) GetUserByID(ctx context.Context, id string) (*user, error) {
    query := `SELECT id, user_identifier, phone_number, email, full_name, password_hash, role 
              FROM users WHERE id = $1`
    
    u := &user{}
    err := r.db.QueryRow(ctx, query, id).Scan(
        &u.ID, &u.UserIdentifier, &u.PhoneNumber, &u.Email, &u.FullName, &u.PasswordHash, &u.Role,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return u, nil
}

func (r *repo) UpdatePassword(ctx context.Context, userID, newHashedPassword string) error {
    query := `UPDATE users SET password_hash = $1 WHERE id = $2`
    _, err := r.db.Exec(ctx, query, newHashedPassword, userID)
    return err
}