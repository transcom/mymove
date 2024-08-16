package sitstatus

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
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
	pastSITs    models.SITServiceItemGroupings
	currentSITs models.SITServiceItemGroupings
	futureSITs  models.SITServiceItemGroupings
}

func newSortedShipmentSITs() SortedShipmentSITs {
	return SortedShipmentSITs{
		pastSITs:    make([]models.SITServiceItemGrouping, 0),
		currentSITs: make([]models.SITServiceItemGrouping, 0),
		futureSITs:  make([]models.SITServiceItemGrouping, 0),
	}
}

// Sort the Shipment's SIT groupings by their summary into either past, current, or future SIT groupings
func SortShipmentSITs(sitGroupings models.SITServiceItemGroupings, today time.Time) SortedShipmentSITs {
	shipmentSITs := newSortedShipmentSITs()
	for _, sitGrouping := range sitGroupings {
		if sitGrouping.Summary.SITEntryDate.After(today) {
			shipmentSITs.futureSITs = append(shipmentSITs.futureSITs, sitGrouping)
		} else if sitGrouping.Summary.SITDepartureDate != nil && sitGrouping.Summary.SITDepartureDate.Before(today) {
			shipmentSITs.pastSITs = append(shipmentSITs.pastSITs, sitGrouping)
		} else {
			shipmentSITs.currentSITs = append(shipmentSITs.currentSITs, sitGrouping)
		}
	}
	return shipmentSITs
}

func Clamp(input, min, max int) (int, error) {
	result := input
	if input < min {
		result = min
	} else if input > max {
		result = max
	}
	if result < min || result > max {
		return result, errors.New("Clamp input is out of scope")
	}
	return result, nil
}

// Retrieve the SIT service item groupings for the provided shipment
// Each SIT grouping has a top-level summary of the grouped SIT
func (f shipmentSITStatus) RetrieveShipmentSIT(appCtx appcontext.AppContext, shipment models.MTOShipment) models.SITServiceItemGroupings {
	var shipmentSITs models.SITServiceItemGroupings

	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// Group SITs based on their entry date
	// By using the SIT entry date to group SIT service items together, we can create the support
	// of multiple Origin/Destination SITs on a shipment. (And eventually partial SITs upon further enhancement)
	// Partial SITs are not yet supported.
	groupedSITs := map[time.Time]*models.SITServiceItemGrouping{} // This creates a map of groupings based on a provided time (Entry date)

	for _, serviceItem := range shipment.MTOServiceItems {
		if serviceItem.Status != models.MTOServiceItemStatusApproved {
			continue // Don't group service items that have not been approved
		}

		entryDate := serviceItem.SITEntryDate
		if entryDate == nil {
			continue // Don't group service items that do not have a SIT Entry date
		}

		// Check if a group for the SIT entry date already exists
		if group, exists := groupedSITs[*entryDate]; exists {
			// Append if it one exists
			group.ServiceItems = append(group.ServiceItems, serviceItem)
		} else {
			// Create a new group for this entry date
			location := OriginSITLocation
			if containsReServiceCode(models.ValidDomesticDestinationSITReServiceCodes, serviceItem.ReService.Code) {
				location = DestinationSITLocation
			}
			newGroup := &models.SITServiceItemGrouping{
				Summary:      models.SITSummary{Location: location},
				ServiceItems: []models.MTOServiceItem{serviceItem},
			}
			groupedSITs[*entryDate] = newGroup
		}
	}

	// Generate summaries for each group and append them to shipmentSITs
	for _, group := range groupedSITs {
		summary := f.generateSITSummary(*group, today)
		if summary != nil {
			group.Summary = *summary
			shipmentSITs = append(shipmentSITs, *group)
		}
	}

	return shipmentSITs
}

// Helper function to take in an MTO service item's ReServiceCode and validate it
// against a given array of codes. This is primarily to support the RetrieveShipmentSIT method
// when SIT groupings are created.
func containsReServiceCode(validCodes []models.ReServiceCode, code models.ReServiceCode) bool {
	for _, validCode := range validCodes {
		if validCode == code {
			return true
		}
	}
	return false
}

