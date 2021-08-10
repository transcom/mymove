package accesscode

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// accessCodeValidator is a service object to validate an access code.
type accessCodeValidator struct {
}

// NewAccessCodeValidator creates a new struct with the service dependencies
func NewAccessCodeValidator() services.AccessCodeValidator {
	return &accessCodeValidator{}
}

// ValidateAccessCode validates an access code based upon the code and move type. A valid access
// code is assumed to have no `service_member_id`
func (v accessCodeValidator) ValidateAccessCode(appCfg appconfig.AppConfig, code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error) {
	ac := models.AccessCode{}
	err := appCfg.DB().
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
