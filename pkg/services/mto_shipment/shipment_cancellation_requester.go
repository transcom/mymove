package mtoshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentCancellationRequester struct {
	router services.ShipmentRouter
}

// NewShipmentCancellationRequester creates a new struct with the service dependencies
func NewShipmentCancellationRequester(router services.ShipmentRouter) services.ShipmentCancellationRequester {
	return &shipmentCancellationRequester{
		router,
	}
}

// RequestShipmentCancellation Requests the shipment diversion
func (f *shipmentCancellationRequester) RequestShipmentCancellation(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := FindShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	err = f.router.RequestCancellation(appCtx, shipment)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndSave(shipment)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "Could not validate shipment while requesting the shipment cancellation.")

		return nil, invalidInputError
	}

	return shipment, err
}
