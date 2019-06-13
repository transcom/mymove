package cli

import (
	"os"
)

func (suite *cliTestSuite) TestDODCertificates() {

	if os.Getenv("TEST_ACC_DOD_CERTIFICATES") != "1" {
		suite.logger.Info("skipping TestDODCertificates")
		return
	}

	suite.Setup(InitCertFlags, []string{})
	suite.NoError(CheckCert(suite.viper))
}
