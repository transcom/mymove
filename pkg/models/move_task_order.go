package models

import "github.com/gofrs/uuid"

type MoveTaskOrder struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type MoveTaskOrders []MoveTaskOrder
