package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
)

type webServerSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger logger
}

func TestWebServerSuite(t *testing.T) {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse([]string{})

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config(v.GetString("env"), v.GetBool("debug-logging"))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	fields := make([]zap.Field, 0)
	if len(gitBranch) > 0 {
		fields = append(fields, zap.String("git_branch", gitBranch))
	}
	if len(gitCommit) > 0 {
		fields = append(fields, zap.String("git_commit", gitCommit))
	}
	logger = logger.With(fields...)
	zap.ReplaceGlobals(logger)

	ss := &webServerSuite{
		viper:  v,
		logger: logger,
	}

	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filename := fmt.Sprintf("%s/config/env/%s.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		logger.Info(fmt.Sprintf("Loading environment variables from file %s", filename))
		ss.applyContext(ss.patchContext(ss.loadContext(filename)))
	}

	suite.Run(t, ss)
}

func (suite *webServerSuite) loadContext(variablesFile string) map[string]string {
	ctx := map[string]string{}
	if len(variablesFile) > 0 {
		// Read contents of variables file into vars
		vars, err := ioutil.ReadFile(variablesFile)
		if err != nil {
			suite.logger.Fatal(fmt.Sprintf("error reading variables from file %s", variablesFile))
		}

		// Adds variables from file into context
		for _, x := range strings.Split(string(vars), "\n") {
			// If a line is empty or starts with #, then skip.
			if len(x) > 0 && x[0] != '#' {
				// Split each line on the first equals sign into []string{name, value}
				pair := strings.SplitAfterN(x, "=", 2)
				ctx[pair[0][0:len(pair[0])-1]] = pair[1]
			}
		}
	}
	return ctx
}

func (suite *webServerSuite) patchContext(ctx map[string]string) map[string]string {
	for k, v := range ctx {
		if strings.HasPrefix(v, "/bin/") {
			ctx[k] = filepath.Join(os.Getenv("TEST_ACC_CWD"), v[1:])
		}
	}
	return ctx
}

func (suite *webServerSuite) applyContext(ctx map[string]string) {
	for k, v := range ctx {
		suite.logger.Info("overriding " + k)
		suite.viper.Set(strings.Replace(strings.ToLower(k), "_", "-", -1), v)
	}
}

func (suite *webServerSuite) TestConfigProtocols() {
	suite.Nil(checkProtocols(suite.viper))
}

func (suite *webServerSuite) TestConfigHosts() {
	suite.Nil(checkHosts(suite.viper))
}

func (suite *webServerSuite) TestConfigPorts() {
	suite.Nil(checkPorts(suite.viper))
}

func (suite *webServerSuite) TestConfigDPS() {
	suite.Nil(checkDPS(suite.viper))
}

func (suite *webServerSuite) TestConfigCSRF() {
	suite.Nil(checkCSRF(suite.viper))
}

func (suite *webServerSuite) TestConfigEmail() {
	suite.Nil(checkEmail(suite.viper))
}

func (suite *webServerSuite) TestConfigGEX() {
	suite.Nil(checkGEX(suite.viper))
}

func (suite *webServerSuite) TestConfigEIAKey() {
	suite.Nil(checkEIAKey(suite.viper))
}
func (suite *webServerSuite) TestConfigEIURL() {
	suite.Nil(checkEIAURL(suite.viper))
}

func (suite *webServerSuite) TestConfigStorage() {
	suite.Nil(checkStorage(suite.viper))
}

func (suite *webServerSuite) TestConfigDatabase() {
	suite.Nil(checkDatabase(suite.viper, suite.logger))
}

func (suite *webServerSuite) TestDODCertificates() {

	if os.Getenv("TEST_ACC_DOD_CERTIFICATES") != "1" {
		suite.logger.Info("skipping TestDODCertificates")
		return
	}

	_, _, err := initDODCertificates(suite.viper, suite.logger)
	suite.Nil(err)
}

func (suite *webServerSuite) TestHoneycomb() {

	if os.Getenv("TEST_ACC_HONEYCOMB") != "1" {
		suite.logger.Info("skipping TestHoneycomb")
		return
	}

	enabled := initHoneycomb(suite.viper, suite.logger)
	suite.True(enabled)
}

func (suite *webServerSuite) TestInitDatabase() {

	if os.Getenv("TEST_ACC_INIT_DATABASE") != "1" {
		suite.logger.Info("skipping TestInitDatabase")
		return
	}

	conn, err := initDatabase(suite.viper, suite.logger)
	suite.Nil(err)
	suite.NotNil(conn)
}

func (suite *webServerSuite) TestStaticReqMethodMiddleware() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := validMethodForStaticMiddleware(handler)

	// For a GET request, we should get a 200
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://mil.example.com/static/something", nil)
	middleware.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// For a HEAD request, we should also get a 200
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("HEAD", "http://mil.example.com/static/something", nil)
	middleware.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// For a POST request, we should get a 405 response
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "http://mil.example.com/static/something", nil)
	middleware.ServeHTTP(rr, req)
	suite.Equal(http.StatusMethodNotAllowed, rr.Code, "handler returned wrong status code")
}
