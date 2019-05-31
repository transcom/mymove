package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// AccessCodeValidator is the service object interface for ValidateAccessCode
type AccessCodeValidator interface {
	ValidateAccessCode(code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error)
}

// AccessCodeFetcher is the service object interface for FetchAccessCode
type AccessCodeFetcher interface {
	FetchAccessCode(serviceMemberID uuid.UUID) (*models.AccessCode, error)
}

// AccessCodeClaimer is the service object interface for ValidateAccessCode
type AccessCodeClaimer interface {
	ClaimAccessCode(code string, serviceMemberID uuid.UUID) (*models.AccessCode, *validate.Errors, error)
}
