package mtoshipment

import (
	"context"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentReweighRequester struct {
	db *pop.Connection
}

// NewShipmentReweighRequester creates a new struct with the service dependencies
func NewShipmentReweighRequester(db *pop.Connection) services.ShipmentReweighRequester {
	return &shipmentReweighRequester{
		db,
	}
}

// RequestShipmentReweigh Requests the shipment reweigh
func (f *shipmentReweighRequester) RequestShipmentReweigh(ctx context.Context, shipmentID uuid.UUID) (*models.Reweigh, error) {
	shipment, err := f.findShipment(shipmentID)
	if err != nil {
		return nil, err
	}

	reweigh, err := f.createReweigh(ctx, shipment, checkReweighAllowed())

	return reweigh, err
}

func (f *shipmentReweighRequester) findShipment(shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment

	err := f.db.Q().
		Eager("Reweigh").
		Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (f *shipmentReweighRequester) createReweigh(ctx context.Context, shipment *models.MTOShipment, checks ...validator) (*models.Reweigh, error) {
	if verr := validateShipment(ctx, shipment, shipment, checks...); verr != nil {
		return nil, verr
	}

	now := time.Now()
	reweigh := models.Reweigh{
		RequestedBy: models.ReweighRequesterTOO,
		RequestedAt: now,
		ShipmentID:  shipment.ID,
	}

	verrs, dbErr := f.db.ValidateAndSave(&reweigh)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "Could not save the reweigh while requesting the reweigh as a TOO.")

		return nil, invalidInputError
	}
	if dbErr != nil {
		return nil, dbErr
	}

	reweigh.Shipment = *shipment

	return &reweigh, nil
}
