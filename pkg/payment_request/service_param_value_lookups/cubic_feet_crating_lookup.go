package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

const (
	thousandthInchesPerFoot = 12000
)

// CubicFeetCratingLookup does lookup for CubicFeetCrating
type CubicFeetCratingLookup struct {
}

func (c CubicFeetCratingLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range keyData.MTOServiceItem.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			lengthFeet := dimension.Length / thousandthInchesPerFoot
			heightFeet := dimension.Height / thousandthInchesPerFoot
			widthFeet := dimension.Width / thousandthInchesPerFoot

			volume := lengthFeet * heightFeet * widthFeet
			return fmt.Sprint(*volume.Int32Ptr()), nil
		}
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to calculate crate volume due to missing crate dimensions")
}
