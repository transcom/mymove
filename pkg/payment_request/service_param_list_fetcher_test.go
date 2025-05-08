package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestHelperSuite) TestFetchServiceParamList() {
	// Make a couple of services
	dlhService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
	dopService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)

	// Make a few keys
	contractCodeKey := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameContractCode,
				Origin: models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)
	contractYearNameKey := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameContractYearName,
				Origin: models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	weightEstimatedKey := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
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
			factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
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

	suite.Len(serviceParams, 13)

	// iterate through serviceParams to find the one with the contractCodeKey
	var foundParam *models.ServiceParam
	for i := range serviceParams {
		if serviceParams[i].ServiceItemParamKeyID == contractCodeKey.ID {
			foundParam = &serviceParams[i]
			break
		}
	}
	suite.NotNil(foundParam, "Expected to find a service param with the contract code key")

	suite.Equal(dlhService.ID, foundParam.ServiceID)
	suite.Equal(contractCodeKey.ID, foundParam.ServiceItemParamKeyID)
	suite.Equal(contractCodeKey.Key, foundParam.ServiceItemParamKey.Key)

}
