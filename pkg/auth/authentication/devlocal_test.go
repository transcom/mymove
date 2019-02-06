package authentication

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *AuthSuite) TestCreateUserHandler() {
	t := suite.T()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/devlocal-auth/create", nil)

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", 1234)
	handler := CreateUserHandler{authContext, suite.DB(), "fake key", false}
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	user := models.User{}
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Error("Could not unmarshal json data into User model.", err)
	}
}
