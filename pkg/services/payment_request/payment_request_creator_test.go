package paymentrequest

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {
	var moveTaskOrder models.Move
	var mtoServiceItem1, mtoServiceItem2, mtoServiceItem3, mtoServiceItemSubmitted, mtoServiceItemRejected models.MTOServiceItem
	var serviceItemParamKey1, serviceItemParamKey2, serviceItemParamKey3 models.ServiceItemParamKey
	var displayParams models.PaymentServiceItemParams

	testPrice := unit.Cents(12345)

	serviceItemPricer := &mocks.ServiceItemPricer{}
	planner := &routemocks.Planner{}
	creator := NewPaymentRequestCreator(planner, serviceItemPricer)

	suite.PreloadData(func() {
		// Create some records we'll need to link to
		moveTaskOrder = factory.BuildMove(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				EndDate: time.Now().Add(time.Hour * 24),
			},
		})
		estimatedWeight := unit.Pound(2048)
		mtoServiceItem1 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
			},
		}, nil)
		mtoServiceItem2 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOP,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
			},
		}, nil)
		mtoServiceItem3 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOP,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					UsesExternalVendor:   true,
				},
			},
			{
				Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
			},
		}, nil)
		mtoServiceItemSubmitted = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOP,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					UsesExternalVendor:   true,
				},
			},
			{
				Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusSubmitted},
			},
		}, nil)
		mtoServiceItemRejected = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOP,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					UsesExternalVendor:   true,
				},
			},
			{
				Model: models.MTOServiceItem{Status: models.MTOServiceItemStatusRejected},
			},
		}, nil)
		serviceItemParamKey1 = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameWeightEstimated,
					Description: "estimated weight",
					Type:        models.ServiceItemParamTypeInteger,
					Origin:      models.ServiceItemParamOriginPrime,
				},
			},
		}, nil)
		serviceItemParamKey2 = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameRequestedPickupDate,
					Description: "requested pickup date",
					Type:        models.ServiceItemParamTypeDate,
					Origin:      models.ServiceItemParamOriginPrime,
				},
			},
		}, nil)
		serviceItemParamKey3 = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameZipPickupAddress,
					Description: "zip pickup address",
					Type:        models.ServiceItemParamTypeString,
					Origin:      models.ServiceItemParamOriginPrime,
				},
			},
		}, nil)

		serviceItemParamKey4 := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameEscalationCompounded,
					Description: "escalation factor",
					Type:        models.ServiceItemParamTypeDecimal,
					Origin:      models.ServiceItemParamOriginPricer,
				},
			},
		}, nil)

		serviceItemParamKey5 := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameContractYearName,
					Description: "contract year name",
					Type:        models.ServiceItemParamTypeString,
					Origin:      models.ServiceItemParamOriginPricer,
				},
			},
		}, nil)

		serviceItemParamKey6 := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameIsPeak,
					Description: "is peak",
					Type:        models.ServiceItemParamTypeBoolean,
					Origin:      models.ServiceItemParamOriginPricer,
				},
			},
		}, nil)

		serviceItemParamKey7 := factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNamePriceRateOrFactor,
					Description: "Price, rate, or factor used in calculation",
					Type:        models.ServiceItemParamTypeDecimal,
					Origin:      models.ServiceItemParamOriginPricer,
				},
			},
		}, nil)

		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
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
		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey2,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey4,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey5,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey6,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem1.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey7,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem2.ReService,
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

		displayParams = models.PaymentServiceItemParams{
			{
				ID:                    uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
				ServiceItemParamKeyID: serviceItemParamKey4.ID,
				ServiceItemParamKey:   serviceItemParamKey4,
				Value:                 "1.000",
			},
			{
				ID:                    uuid.FromStringOrNil("d55d2f35-218c-4b85-b9d1-631949b9d984"),
				ServiceItemParamKeyID: serviceItemParamKey5.ID,
				ServiceItemParamKey:   serviceItemParamKey5,
				Value:                 "Base Contract Year 1",
			},
			{
				ID:                    uuid.FromStringOrNil("d44d2f35-218c-4b85-b9d1-631949b9d984"),
				ServiceItemParamKeyID: serviceItemParamKey6.ID,
				ServiceItemParamKey:   serviceItemParamKey6,
				Value:                 "true",
			},
			{
				ID:                    uuid.FromStringOrNil("d22d2f35-218c-4b85-b9d1-631949b9d984"),
				ServiceItemParamKeyID: serviceItemParamKey7.ID,
				ServiceItemParamKey:   serviceItemParamKey7,
				Value:                 "333.2",
			},
		}
		serviceItemPricer.
			On("PriceServiceItem", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(testPrice, displayParams, nil)
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, nil)
	})

	suite.Run("Payment request is created successfully (using IncomingKey)", func() {
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

		paymentRequestReturn, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)
		suite.FatalNoError(err)

		expectedSequenceNumber := 1
		expectedPaymentRequestNumber := fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, expectedSequenceNumber)
		// Verify some of the data that came back
		suite.Equal(expectedPaymentRequestNumber, paymentRequestReturn.PaymentRequestNumber)
		suite.Equal(expectedSequenceNumber, paymentRequestReturn.SequenceNumber)
		suite.NotEqual(paymentRequestReturn.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequestReturn.PaymentServiceItems), "PaymentServiceItems expect 2")
		suite.Equal(6, len(paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 6")
		suite.Equal(5, len(paymentRequestReturn.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 5")

		if suite.Len(paymentRequestReturn.PaymentServiceItems, 2) {
			for _, paymentServiceItem := range paymentRequestReturn.PaymentServiceItems {
				var pricerDisplayParams models.PaymentServiceItemParams
				suite.NotEqual(paymentServiceItem.ID, uuid.Nil)
				suite.Equal(paymentServiceItem.PriceCents, &testPrice)
				for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
					suite.NotEqual(paymentServiceItemParam.ID, uuid.Nil)
					if paymentServiceItemParam.ServiceItemParamKey.Origin == models.ServiceItemParamOriginPricer {
						pricerDisplayParams = append(pricerDisplayParams, paymentServiceItemParam)
					}
				}
				for i, pricerDisplayParam := range pricerDisplayParams {
					suite.Equal(pricerDisplayParam.Value, displayParams[i].Value)
					suite.Equal(pricerDisplayParam.ServiceItemParamKeyID, displayParams[i].ServiceItemParamKeyID)
				}
			}
		}
	})

	suite.Run("Payment request is created successfully (using ServiceItemParamKeyID)", func() {
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
					PaymentRequest: models.PaymentRequest{
						ReviewedAt: models.TimePointer(time.Now()),
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

		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)
		suite.FatalNoError(err)

		// Verify some of the data that came back
		suite.NotEqual(paymentRequest.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequest.PaymentServiceItems), "PaymentServiceItems expect 2")
		suite.Equal(7, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 7")
		suite.Equal(5, len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 5")

		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
			for _, paymentServiceItem := range paymentRequest.PaymentServiceItems {
				var pricerDisplayParams models.PaymentServiceItemParams
				suite.NotEqual(paymentServiceItem.ID, uuid.Nil)
				suite.Equal(paymentServiceItem.PriceCents, &testPrice)
				for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
					suite.NotEqual(paymentServiceItemParam.ID, uuid.Nil)
					if paymentServiceItemParam.ServiceItemParamKey.Origin == models.ServiceItemParamOriginPricer {
						pricerDisplayParams = append(pricerDisplayParams, paymentServiceItemParam)
					}
				}
				for i, pricerDisplayParam := range pricerDisplayParams {
					suite.Equal(pricerDisplayParam.Value, displayParams[i].Value)
					suite.Equal(pricerDisplayParam.ServiceItemParamKeyID, displayParams[i].ServiceItemParamKeyID)
				}
			}
		}
	})

	suite.Run("Payment request is created successfully (using no IncomingKey data or ServiceItemParamKeyID data)", func() {
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

		paymentRequestResult, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)
		suite.FatalNoError(err)

		// Verify some of the data that came back
		suite.NotEqual(paymentRequestResult.ID, uuid.Nil)
		suite.NotEqual(paymentRequest.ID, uuid.Nil)
		suite.Equal(2, len(paymentRequest.PaymentServiceItems), "PaymentServiceItems expect 2")
		suite.Equal(6, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 6")
		suite.Equal(5, len(paymentRequest.PaymentServiceItems[1].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 5")

		if suite.Len(paymentRequest.PaymentServiceItems, 2) {
			for _, paymentServiceItem := range paymentRequest.PaymentServiceItems {
				var pricerDisplayParams models.PaymentServiceItemParams
				suite.NotEqual(paymentServiceItem.ID, uuid.Nil)
				suite.Equal(paymentServiceItem.PriceCents, &testPrice)
				for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
					suite.NotEqual(paymentServiceItemParam.ID, uuid.Nil)
					if paymentServiceItemParam.ServiceItemParamKey.Origin == models.ServiceItemParamOriginPricer {
						pricerDisplayParams = append(pricerDisplayParams, paymentServiceItemParam)
					}
				}
				for i, pricerDisplayParam := range pricerDisplayParams {
					suite.Equal(pricerDisplayParam.Value, displayParams[i].Value)
					suite.Equal(pricerDisplayParam.ServiceItemParamKeyID, displayParams[i].ServiceItemParamKeyID)
				}
			}
		}
	})

	suite.Run("Payment request fails when MTOShipment uses external vendor", func() {
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
				{
					MTOServiceItemID:         mtoServiceItem3.ID,
					MTOServiceItem:           mtoServiceItem3,
					PaymentServiceItemParams: models.PaymentServiceItemParams{},
				},
			},
		}

		paymentRequestResult, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "paymentRequestCreator.validShipment: Shipment uses external vendor for MTOShipmentID")

		// Verify some of the data that came back
		suite.Nil(paymentRequestResult)
	})

	suite.Run("Payment request fails when pricing", func() {
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
			On("PriceServiceItem", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(unit.Cents(0), nil, errors.New(errMsg))

		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, nil)
		failingCreator := NewPaymentRequestCreator(planner, failingServiceItemPricer)

		_, err := failingCreator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)
		suite.Error(err)
		wrappedErr := errors.Unwrap(err)
		suite.Equal(errMsg, wrappedErr.Error())
	})

	suite.Run("Given a non-existent move task order id, the create should fail", func() {
		badID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
		estimatedWeight := unit.Pound(2048)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		serviceItem.MoveTaskOrderID = badID
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: badID,
			IsFinal:         false,
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					MTOServiceItemID: serviceItem.ID,
					MTOServiceItem:   serviceItem,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
							Value:       "3254",
						},
						{
							IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
							Value:       "2022-03-16",
						},
					},
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
		}
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), fmt.Sprintf("ID: %s not found", badID))
	})

	suite.Run("Given an already paid or requested payment service item, the create should not fail", func() {
		move := factory.BuildMove(suite.DB(), []factory.Customization{}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)
		paymentRequest1 := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPaid,
				},
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusPaid,
				},
			},
			{
				Model:    paymentRequest1,
				LinkOnly: true,
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		var paymentRequests models.PaymentRequests
		paymentRequests = append(paymentRequests, paymentRequest1)
		shipment.MoveTaskOrder.PaymentRequests = paymentRequests

		paymentRequest2 := models.PaymentRequest{
			MoveTaskOrderID: move.ID,
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					MTOServiceItemID: serviceItem.ID,
					MTOServiceItem:   serviceItem,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
							Value:       "3254",
						},
						{
							IncomingKey: models.ServiceItemParamNameRequestedPickupDate.String(),
							Value:       "2022-03-16",
						},
					},
					Status: models.PaymentServiceItemStatusRequested,
				},
			},
		}
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest2)

		suite.NoError(err)
	})

	suite.Run("Given no move task order id, the create should fail", func() {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Nil,
			IsFinal:         false,
		}
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)

		suite.Error(err)
		suite.IsType(apperror.InvalidCreateInputError{}, err)
		suite.Equal("Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create", err.Error())
	})

	type generateInvalidMove func() models.Move
	type generateExpectedErrorMessage func(uuid.UUID, uuid.UUID) string

	invalidOrdersTestData := []struct {
		TestDescription      string
		InvalidMove          generateInvalidMove
		ExpectedError        error
		ExpectedErrorMessage generateExpectedErrorMessage
	}{
		// Orders with nil TAC
		{
			TestDescription: "Given move with orders but no LOA, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				orders := mtoInvalid.Orders
				orders.TAC = nil
				suite.MustSave(&orders)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(ordersID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state Orders on MoveTaskOrder (ID: %s) missing Lines of Accounting TAC", ordersID, mtoID)
			},
		},
		// Orders with blank TAC
		{
			TestDescription: "Given move with orders but blank LOA, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				orders := mtoInvalid.Orders
				blankTAC := ""
				orders.TAC = &blankTAC
				err := suite.DB().Update(&orders)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(ordersID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state Orders on MoveTaskOrder (ID: %s) missing Lines of Accounting TAC", ordersID, mtoID)
			},
		},
		// Orders with no OriginDutyLocation
		{
			TestDescription: "Given move with orders no OriginDutyLocation, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				orders := mtoInvalid.Orders
				orders.OriginDutyLocation = nil
				orders.OriginDutyLocationID = nil
				err := suite.DB().Update(&orders)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(ordersID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state Orders on MoveTaskOrder (ID: %s) missing OriginDutyLocation", ordersID, mtoID)
			},
		},
	}

	for _, testData := range invalidOrdersTestData {
		suite.Run(testData.TestDescription, func() {
			mtoInvalid := testData.InvalidMove()
			mtoServiceItem1.MoveTaskOrderID = mtoInvalid.ID
			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: mtoInvalid.ID,
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

			_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)

			suite.Error(err)
			suite.IsType(testData.ExpectedError, err)
			suite.Equal(testData.ExpectedErrorMessage(mtoInvalid.OrdersID, mtoInvalid.ID), err.Error())
		})
	}

	invalidServiceMemberTestData := []struct {
		TestDescription      string
		InvalidMove          generateInvalidMove
		ExpectedError        error
		ExpectedErrorMessage generateExpectedErrorMessage
	}{
		// ServiceMember with no First Name
		{
			TestDescription: "Given move with service member that has no First Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				sm := mtoInvalid.Orders.ServiceMember
				sm.FirstName = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing First Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank First Name
		{
			TestDescription: "Given move with service member that has blank First Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				sm := mtoInvalid.Orders.ServiceMember
				blankStr := ""
				sm.FirstName = &blankStr
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing First Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with no Last Name
		{
			TestDescription: "Given move with service member that has no Last Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				sm := mtoInvalid.Orders.ServiceMember
				sm.LastName = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Last Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank Last Name
		{
			TestDescription: "Given move with service member that has blank Last Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				sm := mtoInvalid.Orders.ServiceMember
				blankStr := ""
				sm.LastName = &blankStr
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Last Name", serviceMemberID, mtoID)
			},
		},
		// Order with no Grade
		{
			TestDescription: "Given move with order that has no Rank, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				mtoInvalid.Orders.Grade = nil
				sm := mtoInvalid.Orders.ServiceMember
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(_ uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state unable to pick contract because move is not available to prime", mtoID)
			},
		},
		// Order with empty Grade
		{
			TestDescription: "Given move with order that has blank Rank, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				mtoInvalid.Orders.Grade = internalmessages.NewOrderPayGrade("")
				sm := mtoInvalid.Orders.ServiceMember
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(_ uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state unable to pick contract because move is not available to prime", mtoID)
			},
		},
		// ServiceMember with no Affiliation
		{
			TestDescription: "Given move with service member that has no Affiliation, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				sm := mtoInvalid.Orders.ServiceMember
				sm.Affiliation = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Affiliation", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank Affiliation
		{
			TestDescription: "Given move with service member that has blank Affiliation, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := factory.BuildMove(suite.DB(), nil, nil)
				sm := mtoInvalid.Orders.ServiceMember
				blank := models.ServiceMemberAffiliation("")
				sm.Affiliation = &blank
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: apperror.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("ID: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Affiliation", serviceMemberID, mtoID)
			},
		},
	}

	for _, testData := range invalidServiceMemberTestData {
		suite.Run(testData.TestDescription, func() {
			mtoInvalid := testData.InvalidMove()
			mtoServiceItem1.MoveTaskOrderID = mtoInvalid.ID
			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: mtoInvalid.ID,
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
			_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)

			suite.Error(err)
			suite.IsType(testData.ExpectedError, err)
			suite.Equal(testData.ExpectedErrorMessage(mtoInvalid.Orders.ServiceMemberID, mtoInvalid.ID), err.Error())
		})
	}

	suite.Run("Given a non-existent service item id, the create should fail", func() {
		badID := uuid.Must(uuid.NewV4())
		mtoServiceItem1.MoveTaskOrderID = moveTaskOrder.ID
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: badID,
					MTOServiceItem:   mtoServiceItem1,
				},
			},
		}
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found for MTO Service Item")
	})

	suite.Run("Given a submitted (not approved) service item, the create should fail", func() {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItemSubmitted.ID,
				},
			},
		}
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})
	suite.Run("Given a submitted (not approved) service item, the create should fail", func() {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItemRejected.ID,
				},
			},
		}
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})
	suite.Run("Given a non-existent service item param key id, the create should fail", func() {
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
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found Service Item Param Key ID", badID), err.Error())
	})

	suite.Run("Given a non-existent service item param key name, the create should fail", func() {
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
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &invalidPaymentRequest)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal("Not found Service Item Param Key bogus: FETCH_NOT_FOUND", err.Error())
	})

	suite.Run("Payment request numbers increment by 1", func() {
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
		_, err = creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest1)
		suite.FatalNoError(err)

		paymentRequest2 := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							ServiceItemParamKeyID: serviceItemParamKey1.ID,
							Value:                 "3254",
						},
					},
				},
			},
		}
		_, err = creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest2)
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

	suite.Run("Payment request number fails due to empty MTO ReferenceID", func() {

		saveReferenceID := *moveTaskOrder.ReferenceID
		*moveTaskOrder.ReferenceID = ""
		suite.MustSave(&moveTaskOrder)

		// Create new one
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
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest1)
		suite.Contains(err.Error(), "has missing ReferenceID")

		moveTaskOrder.ReferenceID = &saveReferenceID
		suite.MustSave(&moveTaskOrder)
	})

	suite.Run("cannot submit a payment request if any final payment request already exists", func() {
		paymentRequest1 := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         true,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: models.ServiceItemParamNameWeightEstimated.String(),
							Value:       "1000",
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

		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest1)
		suite.FatalNoError(err)
		paymentRequest2 := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem3.ID,
					MTOServiceItem:   mtoServiceItem3,
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
			},
		}

		_, err = creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest2)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "final PaymentRequest has already been submitted")

		// We need to reset this to prevent prevent tests below this one from breaking
		paymentRequest1.IsFinal = false
		suite.MustSave(&paymentRequest1)
	})

	suite.Run("payment request can be created after a final request that was rejected", func() {
		paymentRequest1 := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         true,
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
			},
		}

		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest1)
		suite.FatalNoError(err)

		paymentRequest1.Status = models.PaymentRequestStatusReviewedAllRejected
		suite.MustSave(&paymentRequest1)

		paymentRequest2 := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         true,
			PaymentServiceItems: []models.PaymentServiceItem{
				{
					MTOServiceItemID: mtoServiceItem2.ID,
					MTOServiceItem:   mtoServiceItem2,
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
			},
		}

		_, err = creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest2)
		suite.NoError(err)

		paymentRequest1.IsFinal = false
		suite.MustSave(&paymentRequest1)
		paymentRequest2.IsFinal = false
		suite.MustSave(&paymentRequest2)
	})
	suite.Run("Payment request number fails due to nil MTO ReferenceID", func() {

		saveReferenceID := *moveTaskOrder.ReferenceID
		moveTaskOrder.ReferenceID = nil
		suite.MustSave(&moveTaskOrder)

		// Create new one
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
		_, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest1)
		suite.Contains(err.Error(), "has missing ReferenceID")

		moveTaskOrder.ReferenceID = &saveReferenceID
		suite.MustSave(&moveTaskOrder)
	})

	suite.Run("CreatePaymentRequest should not return params from rate engine", func() {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					MTOServiceItemID: mtoServiceItem1.ID,
					MTOServiceItem:   mtoServiceItem1,
				},
			},
		}

		paymentRequestReturn, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequest)
		suite.FatalNoError(err)
		suite.NotEqual(paymentRequestReturn.ID, uuid.Nil)
		suite.Equal(1, len(paymentRequestReturn.PaymentServiceItems), "PaymentServiceItems expect 1")

		// Verify that none of the returned service item params are from the Pricer
		if suite.Len(paymentRequestReturn.PaymentServiceItems, 1) {
			for _, param := range paymentRequestReturn.PaymentServiceItems[0].PaymentServiceItemParams {
				suite.NotEqual(param.ServiceItemParamKey.Origin, string(models.ServiceItemParamOriginPricer))
			}
		}
	})
}

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequestCheckOnNTSRelease() {
	testStorageFacilityZip := "30907"
	testDestinationZip := "78234"
	testEscalationCompounded := 1.04071
	testDLHRate := unit.Millicents(6000)
	testOriginalWeight := unit.Pound(3652)
	testZip3Distance := 1234

	// ((testOriginalWeight / 100.0) * testZip3Distance * testDLHRate * testEscalationCompounded) / 1000
	testDLHTotalPrice := unit.Cents(279407)

	//
	// Test data setup
	//

	// Make storage facility and destination addresses
	storageFacilityAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "235 Prospect Valley Road SE",
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     testStorageFacilityZip,
			},
		},
	}, nil)
	destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "17 8th St",
				City:           "San Antonio",
				State:          "TX",
				PostalCode:     testDestinationZip,
			},
		},
	}, nil)

	// Make a storage facility
	storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
		{
			Model:    storageFacilityAddress,
			LinkOnly: true,
		},
	}, nil)
	// Contract year, service area, rate area, zip3
	contractYear, serviceArea, _, _ := testdatagen.SetupServiceAreaRateArea(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EscalationCompounded: testEscalationCompounded,
			EndDate:              time.Now().Add(time.Hour * 24),
		},
		ReRateArea: models.ReRateArea{
			Name: "Georgia",
		},
		ReZip3: models.ReZip3{
			Zip3:          storageFacilityAddress.PostalCode[0:3],
			BasePointCity: storageFacilityAddress.City,
			State:         storageFacilityAddress.State,
		},
	})

	// DLH price data
	testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          false,
			PriceMillicents:       testDLHRate,
		},
	})

	// Make move and shipment
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	actualPickupDate := time.Date(testdatagen.GHCTestYear, time.January, 15, 0, 0, 0, 0, time.UTC)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ShipmentType:      models.MTOShipmentTypeHHGOutOfNTSDom,
				PrimeActualWeight: &testOriginalWeight,
				ActualPickupDate:  &actualPickupDate,
			},
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
	}, nil)

	mtoServiceItemDLH := factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDLH, move, shipment, nil, nil)

	// Build up a payment request for the DLH.
	paymentRequestArg := models.PaymentRequest{
		MoveTaskOrderID: move.ID,
		IsFinal:         false,
		PaymentServiceItems: models.PaymentServiceItems{
			{
				MTOServiceItemID: mtoServiceItemDLH.ID,
				MTOServiceItem:   mtoServiceItemDLH,
			},
		},
	}

	//
	// Create the payment request
	//

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		testStorageFacilityZip,
		testDestinationZip,
	).Return(testZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	paymentRequest, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequestArg)
	suite.FatalNoError(err)

	// Make sure we have just the DLH payment service item
	if suite.Len(paymentRequest.PaymentServiceItems, 1) {
		psi := paymentRequest.PaymentServiceItems[0]
		suite.Equal(models.ReServiceCodeDLH, psi.MTOServiceItem.ReService.Code)

		// Validate the calculated price
		suite.Equal(testDLHTotalPrice, *psi.PriceCents)

		// Check some key payment service item parameters that are different for NTS-Release
		referenceDateParam := getPaymentServiceItemParam(psi.PaymentServiceItemParams, models.ServiceItemParamNameReferenceDate)
		actualPickupDateStr := actualPickupDate.Format(ghcrateengine.DateParamFormat)
		if suite.NotNil(referenceDateParam) {
			suite.Equal(actualPickupDateStr, referenceDateParam.Value)
		}
		actualPickupDateParam := getPaymentServiceItemParam(psi.PaymentServiceItemParams, models.ServiceItemParamNameActualPickupDate)
		if suite.NotNil(actualPickupDateParam) {
			suite.Equal(actualPickupDateStr, actualPickupDateParam.Value)
		}

		requestedPickupDateParam := getPaymentServiceItemParam(psi.PaymentServiceItemParams, models.ServiceItemParamNameRequestedPickupDate)
		suite.Nil(requestedPickupDateParam)
	}
}

func getPaymentServiceItemParam(psiParams models.PaymentServiceItemParams, key models.ServiceItemParamName) *models.PaymentServiceItemParam {
	for _, param := range psiParams {
		if param.ServiceItemParamKey.Key == key {
			return &param
		}
	}

	return nil
}
