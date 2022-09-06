package mtoshipment

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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

	if shipment.UsesExternalVendor {
		return &models.MTOShipment{}, apperror.NewConflictError(shipmentID, "shipmentApprover: shipment uses external vendor, cannot be approved for GHC Prime")
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return &models.MTOShipment{}, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
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
			invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")

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
	shipment, err := FindShipment(appCtx, shipmentID, "MoveTaskOrder", "PickupAddress", "DestinationAddress", "StorageFacility")

	if err != nil {
		return nil, err
	}

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "StorageFacility.Address" because
	// StorageFacility is a pointer.
	if shipment.StorageFacility != nil {
		err = appCtx.DB().Load(shipment.StorageFacility, "Address")
	}

	if err != nil {
		return nil, apperror.NewQueryError("MTOShipment", err, "")
	}

	if shipment.ShipmentType == models.MTOShipmentTypePPM {
		err = appCtx.DB().Load(shipment, "PPMShipment")
	}

	if err != nil {
		return nil, apperror.NewQueryError("MTOShipment", err, "")
	}

	return shipment, nil
}

func (f *shipmentApprover) setRequiredDeliveryDate(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.ScheduledPickupDate != nil &&
		shipment.RequiredDeliveryDate == nil &&
		(shipment.PrimeEstimatedWeight != nil || (shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom &&
			shipment.NTSRecordedWeight != nil)) {

		var pickupLocation *models.Address
		var deliveryLocation *models.Address
		var weight int

		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHGIntoNTSDom:
			if shipment.StorageFacility == nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGIntoNTSDom)
			}
			pickupLocation = shipment.PickupAddress
			deliveryLocation = &shipment.StorageFacility.Address
			weight = shipment.PrimeEstimatedWeight.Int()
		case models.MTOShipmentTypeHHGOutOfNTSDom:
			if shipment.StorageFacility == nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGOutOfNTSDom)
			}
			pickupLocation = &shipment.StorageFacility.Address
			deliveryLocation = shipment.DestinationAddress
			weight = shipment.NTSRecordedWeight.Int()
		default:
			pickupLocation = shipment.PickupAddress
			deliveryLocation = shipment.DestinationAddress
			weight = shipment.PrimeEstimatedWeight.Int()
		}
		requiredDeliveryDate, calcErr := CalculateRequiredDeliveryDate(appCtx, f.planner, *pickupLocation, *deliveryLocation, *shipment.ScheduledPickupDate, weight)
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
			invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue creating service items for the shipment")
			return invalidInputError
		}

		if err != nil {
			return err
		}
	}

	return nil
}
