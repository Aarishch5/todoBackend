package models

import "github.com/google/uuid"

type Todo struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"t_name" db:"t_name"`
	Description string    `json:"description" db:"description"`
}
