package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ServiceParamValueLookupsSuite) TestCubicFeetBilledLookup() {
	key := models.ServiceItemParamNameCubicFeetBilled

	suite.Run("successful CubicFeetBilled lookup, above minimum", func() {
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
					Width:     8 * 12 * 1000,
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

		suite.Equal("1029.33", stringValue)
	})

	suite.Run("When domestic crate volume is less than minimum, billed volume should be set to minimum", func() {
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
					Type:      models.DimensionTypeCrate,
					Length:    1000,
					Height:    1000,
					Width:     1000,
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
					Length:    100,
					Height:    100,
					Width:     100,
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

		suite.Equal("4.00", stringValue)
	})

	suite.Run("When international external crate volume is less than minimum, billed volume should be set to minimum", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeICRT,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		cratingDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeCrate,
					Length:    1000,
					Height:    1000,
					Width:     1000,
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
					Length:    100,
					Height:    100,
					Width:     100,
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
		mtoServiceItem.ExternalCrate = models.BoolPointer(true)
		suite.MustSave(&mtoServiceItem)
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		stringValue, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)

		suite.Equal("4.00", stringValue)
	})

	suite.Run("When international non-external crate volume is less than minimum, billed volume should NOT be set to minimum", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeICRT,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		cratingDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeCrate,
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
		itemDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeItem,
					Length:    100,
					Height:    100,
					Width:     100,
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

		suite.Equal("1.00", stringValue)
	})

	suite.Run("missing dimension should error", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)

		suite.Error(err)
		suite.Contains(err.Error(), "missing crate dimensions")
	})
}
