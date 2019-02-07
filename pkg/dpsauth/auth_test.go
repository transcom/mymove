package dpsauth

import (
	"net/http"
	"net/http/httptest"
)

func (suite *dpsAuthSuite) TestSetCookieHandler() {
	secretKey := "secret"
	dpsCookieName := "DPS_COOKIE"
	cookieDomain := "sddctest"
	cookieSecret := []byte("j-7oWD_dOnhVf$PpQLRkMxaLmFDj!aE$")
	cookieExpires := 240
	handler := NewSetCookieHandler(suite.logger, secretKey, cookieDomain, cookieSecret, cookieExpires)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dps_auth/set_cookie", nil)
	token, _ := GenerateToken("uuid", dpsCookieName, "www.example.com", secretKey)
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusSeeOther, rr.Code)

	cookies := rr.HeaderMap["Set-Cookie"]
	suite.Equal(2, len(cookies))

	suite.Contains(cookies[0], dpsCookieName)
	suite.Contains(cookies[0], cookieDomain)
	suite.Contains(cookies[1], "DPSETAROLE=dodcustomer")
	suite.Contains(cookies[1], cookieDomain)
}
