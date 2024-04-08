package mtoserviceitem

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
)

// OriginSITLocation is the constant representing when the shipment in storage occurs at the origin
const OriginSITLocation = "ORIGIN"

// DestinationSITLocation is the constant representing when the shipment in storage occurs at the destination
const DestinationSITLocation = "DESTINATION"

// Number of days of grace period after customer contacts prime for delivery out of SIT
const GracePeriodDays = 5

type mtoServiceItemQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoServiceItemUpdater struct {
	planner          route.Planner
	builder          mtoServiceItemQueryBuilder
	createNewBuilder func() mtoServiceItemQueryBuilder
	moveRouter       services.MoveRouter
	shipmentFetcher  services.MTOShipmentFetcher
	addressCreator   services.AddressCreator
}

// NewMTOServiceItemUpdater returns a new mto service item updater
func NewMTOServiceItemUpdater(planner route.Planner, builder mtoServiceItemQueryBuilder, moveRouter services.MoveRouter, shipmentFetcher services.MTOShipmentFetcher, addressCreator services.AddressCreator) services.MTOServiceItemUpdater {
	// used inside a transaction and mocking		return &mtoServiceItemUpdater{builder: builder}
	createNewBuilder := func() mtoServiceItemQueryBuilder {
		return query.NewQueryBuilder()
	}

	return &mtoServiceItemUpdater{planner, builder, createNewBuilder, moveRouter, shipmentFetcher, addressCreator}
}

func (p *mtoServiceItemUpdater) ApproveOrRejectServiceItem(
	appCtx appcontext.AppContext,
	mtoServiceItemID uuid.UUID,
	status models.MTOServiceItemStatus,
	rejectionReason *string,
	eTag string,
) (*models.MTOServiceItem, error) {
	mtoServiceItem, err := p.findServiceItem(appCtx, mtoServiceItemID)
	if err != nil {
		return &models.MTOServiceItem{}, err
	}

	return p.approveOrRejectServiceItem(appCtx, *mtoServiceItem, status, rejectionReason, eTag, checkMoveStatus(), checkETag())
}

func (p *mtoServiceItemUpdater) ConvertItemToCustomerExpense(
	appCtx appcontext.AppContext,
	shipment *models.MTOShipment,
	customerExpenseReason *string,
	convertToCustomerExpense bool,
) (*models.MTOServiceItem, error) {
	var DOFSITCodeID, DDFSITCodeID uuid.UUID
	DOFSITServiceErr := appCtx.DB().RawQuery(`SELECT id FROM re_services WHERE code = 'DOFSIT'`).First(&DOFSITCodeID) // First get uuid for DOFSIT service code
	if DOFSITServiceErr != nil {
		return nil, apperror.NewNotFoundError(uuid.Nil, "Couldn't find entry for DOFSIT ReService code in re_services table.")
	}
	DDFSITServiceErr := appCtx.DB().RawQuery(`SELECT id FROM re_services WHERE code = 'DOFSIT'`).First(&DDFSITCodeID)
	if DDFSITServiceErr != nil {
		return nil, apperror.NewNotFoundError(uuid.Nil, "Couldn't find entry for DDFSIT ReService code in re_services table.")
	}

	sitStatusService := sitstatus.NewShipmentSITStatus()
	shipmentSITStatus, err := sitStatusService.CalculateShipmentSITStatus(appCtx, *shipment)
	if err != nil {
		return nil, err
	} else if shipmentSITStatus == nil {
		return nil, apperror.NewNotFoundError(shipment.ID, "for current SIT MTO Service Item.")
	}

	// Now get the service item associated with the current mto_shipment
	var SITItem models.MTOServiceItem
	getSITItemErr := appCtx.DB().RawQuery(`SELECT * FROM mto_service_items WHERE id = ?`, shipmentSITStatus.CurrentSIT.ServiceItemID).First(&SITItem)
	if getSITItemErr != nil {
		switch getSITItemErr {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipment.ID, "for MTO Service Item")
		default:
			return nil, getSITItemErr
		}
	}

	eTag := etag.GenerateEtag(SITItem.UpdatedAt)

	// Finally, update the mto_service_item with the members_expense flag set to TRUE
	SITItem.CustomerExpense = true
	mtoServiceItem, err := p.findServiceItem(appCtx, SITItem.ID)
	if err != nil {
		return &models.MTOServiceItem{}, err
	}

	return p.convertItemToCustomerExpense(appCtx, *mtoServiceItem, customerExpenseReason, convertToCustomerExpense, eTag, checkETag())
}