// Helper function to generate the SIT Summary for a group of service items
func (f shipmentSITStatus) generateSITSummary(sit models.SITServiceItemGrouping, today time.Time) *models.SITSummary {
	if sit.ServiceItems == nil {
		// Return nil if there are no service items
		return nil
	}
	// This is where the craziest part of the code should ever be (Besides the grouping section)
	// Due to our service item architecture, SIT is split across many service items
	// and due to existing handlers and service objects, it's possible these SIT service items
	// will have discrepancies and information spread across multiple items.
	// This SIT summary is to make it readable down the line, and handle all complex calculations
	// in one, central location.
	var earliestSITEntryDate *time.Time
	var earliestSITDepartureDate *time.Time
	var earliestSITAuthorizedEndDate *time.Time
	var earliestSITCustomerContacted *time.Time
	var earliestSITRequestedDelivery *time.Time
	var calculatedTotalDaysInSIT *int
	var location string
	var firstDaySITServiceItemID uuid.UUID

	// Begin filling in the holes so we can generate a SIT summary
	for _, sitServiceItem := range sit.ServiceItems {
		// Grab the first location found
		if location == "" {
			if containsReServiceCode(models.ValidDomesticOriginSITReServiceCodes, sitServiceItem.ReService.Code) {
				// Set to Domestic Origin
				location = OriginSITLocation
			}
			if containsReServiceCode(models.ValidDomesticDestinationSITReServiceCodes, sitServiceItem.ReService.Code) {
				// Set to Domestic Destination
				location = DestinationSITLocation
			}
		}

		// Grab the first day SIT service item ID for payment requests
		// TODO: Eventually refactor this out and use the entire group for the payment request
		if (sitServiceItem.ReService.Code == models.ReServiceCodeDOFSIT || sitServiceItem.ReService.Code == models.ReServiceCodeDOASIT) && firstDaySITServiceItemID == uuid.Nil {
			firstDaySITServiceItemID = sitServiceItem.ID
		}

		// Grab the earliest SIT entry date (Granted they should always be the same for a group)
		if earliestSITEntryDate == nil || (sitServiceItem.SITEntryDate != nil && sitServiceItem.SITEntryDate.Before(*earliestSITEntryDate)) {
			earliestSITEntryDate = sitServiceItem.SITEntryDate
		}

		// Grab the earliest SIT Departure Date
		if earliestSITDepartureDate == nil || (sitServiceItem.SITDepartureDate != nil && sitServiceItem.SITDepartureDate.Before(*earliestSITDepartureDate)) {
			earliestSITDepartureDate = sitServiceItem.SITDepartureDate
		}

		// Grab the earliest SIT Authorized End Date
		// based off of the provided earliest SIT entry date
		// retrieving the authorized end date requires a SIT entry date
		if earliestSITAuthorizedEndDate == nil && earliestSITEntryDate != nil {
			daysInSIT := daysInSIT(*earliestSITEntryDate, earliestSITDepartureDate, today)
			calculatedTotalDaysInSIT = &daysInSIT
			earliestSITAuthorizedEndDateValue := CalculateSITAuthorizedEndDate(len(sit.ServiceItems), daysInSIT, *earliestSITEntryDate, *calculatedTotalDaysInSIT)
			earliestSITAuthorizedEndDate = &earliestSITAuthorizedEndDateValue
		}

		// Grab the first Customer Contacted
		if earliestSITCustomerContacted == nil && sitServiceItem.SITCustomerContacted != nil {
			earliestSITCustomerContacted = sitServiceItem.SITCustomerContacted
		}

		// Grab the first Requested Delivery
		if earliestSITRequestedDelivery == nil && sitServiceItem.SITRequestedDelivery != nil {
			earliestSITRequestedDelivery = sitServiceItem.SITRequestedDelivery
		}
	}

	return &models.SITSummary{
		FirstDaySITServiceItemID: firstDaySITServiceItemID,
		Location:                 location,
		DaysInSIT:                *calculatedTotalDaysInSIT, // FIXME: This appears to not be calculating properly at the summary level
		SITEntryDate:             *earliestSITEntryDate,
		SITDepartureDate:         earliestSITDepartureDate,
		SITAuthorizedEndDate:     *earliestSITAuthorizedEndDate,
		SITCustomerContacted:     earliestSITCustomerContacted,
		SITRequestedDelivery:     earliestSITRequestedDelivery,
	}
}

