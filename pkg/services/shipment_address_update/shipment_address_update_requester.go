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
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
)

type shipmentAddressUpdateRequester struct {
	planner         route.Planner
	addressCreator  services.AddressCreator
	moveRouter      services.MoveRouter
	shipmentFetcher services.MTOShipmentFetcher
	services.MTOServiceItemUpdater
	services.MTOServiceItemCreator
}

func NewShipmentAddressUpdateRequester(planner route.Planner, addressCreator services.AddressCreator, moveRouter services.MoveRouter) services.ShipmentAddressUpdateRequester {

	return &shipmentAddressUpdateRequester{
		planner:        planner,
		addressCreator: addressCreator,
		moveRouter:     moveRouter,
	}
}

func (f *shipmentAddressUpdateRequester) isAddressChangeDistanceOver50(appCtx appcontext.AppContext, addressUpdate models.ShipmentAddressUpdate) (bool, error) {

	//We calculate and set the distance between the old and new address
	distance, err := f.planner.ZipTransitDistance(appCtx, addressUpdate.OriginalAddress.PostalCode, addressUpdate.NewAddress.PostalCode)
	if err != nil {
		return false, err
	}

	if distance <= 50 {
		return false, nil
	}
	return true, nil
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

func (f *shipmentAddressUpdateRequester) doesShipmentContainApprovedDestinationSIT(shipment models.MTOShipment) bool {
	if len(shipment.MTOServiceItems) > 0 {
		serviceItems := shipment.MTOServiceItems

		for _, serviceItem := range serviceItems {
			serviceCode := serviceItem.ReService.Code
			status := serviceItem.Status
			if (serviceCode == models.ReServiceCodeDDASIT || serviceCode == models.ReServiceCodeDDDSIT || serviceCode == models.ReServiceCodeDDFSIT || serviceCode == models.ReServiceCodeDDSFSC) &&
				status == models.MTOServiceItemStatusApproved {
				return true
			}
		}
	}
	return false
}

func (f *shipmentAddressUpdateRequester) mapServiceItemWithUpdatedPriceRequirements(originalServiceItem models.MTOServiceItem) models.MTOServiceItem {
	var reService models.ReService
	now := time.Now()

	if originalServiceItem.ReService.Code == models.ReServiceCodeDSH {
		reService = models.ReService{
			Code: models.ReServiceCodeDLH,
		}
	} else if originalServiceItem.ReService.Code == models.ReServiceCodeDLH {
		reService = models.ReService{
			Code: models.ReServiceCodeDSH,
		}
	} else {
		reService = originalServiceItem.ReService
	}

	newServiceItem := models.MTOServiceItem{
		MTOShipmentID:                   originalServiceItem.MTOShipmentID,
		MoveTaskOrderID:                 originalServiceItem.MoveTaskOrderID,
		ReService:                       reService,
		SITEntryDate:                    originalServiceItem.SITEntryDate,
		SITDepartureDate:                originalServiceItem.SITDepartureDate,
		SITPostalCode:                   originalServiceItem.SITPostalCode,
		Reason:                          originalServiceItem.Reason,
		Status:                          models.MTOServiceItemStatusApproved,
		CustomerContacts:                originalServiceItem.CustomerContacts,
		PickupPostalCode:                originalServiceItem.PickupPostalCode,
		SITCustomerContacted:            originalServiceItem.SITCustomerContacted,
		SITRequestedDelivery:            originalServiceItem.SITRequestedDelivery,
		SITOriginHHGOriginalAddressID:   originalServiceItem.SITOriginHHGOriginalAddressID,
		SITOriginHHGActualAddressID:     originalServiceItem.SITOriginHHGActualAddressID,
		SITDestinationOriginalAddressID: originalServiceItem.SITDestinationOriginalAddressID,
		SITDestinationFinalAddressID:    originalServiceItem.SITDestinationFinalAddressID,
		Description:                     originalServiceItem.Description,
		EstimatedWeight:                 originalServiceItem.EstimatedWeight,
		ActualWeight:                    originalServiceItem.ActualWeight,
		Dimensions:                      originalServiceItem.Dimensions,
		ServiceRequestDocuments:         originalServiceItem.ServiceRequestDocuments,
		CreatedAt:                       originalServiceItem.CreatedAt,
		ApprovedAt:                      &now,
	}

	return newServiceItem
}

func checkForApprovedPaymentRequestOnServiceItem(appCtx appcontext.AppContext, mtoShipment models.MTOShipment) (bool, error) {
	mtoShipmentSITPaymentServiceItems := models.PaymentServiceItems{}

	err := appCtx.DB().Q().
		Join("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Join("re_services", "re_services.id = mto_service_items.re_service_id").
		Join("payment_requests", "payment_requests.id = payment_service_items.payment_request_id").
		Eager("MTOServiceItem.ReService", "PaymentServiceItemParams.ServiceItemParamKey").
		Where("mto_service_items.mto_shipment_id = ($1)", mtoShipment.ID).
		Where("payment_requests.status IN ($2, $3, $4, $5)",
			models.PaymentRequestStatusReviewed,
			models.PaymentRequestStatusSentToGex,
			models.PaymentRequestStatusTppsReceived,
			models.PaymentRequestStatusPaid).
		Where("payment_service_items.status != $6", models.PaymentServiceItemStatusDenied).
		Where("re_services.code IN ($7, $8)", models.ReServiceCodeDSH, models.ReServiceCodeDLH).
		All(&mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return false, err
	}

	if len(mtoShipmentSITPaymentServiceItems) != 0 {
		return true, err
	}

	return false, err
}

// RequestShipmentDeliveryAddressUpdate is used to update the destination address of an HHG shipment after it has been approved by the TOO. If this update could result in excess cost for the customer, this service requires the change to go through TOO approval.
func (f *shipmentAddressUpdateRequester) RequestShipmentDeliveryAddressUpdate(appCtx appcontext.AppContext, shipmentID uuid.UUID, newAddress models.Address, contractorRemarks string, eTag string) (*models.ShipmentAddressUpdate, error) {
	var addressUpdate models.ShipmentAddressUpdate
	var shipment models.MTOShipment
	err := appCtx.DB().EagerPreload("MoveTaskOrder", "PickupAddress", "StorageFacility.Address", "MTOServiceItems.ReService", "DestinationAddress", "MTOServiceItems.SITDestinationOriginalAddress").Find(&shipment, shipmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(shipmentID, "looking for shipment")
		}
		return nil, apperror.NewQueryError("MTOShipment", err, "")
	}

	if shipment.MoveTaskOrder.AvailableToPrimeAt == nil {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for moves that are available to the Prime")
	}
	if shipment.ShipmentType != models.MTOShipmentTypeHHG && shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom {
		return nil, apperror.NewUnprocessableEntityError("destination address update requests can only be created for HHG and NTS-Release shipments")
	}
	if eTag != etag.GenerateEtag(shipment.UpdatedAt) {
		return nil, apperror.NewPreconditionFailedError(shipmentID, nil)
	}

	shipmentHasApprovedDestSIT := f.doesShipmentContainApprovedDestinationSIT(shipment)

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
	} else {
		addressUpdate.OriginalAddressID = *shipment.DestinationAddressID
		addressUpdate.OriginalAddress = *shipment.DestinationAddress
	}

	addressUpdate.Status = models.ShipmentAddressUpdateStatusApproved
	addressUpdate.ContractorRemarks = contractorRemarks
	address, err := f.addressCreator.CreateAddress(appCtx, &newAddress)
	if err != nil {
		return nil, err
	}
	addressUpdate.NewAddressID = address.ID
	addressUpdate.NewAddress = *address

	// if the shipment contains destination SIT service items, we need to update the addressUpdate data
	// with the SIT original address and calculate the distances between the old & new shipment addresses
	if shipmentHasApprovedDestSIT {
		serviceItems := shipment.MTOServiceItems
		for _, serviceItem := range serviceItems {
			serviceCode := serviceItem.ReService.Code
			if serviceCode == models.ReServiceCodeDDASIT || serviceCode == models.ReServiceCodeDDDSIT || serviceCode == models.ReServiceCodeDDFSIT || serviceCode == models.ReServiceCodeDDSFSC {
				if serviceItem.SITDestinationOriginalAddressID != nil {
					addressUpdate.SitOriginalAddressID = serviceItem.SITDestinationOriginalAddressID
				}
				if serviceItem.SITDestinationOriginalAddress != nil {
					addressUpdate.SitOriginalAddress = serviceItem.SITDestinationOriginalAddress
				}
			}
			// if we have updated the values we need, no need to keep looping through the service items
			if addressUpdate.SitOriginalAddress != nil && addressUpdate.SitOriginalAddressID != nil {
				break
			}
		}
		if addressUpdate.SitOriginalAddress == nil {
			return nil, apperror.NewUnprocessableEntityError("shipments with approved destination SIT must have a SIT destination original address")
		}
		var distanceBetweenNew int
		var distanceBetweenOld int
		// if there was data already in the table, we want the "new" mileage to be the "old" mileage
		// if there is NOT, then we will calculate the distance between the original SIT dest address & the previous shipment address
		if addressUpdate.NewSitDistanceBetween != nil {
			distanceBetweenOld = *addressUpdate.NewSitDistanceBetween
		} else {
			distanceBetweenOld, err = f.planner.ZipTransitDistance(appCtx, addressUpdate.SitOriginalAddress.PostalCode, addressUpdate.OriginalAddress.PostalCode)
		}
		if err != nil {
			return nil, err
		}

		// calculating distance between the new address update & the SIT
		distanceBetweenNew, err = f.planner.ZipTransitDistance(appCtx, addressUpdate.SitOriginalAddress.PostalCode, addressUpdate.NewAddress.PostalCode)
		if err != nil {
			return nil, err
		}
		addressUpdate.NewSitDistanceBetween = &distanceBetweenNew
		addressUpdate.OldSitDistanceBetween = &distanceBetweenOld
	} else {
		addressUpdate.SitOriginalAddressID = nil
		addressUpdate.SitOriginalAddress = nil
		addressUpdate.NewSitDistanceBetween = nil
		addressUpdate.OldSitDistanceBetween = nil
	}

	contract, err := serviceparamvaluelookups.FetchContract(appCtx, *shipment.MoveTaskOrder.AvailableToPrimeAt)
	if err != nil {
		return nil, err
	}

	updateNeedsTOOReview, err := f.doesDeliveryAddressUpdateChangeServiceArea(appCtx, contract.ID, addressUpdate.OriginalAddress, newAddress)
	if err != nil {
		return nil, err
	}

	if !updateNeedsTOOReview {
		if shipment.ShipmentType == models.MTOShipmentTypeHHG {
			updateNeedsTOOReview, err = f.doesDeliveryAddressUpdateChangeShipmentPricingType(*shipment.PickupAddress, addressUpdate.OriginalAddress, newAddress)
			if err != nil {
				return nil, err
			}
		} else if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
			updateNeedsTOOReview, err = f.doesDeliveryAddressUpdateChangeShipmentPricingType(shipment.StorageFacility.Address, addressUpdate.OriginalAddress, newAddress)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipment type must be either an HHG or NTSr")
		}
	}

	if !updateNeedsTOOReview {
		if shipment.ShipmentType == models.MTOShipmentTypeHHG {
			updateNeedsTOOReview, err = f.doesDeliveryAddressUpdateChangeMileageBracket(appCtx, *shipment.PickupAddress, addressUpdate.OriginalAddress, newAddress)
			if err != nil {
				return nil, err
			}
		} else if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
			updateNeedsTOOReview, err = f.doesDeliveryAddressUpdateChangeMileageBracket(appCtx, shipment.StorageFacility.Address, addressUpdate.OriginalAddress, newAddress)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipment type must be either an HHG or NTSr")
		}
	}

	if !updateNeedsTOOReview {
		updateNeedsTOOReview, err = f.isAddressChangeDistanceOver50(appCtx, addressUpdate)
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

		//Get the move
		var move models.Move
		err := txnAppCtx.DB().Find(&move, shipment.MoveTaskOrderID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(shipment.MoveTaskOrderID, "looking for Move")
			default:
				return apperror.NewQueryError("Move", err, "unable to retrieve move")
			}
		}

		existingMoveStatus := move.Status
		if updateNeedsTOOReview {
			err = f.moveRouter.SendToOfficeUser(appCtx, &shipment.MoveTaskOrder)
			if err != nil {
				return err
			}

			// Only update if the move status has actually changed
			if existingMoveStatus != move.Status {
				err = txnAppCtx.DB().Update(&move)
				if err != nil {
					return err
				}
			}
		} else {
			shipment.DestinationAddressID = &addressUpdate.NewAddressID

			// Update MTO Shipment Destination Service Items
			err = mtoshipment.UpdateDestinationSITServiceItemsAddress(appCtx, &shipment)
			if err != nil {
				return err
			}

			err = mtoshipment.UpdateDestinationSITServiceItemsSITDeliveryMiles(f.planner, appCtx, &shipment, &addressUpdate.NewAddress, updateNeedsTOOReview)
			if err != nil {
				return err
			}
		}

		// If the request needs TOO review, this will just update the UpdatedAt timestamp on the shipment
		verrs, err = appCtx.DB().ValidateAndUpdate(&shipment)
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

