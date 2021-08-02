package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	minCubicFeetBilled = unit.CubicFeet(4.0)
)

// CubicFeetBilledLookup does lookup for CubicFeetBilled
type CubicFeetBilledLookup struct {
	Dimensions models.MTOServiceItemDimensions
}

func (c CubicFeetBilledLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// call CubicFeetCratingLookup
	// convert string to number
	// do the math

	cubicFeet, err := CubicFeetCratingLookup(c).lookup(keyData)
	if err != nil {
		return "", err
	}

	parsedCubicFeetCrating, err := unit.CubicFeetFromString(cubicFeet)
	if err != nil {
		return "", fmt.Errorf("could not convert CubicFeetCratingLookup [%s] to integer", cubicFeet)
	}

	if parsedCubicFeetCrating < minCubicFeetBilled {
		return minCubicFeetBilled.String(), nil
	}

	return cubicFeet, nil
}
