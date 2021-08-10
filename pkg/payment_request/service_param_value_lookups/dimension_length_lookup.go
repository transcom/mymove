package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// DimensionLengthLookup does lookup for DimensionLengthLookup
type DimensionLengthLookup struct {
	Dimensions models.MTOServiceItemDimensions
}

func (d DimensionLengthLookup) lookup(appCfg appconfig.AppConfig, keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range d.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			lengthInches := int(dimension.Length.ToInches())

			return strconv.Itoa(lengthInches), nil
		}
	}

	return "", services.NewConflictError(keyData.MTOServiceItemID, "unable to find length crate dimension")
}
