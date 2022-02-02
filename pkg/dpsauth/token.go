package dpsauth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// GenerateToken generates the DPS auth token passed to the endpoint that sets the cookie
func GenerateToken(loginGovID string, cookieName string, dpsRedirectURL string, secretKey string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: loginGovID, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
		CookieName:       cookieName,
		DPSRedirectURL:   dpsRedirectURL,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "Signing JWT")
	}
	return jwt, nil
}

// ParseToken parses the token string into its claims
func ParseToken(token string, secretKey string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Parsing JWT")
	}

	if parsedToken == nil || !parsedToken.Valid {
		return nil, errors.New("Invalid DPS auth token")
	}

	return parsedToken.Claims.(*Claims), nil
}
