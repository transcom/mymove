package mtoserviceitem

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type createMTOServiceItemQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoServiceItemCreator struct {
	builder          createMTOServiceItemQueryBuilder
	createNewBuilder func() createMTOServiceItemQueryBuilder
	moveRouter       services.MoveRouter
}

// CreateMTOServiceItem creates a MTO Service Item
func (o *mtoServiceItemCreator) CreateMTOServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error) {
	var verrs *validate.Errors
	var err error
	var requestedServiceItems models.MTOServiceItems // used in case additional service items need to be auto-created
	var createdServiceItems models.MTOServiceItems

	var move models.Move
	moveID := serviceItem.MoveTaskOrderID
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveID),
	}
	// check if Move exists
	err = o.builder.FetchOne(appCtx, &move, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil, apperror.NewNotFoundError(moveID, "in Moves")
		default:
			return nil, nil, apperror.NewQueryError("Move", err, "")
		}
	}

	// Service items can only be created if a Move's status is either Approved
	// or Approvals Requested, so check and fail early.
	if move.Status != models.MoveStatusAPPROVED && move.Status != models.MoveStatusAPPROVALSREQUESTED {
		return nil, nil, apperror.NewConflictError(
			move.ID,
			fmt.Sprintf("Cannot create service items before a move has been approved. The current status for the move with ID %s is %s", move.ID, move.Status),
		)
	}

	// find the re service code id
	var reService models.ReService
	reServiceCode := serviceItem.ReService.Code
	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("code", "=", reServiceCode),
	}
	err = o.builder.FetchOne(appCtx, &reService, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service item with code: %s", reServiceCode))
		default:
			return nil, nil, apperror.NewQueryError("ReService", err, "")
		}
	}
	// set re service fields for service item
	serviceItem.ReServiceID = reService.ID
	serviceItem.ReService.Name = reService.Name

	// We can have two service items that come in from a MTO approval that do not have an MTOShipmentID
	// they are MTO level service items. This should capture that and create them accordingly, they are thankfully
	// also rather basic.
	if serviceItem.MTOShipmentID == nil {
		if serviceItem.ReService.Code == models.ReServiceCodeMS || serviceItem.ReService.Code == models.ReServiceCodeCS {
			serviceItem.Status = "APPROVED"
		}
		verrs, err = o.builder.CreateOne(appCtx, serviceItem)
		if verrs != nil {
			return nil, verrs, nil
		}
		if err != nil {
			return nil, nil, err
		}

		createdServiceItems = append(createdServiceItems, *serviceItem)

		return &createdServiceItems, nil, nil
	}

	// By the time the serviceItem model object gets here to the creator it should have a status attached to it.
	// If for some reason that isn't the case we will set it
	if serviceItem.Status == "" {
		serviceItem.Status = models.MTOServiceItemStatusSubmitted
	}

	// TODO: Once customer onboarding is built, we can revisit to figure out which service items goes under each type of shipment
	// check if shipment exists linked by MoveTaskOrderID
	var mtoShipment models.MTOShipment
	mtoShipmentID := *serviceItem.MTOShipmentID
	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoShipmentID),
		query.NewQueryFilter("move_id", "=", moveID),
	}
	err = o.builder.FetchOne(appCtx, &mtoShipment, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil, apperror.NewNotFoundError(mtoShipmentID, fmt.Sprintf("for mtoShipment with moveID: %s", moveID.String()))
		default:
			return nil, nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	if serviceItem.ReService.Code == models.ReServiceCodeDOSHUT || serviceItem.ReService.Code == models.ReServiceCodeDDSHUT {
		if mtoShipment.PrimeEstimatedWeight == nil {
			return nil, verrs, apperror.NewConflictError(mtoShipmentID, fmt.Sprintf("The associated MTOShipment (%s) must have a valid PrimeEstimatedWeight to create this service item.", mtoShipmentID))
		}
	}

	if serviceItem.ReService.Code == models.ReServiceCodeDOASIT {
		// DOASIT must be associated with shipment that has DOFSIT
		serviceItem, err = o.validateSITStandaloneServiceItem(appCtx, serviceItem, models.ReServiceCodeDOFSIT)
		if err != nil {
			return nil, nil, err
		}
	}

	if serviceItem.ReService.Code == models.ReServiceCodeDDASIT {
		// DDASIT must be associated with shipment that has DDFSIT
		serviceItem, err = o.validateSITStandaloneServiceItem(appCtx, serviceItem, models.ReServiceCodeDDFSIT)
		if err != nil {
			return nil, nil, err
		}
	}

	for index := range serviceItem.CustomerContacts {
		createCustContacts := &serviceItem.CustomerContacts[index]
		err = validateTimeMilitaryField(appCtx, createCustContacts.TimeMilitary)
		if err != nil {
			return nil, nil, apperror.NewInvalidInputError(serviceItem.ID, err, nil, err.Error())
		}
	}

	if serviceItem.ReService.Code == models.ReServiceCodeDDDSIT || serviceItem.ReService.Code == models.ReServiceCodeDOPSIT {
		verrs = validate.NewErrors()
		verrs.Add("reServiceCode", fmt.Sprintf("%s cannot be created", serviceItem.ReService.Code))
		return nil, nil, apperror.NewInvalidInputError(serviceItem.ID, nil, verrs,
			fmt.Sprintf("A service item with reServiceCode %s cannot be manually created.", serviceItem.ReService.Code))
	}

	updateShipmentPickupAddress := false
	if serviceItem.ReService.Code == models.ReServiceCodeDDFSIT || serviceItem.ReService.Code == models.ReServiceCodeDOFSIT {
		extraServiceItems, errSIT := o.validateFirstDaySITServiceItem(appCtx, serviceItem)
		if errSIT != nil {
			return nil, nil, errSIT
		}

		// update HHG origin address for ReServiceCodeDOFSIT service item
		if serviceItem.ReService.Code == models.ReServiceCodeDOFSIT {
			// When creating a DOFSIT, the prime must provide an HHG actual address for the move/shift in origin (pickup address)
			if serviceItem.SITOriginHHGActualAddress == nil {
				verrs = validate.NewErrors()
				verrs.Add("reServiceCode", fmt.Sprintf("%s cannot be created", serviceItem.ReService.Code))
				return nil, nil, apperror.NewInvalidInputError(serviceItem.ID, nil, verrs,
					fmt.Sprintf("A service item with reServiceCode %s must have the sitHHGActualOrigin field set.", serviceItem.ReService.Code))
			}

			// update the SIT service item to track/save the HHG original pickup address (that came from the
			// MTO shipment
			serviceItem.SITOriginHHGOriginalAddress = mtoShipment.PickupAddress.Copy()
			serviceItem.SITOriginHHGOriginalAddress.ID = uuid.Nil
			serviceItem.SITOriginHHGOriginalAddressID = nil

			// update the MTO shipment with the new (actual) pickup address
			mtoShipment.PickupAddress = serviceItem.SITOriginHHGActualAddress.Copy()
			mtoShipment.PickupAddress.ID = *mtoShipment.PickupAddressID // Keep to same ID to be updated with new values

			// changes were made to the shipment, needs to be saved to the database
			updateShipmentPickupAddress = true

			// Find the DOPSIT service item and update the SIT related address fields. These fields
			// will be used for pricing when a payment request is created for DOPSIT
			for itemIndex := range *extraServiceItems {
				extraServiceItem := &(*extraServiceItems)[itemIndex]
				if extraServiceItem.ReService.Code == models.ReServiceCodeDOPSIT ||
					extraServiceItem.ReService.Code == models.ReServiceCodeDOASIT {
					extraServiceItem.SITOriginHHGActualAddress = serviceItem.SITOriginHHGActualAddress
					extraServiceItem.SITOriginHHGActualAddressID = serviceItem.SITOriginHHGActualAddressID
					extraServiceItem.SITOriginHHGOriginalAddress = serviceItem.SITOriginHHGOriginalAddress
					extraServiceItem.SITOriginHHGOriginalAddressID = serviceItem.SITOriginHHGOriginalAddressID
				}

			}
		}

		requestedServiceItems = append(requestedServiceItems, *extraServiceItems...)
	}
	requestedServiceItems = append(requestedServiceItems, *serviceItem)

	// create new items in a transaction in case of failure
	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		if err != nil {
			txnAppCtx.Logger().Error(fmt.Sprintf("error starting txn: %v", err))
			return err
		}
		for serviceItemIndex := range requestedServiceItems {
			requestedServiceItem := &requestedServiceItems[serviceItemIndex]

			// create address if ID (UUID) is Nil
			if requestedServiceItem.SITOriginHHGActualAddress != nil {
				address := requestedServiceItem.SITOriginHHGActualAddress
				if address.ID == uuid.Nil {
					verrs, err = o.builder.CreateOne(txnAppCtx, address)
					if verrs != nil || err != nil {
						return fmt.Errorf("failed to save SITOriginHHGActualAddress: %#v %e", verrs, err)
					}
				}
				requestedServiceItem.SITOriginHHGActualAddressID = &address.ID
			}

			// create address if ID (UUID) is Nil
			if requestedServiceItem.SITOriginHHGOriginalAddress != nil {
				address := requestedServiceItem.SITOriginHHGOriginalAddress
				if address.ID == uuid.Nil {
					verrs, err = o.builder.CreateOne(txnAppCtx, address)
					if verrs != nil || err != nil {
						return fmt.Errorf("failed to save SITOriginHHGOriginalAddress: %#v %e", verrs, err)
					}
				}
				requestedServiceItem.SITOriginHHGOriginalAddressID = &address.ID
			}

			verrs, err = o.builder.CreateOne(txnAppCtx, requestedServiceItem)
			if verrs != nil || err != nil {
				return fmt.Errorf("%#v %e", verrs, err)
			}

			createdServiceItems = append(createdServiceItems, *requestedServiceItem)

			// create dimensions if any
			for index := range requestedServiceItem.Dimensions {
				createDimension := &requestedServiceItem.Dimensions[index]
				createDimension.MTOServiceItemID = requestedServiceItem.ID
				verrs, err = o.builder.CreateOne(txnAppCtx, createDimension)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "Failed to create dimensions")
				}
				if err != nil {
					return fmt.Errorf("%e", err)
				}
			}

			// create customer contacts if any
			for index := range requestedServiceItem.CustomerContacts {
				createCustContacts := &requestedServiceItem.CustomerContacts[index]
				createCustContacts.MTOServiceItemID = requestedServiceItem.ID
				verrs, err = o.builder.CreateOne(txnAppCtx, createCustContacts)
				if verrs != nil || err != nil {
					return fmt.Errorf("%#v %e", verrs, err)
				}
			}
		}

		// If updates were made to shipment, save update in the database
		if updateShipmentPickupAddress {
			verrs, err = o.builder.UpdateOne(txnAppCtx, mtoShipment.PickupAddress, nil)
			if verrs != nil || err != nil {
				return fmt.Errorf("failed to update mtoShipment.PickupAddress: %#v %e", verrs, err)
			}
		}

		if _, err = o.moveRouter.ApproveOrRequestApproval(txnAppCtx, move); err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		return nil, nil, transactionErr
	} else if verrs != nil && verrs.HasAny() {
		return nil, verrs, nil
	} else if err != nil {
		return nil, verrs, apperror.NewQueryError("unknown", err, "")
	}

	return &createdServiceItems, nil, nil
}

