package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"database/sql"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// WeightBilledActualLookup does lookup on actual weight billed
type WeightBilledActualLookup struct {
}

func (r WeightBilledActualLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	var value string

	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("MTOShipment").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// Make sure there's an MTOShipment since that's nullable
	mtoShipmentID := mtoServiceItem.MTOShipmentID
	if mtoShipmentID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	// Make sure there's an estimated weight since that's nullable
	estimatedWeight := mtoServiceItem.MTOShipment.PrimeEstimatedWeight
	if estimatedWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find estimated weight for MTOShipmentID [%s]", mtoShipmentID)
	}

	// Make sure there's an actual weight since that's nullable
	actualWeight := mtoServiceItem.MTOShipment.PrimeActualWeight
	if actualWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find estimated weight for MTOShipmentID [%s]", mtoShipmentID)
	}

	value = fmt.Sprintf("%d", int(*actualWeight))

	return value, nil
}
