package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/services"

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

func (c CubicFeetBilledLookup) lookup(appCfg appconfig.AppConfig, keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range c.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			volume := dimension.Volume().ToCubicFeet()
			if volume < minCubicFeetBilled {
				volume = minCubicFeetBilled
			}
			return volume.String(), nil
		}
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to calculate crate volume due to missing crate dimensions")
}