func (p *mtoServiceItemUpdater) findServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
	var serviceItem models.MTOServiceItem
	err := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder",
		"SITDestinationFinalAddress",
		"ReService",
		"SITOriginHHGOriginalAddress",
	).Find(&serviceItem, serviceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(serviceItemID, "while looking for service item")
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	return &serviceItem, nil
}

func (p *mtoServiceItemUpdater) approveOrRejectServiceItem(
	appCtx appcontext.AppContext,
	serviceItem models.MTOServiceItem,
	status models.MTOServiceItemStatus,
	rejectionReason *string,
	eTag string,
	checks ...validator,
) (*models.MTOServiceItem, error) {
	if verr := validateServiceItem(appCtx, &serviceItem, eTag, checks...); verr != nil {
		return nil, verr
	}

	var returnedServiceItem models.MTOServiceItem

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		updatedServiceItem, err := p.updateServiceItem(txnAppCtx, serviceItem, status, rejectionReason)
		if err != nil {
			return err
		}
		move := serviceItem.MoveTaskOrder

		if _, err = p.moveRouter.ApproveOrRequestApproval(txnAppCtx, move); err != nil {
			return err
		}

		returnedServiceItem = *updatedServiceItem

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedServiceItem, nil
}

func (p *mtoServiceItemUpdater) updateServiceItem(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem, status models.MTOServiceItemStatus, rejectionReason *string) (*models.MTOServiceItem, error) {
	serviceItem.Status = status
	now := time.Now()

	if status == models.MTOServiceItemStatusRejected {
		if rejectionReason == nil {
			verrs := validate.NewErrors()
			verrs.Add("rejectionReason", "field must be provided when status is set to REJECTED")
			err := apperror.NewInvalidInputError(serviceItem.ID, nil, verrs, "Invalid input found in the request.")
			return nil, err
		}

		serviceItem.RejectionReason = rejectionReason
		serviceItem.RejectedAt = &now
		// clear field if previously accepted
		serviceItem.ApprovedAt = nil
	} else if status == models.MTOServiceItemStatusApproved {
		// clear fields if previously rejected
		serviceItem.RejectionReason = nil
		serviceItem.RejectedAt = nil
		serviceItem.ApprovedAt = &now

		if serviceItem.MTOShipmentID != nil {
			// Get the shipment destination address
			mtoShipment, err := p.shipmentFetcher.GetShipment(appCtx, *serviceItem.MTOShipmentID, "DestinationAddress", "PickupAddress", "MTOServiceItems.SITOriginHHGOriginalAddress")
			if err != nil {
				return nil, err
			}

			// Check to see if there is already a SIT Destination Original Address
			// by checking for the ID before trying to set one on the service item.
			// If there isn't one, then we set it. We will update all four destination
			// SIT service items that get created
			if (serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
				serviceItem.ReService.Code == models.ReServiceCodeDDSFSC ||
				serviceItem.ReService.Code == models.ReServiceCodeDDASIT ||
				serviceItem.ReService.Code == models.ReServiceCodeDDFSIT) &&
				serviceItem.SITDestinationOriginalAddressID == nil {

				// Set the original address on a service item to the shipment's
				// destination address when approving destination SIT service items
				// Creating a new address record to ensure SITDestinationOriginalAddress
				// doesn't change if shipment destination address is updated
				shipmentDestinationAddress := &models.Address{
					StreetAddress1: mtoShipment.DestinationAddress.StreetAddress1,
					StreetAddress2: mtoShipment.DestinationAddress.StreetAddress2,
					StreetAddress3: mtoShipment.DestinationAddress.StreetAddress3,
					City:           mtoShipment.DestinationAddress.City,
					State:          mtoShipment.DestinationAddress.State,
					PostalCode:     mtoShipment.DestinationAddress.PostalCode,
					Country:        mtoShipment.DestinationAddress.Country,
				}
				shipmentDestinationAddress, err = p.addressCreator.CreateAddress(appCtx, shipmentDestinationAddress)
				if err != nil {
					return nil, err
				}
				serviceItem.SITDestinationOriginalAddressID = &shipmentDestinationAddress.ID
				serviceItem.SITDestinationOriginalAddress = shipmentDestinationAddress

				if serviceItem.SITDestinationFinalAddressID == nil {
					serviceItem.SITDestinationFinalAddressID = &shipmentDestinationAddress.ID
					serviceItem.SITDestinationFinalAddress = shipmentDestinationAddress
				}

				// Calculate SITDeliveryMiles for DDDSIT and DDSFSC origin SIT service items
				if serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
					serviceItem.ReService.Code == models.ReServiceCodeDDSFSC {
					// Destination SIT: distance between shipment destination address & service item ORIGINAL destination address
					milesCalculated, err := p.planner.ZipTransitDistance(appCtx, mtoShipment.DestinationAddress.PostalCode, serviceItem.SITDestinationOriginalAddress.PostalCode)
					if err != nil {
						return nil, err
					}
					serviceItem.SITDeliveryMiles = &milesCalculated
				}

			}
			// Calculate SITDeliveryMiles for DOPSIT and DOSFSC origin SIT service items
			if serviceItem.ReService.Code == models.ReServiceCodeDOPSIT ||
				serviceItem.ReService.Code == models.ReServiceCodeDOSFSC {
				// Origin SIT: distance between shipment pickup address & service item ORIGINAL pickup address
				milesCalculated, err := p.planner.ZipTransitDistance(appCtx, mtoShipment.PickupAddress.PostalCode, serviceItem.SITOriginHHGOriginalAddress.PostalCode)
				if err != nil {
					return nil, err
				}
				serviceItem.SITDeliveryMiles = &milesCalculated
			}
		}
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&serviceItem)
	if e := handleError(serviceItem.ID, verrs, err); e != nil {
		return nil, e
	}

	return &serviceItem, nil
}

