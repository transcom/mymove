package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Customer struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	DodID     *string   `db:"dod_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
