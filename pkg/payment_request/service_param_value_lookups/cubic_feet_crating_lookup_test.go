package serviceparamvaluelookups

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestCubicFeetCratingLookup() {
	key := models.ServiceItemParamNameCubicFeetCrating

	suite.T().Run("successful CubicFeetCrating lookup", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDCRT,
			},
		})
		cratingDimension := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				MTOServiceItemID: mtoServiceItem.ID,
				Type:             models.DimensionTypeCrate,
				// These dimensions are chosen to overflow 32bit ints if multiplied, and give a fractional result
				// when converted to cubic feet.
				Length:    16*12*1000 + 1000,
				Height:    8 * 12 * 1000,
				Width:     8 * 12 * 1000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		})
		itemDimension := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				MTOServiceItemID: mtoServiceItem.ID,
				Type:             models.DimensionTypeItem,
				Length:           12000,
				Height:           12000,
				Width:            12000,
				CreatedAt:        time.Time{},
				UpdatedAt:        time.Time{},
			},
		})
		mtoServiceItem.Dimensions = []models.MTOServiceItemDimension{itemDimension, cratingDimension}
		suite.MustSave(&mtoServiceItem)
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		stringValue, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)

		suite.Equal("1029.33", stringValue)
	})

	suite.T().Run("missing dimension should error", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)

		suite.Error(err)
		suite.Contains(err.Error(), "missing crate dimensions")
	})
}
