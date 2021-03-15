package primeapi

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/swag"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/services/fetch"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestTruncateAll() {

	move := testdatagen.MakeDefaultMove(suite.DB())
	fmt.Println("created move", move.ID, move.ContractorID)

	err := suite.DB().TruncateAll()
	fmt.Println(err)
	fmt.Println("truncated db")

	foundMove := models.Move{}
	err = suite.DB().Find(&foundMove, move.ID.String())
	fmt.Println(err)
	fmt.Println("found move", foundMove.ID, foundMove.ContractorID)
}

func (suite *HandlerSuite) TestFetchMTOUpdatesHandler() {
	// unavailable MTO
	testdatagen.MakeDefaultMove(suite.DB())

	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
	})

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := FetchMTOUpdatesHandler{
		HandlerContext:       context,
		MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}

	suite.T().Run("with mto service item dimensions", func(t *testing.T) {
		reServiceDomCrating := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DCRT",
				Name: "Dom. Crating",
			},
		})

		mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: moveTaskOrder.ID,
			},
			ReService: reServiceDomCrating,
		})

		testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				Type:      models.DimensionTypeItem,
				Length:    1000,
				Height:    1000,
				Width:     1000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			MTOServiceItem: mtoServiceItem1,
		})

		testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				MTOServiceItemID: mtoServiceItem1.ID,
				Type:             models.DimensionTypeCrate,
				Length:           2000,
				Height:           2000,
				Width:            2000,
				CreatedAt:        time.Time{},
				UpdatedAt:        time.Time{},
			},
		})

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

		suite.Equal(1, len(moveTaskOrdersPayload[0].MtoServiceItems()))
		suite.NotEmpty(moveTaskOrdersPayload[0].MtoServiceItems()[0].(*primemessages.MTOServiceItemDomesticCrating).Crate.ID)
		suite.NotEmpty(moveTaskOrdersPayload[0].MtoServiceItems()[0].(*primemessages.MTOServiceItemDomesticCrating).Item.ID)
	})

	suite.T().Run("with mto service item customer contacts", func(t *testing.T) {
		reServiceDesSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DDFSIT",
				Name: "Destination 1st Day SIT",
			},
		})

		mtoServiceItem2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: moveTaskOrder.ID,
			},
			ReService: reServiceDesSIT,
		})

		testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				MTOServiceItemID:           mtoServiceItem2.ID,
				Type:                       models.CustomerContactTypeFirst,
				TimeMilitary:               "0400Z",
				FirstAvailableDeliveryDate: time.Now(),
			},
			ReService: reServiceDesSIT,
		})

		testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				MTOServiceItemID:           mtoServiceItem2.ID,
				Type:                       models.CustomerContactTypeSecond,
				TimeMilitary:               "0400Z",
				FirstAvailableDeliveryDate: time.Now(),
			},
			ReService: reServiceDesSIT,
		})

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
		moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

		suite.Equal(2, len(moveTaskOrdersPayload[0].MtoServiceItems()))

		// get ddfsit service item
		var serviceItemDDFSIT *primemessages.MTOServiceItemDestSIT
		for _, item := range moveTaskOrdersPayload[0].MtoServiceItems() {
			if item.ModelType() == primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT {
				serviceItemDDFSIT = item.(*primemessages.MTOServiceItemDestSIT)
				break
			}
		}

		suite.NotEmpty(serviceItemDDFSIT.ID)
	})
}

func (suite *HandlerSuite) TestFetchMTOUpdatesHandlerPaymentRequest() {
	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())

	// This should create all the other associated records we need.
	paymentServiceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameRequestedPickupDate,
		},
	})

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := FetchMTOUpdatesHandler{HandlerContext: context, MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(suite.DB())}

	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
	suite.Equal(1, len(moveTaskOrdersPayload[0].PaymentRequests))
	suite.Equal(paymentServiceItemParam.PaymentServiceItem.PaymentRequestID.String(), moveTaskOrdersPayload[0].PaymentRequests[0].ID.String())
	suite.Equal(1, len(moveTaskOrdersPayload[0].PaymentRequests[0].PaymentServiceItems))
	suite.Equal(paymentServiceItemParam.PaymentServiceItemID.String(), moveTaskOrdersPayload[0].PaymentRequests[0].PaymentServiceItems[0].ID.String())
	suite.Equal(1, len(moveTaskOrdersPayload[0].PaymentRequests[0].PaymentServiceItems[0].PaymentServiceItemParams))
	suite.Equal(paymentServiceItemParam.ID.String(), moveTaskOrdersPayload[0].PaymentRequests[0].PaymentServiceItems[0].PaymentServiceItemParams[0].ID.String())
	suite.Equal(1, len(moveTaskOrdersPayload[0].MtoShipments))
	suite.Equal(paymentServiceItemParam.PaymentServiceItem.MTOServiceItem.MTOShipmentID.String(), moveTaskOrdersPayload[0].MtoShipments[0].ID.String())
	suite.NotNil(moveTaskOrdersPayload[0].MtoShipments[0].ETag)
}