// checkDuplicateServiceCodes checks if the move or shipment has a duplicate service item with the same code as the one
// requested.
func (o *mtoServiceItemCreator) checkDuplicateServiceCodes(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem) error {
	var duplicateServiceItem models.MTOServiceItem

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("move_id", "=", serviceItem.MoveTaskOrderID),
		query.NewQueryFilter("re_service_id", "=", serviceItem.ReServiceID),
	}
	if serviceItem.MTOShipmentID != nil {
		queryFilters = append(queryFilters, query.NewQueryFilter("mto_shipment_id", "=", serviceItem.MTOShipmentID))
	}

	// We DON'T want to find this service item:
	err := o.builder.FetchOne(appCtx, &duplicateServiceItem, queryFilters)
	if err == nil && duplicateServiceItem.ID != uuid.Nil {
		return apperror.NewConflictError(duplicateServiceItem.ID,
			fmt.Sprintf("for creating a service item. A service item with reServiceCode %s already exists for this move and/or shipment.", serviceItem.ReService.Code))
	} else if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		return err
	}

	return nil
}

// makeExtraSITServiceItem sets up extra SIT service items if a first-day SIT service item is being created.
func (o *mtoServiceItemCreator) makeExtraSITServiceItem(appCtx appcontext.AppContext, firstSIT *models.MTOServiceItem, reServiceCode models.ReServiceCode) (*models.MTOServiceItem, error) {
	var reService models.ReService

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("code", "=", reServiceCode),
	}
	err := o.builder.FetchOne(appCtx, &reService, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service code: %s", reServiceCode))
		default:
			return nil, apperror.NewQueryError("ReService", err, "")
		}
	}

	extraServiceItem := models.MTOServiceItem{
		MTOShipmentID:   firstSIT.MTOShipmentID,
		MoveTaskOrderID: firstSIT.MoveTaskOrderID,
		ReServiceID:     reService.ID,
		ReService:       reService,
		SITEntryDate:    firstSIT.SITEntryDate,
		SITPostalCode:   firstSIT.SITPostalCode,
		Reason:          firstSIT.Reason,
		Status:          models.MTOServiceItemStatusSubmitted,
	}

	return &extraServiceItem, nil
}

