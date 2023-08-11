package cli

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func (suite *cliTestSuite) TestConfigGEXSFTP() {
	suite.Setup(InitGEXSFTPFlags, []string{})

	suite.T().Setenv("GEX_SFTP_PORT", "1234")

	suite.T().Setenv("GEX_SFTP_IP_ADDRESS", "127.0.0.1")

	suite.T().Setenv("GEX_SFTP_USER_ID", "FAKE_USER_ID")

	suite.T().Setenv("GEX_SFTP_PASSWORD", "FAKE PASSWORD")

	// Generate a key for testing
	bitSize := 4096
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	suite.T().Setenv("GEX_PRIVATE_KEY", string(keyPEM))
	suite.Require().Nil(err)

	// generated fake host key to pass parser used following command and only saved the pub key
	//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
	suite.T().Setenv("GEX_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")
	suite.Require().Nil(err)

	suite.Run("pass without errors if the environment is set up with valid values", func() {
		suite.NoError(CheckGEXSFTP(suite.viper))
	})

	suite.Run("fail with an error when given an invalid private key", func() {
		// Override what was set earlier with a fake key
		suite.T().Setenv("GEX_PRIVATE_KEY", "FAKEKEY")

		suite.Error(CheckGEXSFTP(suite.viper), "no key found")
	})

	suite.Run("fail with an error when given an invalid host key", func() {
		// Override what was set earlier with a fake key
		suite.T().Setenv("GEX_SFTP_HOST_KEY", "FAKEKEY")

		suite.Error(CheckGEXSFTP(suite.viper), "no key found")
	})
}
