package invoice

import (
	"os"
	"testing"
	"time"

	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SyncadaSftpReaderSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestSyncadaSftpReaderSuite(t *testing.T) {

	ts := &SyncadaSftpReaderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("syncada_sftp_reader")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

// TODO I am sure there is a better way to construct an os.FileInfo for tests
type FakeFileInfo struct {
	name    string
	modTime time.Time
}

func (f FakeFileInfo) Name() string {
	return f.name
}
func (f FakeFileInfo) Size() int64 {
	return 1024
}
func (f FakeFileInfo) Mode() os.FileMode {
	return 0
}
func (f FakeFileInfo) ModTime() time.Time {
	return f.modTime
}
func (f FakeFileInfo) IsDir() bool {
	return false
}
func (f FakeFileInfo) Sys() interface{} {
	return nil
}

func (suite *SyncadaSftpReaderSuite) TestReadToSyncadaSftp() {
	ediFileInfo := FakeFileInfo{"edifile", time.Now()}
	singleFileInfo := make([]os.FileInfo, 1)
	singleFileInfo[0] = ediFileInfo
	pickupDir := "/foo"
	//ediFilePath := "/foo/edifile"

	suite.T().Run("Nothing should be processed or read from an empty directory", func(t *testing.T) {
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(make([]os.FileInfo, 0), nil)

		processor := &mocks.SyncadaFileProcessor{}
		processor.On("ProcessFile", mock.Anything).Return(nil)

		session := InitNewSyncadaSFTPReaderSession(client, suite.logger)
		_, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
		suite.NoError(err)
		client.AssertCalled(t, "ReadDir", pickupDir)
		processor.AssertNotCalled(t, "ProcessFile", mock.Anything)
		client.AssertNotCalled(t, "Remove", mock.Anything)
	})
	suite.T().Run("File open error should prevent processing and deletion", func(t *testing.T) {
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(singleFileInfo, nil)

		// TODO maybe i should mock SFTPClient.downloadFile instead?
		client.On("Open", mock.Anything).Return(nil, errors.New("ERROR"))

		processor := &mocks.SyncadaFileProcessor{}
		processor.On("ProcessFile", mock.Anything).Return(nil)

		session := InitNewSyncadaSFTPReaderSession(client, suite.logger)
		session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
		client.AssertCalled(t, "ReadDir", pickupDir)
		processor.AssertNotCalled(t, "ProcessFile", mock.Anything)
		client.AssertNotCalled(t, "Remove", mock.Anything)
	})
	//suite.T().Run("File read successfully should be processed and deleted", func(t *testing.T) {
	//	client := &mocks.SFTPClient{}
	//	client.On("ReadDir", mock.Anything).Return(singleFileInfo, nil)
	//
	//	// TODO I've tried a few ways to mock this and have not been successful
	//	client.On("Open", mock.AnythingOfType("string")).Return(bytes.NewBufferString("datadatadata"))
	//
	//	processor := &mocks.SyncadaFileProcessor{}
	//	processor.On("ProcessFile", mock.Anything).Return(nil)
	//
	//	session := InitNewSyncadaSFTPReaderSession(client, suite.logger)
	//	session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
	//	client.AssertCalled(t, "ReadDir", pickupDir)
	//	processor.AssertCalled(t, "ProcessFile", ediFilePath, "datadatadata")
	//	client.AssertCalled(t, "Remove", ediFilePath)
	//})

	// test that we request files from the right directory?
	// test that we don't crash when reading files we dont have permission to read?
	// test skipping older files
	// test that we call process function on files

	// test file reading errors not crashing us
	// make sure we're calling delete on successfully processed files ONLY
}
