package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
)

// FSCWeightBasedDistanceMultiplierLookup does lookup on fuel surcharge related weight based distance multiplier rate based on billed actual weight
type FSCWeightBasedDistanceMultiplierLookup struct {
}

func (r FSCWeightBasedDistanceMultiplierLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	weight, err := WeightBilledActualLookup{}.lookup(keyData)
	if err != nil {
		return "", err
	}

	weightBilledActual, err := strconv.Atoi(weight)
	if err != nil {
		return "", fmt.Errorf("could not convert WeightBilledActualLookup [%s] to integer", weight)
	}

	if weightBilledActual <= 5000 {
		return "0.000417", nil
	} else if weightBilledActual <= 10000 {
		return "0.0006255", nil
	} else if weightBilledActual <= 24000 {
		return "0.000834", nil
	} else {
		return "0.00139", nil
	}
}
