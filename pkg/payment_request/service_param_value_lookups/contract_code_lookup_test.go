package serviceparamvaluelookups

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestContractCodeLookup() {
	key := models.ServiceItemParamNameContractCode

	suite.Run("golden path", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(testdatagen.DefaultContractCode, valueStr)
	})

	suite.Run("golden path with param cache", func() {
		// DLH
		mtoServiceItem1 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		serviceItemParamKey1 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameContractCode,
					Description: "contract code",
					Type:        models.ServiceItemParamTypeString,
					Origin:      models.ServiceItemParamOriginSystem,
				},
			},
		}, nil)

		factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey1,
				LinkOnly: true,
			},
		}, nil)

		paramCache := NewServiceParamsCache()
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem1, uuid.Must(uuid.NewV4()), mtoServiceItem1.MoveTaskOrderID, &paramCache)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		suite.Equal(testdatagen.DefaultContractCode, valueStr)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItem1.MTOShipmentID, key)
		suite.Equal(testdatagen.DefaultContractCode, *paramCacheValue)
	})
}
