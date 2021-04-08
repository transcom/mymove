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

	serviceItemParamKey4 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameEscalationCompounded,
			Description: "escalation factor",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})

	serviceItemParamKey5 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameContractYearName,
			Description: "contract year name",
			Type:        models.ServiceItemParamTypeString,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})

	serviceItemParamKey6 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameIsPeak,
			Description: "is peak",
			Type:        models.ServiceItemParamTypeBoolean,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})

	serviceItemParamKey7 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNamePriceRateOrFactor,
			Description: "Price, rate, or factor used in calculation",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginPricer,
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
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey4.ID,
			ServiceItemParamKey:   serviceItemParamKey4,
		},
	})
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey5.ID,
			ServiceItemParamKey:   serviceItemParamKey5,
		},
	})
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey6.ID,
			ServiceItemParamKey:   serviceItemParamKey6,
		},
	})
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey7.ID,
			ServiceItemParamKey:   serviceItemParamKey7,
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
	displayParams := models.PaymentServiceItemParams{
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
		On("PriceServiceItem", mock.Anything).Return(testPrice, displayParams, nil).
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
		suite.Equal(7, len(paymentRequest.PaymentServiceItems[0].PaymentServiceItemParams), "PaymentServiceItems[1].PaymentServiceItemParams expect 6")
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
			On("PriceServiceItem", mock.Anything).Return(unit.Cents(0), nil, errors.New(errMsg)).
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
		suite.IsType(services.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("id: %s not found for Move", badID), err.Error())
	})

	suite.T().Run("Given no move task order id, the create should fail", func(t *testing.T) {
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: uuid.Nil,
			IsFinal:         false,
		}
		_, err := creator.CreatePaymentRequest(&invalidPaymentRequest)

		suite.Error(err)
		suite.IsType(services.InvalidCreateInputError{}, err)
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
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				orders := mtoInvalid.Orders
				orders.TAC = nil
				suite.MustSave(&orders)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(ordersID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state Orders on MoveTaskOrder (ID: %s) missing Lines of Accounting TAC", ordersID, mtoID)
			},
		},
		// Orders with blank TAC
		{
			TestDescription: "Given move with orders but blank LOA, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				orders := mtoInvalid.Orders
				blankTAC := ""
				orders.TAC = &blankTAC
				err := suite.DB().Update(&orders)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(ordersID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state Orders on MoveTaskOrder (ID: %s) missing Lines of Accounting TAC", ordersID, mtoID)
			},
		},
		// Orders with no OriginDutyStation
		{
			TestDescription: "Given move with orders no OriginDutyStation, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				orders := mtoInvalid.Orders
				orders.OriginDutyStation = nil
				orders.OriginDutyStationID = nil
				err := suite.DB().Update(&orders)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(ordersID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state Orders on MoveTaskOrder (ID: %s) missing OriginDutyStation", ordersID, mtoID)
			},
		},
	}

	for _, testData := range invalidOrdersTestData {
		suite.T().Run(testData.TestDescription, func(t *testing.T) {
			mtoInvalid := testData.InvalidMove()
			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: mtoInvalid.ID,
			}
			_, err := creator.CreatePaymentRequest(&paymentRequest)

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
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				sm.FirstName = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing First Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank First Name
		{
			TestDescription: "Given move with service member that has blank First Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				blankStr := ""
				sm.FirstName = &blankStr
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing First Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with no Last Name
		{
			TestDescription: "Given move with service member that has no Last Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				sm.LastName = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Last Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank Last Name
		{
			TestDescription: "Given move with service member that has blank Last Name, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				blankStr := ""
				sm.LastName = &blankStr
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Last Name", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with no Rank
		{
			TestDescription: "Given move with service member that has no Rank, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				sm.Rank = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Rank", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank Rank
		{
			TestDescription: "Given move with service member that has blank Rank, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				var blank models.ServiceMemberRank
				blank = ""
				sm.Rank = &blank
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Rank", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with no Affiliation
		{
			TestDescription: "Given move with service member that has no Affiliation, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				sm.Affiliation = nil
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Affiliation", serviceMemberID, mtoID)
			},
		},
		// ServiceMember with blank Affiliation
		{
			TestDescription: "Given move with service member that has blank Affiliation, the create should fail",
			InvalidMove: func() models.Move {
				mtoInvalid := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
				sm := mtoInvalid.Orders.ServiceMember
				var blank models.ServiceMemberAffiliation
				blank = ""
				sm.Affiliation = &blank
				err := suite.DB().Update(&sm)
				suite.FatalNoError(err)
				return mtoInvalid
			},
			ExpectedError: services.ConflictError{},
			ExpectedErrorMessage: func(serviceMemberID uuid.UUID, mtoID uuid.UUID) string {
				return fmt.Sprintf("id: %s is in a conflicting state ServiceMember on MoveTaskOrder (ID: %s) missing Affiliation", serviceMemberID, mtoID)
			},
		},
	}

	for _, testData := range invalidServiceMemberTestData {
		suite.T().Run(testData.TestDescription, func(t *testing.T) {
			mtoInvalid := testData.InvalidMove()
			paymentRequest := models.PaymentRequest{
				MoveTaskOrderID: mtoInvalid.ID,
			}
			_, err := creator.CreatePaymentRequest(&paymentRequest)

			suite.Error(err)
			suite.IsType(testData.ExpectedError, err)
			suite.Equal(testData.ExpectedErrorMessage(mtoInvalid.Orders.ServiceMemberID, mtoInvalid.ID), err.Error())
		})
	}

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
		suite.IsType(services.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("id: %s not found for MTO Service Item", badID), err.Error())
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
		suite.IsType(services.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("id: %s not found Service Item Param Key ID", badID), err.Error())
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
		suite.IsType(services.NotFoundError{}, err)
		suite.Equal("not found Service Item Param Key bogus: FETCH_NOT_FOUND", err.Error())
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

	suite.T().Run("Payment request number fails due to empty MTO ReferenceID", func(t *testing.T) {

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
		_, err := creator.CreatePaymentRequest(&paymentRequest1)
		suite.Contains(err.Error(), "has missing ReferenceID")

		moveTaskOrder.ReferenceID = &saveReferenceID
		suite.MustSave(&moveTaskOrder)
	})
	suite.T().Run("payment request created after final request fails", func(t *testing.T) {
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

		_, err := creator.CreatePaymentRequest(&paymentRequest1)
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

		_, err = creator.CreatePaymentRequest(&paymentRequest2)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "final PaymentRequest has already been submitted")

		// We need to reset this to prevent prevent tests below this one from breaking
		paymentRequest1.IsFinal = false
		suite.MustSave(&paymentRequest1)
	})

	suite.T().Run("payment request can be created after a final request that was rejected", func(t *testing.T) {
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

		_, err := creator.CreatePaymentRequest(&paymentRequest1)
		suite.FatalNoError(err)

		paymentRequest1.Status = models.PaymentRequestStatusReviewedAllRejected
		suite.MustSave(&paymentRequest1)

		paymentRequest2 := models.PaymentRequest{
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

		_, err = creator.CreatePaymentRequest(&paymentRequest2)
		suite.NoError(err)

		paymentRequest1.IsFinal = false
		suite.MustSave(&paymentRequest1)
		paymentRequest2.IsFinal = false
		suite.MustSave(&paymentRequest2)
	})
	suite.T().Run("Payment request number fails due to nil MTO ReferenceID", func(t *testing.T) {

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
		_, err := creator.CreatePaymentRequest(&paymentRequest1)
		suite.Contains(err.Error(), "has missing ReferenceID")

		moveTaskOrder.ReferenceID = &saveReferenceID
		suite.MustSave(&moveTaskOrder)
	})

	suite.T().Run("CreatePaymentRequest should not return params from rate engine", func(t *testing.T) {
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

		paymentRequestReturn, err := creator.CreatePaymentRequest(&paymentRequest)
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
