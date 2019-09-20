package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ServeSuite struct {
	testingsuite.PopTestSuite
	logger logger
}

func (suite *ServeSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestServeSuite(t *testing.T) {
	flag := pflag.CommandLine

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config(v.GetString(cli.DbEnvFlag), v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Panic(err)
	}
	hs := &ServeSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func (suite *ServeSuite) TestIsLoggedInWhenNoUserLoggedIn() {
	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(isLoggedIn(suite.logger))

	handler.ServeHTTP(rr, req)

	// expects to return 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

	// expects to return that no one is logged in
	expected := "{\"isLoggedIn\":false}\n"
	suite.Equal(expected, rr.Body.String(), "handler returned wrong body")
}

func (suite *ServeSuite) TestIsLoggedInWhenUserLoggedIn() {
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  loginGovUUID,
		LoginGovEmail: "email@example.com",
		Disabled:      false,
	}
	suite.MustSave(&user)

	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token"}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(isLoggedIn(suite.logger))

	handler.ServeHTTP(rr, req)

	// expects to return 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

	// expects to return that no one is logged in
	expected := "{\"isLoggedIn\":true}\n"
	suite.Equal(expected, rr.Body.String(), "handler returned wrong body")
}
