//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values are used to set/unset environment variables needed for session creation in the unit test's local database
//RA: Setting/unsetting of environment variables does not present any risks and are solely used for unit testing purposes
//RA Developer Status: Mitigated
//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
//RA Validator: jneuner@mitre.org
//RA Modified Severity:
// nolint:errcheck
package invoice

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SyncadaSftpSenderSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestSyncadaSftpSenderSuite(t *testing.T) {

	ts := &SyncadaSftpSenderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("syncada_sftp_sender")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *SyncadaSftpSenderSuite) TestSendToSyncadaSftp() {
	type setupEnvVars func()

	missingCreds := []struct {
		TestEnvironmentVar string
		Setup              setupEnvVars
	}{
		{
			TestEnvironmentVar: "SYNCADA_SFTP_PORT",
			Setup: func() {
				os.Unsetenv("SYNCADA_SFTP_PORT")
				os.Unsetenv("SYNCADA_SFTP_USER_ID")
				os.Unsetenv("SYNCADA_SFTP_IP_ADDRESS")
				os.Unsetenv("SYNCADA_SFTP_PASSWORD")
				os.Unsetenv("SYNCADA_SFTP_INBOUND_DIRECTORY")
				os.Unsetenv("SYNCADA_SFTP_HOST_KEY")
			},
		},
		{
			TestEnvironmentVar: "SYNCADA_SFTP_USER_ID",
			Setup: func() {
				os.Setenv("SYNCADA_SFTP_PORT", "1234")
				os.Unsetenv("SYNCADA_SFTP_USER_ID")
				os.Unsetenv("SYNCADA_SFTP_IP_ADDRESS")
				os.Unsetenv("SYNCADA_SFTP_PASSWORD")
				os.Unsetenv("SYNCADA_SFTP_INBOUND_DIRECTORY")
				os.Unsetenv("SYNCADA_SFTP_HOST_KEY")
			},
		},
		{
			TestEnvironmentVar: "SYNCADA_SFTP_IP_ADDRESS",
			Setup: func() {
				os.Setenv("SYNCADA_SFTP_PORT", "1234")
				os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
				os.Unsetenv("SYNCADA_SFTP_IP_ADDRESS")
				os.Unsetenv("SYNCADA_SFTP_PASSWORD")
				os.Unsetenv("SYNCADA_SFTP_INBOUND_DIRECTORY")
				os.Unsetenv("SYNCADA_SFTP_HOST_KEY")
			},
		},
		{
			TestEnvironmentVar: "SYNCADA_SFTP_PASSWORD",
			Setup: func() {
				os.Setenv("SYNCADA_SFTP_PORT", "1234")
				os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
				os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
				os.Unsetenv("SYNCADA_SFTP_PASSWORD")
				os.Unsetenv("SYNCADA_SFTP_INBOUND_DIRECTORY")
				os.Unsetenv("SYNCADA_SFTP_HOST_KEY")
			},
		},
		{
			TestEnvironmentVar: "SYNCADA_SFTP_INBOUND_DIRECTORY",
			Setup: func() {
				os.Setenv("SYNCADA_SFTP_PORT", "1234")
				os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
				os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
				os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
				os.Unsetenv("SYNCADA_SFTP_INBOUND_DIRECTORY")
				os.Unsetenv("SYNCADA_SFTP_HOST_KEY")
			},
		},
		{
			TestEnvironmentVar: "SYNCADA_SFTP_HOST_KEY",
			Setup: func() {
				os.Setenv("SYNCADA_SFTP_PORT", "1234")
				os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
				os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
				os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
				os.Setenv("SYNCADA_SFTP_INBOUND_DIRECTORY", "/Dropoff")
				os.Unsetenv("SYNCADA_SFTP_HOST_KEY")
			},
		},
	}

	// Test failure if any environment variable is missing
	for _, data := range missingCreds {
		suite.T().Run(fmt.Sprintf("constructor fails if %s is missing", data.TestEnvironmentVar), func(t *testing.T) {
			data.Setup()
			sender, err := InitNewSyncadaSFTPSession()
			suite.Error(err)
			suite.Nil(sender)
			suite.Equal(fmt.Sprintf("Invalid credentials sftp missing %s", data.TestEnvironmentVar), err.Error())
		})
	}

	suite.T().Run("constructor doesn't fail if passed in all env", func(t *testing.T) {
		os.Setenv("SYNCADA_SFTP_PORT", "1234")
		os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
		os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
		os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
		os.Setenv("SYNCADA_SFTP_INBOUND_DIRECTORY", "/Dropoff")
		// generated fake host key to pass parser used following command and only saved the pub key
		//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
		os.Setenv("SYNCADA_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")
		sender, err := InitNewSyncadaSFTPSession()
		suite.NoError(err)
		suite.NotNil(sender)
	})

	suite.T().Run("constructor fails with invalid host key", func(t *testing.T) {
		os.Setenv("SYNCADA_SFTP_PORT", "1234")
		os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
		os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
		os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
		os.Setenv("SYNCADA_SFTP_INBOUND_DIRECTORY", "/Dropoff")
		os.Setenv("SYNCADA_SFTP_HOST_KEY", "FAKE::HOSTKEY::INVALID")
		sender, err := InitNewSyncadaSFTPSession()
		suite.Error(err)
		suite.Nil(sender)
		suite.Equal("Failed to parse host key ssh: no key found", err.Error())
	})
}