// CalculateShipmentSITStatus creates a SIT Status for payload to be used in
// multiple handlers in the `ghcapi` package for the MTOShipment handlers.
func (f shipmentSITStatus) CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) (*services.SITStatus, models.MTOShipment, error) {
	if shipment.MTOServiceItems == nil || len(shipment.MTOServiceItems) == 0 {
		return nil, shipment, nil
	}

	var shipmentSITStatus services.SITStatus

	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	sitGroupings := f.RetrieveShipmentSIT(appCtx, shipment)

	// Sort the SIT groupings into past, current, future
	shipmentSITGroupings := SortShipmentSITs(sitGroupings, today)

	// From our SIT groupings, grab the current one
	// Current SIT can only be either Current or Future at this time
	currentSIT := getCurrentSIT(shipmentSITGroupings)

	// There were no relevant SIT service items for this shipment
	if currentSIT == nil && len(shipmentSITGroupings.pastSITs) == 0 {
		return nil, shipment, nil
	}

	shipmentSITStatus.ShipmentID = shipment.ID
	totalSITAllowance, err := f.CalculateShipmentSITAllowance(appCtx, shipment)
	if err != nil {
		return nil, shipment, err
	}
	totalSITDaysUsedClampedResult, totalDaysUsedErr := Clamp(CalculateTotalDaysInSIT(shipmentSITGroupings, today), 0, totalSITAllowance)
	if totalDaysUsedErr != nil {
		return nil, shipment, err
	}
	shipmentSITStatus.TotalSITDaysUsed = totalSITDaysUsedClampedResult
	shipmentSITStatus.CalculatedTotalDaysInSIT = CalculateTotalDaysInSIT(shipmentSITGroupings, today)
	shipmentSITStatus.TotalDaysRemaining = totalSITAllowance - shipmentSITStatus.TotalSITDaysUsed
	shipmentSITStatus.PastSITs = shipmentSITGroupings.pastSITs

	if currentSIT != nil {
		location := currentSIT.Summary.Location
		firstDaySITServiceItemID := currentSIT.Summary.FirstDaySITServiceItemID
		daysInSIT := daysInSIT(currentSIT.Summary.SITEntryDate, currentSIT.Summary.SITDepartureDate, today)
		sitEntryDate := currentSIT.Summary.SITEntryDate
		sitDepartureDate := currentSIT.Summary.SITDepartureDate
		sitAuthorizedEndDate := CalculateSITAuthorizedEndDate(totalSITAllowance, daysInSIT, sitEntryDate, shipmentSITStatus.CalculatedTotalDaysInSIT)
		var sitCustomerContacted, sitRequestedDelivery *time.Time
		sitCustomerContacted = currentSIT.Summary.SITCustomerContacted
		sitRequestedDelivery = currentSIT.Summary.SITRequestedDelivery

		shipmentSITStatus.CurrentSIT = &services.CurrentSIT{
			ServiceItemID:        firstDaySITServiceItemID,
			Location:             location,
			DaysInSIT:            daysInSIT,
			SITEntryDate:         sitEntryDate,
			SITDepartureDate:     sitDepartureDate,
			SITAuthorizedEndDate: sitAuthorizedEndDate,
			SITCustomerContacted: sitCustomerContacted,
			SITRequestedDelivery: sitRequestedDelivery,
		}

		// update the shipment's OriginSITAuthEndDate or DestinationSITAuthEndDate depending on what currentSIT location is
		if shipmentSITStatus.CurrentSIT != nil {
			if location == OriginSITLocation {
				shipment.OriginSITAuthEndDate = &shipmentSITStatus.CurrentSIT.SITAuthorizedEndDate
			} else {
				shipment.DestinationSITAuthEndDate = &shipmentSITStatus.CurrentSIT.SITAuthorizedEndDate
			}
		}
	}
	return &shipmentSITStatus, shipment, nil
}

