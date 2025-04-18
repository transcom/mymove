package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

type serviceParamSubtestData struct {
	move                     models.Move
	mtoServiceItem1          models.MTOServiceItem
	mtoServiceItem2          models.MTOServiceItem
	mtoService1ServiceParams models.ServiceParams
	mtoService2ServiceParams models.ServiceParams
}

func (suite *PaymentRequestHelperSuite) makeServiceParamTestData() (subtestData *serviceParamSubtestData) {
	subtestData = &serviceParamSubtestData{}
	// Create some records we'll need to link to
	subtestData.move = factory.BuildMove(suite.DB(), nil, nil)
	subtestData.mtoServiceItem1 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDLH,
				Name: "Domestic Linehaul",
			},
		},
	}, nil)

	subtestData.mtoServiceItem2 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOP,
				Name: "Domestic Origin Pickup",
			},
		},
	}, nil)

	serviceItemParamKey1 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameWeightEstimated,
				Description: "estimated weight",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)
	serviceItemParamKey2 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameRequestedPickupDate,
				Description: "requested pickup date",
				Type:        models.ServiceItemParamTypeDate,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)
	serviceItemParamKey3 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameZipPickupAddress,
				Description: "ZIP Pickup Address",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)

	serviceItemParamKey4 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameDistanceZip,
				Description: "ZIP Distance",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)

	mtoServiceItem1Param1 := factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItem1.ReService,
			LinkOnly: true,
		},
		{
			Model:    serviceItemParamKey1,
			LinkOnly: true,
		},
		{
			Model: models.ServiceParam{
				IsOptional: true,
			},
		},
	}, nil)

	mtoServiceItem1Param2 := factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItem1.ReService,
			LinkOnly: true,
		},
		{
			Model:    serviceItemParamKey2,
			LinkOnly: true,
		},
	}, nil)

	mtoServiceItem1Param3 := factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItem1.ReService,
			LinkOnly: true,
		},
		{
			Model:    serviceItemParamKey3,
			LinkOnly: true,
		},
	}, nil)

	mtoServiceItem1Param4 := factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItem1.ReService,
			LinkOnly: true,
		},
		{
			Model:    serviceItemParamKey4,
			LinkOnly: true,
		},
	}, nil)

	mtoServiceItem2Param1 := factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItem2.ReService,
			LinkOnly: true,
		},
		{
			Model:    serviceItemParamKey1,
			LinkOnly: true,
		},
		{
			Model: models.ServiceParam{
				IsOptional: true,
			},
		},
	}, nil)

	subtestData.mtoService1ServiceParams = models.ServiceParams{
		mtoServiceItem1Param1,
		mtoServiceItem1Param2,
		mtoServiceItem1Param3,
		mtoServiceItem1Param4,
	}

	subtestData.mtoService2ServiceParams = models.ServiceParams{
		mtoServiceItem2Param1,
	}

	return subtestData

}

func (suite *PaymentRequestHelperSuite) TestValidServiceParamList() {
	suite.Run("Validate Service Items Params is TRUE (All params present)", func() {
		subtestData := suite.makeServiceParamTestData()

		// This set of service params has required and optional params with no value set (should be valid).
		var paymentServiceItemParams1 models.PaymentServiceItemParams
		for _, p := range subtestData.mtoService1ServiceParams {
			paymentServiceItemParams1 = append(paymentServiceItemParams1,
				models.PaymentServiceItemParam{
					ServiceItemParamKeyID: p.ServiceItemParamKeyID,
					ServiceItemParamKey:   p.ServiceItemParamKey,
				})
		}
		var paymentServiceItemParams2 models.PaymentServiceItemParams
		// This set of service params has only one optional param, and we're setting a value (should be valid).
		for _, p := range subtestData.mtoService2ServiceParams {
			paymentServiceItemParams2 = append(paymentServiceItemParams2,
				models.PaymentServiceItemParam{
					ServiceItemParamKeyID: p.ServiceItemParamKeyID,
					ServiceItemParamKey:   p.ServiceItemParamKey,
					Value:                 "1000",
				})
		}
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: subtestData.move.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID:         subtestData.mtoServiceItem1.ID,
					MTOServiceItem:           subtestData.mtoServiceItem1,
					PaymentServiceItemParams: paymentServiceItemParams1,
				},
				{
					MTOServiceItemID:         subtestData.mtoServiceItem2.ID,
					MTOServiceItem:           subtestData.mtoServiceItem2,
					PaymentServiceItemParams: paymentServiceItemParams2,
				},
			},
		}

		paymentHelper := RequestPaymentHelper{}
		validParamList1, validateMessage1 := paymentHelper.ValidServiceParamList(subtestData.mtoServiceItem1, subtestData.mtoService1ServiceParams, paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams)
		suite.Equal(true, validParamList1, "All params for service item should be present")
		suite.Empty(validateMessage1, "No error message returned")
		validParamList2, validateMessage2 := paymentHelper.ValidServiceParamList(subtestData.mtoServiceItem2, subtestData.mtoService2ServiceParams, paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams)
		suite.Equal(true, validParamList2, "All params for service item should be present")
		suite.Empty(validateMessage2, "No error message returned")
	})

	suite.Run("Validate Service Items Params is FALSE (Params are missing)", func() {
		subtestData := suite.makeServiceParamTestData()

		// for service item 1, deliberately supply NO payment parameters
		paymentServiceItemParams1 := models.PaymentServiceItemParams{}

		var paymentServiceItemParams2 models.PaymentServiceItemParams
		for _, p := range subtestData.mtoService2ServiceParams {
			paymentServiceItemParams2 = append(paymentServiceItemParams2, models.PaymentServiceItemParam{
				ServiceItemParamKeyID: p.ServiceItemParamKeyID,
				ServiceItemParamKey:   p.ServiceItemParamKey,
				Value:                 "1000",
			})
		}

		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: subtestData.move.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID:         subtestData.mtoServiceItem1.ID,
					MTOServiceItem:           subtestData.mtoServiceItem1,
					PaymentServiceItemParams: paymentServiceItemParams1,
				},
				{
					MTOServiceItemID:         subtestData.mtoServiceItem2.ID,
					MTOServiceItem:           subtestData.mtoServiceItem2,
					PaymentServiceItemParams: paymentServiceItemParams2,
				},
			},
		}

		paymentHelper := RequestPaymentHelper{}
		validParamList1, validateMessage1 := paymentHelper.ValidServiceParamList(
			subtestData.mtoServiceItem1,
			subtestData.mtoService1ServiceParams,
			paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams,
		)
		suite.Equal(false, validParamList1, "Expected validation to fail for service item 1 since required params are missing")
		suite.NotEmpty(validateMessage1, "Expected error message listing missing param keys")

		validParamList2, validateMessage2 := paymentHelper.ValidServiceParamList(
			subtestData.mtoServiceItem2,
			subtestData.mtoService2ServiceParams,
			paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams,
		)
		suite.Equal(true, validParamList2, "Expected validation to succeed for service item 2")
		suite.Empty(validateMessage2, "No error message should be returned for a complete set")
	})

}
