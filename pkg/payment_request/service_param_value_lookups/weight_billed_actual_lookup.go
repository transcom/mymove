package serviceparamvaluelookups

import (
	"fmt"
	"math"

	"github.com/gofrs/uuid"

	"database/sql"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
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
	err := db.Eager("ReService", "MTOShipment").Find(&mtoServiceItem, mtoServiceItemID)
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
		return "", fmt.Errorf("could not find actual weight for MTOShipmentID [%s]", mtoShipmentID)
	}

	estimatedWeightCap := math.Round(float64(*estimatedWeight) * 1.10)
	if float64(*actualWeight) > estimatedWeightCap {
		value = fmt.Sprintf("%d", int(estimatedWeightCap))
	} else if fiveHundredMinimumApplies(mtoServiceItem.ReService.Code, *actualWeight) {
		value = "500"
	} else {
		value = fmt.Sprintf("%d", int(*actualWeight))
	}

	return value, nil
}

func fiveHundredMinimumApplies(code models.ReServiceCode, actual unit.Pound) bool {
	switch code {
	case models.ReServiceCodeDLH:
		return int(actual) < 500
	case models.ReServiceCodeDSH:
		return int(actual) < 500
	case models.ReServiceCodeDOP:
		return int(actual) < 500
	case models.ReServiceCodeDDP:
		return int(actual) < 500
	case models.ReServiceCodeDOFSIT:
		return int(actual) < 500
	case models.ReServiceCodeDDFSIT:
		return int(actual) < 500
	case models.ReServiceCodeDOASIT:
		return int(actual) < 500
	case models.ReServiceCodeDDASIT:
		return int(actual) < 500
	case models.ReServiceCodeDOPSIT:
		return int(actual) < 500
	case models.ReServiceCodeDDDSIT:
		return int(actual) < 500
	case models.ReServiceCodeDPK:
		return int(actual) < 500
	case models.ReServiceCodeDUPK:
		return int(actual) < 500
	default:
		return false
	}
}
