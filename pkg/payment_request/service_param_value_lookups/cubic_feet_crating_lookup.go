package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// CubicFeetCratingLookup does lookup for CubicFeetCrating
type CubicFeetCratingLookup struct {
	Dimensions models.MTOServiceItemDimensions
}

func (c CubicFeetCratingLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	// Each service item has an array of dimensions. There is a DB constraint preventing
	// more than one dimension of each type for a given service item, so we just have to
	// look for the first crating dimension.
	for _, dimension := range c.Dimensions {
		if dimension.Type == models.DimensionTypeCrate {
			return dimension.Volume().ToCubicFeet().String(), nil
		}
	}

	return "", apperror.NewConflictError(keyData.MTOServiceItemID, "unable to calculate crate volume due to missing crate dimensions")
}
