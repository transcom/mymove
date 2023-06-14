// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values are used to set/unset environment variables needed for session creation in the unit test's local database
// RA: Setting/unsetting of environment variables does not present any risks and are solely used for unit testing purposes
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package invoice

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SyncadaSftpSenderSuite struct {
	*testingsuite.PopTestSuite
}

func TestSyncadaSftpSenderSuite(t *testing.T) {

	ts := &SyncadaSftpSenderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("syncada_sftp_sender"),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *SyncadaSftpSenderSuite) TestSendToSyncadaSftp() {
	suite.Run("constructor fails if the GEX_SFTP_USER_ID is missing", func() {
		os.Unsetenv("GEX_SFTP_USER_ID")
		sender, err := InitNewSyncadaSFTPSession()
		suite.Error(err)
		suite.Nil(sender)
		suite.Equal("Invalid credentials sftp missing GEX_SFTP_USER_ID", err.Error())
	})

	suite.Run("constructor doesn't fail if passed in all env", func() {
		os.Setenv("GEX_SFTP_USER_ID", "FAKE_USER_ID")
		sender, err := InitNewSyncadaSFTPSession()
		suite.NoError(err)
		suite.NotNil(sender)
	})
}
