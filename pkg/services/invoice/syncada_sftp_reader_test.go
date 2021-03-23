package invoice

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

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

type FakeFile struct {
	contents string
}

func (f FakeFile) Close() error {
	return nil
}

func (f FakeFile) WriteTo(w io.Writer) (int64, error) {
	written, err := w.Write([]byte(f.contents))
	return int64(written), err
}

type FileTestData struct {
	fileInfo FakeFileInfo
	file     FakeFile
	path     string
}

func (suite *SyncadaSftpReaderSuite) TestReadToSyncadaSftp() {
	pickupDir := "/foo"
	ediFileInfo := FakeFileInfo{"edifile", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	singleFileInfo := make([]os.FileInfo, 1)
	singleFileInfo[0] = ediFileInfo

	multipleFileTestData := []FileTestData{
		{
			fileInfo: FakeFileInfo{"file0", time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)},
			file:     FakeFile{"contents0"},
			path:     pickupDir + "/" + "file0",
		},
		{
			fileInfo: FakeFileInfo{"file1", time.Date(2020, 1, 1, 1, 1, 0, 0, time.UTC)},
			file:     FakeFile{"contents1"},
			path:     pickupDir + "/" + "file1",
		},
		{
			fileInfo: FakeFileInfo{"file2", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC)},
			file:     FakeFile{"contents2"},
			path:     pickupDir + "/" + "file2",
		},
	}
	infoForMultipleFiles := make([]os.FileInfo, len(multipleFileTestData))
	for i, f := range multipleFileTestData {
		infoForMultipleFiles[i] = f.fileInfo
	}

	suite.T().Run("Nothing should be processed or read from an empty directory", func(t *testing.T) {
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(make([]os.FileInfo, 0), nil)

		processor := &mocks.SyncadaFileProcessor{}

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		_, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
		suite.NoError(err)
		client.AssertCalled(t, "ReadDir", pickupDir)
		processor.AssertNotCalled(t, "ProcessFile", mock.Anything)
		client.AssertNotCalled(t, "Remove", mock.Anything)
	})

	suite.T().Run("ReadDir error should result in error", func(t *testing.T) {
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(nil, errors.New("ERROR"))
		processor := &mocks.SyncadaFileProcessor{}
		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		_, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
		suite.Error(err)
	})

	suite.T().Run("File open error should prevent processing and deletion", func(t *testing.T) {
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(singleFileInfo, nil)

		client.On("Open", mock.Anything).Return(nil, errors.New("ERROR"))

		processor := &mocks.SyncadaFileProcessor{}
		processor.On("ProcessFile", mock.Anything).Return(nil)

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		_, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
		suite.NoError(err)

		client.AssertCalled(t, "ReadDir", pickupDir)
		processor.AssertNotCalled(t, "ProcessFile", mock.Anything)
		client.AssertNotCalled(t, "Remove", mock.Anything)
	})

	suite.T().Run("File WriteTo error should prevent processing and deletion", func(t *testing.T) {
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(singleFileInfo, nil)

		fileThatWillReturnErrOnWrite := &mocks.SFTPFiler{}
		fileThatWillReturnErrOnWrite.On("WriteTo", mock.Anything).Return(int64(0), errors.New("ERROR"))
		fileThatWillReturnErrOnWrite.On("Close", mock.Anything).Return(nil)

		client.On("Open", mock.Anything).Return(fileThatWillReturnErrOnWrite, nil)

		processor := &mocks.SyncadaFileProcessor{}
		processor.On("ProcessFile", mock.Anything).Return(nil)

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		_, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)
		suite.NoError(err)

		client.AssertCalled(t, "ReadDir", pickupDir)
		client.AssertCalled(t, "Open", mock.Anything)
		processor.AssertNotCalled(t, "ProcessFile", mock.Anything)
		client.AssertNotCalled(t, "Remove", mock.Anything)
	})

	suite.T().Run("If ProcessFile returns error, we should not remove the file", func(t *testing.T) {
		// set up mocks
		client := &mocks.SFTPClient{}
		client.On("ReadDir", mock.Anything).Return(infoForMultipleFiles, nil)
		for _, data := range multipleFileTestData {
			client.On("Open", data.path).Return(data.file, nil)
			client.On("Remove", data.path).Return(nil)
		}
		processor := &mocks.SyncadaFileProcessor{}

		// Mock processing error for one of the files
		processor.On("ProcessFile", multipleFileTestData[1].path, multipleFileTestData[1].file.contents).Return(errors.New("ERROR"))
		// No error for the rest of the files
		processor.On("ProcessFile", mock.Anything, mock.Anything).Return(nil)

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		modTime, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)

		suite.NoError(err)
		suite.Equal(multipleFileTestData[len(multipleFileTestData)-1].fileInfo.ModTime(), modTime)

		// Make sure we called external methods with the right args for every file
		client.AssertCalled(t, "ReadDir", pickupDir)
		for _, data := range multipleFileTestData {
			client.AssertCalled(t, "Open", data.path)
			processor.AssertCalled(t, "ProcessFile", data.path, data.file.contents)
		}
		client.AssertNotCalled(t, "Remove", multipleFileTestData[1].path)
	})

	suite.T().Run("Files read successfully should be processed and deleted", func(t *testing.T) {
		// set up mocks
		client := &mocks.SFTPClient{}
		processor := &mocks.SyncadaFileProcessor{}
		client.On("ReadDir", mock.Anything).Return(infoForMultipleFiles, nil)
		for _, data := range multipleFileTestData {
			client.On("Open", data.path).Return(data.file, nil)
			processor.On("ProcessFile", data.path, data.file.contents).Return(nil)
			client.On("Remove", data.path).Return(nil)
		}

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		modTime, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)

		suite.NoError(err)
		suite.Equal(multipleFileTestData[len(multipleFileTestData)-1].fileInfo.ModTime(), modTime)

		// Make sure we called external methods with the right args for every file
		client.AssertCalled(t, "ReadDir", pickupDir)
		for _, data := range multipleFileTestData {
			client.AssertCalled(t, "Open", data.path)
			processor.AssertCalled(t, "ProcessFile", data.path, data.file.contents)
			client.AssertCalled(t, "Remove", data.path)
		}
	})

	suite.T().Run("Files should not be deleted when deletion flag is not set", func(t *testing.T) {
		// set up mocks
		client := &mocks.SFTPClient{}
		processor := &mocks.SyncadaFileProcessor{}
		client.On("ReadDir", mock.Anything).Return(infoForMultipleFiles, nil)
		for _, data := range multipleFileTestData {
			client.On("Open", data.path).Return(data.file, nil)
			processor.On("ProcessFile", data.path, data.file.contents).Return(nil)
			client.On("Remove", data.path).Return(nil)
		}

		session := NewSyncadaSFTPReaderSession(client, suite.logger, false)
		modTime, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)

		suite.NoError(err)
		suite.Equal(multipleFileTestData[len(multipleFileTestData)-1].fileInfo.ModTime(), modTime)

		// Make sure we open and process all files, and do not delete any of them
		client.AssertCalled(t, "ReadDir", pickupDir)
		for _, data := range multipleFileTestData {
			client.AssertCalled(t, "Open", data.path)
			processor.AssertCalled(t, "ProcessFile", data.path, data.file.contents)
		}
		client.AssertNotCalled(t, "Remove", mock.Anything)
	})

	suite.T().Run("SFTP Remove errors don't cause an error", func(t *testing.T) {
		// set up mocks
		client := &mocks.SFTPClient{}
		processor := &mocks.SyncadaFileProcessor{}
		client.On("ReadDir", mock.Anything).Return(infoForMultipleFiles, nil)
		client.On("Remove", mock.Anything).Return(errors.New("ERROR"))
		for _, data := range multipleFileTestData {
			client.On("Open", data.path).Return(data.file, nil)
			processor.On("ProcessFile", data.path, data.file.contents).Return(nil)
		}

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)
		modTime, err := session.FetchAndProcessSyncadaFiles(pickupDir, time.Time{}, processor)

		suite.NoError(err)
		suite.Equal(multipleFileTestData[len(multipleFileTestData)-1].fileInfo.ModTime(), modTime)

		// Make sure we called external methods with the right args for every file
		for _, data := range multipleFileTestData {
			client.AssertCalled(t, "Remove", data.path)
		}
	})

	suite.T().Run("Files before cutoff time should be skipped", func(t *testing.T) {
		// set up mocks
		client := &mocks.SFTPClient{}
		processor := &mocks.SyncadaFileProcessor{}
		client.On("ReadDir", mock.Anything).Return(infoForMultipleFiles, nil)

		// only need to mock calls for the one file that will be processed
		client.On("Open", multipleFileTestData[2].path).Return(multipleFileTestData[2].file, nil)
		processor.On("ProcessFile", multipleFileTestData[2].path, multipleFileTestData[2].file.contents).Return(nil)
		client.On("Remove", multipleFileTestData[2].path).Return(nil)

		session := NewSyncadaSFTPReaderSession(client, suite.logger, true)

		// We're using the modified time for the second file as the last read time.
		// The files are sorted by modTime, so we should skip the first two files and only process the third
		modTime, err := session.FetchAndProcessSyncadaFiles(pickupDir, multipleFileTestData[1].fileInfo.modTime, processor)

		suite.NoError(err)
		suite.Equal(multipleFileTestData[2].fileInfo.ModTime(), modTime)

		// Files at or before the cutoff time should be skipped
		for _, data := range multipleFileTestData[:2] {
			client.AssertNotCalled(t, "Open", data.path)
		}
		// Last file should be opened
		client.AssertCalled(t, "Open", multipleFileTestData[2].path)
	})
}
