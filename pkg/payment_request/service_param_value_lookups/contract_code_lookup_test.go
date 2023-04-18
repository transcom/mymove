package serviceparamvaluelookups

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestContractCodeLookup() {
	key := models.ServiceItemParamNameContractCode

	suite.Run("golden path", func() {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, valueStr)
	})

	suite.Run("golden path with param cache", func() {
		// DLH
		mtoServiceItem1 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
		}, nil)

		// ContractCode
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

		_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID:             mtoServiceItem1.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey1.ID,
				ServiceItemParamKey:   serviceItemParamKey1,
			},
		})

		paramCache := NewServiceParamsCache()
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem1, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), &paramCache)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, valueStr)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItem1.MTOShipmentID, key)
		suite.Equal(ghcrateengine.DefaultContractCode, *paramCacheValue)
	})
}
