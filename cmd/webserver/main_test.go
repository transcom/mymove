package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/logging"
)

type webServerSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger *webserverLogger
}

func TestWebServerSuite(t *testing.T) {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse([]string{})

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	zapLogger, err := logging.Config(v.GetString("env"), v.GetBool("debug-logging"))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	logger := &webserverLogger{zapLogger}

	ss := &webServerSuite{
		viper:  v,
		logger: logger,
	}

	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filename := fmt.Sprintf("%s/config/env/%s.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		logger.Info(fmt.Sprintf("Loading environment variables from file %s", filename))
		ss.applyContext(ss.loadContext(filename))
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

func (suite *webServerSuite) TestConfigStorage() {
	suite.Nil(checkStorage(suite.viper))
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

func (suite *webServerSuite) TestDatabase() {

	if os.Getenv("TEST_ACC_DATABASE") != "1" {
		suite.logger.Info("skipping TestDatabase")
		return
	}

	_, err := initDatabase(suite.viper, suite.logger)
	suite.Nil(err)
}
