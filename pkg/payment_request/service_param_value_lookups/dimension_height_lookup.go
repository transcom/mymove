package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// DimensionHeightLookup does lookup for DimensionHeightLookup
type DimensionHeightLookup struct {
	Dimensions models.MTOServiceItemDimensions
}

func (d DimensionHeightLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range d.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			heightInches := int(dimension.Height.ToInches())

			return strconv.Itoa(heightInches), nil
		}
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to find height crate dimension")
}
