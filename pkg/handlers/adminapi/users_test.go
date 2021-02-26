package adminapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
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
	user := testdatagen.MakeDefaultUser(suite.DB())
	userIDString := user.ID.String()
	userID := user.ID

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s", userID), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := userop.GetUserParams{
			HTTPRequest: req,
			UserID:      strfmt.UUID(userIDString),
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := GetUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			userservice.NewUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.GetUserOK{}, response)
		okResponse := response.(*userop.GetUserOK)
		suite.Equal(userIDString, okResponse.Payload.ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		user := models.User{ID: userID}
		params := userop.GetUserParams{
			HTTPRequest: req,
			UserID:      strfmt.UUID(userIDString),
		}
		userFetcher := &mocks.UserFetcher{}
		userFetcher.On("FetchUser",
			mock.Anything,
		).Return(user, nil).Once()
		handler := GetUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			userFetcher,
			newQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.GetUserOK{}, response)
		okResponse := response.(*userop.GetUserOK)
		suite.Equal(userIDString, okResponse.Payload.ID.String())
	})

	suite.T().Run("unsuccessful response when fetch fails", func(t *testing.T) {
		params := userop.GetUserParams{
			HTTPRequest: req,
			UserID:      strfmt.UUID(userIDString),
		}
		expectedError := models.ErrFetchNotFound
		userFetcher := &mocks.UserFetcher{}
		userFetcher.On("FetchUser",
			mock.Anything,
		).Return(models.User{}, expectedError).Once()
		handler := GetUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			userFetcher,
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

func (suite *HandlerSuite) TestIndexUsersHandler() {
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
	req := httptest.NewRequest("GET", "/users", nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := userop.IndexUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.IndexUsersOK{}, response)
		okResponse := response.(*userop.IndexUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := userop.IndexUsersParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		userListFetcher := &mocks.ListFetcher{}
		userListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		userListFetcher.On("FetchRecordCount",
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexUsersHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
	milSessionID := "mil-session"
	adminSessionID := "admin-session"
	officeSessionID := "office-session"

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
	userID := user.ID

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	endpoint := fmt.Sprintf("/users/%s", userID)
	req := httptest.NewRequest("PUT", endpoint, nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	revokeMilSession := true
	revokeAdminSession := false
	revokeOfficeSession := true
	active := true
	deactivate := false
	revoke := true

	params := userop.UpdateUserParams{
		HTTPRequest: req,
		User: &adminmessages.UserUpdatePayload{
			RevokeMilSession:    &revokeMilSession,    // true
			RevokeAdminSession:  &revokeAdminSession,  // false
			RevokeOfficeSession: &revokeOfficeSession, // true
			Active:              &active,              // required value, true
		},
		UserID: strfmt.UUID(userID.String()),
	}

	sessionManagers := setupSessionManagers()
	handlerContext := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handlerContext.SetSessionManagers(sessionManagers)
	queryBuilder := query.NewQueryBuilder(suite.DB())
	officeUpdater := officeuser.NewOfficeUserUpdater(queryBuilder)
	adminUpdater := adminuser.NewAdminUserUpdater(queryBuilder)

	suite.T().Run("Successful userSessionRevocation", func(t *testing.T) {
		// Under test: UsereSessionRevocation, userUpdater
		// Set up:     We revoke two sessions from an existing user.
		//			   The user remains active.
		// Expected outcome:
		//             CurrentMilSessionID and CurrentOfficeSessionID are cleared from db.
		//			   CurrentAdminSessionID is still present.
		//             Active status is true.```
		queryBuilder := query.NewQueryBuilder(suite.DB())

		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), userID)
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		// User is still active
		suite.Equal(true, foundUser.Active)

	})

	suite.T().Run("Successful userSessionRevocation and status update", func(t *testing.T) {
		// Under test: UsereSessionRevocation, UserUpdater
		// Set up:     We pass in payload to revoke two sessions from an existing user
		//			   and deactivate the user.
		// Expected outcome:
		//             Active status is false.
		//			   *All* sessions are revoked because the user is deactivated.

		queryBuilder := query.NewQueryBuilder(suite.DB())

		// Create a fresh user to reset all values
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		userID := user.ID

		params := userop.UpdateUserParams{
			HTTPRequest: req,
			User: &adminmessages.UserUpdatePayload{
				Active:              &deactivate,
				RevokeMilSession:    &revoke,
				RevokeOfficeSession: &revoke,
			},
			UserID: strfmt.UUID(userID.String()),
		}

		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), userID)
		// The user is deactivated and all user sessions are revoked.
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(false, foundUser.Active)

	})

	suite.T().Run("Successful user deactivate, no sessions passed in", func(t *testing.T) {
		// Under test: UpdateUser
		// Set up:     The user is active with sessions. No session properties are
		//			   included in the payload.
		// Expected outcome:
		//             Active status is false.
		//			   All sessions are revoked because the user is deactivated.

		queryBuilder := query.NewQueryBuilder(suite.DB())

		// Create a fresh user to reset all values
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 true,
			},
		})

		userID := user.ID

		params := userop.UpdateUserParams{
			HTTPRequest: req,
			User: &adminmessages.UserUpdatePayload{
				Active: &deactivate,
			},
			UserID: strfmt.UUID(userID.String()),
		}

		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), userID)
		// The user is deactivated and all user sessions are revoked.
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(false, foundUser.Active)

	})

	suite.T().Run("Successful user activate, no sessions passed in", func(t *testing.T) {
		// Test UserUpdater
		// Under test: UpdateUser
		// Set up:     The user is marked active, but no session information
		//			   is included in the payload.
		// Expected outcome:
		//             Active status is true.
		// 			   Session IDs remain.
		//

		queryBuilder := query.NewQueryBuilder(suite.DB())

		// Create a new user that is inactive but has session values
		user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				CurrentMilSessionID:    milSessionID,
				CurrentAdminSessionID:  adminSessionID,
				CurrentOfficeSessionID: officeSessionID,
				Active:                 false,
			},
		})
		userID := user.ID

		// Manually update Active because of an issue with mergeModels in MakeUser
		user.Active = deactivate
		suite.DB().ValidateAndUpdate(&user)

		params := userop.UpdateUserParams{
			HTTPRequest: req,
			User: &adminmessages.UserUpdatePayload{
				Active: &active,
			},
			UserID: strfmt.UUID(userID.String()),
		}

		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), userID)

		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal(milSessionID, foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal(officeSessionID, foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)

	})

	suite.T().Run("Failed update with RevokeUserSession, successful update with userUpdater", func(t *testing.T) {
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

		userID := user.ID

		params := userop.UpdateUserParams{
			HTTPRequest: req,
			User: &adminmessages.UserUpdatePayload{
				Active: &deactivate,
			},
			UserID: strfmt.UUID(userID.String()),
		}

		userRevocation := &mocks.UserSessionRevocation{}
		queryBuilder := query.NewQueryBuilder(suite.DB())

		err := validate.NewErrors()

		userRevocation.On("RevokeUserSession",
			mock.Anything,
			params.User,
			sessionManagers[0].Store,
		).Return(nil, err, nil).Once()

		handler := UpdateUserHandler{
			handlerContext,
			userRevocation,
			userservice.NewUserUpdater(queryBuilder, officeUpdater, adminUpdater),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))
		handler.Handle(params)
		foundUser, _ := models.GetUser(suite.DB(), userID)

		// Session update fails, active update succeeds
		suite.Error(err, "Error saving user")
		suite.Equal(false, foundUser.Active)
	})

	suite.T().Run("Successful update with RevokeUserSession, failed update with userUpdater", func(t *testing.T) {
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

		userID := user.ID

		params := userop.UpdateUserParams{
			HTTPRequest: req,
			User: &adminmessages.UserUpdatePayload{
				Active:              &deactivate,
				RevokeAdminSession:  &revoke,
				RevokeMilSession:    &revoke,
				RevokeOfficeSession: &revoke,
			},
			UserID: strfmt.UUID(userID.String()),
		}

		userUpdater := &mocks.UserUpdater{}
		queryBuilder := query.NewQueryBuilder(suite.DB())

		err := validate.NewErrors()

		userUpdater.On("UpdateUser",
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, err).Once()

		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userUpdater,
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))
		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), userID)

		// Session update succeeds, active update fails
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)
	})

	suite.T().Run("Failed update with both RevokeUserSession and userUpdater", func(t *testing.T) {
		// Mocked: UserUpdater, RevokeUserSession
		// Set up:     RevokeUser and updateUser return an err.
		// Expected outcome:
		//             Active status is unchanged and true.
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

		userID := user.ID

		params := userop.UpdateUserParams{
			HTTPRequest: req,
			User: &adminmessages.UserUpdatePayload{
				Active:              &deactivate,
				RevokeAdminSession:  &revoke,
				RevokeMilSession:    &revoke,
				RevokeOfficeSession: &revoke,
			},
			UserID: strfmt.UUID(userID.String()),
		}

		userUpdater := &mocks.UserUpdater{}
		userRevocation := &mocks.UserSessionRevocation{}
		err := validate.NewErrors()

		userRevocation.On("RevokeUserSession",
			mock.Anything,
			params.User,
			sessionManagers[0].Store,
		).Return(nil, err, nil).Once()

		userUpdater.On("UpdateUser",
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, err).Once()

		handler := UpdateUserHandler{
			handlerContext,
			userRevocation,
			userUpdater,
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))
		response := handler.Handle(params)

		foundUser, _ := models.GetUser(suite.DB(), userID)

		// Session update succeeds, active update fails
		suite.IsType(&userop.UpdateUserInternalServerError{}, response)
		suite.Equal(milSessionID, foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal(officeSessionID, foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)
	})

}
