package uploader_test

import (
	"io"
	"log"
	"os"
	"path"
	"testing"

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
	suite.DB().TruncateAll()
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
	fixtureDir := "testdata"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatalf("failed to get current directory: %s", err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)
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
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
		storer:       storageTest.NewFakeS3Storage(true),
	}

	suite.Run(t, hs)
}

func (suite *UploaderSuite) TestUploadFromLocalFile() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	up := uploader.NewUploader(suite.DB(), suite.logger, suite.storer)
	file := suite.fixture("test.pdf")

	upload, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, file, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(upload.ContentType, "application/pdf")
	suite.Equal(upload.Checksum, "nOE6HwzyE4VEDXn67ULeeA==")
}

func (suite *UploaderSuite) TestUploadFromLocalFileZeroLength() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	up := uploader.NewUploader(suite.DB(), suite.logger, suite.storer)
	file := suite.fixture("empty.pdf")

	upload, verrs, err := up.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, file, uploader.AllowedTypesPDF)
	suite.Equal(err, uploader.ErrZeroLengthFile)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(upload, "returned an upload when erroring")
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

	up := uploader.NewUploader(suite.DB(), suite.logger, suite.storer)
	file := suite.fixture("test.pdf")
	fixtureFileInfo, err := file.Stat()
	suite.Nil(err)

	// Create file and upload
	upload, verrs, err := up.CreateUpload(userID, &file, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.Empty(verrs.Error(), "verrs returned error")
	suite.NotNil(upload, "failed to create upload structure")
	file.Close()

	// Download file and test size
	download, err := up.Download(upload)
	suite.Nil(err)
	defer download.Close()

	outputFile, err := suite.helperNewTempFile()
	suite.Nil(err)
	defer outputFile.Close()

	written, err := io.Copy(outputFile, download)
	suite.Nil(err)
	suite.NotEqual(0, written)

	info, err := outputFile.Stat()
	suite.Equal(fixtureFileInfo.Size(), info.Size())

	// Delete file previously uploaded
	err = up.Storer.Delete(upload.StorageKey)
	suite.Nil(err)
}
