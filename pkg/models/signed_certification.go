package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// SignedCertification represents users acceptance
type SignedCertification struct {
	ID                uuid.UUID `json:"id" db:"id"`
	SubmittingUserID  uuid.UUID `json:"submitting_user_id" db:"submitting_user_id"`
	MoveID            uuid.UUID `json:"move_id" db:"move_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	CertificationText string    `json:"certification_text" db:"certification_text"`
	Signature         string    `json:"signature" db:"signature"`
	Date              time.Time `json:"date" db:"date"`
}

// SignedCertifications is not required by pop and may be deleted
type SignedCertifications []SignedCertification

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *SignedCertification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.CertificationText, Name: "CertificationText"},
		&validators.StringIsPresent{Field: s.Signature, Name: "Signature"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *SignedCertification) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *SignedCertification) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