// NewMTOServiceItemCreator returns a new MTO service item creator
func NewMTOServiceItemCreator(builder createMTOServiceItemQueryBuilder, moveRouter services.MoveRouter) services.MTOServiceItemCreator {
	// used inside a transaction and mocking
	createNewBuilder := func() createMTOServiceItemQueryBuilder {
		return query.NewQueryBuilder()
	}

	return &mtoServiceItemCreator{builder: builder, createNewBuilder: createNewBuilder, moveRouter: moveRouter}
}

func validateTimeMilitaryField(appCtx appcontext.AppContext, timeMilitary string) error {
	if len(timeMilitary) == 0 {
		return nil
	} else if len(timeMilitary) != 5 {
		return fmt.Errorf("timeMilitary must be in format HHMMZ")
	}

	hours := timeMilitary[:2]
	minutes := timeMilitary[2:4]
	suffix := timeMilitary[len(timeMilitary)-1:]

	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		return fmt.Errorf("timeMilitary must have a valid number for hours")
	}

	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		return fmt.Errorf("timeMilitary must have a valid number for minutes")
	}

	if !(0 <= hoursInt) || !(hoursInt < 24) {
		return fmt.Errorf("timeMilitary hours must be between 00 and 23")
	}
	if !(0 <= minutesInt) || !(minutesInt < 60) {
		return fmt.Errorf("timeMilitary minutes must be between 00 and 59")
	}

	if suffix != "Z" {
		return fmt.Errorf("timeMilitary must end with 'Z'")
	}

	return nil
}

