package models

import "github.com/gofrs/uuid"

type ServiceItem struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type ServiceItems []ServiceItem
