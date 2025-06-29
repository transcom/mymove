// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package adminapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	userservice "github.com/transcom/mymove/pkg/services/user"
)

func (suite *HandlerSuite) TestGetUserHandler() {
	suite.Run("integration test ok response", func() {
		user := factory.BuildDefaultUser(suite.DB())

		params := userop.GetUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/users/%s", user.ID)),
			UserID:      strfmt.UUID(user.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetUserHandler{
			suite.NewHandlerConfig(),
			userservice.NewUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.GetUserOK{}, response)
		okResponse := response.(*userop.GetUserOK)
		suite.Equal(user.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		user := factory.BuildDefaultUser(suite.DB())
		params := userop.GetUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/users/%s", user.ID)),
			UserID:      strfmt.UUID(user.ID.String()),
		}
		expectedError := models.ErrFetchNotFound
		userFetcher := &mocks.UserFetcher{}
		userFetcher.On("FetchUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.User{}, expectedError).Once()
		handler := GetUserHandler{
			suite.NewHandlerConfig(),
			userFetcher,
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

func (suite *HandlerSuite) TestIndexUsersHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		users := models.Users{
			factory.BuildDefaultUser(suite.DB()),
			factory.BuildDefaultUser(suite.DB()),
		}
		params := userop.IndexUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexUsersHandler{
			HandlerConfig:  suite.NewHandlerConfig(),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.IndexUsersOK{}, response)
		okResponse := response.(*userop.IndexUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(users[0].ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := userop.IndexUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/users"),
		}
		expectedError := models.ErrFetchNotFound
		userListFetcher := &mocks.ListFetcher{}
		userListFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		userListFetcher.On("FetchRecordCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexUsersHandler{
			HandlerConfig:  suite.NewHandlerConfig(),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    userListFetcher,
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

func (suite *HandlerSuite) TestUpdateUserHandler() {
	// Set constants for the session names
	milSessionID := "mil-session"
	adminSessionID := "admin-session"
	officeSessionID := "office-session"

	activeUserTrait := func() []factory.Customization {
		return []factory.Customization{
			{
				Model: models.User{
					CurrentMilSessionID:    milSessionID,
					CurrentAdminSessionID:  adminSessionID,
					CurrentOfficeSessionID: officeSessionID,
					Active:                 true,
				},
			},
		}
	}

	// Create a handler and service object instances to test
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	queryBuilder := query.NewQueryBuilder()
	officeUpdater := officeuser.NewOfficeUserUpdater(queryBuilder)
	adminUpdater := adminuser.NewAdminUserUpdater(queryBuilder)

	setupHandler := func() UpdateUserHandler {
		handlerConfig := suite.NewHandlerConfig()

		return UpdateUserHandler{
			handlerConfig,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater, suite.TestNotificationSender()),
			newQueryFilter,
		}
	}

	suite.Run("Successful userSessionRevocation", func() {
		// Under test: UsereSessionRevocation, userUpdater
		// Mocked: 	   QueryFilter
		// Set up:     We revoke two sessions from an existing user.
		//			   The user active state is set to nil, so user should stay active.
		// Expected outcome:
		//             CurrentMilSessionID and CurrentOfficeSessionID are cleared from db.
		//			   CurrentAdminSessionID is still present.
		//             Active status is true.```

		// Create a user that has multiple sessions
		// and is labeled an Active user of the system
		user := factory.BuildUser(suite.DB(), nil, []factory.Trait{
			activeUserTrait,
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				RevokeMilSession:    models.BoolPointer(true),
				RevokeAdminSession:  models.BoolPointer(false),
				RevokeOfficeSession: models.BoolPointer(true),
				Active:              nil,
				OktaEmail:           &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := setupHandler().Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), user.ID)
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		// User is still active
		suite.Equal(true, foundUser.Active)

	})

	suite.Run("Successful userSessionRevocation and status update", func() {
		// Under test: UsereSessionRevocation, UserUpdater
		// Set up:     We pass in payload to revoke two sessions from an existing user
		//			   and deactivate the user.
		// Expected outcome:
		//             Active status is false.
		//			   *All* sessions are revoked because the user is deactivated.

		// Create a fresh user to reset all values
		user := factory.BuildUser(suite.DB(), nil, []factory.Trait{
			activeUserTrait,
		})

		// Create the update to revoke 2 sessions and deactivate the user
		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				Active:              models.BoolPointer(false),
				RevokeMilSession:    models.BoolPointer(true),
				RevokeOfficeSession: models.BoolPointer(true),
				OktaEmail:           &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		// Send request
		suite.NoError(params.User.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		// Check response
		foundUser, _ := models.GetUser(suite.DB(), user.ID)
		// The user is deactivated and all user sessions are revoked.
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(false, foundUser.Active)

	})

	suite.Run("Successful user deactivate, no sessions passed in", func() {
		// Under test: UpdateUser
		// Set up:     The user is active with sessions. No session properties are
		//			   included in the payload.
		// Expected outcome:
		//             Active status is false.
		//			   All sessions are revoked because the user is deactivated.

		// Create a fresh user to reset all values
		user := factory.BuildUser(suite.DB(), nil, []factory.Trait{
			activeUserTrait,
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				Active:    models.BoolPointer(false),
				OktaEmail: &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := setupHandler().Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), user.ID)
		// The user is deactivated and all user sessions are revoked.
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(false, foundUser.Active)

	})

	suite.Run("Successful user activate, no sessions passed in", func() {
		// Test UserUpdater
		// Under test: UpdateUser
		// Set up:     The user is marked active, but no session information
		//			   is included in the payload.
		// Expected outcome:
		//             Active status is true.
		// 			   Session IDs remain.
		//

		// Create a new user that is inactive but has session values
		user := factory.BuildUser(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					CurrentMilSessionID:    milSessionID,
					CurrentAdminSessionID:  adminSessionID,
					CurrentOfficeSessionID: officeSessionID,
					Active:                 false,
				},
			},
		}, nil)

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				Active:    models.BoolPointer(true),
				OktaEmail: &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		suite.NoError(params.User.Validate(strfmt.Default))
		response := setupHandler().Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), user.ID)

		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal(milSessionID, foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal(officeSessionID, foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)

	})

	suite.Run("Failed update with RevokeUserSession, successful update with userUpdater", func() {
		// Under test: UpdateUser
		// Mocked: UserSessionRevocation
		// Set up:     The session revocation fails, and the user's
		// 			   active status is successfully updated to false.
		// Expected outcome:
		//             Active status is false.
		// 			   Err returned for RevokeUserSession

		// Create a fresh user to reset all values
		user := factory.BuildUser(suite.DB(), nil, []factory.Trait{
			activeUserTrait,
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				Active:    models.BoolPointer(false),
				OktaEmail: &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		userRevocation := &mocks.UserSessionRevocation{}

		err := validate.NewErrors()

		handler := setupHandler()
		userRevocation.On("RevokeUserSession",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.User,
			handler.HandlerConfig.SessionManagers(),
		).Return(nil, err, nil).Once()

		handler.UserSessionRevocation = userRevocation

		suite.NoError(params.User.Validate(strfmt.Default))
		handler.Handle(params)
		foundUser, _ := models.GetUser(suite.DB(), user.ID)

		// Session update fails, active update succeeds
		suite.Error(err, "Error saving user")
		suite.Equal(false, foundUser.Active)
	})

	suite.Run("failed update with userUpdater", func() {
		user := factory.BuildUser(suite.DB(), nil, []factory.Trait{
			activeUserTrait,
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				Active:              models.BoolPointer(false),
				RevokeAdminSession:  models.BoolPointer(true),
				RevokeMilSession:    models.BoolPointer(true),
				RevokeOfficeSession: models.BoolPointer(true),
				OktaEmail:           &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		// Create a mock updater that returns an error
		userUpdater := &mocks.UserUpdater{}
		err := validate.NewErrors()

		userUpdater.On("UpdateUser",
			mock.AnythingOfType("*appcontext.appContext"),
			user.ID,
			mock.AnythingOfType("*models.User"),
		).Return(nil, nil, err).Once()

		handler := setupHandler()
		handler.UserUpdater = userUpdater

		suite.NoError(params.User.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&userop.UpdateUserInternalServerError{}, response)
	})

	suite.Run("Failed update with both RevokeUserSession and userUpdater", func() {
		// Mocked: UserUpdater, RevokeUserSession
		// Set up:     RevokeUser and updateUser return an err.
		// Expected outcome:
		//             Active status is unchanged and nil.
		// 			   User sessions are not revoked.

		// Create a fresh user to reset all values
		user := factory.BuildUser(suite.DB(), nil, []factory.Trait{
			activeUserTrait,
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdate{
				Active:              nil,
				RevokeAdminSession:  models.BoolPointer(true),
				RevokeMilSession:    models.BoolPointer(true),
				RevokeOfficeSession: models.BoolPointer(true),
				OktaEmail:           &user.OktaEmail,
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		// Create a mock that returns error on user session revocationand on user update
		userUpdater := &mocks.UserUpdater{}
		userRevocation := &mocks.UserSessionRevocation{}
		err := validate.NewErrors()

		handler := setupHandler()
		userRevocation.On("RevokeUserSession",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.User,
			handler.HandlerConfig.SessionManagers(),
		).Return(nil, err, nil).Once()

		userUpdater.On("UpdateUser",
			mock.AnythingOfType("*appcontext.appContext"),
			user.ID,
			mock.AnythingOfType("*models.User"),
		).Return(nil, nil, err).Once()

		handler.UserUpdater = userUpdater
		handler.UserSessionRevocation = userRevocation

		suite.NoError(params.User.Validate(strfmt.Default))
		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), user.ID)

		// Session update succeeds, active update fails
		suite.IsType(&userop.UpdateUserInternalServerError{}, response)
		suite.Equal(milSessionID, foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal(officeSessionID, foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)
	})
}

func (suite *HandlerSuite) TestDeleteUsersHandler() {
	suite.Run("deleted requested users results in no content (successful) response", func() {
		status := models.OfficeUserStatusAPPROVED
		userRole := roles.Role{
			RoleType: roles.RoleTypeTOO,
		}
		testOfficeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					Status: &status,
				},
			},
			{
				Model: models.User{
					Roles: []roles.Role{userRole},
				},
			},
		}, nil)
		testUser := testOfficeUser.User

		params := userop.DeleteUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/users/%s", testUser.ID)),
			UserID:      *handlers.FmtUUID(testUser.ID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteUserHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			UserDeleter:   userservice.NewUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&userop.DeleteUserNoContent{}, response)

		var dbUser models.User
		err := suite.DB().Where("id = ?", testUser.ID).First(&dbUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		var dbOfficeUser models.OfficeUser
		err = suite.DB().Where("user_id = ?", testUser.ID).First(&dbOfficeUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		// .All does not return a sql no rows error, so we will verify that the struct is empty
		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		var userPrivileges []models.UsersPrivileges
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")
	})

	suite.Run("get an error when the user does not exist", func() {
		userID := uuid.Must(uuid.NewV4())

		params := userop.DeleteUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/users/%s", userID)),
			UserID:      *handlers.FmtUUID(userID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteUserHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			UserDeleter:   userservice.NewUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&userop.DeleteUserNotFound{}, response)
	})

	suite.Run("error response when a user is not in the admin application", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		userID := officeUser.ID
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", userID), nil)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          userID,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := userop.DeleteUserParams{
			HTTPRequest: req.WithContext(ctx),
			UserID:      *handlers.FmtUUID(userID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteUserHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			UserDeleter:   userservice.NewUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&userop.DeleteUserUnauthorized{}, response)
	})

	suite.Run("get an error when the user has a move", func() {
		userRole := roles.Role{
			RoleType: roles.RoleTypeCustomer,
		}
		testUser := factory.BuildUserAndUsersRoles(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{userRole},
				},
			},
		}, nil)
		_ = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    testUser,
				LinkOnly: true,
			},
		}, nil)
		userID := testUser.ID

		params := userop.DeleteUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/users/%s", userID)),
			UserID:      *handlers.FmtUUID(userID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteUserHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			UserDeleter:   userservice.NewUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&userop.DeleteUserConflict{}, response)
	})

	suite.Run("get an error when the user is an Admin", func() {
		userRole := roles.Role{
			RoleType: roles.RoleTypeHQ,
		}
		testUser := factory.BuildUserAndUsersRoles(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{userRole},
				},
			},
		}, nil)
		_ = factory.BuildAdminUser(suite.DB(), []factory.Customization{
			{
				Model: models.AdminUser{
					Active: true,
					UserID: &testUser.ID,
					Email:  testUser.OktaEmail,
				},
			},
			{
				Model:    testUser,
				LinkOnly: true,
			},
		}, nil)
		userID := testUser.ID

		params := userop.DeleteUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/users/%s", userID)),
			UserID:      *handlers.FmtUUID(userID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteUserHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			UserDeleter:   userservice.NewUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&userop.DeleteUserForbidden{}, response)
	})

}
