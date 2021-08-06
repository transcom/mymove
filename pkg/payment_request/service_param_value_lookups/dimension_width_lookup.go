package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// DimensionWidthLookup does lookup for DimensionWidthLookup
type DimensionWidthLookup struct {
	Dimensions models.MTOServiceItemDimensions
}

func (d DimensionWidthLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range d.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			widthInches := int(dimension.Width.ToInches())

			return strconv.Itoa(widthInches), nil
		}
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to find width crate dimension")
}
