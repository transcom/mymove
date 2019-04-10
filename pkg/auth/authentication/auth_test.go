package authentication

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/honeycombio/beeline-go/trace"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const (
	// TspTestHost
	TspTestHost string = "tsp.example.com"
	// OfficeTestHost
	OfficeTestHost string = "office.example.com"
	// MilTestHost
	MilTestHost string = "mil.example.com"
	// OrdersTestHost
	OrdersTestHost string = "orders.example.com"
	// DpsTestHost
	DpsTestHost string = "dps.example.com"
	// SddcTestHost
	SddcTestHost string = "sddc.example.com"
	// FakeRSAKey generated with `bin/generate-devlocal-cert.sh -o Test -u Application -n test.mil -f test`
	FakeRSAKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDB8iPy8nfNMBR6
6rlYOS9VyZYo2uS5AQ03yGDAOzID7/84P3KnvaIW2AGyYsyNNs/j+qcFm2Cr6LK+
bdDJaKHBMAwiOn3BIBkwobWeQ1R12odtGCJyEzjH0kHE3Trtw6ID8tVtzPhfIBhe
p3/lkRuERGC3JJ3gxZYLoh6CJRD5FRrBqQg93Dm4dKfIl2AftLi68o9zSYJoO2qC
iEPeco0+hN6Chu9qRwP3jhswYQbFDu65eW4vHVlKXW6E+34eGKV/cCO0Q86Gi7/8
XZ2tQjCE3eJybmy9BCQi2hM3VKOxHzThhYjpA2ae8wE2Ucm43giSC+L0b1jf279y
dY3S3b/RAgMBAAECggEAYoo2va94MyakoTc1aJ/Vbw73XlapM15XaupCTilFZj7A
O8Hw7U0qV9T0N8B/EZix07F8vxqM6YtXle2R0WN6G//filyRnFhEtDLVZk3rUd3w
RPuoNLGTfeNUS0PkNv3ZCYyN6DXmU96owx7zmp45juB3C1ZtaNC7RbnfKlzO3N56
eZuqVcgarA26JEpycyiEF1yWnRlwEpYFoeHWIyBNb/ssFrJO7fkQfqyaVR2MBsZa
HS43OFWhd8q43tmeFMpcHuh3j/AZ4TsvAPGMtcHRbyAeVD+7X9I6tdkmD58Gz5Di
HcuSQ3Y2GewtC0Uua+Fu+SMSnx7mxX9zafm2a1cuqQKBgQD4A3nIPnIrdr23W1u0
6XT9Ikb0sc6aMXHBb2HM+/HetulKIQ9O/ajZHHQqFdQlo0RjVgPCSJ9R860Lak29
3zPwVCjcs6lsf1QLlijxnZYHl8XpZ11bOf1QmSovGE9Qs06cl5ty8A7OC+dpwo4t
Yyi3J2jDGxFO8hRhL4my6varcwKBgQDIMPPCcGlMe73fU/78/HEocjV/1ZOXqEt7
GbRjMho1s1k+56c4G/wNLn5y7Y9oYSN9UqKswdgS5ALYWg5aY9LpCfgGOAmGMskt
lDEnUq2oV5/D3oF06FwJpX0OyNQKMgzrJXmXpfNWp7lpyfJPlWH04KpShyN4poX3
Pp9mrwdeqwKBgHVLl4YX2oEp2FHmeDnYi8bINky15yNPrSAx4ExE/8A4O58egZH3
L6r25Q2eY0YlsEtWu9Jf7FGi8D1M2lWpQXQxKV4v7jntAj+0lcqnn/QZWLWpeCKU
C3TZ63R4h9J/6vbuUMuMM0RJpvmC1SEsG257yfU0UPxIS1EnXXVr4Jt3AoGBAIdm
RJhQO4gVcZipUR9/BnIavQCXTdoXY+YAvrcQ3hVQFp6rQ7h5hQLNXY0SDBrHCJ/s
0kYSXbh5K0t1rZuJRM+FhJGAOUDg/JytTImSLA5eJZru1ZRizE1h9rGXN4Ml0wMA
N7tP7MPBcXCRvCgDm1tq0Qg8istBpf5SBrIG0+89AoGBAPbRiOsEZKGCfk/umkTp
0iPf4YhQWcRX8hQXdOQUlTyE1mXQRxQ8isSMF5FOfmpJufo2by5MmKoSK/DmquER
8EZVAV6/L2/k+6JcrMtdcNb0zklGOT4CqUtg1UM619dy2+MeOWiYvP3gJsyfSffV
NeWNl8nWD+2zOcRiBri5uUB8
-----END RSA PRIVATE KEY-----`
)

// ApplicationTestServername is a collection of the test servernames
func ApplicationTestServername() auth.ApplicationServername {
	appnames := auth.ApplicationServername{
		MilServername:    MilTestHost,
		OfficeServername: OfficeTestHost,
		TspServername:    TspTestHost,
		OrdersServername: OrdersTestHost,
		DpsServername:    DpsTestHost,
		SddcServername:   SddcTestHost,
	}
	return appnames
}

type AuthSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *AuthSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestAuthSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &AuthSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func fakeLoginGovProvider(logger Logger) LoginGovProvider {
	return NewLoginGovProvider("fakeHostname", "secret_key", logger)
}

func (suite *AuthSuite) TestGenerateNonce() {
	t := suite.T()
	nonce := generateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}

func (suite *AuthSuite) TestAuthorizationLogoutHandler() {
	t := suite.T()

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	callbackPort := 1234

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/auth/logout", OfficeTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := LogoutHandler{authContext, "fake key", false, false}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d wanted %d", status, http.StatusOK)
	}

	redirectURL, err := url.Parse(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}
	params := redirectURL.Query()

	postRedirectURI, err := url.Parse(params["post_logout_redirect_uri"][0])
	suite.Nil(err)
	suite.Equal(OfficeTestHost, postRedirectURI.Hostname())
	suite.Equal(strconv.Itoa(callbackPort), postRedirectURI.Port())
	token := params["id_token_hint"][0]
	suite.Equal(fakeToken, token, "handler id_token")
}

func (suite *AuthSuite) TestRequireAuthMiddleware() {
	// Given: a logged in user
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  loginGovUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.MustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token"}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	var handlerSession *auth.Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerSession = auth.SessionFromRequestContext(r)
	})
	middleware := UserAuthMiddleware(suite.logger)(handler)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're logged in
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, user.ID, "the authenticated user is different from expected")
}

func (suite *AuthSuite) TestRequireAuthMiddlewareUnauthorized() {
	t := suite.T()

	// Given: No logged in users
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(suite.logger)(handler)

	middleware.ServeHTTP(rr, req)

	// We should receive an unauthorized response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusUnauthorized)
	}
}

func (suite *AuthSuite) TestAuthorizeDisableUser() {
	userIdentity := models.UserIdentity{
		Disabled: true,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", OfficeTestHost), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           "disabled@example.com",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	h := CallbackHandler{
		authContext,
		suite.DB(),
		"fake key",
		false,
		false,
	}
	rr := httptest.NewRecorder()
	span := trace.Span{}
	authorizeKnownUser(&userIdentity, h, &session, rr, &span, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "authorizer did not recognize disabled user")
}

func (suite *AuthSuite) TestAuthKnownSingleRoleOffice() {
	officeUserID := uuid.Must(uuid.NewV4())
	tspUserID := uuid.Must(uuid.NewV4())
	userIdentity := models.UserIdentity{
		Disabled:     false,
		OfficeUserID: &officeUserID,
		TspUserID:    &tspUserID,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", OfficeTestHost), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	h := CallbackHandler{
		authContext,
		suite.DB(),
		"fake key",
		false,
		false,
	}
	rr := httptest.NewRecorder()
	span := trace.Span{}
	authorizeKnownUser(&userIdentity, h, &session, rr, &span, req.WithContext(ctx), "")

	// Office app, so should only have office ID information
	suite.Equal(officeUserID, session.OfficeUserID)
	suite.Equal(uuid.Nil, session.TspUserID)
}

func (suite *AuthSuite) TestAuthKnownSingleRoleTSP() {
	officeUserID := uuid.Must(uuid.NewV4())
	tspUserID := uuid.Must(uuid.NewV4())
	userIdentity := models.UserIdentity{
		Disabled:     false,
		OfficeUserID: &officeUserID,
		TspUserID:    &tspUserID,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", TspTestHost), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	h := CallbackHandler{
		authContext,
		suite.DB(),
		"fake key",
		false,
		false,
	}
	rr := httptest.NewRecorder()
	span := trace.Span{}
	authorizeKnownUser(&userIdentity, h, &session, rr, &span, req.WithContext(ctx), "")

	// TSP app, so should only have TSP ID information
	suite.Equal(tspUserID, session.TspUserID)
	suite.Equal(uuid.Nil, session.OfficeUserID)
}
