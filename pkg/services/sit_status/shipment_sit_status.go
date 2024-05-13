package sitstatus

import (
	"time"

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

// CalculateShipmentSITStatus creates a SIT Status for payload to be used in
// multiple handlers in the `ghcapi` package for the MTOShipment handlers.
func (f shipmentSITStatus) CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) (*services.SITStatus, models.MTOShipment, error) {
	if shipment.MTOServiceItems == nil || len(shipment.MTOServiceItems) == 0 {
		return nil, shipment, nil
	}

	var shipmentSITStatus services.SITStatus

	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	shipmentSITs := SortShipmentSITs(shipment, today)

	currentSIT := getCurrentSIT(shipmentSITs)

	// There were no relevant SIT service items for this shipment
	if currentSIT == nil && len(shipmentSITs.pastSITs) == 0 {
		return nil, shipment, nil
	}

	shipmentSITStatus.ShipmentID = shipment.ID
	totalSITAllowance, err := f.CalculateShipmentSITAllowance(appCtx, shipment)
	if err != nil {
		return nil, shipment, err
	}
	totalSITDaysUsedClampedResult, totalDaysUsedErr := Clamp(CalculateTotalDaysInSIT(shipmentSITs, today), 0, totalSITAllowance)
	if totalDaysUsedErr != nil {
		return nil, shipment, err
	}
	shipmentSITStatus.TotalSITDaysUsed = totalSITDaysUsedClampedResult
	shipmentSITStatus.CalculatedTotalDaysInSIT = CalculateTotalDaysInSIT(shipmentSITs, today)
	shipmentSITStatus.TotalDaysRemaining = totalSITAllowance - shipmentSITStatus.TotalSITDaysUsed
	shipmentSITStatus.PastSITs = shipmentSITs.pastSITs

	if currentSIT != nil {
		location := DestinationSITLocation
		if currentSIT.ReService.Code == models.ReServiceCodeDOFSIT || currentSIT.ReService.Code == models.ReServiceCodeDOASIT {
			location = OriginSITLocation
		}
		daysInSIT := daysInSIT(*currentSIT, today)
		sitEntryDate := *currentSIT.SITEntryDate
		sitDepartureDate := currentSIT.SITDepartureDate
		sitAuthorizedEndDate := CalculateSITAuthorizedEndDate(totalSITAllowance, daysInSIT, sitEntryDate, shipmentSITStatus.CalculatedTotalDaysInSIT)
		var sitCustomerContacted, sitRequestedDelivery *time.Time
		sitCustomerContacted = currentSIT.SITCustomerContacted
		sitRequestedDelivery = currentSIT.SITRequestedDelivery

		doaSIT := getAdditionalSIT(shipmentSITs, shipment, today)

		if doaSIT != nil {
			sitCustomerContacted = doaSIT.SITCustomerContacted
			sitRequestedDelivery = doaSIT.SITRequestedDelivery
			sitDepartureDate = doaSIT.SITDepartureDate
		}

		shipmentSITStatus.CurrentSIT = &services.CurrentSIT{
			ServiceItemID:        currentSIT.ID,
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
			transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
				verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(shipment.ID, err, verrs, "invalid input found while updating the shipment sit auth date")
				} else if err != nil {
					return apperror.NewQueryError("shipment", err, "")
				}

				return nil
			})

			if transactionError != nil {
				return nil, shipment, transactionError
			}
		}
	}

	return &shipmentSITStatus, shipment, nil
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

// Private function getAdditionalSIT is used to return the current SIT
// service item with the reServiceCode of DOASIT or DDASIT
func getAdditionalSIT(shipmentSITs SortedShipmentSITs, shipment models.MTOShipment, today time.Time) *models.MTOServiceItem {
	for _, serviceItem := range shipment.MTOServiceItems {
		// only departure SIT service items have a departure date
		if code := serviceItem.ReService.Code; (code == models.ReServiceCodeDOASIT || code == models.ReServiceCodeDDASIT) &&
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

	if len(shipmentSITs.currentSITs) == 0 {
		return nil
	}

	for _, serviceItem := range shipmentSITs.currentSITs {
		if code := serviceItem.ReService.Code; code == models.ReServiceCodeDOASIT || code == models.ReServiceCodeDDASIT {
			return &serviceItem
		}
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

// adds up all the days from pastSITs
func CalculateTotalPastDaysInSIT(shipmentSITs SortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, serviceItem := range shipmentSITs.pastSITs {
		totalDays += daysInSIT(serviceItem, today)
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
