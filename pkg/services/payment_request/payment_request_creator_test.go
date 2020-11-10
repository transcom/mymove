package paymentrequest

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {
	// Create some records we'll need to link to
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	estimatedWeight := unit.Pound(2048)
	mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		ReService: models.ReService{
			Code: "DLH",
		},
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
		},
	})
	mtoServiceItem2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		ReService: models.ReService{
			Code: "DOP",
		},
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
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
			Key:         models.ServiceItemParamNameCanStandAlone,
			Description: "can stand alone",
			Type:        models.ServiceItemParamTypeString,
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

	testPrice := unit.Cents(12345)
	serviceItemPricer := &mocks.ServiceItemPricer{}
	serviceItemPricer.
		On("PriceServiceItem", mock.Anything).Return(testPrice, nil).
		On("UsingConnection", mock.Anything).Return(serviceItemPricer)

	planner := &routemocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.Anything,
		mock.Anything,
	).Return(0, nil)
	creator := NewPaymentRequestCreator(suite.DB(), planner, serviceItemPricer)

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
							IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
							Value:       "3254",
						},
						{
							IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
							Value:       "2019-12-16",
						},
					},
				},
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
							Value:       "7722",
						},
					},
				},
			},
		}

		paymentRequestReturn, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.FatalNoError(err)

		expectedSequenceNumber := 1
		expectedPaymentRequestNumber := fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, expectedSequenceNumber)
		// Verify some of the data that came back
		suite.Equal(expectedPaymentRequestNumber, paymentRequestReturn.PaymentRequestNumber)
		suite.Equal(expectedSequenceNumber, paymentRequestReturn.SequenceNumber)
		suite.NotEqual(paymentRequestReturn.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequestReturn.PaymentServiceItems), "PaymentServiceItems expect 2")
		if suite.Len(paymentRequestReturn.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequestReturn.PaymentServiceItems[0].ID, uuid.Nil)
			suite.Equal(*paymentRequestReturn.PaymentServiceItems[0].PriceCents, testPrice)
			suite.Equal(2, len(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItemParams expect 2")
			if suite.Len(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequestReturn.PaymentServiceItems[1].ID, uuid.Nil)
			suite.Equal(*paymentRequestReturn.PaymentServiceItems[1].PriceCents, testPrice)
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
						{
							ServiceItemParamKeyID: serviceItemParamKey3.ID,
							Value:                 "foobar",
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

		// Verify some of the data that came back
		suite.NotEqual(paymentRequest.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequest.PaymentServiceItems), "PaymentServiceItems expect 2")
		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequest.PaymentServiceItems[0].ID, uuid.Nil)
			suite.Equal(*paymentRequest.PaymentServiceItems[0].PriceCents, testPrice)
			suite.Equal(3, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItemParams expect 3")
			if suite.Len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams, 3) {
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
				suite.NotEqual(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams[2].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequest.PaymentServiceItems[1].ID, uuid.Nil)
			suite.Equal(*paymentRequest.PaymentServiceItems[1].PriceCents, testPrice)
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

		// Verify some of the data that came back
		suite.NotEqual(paymentRequestResult.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequest.PaymentServiceItems), "PaymentServiceItems expect 2")
		if suite.Len(paymentRequestResult.PaymentServiceItems, 2) {
			suite.NotEqual(paymentRequestResult.PaymentServiceItems[0].ID, uuid.Nil)
			suite.Equal(*paymentRequestResult.PaymentServiceItems[0].PriceCents, testPrice)
			suite.Equal(2, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItemParams expect 2")
			if suite.Len(paymentRequestResult.PaymentServiceItems[0].PaymentServiceItemParams, 2) {
				suite.NotEqual(paymentRequestResult.PaymentServiceItems[0].PaymentServiceItemParams[0].ID, uuid.Nil)
				suite.NotEqual(paymentRequestResult.PaymentServiceItems[0].PaymentServiceItemParams[1].ID, uuid.Nil)
			}
			suite.NotEqual(paymentRequestResult.PaymentServiceItems[1].ID, uuid.Nil)
			suite.Equal(*paymentRequestResult.PaymentServiceItems[1].PriceCents, testPrice)
			suite.Equal(1, len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 1")
			if suite.Len(paymentRequestResult.PaymentServiceItems[1].PaymentServiceItemParams, 1) {
				suite.NotEqual(paymentRequestResult.PaymentServiceItems[1].PaymentServiceItemParams[0].ID, uuid.Nil)
			}
		}
	})

	suite.T().Run("Payment request fails when pricing", func(t *testing.T) {
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

		errMsg := "pricing failed"
		failingServiceItemPricer := &mocks.ServiceItemPricer{}
		failingServiceItemPricer.
			On("PriceServiceItem", mock.Anything).Return(unit.Cents(0), errors.New(errMsg)).
			On("UsingConnection", mock.Anything).Return(failingServiceItemPricer)

		planner := &routemocks.Planner{}
		planner.On("Zip5TransitDistanceLineHaul",
			mock.Anything,
			mock.Anything,
		).Return(0, nil)
		failingCreator := NewPaymentRequestCreator(suite.DB(), planner, failingServiceItemPricer)

		_, err := failingCreator.CreatePaymentRequest(&paymentRequest)
		suite.Error(err)
		wrappedErr := errors.Unwrap(err)
		suite.Equal(errMsg, wrappedErr.Error())
	})

	suite.T().Run("Given a non-existent move task order id, the create should fail", func(t *testing.T) {
		badID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: badID,
			IsFinal:         false,
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)

		suite.Error(err)
		_, ok := err.(services.NotFoundError)
		suite.Equal(true, ok)
	})

	suite.T().Run("Given no move task order id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Nil,
			IsFinal:         false,
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)

		suite.Error(err)
		_, ok := err.(services.InvalidCreateInputError)
		suite.Equal(true, ok)
		suite.Equal("Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create", err.Error())
	})

	suite.T().Run("Given move with orders but no LOA, the create should fail", func(t *testing.T) {
		mtoInvalidOrders := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		orders := mtoInvalidOrders.Orders
		orders.TAC = nil
		suite.MustSave(&orders)
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: mtoInvalidOrders.ID,
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
		_, err := creator.CreatePaymentRequest(&paymentRequest)

		suite.Error(err)
		_, ok := err.(*services.BadDataError)
		suite.Equal(true, ok)
		suite.Equal(fmt.Sprintf("Data received from requester is bad: BAD_DATA: MoveTaskOrder (ID: %s) Orders (ID: %s) missing Lines of Accounting TAC", mtoInvalidOrders.ID, mtoInvalidOrders.OrdersID), err.Error())
	})

	suite.T().Run("Given move with orders blank LOA, the create should fail", func(t *testing.T) {
		mtoInvalidOrders := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		orders := mtoInvalidOrders.Orders
		blankTAC := ""
		orders.TAC = &blankTAC
		err := suite.DB().Update(&orders)
		suite.FatalNoError(err)
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: mtoInvalidOrders.ID,
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
		_, err = creator.CreatePaymentRequest(&paymentRequest)

		suite.Error(err)
		_, ok := err.(*services.BadDataError)
		suite.Equal(true, ok)
		suite.Equal(fmt.Sprintf("Data received from requester is bad: BAD_DATA: MoveTaskOrder (ID: %s) Orders (ID: %s) missing Lines of Accounting TAC", mtoInvalidOrders.ID, mtoInvalidOrders.OrdersID), err.Error())
	})

	suite.T().Run("Given move with orders no OriginDutyStation, the create should fail", func(t *testing.T) {
		mtoInvalidOrders := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		orders := mtoInvalidOrders.Orders
		orders.OriginDutyStation = nil
		orders.OriginDutyStationID = nil
		err := suite.DB().Update(&orders)
		suite.FatalNoError(err)
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: mtoInvalidOrders.ID,
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
		_, err = creator.CreatePaymentRequest(&paymentRequest)

		suite.Error(err)
		_, ok := err.(*services.BadDataError)
		suite.Equal(true, ok)
		suite.Equal(fmt.Sprintf("Data received from requester is bad: BAD_DATA: MoveTaskOrder (ID: %s) Orders (ID: %s) missing OriginDutyStation", mtoInvalidOrders.ID, mtoInvalidOrders.OrdersID), err.Error())
	})

	suite.T().Run("Given a non-existent service item id, the create should fail", func(t *testing.T) {
		badID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
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
		_, ok := err.(services.NotFoundError)
		suite.Equal(true, ok)
	})

	suite.T().Run("Given a non-existent service item param key id, the create should fail", func(t *testing.T) {
		badID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
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
		_, ok := err.(services.NotFoundError)
		suite.Equal(true, ok)
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
		_, ok := err.(*services.BadDataError)
		suite.Equal(true, ok)
	})

	suite.T().Run("Payment request numbers increment by 1", func(t *testing.T) {
		// Determine the max sequence number we already have for this MTO ID
		var max int
		err := suite.DB().RawQuery("SELECT COALESCE(MAX(sequence_number),0) FROM payment_requests WHERE move_id = $1", moveTaskOrder.ID).First(&max)
		suite.FatalNoError(err)

		// Create two new ones
		paymentRequest1 := models.PaymentRequest{
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
					},
				},
			},
		}
		_, err = creator.CreatePaymentRequest(&paymentRequest1)
		suite.FatalNoError(err)

		paymentRequest2 := models.PaymentRequest{
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
					},
				},
			},
		}
		_, err = creator.CreatePaymentRequest(&paymentRequest2)
		suite.FatalNoError(err)

		// Verify expected payment request numbers
		expectedSequenceNumber1 := max + 1
		expectedPaymentRequestNumber1 := fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, expectedSequenceNumber1)
		suite.Equal(expectedPaymentRequestNumber1, paymentRequest1.PaymentRequestNumber)
		suite.Equal(expectedSequenceNumber1, paymentRequest1.SequenceNumber)

		expectedSequenceNumber2 := max + 2
		expectedPaymentRequestNumber2 := fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, expectedSequenceNumber2)
		suite.Equal(expectedPaymentRequestNumber2, paymentRequest2.PaymentRequestNumber)
		suite.Equal(expectedSequenceNumber2, paymentRequest2.SequenceNumber)
	})
}
