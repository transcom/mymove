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

	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/uploader"
)

type PaperworkSuite struct {
	testingsuite.PopTestSuite
	userUploader *uploader.UserUploader
	filesToClose []afero.File
}

func (suite *PaperworkSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used to clean up file created for unit test
		//RA: Given the functions causing the lint errors are used to clean up local storage space after a unit test, it does not present a risk
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: N/A
		// nolint:errcheck
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
	storer := storageTest.NewFakeS3Storage(true)

	popSuite := testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction())
	newUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Panic(err)
	}
	hs := &PaperworkSuite{
		PopTestSuite: popSuite,
		userUploader: newUploader,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
