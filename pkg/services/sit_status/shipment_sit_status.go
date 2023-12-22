package sitstatus

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

// OriginSITLocation is the constant representing when the shipment in storage occurs at the origin
const OriginSITLocation = "ORIGIN"

// DestinationSITLocation is the constant representing when the shipment in storage occurs at the destination
const DestinationSITLocation = "DESTINATION"

// Number of days of grace period after customer contacts prime for delivery out of SIT
const GracePeriodDays = 5

type shipmentSITStatus struct {
}

// NewShipmentSITStatus creates a new instance of the service object that implements calculating a shipments SIT summary
func NewShipmentSITStatus() services.ShipmentSITStatus {
	return &shipmentSITStatus{}
}

type SortedShipmentSITs struct {
	pastSITs    []models.MTOServiceItem
	currentSITs []models.MTOServiceItem
	futureSITs  []models.MTOServiceItem
}

func newSortedShipmentSITs() SortedShipmentSITs {
	return SortedShipmentSITs{
		pastSITs:    make([]models.MTOServiceItem, 0),
		currentSITs: make([]models.MTOServiceItem, 0),
		futureSITs:  make([]models.MTOServiceItem, 0),
	}
}

func SortShipmentSITs(shipment models.MTOShipment, today time.Time) SortedShipmentSITs {
	shipmentSITs := newSortedShipmentSITs()

	// TODO change service codes here to DOFSIT & DDFSIT and see what breaks
	for _, serviceItem := range shipment.MTOServiceItems {
		// only departure SIT service items have a departure date
		if code := serviceItem.ReService.Code; (code == models.ReServiceCodeDOFSIT || code == models.ReServiceCodeDDFSIT) &&
			serviceItem.Status == models.MTOServiceItemStatusApproved {
			if serviceItem.SITEntryDate.After(today) {
				shipmentSITs.futureSITs = append(shipmentSITs.futureSITs, serviceItem)
			} else if serviceItem.SITDepartureDate != nil && serviceItem.SITDepartureDate.Before(today) {
				shipmentSITs.pastSITs = append(shipmentSITs.pastSITs, serviceItem)
			} else {
				shipmentSITs.currentSITs = append(shipmentSITs.currentSITs, serviceItem)
			}
		}
	}
	return shipmentSITs
}

func clamp(input, min, max int) int {
	if input < min {
		return min
	} else if input > max {
		return max
	}
	return input
}

// CalculateShipmentSITStatus creates a SIT Status for payload to be used in
// multiple handlers in the `ghcapi` package for the MTOShipment handlers.
func (f shipmentSITStatus) CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) (*services.SITStatus, error) {
	if shipment.MTOServiceItems == nil || len(shipment.MTOServiceItems) == 0 {
		return nil, nil
	}

	var shipmentSITStatus services.SITStatus

	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	shipmentSITs := SortShipmentSITs(shipment, today)

	currentSIT := getCurrentSIT(shipmentSITs)

	// There were no relevant SIT service items for this shipment
	if currentSIT == nil && len(shipmentSITs.pastSITs) == 0 {
		return nil, nil
	}

	shipmentSITStatus.ShipmentID = shipment.ID
	totalSITAllowance, err := f.CalculateShipmentSITAllowance(appCtx, shipment)
	if err != nil {
		return nil, err
	}
	shipmentSITStatus.TotalSITDaysUsed = clamp(CalculateTotalDaysInSIT(shipmentSITs, today), 0, totalSITAllowance)
	shipmentSITStatus.CalculatedTotalDaysInSIT = CalculateTotalDaysInSIT(shipmentSITs, today)
	shipmentSITStatus.TotalDaysRemaining = totalSITAllowance - shipmentSITStatus.TotalSITDaysUsed
	shipmentSITStatus.PastSITs = shipmentSITs.pastSITs

	if currentSIT != nil {
		location := DestinationSITLocation
		if currentSIT.ReService.Code == models.ReServiceCodeDOFSIT {
			location = OriginSITLocation
		}
		daysInSIT := daysInSIT(*currentSIT, today)
		sitEntryDate := *currentSIT.SITEntryDate
		sitDepartureDate := currentSIT.SITDepartureDate
		sitAllowanceEndDate := CalculateSITAllowanceEndDate(shipmentSITStatus.TotalDaysRemaining, sitEntryDate, today)
		var sitCustomerContacted, sitRequestedDelivery *time.Time
		sitCustomerContacted = currentSIT.SITCustomerContacted
		sitRequestedDelivery = currentSIT.SITRequestedDelivery

		shipmentSITStatus.CurrentSIT = &services.CurrentSIT{
			ServiceItemID:        currentSIT.ID,
			Location:             location,
			DaysInSIT:            daysInSIT,
			SITEntryDate:         sitEntryDate,
			SITDepartureDate:     sitDepartureDate,
			SITAllowanceEndDate:  sitAllowanceEndDate,
			SITCustomerContacted: sitCustomerContacted,
			SITRequestedDelivery: sitRequestedDelivery,
		}
	}

	return &shipmentSITStatus, nil
}

