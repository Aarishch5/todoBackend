package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"t_name" db:"t_name"`
	Description string    `json:"description" db:"description"`
}

func (t *Todo) Validate() error {
	t.Name = strings.TrimSpace(t.Name)

	if t.Name == "" {
		return errors.New("name is required")
	}
	if len(t.Name) > 200 {
		return errors.New("name must be at most 200 characters")
	}
	if len(t.Description) > 1000 {
		return errors.New("description must be at most 1000 characters")
	}
	return nil
}
