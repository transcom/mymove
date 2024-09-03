package paymentrequest

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// PaymentRequestServiceSuite is a suite for testing payment requests
type PaymentRequestServiceSuite struct {
	*testingsuite.PopTestSuite
	fs     *afero.Afero
	storer storage.FileStorer
}

func TestPaymentRequestServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &PaymentRequestServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		fs:           file,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *PaymentRequestServiceSuite) openLocalFile(path string) (afero.File, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		suite.Logger().Fatal("Error opening local file", zap.Error(err))
	}

	outputFile, err := suite.fs.Create(path)
	if err != nil {
		suite.Logger().Fatal("Error creating afero file", zap.Error(err))
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		suite.Logger().Fatal("Error copying to afero file", zap.Error(err))
	}

	return outputFile, nil
}
