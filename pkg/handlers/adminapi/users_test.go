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
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
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
	userActive := true

	params := userop.UpdateUserParams{
		HTTPRequest: req,
		User: &adminmessages.UserUpdatePayload{
			RevokeMilSession:    &revokeMilSession,
			RevokeAdminSession:  &revokeAdminSession,
			RevokeOfficeSession: &revokeOfficeSession,
			Active:              &userActive,
		},
		UserID: strfmt.UUID(userID.String()),
	}

	sessionManagers := setupSessionManagers()
	handlerContext := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handlerContext.SetSessionManagers(sessionManagers)

	suite.T().Run("Successful update of sessions only", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))

		response := handler.Handle(params)
		foundUser, _ := models.GetUser(suite.DB(), userID)

		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(true, foundUser.Active)
	})

	suite.T().Run("Successful update of setting active to false", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		userActive := false
		params.User.Active = &userActive

		handler := UpdateUserHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			userservice.NewUserUpdater(queryBuilder),
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))
		response := handler.Handle(params)
		foundUser, _ := models.GetUser(suite.DB(), userID)

		// When we set Active to false, all user sessions
		// should be revoked, regardless of what params were
		// passed in the payload.
		suite.IsType(&userop.UpdateUserOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal("", foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
		suite.Equal(false, foundUser.Active)
	})

	suite.T().Run("Failed update", func(t *testing.T) {
		userRevocation := &mocks.UserSessionRevocation{}
		userUpdater := &mocks.UserUpdater{}

		userRevocation.On("RevokeUserSession",
			mock.Anything,
			params.User,
			sessionManagers[0].Store,
		).Return(&user, nil, nil).Once()

		userUpdater.On("UpdateUser",
			mock.Anything,
			params.User,
		).Return(&user, nil, nil).Once()

		handler := UpdateUserHandler{
			handlerContext,
			userRevocation,
			userUpdater,
			newQueryFilter,
		}

		suite.NoError(params.User.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&userop.UpdateUserOK{}, response)
	})

	userRevocation := &mocks.UserSessionRevocation{}
	userUpdater := &mocks.UserUpdater{}
	err := validate.NewErrors()

	userRevocation.On("RevokeUserSession",
		mock.Anything,
		params.User,
		sessionManagers[0].Store,
	).Return(nil, err, nil).Once()

	userUpdater.On("UpdateUser",
		mock.Anything,
		params.User,
	).Return(nil, err, nil).Once()

	handler := UpdateUserHandler{
		handlerContext,
		userRevocation,
		userUpdater,
		newQueryFilter,
	}

	suite.NoError(params.User.Validate(strfmt.Default))
	handler.Handle(params)

	suite.Error(err, "Error saving user")

}
