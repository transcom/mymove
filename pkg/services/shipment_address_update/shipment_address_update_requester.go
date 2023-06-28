package shipmentaddressupdate

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	serviceparamvaluelookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentAddressUpdateRequester struct {
	planner           route.Planner
	addressCreator    services.AddressCreator
	moveRouter        services.MoveRouter
	shipmentSITStatus services.ShipmentSITStatus
}

func NewShipmentAddressUpdateRequester(planner route.Planner, addressCreator services.AddressCreator, moveRouter services.MoveRouter, shipmentSITStatus services.ShipmentSITStatus) services.ShipmentAddressUpdateRequester {

	return &shipmentAddressUpdateRequester{
		planner:           planner,
		addressCreator:    addressCreator,
		shipmentSITStatus: shipmentSITStatus,
		moveRouter:        moveRouter,
	}
}

func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeServiceArea(appCtx appcontext.AppContext, contractID uuid.UUID, originalDeliveryAddress models.Address, newDeliveryAddress models.Address) (bool, error) {
	var existingServiceArea models.ReZip3
	var actualServiceArea models.ReZip3

	originalZip := originalDeliveryAddress.PostalCode[0:3]
	destinationZip := newDeliveryAddress.PostalCode[0:3]

	if originalZip == destinationZip {
		// If the ZIP hasn't changed, we must be in the same service area
		return false, nil
	}

	err := appCtx.DB().Where("zip3 = ?", originalZip).Where("contract_id = ?", contractID).First(&existingServiceArea)
	if err != nil {
		return false, err
	}

	err = appCtx.DB().Where("zip3 = ?", destinationZip).Where("contract_id = ?", contractID).First(&actualServiceArea)
	if err != nil {
		return false, err
	}

	if existingServiceArea.DomesticServiceAreaID != actualServiceArea.DomesticServiceAreaID {
		return true, nil
	}
	return false, nil
}

// mileage bracket change (only applicable for linehaul)
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeMileageBracket(appCtx appcontext.AppContext, originalPickupAddress models.Address, originalDeliveryAddress, newDeliveryAddress models.Address) (bool, error) {

	// Mileage brackets are taken from the pricing spreadsheet, "2a) Domestic Linehaul Prices"
	// They are: [0, 250], [251, 500], [501, 1000], [1001, 1500], [1501-2000], [2001, 2500], [2501, 3000], [3001, 3500], [3501, 4000], and [4001, infinity)
	// We will handle the maximum bracket (>=4001 miles) separately.
	var milesLower = [9]int{0, 251, 501, 1001, 1501, 2001, 2501, 3001, 3501}
	var milesUpper = [9]int{250, 500, 1000, 1500, 2000, 2500, 3000, 3500, 4000}

	if originalDeliveryAddress.PostalCode == newDeliveryAddress.PostalCode {
		return false, nil
	}

	previousDistance, err := f.planner.ZipTransitDistance(appCtx, originalPickupAddress.PostalCode, originalDeliveryAddress.PostalCode)
	if err != nil {
		return false, err
	}
	newDistance, err := f.planner.ZipTransitDistance(appCtx, originalPickupAddress.PostalCode, newDeliveryAddress.PostalCode)
	if err != nil {
		return false, err
	}

	if previousDistance == newDistance {
		return false, nil
	}

	for index, lowerLimit := range milesLower {
		upperLimit := milesUpper[index]

		// Find the mileage bracket that the original shipment's distance falls into
		if previousDistance >= lowerLimit && previousDistance <= upperLimit {

			// If the new distance after the address change falls in a different bracket, then there could be a pricing change
			newDistanceIsInSameBracket := newDistance >= lowerLimit && newDistance <= upperLimit
			return !newDistanceIsInSameBracket, nil
		}
	}

	// if we get past the loop, then the original distance must be >=4001 miles, so we just have to check if
	// the new distance is also in this last bracket.
	if newDistance >= 4001 {
		return false, nil
	}
	return true, nil
}

// doesDeliveryAddressUpdateChangeShipmentPricingType checks if an address update would change a move from shorthaul to linehaul pricing or vice versa
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeShipmentPricingType(originalPickupAddress models.Address, originalDeliveryAddress models.Address, newDeliveryAddress models.Address) (bool, error) {
	originalZip := originalPickupAddress.PostalCode[0:3]
	originalDestinationZip := originalDeliveryAddress.PostalCode[0:3]
	newDestinationZip := newDeliveryAddress.PostalCode[0:3]

	isOriginalRouteShorthaul := originalZip == originalDestinationZip

	isNewRouteShorthaul := originalZip == newDestinationZip

	if isOriginalRouteShorthaul == isNewRouteShorthaul {
		return false, nil
	}
	return true, nil
}

