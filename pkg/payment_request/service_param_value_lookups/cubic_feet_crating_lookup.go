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
	// Each service item has an array of dimensions. We expect there to be at most one
	// dimension for the crate size.
	for _, dimension := range keyData.MTOServiceItem.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			// is there a classier way to work with unit.ThousandthInches?
			lengthFeet := dimension.Length / thousandthInchesPerFoot
			heightFeet := dimension.Height / thousandthInchesPerFoot
			widthFeet := dimension.Width / thousandthInchesPerFoot

			volume := lengthFeet * heightFeet * widthFeet
			return fmt.Sprint(*volume.Int32Ptr()), nil
		}
	}
	// TODO maybe add a check for multiple crating dimensions

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to calculate crate volume due to missing crate dimensions")
}
