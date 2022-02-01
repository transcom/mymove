package dpsauth

import (
	"net/url"

	"github.com/gofrs/uuid"
)

func (suite *dpsAuthSuite) TestToken() {
	t := suite.T()
	userID := uuid.Must(uuid.NewV4()).String()
	cookieName := "TEST_COOKIE"
	redirectURL := "https://www.example.com"
	secretKey := "abcd"

	token, err := GenerateToken(userID, cookieName, redirectURL, secretKey)
	if err != nil {
		t.Error("Error generating token", err)
		return
	}

	escaped := url.QueryEscape(token)
	claims, err := ParseToken(escaped, secretKey)
	if err != nil {
		t.Error("Error parsing claims from token", err)
		return
	}

	suite.Equal(userID, claims.RegisteredClaims.Subject)
	suite.Equal(cookieName, claims.CookieName)
	suite.Equal(redirectURL, claims.DPSRedirectURL)
}
