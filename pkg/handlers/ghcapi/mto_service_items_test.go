package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/validate"
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

func (suite *HandlerSuite) TestCreateMTOServiceItemHandler() {
	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	reServiceID, _ := uuid.NewV4()
	mtoShipmentID, _ := uuid.NewV4()

	serviceItem := models.MTOServiceItem{
		ID: serviceItemID, MoveTaskOrderID: moveTaskOrderID, ReServiceID: reServiceID, MTOShipmentID: &mtoShipmentID,
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/%s/mto_service_items", moveTaskOrderID.String()), nil)
	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoserviceitemop.CreateMTOServiceItemParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(serviceItem.MoveTaskOrderID),
		CreateMTOServiceItemBody: mtoserviceitemop.CreateMTOServiceItemBody{
			ReServiceID:   handlers.FmtUUID(serviceItem.ReServiceID),
			MtoShipmentID: handlers.FmtUUIDPtr(serviceItem.MTOShipmentID),
		},
	}

	serviceItemCreator := &mocks.MTOServiceItemCreator{}

	suite.T().Run("Successful create", func(t *testing.T) {
		var serviceItems models.MTOServiceItems
		serviceItems = append(serviceItems, serviceItem)

		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(&serviceItems, nil, nil).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemCreated{}, response)
	})

	suite.T().Run("Failed create: InternalServiceError", func(t *testing.T) {
		err := errors.New("cannot create service item")
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, validate.NewErrors(), err).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemInternalServerError{}, response)
	})

	suite.T().Run("Failed create: UnprocessableEntity", func(t *testing.T) {
		verrs := validate.NewErrors()
		verrs.Add("error", "error test")
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, verrs, nil).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.T().Run("Failed create: UnprocessableEntity - UUID parsing error", func(t *testing.T) {
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		newParams := mtoserviceitemop.CreateMTOServiceItemParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItem.MoveTaskOrderID),
			CreateMTOServiceItemBody: mtoserviceitemop.CreateMTOServiceItemBody{
				ReServiceID:   handlers.FmtUUID(serviceItem.ReServiceID),
				MtoShipmentID: handlers.FmtUUIDPtr(serviceItem.MTOShipmentID),
			},
		}
		newParams.MoveTaskOrderID = "blah"

		response := handler.Handle(newParams)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.T().Run("Failed create: UnprocessableEntity - Violates foreign key constraints", func(t *testing.T) {
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, validate.NewErrors(), errors.New(models.ViolatesForeignKeyConstraint)).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemNotFound{}, response)
	})
}

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
	requestUser := testdatagen.MakeDefaultUser(suite.DB())
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
