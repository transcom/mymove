package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentApprover struct {
	router    services.ShipmentRouter
	siCreator services.MTOServiceItemCreator
	planner   route.Planner
}

// NewShipmentApprover creates a new struct with the service dependencies
func NewShipmentApprover(router services.ShipmentRouter, siCreator services.MTOServiceItemCreator, planner route.Planner) services.ShipmentApprover {
	return &shipmentApprover{
		router,
		siCreator,
		planner,
	}
}

// ApproveShipment Approves the shipment
func (f *shipmentApprover) ApproveShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, services.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	err = f.router.Approve(appCtx, shipment)
	if err != nil {
		return nil, err
	}

	err = f.setRequiredDeliveryDate(appCtx, shipment)
	if err != nil {
		return nil, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		verrs, err := txnAppCtx.DB().ValidateAndSave(shipment)
		if verrs != nil && verrs.HasAny() {
			invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")

			return invalidInputError
		}
		if err != nil {
			return err
		}

		// after approving shipment, shipment level service items must be created
		err = f.createShipmentServiceItems(txnAppCtx, shipment)
		if err != nil {
			return err
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return shipment, nil
}

func (f *shipmentApprover) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Eager("MoveTaskOrder", "PickupAddress", "DestinationAddress").Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (f *shipmentApprover) setRequiredDeliveryDate(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.ScheduledPickupDate != nil &&
		shipment.RequiredDeliveryDate == nil &&
		shipment.PrimeEstimatedWeight != nil {
		requiredDeliveryDate, calcErr := CalculateRequiredDeliveryDate(appCtx, f.planner, *shipment.PickupAddress, *shipment.DestinationAddress, *shipment.ScheduledPickupDate, shipment.PrimeEstimatedWeight.Int())
		if calcErr != nil {
			return calcErr
		}

		shipment.RequiredDeliveryDate = requiredDeliveryDate
	}

	return nil
}

func (f *shipmentApprover) createShipmentServiceItems(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	reServiceCodes := reServiceCodesForShipment(*shipment)
	serviceItemsToCreate := constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
	for _, serviceItem := range serviceItemsToCreate {
		copyOfServiceItem := serviceItem // Make copy to avoid implicit memory aliasing of items from a range statement.
		_, verrs, err := f.siCreator.CreateMTOServiceItem(appCtx, &copyOfServiceItem)

		if verrs != nil && verrs.HasAny() {
			invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue creating service items for the shipment")
			return invalidInputError
		}

		if err != nil {
			return err
		}
	}

	return nil
}
