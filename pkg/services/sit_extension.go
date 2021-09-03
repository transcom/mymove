package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//SITExtensionApprover is the service object interface for approving a SIT extension
//go:generate mockery --name SITExtension --disable-version-string
type SITExtensionApprover interface {
	ApproveSITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, approvedDays int, officeRemarks *string, eTag string) (*models.MTOShipment, error)
}

//SITExtensionDenier is the service object interface for denying a SIT extension
//go:generate mockery --name SITExtension --disable-version-string
type SITExtensionDenier interface {
	DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOShipment, error)
}
