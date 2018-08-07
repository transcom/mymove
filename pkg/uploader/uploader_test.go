package uploader

import (
	"io"
	"log"
	"os"
	"path"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type UploaderSuite struct {
	suite.Suite
	db           *pop.Connection
	logger       *zap.Logger
	storer       storage.FileStorer
	filesToClose []afero.File
	fs           *afero.Afero
}

func (suite *UploaderSuite) SetupTest() {
	var fs = afero.NewMemMapFs()
	suite.fs = &afero.Afero{Fs: fs}
	suite.db.TruncateAll()
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
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &UploaderSuite{
		db:     db,
		logger: logger,
		storer: storageTest.NewFakeS3Storage(true),
	}

	suite.Run(t, hs)
}

func (suite *UploaderSuite) TestUploadFromLocalFile() {
	document := testdatagen.MakeDefaultDocument(suite.db)

	up := NewUploader(suite.db, suite.logger, suite.storer)
	file := suite.fixture("test.pdf")

	upload, verrs, err := up.CreateUpload(&document.ID, document.ServiceMember.UserID, file)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(upload.ContentType, "application/pdf")
	suite.Equal(upload.Checksum, "nOE6HwzyE4VEDXn67ULeeA==")
}

func (suite *UploaderSuite) TestUploadFromLocalFileZeroLength() {
	document := testdatagen.MakeDefaultDocument(suite.db)

	up := NewUploader(suite.db, suite.logger, suite.storer)
	file := suite.fixture("empty.pdf")

	upload, verrs, err := up.CreateUpload(&document.ID, document.ServiceMember.UserID, file)
	suite.Equal(err, ErrZeroLengthFile)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(upload, "returned an upload when erroring")
}
