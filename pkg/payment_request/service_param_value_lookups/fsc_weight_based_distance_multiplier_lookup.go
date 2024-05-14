package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

const weightBasedDistanceMultiplierLevelOne = "0.000417"
const weightBasedDistanceMultiplierLevelTwo = "0.0006255"
const weightBasedDistanceMultiplierLevelThree = "0.000834"
const weightBasedDistanceMultiplierLevelFour = "0.00139"

// FSCWeightBasedDistanceMultiplierLookup does lookup on fuel surcharge related weight based distance multiplier rate based on billed weight
type FSCWeightBasedDistanceMultiplierLookup struct {
	MTOShipment models.MTOShipment
}

func (r FSCWeightBasedDistanceMultiplierLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	weight, err := WeightBilledLookup(r).lookup(appCtx, keyData)
	if err != nil {
		return "", err
	}

	weightBilled, err := strconv.Atoi(weight)
	if err != nil {
		return "", fmt.Errorf("could not convert WeightBilledLookup [%s] to integer", weight)
	}

	if weightBilled <= 5000 {
		return weightBasedDistanceMultiplierLevelOne, nil
	} else if weightBilled <= 10000 {
		return weightBasedDistanceMultiplierLevelTwo, nil
	} else if weightBilled <= 24000 {
		return weightBasedDistanceMultiplierLevelThree, nil
		//nolint:revive
	} else {
		return weightBasedDistanceMultiplierLevelFour, nil
	}
}
