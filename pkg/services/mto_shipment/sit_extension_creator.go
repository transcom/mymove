package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/services"
)

type sitExtensionCreator struct {
}

// NewSITExtensionCreator creates a new struct with the service dependencies
func NewSITExtensionCreator() services.SITExtensionCreator {
	return &sitExtensionCreator{}
}

// todo: should this return a shipment or just the sit extension? previous work did shipment
// CreateSITExtension creates a SIT Extension
func (f *sitExtensionCreator) CreateSITExtension(appCtx appcontext.AppContext, sitExtension *models.SITExtension, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	// todo: this is just so it passes now, actually write this
	return shipment, nil
}

func (f *sitExtensionCreator) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}