func (p *mtoServiceItemUpdater) convertItemToCustomerExpense(
	appCtx appcontext.AppContext,
	serviceItem models.MTOServiceItem,
	customerExpenseReason *string,
	convertToCustomerExpense bool,
	eTag string,
	checks ...validator,
) (*models.MTOServiceItem, error) {
	if verr := validateServiceItem(appCtx, &serviceItem, eTag, checks...); verr != nil {
		return nil, verr
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		serviceItem.CustomerExpense = convertToCustomerExpense
		serviceItem.CustomerExpenseReason = customerExpenseReason
		verrs, err := appCtx.DB().ValidateAndUpdate(&serviceItem)
		e := handleError(serviceItem.ID, verrs, err)
		return e
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &serviceItem, nil
}

// UpdateMTOServiceItemBasic updates the MTO Service Item using base validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemBasic(
	appCtx appcontext.AppContext,
	mtoServiceItem *models.MTOServiceItem,
	eTag string,
) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(appCtx, mtoServiceItem, eTag, UpdateMTOServiceItemBasicValidator)
}

// UpdateMTOServiceItemPrime updates the MTO Service Item using Prime API validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemPrime(
	appCtx appcontext.AppContext,
	mtoServiceItem *models.MTOServiceItem,
	planner route.Planner,
	shipment models.MTOShipment,
	eTag string,
) (*models.MTOServiceItem, error) {
	updatedServiceItem, err := p.UpdateMTOServiceItem(appCtx, mtoServiceItem, eTag, UpdateMTOServiceItemPrimeValidator)

	if updatedServiceItem != nil {
		code := updatedServiceItem.ReService.Code

		// If this is an update to an Origin SIT or Destination SIT service item we need to recalculate the
		// Authorized End Date and Required Delivery Date
		if (code == models.ReServiceCodeDOASIT || code == models.ReServiceCodeDDASIT) &&
			updatedServiceItem.Status == models.MTOServiceItemStatusApproved {
			err = calculateSITDates(appCtx, mtoServiceItem, shipment, planner)
		}
	}

	return updatedServiceItem, err
}

