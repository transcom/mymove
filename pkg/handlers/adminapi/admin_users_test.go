package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	adminuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/admin_users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexAdminUsersHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		AdminUser: models.AdminUser{
			ID: id,
		},
	}
	testdatagen.MakeAdminUser(suite.DB(), assertions)
	testdatagen.MakeDefaultAdminUser(suite.DB())

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", "/admin_users", nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexAdminUsersHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:       query.NewQueryFilter,
			AdminUserListFetcher: adminuser.NewAdminUserListFetcher(queryBuilder),
			NewPagination:        pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&adminuserop.IndexAdminUsersOK{}, response)
		okResponse := response.(*adminuserop.IndexAdminUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		adminUser := models.AdminUser{ID: id}
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: req,
		}
		adminUserListFetcher := &mocks.AdminUserListFetcher{}
		adminUserListFetcher.On("FetchAdminUserList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.AdminUsers{adminUser}, nil).Once()
		handler := IndexAdminUsersHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:       newQueryFilter,
			AdminUserListFetcher: adminUserListFetcher,
			NewPagination:        pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&adminuserop.IndexAdminUsersOK{}, response)
		okResponse := response.(*adminuserop.IndexAdminUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		adminUserListFetcher := &mocks.AdminUserListFetcher{}
		adminUserListFetcher.On("FetchAdminUserList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexAdminUsersHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:       newQueryFilter,
			AdminUserListFetcher: adminUserListFetcher,
			NewPagination:        pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
