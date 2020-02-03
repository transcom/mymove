package uploader_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/storage/mocks"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/uploader"
)

type UploaderSuite struct {
	testingsuite.PopTestSuite
	logger       uploader.Logger
	storer       storage.FileStorer
	filesToClose []afero.File
	fs           *afero.Afero
}

func (suite *UploaderSuite) SetupTest() {
	var fs = afero.NewMemMapFs()
	suite.fs = &afero.Afero{Fs: fs}
}

func (suite *UploaderSuite) openLocalFile(path string) (afero.File, error) {
	file, err := os.Open(path)
	if err != nil {
		suite.logger.Fatal("Error opening local file", zap.Error(err))
	}

	outputFile, err := suite.fs.Create(path)
	if err != nil {
		suite.logger.Fatal("Error creating afero file", zap.Error(err))
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		suite.logger.Fatal("Error copying to afero file", zap.Error(err))
	}

	return outputFile, nil
}

func (suite *UploaderSuite) fixture(name string) afero.File {
	fixtureDir := "testdatagen/testdata"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatalf("failed to get current directory: %s", err)
	}

	fixturePath := path.Join(cwd, "..", fixtureDir, name)
	file, err := suite.openLocalFile(fixturePath)

	if err != nil {
		suite.T().Fatalf("failed to create a fixture file: %s", err)
	}
	suite.closeFile(file)
	return file
}

func (suite *UploaderSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Close()
	}
}

func (suite *UploaderSuite) closeFile(file afero.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

func TestUploaderSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &UploaderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
		storer:       storageTest.NewFakeS3Storage(true),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *UploaderSuite) TestUploaderExceedsFileSizeLimit() {
	_, err := uploader.NewUploader(suite.DB(), suite.logger, suite.storer, 251*uploader.MB)
	suite.Error(err)
	suite.Equal(uploader.ErrFileSizeLimitExceedsMax, err)
}

func (suite *UploaderSuite) TestUploadFromLocalFile() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	up, err := uploader.NewUploader(suite.DB(), suite.logger, suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")

	upload, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(upload.ContentType, "application/pdf")
	suite.Equal(upload.Checksum, "nOE6HwzyE4VEDXn67ULeeA==")
}

func (suite *UploaderSuite) TestUploadFromLocalFileZeroLength() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	up, err := uploader.NewUploader(suite.DB(), suite.logger, suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(0 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	upload, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.Equal(uploader.ErrZeroLengthFile, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(upload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestUploadFromLocalFileWrongContentType() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	up, err := uploader.NewUploader(suite.DB(), suite.logger, suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(1 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	upload, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.NoError(err)
	suite.True(verrs.HasAny(), "invalid content type for upload")
	suite.Nil(upload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestTooLargeUploadFromLocalFile() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	up, err := uploader.NewUploader(suite.DB(), suite.logger, suite.storer, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(26 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	_, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, uploader.File{File: f}, uploader.AllowedTypesAny)
	suite.Error(err)
	suite.IsType(uploader.ErrTooLarge{}, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestStorerCalledWithTags() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	fakeS3 := &mocks.FileStorer{}
	up, err := uploader.NewUploader(suite.DB(), suite.logger, fakeS3, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(5 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	tags := "metaDataTag=value"
	fakeS3.On("Store",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		&tags).Return(&storage.StoreResult{}, nil)
	// assert tags are passed along to storer
	_, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, uploader.File{File: f, Tags: &tags}, uploader.AllowedTypesAny)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) createFileOfArbitrarySize(size uint64) (afero.File, func(), error) {
	data := make([]byte, size, size)
	tmpFileName := "tmpfile"
	f, err := suite.fs.Create(tmpFileName)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	_, err = f.Write(data)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	cleanup := func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Println("error closing file")
		}
		if removeErr := suite.fs.Remove(tmpFileName); removeErr != nil {
			log.Println("error removing file")
		}
	}
	return f, cleanup, err
}

func (suite *UploaderSuite) helperNewTempFile() (afero.File, error) {
	outputFile, err := suite.fs.TempFile("/tmp/milmoves/", "TestCreateUploadNoDocument")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return outputFile, nil
}

func (suite *UploaderSuite) TestCreateUploadNoDocument() {
	document := testdatagen.MakeDefaultDocument(suite.DB())
	userID := document.ServiceMember.UserID

	up, err := uploader.NewUploader(suite.DB(), suite.logger, suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")
	fixtureFileInfo, err := file.Stat()
	suite.NoError(err)

	// Create file and upload
	upload, verrs, err := up.CreateUpload(userID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.Empty(verrs.Error(), "verrs returned error")
	suite.NotNil(upload, "failed to create upload structure")
	file.Close()

	// Download file and test size
	download, err := up.Download(upload)
	suite.NoError(err)
	defer download.Close()

	outputFile, err := suite.helperNewTempFile()
	suite.NoError(err)
	defer outputFile.Close()

	written, err := io.Copy(outputFile, download)
	suite.NoError(err)
	suite.NotEqual(0, written)

	info, err := outputFile.Stat()
	suite.Equal(fixtureFileInfo.Size(), info.Size())
	suite.NoError(err)

	// Delete file previously uploaded
	err = up.Storer.Delete(upload.StorageKey)
	suite.NoError(err)
}
