package accesscode

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// accessCodeFetcher is a service object to fetch an access code.
type accessCodeFetcher struct {
	DB *pop.Connection
}

// NewAccessCodeFetcher creates a new struct with the service dependencies.
func NewAccessCodeFetcher(db *pop.Connection) services.AccessCodeFetcher {
	return &accessCodeFetcher{db}
}

// FetchAccessCode fetches an access code based upon the service member id.
func (f accessCodeFetcher) FetchAccessCode(serviceMemberID uuid.UUID) (*models.AccessCode, error) {
	ac := models.AccessCode{}
	err := f.DB.
		Where("service_member_id = ?", serviceMemberID).
		First(&ac)

	if err != nil {
		return &ac, err
	}

	return &ac, nil
}