// Calculate Required Delivery Date(RDD) from customer contact and requested delivery dates
// The RDD is calculated using the following business logic:
// If the SIT Departure Date is the same day or after the Customer Contact Date + GracePeriodDays then the RDD is Customer Contact Date + GracePeriodDays + GHC Transit Time
// If however the SIT Departure Date is before the Customer Contact Date + GracePeriodDays then the RDD is SIT Departure Date + GHC Transit Time
func calculateOriginSITRequiredDeliveryDate(appCtx appcontext.AppContext, shipment models.MTOShipment, planner route.Planner,
	sitCustomerContacted *time.Time, sitDepartureDate *time.Time) (*time.Time, error) {
	// Get a distance calculation between pickup and destination addresses.
	distance, err := planner.ZipTransitDistance(appCtx, shipment.PickupAddress.PostalCode, shipment.DestinationAddress.PostalCode)

	if err != nil {
		return nil, apperror.NewUnprocessableEntityError("cannot calculate distance between pickup and destination addresses")
	}

	weight := shipment.PrimeEstimatedWeight

	if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
		weight = shipment.NTSRecordedWeight
	}

	// Query the ghc_domestic_transit_times table for the max transit time using the distance between location
	// and the weight to determine the number of days for transit
	var ghcDomesticTransitTime models.GHCDomesticTransitTime
	err = appCtx.DB().Where("distance_miles_lower <= ? "+
		"AND distance_miles_upper >= ? "+
		"AND weight_lbs_lower <= ? "+
		"AND (weight_lbs_upper >= ? OR weight_lbs_upper = 0)",
		distance, distance, weight, weight).First(&ghcDomesticTransitTime)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipment.ID, fmt.Sprintf(
				"failed to find transit time for shipment of %d lbs weight and %d mile distance", weight.Int(), distance))
		default:
			return nil, apperror.NewQueryError("CalculateSITAllowanceRequestedDates", err, "failed to query for transit time")
		}
	}

	var requiredDeliveryDate time.Time
	customerContactDatePlusFive := sitCustomerContacted.AddDate(0, 0, GracePeriodDays)

	// we calculate required delivery date here using customer contact date and transit time
	if sitDepartureDate.Before(customerContactDatePlusFive) {
		requiredDeliveryDate = sitDepartureDate.AddDate(0, 0, ghcDomesticTransitTime.MaxDaysTransitTime)
	} else if sitDepartureDate.After(customerContactDatePlusFive) || sitDepartureDate.Equal(customerContactDatePlusFive) {
		requiredDeliveryDate = customerContactDatePlusFive.AddDate(0, 0, ghcDomesticTransitTime.MaxDaysTransitTime)
	}

	// Weekends and holidays are not allowable dates, find the next available workday
	var calendar = dates.NewUSCalendar()

	actual, observed, _ := calendar.IsHoliday(requiredDeliveryDate)

	if actual || observed || !calendar.IsWorkday(requiredDeliveryDate) {
		requiredDeliveryDate = dates.NextWorkday(*calendar, requiredDeliveryDate)
	}

	return &requiredDeliveryDate, nil
}

// Calculate the Required Delivery Date for the service item based on business logic using the
// Customer Contact Date, Customer Requested Delivery Date, and SIT Departure Date
func calculateSITDates(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, shipment models.MTOShipment,
	planner route.Planner) error {
	location := DestinationSITLocation

	if serviceItem.ReService.Code == models.ReServiceCodeDOASIT {
		location = OriginSITLocation
	}

	sitDepartureDate := serviceItem.SITDepartureDate

	if location == OriginSITLocation {

		if sitDepartureDate != nil {
			requiredDeliveryDate, err := calculateOriginSITRequiredDeliveryDate(appCtx, shipment, planner,
				serviceItem.SITCustomerContacted, sitDepartureDate)

			if err != nil {
				return err
			}

			shipment.RequiredDeliveryDate = requiredDeliveryDate
		} else {
			return apperror.NewNotFoundError(shipment.ID, "sit departure date not found, cannot update Required Delivery Date")
		}
	}

	// For Origin SIT we need to update the Required Delivery Date which is stored with the shipment instead of the service item
	if location == OriginSITLocation {
		var verrs *validate.Errors

		verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(shipment.ID, err, verrs, "invalid input found while updating dates of shipment")
		} else if err != nil {
			return apperror.NewQueryError("Shipment", err, "")
		}
	}

	return nil
}

