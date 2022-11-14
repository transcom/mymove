package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
)

// SignedCertificationType represents the types of certificates
type SignedCertificationType string

const (
	// SignedCertificationTypePPM captures enum value "PPM" - deprecated 10/2020
	SignedCertificationTypePPM SignedCertificationType = "PPM"

	// SignedCertificationTypePPMPAYMENT captures enum value "PPM_PAYMENT"
	SignedCertificationTypePPMPAYMENT SignedCertificationType = "PPM_PAYMENT"

	// SignedCertificationTypeHHG captures enum value "HHG" - deprecated 10/2020
	SignedCertificationTypeHHG SignedCertificationType = "HHG"

	// SignedCertificationTypeSHIPMENT captures enum value "SHIPMENT" for all shipment types from 10/2020
	SignedCertificationTypeSHIPMENT SignedCertificationType = "SHIPMENT"
)

var AllowedSignedCertificationTypes = []string{
	string(SignedCertificationTypePPMPAYMENT),
	string(SignedCertificationTypeSHIPMENT),
}

// SignedCertification represents users acceptance
type SignedCertification struct {
	ID                       uuid.UUID                `json:"id" db:"id"`
	SubmittingUserID         uuid.UUID                `json:"submitting_user_id" db:"submitting_user_id"`
	MoveID                   uuid.UUID                `json:"move_id" db:"move_id"`
	PersonallyProcuredMoveID *uuid.UUID               `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	PpmID                    *uuid.UUID               `json:"ppm_id" db:"ppm_id"`
	CertificationType        *SignedCertificationType `json:"certification_type" db:"certification_type"`
	CreatedAt                time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time                `json:"updated_at" db:"updated_at"`
	CertificationText        string                   `json:"certification_text" db:"certification_text"`
	Signature                string                   `json:"signature" db:"signature"`
	Date                     time.Time                `json:"date" db:"date"`
}

// SignedCertifications is not required by pop and may be deleted
type SignedCertifications []SignedCertification

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (s *SignedCertification) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "SubmittingUserID", Field: s.SubmittingUserID},
		&validators.UUIDIsPresent{Name: "MoveID", Field: s.MoveID},
		&OptionalUUIDIsPresent{Name: "PersonallyProcuredMoveID", Field: s.PersonallyProcuredMoveID},
		&OptionalUUIDIsPresent{Name: "PpmID", Field: s.PpmID},
		&OptionalStringInclusion{Name: "CertificationType", Field: (*string)(s.CertificationType), List: AllowedSignedCertificationTypes},
		&validators.StringIsPresent{Name: "CertificationText", Field: s.CertificationText},
		&validators.StringIsPresent{Name: "Signature", Field: s.Signature},
		&validators.TimeIsPresent{Name: "Date", Field: s.Date},
	), nil
}

// DEPRECATED - This can be removed when the PPM Shipment Summary Worksheet is updated
// to use the new PPM shipment table
// FetchSignedCertificationsPPMPayment Fetches and Validates a PPM Payment Signature
func FetchSignedCertificationsPPMPayment(db *pop.Connection, session *auth.Session, id uuid.UUID) (*SignedCertification, error) {
	var signedCertification SignedCertification
	err := db.Where("move_id = $1", id.String()).Where("certification_type = $2", SignedCertificationTypePPMPAYMENT).First(&signedCertification)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
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
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
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
