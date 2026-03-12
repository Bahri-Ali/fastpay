package auth 

import (
	"time"
)

type UserRole string 
const(
	RoleNormal   UserRole = "normal"
    RoleMerchant UserRole = "merchant"
    RoleChild    UserRole = "child"
)

type User struct {
    ID             string    `json:"id" gorm:"type:uuid;unique;primaryKey;default:gen_random_uuid()"` // Technical ID (UUID)
    UserIdentifier string    `json:"user_identifier" gorm:"unique;not null"`                  // Custom Public ID (30 digits)
    PhoneNumber    string    `json:"phone_number" gorm:"unique;not null"`
    Email          *string   `json:"email" gorm:"unique"`
    PasswordHash   string    `json:"-"`
    FullName       string    `json:"full_name"`
    Role           UserRole  `json:"role" gorm:"type:user_role;default:'normal'"`
    ParentID       *string   `json:"parent_id" gorm:"type:uuid"` 
    IsActive       bool      `json:"is_active" gorm:"default:true"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type Session struct {
    ID            string    `json:"id" db:"id"`
    UserID        string    `json:"user_id" db:"user_id"`
    TokenHash     string    `json:"-" db:"token_hash"` 
    ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
    LastActivity  time.Time `json:"last_activity" db:"last_activity"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

func (User) users() string {
    return "users"
}