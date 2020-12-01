package serviceparamvaluelookups

import (
	"fmt"
	"math"

	"github.com/transcom/mymove/pkg/models"
)

// WeightBilledActualLookup does lookup on actual weight billed
type WeightBilledActualLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightBilledActualLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Make sure there's an estimated weight since that's nullable
	estimatedWeight := r.MTOShipment.PrimeEstimatedWeight
	if estimatedWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find estimated weight for MTOShipmentID [%s]", r.MTOShipment.ID)
	}

	// Make sure there's an actual weight since that's nullable
	actualWeight := r.MTOShipment.PrimeActualWeight
	if actualWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find actual weight for MTOShipmentID [%s]", r.MTOShipment.ID)
	}

	var value string
	estimatedWeightCap := math.Round(float64(*estimatedWeight) * 1.10)
	if float64(*actualWeight) > estimatedWeightCap {
		value = fmt.Sprintf("%d", int(estimatedWeightCap))
	} else {
		value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(*actualWeight))
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
