package adminapi

import (
	"net/http/httptest"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func (suite *HandlerSuite) TestConnectAdminPG() {
	defer func() { adminDBURL = "" }()
	adminDBURL = "badurl"
	conn := connectAdminDB(suite.Logger())
	suite.Nil(conn)

	adminDBURL = "postgres://postgres@localhost:5432/db_test"
	dbSSLRootCertFile = "config/tls/devlocal-ca.pem"
	conn = connectAdminDB(suite.Logger())
	suite.Nil(conn)
}

func (suite *HandlerSuite) TestDebugMiddleware() {
	m := NewAdminDebugMiddleware(suite.Logger())
	h := NewEnvHandler(suite.Logger())
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/admin/debug/env", strings.NewReader(`{"name":"DB_USER"}`))
	mh := m(h)
	mh.ServeHTTP(w, req)
	suite.Equal(403, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/debug/env", strings.NewReader(`{"name":"DB_USER"}`))
	req.Header.Add("X-Admin-Debug", os.Getenv("DB_HOST"))
	mh.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)
}

func (suite *HandlerSuite) TestEnvHandler() {
	h := NewEnvHandler(suite.Logger())
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/admin/debug/env", strings.NewReader(`{"name":"DB_USER"}`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/debug/env", strings.NewReader(`""`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)
}

func (suite *HandlerSuite) TestPGHandler() {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	NewAdminDB(v)
	h := NewPGHandler(suite.Logger())
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/admin/debug/pg", strings.NewReader(`{"query":"SELECT now()"}`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/debug/pg", strings.NewReader(`{"query":"SELECT"}`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/debug/pg", strings.NewReader(`{"exec":"COMMENT ON TABLE addresses IS 'funky'"}`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/debug/pg", strings.NewReader(`{"exec":"COMMENT bad"}`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/debug/pg", strings.NewReader(`""`))
	h.ServeHTTP(w, req)
	suite.Equal(200, w.Result().StatusCode)
}
