package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestHelperSuite) TestValidServiceParamList() {
	// Create some records we'll need to link to
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		ReService: models.ReService{
			Code: "DLH",
			Name: "Domestic Linehaul",
		},
	})
	mtoServiceItem2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		ReService: models.ReService{
			Code: "DOP",
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
			Key:         models.ServiceItemParamNameDistanceZip3,
			Description: "ZIP 3 Distance",
			Type:        models.ServiceItemParamTypeString,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	mtoServiceItem1Param1 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
		},
	})

	mtoServiceItem1Param2 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey2.ID,
			ServiceItemParamKey:   serviceItemParamKey2,
		},
	})

	mtoServiceItem1Param3 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey3.ID,
			ServiceItemParamKey:   serviceItemParamKey3,
		},
	})

	mtoServiceItem1Param4 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey4.ID,
			ServiceItemParamKey:   serviceItemParamKey4,
		},
	})

	mtoServiceItem2Param1 := testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem2.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
		},
	})

	mtoService1ServiceParams := models.ServiceParams{
		mtoServiceItem1Param1,
		mtoServiceItem1Param2,
		mtoServiceItem1Param3,
		mtoServiceItem1Param4,
	}

	mtoService2ServiceParams := models.ServiceParams{
		mtoServiceItem2Param1,
	}

	suite.T().Run("Validate Service Items Params is TRUE (All params present)", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: mtoServiceItem1Param1.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem1Param1.ServiceItemParamKey,
						},
						{
							ServiceItemParamKeyID: mtoServiceItem1Param2.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem1Param2.ServiceItemParamKey,
						},
						{
							ServiceItemParamKeyID: mtoServiceItem1Param3.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem1Param3.ServiceItemParamKey,
						},
						{
							ServiceItemParamKeyID: mtoServiceItem1Param4.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem1Param4.ServiceItemParamKey,
						},
					},
				},
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: mtoServiceItem2Param1.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem2Param1.ServiceItemParamKey,
						},
					},
				},
			},
		}

		paymentHelper := RequestPaymentHelper{DB: suite.DB()}
		validParamList1, validateMessage1 := paymentHelper.ValidServiceParamList(mtoServiceItem1, mtoService1ServiceParams, paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams)
		suite.Equal(true, validParamList1, "All params for service item should be present")
		suite.Empty(validateMessage1, "No error message returned")
		validParamList2, validateMessage2 := paymentHelper.ValidServiceParamList(mtoServiceItem2, mtoService2ServiceParams, paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams)
		suite.Equal(true, validParamList2, "All params for service item should be present")
		suite.Empty(validateMessage2, "No error message returned")
	})

	suite.T().Run("Validate Service Items Params is FALSE (Params are missing)", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: mtoServiceItem1Param1.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem1Param1.ServiceItemParamKey,
						},
					},
				},
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: mtoServiceItem2Param1.ServiceItemParamKeyID,
							ServiceItemParamKey:   mtoServiceItem2Param1.ServiceItemParamKey,
						},
					},
				},
			},
		}

		paymentHelper := RequestPaymentHelper{DB: suite.DB()}
		validParamList1, validateMessage1 := paymentHelper.ValidServiceParamList(mtoServiceItem1, mtoService1ServiceParams, paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams)
		suite.Equal(false, validParamList1, "All params for service item should be present")
		suite.NotEmpty(validateMessage1, "Error message with list of missing param keys")
		validParamList2, validateMessage2 := paymentHelper.ValidServiceParamList(mtoServiceItem2, mtoService2ServiceParams, paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams)
		suite.Equal(true, validParamList2, "All params for service item should be present")
		suite.Empty(validateMessage2, "No error message returned")
	})
}
