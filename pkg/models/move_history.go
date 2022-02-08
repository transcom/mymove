package models

import (
	"github.com/gofrs/uuid"
)

// MoveHistory captures a move's audit history.
// This struct doesn't have a database table it is used for the ghc.yaml/GHC API.
// NO DATABASE TABLE
type MoveHistory struct {
	ID             uuid.UUID
	Locator        string
	ReferenceID    *string
	AuditHistories AuditHistories
}
