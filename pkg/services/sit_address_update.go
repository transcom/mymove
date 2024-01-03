package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ApprovedSITAddressUpdateRequestCreator Interface for the service object that creates an approved SIT Address Update
//
//go:generate mockery --name ApprovedSITAddressUpdateRequestCreator
type ApprovedSITAddressUpdateRequestCreator interface {
	CreateApprovedSITAddressUpdate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error)
}

// SITAddressUpdateRequestCreator creates a SIT Address Update Request with a distance greater than 50 miles
type SITAddressUpdateRequestCreator interface {
	CreateSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error)
}
