package cli

import (
	"os"
)

func (suite *cliTestSuite) TestConfigGEXSFTP() {
	suite.Setup(InitGEXSFTPFlags, []string{})

	err := os.Setenv("GEX_SFTP_PORT", "1234")
	suite.FatalNoError(err)
	err = os.Setenv("GEX_SFTP_IP_ADDRESS", "127.0.0.1")
	suite.FatalNoError(err)
	// generated fake host key to pass parser used following command and only saved the pub key
	//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
	err = os.Setenv("GEX_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")
	suite.FatalNoError(err)

	suite.NoError(CheckGEXSFTP(suite.viper))
}
