package sitextension

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type sitExtensionCreator struct {
	checks []sitExtensionValidator
}

// NewSitExtensionCreator creates a new struct with the service dependencies
func NewSitExtensionCreator() services.SITExtensionCreator {
	return &sitExtensionCreator{
		checks: []sitExtensionValidator{
			checkShipmentID(),
			checkRequiredFields(),
			checkSITExtensionPending(),
		},
	}
}

// CreateSITExtension creates a SIT extension
func (f *sitExtensionCreator) CreateSITExtension(appCtx appcontext.AppContext, sitExtension *models.SITExtension) (*models.SITExtension, error) {
	// Get existing shipment info
	shipment := &models.MTOShipment{}
	// Find the shipment, return error if not found
	err := appCtx.DB().Find(shipment, sitExtension.MTOShipmentID)

	if err != nil {
		return nil, services.NewNotFoundError(sitExtension.MTOShipmentID, "while looking for MTOShipment")
	}

	// Set status to pending if none is provided
	if sitExtension.Status == "" {
		sitExtension.Status = models.SITExtensionStatusPending
	}

	err = validateSITExtension(appCtx, *sitExtension, shipment, f.checks...)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndCreate(sitExtension)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the SIT extension.")
	} else if err != nil {
		return nil, services.NewQueryError("SITExtension", err, "")
	}

	return sitExtension, nil
}
