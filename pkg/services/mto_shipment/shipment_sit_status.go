package mtoshipment

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// OriginSITLocation is the constant representing when the shipment in storage occurs at the origin
const OriginSITLocation = "At origin"

// DestinationSITLocation is the constant representing when the shipment in storage occurs at the destination
const DestinationSITLocation = "At destination"

type shipmentSITStatus struct {
}

// NewShipmentSITStatus creates a new instance of the service object that implements calculating a shipments SIT summary
func NewShipmentSITStatus() services.ShipmentSITStatus {
	return &shipmentSITStatus{}
}

func (f shipmentSITStatus) CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) *services.SITStatus {
	if shipment.MTOServiceItems == nil || len(shipment.MTOServiceItems) == 0 {
		return nil
	}

	var shipmentSITStatus services.SITStatus
	var mostRecentSIT *models.MTOServiceItem

	// Collect all Departure SITs from origin and destination and find the most recent SIT service item
	for i, serviceItem := range shipment.MTOServiceItems {
		// only departure SIT service items have a departure date
		if code := serviceItem.ReService.Code; code == models.ReServiceCodeDOPSIT || code == models.ReServiceCodeDDDSIT {
			if mostRecentSIT == nil {
				shipmentSITStatus.DaysInSIT = daysInSIT(serviceItem)
				shipmentSITStatus.TotalSITDaysUsed += shipmentSITStatus.DaysInSIT
				mostRecentSIT = &shipment.MTOServiceItems[i]
			} else if mostRecentSIT.SITEntryDate.Before(*serviceItem.SITEntryDate) {
				shipmentSITStatus.PastSITs = append(shipmentSITStatus.PastSITs, *mostRecentSIT)
				shipmentSITStatus.DaysInSIT = daysInSIT(*mostRecentSIT)
				shipmentSITStatus.TotalSITDaysUsed += shipmentSITStatus.DaysInSIT
				mostRecentSIT = &shipment.MTOServiceItems[i]
			} else {
				shipmentSITStatus.TotalSITDaysUsed += daysInSIT(serviceItem)
				shipmentSITStatus.PastSITs = append(shipmentSITStatus.PastSITs, shipment.MTOServiceItems[i])
			}
		}
	}

	// There were no departure SIT service items for this shipment
	if mostRecentSIT == nil {
		return nil
	}

	shipmentSITStatus.ShipmentID = shipment.ID

	if mostRecentSIT.ReService.Code == models.ReServiceCodeDOPSIT {
		shipmentSITStatus.Location = OriginSITLocation
	} else {
		shipmentSITStatus.Location = DestinationSITLocation
	}

	shipmentSITStatus.SITEntryDate = *mostRecentSIT.SITEntryDate
	shipmentSITStatus.SITDepartureDate = mostRecentSIT.SITDepartureDate

	today := time.Now()
	if mostRecentSIT.SITDepartureDate != nil && today.Before(*mostRecentSIT.SITDepartureDate) {
		daysRemaining := int(mostRecentSIT.SITDepartureDate.Sub(today).Hours()) / 24
		shipmentSITStatus.DaysRemaining = &daysRemaining
	}

	return &shipmentSITStatus
}

func daysInSIT(serviceItem models.MTOServiceItem) int {
	today := time.Now()
	if serviceItem.SITDepartureDate != nil && serviceItem.SITDepartureDate.Before(today) {
		return int(serviceItem.SITDepartureDate.Sub(*serviceItem.SITEntryDate).Hours()) / 24
	} else if serviceItem.SITEntryDate.Before(today) {
		return int(today.Sub(*serviceItem.SITEntryDate).Hours()) / 24
	}

	return 0
}

func (f shipmentSITStatus) CalculateShipmentsSITStatuses(appCtx appcontext.AppContext, shipments []models.MTOShipment) map[string]services.SITStatus {
	shipmentsSITStatuses := map[string]services.SITStatus{}

	for _, shipment := range shipments {
		shipmentSITStatus := f.CalculateShipmentSITStatus(appCtx, shipment)
		if shipmentSITStatus != nil {
			shipmentsSITStatuses[shipment.ID.String()] = *shipmentSITStatus
		}
	}

	return shipmentsSITStatuses
}
