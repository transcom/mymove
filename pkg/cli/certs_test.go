package cli

import (
	"fmt"
	"os"
)

func (suite *cliTestSuite) TestDODCertificates() {

	if os.Getenv("TEST_ACC_DOD_CERTIFICATES") != "1" {
		suite.logger.Info("skipping TestDODCertificates")
		return
	}

	flagSet := []string{
		fmt.Sprintf("--%s=%s", DevlocalCAFlag, "github.com/transcom/mymove/config/tls/devlocal-ca.pem"),
	}
	suite.Setup(InitCertFlags, flagSet)
	_, _, err := InitDoDCertificates(suite.viper, suite.logger)
	suite.Nil(err)
}
