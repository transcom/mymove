package dpsauth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const prefix = "mymove-"

// LoginGovIDToCookie takes the Login.gov UUID of the current user and returns the cookie.
func LoginGovIDToCookie(userID string, cookieSecret []byte, cookieExpires int) (*http.Cookie, error) {
	expirationTime := time.Now().Add(time.Minute * time.Duration(cookieExpires))

	claims := &jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString(cookieSecret)
	if err != nil {
		return nil, errors.Wrap(err, "Signing JWT")
	}

	cookie := http.Cookie{Value: prefix + jwt, Expires: expirationTime}
	return &cookie, nil
}

// CookieToLoginGovID takes a cookie value and returns the Login.gov UUID only if it's a
// valid, unexpired cookie.
func CookieToLoginGovID(cookieValue string, cookieSecret []byte) (string, error) {
	if !strings.HasPrefix(cookieValue, prefix) {
		return "", &ErrInvalidCookie{errMessage: "Invalid cookie: missing prefix"}
	}
	token, err := jwt.ParseWithClaims(cookieValue[len(prefix):], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cookieSecret, nil
	})

	if err != nil {
		return "", &ErrInvalidCookie{errMessage: fmt.Sprintf("Invalid cookie: unable to parse JWT - %s", err.Error())}
	}

	if token == nil || !token.Valid {
		return "", &ErrInvalidCookie{errMessage: "Invalid cookie: failed JWT validation"}
	}

	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Subject, nil
}
