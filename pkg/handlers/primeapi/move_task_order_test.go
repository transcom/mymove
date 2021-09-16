package primeapi

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/go-openapi/strfmt"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/services/fetch"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListMovesHandlerReturnsUpdated() {
	now := time.Now()
	lastFetch := now.Add(-time.Second)

	move := testdatagen.MakeAvailableMove(suite.DB())

	// this move should not be returned
	olderMove := testdatagen.MakeAvailableMove(suite.DB())

	// Pop will overwrite UpdatedAt when saving a model, so use SQL to set it in the past
	suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
		now.Add(-2*time.Second), olderMove.ID).Exec())
	suite.Require().NoError(suite.DB().RawQuery("UPDATE orders SET updated_at=$1 WHERE id=$2;",
		now.Add(-10*time.Second), olderMove.OrdersID).Exec())

	since := handlers.FmtDateTime(lastFetch)
	request := httptest.NewRequest("GET", fmt.Sprintf("/moves?since=%s", since.String()), nil)
	params := movetaskorderops.ListMovesParams{HTTPRequest: request, Since: since}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := ListMovesHandler{HandlerContext: context, MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher()}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	listMovesResponse := response.(*movetaskorderops.ListMovesOK)
	movesList := listMovesResponse.Payload

	suite.Equal(1, len(movesList))
	suite.Equal(move.ID.String(), movesList[0].ID.String())
}

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMoveTaskOrderHandler{context,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	suite.T().Run("Success with Prime-available move by ID", func(t *testing.T) {
		successMove := testdatagen.MakeAvailableMove(suite.DB())
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.ID.String(),
		}
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload
		suite.Equal(movePayload.ID.String(), successMove.ID.String())
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt) // checks that the date is not 0001-01-01
	})

	suite.T().Run("Success with Prime-available move by Locator", func(t *testing.T) {
		successMove := testdatagen.MakeAvailableMove(suite.DB())
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload
		suite.Equal(movePayload.ID.String(), successMove.ID.String())
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt) // checks that the date is not 0001-01-01
	})

	suite.T().Run("Success returns reweighs on shipments if they exist", func(t *testing.T) {
		successMove := testdatagen.MakeAvailableMove(suite.DB())
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			Move: successMove,
		})

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload
		reweighPayload := movePayload.MtoShipments[0].Reweigh
		suite.Equal(movePayload.ID.String(), successMove.ID.String())
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), reweighPayload.ID)
	})

	suite.T().Run("Failure 'Not Found' for non-available move", func(t *testing.T) {
		failureMove := testdatagen.MakeDefaultMove(suite.DB()) // default is not available to Prime
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      failureMove.ID.String(),
		}
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderNotFound{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderNotFound)
		movePayload := moveResponse.Payload
		suite.Contains(*movePayload.Detail, failureMove.ID.String())
	})
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
		queryBuilder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(queryBuilder)
		moveRouter := moverouter.NewMoveRouter()
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		updater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter)
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

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

		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		fetcher := fetch.NewFetcher(queryBuilder)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		updater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter)
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
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
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
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
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
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NewInvalidInputError(mto.ID, nil, validate.NewErrors(), ""))

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationUnprocessableEntity{}, response)
	})
}
