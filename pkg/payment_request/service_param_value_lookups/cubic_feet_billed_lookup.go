package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	//  "github.com/transcom/mymove/pkg/unit"
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

	_, err = strconv.Atoi(cubicFeet) // convert to new unit CubicFeet //cubicFeetBilled
	if err != nil {
		return "", fmt.Errorf("could not convert CubicFeetCratingLookup [%s] to integer", cubicFeet)
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "")

}
