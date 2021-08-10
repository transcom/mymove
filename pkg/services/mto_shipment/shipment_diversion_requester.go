package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentDiversionRequester struct {
	router services.ShipmentRouter
}

// NewShipmentDiversionRequester creates a new struct with the service dependencies
func NewShipmentDiversionRequester(router services.ShipmentRouter) services.ShipmentDiversionRequester {
	return &shipmentDiversionRequester{
		router,
	}
}

// RequestShipmentDiversion Requests the shipment diversion
func (f *shipmentDiversionRequester) RequestShipmentDiversion(appCfg appconfig.AppConfig, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCfg, shipmentID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, services.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	err = f.router.RequestDiversion(appCfg, shipment)
	if err != nil {
		return nil, err
	}

	verrs, err := appCfg.DB().ValidateAndSave(shipment)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "Could not validate shipment while requesting the diversion.")

		return nil, invalidInputError
	}

	return shipment, err
}

func (f *shipmentDiversionRequester) findShipment(appCfg appconfig.AppConfig, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCfg.DB().Q().Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}
