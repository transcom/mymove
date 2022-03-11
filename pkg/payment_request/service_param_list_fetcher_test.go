package paymentrequest

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestHelperSuite) TestFetchServiceParamList() {
	// Make a couple of services
	dlhService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDLH,
		},
	})
	dopService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOP,
		},
	})

	// Make a few keys
	contractCodeKey := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:    models.ServiceItemParamNameContractCode,
			Origin: models.ServiceItemParamOriginSystem,
		},
	})
	contractYearNameKey := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:    models.ServiceItemParamNameContractYearName,
			Origin: models.ServiceItemParamOriginPricer,
		},
	})
	weightEstimatedKey := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:    models.ServiceItemParamNameWeightEstimated,
			Origin: models.ServiceItemParamOriginPrime,
		},
	})

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
			testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
				ServiceParam: models.ServiceParam{
					ServiceID:             serviceKey.service.ID,
					ServiceItemParamKeyID: key.ID,
				},
			})
		}
	}

	// Make an MTO service item for DLH
	dlhServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		ReService: dlhService,
		Stub:      true,
	})

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
