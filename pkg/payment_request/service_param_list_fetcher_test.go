package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestHelperSuite) TestFetchServiceParamList() {
	// Make a couple of services
	dlhService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
	dopService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOP)

	// Make a few keys
	contractCodeKey := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameContractCode,
				Origin: models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)
	contractYearNameKey := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameContractYearName,
				Origin: models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	weightEstimatedKey := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameWeightEstimated,
				Origin: models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)

	// Make the service param associations
	var serviceKeysAssociation = []struct {
		service models.ReService
		keys    models.ServiceItemParamKeys
	}{
		{dlhService, models.ServiceItemParamKeys{contractCodeKey, contractYearNameKey}},
		{dopService, models.ServiceItemParamKeys{contractYearNameKey, weightEstimatedKey}},
	}

	for _, serviceKey := range serviceKeysAssociation {
		for _, key := range serviceKey.keys {
			factory.BuildServiceParam(suite.DB(), []factory.Customization{
				{
					Model:    serviceKey.service,
					LinkOnly: true,
				},
				{
					Model:    key,
					LinkOnly: true,
				},
			}, nil)
		}
	}

	// Make an MTO service item for DLH
	dlhServiceItem := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model:    dlhService,
			LinkOnly: true,
		},
	}, nil)

	helper := RequestPaymentHelper{}
	serviceParams, err := helper.FetchServiceParamList(suite.AppContextForTest(), dlhServiceItem)
	suite.NoError(err)

	// We should get back only the contract code key since the contract year has origin of pricer
	if suite.Len(serviceParams, 1) {
		suite.Equal(dlhService.ID, serviceParams[0].ServiceID)
		suite.Equal(contractCodeKey.ID, serviceParams[0].ServiceItemParamKeyID)
		// Make sure we can read something off the ServiceItemParamKey association since it should have loaded
		suite.Equal(contractCodeKey.Key, serviceParams[0].ServiceItemParamKey.Key)
	}
}
