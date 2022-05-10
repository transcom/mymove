package adminapi

import (
	"fmt"
	"net/http"

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
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		adminUsers := models.AdminUsers{
			testdatagen.MakeDefaultAdminUser(suite.DB()),
			testdatagen.MakeDefaultAdminUser(suite.DB()),
		}
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/admin_users"),
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
		suite.Equal(adminUsers[0].ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		params := adminuserop.IndexAdminUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/admin_users"),
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
			NewQueryFilter:       newMockQueryFilterBuilder(&mocks.QueryFilter{}),
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
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		adminUser := testdatagen.MakeDefaultAdminUser(suite.DB())
		params := adminuserop.GetAdminUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/admin_users/%s", adminUser.ID)),
			AdminUserID: strfmt.UUID(adminUser.ID.String()),
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
		suite.Equal(adminUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("successful response", func() {
		adminUser := models.AdminUser{ID: uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763")}
		params := adminuserop.GetAdminUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/admin_users/%s", adminUser.ID)),
			AdminUserID: strfmt.UUID(adminUser.ID.String()),
		}
		adminUserFetcher := &mocks.AdminUserFetcher{}
		adminUserFetcher.On("FetchAdminUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(adminUser, nil).Once()
		handler := GetAdminUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			adminUserFetcher,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		suite.IsType(&adminuserop.GetAdminUserOK{}, response)
		okResponse := response.(*adminuserop.GetAdminUserOK)
		suite.Equal(adminUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		adminUser := testdatagen.MakeDefaultAdminUser(suite.DB())
		params := adminuserop.GetAdminUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/admin_users/%s", adminUser.ID)),
			AdminUserID: strfmt.UUID(adminUser.ID.String()),
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
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
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
	organizationID, _ := uuid.NewV4()
	adminUser := models.AdminUser{
		ID:             uuid.Nil,
		OrganizationID: &organizationID,
		UserID:         nil,
		Role:           models.SystemAdminRole,
		Active:         true,
	}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.Run("Successful create", func() {
		params := adminuserop.CreateAdminUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/admin_users"),
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
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/admin_users"),
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
	adminUserID := uuid.Must(uuid.NewV4())
	adminUser := models.AdminUser{ID: adminUserID, FirstName: "Leo", LastName: "Spaceman"}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.Run("Successful update", func() {
		params := adminuserop.UpdateAdminUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/admin_users/%s", adminUserID)),
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
		// TESTCASE SCENARIO
		// Under test: UpdateAdminUserHandler
		// Mocked: UpdateAdminUser
		// Set up: UpdateAdminUser is mocked to return validation errors as if an error was encountered
		// Expected outcome: The handler should see the validation errors and also return an error
		params := adminuserop.UpdateAdminUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/admin_users/%s", adminUserID)),
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
