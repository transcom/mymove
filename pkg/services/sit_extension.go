package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//SITExtension is the service object interface for approving and denying a SIT extension
//go:generate mockery --name SITExtension --disable-version-string
type SITExtension interface {
	ApproveSITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, approvedDays *int, officeRemarks *string, eTag string) (*models.MTOShipment, error)
}
