package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// ActualPickupDateLookup does lookup on actual pickup date
type ActualPickupDateLookup struct {
	MTOShipment models.MTOShipment
}

func (r ActualPickupDateLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	db := appCtx.DB()

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("MTOShipment").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", apperror.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// Make sure there's an MTOShipment since that's nullable
	mtoShipmentID := mtoServiceItem.MTOShipmentID
	if mtoShipmentID == nil {
		return "", apperror.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	// Make sure there's a actual pickup date since that's nullable
	requestedPickupDate := mtoServiceItem.MTOShipment.ActualPickupDate
	if requestedPickupDate == nil {
		return "", fmt.Errorf("could not find an actual pickup date for MTOShipmentID [%s]", mtoShipmentID)
	}

	return requestedPickupDate.Format(ghcrateengine.DateParamFormat), nil
}
