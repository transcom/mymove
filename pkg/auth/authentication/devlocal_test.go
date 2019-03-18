package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func getCookie(name string, cookies []*http.Cookie) (*http.Cookie, error) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}
	return nil, errors.Errorf("Unable to find cookie: %s", name)
}

func (suite *AuthSuite) TestCreateUserHandler() {
	t := suite.T()

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	callbackPort := 1234

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", MilTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		UserID:          fakeUUID,
		Hostname:        MilTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateUserHandler(authContext, suite.DB(), "fake key", false, false)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("mil_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}

	user := models.User{}
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Error("Could not unmarshal json data into User model.", err)
	}
}

func (suite *AuthSuite) TestCreateAndLoginUserHandler() {
	t := suite.T()

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	callbackPort := 1234

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/devlocal-auth/new", MilTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		UserID:          fakeUUID,
		Hostname:        MilTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := NewCreateAndLoginUserHandler(authContext, suite.DB(), "fake key", false, false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req.WithContext(ctx))

	suite.Equal(http.StatusSeeOther, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusSeeOther)
	}

	cookies := rr.Result().Cookies()
	if _, err := getCookie("mil_session_token", cookies); err != nil {
		t.Error("could not find session token in response")
	}
}