func (suite *HandlerSuite) TestFetchMTOUpdatesHandlerMinimal() {
	// Creates a move task order with one minimal shipment and no payment requests
	// or service items
	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())

	testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
	})

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := FetchMTOUpdatesHandler{HandlerContext: context, MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(suite.DB())}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
	suite.NotNil(moveTaskOrdersPayload[0].MtoShipments[0].ETag)
}

func (suite *HandlerSuite) TestListMoveTaskOrdersHandlerReturnsUpdated() {
	now := time.Now()
	lastFetch := now.Add(-time.Second)

	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())

	// this MTO should not be returned
	olderMoveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())

	// Pop will overwrite UpdatedAt when saving a model, so use SQL to set it in the past
	suite.NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
		now.Add(-2*time.Second), olderMoveTaskOrder.ID).Exec())

	since := lastFetch.Unix()
	request := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders?since=%d", lastFetch.Unix()), nil)

	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request, Since: &since}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := FetchMTOUpdatesHandler{HandlerContext: context, MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(suite.DB())}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
}

func (suite *HandlerSuite) makeAvailableMoveWithAddress(addressToSet models.Address) models.Move {
	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: addressToSet,
	})

	newDutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			AddressID: address.ID,
			Address:   address,
		},
	})

	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			NewDutyStationID: newDutyStation.ID,
			NewDutyStation:   newDutyStation,
		},
	})

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: swag.Time(time.Now()),
			Status:             models.MoveStatusAPPROVED,
		},
		Order: order,
	})

	return move
}

func (suite *HandlerSuite) equalAddress(expected models.Address, actual *primemessages.Address) {
	suite.Equal(expected.ID.String(), actual.ID.String())
	suite.Equal(expected.StreetAddress1, *actual.StreetAddress1)
	suite.Equal(*expected.StreetAddress2, *actual.StreetAddress2)
	suite.Equal(*expected.StreetAddress3, *actual.StreetAddress3)
	suite.Equal(expected.City, *actual.City)
	suite.Equal(expected.State, *actual.State)
	suite.Equal(expected.PostalCode, *actual.PostalCode)
	suite.Equal(*expected.Country, *actual.Country)
}

func (suite *HandlerSuite) equalPaymentRequest(expected models.PaymentRequest, actual *primemessages.PaymentRequest) {
	suite.Equal(expected.ID.String(), actual.ID.String())
	suite.Equal(expected.MoveTaskOrderID.String(), actual.MoveTaskOrderID.String())
	suite.Equal(expected.IsFinal, *actual.IsFinal)
	suite.Equal(expected.Status.String(), string(actual.Status))
	suite.Equal(expected.RejectionReason, actual.RejectionReason)
	suite.Equal(expected.PaymentRequestNumber, actual.PaymentRequestNumber)
}

