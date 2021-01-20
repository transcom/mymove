package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	customeruserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/customer_users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexCustomerUsersHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	// uuidString := "f0ddc118-3f7e-476b-b8be-0f964a5feee2"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		User: models.User{
			ID: id,
		},
	}
	testdatagen.MakeUser(suite.DB(), assertions)
	testdatagen.MakeDefaultUser(suite.DB())

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/customer_users", nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := customeruserop.IndexCustomerUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexCustomerUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&customeruserop.IndexCustomerUsersOK{}, response)
		okResponse := response.(*customeruserop.IndexCustomerUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := customeruserop.IndexCustomerUsersParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		customerUserListFetcher := &mocks.ListFetcher{}
		customerUserListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		customerUserListFetcher.On("FetchRecordCount",
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexCustomerUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    customerUserListFetcher,
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
