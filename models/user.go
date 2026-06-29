package models

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"-" db:"password"`
}
