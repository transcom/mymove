package models

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// SignedCertificationType represents the types of certificates
type SignedCertificationType string

const (
	// SignedCertificationTypePPM captures enum value "PPM"
	SignedCertificationTypePPM SignedCertificationType = "PPM"

	// SignedCertificationTypePPMPAYMENT captures enum value "PPM_PAYMENT"
	SignedCertificationTypePPMPAYMENT SignedCertificationType = "PPM_PAYMENT"

	// SignedCertificationTypeHHG captures enum value "HHG"
	SignedCertificationTypeHHG SignedCertificationType = "HHG"
)

var signedCertifications = []string{
	string(SignedCertificationTypePPM),
	string(SignedCertificationTypePPMPAYMENT),
	string(SignedCertificationTypeHHG),
}

// SignedCertification represents users acceptance
type SignedCertification struct {
	ID                       uuid.UUID                `json:"id" db:"id"`
	SubmittingUserID         uuid.UUID                `json:"submitting_user_id" db:"submitting_user_id"`
	MoveID                   uuid.UUID                `json:"move_id" db:"move_id"`
	PersonallyProcuredMoveID *uuid.UUID               `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	ShipmentID               *uuid.UUID               `json:"shipment_id" db:"shipment_id"`
	CertificationType        *SignedCertificationType `json:"certification_type" db:"certification_type"`
	CreatedAt                time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time                `json:"updated_at" db:"updated_at"`
	CertificationText        string                   `json:"certification_text" db:"certification_text"`
	Signature                string                   `json:"signature" db:"signature"`
	Date                     time.Time                `json:"date" db:"date"`
}

// SignedCertifications is not required by pop and may be deleted
type SignedCertifications []SignedCertification

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *SignedCertification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var ptrCertificationType *string
	if s.CertificationType != nil {
		certificationType := string(*s.CertificationType)
		ptrCertificationType = &certificationType
	}

	return validate.Validate(
		&validators.StringIsPresent{Field: s.CertificationText, Name: "CertificationText"},
		&validators.StringIsPresent{Field: s.Signature, Name: "Signature"},
		&OptionalStringInclusion{Field: ptrCertificationType, Name: "CertificationType", List: signedCertifications},
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

// FetchSignedCertificationsPPMPayment Fetches and Validates a PPM Payment Signature
func FetchSignedCertificationsPPMPayment(db *pop.Connection, session *auth.Session, id uuid.UUID) (*SignedCertification, error) {
	var signedCertification SignedCertification
	err := db.Where("move_id = $1", id.String()).Where("certification_type = $2", SignedCertificationTypePPMPAYMENT).First(&signedCertification)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			msg := fmt.Sprintf("signed_certification: with move_id: %s and certification_type: %s not found", id.String(), SignedCertificationTypePPMPAYMENT)
			return nil, errors.Wrap(ErrFetchNotFound, msg)
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, errors.Wrap(err, "signed_certification: unable to fetch signed certification")
	}
	// Validate the move is associated to the logged-in service member
	_, fetchErr := FetchMove(db, session, id)
	if fetchErr != nil {
		return nil, errors.Wrap(ErrFetchForbidden, "signed_certification: unauthorized access")
	}

	return &signedCertification, nil
}

// FetchSignedCertifications Fetches and Validates a all signed certifications associated with a move
func FetchSignedCertifications(db *pop.Connection, session *auth.Session, id uuid.UUID) ([]*SignedCertification, error) {
	var signedCertification []*SignedCertification
	err := db.Where("move_id = $1", id.String()).All(&signedCertification)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			msg := fmt.Sprintf("signed_certification: with move_id: %s not found", id.String())
			return nil, errors.Wrap(ErrFetchNotFound, msg)
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, errors.Wrap(err, "signed_certification: unable to fetch signed certification")
	}
	// Validate the move is associated to the logged-in service member
	_, fetchErr := FetchMove(db, session, id)
	if fetchErr != nil {
		return nil, errors.Wrap(ErrFetchForbidden, "signed_certification: unauthorized access")
	}

	return signedCertification, nil
}
