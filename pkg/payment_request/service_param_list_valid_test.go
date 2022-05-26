package paymentrequest

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
	subtestData.move = testdatagen.MakeDefaultMove(suite.DB())
	subtestData.mtoServiceItem1 = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: subtestData.move,
		ReService: models.ReService{
			Code: models.ReServiceCodeDLH,
			Name: "Domestic Linehaul",
		},
	})
	subtestData.mtoServiceItem2 = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: subtestData.move,
		ReService: models.ReService{
			Code: models.ReServiceCodeDOP,
			Name: "Domestic Origin Pickup",
		},
	})
	serviceItemParamKey1 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameWeightEstimated,
			Description: "estimated weight",
			Type:        models.ServiceItemParamTypeInteger,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})
	serviceItemParamKey2 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameRequestedPickupDate,
			Description: "requested pickup date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})
	serviceItemParamKey3 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameZipPickupAddress,
			Description: "ZIP Pickup Address",
			Type:        models.ServiceItemParamTypeString,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	serviceItemParamKey4 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameDistanceZip,
			Description: "ZIP Distance",
			Type:        models.ServiceItemParamTypeString,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	mtoServiceItem1Param1 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
			IsOptional:            true,
		},
	})

	mtoServiceItem1Param2 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey2.ID,
			ServiceItemParamKey:   serviceItemParamKey2,
		},
	})

	mtoServiceItem1Param3 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey3.ID,
			ServiceItemParamKey:   serviceItemParamKey3,
		},
	})

	mtoServiceItem1Param4 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey4.ID,
			ServiceItemParamKey:   serviceItemParamKey4,
		},
	})

	mtoServiceItem2Param1 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItem2.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
			IsOptional:            true,
		},
	})

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
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: subtestData.move.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: subtestData.mtoServiceItem1.ID,
					MTOServiceItem:   subtestData.mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: subtestData.mtoService1ServiceParams[0].ServiceItemParamKeyID,
							ServiceItemParamKey:   subtestData.mtoService1ServiceParams[0].ServiceItemParamKey,
						},
					},
				},
				{
					MTOServiceItemID: subtestData.mtoServiceItem2.ID,
					MTOServiceItem:   subtestData.mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: subtestData.mtoService2ServiceParams[0].ServiceItemParamKeyID,
							ServiceItemParamKey:   subtestData.mtoService2ServiceParams[0].ServiceItemParamKey,
						},
					},
				},
			},
		}

		paymentHelper := RequestPaymentHelper{}
		validParamList1, validateMessage1 := paymentHelper.ValidServiceParamList(subtestData.mtoServiceItem1, subtestData.mtoService1ServiceParams, paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams)
		suite.Equal(false, validParamList1, "All params for service item should be present")
		suite.NotEmpty(validateMessage1, "Error message with list of missing param keys")
		validParamList2, validateMessage2 := paymentHelper.ValidServiceParamList(subtestData.mtoServiceItem2, subtestData.mtoService2ServiceParams, paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams)
		suite.Equal(true, validParamList2, "All params for service item should be present")
		suite.Empty(validateMessage2, "No error message returned")
	})
}
