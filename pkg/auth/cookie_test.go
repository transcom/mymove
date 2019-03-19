package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func createRandomRSAPEM() (s string, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		err = errors.Wrap(err, "failed to generate key")
		return
	}

	asn1 := x509.MarshalPKCS1PrivateKey(priv)
	privBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: asn1,
	})
	s = string(privBytes[:])

	return
}

func getHandlerParamsWithToken(ss string, expiry time.Time) (*httptest.ResponseRecorder, *http.Request) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://mil.example.com/protected", nil)

	appName, _ := ApplicationName(req.Host, MilTestHost, OfficeTestHost, TspTestHost)

	// Set a secure cookie on the request
	cookieName := fmt.Sprintf("%s_%s", strings.ToLower(string(appName)), UserSessionCookieName)
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
	t := suite.T()
	fakeToken := "some_token"
	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Error("error creating RSA key", err)
	}

	var resultingSession *Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultingSession = SessionFromRequestContext(r)
	})
	middleware := SessionCookieMiddleware(suite.logger, pem, false, MilTestHost, OfficeTestHost, TspTestHost, false)(handler)

	expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)
	rr, req := getHandlerParamsWithToken(fakeToken, expiry)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And there should be no token passed through
	suite.NotNil(resultingSession, "Session should not be nil")
	suite.Equal("", resultingSession.IDToken, "Expected empty IDToken from bad cookie")
}

func (suite *authSuite) TestSessionCookieMiddlewareWithValidToken() {
	t := suite.T()
	email := "some_email@domain.com"
	idToken := "fake_id_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)
	incomingSession := Session{
		UserID:  fakeUUID,
		Email:   email,
		IDToken: idToken,
	}
	ss, err := signTokenStringWithUserInfo(expiry, &incomingSession, pem)
	if err != nil {
		t.Fatal(err)
	}
	rr, req := getHandlerParamsWithToken(ss, expiry)

	var resultingSession *Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultingSession = SessionFromRequestContext(r)
	})
	middleware := SessionCookieMiddleware(suite.logger, pem, false, MilTestHost, OfficeTestHost, TspTestHost, false)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And there should be an ID token in the request context
	suite.NotNil(resultingSession)
	suite.Equal(idToken, resultingSession.IDToken, "handler returned wrong id_token")

	// And the cookie should be renewed
	setCookies := rr.HeaderMap["Set-Cookie"]
	suite.Equal(1, len(setCookies), "expected cookie to be set")
}

func (suite *authSuite) TestSessionCookieMiddlewareWithExpiredToken() {
	t := suite.T()
	email := "some_email@domain.com"
	idToken := "fake_id_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	expiry := GetExpiryTimeFromMinutes(-1)
	incomingSession := Session{
		UserID:  fakeUUID,
		Email:   email,
		IDToken: idToken,
	}
	ss, err := signTokenStringWithUserInfo(expiry, &incomingSession, pem)
	if err != nil {
		t.Fatal(err)
	}

	var resultingSession *Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultingSession = SessionFromRequestContext(r)
	})
	middleware := SessionCookieMiddleware(suite.logger, pem, false, MilTestHost, OfficeTestHost, TspTestHost, false)(handler)

	rr, req := getHandlerParamsWithToken(ss, expiry)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And there should be no token passed through
	// And there should be no token passed through
	suite.NotNil(resultingSession)
	suite.Equal("", resultingSession.IDToken, "Expected empty IDToken from expired")
	suite.Equal(uuid.Nil, resultingSession.UserID, "Expected no UUID from expired cookie")

	// And the cookie should be set
	setCookies := rr.HeaderMap["Set-Cookie"]
	suite.Equal(1, len(setCookies), "expected cookie to be set")
}

