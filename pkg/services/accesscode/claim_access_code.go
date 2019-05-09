package accesscode

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// claimAccessCode is a service object to validate an access code.
type claimAccessCode struct {
	DB *pop.Connection
}

// NewAccessCodeClaimer creates a new struct with the service dependencies
func NewAccessCodeClaimer(db *pop.Connection) services.AccessCodeClaimer {
	return &claimAccessCode{db}
}

// ClaimAccessCode validates an access code based upon the code and move type. A valid access
// code is assumed to have no `service_member_id`
func (v claimAccessCode) ClaimAccessCode(code string, serviceMemberID uuid.UUID) (*models.AccessCode, error) {
	ac := models.AccessCode{}

	err := v.DB.
		Where("code = ?", code).
		First(&ac)

	// service object for fetching an access code by code and put into a model

	// service object for updating instance in the DB

	if err != nil {
		return &ac, err
	}

	return &ac, nil
}
