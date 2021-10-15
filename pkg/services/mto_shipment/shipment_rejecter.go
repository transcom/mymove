package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentRejecter struct {
	router services.ShipmentRouter
}

// NewShipmentRejecter creates a new struct with the service dependencies
func NewShipmentRejecter(router services.ShipmentRouter) services.ShipmentRejecter {
	return &shipmentRejecter{
		router,
	}
}

// RejectShipment rejects the shipment
func (f *shipmentRejecter) RejectShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string, reason *string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	err = f.router.Reject(appCtx, shipment, reason)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndSave(shipment)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "Could not validate shipment while rejecting the shipment.")

		return nil, invalidInputError
	}

	return shipment, err
}

func (f *shipmentRejecter) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, apperror.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}
