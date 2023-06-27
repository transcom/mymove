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
	return false, nil
}

// mileage bracket change (only applicable for linehaul)
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeMileageBracket(appCtx appcontext.AppContext, contractID uuid.UUID, originalPickupAddress models.Address, originalDeliveryAddress, newDeliveryAddress models.Address) (bool, error) {
	// either look up both distances, and look up in hard coded list of brackets
	// or look up the linehaul price record for both and compare miles_upper and miles_lower
	//   this needs weight and isPeak as well.
	//   unless we can assume mileage brackets don't change within a contract, we could maybe aggregate and skip?
	return false, nil
}

// doesDeliveryAddressUpdateChangeShipmentPricingType checks if an address update would change a move from shorthaul to linehaul pricing or vice versa
func (f *shipmentAddressUpdateRequester) doesDeliveryAddressUpdateChangeShipmentPricingType(appCtx appcontext.AppContext, originalPickupAddress models.Address, originalDeliveryAddress models.Address, newDeliveryAddress models.Address) (bool, error) {
	return false, nil
}

// RequestShipmentDeliveryAddressUpdate is used to update the destination address of an HHG shipment without SIT after it has been approved by the TOO. If this update could result in excess cost for the customer, this service requires the change to go through TOO approval.
func (f *shipmentAddressUpdateRequester) RequestShipmentDeliveryAddressUpdate(appCtx appcontext.AppContext, shipmentID uuid.UUID, newAddress models.Address, contractorRemarks string) (*models.ShipmentAddressUpdate, error) {
	var addressUpdate models.ShipmentAddressUpdate
	var shipment models.MTOShipment
	err := appCtx.DB().EagerPreload("MoveTaskOrder", "PickupAddress", "MTOServiceItems", "MTOServiceItems.ReService").Find(&shipment, shipmentID)

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

	err = appCtx.DB().EagerPreload("OriginalAddress", "NewAddress").Where("shipment_id = ?", shipmentID).First(&addressUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find an existing update, we'll need to make a new one
			isThereAnExistingUpdate = false
			addressUpdate.OriginalAddressID = *shipment.DestinationAddressID
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
		if updateNeedsTOOReview {
			err = f.moveRouter.SendToOfficeUser(appCtx, &shipment.MoveTaskOrder)
			if err != nil {
				return err
			}

		} else {
			shipment.DestinationAddressID = &addressUpdate.NewAddressID
			verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(
					shipment.ID, err, verrs, "Invalid input found while saving updated destination address on shipment")
			}
			if err != nil {
				return err
			}
		}

		return nil
	})
	if transactionError != nil {
		return nil, transactionError
	}

	return &addressUpdate, nil
}
