package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/services"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListMTOServiceItemHandler() {
	reServiceID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	mtoShipmentID, _ := uuid.NewV4()

	mto := testdatagen.MakeDefaultMove(suite.DB())
	reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   reServiceID,
			Code: "TEST10000",
		},
	})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{ID: mtoShipmentID},
	})
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: serviceItemID, MoveTaskOrderID: mto.ID, ReServiceID: reService.ID, MTOShipmentID: &mtoShipment.ID,
		},
	})
	serviceItems := models.MTOServiceItems{serviceItem}

	req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mto.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoserviceitemop.ListMTOServiceItemsParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(serviceItem.MoveTaskOrderID),
	}

	suite.T().Run("Successful list fetch - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := ListMTOServiceItemsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			listFetcher,
			fetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsOK{}, response)

		okResponse := response.(*mtoserviceitemop.ListMTOServiceItemsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(serviceItems[0].ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.T().Run("Failure list fetch - Internal Server Error", func(t *testing.T) {
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOServiceItemsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		internalServerErr := errors.New("ServerError")

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(internalServerErr)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsInternalServerError{}, response)
	})

	suite.T().Run("Failure list fetch - 404 Not Found - Move Task Order ID", func(t *testing.T) {
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOServiceItemsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		notfound := errors.New("Not found error")

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(notfound)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsNotFound{}, response)
	})
}

func (suite *HandlerSuite) createServiceItem() (models.MTOServiceItem, models.Move) {
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	return serviceItem, move
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandler() {
	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status",
		moveTaskOrderID, serviceItemID), nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      req,
		IfMatch:          etag.GenerateEtag(time.Now()),
		Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "APPROVED"},
		MoveTaskOrderID:  moveTaskOrderID.String(),
		MtoServiceItemID: serviceItemID.String(),
	}

	// With this first set of tests we'll use mocked service object responses so that we can make sure the handler
	// is returning the right HTTP code given a set of circumstances.
	suite.T().Run("404 - not found response", func(t *testing.T) {
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(errors.New("Not found error")).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusNotFound{}, response)
	})

	suite.T().Run("200 - success response", func(t *testing.T) {
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("UpdateMTOServiceItemStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&models.MTOServiceItem{ID: serviceItemID}, nil).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
	})

	suite.T().Run("412 - precondition failed response", func(t *testing.T) {
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("UpdateMTOServiceItemStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NewPreconditionFailedError(serviceItemID, errors.New("oh no"))).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusPreconditionFailed{}, response)
	})

	suite.T().Run("500 - internal server error response", func(t *testing.T) {
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("UpdateMTOServiceItemStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("oh no")).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusInternalServerError{}, response)
	})

	suite.T().Run("422 - unprocessable entity response", func(t *testing.T) {
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		params.MtoServiceItemID = ""
		fetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusUnprocessableEntity{}, response)
	})

	// With this we'll do a happy path integration test to ensure that the use of the service object
	// by the handler is working as expected.
	suite.T().Run("Successful rejected status update - Integration test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		mtoServiceItem, move := suite.createServiceItem()
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status",
			moveTaskOrderID, serviceItemID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		rejectionReason := "No justification given"
		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "REJECTED", RejectionReason: &rejectionReason},
			MoveTaskOrderID:  move.ID.String(),
			MtoServiceItemID: mtoServiceItem.ID.String(),
		}

		fetcher := fetch.NewFetcher(queryBuilder)
		moveRouter := moverouter.NewMoveRouter(suite.DB(), suite.TestLogger())
		mtoServiceItemStatusUpdater := mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter)

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: mtoServiceItemStatusUpdater,
			Fetcher:               fetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		okResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
		suite.Equal(string(models.MTOServiceItemStatusRejected), string(okResponse.Payload.Status))
		suite.NotNil(okResponse.Payload.RejectedAt)
		suite.Equal(rejectionReason, *okResponse.Payload.RejectionReason)
	})

	// With this we'll do a happy path integration test to ensure that the use of the service object
	// by the handler is working as expected.
	suite.T().Run("Successful status update of MTO service item and event trigger", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		moveRouter := moverouter.NewMoveRouter(suite.DB(), suite.TestLogger())
		mtoServiceItem, availableMove := suite.createServiceItem()
		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		availableMoveID := availableMove.ID
		mtoServiceItemID := mtoServiceItem.ID

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status", availableMoveID, mtoServiceItemID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "APPROVED"},
			MoveTaskOrderID:  availableMoveID.String(),
			MtoServiceItemID: mtoServiceItemID.String(),
		}

		fetcher := fetch.NewFetcher(queryBuilder)
		mtoServiceItemStatusUpdater := mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter)

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOServiceItemUpdater: mtoServiceItemStatusUpdater,
			Fetcher:               fetcher,
		}

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		handler.SetTraceID(traceID)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		okResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
		suite.Equal(string(models.MTOServiceItemStatusApproved), string(okResponse.Payload.Status))
		suite.NotNil(okResponse.Payload.ApprovedAt)
		suite.HasWebhookNotification(mtoServiceItemID, traceID)

		impactedMove := models.Move{}
		_ = suite.DB().Find(&impactedMove, okResponse.Payload.MoveTaskOrderID)
		suite.Equal(models.MoveStatusAPPROVED, impactedMove.Status)
	})
}
