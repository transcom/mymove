package auth

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

func (suite *authSuite) SetupTest() {
	gob.Register(Session{})
}

func setupSessionManagers() AppSessionManagers {
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

	return AppSessionManagers{
		mil: ScsSessionManagerWrapper{
			ScsSessionManager: milSession,
		},
		office: ScsSessionManagerWrapper{
			ScsSessionManager: officeSession,
		},
		admin: ScsSessionManagerWrapper{
			ScsSessionManager: adminSession,
		},
	}
}

func getHandlerParamsWithToken(ss string, expiry time.Time) (*httptest.ResponseRecorder, *http.Request) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://mil.example.com/protected", nil)

	appnames := ApplicationTestServername()
	appName, _ := ApplicationName(req.Host, appnames)

	// Set a secure cookie on the request
	cookieName := fmt.Sprintf("%s_%s", strings.ToLower(string(appName)), "session_token")
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   ss,
		Path:    "/",
		Expires: expiry,
	}
	req.AddCookie(&cookie)
	return rr, req
}

func (suite *authSuite) TestSessionCookieMiddlewareWithBadToken() {
	fakeToken := "some_token"
	sessionManagers := setupSessionManagers()
	milSession := sessionManagers.mil

	var resultingSession *Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultingSession = SessionFromRequestContext(r)
	})
	appnames := ApplicationTestServername()
	middleware := SessionCookieMiddleware(suite.logger, appnames, sessionManagers)(handler)

	expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)
	rr, req := getHandlerParamsWithToken(fakeToken, expiry)

	milSession.LoadAndSave(middleware).ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And there should be no token passed through
	suite.NotNil(resultingSession, "Session should not be nil")
	suite.Equal("", resultingSession.IDToken, "Expected empty IDToken from bad cookie")
}

func (suite *authSuite) TestMaskedCSRFMiddleware() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := MaskedCSRFMiddleware(suite.logger, false)(handler)
	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And the cookie should be added to the session
	setCookies := rr.Result().Cookies()
	suite.Equal(1, len(setCookies), "expected cookie to be set")
}

func (suite *authSuite) TestMaskedCSRFMiddlewareCreatesNewToken() {
	expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Set a secure cookie on the request
	cookie := http.Cookie{
		Name:    MaskedGorillaCSRFToken,
		Value:   "fakecsrftoken",
		Path:    "/",
		Expires: expiry,
	}
	req.AddCookie(&cookie)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := MaskedCSRFMiddleware(suite.logger, false)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// No new cookie should be added to the session
	setCookies := rr.Result().Cookies()
	suite.Equal(1, len(setCookies), "expected a new cookie to be set")
}

func (suite *authSuite) TestMiddlewareConstructor() {
	appnames := ApplicationTestServername()
	sessionManagers := setupSessionManagers()

	adm := SessionCookieMiddleware(suite.logger, appnames, sessionManagers)
	suite.NotNil(adm)
}

func (suite *authSuite) TestMiddlewareMilApp() {
	rr := httptest.NewRecorder()

	appnames := ApplicationTestServername()
	milMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.True(session.IsMilApp(), "first should be milmove app")
		suite.False(session.IsOfficeApp(), "first should not be office app")
		suite.False(session.IsAdminApp(), "first should not be admin app")
		suite.Equal(appnames.MilServername, session.Hostname)
	})
	sessionManagers := setupSessionManagers()
	milSession := sessionManagers.mil
	milMoveMiddleware := SessionCookieMiddleware(suite.logger, appnames, sessionManagers)(milMoveTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", appnames.MilServername), nil)
	milSession.LoadAndSave(milMoveMiddleware).ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", appnames.MilServername), nil)
	milSession.LoadAndSave(milMoveMiddleware).ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(appnames.MilServername)), nil)
	milSession.LoadAndSave(milMoveMiddleware).ServeHTTP(rr, req)
}

func (suite *authSuite) TestMiddlwareOfficeApp() {
	rr := httptest.NewRecorder()

	appnames := ApplicationTestServername()
	officeTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMilApp(), "should not be milmove app")
		suite.True(session.IsOfficeApp(), "should be office app")
		suite.False(session.IsAdminApp(), "should not be admin app")
		suite.Equal(appnames.OfficeServername, session.Hostname)
	})
	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers.office
	officeMiddleware := SessionCookieMiddleware(suite.logger, appnames, sessionManagers)(officeTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", appnames.OfficeServername), nil)
	officeSession.LoadAndSave(officeMiddleware).ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", appnames.OfficeServername), nil)
	officeSession.LoadAndSave(officeMiddleware).ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(appnames.OfficeServername)), nil)
	officeSession.LoadAndSave(officeMiddleware).ServeHTTP(rr, req)
}

func (suite *authSuite) TestMiddlwareAdminApp() {
	rr := httptest.NewRecorder()

	appnames := ApplicationTestServername()
	adminTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMilApp(), "should not be milmove app")
		suite.False(session.IsOfficeApp(), "should not be office app")
		suite.True(session.IsAdminApp(), "should be admin app")
		suite.Equal(AdminTestHost, session.Hostname)
	})
	sessionManagers := setupSessionManagers()
	adminSession := sessionManagers.admin
	adminMiddleware := SessionCookieMiddleware(suite.logger, appnames, sessionManagers)(adminTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", AdminTestHost), nil)
	adminSession.LoadAndSave(adminMiddleware).ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", AdminTestHost), nil)
	adminSession.LoadAndSave(adminMiddleware).ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(AdminTestHost)), nil)
	adminSession.LoadAndSave(adminMiddleware).ServeHTTP(rr, req)
}

func (suite *authSuite) TestMiddlewareBadApp() {
	rr := httptest.NewRecorder()

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	appnames := ApplicationTestServername()
	sessionManagers := setupSessionManagers()
	milSession := sessionManagers.mil
	noAppMiddleware := SessionCookieMiddleware(suite.logger, appnames, sessionManagers)(noAppTestHandler)

	req := httptest.NewRequest("GET", "http://totally.bogus.hostname/some_url", nil)
	milSession.LoadAndSave(noAppMiddleware).ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
