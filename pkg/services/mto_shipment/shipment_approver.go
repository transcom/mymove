package mtoshipment

import (
	"math"
	"slices"

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
	router               services.ShipmentRouter
	siCreator            services.MTOServiceItemCreator
	planner              route.Planner
	moveWeights          services.MoveWeights
	moveTaskOrderUpdater services.MoveTaskOrderUpdater
	moveRouter           services.MoveRouter
}

// NewShipmentApprover creates a new struct with the service dependencies
func NewShipmentApprover(router services.ShipmentRouter, siCreator services.MTOServiceItemCreator, planner route.Planner, moveWeights services.MoveWeights, moveTaskOrderUpdater services.MoveTaskOrderUpdater, moveRouter services.MoveRouter) services.ShipmentApprover {
	return &shipmentApprover{
		router,
		siCreator,
		planner,
		moveWeights,
		moveTaskOrderUpdater,
		moveRouter,
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

	// RequestedPickupDate must be in the future if set
	err = MTOShipmentHasValidRequestedPickupDate().Validate(appCtx, shipment, nil)
	if err != nil {
		return nil, err
	}

	err = f.router.Approve(appCtx, shipment)
	if err != nil {
		return nil, err
	}

	err = f.setRequiredDeliveryDate(appCtx, shipment)
	if err != nil {
		return nil, err
	}

	// if the shipment has an estimated weight at time of approval
	// recalculate the authorized weight to include the newly authorized shipment
	// and check for excess weight
	if shipment.PrimeEstimatedWeight != nil || shipment.NTSRecordedWeight != nil {
		err = f.updateAuthorizedWeight(appCtx, shipment)
		if err != nil {
			return nil, err
		}

		// changes to estimated weight need to run thecheck for excess weight
		_, verrs, err := f.moveWeights.CheckExcessWeight(appCtx, shipment.MoveTaskOrderID, *shipment)
		if verrs != nil && verrs.HasAny() {
			return nil, errors.New(verrs.Error())
		}
		if err != nil {
			return nil, err
		}
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// create international shipment service items before approving
		// we use a database proc to create the basic auto-approved service items
		internationalShipmentTypes := []models.MTOShipmentType{models.MTOShipmentTypeHHG, models.MTOShipmentTypeHHGIntoNTS, models.MTOShipmentTypeHHGOutOfNTS, models.MTOShipmentTypeUnaccompaniedBaggage}
		if slices.Contains(internationalShipmentTypes, shipment.ShipmentType) && shipment.MarketCode == models.MarketCodeInternational {
			err := models.CreateApprovedServiceItemsForShipment(appCtx.DB(), shipment)
			if err != nil {
				return err
			}

			// Update the service item pricing if we have the estimated weight
			if shipment.PrimeEstimatedWeight != nil {
				portZip, portType, err := models.GetPortLocationInfoForShipment(appCtx.DB(), shipment.ID)
				if err != nil {
					return err
				}
				// if we don't have the port data, then we won't worry about pricing
				if portZip != nil && portType != nil {
					var pickupZip string
					var destZip string
					// if the port type is POEFSC this means the shipment is CONUS -> OCONUS (pickup -> port)
					// if the port type is PODFSC this means the shipment is OCONUS -> CONUS (port -> destination)
					if *portType == models.ReServiceCodePOEFSC.String() {
						pickupZip = shipment.PickupAddress.PostalCode
						destZip = *portZip
					} else if *portType == models.ReServiceCodePODFSC.String() {
						pickupZip = *portZip
						destZip = shipment.DestinationAddress.PostalCode
					}
					// we need to get the mileage first, the db proc will consume that
					mileage, err := f.planner.ZipTransitDistance(appCtx, pickupZip, destZip)
					if err != nil {
						return err
					}

					// update the service item pricing if relevant fields have changed
					err = models.UpdateEstimatedPricingForShipmentBasicServiceItems(appCtx.DB(), shipment, &mileage)
					if err != nil {
						return err
					}
				} else {
					// if we don't have the port data, that's okay - we can update the other service items except for PODFSC/POEFSC
					err = models.UpdateEstimatedPricingForShipmentBasicServiceItems(appCtx.DB(), shipment, nil)
					if err != nil {
						return err
					}
				}
			}
		} else {
			// after approving shipment, shipment level service items must be created (this is for domestic shipments only)
			err = f.createShipmentServiceItems(txnAppCtx, shipment)
			if err != nil {
				return err
			}
		}

		verrs, err := txnAppCtx.DB().ValidateAndSave(shipment)
		if verrs != nil && verrs.HasAny() {
			invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")

			return invalidInputError
		}
		if err != nil {
			return err
		}

		var move models.Move
		move.ID = shipment.MoveTaskOrderID
		// re-evaluate move status
		if _, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, move); err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return shipment, nil
}

// ApproveShipments Approves one or more shipments in one transaction
func (f *shipmentApprover) ApproveShipments(appCtx appcontext.AppContext, shipments []services.ShipmentIdWithEtag) (*[]models.MTOShipment, error) {
	var approvedShipments []models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		for _, shipment := range shipments {
			shipmentID := shipment.ShipmentID
			eTag := shipment.ETag

			approvedShipment, err := f.ApproveShipment(txnAppCtx, shipmentID, eTag)
			if err != nil {
				return err
			}

			approvedShipments = append(approvedShipments, *approvedShipment)
		}

		return nil
	})

	return &approvedShipments, transactionError
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
		err = appCtx.DB().Load(shipment.StorageFacility, "Address", "Address.Country")
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
		(shipment.PrimeEstimatedWeight != nil || (shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS &&
			shipment.NTSRecordedWeight != nil)) {

		var pickupLocation *models.Address
		var deliveryLocation *models.Address
		var weight *int

		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHGIntoNTS:
			if shipment.StorageFacility == nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGIntoNTS)
			}
			pickupLocation = shipment.PickupAddress
			deliveryLocation = &shipment.StorageFacility.Address
			if shipment.PrimeEstimatedWeight != nil {
				weight = models.IntPointer(shipment.PrimeEstimatedWeight.Int())
			}
		case models.MTOShipmentTypeHHGOutOfNTS:
			if shipment.StorageFacility == nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGOutOfNTS)
			}
			pickupLocation = &shipment.StorageFacility.Address
			deliveryLocation = shipment.DestinationAddress
			if shipment.NTSRecordedWeight != nil {
				weight = models.IntPointer(shipment.NTSRecordedWeight.Int())
			}
		default:
			pickupLocation = shipment.PickupAddress
			deliveryLocation = shipment.DestinationAddress
			if shipment.PrimeEstimatedWeight != nil {
				weight = models.IntPointer(shipment.PrimeEstimatedWeight.Int())
			}
		}
		requiredDeliveryDate, calcErr := CalculateRequiredDeliveryDate(appCtx, f.planner, *pickupLocation, *deliveryLocation, *shipment.ScheduledPickupDate, weight, shipment.MoveTaskOrderID, shipment.ShipmentType)
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

