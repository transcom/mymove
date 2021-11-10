package dpsauth

import (
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type dpsAuthSuite struct {
	testingsuite.PopTestSuite
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
	s := &dpsAuthSuite{}
	suite.Run(t, s)
}
