package serviceparamvaluelookups

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// MTOEarliestRequestedPickupLookup does lookup on the MTOEarliestRequestedPickup timestamp
type MTOEarliestRequestedPickupLookup struct {
}

func (m MTOEarliestRequestedPickupLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	db := appCtx.DB()

	// Get the MoveTaskOrder
	moveTaskOrderID := keyData.MoveTaskOrderID
	var moveTaskOrder models.Move
	err := db.EagerPreload("MTOShipments").Find(&moveTaskOrder, moveTaskOrderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", apperror.NewNotFoundError(moveTaskOrderID, "looking for MoveTaskOrderID")
		default:
			return "", apperror.NewQueryError("Move", err, "")
		}
	}

	var earliestPickupDate *time.Time
	for _, shipment := range moveTaskOrder.MTOShipments {
		if shipment.ShipmentType != models.MTOShipmentTypePPM && earliestPickupDate == nil && shipment.DeletedAt == nil {
			earliestPickupDate = shipment.RequestedPickupDate
		}

		if shipment.ShipmentType != models.MTOShipmentTypePPM && shipment.RequestedPickupDate.Before(*earliestPickupDate) && shipment.DeletedAt == nil {
			earliestPickupDate = shipment.RequestedPickupDate
		}
	}

	utcMidnight := models.TimePointer(time.Date(
		earliestPickupDate.Year(),
		earliestPickupDate.Month(),
		earliestPickupDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	))

	earliestPickupDate = utcMidnight
	return (*earliestPickupDate).Format(ghcrateengine.TimestampParamFormat), nil
}
