package paymentrequest

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {
	// Create some records we'll need to link to
	moveTaskOrder := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: moveTaskOrder,
		ReService: models.ReService{
			Code: "DLH",
		},
	})
	mtoServiceItem2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: moveTaskOrder,
		ReService: models.ReService{
			Code: "DOP",
		},
	})
	serviceItemParamKey1 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         "WeightEstimated",
			Description: "estimated weight",
			Type:        models.ServiceItemParamTypeInteger,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})
	serviceItemParamKey2 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         "RequestedPickupDate",
			Description: "requested pickup date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
		},
	})

	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey2.ID,
			ServiceItemParamKey:   serviceItemParamKey2,
		},
	})

	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem2.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
		},
	})

	/*
		_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID: mtoServiceItem2.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey2.ID,
				ServiceItemParamKey: serviceItemParamKey2,
			},
		})
	*/

	creator := NewPaymentRequestCreator(suite.DB())

	suite.T().Run("Payment request is created successfully (using IncomingKey)", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "WeightEstimated",
							Value:       "3254",
						},
						{
							IncomingKey: "RequestedPickupDate",
							Value:       "2019-12-16",
						},
					},
				},
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "WeightEstimated",
							Value:       "7722",
						},
					},
				},
			},
		}

		paymentRequestReturn, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.FatalNoError(err)

		//fmt.Printf("\n\npaymentRequestReturn\n%v\n\n", paymentRequestReturn)
		for _, item := range paymentRequestReturn.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}

		// Verify some of the data that came back
		suite.NotEqual(paymentRequestReturn.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequestReturn.PaymentServiceItems), "PaymentServiceItems expect 2")
		if suite.Len(paymentRequestReturn.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequestReturn.PaymentServiceItems[0].ID, uuid.Nil)
			suite.Equal(2, len(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItemParams expect 2")
			if suite.Len(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequestReturn.PaymentServiceItems[1].ID, uuid.Nil)
			suite.Equal(1, len(paymentRequestReturn.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 1")
			if suite.Len(paymentRequestReturn.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
				suite.NotEqual(paymentRequestReturn.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
			}
		}
	})

	suite.T().Run("Payment request is created successfully (using ServiceItemParamKeyID)", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: serviceItemParamKey1.ID,
							Value:                 "3254",
						},
						{
							ServiceItemParamKeyID: serviceItemParamKey2.ID,
							Value:                 "2019-12-16",
						},
					},
				},
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: serviceItemParamKey1.ID,
							Value:                 "7722",
						},
					},
				},
			},
		}

		_, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.FatalNoError(err)

		for _, item := range paymentRequest.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}

		// Verify some of the data that came back
		suite.NotEqual(paymentRequest.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequest.PaymentServiceItems), "PaymentServiceItems expect 2")
		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequest.PaymentServiceItems[0].ID, uuid.Nil)
			suite.Equal(2, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItemParams expect 2")
			if suite.Len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequest.PaymentServiceItems[1].ID, uuid.Nil)
			suite.Equal(1, len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 1")
			if suite.Len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
			}
		}
	})

	suite.T().Run("Payment request is created successfully (using no IncomingKey data or ServiceItemParamKeyID data)", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID:         mtoServiceItem1.ID,
					MTOServiceItem:           mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{},
				},
				{
					MTOServiceItemID:         mtoServiceItem2.ID,
					MTOServiceItem:           mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{},
				},
			},
		}

		paymentRequestResult, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.FatalNoError(err)

		for _, item := range paymentRequest.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}

		// Verify some of the data that came back
		suite.NotEqual(paymentRequestResult.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequest.PaymentServiceItems), "PaymentServiceItems expect 2")
		if suite.Len(paymentRequestResult.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequestResult.PaymentServiceItems[0].ID, uuid.Nil)
			suite.Equal(2, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItemParams expect 2")
			if suite.Len(paymentRequestResult.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequestResult.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequestResult.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequestResult.PaymentServiceItems[1].ID, uuid.Nil)
			suite.Equal(1, len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 1")
			if suite.Len(paymentRequestResult.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
				suite.NotEqual(paymentRequestResult.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
			}
		}
	})

	badID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")

	suite.T().Run("Given a non-existent move task order id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: badID,
			IsFinal:         false,
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)

		for _, item := range invalidPaymentRequest.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}
	})

	suite.T().Run("Given a non-existent service item id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: badID,
				},
			},
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)

		for _, item := range invalidPaymentRequest.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}
	})

	suite.T().Run("Given a non-existent service item param key id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: badID,
							Value:                 "3254",
						},
					},
				},
			},
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)

		for _, item := range invalidPaymentRequest.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}
	})

	suite.T().Run("Given a non-existent service item param key name, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "bogus",
							Value:       "3254",
						},
					},
				},
			},
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)

		for _, item := range invalidPaymentRequest.PaymentServiceItems {
			fmt.Printf("********\n\nMTO Service Item: <%s>\n", item.MTOServiceItem.ReService.Code)
			fmt.Printf("MTO Service Item ID: <%s>\n", item.MTOServiceItem.ID.String())

			for _, param := range item.PaymentServiceItemParams {
				fmt.Printf("param: incoming key <%s>\n", param.IncomingKey)
				fmt.Printf("param: key <%s>\n", param.ServiceItemParamKey.Key)
			}
		}
	})
}
