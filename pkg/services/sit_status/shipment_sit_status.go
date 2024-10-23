package sitstatus

import (
	"sort"
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
	PastSITs    models.SITServiceItemGroupings
	CurrentSITs models.SITServiceItemGroupings // Takes an array but at this time only a single CurrentSIT is supported. This could potentially be used for partial delivery current SITs
	FutureSITs  models.SITServiceItemGroupings
}

func newSortedShipmentSITs() SortedShipmentSITs {
	return SortedShipmentSITs{
		PastSITs:    make([]models.SITServiceItemGrouping, 0),
		CurrentSITs: make([]models.SITServiceItemGrouping, 0),
		FutureSITs:  make([]models.SITServiceItemGrouping, 0),
	}
}

// Sort the Shipment's SIT groupings by their summary into either past, current, or future SIT groupings
func SortShipmentSITs(sitGroupings models.SITServiceItemGroupings, today time.Time) SortedShipmentSITs {
	shipmentSITs := newSortedShipmentSITs()
	for _, sitGrouping := range sitGroupings {
		if sitGrouping.Summary.SITEntryDate.After(today) {
			shipmentSITs.FutureSITs = append(shipmentSITs.FutureSITs, sitGrouping)
		} else if sitGrouping.Summary.SITDepartureDate != nil && sitGrouping.Summary.SITDepartureDate.Before(today) {
			shipmentSITs.PastSITs = append(shipmentSITs.PastSITs, sitGrouping)
		} else {
			shipmentSITs.CurrentSITs = append(shipmentSITs.CurrentSITs, sitGrouping)
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
func (f shipmentSITStatus) RetrieveShipmentSIT(appCtx appcontext.AppContext, shipment models.MTOShipment) (models.SITServiceItemGroupings, error) {
	var shipmentSITs models.SITServiceItemGroupings
	var daysUsed int // Sum of all days used, this will increase after each SIT summary is generated for each group
	// It is passed to the next grouping (sorted by entry date) in order to track proper authorized end dates
	// by subtracting the days used from the allowance
	var mostRecentSITAuthorizedEndDate *time.Time

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

	// Get the total SIT allowance for this shipment
	totalSITAllowance, err := f.CalculateShipmentSITAllowance(appCtx, shipment)
	if err != nil {
		return nil, err
	}

	// Sort by entry date
	var entryDates []time.Time
	for entryDate := range groupedSITs {
		// Extract map keys and insert into a slice
		// This is done due to flaky issues that come from
		// a map being unordered.
		// By using a slice, we can ensure a proper sort
		entryDates = append(entryDates, entryDate)
	}
	sort.Slice(entryDates, func(i, j int) bool {
		return entryDates[i].Before(entryDates[j])
	})

	// Generate summaries for each group and append them to shipmentSITs
	// Use the sorted entry date keys
	for _, entryDate := range entryDates {
		// We know all groups will be properly iterated over because
		// entryDates are the extracted map keys from teh groupedSITs map
		group := groupedSITs[entryDate]
		summary := f.generateSITSummary(*group, today, totalSITAllowance, daysUsed, mostRecentSITAuthorizedEndDate)
		if summary != nil {
			group.Summary = *summary
			mostRecentSITAuthorizedEndDate = &group.Summary.SITAuthorizedEndDate
			shipmentSITs = append(shipmentSITs, *group)
			daysUsed += summary.DaysInSIT
		}
	}

	return shipmentSITs, nil
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
// This is where the craziest part of the SIT code should ever be (Besides the grouping section)
// Due to our service item architecture, SIT is split across many service items
// and due to existing handlers and service objects, it's possible these SIT service items
// will have discrepancies and information spread across multiple items.
// This SIT summary is to make it readable down the line, and handle all complex calculations
// in one, central location.
//   - sit: The SIT service item grouping we are generating the summary for.
//   - today: The current date, used for calculating days in SIT.
//   - totalSitAllowance: The total number of days allowed for this shipment.
//   - daysUsedSoFar: The number of SIT days that have already been used in previous SIT groupings.
//     This is used to ensure the remaining SIT allowance is correctly applied to this grouping, and is
//     generated by sorting SIT groupings by entry date, and adding to the sum of daysUsed
func (f shipmentSITStatus) generateSITSummary(sit models.SITServiceItemGrouping, today time.Time, totalSitAllowance int, daysUsedSoFar int, mostRecentSITAuthorizedEndDate *time.Time) *models.SITSummary {
	if sit.ServiceItems == nil {
		// Return nil if there are no service items
		return nil
	}

	var earliestSITEntryDate *time.Time
	var earliestSITDepartureDate *time.Time
	var earliestSITAuthorizedEndDate *time.Time
	var earliestSITCustomerContacted *time.Time
	var earliestSITRequestedDelivery *time.Time
	var calculatedTotalDaysInSIT int
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
		if (sitServiceItem.ReService.Code == models.ReServiceCodeDOFSIT) && firstDaySITServiceItemID == uuid.Nil {
			firstDaySITServiceItemID = sitServiceItem.ID
		}

		if (sitServiceItem.ReService.Code == models.ReServiceCodeDDFSIT) && firstDaySITServiceItemID == uuid.Nil {
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
	}

	// Calculate the days in SIT and authorized end date based on the earliest SIT entry date and earliest SIT departure date if any were discovered from the SIT group
	if earliestSITEntryDate != nil {
		calculatedTotalDaysInSIT = daysInSIT(*earliestSITEntryDate, earliestSITDepartureDate, today)
		earliestSITAuthorizedEndDateValue := calculateSITAuthorizedEndDate(totalSitAllowance, calculatedTotalDaysInSIT, *earliestSITEntryDate, daysUsedSoFar, earliestSITDepartureDate, mostRecentSITAuthorizedEndDate)
		earliestSITAuthorizedEndDate = &earliestSITAuthorizedEndDateValue
	}

	// Grab the first Customer Contacted
	for _, sitServiceItem := range sit.ServiceItems {
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
		DaysInSIT:                calculatedTotalDaysInSIT,
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
	if len(shipment.MTOServiceItems) == 0 {
		return nil, shipment, nil
	}

	var shipmentSITStatus services.SITStatus

	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	sitGroupings, err := f.RetrieveShipmentSIT(appCtx, shipment)
	if err != nil {
		return nil, shipment, err
	}

	// Sort the SIT groupings into past, current, future
	shipmentSITGroupings := SortShipmentSITs(sitGroupings, today)

	// From our SIT groupings, grab the current one
	// Current SIT can only be either Current or Future at this time
	currentSIT := getCurrentSIT(shipmentSITGroupings)

	// There were no relevant SIT service items for this shipment
	if currentSIT == nil && len(shipmentSITGroupings.PastSITs) == 0 {
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
	shipmentSITStatus.PastSITs = shipmentSITGroupings.PastSITs

	if currentSIT != nil {
		location := currentSIT.Summary.Location
		firstDaySITServiceItemID := currentSIT.Summary.FirstDaySITServiceItemID
		daysInSIT := daysInSIT(currentSIT.Summary.SITEntryDate, currentSIT.Summary.SITDepartureDate, today)
		sitEntryDate := currentSIT.Summary.SITEntryDate
		sitDepartureDate := currentSIT.Summary.SITDepartureDate
		var sitCustomerContacted, sitRequestedDelivery *time.Time
		sitCustomerContacted = currentSIT.Summary.SITCustomerContacted
		sitRequestedDelivery = currentSIT.Summary.SITRequestedDelivery

		shipmentSITStatus.CurrentSIT = &services.CurrentSIT{
			ServiceItemID:        firstDaySITServiceItemID,
			Location:             location,
			DaysInSIT:            daysInSIT,
			SITEntryDate:         sitEntryDate,
			SITDepartureDate:     sitDepartureDate,
			SITAuthorizedEndDate: currentSIT.Summary.SITAuthorizedEndDate,
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
	if len(shipmentSITs.CurrentSITs) > 0 {
		return getEarliestSIT(shipmentSITs.CurrentSITs)
	} else if len(shipmentSITs.FutureSITs) > 0 {
		return getEarliestSIT(shipmentSITs.FutureSITs)
	}
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
	// Per B-20967, the last day should now always be inclusive even if there is a SIT departure date
	if sitDepartureDate != nil && sitDepartureDate.Before(today) {
		days = int(sitDepartureDate.Sub(sitEntryDate).Hours())/24 + 1
	} else if sitEntryDate.Before(today) || sitEntryDate.Equal(today) {
		days = int(today.Sub(sitEntryDate).Hours())/24 + 1
	}
	return days
}

func CalculateTotalDaysInSIT(shipmentSITs SortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, pastSIT := range shipmentSITs.PastSITs {
		totalDays += daysInSIT(pastSIT.Summary.SITEntryDate, pastSIT.Summary.SITDepartureDate, today)
	}
	for _, currentSIT := range shipmentSITs.CurrentSITs {
		totalDays += daysInSIT(currentSIT.Summary.SITEntryDate, currentSIT.Summary.SITDepartureDate, today)
	}
	return totalDays
}

// adds up all the days from pastSITs
func CalculateTotalPastDaysInSIT(shipmentSITs SortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, pastSIT := range shipmentSITs.PastSITs {
		totalDays += daysInSIT(pastSIT.Summary.SITEntryDate, pastSIT.Summary.SITDepartureDate, today)
	}
	return totalDays
}

// totalDaysUsedByPreviousSITs is the current SUM of past SITs days used
// currentDaysUsedWithinThisSIT is the amount of days spent in SIT for the current SIT so far
// This is a private helper func for the generation of the SIT summaries
func calculateSITAuthorizedEndDate(totalSITAllowance int, currentDaysUsedWithinThisSIT int, sitEntryDate time.Time, totalDaysUsedByPreviousSITs int, sitDepartureDate *time.Time, mostRecentSITAuthorizedEndDate *time.Time) time.Time {
	// Get the remaining allowance for this SIT's authorized end date by
	// subtracting the allowance by the days used by previous SITs.
	remainingSITAllowance := totalSITAllowance - totalDaysUsedByPreviousSITs
	// Use clamp to determine if this SIT was created already in violation of the authorized
	// amount of days. This is a super rare case, but deserves the check nonetheless
	_, err := Clamp(remainingSITAllowance-currentDaysUsedWithinThisSIT, 0, remainingSITAllowance)
	if err != nil && mostRecentSITAuthorizedEndDate != nil {
		// This error is only triggered when past SITs have already exceeded the allowed SIT days.
		// In this case, the new SIT is starting already in violation of the SIT allowance.
		// This should never be able to occur due to service object prevention, but just in case we default to
		// the previous SITs authorized end date.
		return *mostRecentSITAuthorizedEndDate
	}

	// The authorized end date will be the sum of days remaining in the allowance from the entry date
	// Subtract the last day to be inclusive of counting it
	// Eg, this func successfully counts totalSITAllowance days from SITEntryDate, but per customer requirements, they will
	// then count that last day too. So if given 90 days of allowance, we'd get Aug 20 2024 thru Nov 18 2024. 90 days as expected
	// but then if you count the last day, it gets 91. Thus making the calculation incorrect. This is why we subtract a day, to be inclusive of it
	sitAuthorizedEndDate := sitEntryDate.AddDate(0, 0, remainingSITAllowance-1)
	// Ensure that the authorized end date does not go before the entry date
	if sitAuthorizedEndDate.Before(sitEntryDate) {
		sitAuthorizedEndDate = sitEntryDate
	}
	// Now that we have our authorized end date, we need to compare it to the departure date.
	// If the SIT departure date is set and it is before the currently authorized end date
	// then the original SIT authorized end date should be updated to the departure date
	if sitDepartureDate != nil && (sitDepartureDate.Before(sitAuthorizedEndDate) && sitDepartureDate.After(sitEntryDate)) {
		sitAuthorizedEndDate = *sitDepartureDate
	}
	return sitAuthorizedEndDate
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
