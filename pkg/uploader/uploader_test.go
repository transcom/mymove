package uploader

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/gobuffalo/pop"
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
	filesToClose []*runtime.File
}

func (suite *UploaderSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *UploaderSuite) fixture(name string) *runtime.File {
	fixtureDir := "testdata"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatalf("failed to get current directory: %s", err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)
	file, err := NewLocalFile(fixturePath)

	if err != nil {
		suite.T().Fatalf("failed to create a fixture file: %s", err)
	}
	suite.closeFile(file)
	return file
}

func (suite *UploaderSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Data.Close()
	}
}

func (suite *UploaderSuite) closeFile(file *runtime.File) {
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

	upload, verrs, err := up.CreateUpload(document.ID, document.ServiceMember.UserID, file)
	suite.Nil(err, "failed to create upload")
	suite.Nil(verrs, "failed to validate upload")
	suite.Equal(upload.ContentType, "application/pdf")
	suite.Equal(upload.Checksum, "nOE6HwzyE4VEDXn67ULeeA==")
}

func (suite *UploaderSuite) TestUploadFromLocalFileZeroLength() {
	document := testdatagen.MakeDefaultDocument(suite.db)

	up := NewUploader(suite.db, suite.logger, suite.storer)
	file := suite.fixture("empty.pdf")

	upload, verrs, err := up.CreateUpload(document.ID, document.ServiceMember.UserID, file)
	suite.Equal(err, ErrZeroLengthFile)
	suite.Nil(verrs, "failed to validate upload")
	suite.Nil(upload, "returned an upload when erroring")
}
