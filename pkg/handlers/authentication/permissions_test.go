package authentication

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	roleService "github.com/transcom/mymove/pkg/services/roles"
)

func (suite *AuthSuite) TestSwitchRolesSuccess() {
	tooPerms := GetPermissionsForRole(roles.RoleTypeTOO)

	setupOfficeUserAndIdentity := func(userRoles []roles.RoleType) (models.OfficeUser, roles.Role) {
		officeUser := factory.BuildOfficeUserWithRoles(
			suite.DB(),
			factory.GetTraitActiveOfficeUser(),
			userRoles,
		)

		// default role is HQ because of alphabetical sorting
		identity, err := models.FetchUserIdentity(suite.DB(), officeUser.User.OktaID)
		suite.FatalNoError(err)
		defaultRole, err := identity.Roles.Default()
		suite.FatalNoError(err)
		return officeUser, *defaultRole
	}

	suite.Run("Can successfully switch roles from HQ to TOO - 200", func() {
		officeUser, defaultRole := setupOfficeUserAndIdentity([]roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeHQ})

		// Handler session starts with HQ
		handlerSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *officeUser.UserID,
			IDToken:         "fake token",
			ActiveRole:      defaultRole,
		}

		// Patch it to be TOO
		hc := suite.NewHandlerConfig()
		sessionMgr := hc.SessionManagers().Office
		req := httptest.NewRequest("PATCH", "/auth/activeRole", nil)
		req = suite.SetupSessionRequest(req, &handlerSession, sessionMgr)
		handler := sessionMgr.LoadAndSave(NewActiveRoleUpdateHandler(
			suite.AuthContext(),
			hc,
			roleService.NewRolesFetcher(),
		))
		payload, _ := json.Marshal(map[string]string{"roleType": string(roles.RoleTypeTOO)})
		req.Body = io.NopCloser(bytes.NewReader(payload))
		req.Host = "office.example.com"

		rr := httptest.NewRecorder()

		// Submit
		handler.ServeHTTP(rr, req)

		// Expect 200 OK
		suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

		// Expect our session active role to now be TOO
		// memory
		suite.Equal(roles.RoleTypeTOO, handlerSession.ActiveRole.RoleType, "Handler memory session did not switch roles")
		// session store
		storeCtx, err := sessionMgr.Load(req.Context(), handlerSession.IDToken)
		suite.FatalNoError(err)
		suite.NotEmpty(storeCtx)

		unCastedServerSession := sessionMgr.Get(storeCtx, "session")
		castedServerSession, ok := unCastedServerSession.(*auth.Session)
		suite.FatalTrue(ok)
		suite.Equal(handlerSession.IDToken, castedServerSession.IDToken)
		suite.Equal(roles.RoleTypeTOO, castedServerSession.ActiveRole.RoleType)
		// Perms went from HQ -> TOO
		suite.Equal(tooPerms, handlerSession.Permissions)
		suite.Equal(tooPerms, castedServerSession.Permissions)
	})

	suite.Run("Error if context can't find session key - 500", func() {
		// In the scenario the session manager and handler context differ
		// in session keys, something bad has gone wrong, almost as if
		// two session managers were created by the backend
		officeUser, defaultRole := setupOfficeUserAndIdentity([]roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeHQ})

		// Handler session starts with HQ
		handlerSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *officeUser.UserID,
			IDToken:         "fake token",
			ActiveRole:      defaultRole,
		}

		// Patch it to be TOO

		sessionMgr := suite.NewHandlerConfig().SessionManagers().Office
		req := httptest.NewRequest("PATCH", "/auth/activeRole", nil)
		req = suite.SetupSessionRequest(req, &handlerSession, sessionMgr)
		handler := sessionMgr.LoadAndSave(NewActiveRoleUpdateHandler(
			suite.AuthContext(),
			suite.NewHandlerConfig(), // Creates a second, conflicting session manager causing the failure
			roleService.NewRolesFetcher(),
		))
		payload, _ := json.Marshal(map[string]string{"roleType": string(roles.RoleTypeTOO)})
		req.Body = io.NopCloser(bytes.NewReader(payload))
		req.Host = "office.example.com"

		rr := httptest.NewRecorder()

		// Submit
		handler.ServeHTTP(rr, req)

		// Expect 500
		suite.Equal(http.StatusInternalServerError, rr.Code, "handler returned the wrong status code")

		// Expect our session active role to not be TOO
		// memory
		suite.NotEqual(roles.RoleTypeTOO, handlerSession.ActiveRole.RoleType, "Handler memory session switched roles when it shouldn't have")
		// session store
		storeCtx, err := sessionMgr.Load(req.Context(), handlerSession.IDToken)
		suite.FatalNoError(err)
		suite.NotEmpty(storeCtx)

		unCastedServerSession := sessionMgr.Get(storeCtx, "session")
		castedServerSession, ok := unCastedServerSession.(*auth.Session)
		suite.FatalTrue(ok)
		suite.Equal(handlerSession.IDToken, castedServerSession.IDToken)
		suite.NotEqual(roles.RoleTypeTOO, castedServerSession.ActiveRole.RoleType, "Handler memory session switched roles when it shouldn't have")
		// Perms should not be TOO
		suite.NotEqual(tooPerms, handlerSession.Permissions)
		suite.NotEqual(tooPerms, castedServerSession.Permissions)
	})

	suite.Run("Error if requested a role not assigned to the user - 403", func() {
		// HQ user going to request TOO
		officeUser, defaultRole := setupOfficeUserAndIdentity([]roles.RoleType{roles.RoleTypeHQ})

		// Handler session starts with HQ
		handlerSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *officeUser.UserID,
			IDToken:         "fake token",
			ActiveRole:      defaultRole,
		}

		// Patch it to be TOO
		sessionMgr := suite.NewHandlerConfig().SessionManagers().Office
		req := httptest.NewRequest("PATCH", "/auth/activeRole", nil)
		req = suite.SetupSessionRequest(req, &handlerSession, sessionMgr)
		handler := sessionMgr.LoadAndSave(NewActiveRoleUpdateHandler(
			suite.AuthContext(),
			suite.NewHandlerConfig(), // Creates a second, conflicting session manager causing the failure
			roleService.NewRolesFetcher(),
		))
		payload, _ := json.Marshal(map[string]string{"roleType": string(roles.RoleTypeTOO)})
		req.Body = io.NopCloser(bytes.NewReader(payload))
		req.Host = "office.example.com"

		rr := httptest.NewRecorder()

		// Submit
		handler.ServeHTTP(rr, req)

		// Expect 403
		suite.Equal(http.StatusForbidden, rr.Code, "handler returned the wrong status code")

		// Expect our session active role to not be TOO
		// memory
		suite.NotEqual(roles.RoleTypeTOO, handlerSession.ActiveRole.RoleType, "Handler memory session switched roles when it shouldn't have")
		// session store
		storeCtx, err := sessionMgr.Load(req.Context(), handlerSession.IDToken)
		suite.FatalNoError(err)
		suite.NotEmpty(storeCtx)

		unCastedServerSession := sessionMgr.Get(storeCtx, "session")
		castedServerSession, ok := unCastedServerSession.(*auth.Session)
		suite.FatalTrue(ok)
		suite.Equal(handlerSession.IDToken, castedServerSession.IDToken)
		suite.NotEqual(roles.RoleTypeTOO, castedServerSession.ActiveRole.RoleType, "Handler memory session switched roles when it shouldn't have")
		// Perms should not be TOO
		suite.NotEqual(tooPerms, handlerSession.Permissions)
		suite.NotEqual(tooPerms, castedServerSession.Permissions)
	})
}