/*
Private function that takes a list of sitServiceItems and returns the one
that enters SIT on the earliest date
*/
func getEarliestSIT(sitServiceItems []models.MTOServiceItem) *models.MTOServiceItem {
	if len(sitServiceItems) == 0 {
		return nil
	}
	earliest := sitServiceItems[0]
	for _, serviceItem := range sitServiceItems {
		if serviceItem.SITEntryDate.Before(*earliest.SITEntryDate) {
			earliest = serviceItem
		}
	}
	return &earliest
}

/*
Private function that returns the most relevant current or upcoming SIT.
SIT service items that have already started are prioritized, followed by SIT
service items that start in the future.
*/
func getCurrentSIT(shipmentSITs SortedShipmentSITs) *models.MTOServiceItem {
	if len(shipmentSITs.currentSITs) > 0 {
		return getEarliestSIT(shipmentSITs.currentSITs)
	} else if len(shipmentSITs.futureSITs) > 0 {
		return getEarliestSIT(shipmentSITs.futureSITs)
	}
	return nil
}

// Private function daysInSIT is used to calculate the number of days an item
// is in SIT using a serviceItem and the current day.
//
// If the service item has a departure date and SIT entry date is in the past,
// then the return value is the SITDepartureDate - SITEntryDate.
//
// If there is no departure date and the SIT entry date in the past, then the
// return value is Today - SITEntryDate.
func daysInSIT(serviceItem models.MTOServiceItem, today time.Time) int {
	if serviceItem.SITDepartureDate != nil && serviceItem.SITDepartureDate.Before(today) {
		return int(serviceItem.SITDepartureDate.Sub(*serviceItem.SITEntryDate).Hours()) / 24
	} else if serviceItem.SITEntryDate.Before(today) {
		return int(today.Sub(*serviceItem.SITEntryDate).Hours()) / 24
	}

	return 0
}

func CalculateTotalDaysInSIT(shipmentSITs SortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, serviceItem := range shipmentSITs.pastSITs {
		totalDays += daysInSIT(serviceItem, today)
	}
	for _, serviceItem := range shipmentSITs.currentSITs {
		totalDays += daysInSIT(serviceItem, today)
	}
	return totalDays
}

func CalculateSITAllowanceEndDate(totalDaysRemaining int, sitEntryDate time.Time, today time.Time) time.Time {
	//current SIT
	if sitEntryDate.Before(today) {
		return today.AddDate(0, 0, totalDaysRemaining)
	}
	// future SIT
	return sitEntryDate.AddDate(0, 0, totalDaysRemaining)
}

func (f shipmentSITStatus) CalculateShipmentsSITStatuses(appCtx appcontext.AppContext, shipments []models.MTOShipment) map[string]services.SITStatus {
	shipmentsSITStatuses := map[string]services.SITStatus{}

	for _, shipment := range shipments {
		shipmentSITStatus, _ := f.CalculateShipmentSITStatus(appCtx, shipment)
		if shipmentSITStatus != nil {
			shipmentsSITStatuses[shipment.ID.String()] = *shipmentSITStatus
		}
	}

	return shipmentsSITStatuses
}

// CalculateShipmentSITAllowance finds the number of days allowed in SIT for a shipment based on its entitlement and any approved SIT extensions
func (f shipmentSITStatus) CalculateShipmentSITAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment) (int, error) {
	entitlement, err := fetchEntitlement(appCtx, shipment)
	if err != nil {
		return 0, apperror.NewNotFoundError(shipment.ID, "shipment is missing entitlement")
	}

	totalSITAllowance := 0
	if entitlement.StorageInTransit != nil {
		totalSITAllowance = *entitlement.StorageInTransit
	}
	for _, ext := range shipment.SITDurationUpdates {
		if ext.ApprovedDays != nil {
			totalSITAllowance += *ext.ApprovedDays
		}
	}
	return totalSITAllowance, nil
}

