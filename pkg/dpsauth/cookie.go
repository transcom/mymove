package dpsauth

import (
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var cookieExpiresInMinutes = initCookieExpiration()
var secretKey = initKey()

const prefix = "mymove-"

// ServiceMemberIDToCookie takes the service member UUID of the current user and returns the cookie value.
func ServiceMemberIDToCookie(userID string) (string, error) {
	expiration, err := strconv.Atoi(cookieExpiresInMinutes)
	if err != nil {
		return "", errors.Wrap(err, "Converting DPS_COOKIE_EXPIRES_IN_MINUTES to int")
	}
	expirationTime := time.Now().Add(time.Minute * time.Duration(expiration)).Unix()

	claims := &jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.Wrap(err, "Signing JWT")
	}
	return prefix + jwt, nil
}

// CookieToServiceMemberID takes a cookie value and returns the service member's UUID only if it's a
// valid, unexpired cookie.
func CookieToServiceMemberID(cookieValue string) (string, error) {
	if !strings.HasPrefix(cookieValue, prefix) {
		return "", errors.New("Invalid cookie: missing prefix")
	}
	token, err := jwt.ParseWithClaims(cookieValue[len(prefix):], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || token == nil || !token.Valid {
		return "", errors.Wrap(err, "Failed token validation")
	}

	claims := token.Claims.(*jwt.StandardClaims)

	if time.Now().Unix() > claims.ExpiresAt {
		return "", errors.New("Cookie is expired")
	}

	return claims.Subject, nil
}

func initCookieExpiration() string {
	return os.Getenv("DPS_COOKIE_EXPIRES_IN_MINUTES")
}

func initKey() []byte {
	return []byte(os.Getenv("DPS_AUTH_COOKIE_SECRET_KEY"))
}
