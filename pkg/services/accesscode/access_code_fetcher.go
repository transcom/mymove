package accesscode

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// accessCodeFetcher is a service object to fetch an access code.
type accessCodeFetcher struct {
}

// NewAccessCodeFetcher creates a new struct with the service dependencies.
func NewAccessCodeFetcher() services.AccessCodeFetcher {
	return &accessCodeFetcher{}
}

// FetchAccessCode fetches an access code based upon the service member id.
func (f accessCodeFetcher) FetchAccessCode(appCfg appconfig.AppConfig, serviceMemberID uuid.UUID) (*models.AccessCode, error) {
	ac := models.AccessCode{}
	err := appCfg.DB().
		Where("service_member_id = ?", serviceMemberID).
		First(&ac)

	if err != nil {
		return &ac, err
	}

	return &ac, nil
}
