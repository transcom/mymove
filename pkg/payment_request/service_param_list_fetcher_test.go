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
	requestedPickupDateKey := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameRequestedPickupDate,
				Origin: models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)
	requestedPickupDateNameKey := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNameRequestedPickupDate,
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
		{dlhService, models.ServiceItemParamKeys{requestedPickupDateKey, requestedPickupDateNameKey}},
		{dopService, models.ServiceItemParamKeys{requestedPickupDateNameKey, weightEstimatedKey}},
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

	// We should get back all params applicable to DLH where origin is not PRICER
	if suite.Len(serviceParams, 13) {
		suite.Equal(dlhService.ID, serviceParams[0].ServiceID)
		suite.Equal(requestedPickupDateKey.ID, serviceParams[0].ServiceItemParamKeyID)
		// Make sure we can read something off the ServiceItemParamKey association since it should have loaded
		suite.Equal(requestedPickupDateKey.Key, serviceParams[0].ServiceItemParamKey.Key)
	}
}
