package serviceparamvaluelookups

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestDimensionLengthLookup() {
	key := models.ServiceItemParamNameDimensionLength

	suite.Run("successful DimensionLength lookup", func() {
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
				Width:     9 * 12 * 1000,
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
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		stringValue, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)

		suite.Equal("193", stringValue)
		suite.Equal("193", strconv.Itoa(int(cratingDimension.Length.ToInches())))
	})

	suite.Run("missing dimension should error", func() {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)

		suite.Error(err)
		suite.Contains(err.Error(), "unable to find length crate dimension")
	})
}
