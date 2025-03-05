package serviceparamvaluelookups

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ServiceParamValueLookupsSuite) TestDimensionLengthLookup() {
	key := models.ServiceItemParamNameDimensionLength

	suite.Run("successful DimensionLength lookup", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDCRT,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		cratingDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type: models.DimensionTypeCrate,
					// These dimensions are chosen to overflow 32bit ints if multiplied, and give a fractional result
					// when converted to cubic feet.
					Length:    16*12*1000 + 1000,
					Height:    8 * 12 * 1000,
					Width:     9 * 12 * 1000,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			{
				Model:    mtoServiceItem,
				LinkOnly: true,
			},
		}, nil)
		itemDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeItem,
					Length:    12000,
					Height:    12000,
					Width:     12000,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			{
				Model:    mtoServiceItem,
				LinkOnly: true,
			},
		}, nil)
		mtoServiceItem.Dimensions = []models.MTOServiceItemDimension{itemDimension, cratingDimension}
		suite.MustSave(&mtoServiceItem)
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		stringValue, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)

		suite.Equal("193", stringValue)
		suite.Equal("193", strconv.Itoa(int(cratingDimension.Length.ToInches())))
	})

	suite.Run("missing dimension should error", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)

		suite.Error(err)
		suite.Contains(err.Error(), "unable to find length crate dimension")
	})
}
