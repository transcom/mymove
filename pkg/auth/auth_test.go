package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/gorilla/context"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

const testHostname = "hostname"

func TestGenerateNonce(t *testing.T) {
	nonce := generateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}

var dbConnection *pop.Connection

func setupDBConnection() {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	dbConnection = conn
}

func getHandlerParamsWithCookie(ss string) (*httptest.ResponseRecorder, *http.Request) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	// Set a secure cookie on the recorder
	s, _ := cookieStore.Get(req, JwtCookieName)
	s.AddFlash(ss)
	s.IsNew = false
	s.Save(req, rr)

	// Calling Get sets the Set-Cookie header, so let's grab that cookie
	// and set it as the Cookie header to trick the middleware
	cookies := rr.Header()["Set-Cookie"]
	req.Header.Set("Cookie", cookies[0])
	// And refresh the recorder to get rid of the Set-Cookie header
	rr = httptest.NewRecorder()

	return rr, req
}

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

func TestMain(m *testing.M) {
	setupDBConnection()
	os.Exit(m.Run())
}

func TestAuthorizationLogoutHandler(t *testing.T) {
	fakeToken := "some_token"
	responsePattern := regexp.MustCompile(`href="(.+)"`)
	req, err := http.NewRequest("GET", "/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthorizationLogoutHandler(fmt.Sprintf("http://%s", testHostname)))

	context.Set(req, "id_token", fakeToken)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusTemporaryRedirect)
	}

	redirectURL, err := url.Parse(responsePattern.FindStringSubmatch(rr.Body.String())[1])
	if err != nil {
		t.Fatal(err)
	}
	params := redirectURL.Query()

	postRedirectURI, err := url.Parse(params["post_logout_redirect_uri"][0])
	if err != nil {
		t.Fatal(err)
	}

	if testHostname != postRedirectURI.Host {
		t.Errorf("handler returned wrong redirect URI hostname: got %v wanted %v", postRedirectURI.Host, testHostname)
	}

	fmt.Println(params)
	if token := params["id_token_hint"][0]; token != fakeToken {
		t.Errorf("handler returned wrong id_token: got %v wanted %v", token, fakeToken)
	}
}

func TestUserAuthMiddlewareWithBadToken(t *testing.T) {
	fakeToken := "some_token"
	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Error("error creating RSA key", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(pem, testHostname)(handler)

	rr, req := getHandlerParamsWithCookie(fakeToken)

	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusTemporaryRedirect)
	}

	if incomingToken, ok := context.Get(req, "id_token").(*string); ok {
		t.Errorf("expected id_token to be nil, got %v", incomingToken)
	}
}

func TestUserAuthMiddlewareWithValidToken(t *testing.T) {
	email := "some_email@domain.com"
	idToken := "fake_id_token"

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	// Brand new token, shouldn't be renewed
	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	ss, err := signedTokenStringWithUserInfo(email, idToken, expiry, pem)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(pem, testHostname)(handler)

	rr, req := getHandlerParamsWithCookie(ss)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be an ID token in the request context
	if incomingToken, ok := context.Get(req, "id_token").(string); !ok || incomingToken != idToken {
		t.Errorf("handler returned wrong id_token: got %v, wanted %v", incomingToken, idToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func TestUserAuthMiddlewareWithRenewalToken(t *testing.T) {
	email := "some_email@domain.com"
	idToken := "fake_id_token"

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	// Token will expire in 1 minute, should be renewed
	expiry := getExpiryTimeFromMinutes(1)
	ss, err := signedTokenStringWithUserInfo(email, idToken, expiry, pem)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
	})
	middleware := UserAuthMiddleware(pem, testHostname)(handler)

	rr, req := getHandlerParamsWithCookie(ss)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be an ID token in the request context
	if incomingToken, ok := context.Get(req, "id_token").(string); !ok || incomingToken != idToken {
		t.Errorf("handler returned wrong id_token: got %v, wanted %v", incomingToken, idToken)
	}

	// And the cookie should be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 1 {
		t.Errorf("expected 1 cookie to be set, got %v", len(setCookies))
	}
}

func TestUserAuthMiddlewareWithExpiredToken(t *testing.T) {
	email := "some_email@domain.com"
	idToken := "fake_id_token"

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	expiry := getExpiryTimeFromMinutes(-1)
	ss, err := signedTokenStringWithUserInfo(email, idToken, expiry, pem)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(pem, testHostname)(handler)

	rr, req := getHandlerParamsWithCookie(ss)

	middleware.ServeHTTP(rr, req)

	// We should be redirected to the landing page
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusTemporaryRedirect)
	}

	// And there should be no token passed through
	if incomingToken, ok := context.Get(req, "id_token").(string); ok {
		t.Errorf("expected id_token to be nil, got %v", incomingToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}
