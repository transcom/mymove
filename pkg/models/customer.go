package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Customer struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	User      User      `belongs_to:"users"`
}
