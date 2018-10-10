package dpsauth

import (
	"net/url"
	"os"
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/suite"
)

type dpsAuthSuite struct {
	suite.Suite
}

func (suite *dpsAuthSuite) SetupTest() {
	key := os.Getenv("DPS_AUTH_COOKIE_SECRET_KEY")
	if len(key) == 0 {
		suite.T().Fatal("You must set the DPS_AUTH_COOKIE_SECRET_KEY environment variable to run this test")
	}
}

func (suite *dpsAuthSuite) TestCookie() {
	t := suite.T()
	userID := uuid.Must(uuid.NewV4()).String()
	cookie, err := UserIDToCookie(userID)
	if err != nil {
		t.Error("Error generating cookie value from user ID", err)
	}

	// Mimic cookie being passed back as an API param
	escaped := url.QueryEscape(cookie)
	userIDFromCookie, err := CookieToUserID(escaped)
	if err != nil {
		t.Error("Error extracting user ID from cookie value", err)
	}
	suite.Equal(userID, userIDFromCookie)
}

func TestDPSAuthSuite(t *testing.T) {
	s := &dpsAuthSuite{}
	suite.Run(t, s)
}