// when a TOO approves a shipment, if it was created by PRIME and an estimated weight exists
// add that to the authorized weight
func (f *shipmentApprover) updateAuthorizedWeight(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	var move models.Move
	err := appCtx.DB().EagerPreload(
		"MTOShipments",
		"Orders.Entitlement",
	).Find(&move, shipment.MoveTaskOrderID)

	if err != nil {
		return apperror.NewQueryError("Move", err, "unable to find Move")
	}

	var dBAuthorizedWeight int
	if shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
		dBAuthorizedWeight = int(*shipment.PrimeEstimatedWeight)
	} else {
		dBAuthorizedWeight = int(*shipment.NTSRecordedWeight)
	}
	if len(move.MTOShipments) != 0 {
		for _, mtoShipment := range move.MTOShipments {
			if mtoShipment.Status == models.MTOShipmentStatusApproved && mtoShipment.ID != shipment.ID {
				if mtoShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
					//uses PrimeEstimatedWeight for HHG and NTS shipments
					if mtoShipment.PrimeEstimatedWeight != nil {
						dBAuthorizedWeight += int(*mtoShipment.PrimeEstimatedWeight)
					}
				} else {
					//used NTSRecordedWeight for NTSRShipments
					if mtoShipment.NTSRecordedWeight != nil {
						dBAuthorizedWeight += int(*mtoShipment.NTSRecordedWeight)
					}
				}
			}
		}
	}
	dBAuthorizedWeight = int(math.Round(float64(dBAuthorizedWeight) * 1.10))

	entitlement := move.Orders.Entitlement
	entitlement.DBAuthorizedWeight = &dBAuthorizedWeight
	verrs, err := appCtx.DB().ValidateAndUpdate(entitlement)

	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")
		return invalidInputError
	}
	if err != nil {
		return err
	}

	return nil
}