// Check if and address has and ID, if it does, it needs to match OG SIT
func (o *mtoServiceItemCreator) validateSITStandaloneServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, reServiceCode models.ReServiceCode) (*models.MTOServiceItem, error) {
	var mtoServiceItem models.MTOServiceItem
	var mtoShipmentID uuid.UUID
	var validReService models.ReService
	mtoShipmentID = *serviceItem.MTOShipmentID

	queryFilter := []services.QueryFilter{
		query.NewQueryFilter("code", "=", reServiceCode),
	}

	// Fetch the ID for the ReServiceCode passed in, so we can check the shipment for its existence
	err := o.builder.FetchOne(appCtx, &validReService, queryFilter)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service code: %s", validReService.Code))
		default:
			return nil, apperror.NewQueryError("ReService", err, "")
		}
	}

	mtoServiceItemQueryFilter := []services.QueryFilter{
		query.NewQueryFilter("mto_shipment_id", "=", mtoShipmentID),
		query.NewQueryFilter("re_service_id", "=", validReService.ID),
	}
	// Fetch the required first-day SIT item for the shipment
	err = o.builder.FetchOne(appCtx, &mtoServiceItem, mtoServiceItemQueryFilter)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("No matching first-day SIT service item found for shipment: %s", mtoShipmentID))
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	verrs := validate.NewErrors()

	// check if the address IDs are nil, if not they need to match the orginal SIT address
	if serviceItem.SITOriginHHGOriginalAddress != nil && serviceItem.SITOriginHHGOriginalAddress.ID != mtoServiceItem.SITOriginHHGOriginalAddress.ID {
		verrs.Add("SITOriginHHGOriginalAddressID", fmt.Sprintf("%s invalid SITOriginHHGOriginalAddressID", serviceItem.ReService.Code))
	}

	if serviceItem.SITOriginHHGActualAddress != nil && serviceItem.SITOriginHHGActualAddress.ID != mtoServiceItem.SITOriginHHGActualAddress.ID {
		verrs.Add("SITOriginHHGActualAddress", fmt.Sprintf("%s invalid SITOriginHHGActualAddressID", serviceItem.ReService.Code))
	}

	if verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(serviceItem.ID, nil, verrs,
			fmt.Sprintf("There was invalid input in the standalone service item %s", serviceItem.ID))

	}

	// If the required first-day SIT item exists, we can update the related
	// service item passed in with the parent item's field values

	serviceItem.SITEntryDate = mtoServiceItem.SITEntryDate
	serviceItem.SITDepartureDate = mtoServiceItem.SITDepartureDate
	serviceItem.SITPostalCode = mtoServiceItem.SITPostalCode
	serviceItem.Reason = mtoServiceItem.Reason

	return serviceItem, nil
}

