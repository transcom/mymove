//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package adminapi

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/stretchr/testify/mock"

	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	userservice "github.com/transcom/mymove/pkg/services/user"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func setupSessionManagers() [3]*scs.SessionManager {
	var milSession, adminSession, officeSession *scs.SessionManager
	store := memstore.New()
	milSession = scs.New()
	milSession.Store = store
	milSession.Cookie.Name = "mil_session_token"

	adminSession = scs.New()
	adminSession.Store = store
	adminSession.Cookie.Name = "admin_session_token"

	officeSession = scs.New()
	officeSession.Store = store
	officeSession.Cookie.Name = "office_session_token"

	return [3]*scs.SessionManager{milSession, adminSession, officeSession}
}

func (suite *HandlerSuite) TestGetUserHandler() {
	suite.Run("integration test ok response", func() {
		user := testdatagen.MakeDefaultUser(suite.DB())
		params := userop.GetUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/users/%s", user.ID)),
			UserID:      strfmt.UUID(user.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			userservice.NewUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.GetUserOK{}, response)
		okResponse := response.(*userop.GetUserOK)
		suite.Equal(user.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		user := testdatagen.MakeDefaultUser(suite.DB())
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
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
			testdatagen.MakeDefaultUser(suite.DB()),
			testdatagen.MakeDefaultUser(suite.DB()),
		}
		params := userop.IndexUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	// Create a handler and service object instances to test
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	sessionManagers := setupSessionManagers()
	queryBuilder := query.NewQueryBuilder()
	officeUpdater := officeuser.NewOfficeUserUpdater(queryBuilder)
	adminUpdater := adminuser.NewAdminUserUpdater(queryBuilder)

	setupHandler := func() UpdateUserHandler {
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		handlerContext.SetSessionManagers(sessionManagers)

		return UpdateUserHandler{
			handlerContext,
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
		assertions := testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		}
		user := testdatagen.MakeUser(suite.DB(), assertions)

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				RevokeMilSession:    swag.Bool(true),
				RevokeAdminSession:  swag.Bool(false),
				RevokeOfficeSession: swag.Bool(true),
				Active:              nil,
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
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		// Create the update to revoke 2 sessions and deactivate the user
		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				Active:              swag.Bool(false),
				RevokeMilSession:    swag.Bool(true),
				RevokeOfficeSession: swag.Bool(true),
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
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				Active: swag.Bool(false),
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
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 false,
			},
		})

		// Manually update Active because of an issue with mergeModels in MakeUser
		suite.DB().ValidateAndUpdate(&user)

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				Active: swag.Bool(true),
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
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				Active: swag.Bool(false),
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		userRevocation := &mocks.UserSessionRevocation{}

		err := validate.NewErrors()

		userRevocation.On("RevokeUserSession",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.User,
			sessionManagers[0].Store,
		).Return(nil, err, nil).Once()

		handler := setupHandler()
		handler.UserSessionRevocation = userRevocation

		suite.NoError(params.User.Validate(strfmt.Default))
		handler.Handle(params)
		foundUser, _ := models.GetUser(suite.DB(), user.ID)

		// Session update fails, active update succeeds
		suite.Error(err, "Error saving user")
		suite.Equal(false, foundUser.Active)
	})

	suite.Run("Successful update with RevokeUserSession, failed update with userUpdater", func() {
		// Under test: UserSessionRevocation
		// Mocked: UserUpdater
		// Set up:     The session revocation succeeds, and updateUser
		// 			   returns an error.
		// Expected outcome:
		//             Active status is unchanged and true.
		// 			   User sessions are revoked.

		// Create a fresh user to reset all values
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				Active:              swag.Bool(false),
				RevokeAdminSession:  swag.Bool(true),
				RevokeMilSession:    swag.Bool(true),
				RevokeOfficeSession: swag.Bool(true),
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

		foundUser, _ := models.GetUser(suite.DB(), user.ID)

		// Session update succeeds, active update fails
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)
	})

	suite.Run("Failed update with both RevokeUserSession and userUpdater", func() {
		// Mocked: UserUpdater, RevokeUserSession
		// Set up:     RevokeUser and updateUser return an err.
		// Expected outcome:
		//             Active status is unchanged and nil.
		// 			   User sessions are not revoked.

		// Create a fresh user to reset all values
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		params := userop.UpdateUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/users/%s", user.ID)),
			User: &adminmessages.UserUpdatePayload{
				Active:              nil,
				RevokeAdminSession:  swag.Bool(true),
				RevokeMilSession:    swag.Bool(true),
				RevokeOfficeSession: swag.Bool(true),
			},
			UserID: strfmt.UUID(user.ID.String()),
		}

		// Create a mock that returns error on user session revocationand on user update
		userUpdater := &mocks.UserUpdater{}
		userRevocation := &mocks.UserSessionRevocation{}
		err := validate.NewErrors()

		userRevocation.On("RevokeUserSession",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.User,
			sessionManagers[0].Store,
		).Return(nil, err, nil).Once()

		userUpdater.On("UpdateUser",
			mock.AnythingOfType("*appcontext.appContext"),
			user.ID,
			mock.AnythingOfType("*models.User"),
		).Return(nil, nil, err).Once()

		handler := setupHandler()
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
