package mtoshipment

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentRejecter struct {
	db     *pop.Connection
	router services.ShipmentRouter
}

// NewShipmentRejecter creates a new struct with the service dependencies
func NewShipmentRejecter(db *pop.Connection, router services.ShipmentRouter) services.ShipmentRejecter {
	return &shipmentRejecter{
		db,
		router,
	}
}

// RejectShipment rejects the shipment
func (f *shipmentRejecter) RejectShipment(shipmentID uuid.UUID, eTag string, reason *string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(shipmentID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, services.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	err = f.router.Reject(shipment, reason)
	if err != nil {
		return nil, err
	}

	verrs, err := f.db.ValidateAndSave(shipment)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "Could not validate shipment while rejecting the shipment.")

		return nil, invalidInputError
	}

	return shipment, err
}

func (f *shipmentRejecter) findShipment(shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := f.db.Q().Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}
