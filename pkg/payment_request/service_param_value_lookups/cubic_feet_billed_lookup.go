package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	minCubicFeetBilled = unit.CubicFeet(4.0)
)

// CubicFeetBilledLookup does lookup for CubicFeetBilled
type CubicFeetBilledLookup struct {
	Dimensions  models.MTOServiceItemDimensions
	ServiceItem models.MTOServiceItem
}

func (c CubicFeetBilledLookup) lookup(_ appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	isIntlCrateUncrate := c.ServiceItem.ReService.Code == models.ReServiceCodeICRT || c.ServiceItem.ReService.Code == models.ReServiceCodeIUCRT
	isExternalCrate := c.ServiceItem.ExternalCrate != nil && *c.ServiceItem.ExternalCrate

	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range c.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			volume := dimension.Volume().ToCubicFeet()
			if (!isIntlCrateUncrate || isExternalCrate) && volume < minCubicFeetBilled {
				volume = minCubicFeetBilled
			}
			return volume.String(), nil
		}
	}

	return "", apperror.NewConflictError(keyData.MTOServiceItemID, "unable to calculate crate volume due to missing crate dimensions")
}
