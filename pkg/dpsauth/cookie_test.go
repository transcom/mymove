package dpsauth

import (
	"log"
	"net/url"
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

func (suite *dpsAuthSuite) TestCookie() {
	t := suite.T()
	userID := uuid.Must(uuid.NewV4()).String()
	cookieSecret := []byte("j-7oWD_dOnhVf$PpQLRkMxaLmFDj!aE$")
	cookieExpires := 240
	cookie, err := LoginGovIDToCookie(userID, cookieSecret, cookieExpires)
	if err != nil {
		t.Error("Error generating cookie value from user ID", err)
	}

	// Mimic cookie being passed back in an API call via query param
	escaped := url.QueryEscape(cookie.Value)
	userIDFromCookie, err := CookieToLoginGovID(escaped, cookieSecret)
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
