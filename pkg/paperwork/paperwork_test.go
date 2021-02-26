package paperwork

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/uploader"
)

type PaperworkSuite struct {
	testingsuite.PopTestSuite
	logger       Logger
	userUploader *uploader.UserUploader
	filesToClose []afero.File
}

func (suite *PaperworkSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Close()
	}
}

func (suite *PaperworkSuite) closeFile(file afero.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

func (suite *PaperworkSuite) openLocalFile(path string, fs *afero.Afero) (afero.File, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	outputFile, err := fs.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating afero file")
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		return nil, errors.Wrap(err, "error copying over file contents")
	}

	suite.closeFile(outputFile)

	return outputFile, nil
}

func TestPaperworkSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	storer := storageTest.NewFakeS3Storage(true)

	popSuite := testingsuite.NewPopTestSuite(testingsuite.CurrentPackage())
	newUploader, err := uploader.NewUserUploader(popSuite.DB(), logger, storer, 25*uploader.MB)
	if err != nil {
		log.Panic(err)
	}
	hs := &PaperworkSuite{
		PopTestSuite: popSuite,
		logger:       logger,
		userUploader: newUploader,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