// UpdateMTOServiceItem updates the given service item
func (p *mtoServiceItemUpdater) UpdateMTOServiceItem(
	appCtx appcontext.AppContext,
	mtoServiceItem *models.MTOServiceItem,
	eTag string,
	validatorKey string,
) (*models.MTOServiceItem, error) {
	// Find the service item, return error if not found
	oldServiceItem, err := models.FetchServiceItem(appCtx.DB(), mtoServiceItem.ID)
	if err != nil {
		switch err {
		case models.ErrFetchNotFound:
			return nil, apperror.NewNotFoundError(mtoServiceItem.ID, "while looking for MTOServiceItem")
		default:
			return nil, apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	checker := movetaskorder.NewMoveTaskOrderChecker()
	serviceItemData := updateMTOServiceItemData{
		updatedServiceItem:  *mtoServiceItem,
		oldServiceItem:      oldServiceItem,
		availabilityChecker: checker,
		verrs:               validate.NewErrors(),
	}

	validServiceItem, err := ValidateUpdateMTOServiceItem(appCtx, &serviceItemData, validatorKey)
	if err != nil {
		return nil, err
	}

	// If we have any Customer Contacts we need to make sure that they are associated with
	// all related destination SIT service items. This is especially important if we are creating new Customer Contacts.
	if len(validServiceItem.CustomerContacts) > 0 {
		relatedServiceItems, fetchErr := models.FetchRelatedDestinationSITServiceItems(appCtx.DB(), validServiceItem.ID)
		if fetchErr != nil {
			return nil, fetchErr
		}
		for i := range validServiceItem.CustomerContacts {
			validServiceItem.CustomerContacts[i].MTOServiceItems = relatedServiceItems
		}
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldServiceItem.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(validServiceItem.ID, nil)
	}

	// Create address record (if needed) and update service item in a single transaction
	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if validServiceItem.SITDestinationFinalAddress != nil {
			if validServiceItem.SITDestinationFinalAddressID == nil || *validServiceItem.SITDestinationFinalAddressID == uuid.Nil {
				validServiceItem.SITDestinationFinalAddress, err = p.addressCreator.CreateAddress(txnAppCtx, validServiceItem.SITDestinationFinalAddress)
				if err != nil {
					// If the error is something else (this is unexpected), we create a QueryError
					return apperror.NewQueryError("MTOServiceItem", err, "")
				}
			}
			validServiceItem.SITDestinationFinalAddressID = &validServiceItem.SITDestinationFinalAddress.ID
		}
		for index := range validServiceItem.CustomerContacts {
			validCustomerContact := &validServiceItem.CustomerContacts[index]
			if validCustomerContact.ID == uuid.Nil {
				verrs, createErr := p.builder.CreateOne(txnAppCtx, validCustomerContact)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(
						validServiceItem.ID, createErr, verrs, "Invalid input found while creating a Customer Contact for service item.")
				} else if createErr != nil {
					// If the error is something else (this is unexpected), we create a QueryError
					return apperror.NewQueryError("MTOServiceItem", createErr, "")
				}
			} else {
				verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(validCustomerContact)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(validServiceItem.ID, updateErr, verrs, "Invalid input found while updating customer contact for the service item.")
				} else if updateErr != nil {
					// If the error is something else (this is unexpected), we create a QueryError
					return apperror.NewQueryError("MTOServiceItem", updateErr, "")
				}
			}
		}

		// Make the update and create a InvalidInputError if there were validation issues
		verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(validServiceItem)

		// If there were validation errors create an InvalidInputError type
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(validServiceItem.ID, updateErr, verrs, "Invalid input found while updating the service item.")
		} else if updateErr != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("MTOServiceItem", updateErr, "")
		}
		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return validServiceItem, nil
}

// ValidateUpdateMTOServiceItem checks the provided serviceItemData struct against the validator indicated by validatorKey.
// Defaults to base validation if the empty string is entered as the key.
// Returns an MTOServiceItem that has been set up for update.
func ValidateUpdateMTOServiceItem(appCtx appcontext.AppContext, serviceItemData *updateMTOServiceItemData, validatorKey string) (*models.MTOServiceItem, error) {
	if validatorKey == "" {
		validatorKey = UpdateMTOServiceItemBasicValidator
	}
	validator, ok := UpdateMTOServiceItemValidators[validatorKey]
	if !ok {
		err := fmt.Errorf("validator key %s was not found in update MTO Service Item validators", validatorKey)
		return nil, err
	}
	err := validator.validate(appCtx, serviceItemData)
	if err != nil {
		return nil, err
	}

	newServiceItem := serviceItemData.setNewMTOServiceItem()

	return newServiceItem, nil
}
