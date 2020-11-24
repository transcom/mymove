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
			},
		},
	}

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
		sender, err := InitNewSyncadaSFTPSession()
		suite.NoError(err)
		suite.NotNil(sender)
	})
}