func (suite *authSuite) TestSessionCookiePR161162731() {
	t := suite.T()
	email := "some_email@domain.com"
	idToken := "fake_id_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	expiry := GetExpiryTimeFromMinutes(SessionExpiryInMinutes)
	incomingSession := Session{
		UserID:  fakeUUID,
		Email:   email,
		IDToken: idToken,
	}
	ss, err := signTokenStringWithUserInfo(expiry, &incomingSession, pem)
	if err != nil {
		t.Fatal(err)
	}
	rr, req := getHandlerParamsWithToken(ss, expiry)

	var resultingSession *Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultingSession = SessionFromRequestContext(r)
		WriteSessionCookie(w, resultingSession, "freddy", false, suite.logger, false)
	})
	middleware := SessionCookieMiddleware(suite.logger, pem, false, MilTestHost, OfficeTestHost, TspTestHost, false)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And there should be an ID token in the request context
	suite.NotNil(resultingSession)
	suite.Equal(idToken, resultingSession.IDToken, "handler returned wrong id_token")

	// And the cookie should be renewed
	setCookies := rr.HeaderMap["Set-Cookie"]
	suite.Equal(1, len(setCookies), "expected cookie to be set")
}

func (suite *authSuite) TestMaskedCSRFMiddleware() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := MaskedCSRFMiddleware(suite.logger, false, false)(handler)
	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// And the cookie should be added to the session
	setCookies := rr.HeaderMap["Set-Cookie"]
	suite.Equal(1, len(setCookies), "expected cookie to be set")
}

func (suite *authSuite) TestMaskedCSRFMiddlewareCreatesOneToken() {
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
	middleware := MaskedCSRFMiddleware(suite.logger, false, false)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// No new cookie should be added to the session
	setCookies := rr.HeaderMap["Set-Cookie"]
	suite.Equal(0, len(setCookies), "expected no new cookie to be set")
}

func (suite *authSuite) TestMiddlewareConstructor() {
	adm := SessionCookieMiddleware(suite.logger, "secret", false, MilTestHost, OfficeTestHost, TspTestHost, false)
	suite.NotNil(adm)
}

func (suite *authSuite) TestMiddlewareMilApp() {
	rr := httptest.NewRecorder()

	milMoveTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.True(session.IsMilApp(), "first should be milmove app")
		suite.False(session.IsOfficeApp(), "first should not be office app")
		suite.False(session.IsTspApp(), "first should not be tsp app")
		suite.Equal(MilTestHost, session.Hostname)
	})
	milMoveMiddleware := SessionCookieMiddleware(suite.logger, "secret", false, MilTestHost, OfficeTestHost, TspTestHost, false)(milMoveTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", MilTestHost), nil)
	milMoveMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", MilTestHost), nil)
	milMoveMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(MilTestHost)), nil)
	milMoveMiddleware.ServeHTTP(rr, req)
}

func (suite *authSuite) TestMiddlwareOfficeApp() {
	rr := httptest.NewRecorder()

	officeTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMilApp(), "should not be milmove app")
		suite.True(session.IsOfficeApp(), "should be office app")
		suite.False(session.IsTspApp(), "should not be tsp app")
		suite.Equal(OfficeTestHost, session.Hostname)
	})
	officeMiddleware := SessionCookieMiddleware(suite.logger, "secret", false, MilTestHost, OfficeTestHost, TspTestHost, false)(officeTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", OfficeTestHost), nil)
	officeMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", OfficeTestHost), nil)
	officeMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(OfficeTestHost)), nil)
	officeMiddleware.ServeHTTP(rr, req)
}

func (suite *authSuite) TestMiddlwareTspApp() {
	rr := httptest.NewRecorder()

	tspTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequestContext(r)
		suite.False(session.IsMilApp(), "should not be milmove app")
		suite.False(session.IsOfficeApp(), "should not be office app")
		suite.True(session.IsTspApp(), "should be tsp app")
		suite.Equal(TspTestHost, session.Hostname)
	})
	tspMiddleware := SessionCookieMiddleware(suite.logger, "secret", false, MilTestHost, OfficeTestHost, TspTestHost, false)(tspTestHandler)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/some_url", TspTestHost), nil)
	tspMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", TspTestHost), nil)
	tspMiddleware.ServeHTTP(rr, req)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://%s:8080/some_url", strings.ToUpper(TspTestHost)), nil)
	tspMiddleware.ServeHTTP(rr, req)
}

func (suite *authSuite) TestMiddlewareBadApp() {
	rr := httptest.NewRecorder()

	noAppTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Fail("Should not be called")
	})
	noAppMiddleware := SessionCookieMiddleware(suite.logger, "secret", false, MilTestHost, OfficeTestHost, TspTestHost, false)(noAppTestHandler)

	req := httptest.NewRequest("GET", "http://totally.bogus.hostname/some_url", nil)
	noAppMiddleware.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code, "Should get an error ")
}