func (suite *HandlerSuite) TestFetchMTOUpdatesHandlerLoopIteratorPointer() {
	// Create two moves with different addresses.
	move1 := suite.makeAvailableMoveWithAddress(models.Address{
		StreetAddress1: "1 First St",
		StreetAddress2: swag.String("Apt 1"),
		StreetAddress3: swag.String("Suite A"),
		City:           "Augusta",
		State:          "GA",
		PostalCode:     "30907",
		Country:        swag.String("US"),
	})

	move2 := suite.makeAvailableMoveWithAddress(models.Address{
		StreetAddress1: "2 Second St",
		StreetAddress2: swag.String("Apt 2"),
		StreetAddress3: swag.String("Suite B"),
		City:           "Columbia",
		State:          "SC",
		PostalCode:     "29212",
		Country:        swag.String("United States"),
	})

	// Create two payment requests on the second move.
	paymentRequest1 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: move2,
		PaymentRequest: models.PaymentRequest{
			IsFinal:        false,
			SequenceNumber: 1,
		},
	})

	paymentRequest2 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: move2,
		PaymentRequest: models.PaymentRequest{
			IsFinal:        true,
			SequenceNumber: 2,
		},
	})

	// Setup and call the handler.
	request := httptest.NewRequest("GET", "/move-task-orders", nil)
	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := FetchMTOUpdatesHandler{
		HandlerContext:       context,
		MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	// Check the addresses across the two moves.
	// NOTE: The payload isn't ordered, so I have to associate the correct move.
	suite.FatalFalse(len(moveTaskOrdersPayload) != 2)
	move1Payload := moveTaskOrdersPayload[0]
	move2Payload := moveTaskOrdersPayload[1]
	if move1Payload.ID.String() != move1.ID.String() {
		move1Payload = moveTaskOrdersPayload[1]
		move2Payload = moveTaskOrdersPayload[0]
	}

	suite.equalAddress(move1.Orders.NewDutyStation.Address, move1Payload.MoveOrder.DestinationDutyStation.Address)
	suite.equalAddress(move2.Orders.NewDutyStation.Address, move2Payload.MoveOrder.DestinationDutyStation.Address)

	// Check the two payment requests across the second move.
	// NOTE: The payload isn't ordered, so I have to associate the correct payment request.
	paymentRequestsPayload := move2Payload.PaymentRequests
	suite.FatalFalse(len(paymentRequestsPayload) != 2)
	paymentRequest1Payload := paymentRequestsPayload[0]
	paymentRequest2Payload := paymentRequestsPayload[1]
	if paymentRequest1Payload.ID.String() != paymentRequest1.ID.String() {
		paymentRequest1Payload = paymentRequestsPayload[1]
		paymentRequest2Payload = paymentRequestsPayload[0]
	}

	suite.equalPaymentRequest(paymentRequest1, paymentRequest1Payload)
	suite.equalPaymentRequest(paymentRequest2, paymentRequest2Payload)
}

func (suite *HandlerSuite) TestUpdateMTOPostCounselingInfo() {
	mto := testdatagen.MakeAvailableMove(suite.DB())

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	eTag := base64.StdEncoding.EncodeToString([]byte(mto.UpdatedAt.Format(time.RFC3339Nano)))

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", mto.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	ppmType := "FULL"
	params := movetaskorderops.UpdateMTOPostCounselingInformationParams{
		HTTPRequest:     req,
		MoveTaskOrderID: mto.ID.String(),
		Body: movetaskorderops.UpdateMTOPostCounselingInformationBody{
			PpmType:            ppmType,
			PpmEstimatedWeight: 3000,
			PointOfContact:     "user@prime.com",
		},
		IfMatch: eTag,
	}

	suite.T().Run("Successful patch - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)
		updater := movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator)
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			fetcher,
			updater,
			mtoChecker,
		}

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationOK{}, response)

		okResponse := response.(*movetaskorderops.UpdateMTOPostCounselingInformationOK)
		suite.Equal(mto.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)
		suite.Equal(okResponse.Payload.PpmType, "FULL")
		suite.Equal(okResponse.Payload.PpmEstimatedWeight, int64(3000))
	})

	suite.T().Run("Unsuccessful patch - Integration Test - patch fail MTO not available", func(t *testing.T) {
		defaultMTO := testdatagen.MakeDefaultMove(suite.DB())

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		eTag := base64.StdEncoding.EncodeToString([]byte(defaultMTO.UpdatedAt.Format(time.RFC3339Nano)))

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", defaultMTO.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		ppmType := "FULL"
		defaultMTOParams := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: defaultMTO.ID.String(),
			Body: movetaskorderops.UpdateMTOPostCounselingInformationBody{
				PpmType:            ppmType,
				PpmEstimatedWeight: 3000,
				PointOfContact:     "user@prime.com",
			},
			IfMatch: eTag,
		}

		mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)
		updater := movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator)
		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			fetcher,
			updater,
			mtoChecker,
		}

		response := handler.Handle(defaultMTOParams)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationNotFound{}, response)
	})

	suite.T().Run("Patch failure - 500", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationInternalServerError{}, response)
	})

	suite.T().Run("Patch failure - 404", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationNotFound{}, response)
	})

	suite.T().Run("Patch failure - 422", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NewInvalidInputError(mto.ID, nil, validate.NewErrors(), ""))

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationUnprocessableEntity{}, response)
	})
}