/*
Helper function that takes in the shipment's SIT groupings and returns the group which
enters SIT on the earliest date based on their SIT Summary
*/
func getEarliestSIT(sitGroupings models.SITServiceItemGroupings) *models.SITServiceItemGrouping {
	if len(sitGroupings) == 0 {
		return nil
	}
	earliest := sitGroupings[0]
	for _, sit := range sitGroupings {
		if sit.Summary.SITEntryDate.Before(earliest.Summary.SITEntryDate) {
			earliest = sit
		}
	}
	return &earliest
}

/*
Private function that returns the most relevant current or upcoming SIT.
SIT service items that have already started are prioritized, followed by SIT
service items that start in the future.
*/
func getCurrentSIT(shipmentSITs SortedShipmentSITs) *models.SITServiceItemGrouping {
	if len(shipmentSITs.currentSITs) > 0 {
		return getEarliestSIT(shipmentSITs.currentSITs)
	} else if len(shipmentSITs.futureSITs) > 0 {
		return getEarliestSIT(shipmentSITs.futureSITs)
	} /* else if len(shipmentSITs.pastSITs) > 0 {
		// TODO: Enhance
		// This is a temporary check of to return the earliest
		// past SIT if there are no current or future SITs to choose from.
		// This is done because at this time the UI can only handle 1 'active' SIT at a time
		// and the customer has deemed that if a new SIT has not been implemented, then we still want
		// to display the old SIT as the 'current' SIT, even though it's in the past.
		return getEarliestSIT(shipmentSITs.pastSITs)
	} */
	return nil
}

// Private function daysInSIT is used to calculate the number of days an item
// is in SIT using a serviceItem and the current day.
//
// If the service item has a departure date and SIT entry date is in the past,
// then the return value is the SITDepartureDate - SITEntryDate in hours, then converted to days.
//
// If there is no departure date and the SIT entry date in the past, then the
// return value is Today - SITEntryDate, adding 1 to include today.
func daysInSIT(sitEntryDate time.Time, sitDepartureDate *time.Time, today time.Time) int {
	var days int
	if sitDepartureDate != nil && sitDepartureDate.Before(today) {
		days = int(sitDepartureDate.Sub(sitEntryDate).Hours()) / 24
	} else if sitEntryDate.Before(today) || sitEntryDate.Equal(today) {
		days = int(today.Sub(sitEntryDate).Hours())/24 + 1
	}
	return days
}

func CalculateTotalDaysInSIT(shipmentSITs SortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, pastSIT := range shipmentSITs.pastSITs {
		totalDays += daysInSIT(pastSIT.Summary.SITEntryDate, pastSIT.Summary.SITDepartureDate, today)
	}
	for _, currentSIT := range shipmentSITs.currentSITs {
		totalDays += daysInSIT(currentSIT.Summary.SITEntryDate, currentSIT.Summary.SITDepartureDate, today)
	}
	return totalDays
}

// adds up all the days from pastSITs
func CalculateTotalPastDaysInSIT(shipmentSITs SortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, pastSIT := range shipmentSITs.pastSITs {
		totalDays += daysInSIT(pastSIT.Summary.SITEntryDate, pastSIT.Summary.SITDepartureDate, today)
	}
	return totalDays
}

func CalculateSITAuthorizedEndDate(totalSITAllowance int, currentDaysInSIT int, sitEntryDate time.Time, calculatedTotalDaysInSIT int) time.Time {
	return sitEntryDate.AddDate(0, 0, (totalSITAllowance - (calculatedTotalDaysInSIT - currentDaysInSIT)))
}

func (f shipmentSITStatus) CalculateShipmentsSITStatuses(appCtx appcontext.AppContext, shipments []models.MTOShipment) map[string]services.SITStatus {
	shipmentsSITStatuses := map[string]services.SITStatus{}

	for _, shipment := range shipments {
		shipmentSITStatus, _, _ := f.CalculateShipmentSITStatus(appCtx, shipment)
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
