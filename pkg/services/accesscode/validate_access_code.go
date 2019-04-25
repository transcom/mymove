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
// code is assumed to have no `user_id`
func (v validateAccessCode) ValidateAccessCode(code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error) {
	ac := models.AccessCode{}

	err := v.DB.
		Where("code = ?", code).
		Where("user_id IS NULL").
		Where("move_type = ?", moveType).
		First(&ac)

	if err != nil {
		return &ac, false, err
	}

	return &ac, true, nil
}
