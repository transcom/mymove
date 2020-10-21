package primeapi

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/validate"
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

type FeatureFlag struct {
	Name   string
	Active bool
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
		var serviceItemDDFSIT *primemessages.MTOServiceItemDDFSIT
		for _, item := range moveTaskOrdersPayload[0].MtoServiceItems() {
			if item.ModelType() == primemessages.MTOServiceItemModelTypeMTOServiceItemDDFSIT {
				serviceItemDDFSIT = item.(*primemessages.MTOServiceItemDDFSIT)
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

func (suite *HandlerSuite) TestUpdateMTOPostCounselingInfo() {
	mto := testdatagen.MakeDefaultMove(suite.DB())

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
		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			fetcher,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationOK{}, response)

		okResponse := response.(*movetaskorderops.UpdateMTOPostCounselingInformationOK)
		suite.Equal(mto.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)
		suite.Equal(okResponse.Payload.PpmType, "FULL")
		suite.Equal(okResponse.Payload.PpmEstimatedWeight, int64(3000))
	})
	suite.T().Run("Patch failure - 500", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
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
		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
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
		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
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
