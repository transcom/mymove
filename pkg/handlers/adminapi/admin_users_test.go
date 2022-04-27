package adminapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/gen/adminmessages"

	"github.com/go-openapi/strfmt"
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
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	setupRequest := func() *http.Request {
		assertions := testdatagen.Assertions{
			AdminUser: models.AdminUser{
				ID: id,
			},
		}
		testdatagen.MakeAdminUser(suite.DB(), assertions)
		testdatagen.MakeDefaultAdminUser(suite.DB())

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("GET", "/admin_users", nil)
		return suite.AuthenticateAdminRequest(req, requestUser)
	}

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		req := setupRequest()
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexAdminUsersHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter:       query.NewQueryFilter,
			AdminUserListFetcher: adminuser.NewAdminUserListFetcher(queryBuilder),
			NewPagination:        pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&adminuserop.IndexAdminUsersOK{}, response)
		okResponse := response.(*adminuserop.IndexAdminUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(id.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("successful response", func() {
		adminUser := models.AdminUser{ID: id}
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: setupRequest(),
		}
		adminUserListFetcher := &mocks.AdminUserListFetcher{}
		adminUserListFetcher.On("FetchAdminUserList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.AdminUsers{adminUser}, nil).Once()
		adminUserListFetcher.On("FetchAdminUserCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(1, nil).Once()
		handler := IndexAdminUsersHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	suite.Run("unsuccesful response when fetch fails", func() {
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: setupRequest(),
		}
		expectedError := models.ErrFetchNotFound
		adminUserListFetcher := &mocks.AdminUserListFetcher{}
		adminUserListFetcher.On("FetchAdminUserList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		adminUserListFetcher.On("FetchAdminUserCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexAdminUsersHandler{
			HandlerContext:       handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

func (suite *HandlerSuite) TestGetAdminUserHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	setupRequest := func() *http.Request {
		assertions := testdatagen.Assertions{
			AdminUser: models.AdminUser{
				ID: id,
			},
		}
		testdatagen.MakeAdminUser(suite.DB(), assertions)

		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("GET", fmt.Sprintf("/admin_users/%s", id), nil)
		return suite.AuthenticateUserRequest(req, requestUser)
	}

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		params := adminuserop.GetAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUserID: strfmt.UUID(uuidString),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminuser.NewAdminUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&adminuserop.GetAdminUserOK{}, response)
		okResponse := response.(*adminuserop.GetAdminUserOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	suite.Run("successful response", func() {
		adminUser := models.AdminUser{ID: id}
		params := adminuserop.GetAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUserID: strfmt.UUID(uuidString),
		}
		adminUserFetcher := &mocks.AdminUserFetcher{}
		adminUserFetcher.On("FetchAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(adminUser, nil).Once()
		handler := GetAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserFetcher,
			newQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&adminuserop.GetAdminUserOK{}, response)
		okResponse := response.(*adminuserop.GetAdminUserOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		params := adminuserop.GetAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUserID: strfmt.UUID(uuidString),
		}
		expectedError := models.ErrFetchNotFound
		adminUserFetcher := &mocks.AdminUserFetcher{}
		adminUserFetcher.On("FetchAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.AdminUser{}, expectedError).Once()
		handler := GetAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserFetcher,
			newQueryFilter,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestCreateAdminUserHandler() {
	setupRequest := func() *http.Request {
		req := httptest.NewRequest("POST", "/admin_users", nil)
		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		return suite.AuthenticateUserRequest(req, requestUser)
	}

	organizationID, _ := uuid.NewV4()
	adminUserID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	adminUser := models.AdminUser{
		ID:             adminUserID,
		OrganizationID: &organizationID,
		UserID:         nil,
		Role:           models.SystemAdminRole,
		Active:         true,
	}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.Run("Successful create", func() {
		params := adminuserop.CreateAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUser: &adminmessages.AdminUserCreatePayload{
				FirstName:      adminUser.FirstName,
				LastName:       adminUser.LastName,
				OrganizationID: strfmt.UUID(adminUser.OrganizationID.String()),
			},
		}

		adminUserCreator := &mocks.AdminUserCreator{}
		adminUserCreator.On("CreateAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			&adminUser,
			mock.Anything).Return(&adminUser, nil, nil).Once()

		handler := CreateAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&adminuserop.CreateAdminUserCreated{}, response)
	})

	suite.Run("Failed create", func() {
		params := adminuserop.CreateAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUser: &adminmessages.AdminUserCreatePayload{
				FirstName:      adminUser.FirstName,
				LastName:       adminUser.LastName,
				OrganizationID: strfmt.UUID(adminUser.OrganizationID.String()),
			},
		}

		adminUserCreator := &mocks.AdminUserCreator{}
		adminUserCreator.On("CreateAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			&adminUser,
			mock.Anything).Return(&adminUser, nil, nil).Once()

		handler := CreateAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&adminuserop.CreateAdminUserCreated{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateAdminUserHandler() {
	adminUserID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	adminUser := models.AdminUser{ID: adminUserID, FirstName: "Leo", LastName: "Spaceman"}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	setupRequest := func() *http.Request {
		req := httptest.NewRequest("PUT", fmt.Sprintf("/admin_users/%s", adminUserID), nil)
		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		return suite.AuthenticateUserRequest(req, requestUser)
	}

	suite.Run("Successful update", func() {
		params := adminuserop.UpdateAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUser: &adminmessages.AdminUserUpdatePayload{
				FirstName: &adminUser.FirstName,
				LastName:  &adminUser.LastName,
			},
		}

		adminUserUpdater := &mocks.AdminUserUpdater{}
		adminUserUpdater.On("UpdateAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.AdminUser,
		).Return(&adminUser, nil, nil).Once()

		handler := UpdateAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserUpdater,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&adminuserop.UpdateAdminUserOK{}, response)
	})

	suite.Run("Failed update", func() {
		params := adminuserop.UpdateAdminUserParams{
			HTTPRequest: setupRequest(),
			AdminUser: &adminmessages.AdminUserUpdatePayload{
				FirstName: &adminUser.FirstName,
				LastName:  &adminUser.LastName,
			},
		}

		adminUserUpdater := &mocks.AdminUserUpdater{}
		err := validate.NewErrors()

		adminUserUpdater.On("UpdateAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.AdminUser,
		).Return(nil, err, nil).Once()

		handler := UpdateAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserUpdater,
			newQueryFilter,
		}

		handler.Handle(params)
		suite.Error(err, "Error saving user")
	})
}
