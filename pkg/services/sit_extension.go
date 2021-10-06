package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//SITExtensionApprover is the service object interface for approving a SIT extension
//go:generate mockery --name SITExtensionApprover --disable-version-string
type SITExtensionApprover interface {
	ApproveSITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, approvedDays int, officeRemarks *string, eTag string) (*models.MTOShipment, error)
}

//SITExtensionDenier is the service object interface for denying a SIT extension
//go:generate mockery --name SITExtensionDenier --disable-version-string
type SITExtensionDenier interface {
	DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOShipment, error)
}

//SITExtensionCreator creates a SIT extension
type SITExtensionCreator interface {
	CreateSITExtension(appCtx appcontext.AppContext, sitExtension *models.SITExtension) (*models.SITExtension, error)
}

//SITExtensionCreatorAsTOO is the service object interface to create an approved SIT extension
//go:generate mockery --name SITExtensionCreatorAsTOO --disable-version-string
type SITExtensionCreatorAsTOO interface {
	CreateSITExtensionAsTOO(appCtx appcontext.AppContext, sitExtension *models.SITExtension, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}
