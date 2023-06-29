package mtoshipment

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// OriginSITLocation is the constant representing when the shipment in storage occurs at the origin
const OriginSITLocation = "ORIGIN"

// DestinationSITLocation is the constant representing when the shipment in storage occurs at the destination
const DestinationSITLocation = "DESTINATION"

type shipmentSITStatus struct {
}

// NewShipmentSITStatus creates a new instance of the service object that implements calculating a shipments SIT summary
func NewShipmentSITStatus() services.ShipmentSITStatus {
	return &shipmentSITStatus{}
}

type sortedShipmentSITs struct {
	pastSITs    []models.MTOServiceItem
	currentSITs []models.MTOServiceItem
	futureSITs  []models.MTOServiceItem
	pendingSITs []models.MTOServiceItem
}

func newSortedShipmentSITs() sortedShipmentSITs {
	return sortedShipmentSITs{
		pastSITs:    make([]models.MTOServiceItem, 0),
		currentSITs: make([]models.MTOServiceItem, 0),
		futureSITs:  make([]models.MTOServiceItem, 0),
		pendingSITs: make([]models.MTOServiceItem, 0),
	}
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

	shipmentSITs := newSortedShipmentSITs()

	for _, serviceItem := range shipment.MTOServiceItems {
		// only departure SIT service items have a departure date
		if code := serviceItem.ReService.Code; code == models.ReServiceCodeDOPSIT || code == models.ReServiceCodeDDDSIT {
			if serviceItem.Status == models.MTOServiceItemStatusApproved {
				if serviceItem.SITEntryDate.After(today) {
					shipmentSITs.futureSITs = append(shipmentSITs.futureSITs, serviceItem)
				} else if serviceItem.SITDepartureDate != nil && serviceItem.SITDepartureDate.Before(today) {
					shipmentSITs.pastSITs = append(shipmentSITs.pastSITs, serviceItem)
				} else {
					shipmentSITs.currentSITs = append(shipmentSITs.currentSITs, serviceItem)
				}
			} else if serviceItem.Status == models.MTOServiceItemStatusSubmitted {
				shipmentSITs.pendingSITs = append(shipmentSITs.pendingSITs, serviceItem)
			}
		}
	}

	currentSIT := getCurrentSIT(shipmentSITs)

	// There were no relevant SIT service items for this shipment
	if currentSIT == nil && len(shipmentSITStatus.PastSITs) == 0 {
		return nil, nil
	}

	shipmentSITStatus.ShipmentID = shipment.ID
	totalSITAllowance, err := f.CalculateShipmentSITAllowance(appCtx, shipment)
	if err != nil {
		return nil, err
	}
	shipmentSITStatus.TotalSITDaysUsed = calculateTotalDaysInSIT(shipmentSITs, today)
	shipmentSITStatus.TotalDaysRemaining = totalSITAllowance - shipmentSITStatus.TotalSITDaysUsed

	if currentSIT != nil {
		if currentSIT.ReService.Code == models.ReServiceCodeDOPSIT {
			shipmentSITStatus.CurrentSIT.Location = OriginSITLocation
		} else {
			shipmentSITStatus.CurrentSIT.Location = DestinationSITLocation
		}
		shipmentSITStatus.CurrentSIT.SITEntryDate = *currentSIT.SITEntryDate
		shipmentSITStatus.CurrentSIT.SITDepartureDate = currentSIT.SITDepartureDate
		shipmentSITStatus.CurrentSIT.SITAllowanceEndDate = calculateSITAllowanceEndDate(shipmentSITStatus, today)
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
service items that start in the future, and then SIT service items that have
yet to be approved.
*/
func getCurrentSIT(shipmentSITs sortedShipmentSITs) *models.MTOServiceItem {
	if len(shipmentSITs.currentSITs) > 0 {
		return getEarliestSIT(shipmentSITs.currentSITs)
	} else if len(shipmentSITs.futureSITs) > 0 {
		return getEarliestSIT(shipmentSITs.futureSITs)
	} else if len(shipmentSITs.pendingSITs) > 0 {
		return getEarliestSIT(shipmentSITs.pendingSITs)
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

func calculateTotalDaysInSIT(shipmentSITs sortedShipmentSITs, today time.Time) int {
	totalDays := 0
	for _, serviceItem := range shipmentSITs.pastSITs {
		totalDays += daysInSIT(serviceItem, today)
	}
	for _, serviceItem := range shipmentSITs.currentSITs {
		totalDays += daysInSIT(serviceItem, today)
	}
	return totalDays
}

func calculateSITAllowanceEndDate(shipmentSITStatus services.SITStatus, today time.Time) time.Time {
	//current SIT
	if shipmentSITStatus.CurrentSIT.SITEntryDate.Before(today) {
		return today.AddDate(0, 0, shipmentSITStatus.TotalDaysRemaining)
	}
	// future SIT
	return shipmentSITStatus.CurrentSIT.SITEntryDate.AddDate(0, 0, shipmentSITStatus.TotalDaysRemaining)

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
