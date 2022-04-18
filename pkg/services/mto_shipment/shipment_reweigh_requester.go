package mtoshipment

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/db/utilities"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentReweighRequester struct {
}

// NewShipmentReweighRequester creates a new struct with the service dependencies
func NewShipmentReweighRequester() services.ShipmentReweighRequester {
	return &shipmentReweighRequester{}
}

// RequestShipmentReweigh Requests the shipment reweigh
func (f *shipmentReweighRequester) RequestShipmentReweigh(appCtx appcontext.AppContext, shipmentID uuid.UUID, requester models.ReweighRequester) (*models.Reweigh, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	reweigh, err := f.createReweigh(appCtx, shipment, requester, checkReweighAllowed())

	return reweigh, err
}

func (f *shipmentReweighRequester) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment

	err := appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope()).
		Eager("Reweigh").
		Find(&shipment, shipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipmentID, "while looking for shipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}

func (f *shipmentReweighRequester) createReweigh(appCtx appcontext.AppContext, shipment *models.MTOShipment, requester models.ReweighRequester, checks ...validator) (*models.Reweigh, error) {
	if verr := validateShipment(appCtx, shipment, shipment, checks...); verr != nil {
		return nil, verr
	}

	now := time.Now()
	reweigh := models.Reweigh{
		RequestedBy: requester,
		RequestedAt: now,
		ShipmentID:  shipment.ID,
	}

	verrs, dbErr := appCtx.DB().ValidateAndSave(&reweigh)
	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "Could not save the reweigh while requesting the reweigh as a TOO.")

		return nil, invalidInputError
	}
	if dbErr != nil {
		return nil, dbErr
	}

	return &reweigh, nil
}
