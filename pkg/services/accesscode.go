package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"

	"github.com/transcom/mymove/pkg/models"
)

// AccessCodeValidator is the service object interface for ValidateAccessCode
//go:generate mockery --name AccessCodeValidator --disable-version-string
type AccessCodeValidator interface {
	ValidateAccessCode(code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error)
}

// AccessCodeFetcher is the service object interface for FetchAccessCode
//go:generate mockery --name AccessCodeFetcher --disable-version-string
type AccessCodeFetcher interface {
	FetchAccessCode(serviceMemberID uuid.UUID) (*models.AccessCode, error)
}

// AccessCodeClaimer is the service object interface for ValidateAccessCode
//go:generate mockery --name AccessCodeClaimer --disable-version-string
type AccessCodeClaimer interface {
	ClaimAccessCode(code string, serviceMemberID uuid.UUID) (*models.AccessCode, *validate.Errors, error)
}

// OfficeAccessCodes is a slice of access codes for the office
type OfficeAccessCodes []adminmessages.AccessCode

// AccessCodeListFetcher is the service object interface for FetchAccessCodeList
//go:generate mockery --name AccessCodeListFetcher --disable-version-string
type AccessCodeListFetcher interface {
	FetchAccessCodeList(filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.AccessCodes, error)
	FetchAccessCodeCount(filters []QueryFilter) (int, error)
}
