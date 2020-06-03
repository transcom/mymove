package serviceparamvaluelookups

import (
	"fmt"
	"math"

	"github.com/gofrs/uuid"

	"database/sql"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// WeightBilledActualLookup does lookup on actual weight billed
type WeightBilledActualLookup struct {
}

func (r WeightBilledActualLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
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

	var value string
	estimatedWeightCap := math.Round(float64(*estimatedWeight) * 1.10)
	if float64(*actualWeight) > estimatedWeightCap {
		value = fmt.Sprintf("%d", int(estimatedWeightCap))
	} else {
		value = applyMinimum(mtoServiceItem.ReService.Code, mtoServiceItem.MTOShipment.ShipmentType, int(*actualWeight))
	}

	return value, nil
}

// Looks at code and applies minimum if necessary, otherwise returns actual
func applyMinimum(code models.ReServiceCode, shipmentType models.MTOShipmentType, actual int) string {
	result := actual
	switch shipmentType {
	case models.MTOShipmentTypeInternationalUB:
		switch code {
		case models.ReServiceCodeIOSHUT,
			models.ReServiceCodeIDSHUT:
			if int(actual) < 300 {
				result = 300
			}
		}
	default:
		switch code {
		case models.ReServiceCodeDLH,
			models.ReServiceCodeDSH,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDOFSIT,
			models.ReServiceCodeDDFSIT,
			models.ReServiceCodeDOASIT,
			models.ReServiceCodeDDASIT,
			models.ReServiceCodeDOPSIT,
			models.ReServiceCodeDDDSIT,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
			models.ReServiceCodeDOSHUT,
			models.ReServiceCodeDDSHUT,
			models.ReServiceCodeIOOLH,
			models.ReServiceCodeICOLH,
			models.ReServiceCodeIOCLH,
			models.ReServiceCodeIHPK,
			models.ReServiceCodeIHUPK,
			models.ReServiceCodeIOFSIT,
			models.ReServiceCodeIDFSIT,
			models.ReServiceCodeIOASIT,
			models.ReServiceCodeIDASIT,
			models.ReServiceCodeIOPSIT,
			models.ReServiceCodeIDDSIT,
			models.ReServiceCodeIOSHUT,
			models.ReServiceCodeIDSHUT:
			if int(actual) < 500 {
				result = 500
			}
		case models.ReServiceCodeIOOUB,
			models.ReServiceCodeICOUB,
			models.ReServiceCodeIOCUB,
			models.ReServiceCodeIUBPK,
			models.ReServiceCodeIUBUPK:
			if int(actual) < 300 {
				result = 300
			}
		}
	}
	return fmt.Sprintf("%d", result)
}
