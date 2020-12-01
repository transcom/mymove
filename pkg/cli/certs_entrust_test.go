package cli

import "os"

func (suite *cliTestSuite) TestEntrustCertificates() {

	if os.Getenv("TEST_ACC_DOD_CERTIFICATES") != "1" {
		suite.logger.Info("skipping TestEntrustCertificates")
		return
	}

	suite.Setup(InitEntrustCertFlags, []string{})
	suite.NoError(CheckEntrustCert(suite.viper))
}
