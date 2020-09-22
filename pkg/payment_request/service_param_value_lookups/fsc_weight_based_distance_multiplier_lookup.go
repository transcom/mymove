package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

const weightBasedDistanceMultiplierLevelOne = "0.000417"
const weightBasedDistanceMultiplierLevelTwo = "0.0006255"
const weightBasedDistanceMultiplierLevelThree = "0.000834"
const weightBasedDistanceMultiplierLevelFour = "0.00139"

// FSCWeightBasedDistanceMultiplierLookup does lookup on fuel surcharge related weight based distance multiplier rate based on billed actual weight
type FSCWeightBasedDistanceMultiplierLookup struct {
	MTOShipment models.MTOShipment
}

func (r FSCWeightBasedDistanceMultiplierLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	weight, err := WeightBilledActualLookup{
		MTOShipment: r.MTOShipment,
	}.lookup(keyData)
	if err != nil {
		return "", err
	}

	weightBilledActual, err := strconv.Atoi(weight)
	if err != nil {
		return "", fmt.Errorf("could not convert WeightBilledActualLookup [%s] to integer", weight)
	}

	if weightBilledActual <= 5000 {
		return weightBasedDistanceMultiplierLevelOne, nil
	} else if weightBilledActual <= 10000 {
		return weightBasedDistanceMultiplierLevelTwo, nil
	} else if weightBilledActual <= 24000 {
		return weightBasedDistanceMultiplierLevelThree, nil
	} else {
		return weightBasedDistanceMultiplierLevelFour, nil
	}
}
