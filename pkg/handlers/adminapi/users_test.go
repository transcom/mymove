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
	"github.com/stretchr/testify/mock"

	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
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

func (suite *HandlerSuite) TestRevokeUserSessionHandler() {
	milSessionID := "mil-session"
	adminSessionID := "admin-session"
	officeSessionID := "office-session"

	assertions := testdatagen.Assertions{
		User: models.User{
			CurrentMilSessionID:    milSessionID,
			CurrentAdminSessionID:  adminSessionID,
			CurrentOfficeSessionID: officeSessionID,
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

	params := userop.RevokeUserSessionParams{
		HTTPRequest: req,
		User: &adminmessages.UserRevokeSessionPayload{
			RevokeMilSession:    &revokeMilSession,
			RevokeAdminSession:  &revokeAdminSession,
			RevokeOfficeSession: &revokeOfficeSession,
		},
		UserID: strfmt.UUID(userID.String()),
	}

	sessionManagers := setupSessionManagers()
	handlerContext := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handlerContext.SetSessionManagers(sessionManagers)

	suite.T().Run("Successful update", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := RevokeUserSessionHandler{
			handlerContext,
			userservice.NewUserSessionRevocation(queryBuilder),
			newQueryFilter,
		}

		response := handler.Handle(params)
		foundUser, _ := models.GetUser(suite.DB(), userID)

		suite.IsType(&userop.RevokeUserSessionOK{}, response)
		suite.Equal("", foundUser.CurrentMilSessionID)
		suite.Equal(adminSessionID, foundUser.CurrentAdminSessionID)
		suite.Equal("", foundUser.CurrentOfficeSessionID)
	})

	suite.T().Run("Failed update", func(t *testing.T) {
		userUpdater := &mocks.UserSessionRevocation{}

		userUpdater.On("RevokeUserSession",
			mock.Anything,
			params.User,
			sessionManagers[0].Store,
		).Return(&user, nil, nil).Once()

		handler := RevokeUserSessionHandler{
			handlerContext,
			userUpdater,
			newQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.RevokeUserSessionOK{}, response)
	})

	userUpdater := &mocks.UserSessionRevocation{}
	err := validate.NewErrors()

	userUpdater.On("RevokeUserSession",
		mock.Anything,
		params.User,
		sessionManagers[0].Store,
	).Return(nil, err, nil).Once()

	handler := RevokeUserSessionHandler{
		handlerContext,
		userUpdater,
		newQueryFilter,
	}

	handler.Handle(params)

	suite.Error(err, "Error saving user")
}
