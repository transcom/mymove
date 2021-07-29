package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// const (
// 	thousandthInchesPerFoot = 12000
// )

// CubicFeetCratingLookup does lookup for CubicFeetCrating
type CubicFeetCratingLookup struct {
	Dimensions models.MTOServiceItemDimensions
}

func (c CubicFeetCratingLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range c.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			volume := dimension.Volume() // dimensions.Volume().ToFeet()
			return fmt.Sprintf("%.2f", float64(volume)), nil
		}
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to calculate crate volume due to missing crate dimensions")
}
