package paymentrequest

import (
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
			Key:         "weight",
			Description: "actual weight",
			Type:        models.ServiceItemParamTypeInteger,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})
	serviceItemParamKey2 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         "pickup",
			Description: "requested pickup date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	creator := NewPaymentRequestCreator(suite.DB())

	suite.T().Run("Payment request is created successfully (using IncomingKey)", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "weight",
							Value:       "3254",
						},
						{
							IncomingKey: "pickup",
							Value:       "2019-12-16",
						},
					},
				},
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "weight",
							Value:       "7722",
						},
					},
				},
			},
		}

		_, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.FatalNoError(err)

		// Verify some of the data that came back
		suite.NotEqual(paymentRequest.ID, uuid.Nil)
		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequest.PaymentServiceItems[0].ID, uuid.Nil)
			if suite.Len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequest.PaymentServiceItems[1].ID, uuid.Nil)
			if suite.Len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
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

		// Verify some of the data that came back
		suite.NotEqual(paymentRequest.ID, uuid.Nil)
		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequest.PaymentServiceItems[0].ID, uuid.Nil)
			if suite.Len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequest.PaymentServiceItems[1].ID, uuid.Nil)
			if suite.Len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
			}
		}
	})
	/*
	       TODO add in new test for being able to pull out params when none are provided
	   	suite.T().Run("Payment request is created successfully (using no IncomingKey data or ServiceItemParamKeyID data)", func(t *testing.T) {
	   		paymentRequest := models.PaymentRequest{
	   			MoveTaskOrderID: moveTaskOrder.ID,
	   			IsFinal:         false,
	   		}

	   		_, err := creator.CreatePaymentRequest(&paymentRequest)
	   		suite.FatalNoError(err)

	   		// Verify some of the data that came back
	   		suite.NotEqual(paymentRequest.ID, uuid.Nil)
	   		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
	   			suite.NotEqual(paymentRequest.PaymentServiceItems[0].ID, uuid.Nil)
	   			if suite.Len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
	   				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
	   				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
	   			}
	   			suite.NotEqual(paymentRequest.PaymentServiceItems[1].ID, uuid.Nil)
	   			if suite.Len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
	   				suite.NotEqual(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
	   			}
	   		}
	   	})

	*/

	badID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")

	suite.T().Run("Given a non-existent move task order id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: badID,
			IsFinal:         false,
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)
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
	})

	suite.T().Run("Given a non-existent service item param key id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
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
	})

	suite.T().Run("Given a non-existent service item param key name, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
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
	})
}
