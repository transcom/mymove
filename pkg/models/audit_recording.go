package models

import (
	"time"

	"github.com/gobuffalo/pop/slices"
	"github.com/gofrs/uuid"
)

// AuditRecording represents an event we'd like to keep track of
type AuditRecording struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	EventName    string     `json:"event_name" db:"event_name"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	RecordType   string     `json:"record_type" db:"record_type"`
	RecordData   slices.Map `json:"record_data" db:"record_data"`
	Payload      slices.Map `json:"payload" db:"payload"`
	Metadata     slices.Map `json:"metadata" db:"metadata"`
	MoveID       *uuid.UUID `json:"move_id" db:"move_id"`
	Move         Move       `belongs_to:"moves"`
	UserID       *uuid.UUID `json:"user_id" db:"user_id"`
	User         User       `belongs_to:"users"`
	FirstName    *string    `json:"first_name" db:"first_name"`
	LastName     *string    `json:"last_name" db:"last_name"`
	Email        *string    `json:"email" db:"email"`
	ClientCertID *uuid.UUID `json:"client_cert_id" db:"client_cert_id"`
	ClientCert   ClientCert `belongs_to:"client_certs"`
}
