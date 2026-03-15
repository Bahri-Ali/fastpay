package user

import "time"

type User struct {
    PhoneNumber    string    `json:"phone_number" gorm:"unique;not null"`
    Email          *string   `json:"email" gorm:"unique"`
    FullName       string    `json:"full_name"` 
    CreatedAt      time.Time `json:"created_at"`
}