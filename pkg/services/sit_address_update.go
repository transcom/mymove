package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ApprovedSITAddressUpdateCreator Interface for the service object that creates a SIT Address Update
//
//go:generate mockery --name SITAddressUpdateCreator
type ApprovedSITAddressUpdateCreator interface {
	CreateApprovedSITAddressUpdate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error)
}
