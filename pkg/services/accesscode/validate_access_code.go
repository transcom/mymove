package accesscode

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// validateAccessCode is a service object to validate an access code.
type validateAccessCode struct {
	DB *pop.Connection
}

// NewAccessCodeValidator creates a new struct with the service dependencies
func NewAccessCodeValidator(db *pop.Connection) services.AccessCodeValidator {
	return &validateAccessCode{db}
}

// ValidateAccessCode validates an access code based upon the code and move type. A valid access
// code is assumed to have no `service_member_id`
func (v validateAccessCode) ValidateAccessCode(code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error) {
	ac := models.AccessCode{}
	err := v.DB.
		Where("code = ?", code).
		Where("move_type = ?", moveType).
		First(&ac)

	if err != nil {
		return &ac, false, err
	}

	if ac.ServiceMemberID != nil || ac.ClaimedAt != nil {
		return &ac, false, nil
	}

	return &ac, true, nil
}