func (f *shipmentAddressUpdateRequester) ReviewShipmentAddressChange(appCtx appcontext.AppContext, shipmentID uuid.UUID, tooApprovalStatus models.ShipmentAddressUpdateStatus, tooRemarks string) (*models.ShipmentAddressUpdate, error) {
	var shipment models.MTOShipment
	var addressUpdate models.ShipmentAddressUpdate

	err := appCtx.DB().EagerPreload("Shipment", "Shipment.MoveTaskOrder", "Shipment.MTOServiceItems", "Shipment.PickupAddress", "OriginalAddress", "NewAddress", "SitOriginalAddress", "Shipment.DestinationAddress", "Shipment.StorageFacility.Address").Where("shipment_id = ?", shipmentID).First(&addressUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(shipmentID, "looking for shipment address update")
		}
		return nil, apperror.NewQueryError("ShipmentAddressUpdate", err, "")
	}

	shipment = addressUpdate.Shipment

	if tooApprovalStatus == models.ShipmentAddressUpdateStatusApproved {
		queryBuilder := query.NewQueryBuilder()
		serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(f.planner, queryBuilder, f.moveRouter, f.shipmentFetcher, f.addressCreator)
		serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(f.planner, queryBuilder, f.moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		addressUpdate.Status = models.ShipmentAddressUpdateStatusApproved
		addressUpdate.OfficeRemarks = &tooRemarks
		shipment.DestinationAddress = &addressUpdate.NewAddress
		shipment.DestinationAddressID = &addressUpdate.NewAddressID

		var haulPricingTypeHasChanged bool
		if shipment.ShipmentType == models.MTOShipmentTypeHHG {
			haulPricingTypeHasChanged, err = f.doesDeliveryAddressUpdateChangeShipmentPricingType(*shipment.PickupAddress, addressUpdate.OriginalAddress, addressUpdate.NewAddress)
			if err != nil {
				return nil, err
			}
		} else if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
			haulPricingTypeHasChanged, err = f.doesDeliveryAddressUpdateChangeShipmentPricingType(shipment.StorageFacility.Address, addressUpdate.OriginalAddress, addressUpdate.NewAddress)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, apperror.NewInvalidInputError(shipment.ID, nil, nil, "Shipment type must be either an HHG or NTSr")
		}

		var shipmentDetails models.MTOShipment
		err = appCtx.DB().EagerPreload("MoveTaskOrder", "MTOServiceItems.ReService").Find(&shipmentDetails, shipmentID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, apperror.NewNotFoundError(shipmentID, "looking for shipment")
			}
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}

		// If the pricing type has changed then we automatically reject the DLH or DSH service item on the shipment since it is now inaccurate
		var approvedPaymentRequestsExistsForServiceItem bool
		if haulPricingTypeHasChanged && len(shipment.MTOServiceItems) > 0 {
			serviceItems := shipment.MTOServiceItems
			autoRejectionRemark := "Automatically rejected due to change in destination address affecting the ZIP code qualification for short haul / line haul."
			var regeneratedServiceItems models.MTOServiceItems

			for i, serviceItem := range shipmentDetails.MTOServiceItems {
				if (serviceItem.ReService.Code == models.ReServiceCodeDSH || serviceItem.ReService.Code == models.ReServiceCodeDLH) && serviceItem.Status != models.MTOServiceItemStatusRejected {
					// check if a payment request for the DSH or DLH service item exists and status is approved, paid, or sent to GEX
					approvedPaymentRequestsExistsForServiceItem, err = checkForApprovedPaymentRequestOnServiceItem(appCtx, shipment)
					if err != nil {
						return nil, apperror.NewQueryError("ServiceItemPaymentRequests", err, "")
					}

					shipmentHasApprovedDestSIT := f.doesShipmentContainApprovedDestinationSIT(shipmentDetails)

					// do NOT regenerate any service items if the following conditions exist:
					// payment has already been approved for DLH or DSH service item
					// destination SIT is on shipment and any of the service items have an appproved status
					if !approvedPaymentRequestsExistsForServiceItem && !shipmentHasApprovedDestSIT {
						rejectedServiceItem, updateErr := serviceItemUpdater.ApproveOrRejectServiceItem(appCtx, serviceItem.ID, models.MTOServiceItemStatusRejected, &autoRejectionRemark, etag.GenerateEtag(serviceItem.UpdatedAt))
						if updateErr != nil {
							return nil, updateErr
						}
						copyOfServiceItem := f.mapServiceItemWithUpdatedPriceRequirements(*rejectedServiceItem)
						serviceItems[i] = *rejectedServiceItem

						// Regenerate approved service items to replace the rejected ones.
						// Ensure that the updated pricing is applied (e.g. DLH -> DSH, DSH -> DLH etc.)
						regeneratedServiceItem, _, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &copyOfServiceItem)
						if createErr != nil {
							return nil, createErr
						}
						regeneratedServiceItems = append(regeneratedServiceItems, *regeneratedServiceItem...)
						break
					}

				}

			}

			// Append the auto-generated service items to the shipment service items slice
			if len(regeneratedServiceItems) > 0 {
				addressUpdate.Shipment.MTOServiceItems = append(addressUpdate.Shipment.MTOServiceItems, regeneratedServiceItems...)
			}
		}

		// handling NTS shipments that don't have a pickup address
		var pickupAddress models.Address
		if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
			pickupAddress = shipment.StorageFacility.Address
		} else {
			pickupAddress = *shipment.PickupAddress
		}
		// need to assess if the shipment's market code should change
		// when populating the market_code column, it is considered domestic if both pickup & the NEW dest are CONUS addresses
		if pickupAddress.IsOconus != nil && addressUpdate.NewAddress.IsOconus != nil {
			newAddress := addressUpdate.NewAddress
			if !*pickupAddress.IsOconus && !*newAddress.IsOconus {
				marketCodeDomestic := models.MarketCodeDomestic
				shipment.MarketCode = marketCodeDomestic
			} else {
				marketCodeInternational := models.MarketCodeInternational
				shipment.MarketCode = marketCodeInternational
			}
		}
	}

	if tooApprovalStatus == models.ShipmentAddressUpdateStatusRejected {
		addressUpdate.Status = models.ShipmentAddressUpdateStatusRejected
		addressUpdate.OfficeRemarks = &tooRemarks
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		verrs, txnErr := appCtx.DB().ValidateAndSave(&addressUpdate)
		if verrs.HasAny() {
			return apperror.NewInvalidInputError(addressUpdate.ID, txnErr, verrs, "unable to save ShipmentAddressUpdate")
		}
		if txnErr != nil {
			return apperror.NewQueryError("ShipmentAddressUpdate", txnErr, "error saving shipment address update request")
		}

		verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(
				shipment.ID, err, verrs, "Invalid input found while updating shipment")
		}
		if err != nil {
			return err
		}

		if len(shipment.MTOServiceItems) > 0 {
			err = mtoshipment.UpdateDestinationSITServiceItemsAddress(appCtx, &shipment)
		}
		if err != nil {
			return err
		}

		if len(shipment.MTOServiceItems) > 0 {
			err = mtoshipment.UpdateDestinationSITServiceItemsSITDeliveryMiles(f.planner, appCtx, &shipment, &addressUpdate.NewAddress, true)
		}
		if err != nil {
			return err
		}

		_, err = f.moveRouter.ApproveOrRequestApproval(appCtx, shipment.MoveTaskOrder)
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
