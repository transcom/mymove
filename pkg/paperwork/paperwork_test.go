package paperwork

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

type PaperworkSuite struct {
	suite.Suite
	db           *pop.Connection
	logger       *zap.Logger
	uploader     *uploader.Uploader
	filesToClose []afero.File
}

func (suite *PaperworkSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *PaperworkSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
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
	file, err := os.Open(path)
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

func (suite *PaperworkSuite) FatalNil(err error, messages ...string) {
	t := suite.T()
	t.Helper()
	if err != nil {
		if len(messages) > 0 {
			t.Fatalf("%s: %s", strings.Join(messages, ","), err.Error())
		} else {
			t.Fatal(err.Error())
		}
	}
}

func TestPaperworkSuite(t *testing.T) {
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
	storer := storageTest.NewFakeS3Storage(true)

	hs := &PaperworkSuite{
		db:       db,
		logger:   logger,
		uploader: uploader.NewUploader(db, logger, storer),
	}

	suite.Run(t, hs)
}
