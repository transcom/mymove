package shipmentaddressupdate

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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

// service area change
// need old and new dest zips (destination service area?)
// i guess this changes unpack price and stuff like that, but not linehaul price?
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeServiceArea(appCtx appcontext.AppContext, contractID uuid.UUID, originalDeliveryAddress models.Address, newDeliveryAddress models.Address) (bool, error) {

	var existingServiceArea models.ReZip3
	var actualServiceArea models.ReZip3

	var originalZip string
	var destinationZip string

	originalZip = originalDeliveryAddress.PostalCode[0:3]
	destinationZip = newDeliveryAddress.PostalCode[0:3]

	if originalZip == destinationZip {
		actualServiceArea.DomesticServiceAreaID = existingServiceArea.DomesticServiceAreaID
		return false, nil
	}

	err := appCtx.DB().Where("zip3 = ?", originalZip).First(&existingServiceArea)
	if err != nil {
		return false, err
	}

	err = appCtx.DB().Where("zip3 = ?", destinationZip).First(&actualServiceArea)
	if err != nil {
		return false, err
	}

	if existingServiceArea.DomesticServiceAreaID != actualServiceArea.DomesticServiceAreaID {
		return true, nil
	}
	return false, nil
}

// mileage bracket change (only applicable for linehaul)
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeMileageBracket(appCtx appcontext.AppContext, contractID uuid.UUID, originalPickupAddress models.Address, originalDeliveryAddress, newDeliveryAddress models.Address) (bool, error) {
	// either look up both distances, and look up in hard coded list of brackets
	// or look up the linehaul price record for both and compare miles_upper and miles_lower
	//   this needs weight and isPeak as well.
	//   unless we can assume mileage brackets don't change within a contract, we could maybe aggregate and skip?

	var milesUpper = [9]int{250, 500, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	var milesLower = [9]int{0, 251, 501, 1001, 1501, 2001, 2501, 3001, 3501}

	if originalDeliveryAddress.PostalCode == newDeliveryAddress.PostalCode {
		return false, nil
	}

	perviousDistance, err := f.planner.ZipTransitDistance(appCtx, originalPickupAddress.PostalCode, originalDeliveryAddress.PostalCode)
	if err != nil {
		return false, nil
	}
	newDistance, err := f.planner.ZipTransitDistance(appCtx, originalPickupAddress.PostalCode, newDeliveryAddress.PostalCode)
	if err != nil {
		return false, nil
	}

	if perviousDistance == newDistance {
		return false, nil
	}

	for index, lowerLimit := range milesLower {

		upperLimit := milesUpper[index]

		if perviousDistance >= lowerLimit && perviousDistance <= upperLimit {

			if newDistance >= lowerLimit && newDistance <= upperLimit {
				return false, nil
			}
			return true, nil
		}
	}

	if newDistance >= 4001 {
		return false, nil
	}
	return true, nil
}

// doesDeliveryAddressUpdateChangeShipmentPricingType checks if an address update would change a move from shorthaul to linehaul pricing or vice versa
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeShipmentPricingType(appCtx appcontext.AppContext, originalPickupAddress models.Address, originalDeliveryAddress models.Address, newDeliveryAddress models.Address) (bool, error) {

	var originalZip models.ReZip3
	var originalDestinationZip models.ReZip3
	var newDestinationZip models.ReZip3

	originalZip.Zip3 = originalPickupAddress.PostalCode[0:3]
	originalDestinationZip.Zip3 = originalDeliveryAddress.PostalCode[0:3]
	newDestinationZip.Zip3 = newDeliveryAddress.PostalCode[0:3]

	isoriginalrouteshorthaul := originalZip.Zip3 == originalDestinationZip.Zip3

	isnewrouteshorthaul := originalDestinationZip.Zip3 == newDestinationZip.Zip3

	if isoriginalrouteshorthaul == isnewrouteshorthaul {
		return false, nil
	}
	return true, nil
}

// RequestShipmentDeliveryAddressUpdate is used to update the destination address of an HHG shipment without SIT after it has been approved by the TOO. If this update could result in excess cost for the customer, this service requires the change to go through TOO approval.
func (f *shipmentAddressUpdateRequester) RequestShipmentDeliveryAddressUpdate(appCtx appcontext.AppContext, shipmentID uuid.UUID, newAddress models.Address, contractorRemarks string) (*models.ShipmentAddressUpdate, error) {
	var addressUpdate models.ShipmentAddressUpdate
	var shipment models.MTOShipment
	err := appCtx.DB().EagerPreload("MoveTaskOrder", "PickupAddress", "MTOServiceItems", "MTOServiceItems.ReService", "DestinationAddress").Find(&shipment, shipmentID)

	if shipment.ShipmentType != models.MTOShipmentTypeHHG {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for HHG shipments")
	}
	sitStatus, err := f.shipmentSITStatus.CalculateShipmentSITStatus(appCtx, shipment)
	if err != nil {
		return nil, err
	}
	if sitStatus != nil {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for shipments that do not use SIT")
	}

	isThereAnExistingUpdate := true
	err = appCtx.DB().Where("shipment_id = ?", shipmentID).First(&addressUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find an existing update, we'll need to make a new one
			isThereAnExistingUpdate = false
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

	changesServiceArea, err := f.doesDeliveryAddressUpdateChangeServiceArea(appCtx, contract.ID, addressUpdate.OriginalAddress, newAddress)
	if err != nil {
		return nil, err
	}

	changesMileageBracket, err := f.doesDeliveryAddressUpdateChangeMileageBracket(appCtx, contract.ID, *shipment.PickupAddress, addressUpdate.OriginalAddress, newAddress)
	if err != nil {
		return nil, err
	}

	changesShipmentPricingType, err := f.doesDeliveryAddressUpdateChangeShipmentPricingType(appCtx, *shipment.PickupAddress, addressUpdate.OriginalAddress, newAddress)
	if err != nil {
		return nil, err
	}

	updateNeedsTOOReview := changesServiceArea || changesMileageBracket || changesShipmentPricingType
	if updateNeedsTOOReview {
		addressUpdate.Status = models.ShipmentAddressUpdateStatusRequested
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if isThereAnExistingUpdate {
			verrs, txnErr := appCtx.DB().ValidateAndSave(&addressUpdate)
			if verrs.HasAny() {
				return apperror.NewInvalidInputError(addressUpdate.ID, txnErr, verrs, "unable to save ShipmentAddressUpdate")
			}
			if txnErr != nil {
				return apperror.NewQueryError("ShipmentAddressUpdate", txnErr, "error saving shipment address update request")
			}
		} else {
			verrs, txnErr := appCtx.DB().ValidateAndCreate(&addressUpdate)
			if verrs.HasAny() {
				return apperror.NewInvalidInputError(uuid.Nil, txnErr, verrs, "unable to create ShipmentAddressUpdate")
			}
			if txnErr != nil {
				return apperror.NewQueryError("ShipmentAddressUpdate", txnErr, "error creating shipment address update request")
			}
		}

		err = f.moveRouter.SendToOfficeUser(appCtx, &shipment.MoveTaskOrder)
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
