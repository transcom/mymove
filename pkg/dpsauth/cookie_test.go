package dpsauth

import (
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type dpsAuthSuite struct {
	testingsuite.BaseTestSuite
	logger *zap.Logger
}

func (suite *dpsAuthSuite) SetupSuite() {
	key := os.Getenv("DPS_AUTH_COOKIE_SECRET_KEY")
	exp := os.Getenv("DPS_COOKIE_EXPIRES_IN_MINUTES")
	if len(key) == 0 || len(exp) == 0 {
		suite.T().Fatal("You must set the DPS_AUTH_COOKIE_SECRET_KEY and DPS_COOKIE_EXPIRES_IN_MINUTES environment variables to run this test")
	}
}

func (suite *dpsAuthSuite) TestCookie() {
	t := suite.T()
	userID := uuid.Must(uuid.NewV4()).String()
	cookie, err := LoginGovIDToCookie(userID)
	if err != nil {
		t.Error("Error generating cookie value from user ID", err)
	}

	// Mimic cookie being passed back in an API call via query param
	escaped := url.QueryEscape(cookie.Value)
	userIDFromCookie, err := CookieToLoginGovID(escaped)
	if err != nil {
		t.Error("Error extracting user ID from cookie value", err)
	}
	suite.Equal(userID, userIDFromCookie)
}

func TestDPSAuthSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	s := &dpsAuthSuite{logger: logger}
	suite.Run(t, s)
}
