package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// SITAddressUpdateRequestCreator creates a SIT Address Update Request with a distance greater than 50 miles
type SITAddressUpdateRequestCreator interface {
	CreateSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error)
}