func fetchEntitlement(appCtx appcontext.AppContext, mtoShipment models.MTOShipment) (*models.Entitlement, error) {
	var move models.Move
	err := appCtx.DB().Q().EagerPreload("Orders.Entitlement").Find(&move, mtoShipment.MoveTaskOrderID)

	if err != nil {
		return nil, err
	}

	return move.Orders.Entitlement, nil
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

func (f shipmentSITStatus) CalculateSITAllowanceRequestedDates(appCtx appcontext.AppContext, shipment models.MTOShipment, planner route.Planner,
	sitCustomerContacted *time.Time, sitRequestedDelivery *time.Time, eTag string) (*services.SITStatus, error) {
	existingETag := etag.GenerateEtag(shipment.UpdatedAt)

	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(shipment.ID, errors.New("the if-match header value did not match the etag for this record"))
	}

	if shipment.MTOServiceItems == nil || len(shipment.MTOServiceItems) == 0 {
		return nil, apperror.NewNotFoundError(shipment.ID, "shipment is missing MTO Service Items")
	}

	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	shipmentSITs := SortShipmentSITs(shipment, today)

	currentSIT := getCurrentSIT(shipmentSITs)

	// There were no relevant SIT service items for this shipment
	if currentSIT == nil {
		return nil, apperror.NewNotFoundError(shipment.ID, "shipment is missing current SIT")
	}
	var shipmentSITStatus services.SITStatus
	currentSIT.SITCustomerContacted = sitCustomerContacted
	currentSIT.SITRequestedDelivery = sitRequestedDelivery
	shipmentSITStatus.ShipmentID = shipment.ID
	location := DestinationSITLocation

	if currentSIT.ReService.Code == models.ReServiceCodeDOFSIT {
		location = OriginSITLocation
	}

	daysInSIT := daysInSIT(*currentSIT, today)
	sitEntryDate := *currentSIT.SITEntryDate
	sitDepartureDate := currentSIT.SITDepartureDate

	// Calculate sitAllowanceEndDate and required delivery date based on sitCustomerContacted and sitRequestedDelivery
	// using the below business logic.
	sitAllowanceEndDate := sitDepartureDate

	if location == OriginSITLocation {
		// Origin SIT: sitAllowanceEndDate should be GracePeriodDays days after sitCustomerContacted or the sitDepartureDate whichever is earlier.
		calculatedAllowanceEndDate := sitCustomerContacted.AddDate(0, 0, GracePeriodDays)

		if sitDepartureDate == nil || calculatedAllowanceEndDate.Before(*sitDepartureDate) {
			sitAllowanceEndDate = &calculatedAllowanceEndDate
		}

		if sitDepartureDate != nil {
			requiredDeliveryDate, err := calculateOriginSITRequiredDeliveryDate(appCtx, shipment, planner, sitCustomerContacted, sitDepartureDate)

			if err != nil {
				return nil, err
			}

			shipment.RequiredDeliveryDate = requiredDeliveryDate
		} else {
			return nil, apperror.NewNotFoundError(shipment.ID, "sit departure date not found")
		}

	} else if location == DestinationSITLocation {
		// Destination SIT: sitAllowanceEndDate should be GracePeriodDays days after sitRequestedDelivery or the sitDepartureDate whichever is earlier.
		calculatedAllowanceEndDate := sitRequestedDelivery.AddDate(0, 0, GracePeriodDays)

		if sitDepartureDate == nil || calculatedAllowanceEndDate.Before(*sitDepartureDate) {
			sitAllowanceEndDate = &calculatedAllowanceEndDate
		}
	}

	shipmentSITStatus.CurrentSIT = &services.CurrentSIT{
		Location:             location,
		DaysInSIT:            daysInSIT,
		SITEntryDate:         sitEntryDate,
		SITDepartureDate:     sitDepartureDate,
		SITAllowanceEndDate:  *sitAllowanceEndDate,
		SITCustomerContacted: sitCustomerContacted,
		SITRequestedDelivery: sitRequestedDelivery,
	}

	var verrs *validate.Errors
	var err error

	if location == OriginSITLocation {
		verrs, err = appCtx.DB().ValidateAndUpdate(&shipment)

		if verrs != nil && verrs.HasAny() {
			return nil, apperror.NewInvalidInputError(shipment.ID, err, verrs, "invalid input found while updating dates of shipment")
		} else if err != nil {
			return nil, apperror.NewQueryError("Shipment", err, "")
		}
	}

	verrs, err = appCtx.DB().ValidateAndUpdate(currentSIT)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(currentSIT.ID, err, verrs, "invalid input found while updating current sit service item")
	} else if err != nil {
		return nil, apperror.NewQueryError("Service item", err, "")
	}

	return &shipmentSITStatus, nil
}
