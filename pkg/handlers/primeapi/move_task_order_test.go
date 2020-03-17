package primeapi

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

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

func (suite *HandlerSuite) TestListMoveTaskOrdersHandler() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			IsAvailableToPrime: true,
		},
	})

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
		},
	})

	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrder.ID,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: moveTaskOrder,
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: moveTaskOrder,
	})

	// unavailable MTO
	testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := FetchMTOUpdatesHandler{HandlerContext: context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
	suite.Equal(1, len(moveTaskOrdersPayload[0].PaymentRequests))
	suite.Equal(1, len(moveTaskOrdersPayload[0].MtoServiceItems()))
	suite.Equal(2, len(moveTaskOrdersPayload[0].MtoShipments))
	suite.NotNil(moveTaskOrdersPayload[0].MtoShipments[0].ETag)
}

func (suite *HandlerSuite) TestListMoveTaskOrdersHandlerReturnsUpdated() {
	now := time.Now()
	lastFetch := now.Add(-time.Second)

	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			IsAvailableToPrime: true,
		},
	})

	// this MTO should not be returned
	olderMoveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			IsAvailableToPrime: true,
		},
	})

	// Pop will overwrite UpdatedAt when saving a model, so use SQL to set it in the past
	suite.NoError(suite.DB().RawQuery("UPDATE move_task_orders SET updated_at=? WHERE id=?",
		now.Add(-2*time.Second), olderMoveTaskOrder.ID).Exec())

	since := lastFetch.Unix()
	request := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders?since=%d", lastFetch.Unix()), nil)

	params := movetaskorderops.FetchMTOUpdatesParams{HTTPRequest: request, Since: &since}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := FetchMTOUpdatesHandler{HandlerContext: context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.FetchMTOUpdatesOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
}

func (suite *HandlerSuite) TestUpdateMTOPostCounselingInfo() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
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
		},
		IfMatch: eTag,
	}

	suite.T().Run("Successful patch - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		updater := movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder)
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
