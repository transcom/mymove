package cli

import (
	"os"
	"path/filepath"
	"runtime"
)

func (suite *cliTestSuite) TestConfigAuth() {
	// create fake data for the auth test
	// use https://pkg.go.dev/testing#B.Setenv for cleanup
	suite.T().Setenv("LOGIN_GOV_ADMIN_CLIENT_ID",
		"urn:gov:gsa:openidconnect.profiles:sp:sso:dod:mymoveadminfaketest")
	suite.T().Setenv("LOGIN_GOV_OFFICE_CLIENT_ID",
		"urn:gov:gsa:openidconnect.profiles:sp:sso:dod:mymoveofficefaketest")
	suite.T().Setenv("LOGIN_GOV_MY_CLIENT_ID",
		"urn:gov:gsa:openidconnect.profiles:sp:sso:dod:mymovemilfaketest")

	// get the directory of the current test file, use that to get
	// path to data file
	_, filename, _, ok := runtime.Caller(0)
	suite.Require().True(ok)
	dirname := filepath.Dir(filename)
	data, err := os.ReadFile(filepath.Join(dirname, "../../config/tls/devlocal-mtls.key"))
	suite.Require().NoError(err)
	suite.T().Setenv("LOGIN_GOV_SECRET_KEY", string(data))

	suite.Setup(InitAuthFlags, []string{})
	suite.NoError(CheckAuth(suite.viper))
}