// check if an address has an ID
func (o *mtoServiceItemCreator) validateFirstDaySITServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, error) {
	var extraServiceItems models.MTOServiceItems
	var extraServiceItem *models.MTOServiceItem

	// check if there's another First Day SIT item for this shipment
	err := o.checkDuplicateServiceCodes(appCtx, serviceItem)
	if err != nil {
		return nil, err
	}

	verrs := validate.NewErrors()

	// check if the address IDs are nil
	if serviceItem.SITOriginHHGOriginalAddress != nil && serviceItem.SITOriginHHGOriginalAddress.ID != uuid.Nil {
		verrs.Add("SITOriginHHGOriginalAddressID", fmt.Sprintf("%s invalid SITOriginHHGOriginalAddressID", serviceItem.SITOriginHHGOriginalAddress.ID))
	}

	if serviceItem.SITOriginHHGActualAddress != nil && serviceItem.SITOriginHHGActualAddress.ID != uuid.Nil {
		verrs.Add("SITOriginHHGActualAddress", fmt.Sprintf("%s invalid SITOriginHHGActualAddressID", serviceItem.SITOriginHHGActualAddress.ID))
	}

	if verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(serviceItem.ID, nil, verrs,
			fmt.Sprintf("There was invalid input in the service item %s", serviceItem.ID))
	}

	// create the extra service items for first day SIT
	var reServiceCodes []models.ReServiceCode

	switch serviceItem.ReService.Code {
	case models.ReServiceCodeDDFSIT:
		reServiceCodes = append(reServiceCodes, models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT)
	case models.ReServiceCodeDOFSIT:
		reServiceCodes = append(reServiceCodes, models.ReServiceCodeDOASIT, models.ReServiceCodeDOPSIT)
	default:
		verrs := validate.NewErrors()
		verrs.Add("reServiceCode", fmt.Sprintf("%s invalid code", serviceItem.ReService.Code))
		return nil, apperror.NewInvalidInputError(serviceItem.ID, nil, verrs,
			fmt.Sprintf("No additional items can be created for this service item with code %s", serviceItem.ReService.Code))

	}

	for _, code := range reServiceCodes {
		extraServiceItem, err = o.makeExtraSITServiceItem(appCtx, serviceItem, code)
		if err != nil {
			return nil, err
		}
		if extraServiceItem != nil {
			extraServiceItems = append(extraServiceItems, *extraServiceItem)
		}
	}

	return &extraServiceItems, nil
}