// RequestShipmentDeliveryAddressUpdate is used to update the destination address of an HHG shipment without SIT after it has been approved by the TOO. If this update could result in excess cost for the customer, this service requires the change to go through TOO approval.
func (f *shipmentAddressUpdateRequester) RequestShipmentDeliveryAddressUpdate(appCtx appcontext.AppContext, shipmentID uuid.UUID, newAddress models.Address, contractorRemarks string, eTag string) (*models.ShipmentAddressUpdate, error) {
	var addressUpdate models.ShipmentAddressUpdate
	var shipment models.MTOShipment
	err := appCtx.DB().EagerPreload("MoveTaskOrder", "PickupAddress", "MTOServiceItems.ReService", "DestinationAddress").Find(&shipment, shipmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(shipmentID, "looking for shipment")
		}
		return nil, apperror.NewQueryError("MTOShipment", err, "")
	}

	if shipment.MoveTaskOrder.AvailableToPrimeAt == nil {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for moves that are available to the Prime")
	}
	if shipment.ShipmentType != models.MTOShipmentTypeHHG {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for HHG shipments")
	}
	if eTag != etag.GenerateEtag(shipment.UpdatedAt) {
		return nil, apperror.NewPreconditionFailedError(shipmentID, nil)
	}

	sitStatus, err := f.shipmentSITStatus.CalculateShipmentSITStatus(appCtx, shipment)
	if err != nil {
		return nil, err
	}
	if sitStatus != nil {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for shipments that do not use SIT")
	}

	err = appCtx.DB().EagerPreload("OriginalAddress", "NewAddress").Where("shipment_id = ?", shipmentID).First(&addressUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find an existing update, we'll need to make a new one
			addressUpdate.OriginalAddressID = *shipment.DestinationAddressID
			addressUpdate.OriginalAddress = *shipment.DestinationAddress
			addressUpdate.ShipmentID = shipmentID
			addressUpdate.OfficeRemarks = nil
		} else {
			return nil, err
		}
	}

	addressUpdate.Status = models.ShipmentAddressUpdateStatusApproved
	addressUpdate.ContractorRemarks = contractorRemarks
	address, err := f.addressCreator.CreateAddress(appCtx, &newAddress)
	if err != nil {
		return nil, err
	}
	addressUpdate.NewAddressID = address.ID

	contract, err := serviceparamvaluelookups.FetchContract(appCtx, *shipment.MoveTaskOrder.AvailableToPrimeAt)
	if err != nil {
		return nil, err
	}

	updateNeedsTOOReview, err := f.doesDeliveryAddressUpdateChangeServiceArea(appCtx, contract.ID, addressUpdate.OriginalAddress, newAddress)
	if err != nil {
		return nil, err
	}

	if !updateNeedsTOOReview {
		updateNeedsTOOReview, err = f.doesDeliveryAddressUpdateChangeShipmentPricingType(*shipment.PickupAddress, addressUpdate.OriginalAddress, newAddress)
		if err != nil {
			return nil, err
		}
	}

	if !updateNeedsTOOReview {
		updateNeedsTOOReview, err = f.doesDeliveryAddressUpdateChangeMileageBracket(appCtx, *shipment.PickupAddress, addressUpdate.OriginalAddress, newAddress)
		if err != nil {
			return nil, err
		}
	}

	if updateNeedsTOOReview {
		addressUpdate.Status = models.ShipmentAddressUpdateStatusRequested
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, txnErr := appCtx.DB().ValidateAndSave(&addressUpdate)
		if verrs.HasAny() {
			return apperror.NewInvalidInputError(addressUpdate.ID, txnErr, verrs, "unable to save ShipmentAddressUpdate")
		}
		if txnErr != nil {
			return apperror.NewQueryError("ShipmentAddressUpdate", txnErr, "error saving shipment address update request")
		}

		if updateNeedsTOOReview {
			err = f.moveRouter.SendToOfficeUser(appCtx, &shipment.MoveTaskOrder)
			if err != nil {
				return err
			}

			// If the update needs review, we need to manually make sure the etag gets updated
			shipment.UpdatedAt = time.Now()
		} else {
			shipment.DestinationAddressID = &addressUpdate.NewAddressID
		}

		verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(
				shipment.ID, err, verrs, "Invalid input found while updating shipment")
		}
		if err != nil {
			return err
		}

		return nil
	})
	if transactionError != nil {
		return nil, transactionError
	}

	return &addressUpdate, nil
}
