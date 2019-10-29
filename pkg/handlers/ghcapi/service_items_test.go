package ghcapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	serviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	serviceitem "github.com/transcom/mymove/pkg/services/service_item"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListServiceItemsHandler() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	uuidString := mto.ID.String()
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		ServiceItem: models.ServiceItem{
			MoveTaskOrderID: id,
		},
	}
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), assertions)

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/service_items", uuidString), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := serviceitemop.ListServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: uuidString,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := ListServiceItemsHandler{
			HandlerContext:         handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:         query.NewQueryFilter,
			ServiceItemListFetcher: serviceitem.NewServiceItemListFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&serviceitemop.ListServiceItemsOK{}, response)
		okResponse := response.(*serviceitemop.ListServiceItemsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(serviceItem.ID.String(), okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		serviceItem := models.ServiceItem{ID: id}
		params := serviceitemop.ListServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: uuidString,
		}
		serviceItemListFetcher := &mocks.ServiceItemListFetcher{}
		serviceItemListFetcher.On("FetchServiceItemList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.ServiceItems{serviceItem}, nil).Once()
		handler := ListServiceItemsHandler{
			HandlerContext:         handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:         newQueryFilter,
			ServiceItemListFetcher: serviceItemListFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&serviceitemop.ListServiceItemsOK{}, response)
		okResponse := response.(*serviceitemop.ListServiceItemsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := serviceitemop.ListServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: uuidString,
		}
		expectedError := models.ErrFetchNotFound
		serviceItemListFetcher := &mocks.ServiceItemListFetcher{}
		serviceItemListFetcher.On("FetchServiceItemList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := ListServiceItemsHandler{
			HandlerContext:         handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:         newQueryFilter,
			ServiceItemListFetcher: serviceItemListFetcher,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestCreateServiceItemHandler() {
	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	serviceItem := models.ServiceItem{ID: serviceItemID, MoveTaskOrderID: moveTaskOrderID}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/%s/service_items", moveTaskOrderID.String()), nil)
	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := serviceitemop.CreateServiceItemParams{
		HTTPRequest:     req,
		MoveTaskOrderID: serviceItem.MoveTaskOrderID.String(),
	}

	suite.T().Run("Successful create", func(t *testing.T) {
		serviceItemCreator := &mocks.ServiceItemCreator{}

		serviceItemCreator.On("CreateServiceItem",
			&serviceItem,
			mock.Anything).Return(&serviceItem, nil, nil).Once()

		handler := CreateServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&serviceitemop.CreateServiceItemCreated{}, response)
	})

	suite.T().Run("Failed create", func(t *testing.T) {
		serviceItemCreator := &mocks.ServiceItemCreator{}

		serviceItemCreator.On("CreateServiceItem",
			&serviceItem,
			mock.Anything).Return(&serviceItem, nil, nil).Once()

		handler := CreateServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&serviceitemop.CreateServiceItemCreated{}, response)
	})
}
