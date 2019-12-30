package paymentrequest

import (
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/storage/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaymentRequestServiceSuite) openLocalFile(path string) (afero.File, error) {
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

func (suite *PaymentRequestServiceSuite) TestCreateUploadSuccess() {
	storer := &mocks.FileStorer{}
	storer.On("Store",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*mem.File"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*string"),
	).Return(&storage.StoreResult{}, nil).Once()

	activeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{}) // temp user-- will need to be connected to prime
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	file, err := suite.openLocalFile("../../uploader/testdata/test.pdf")
	suite.NoError(err)

	uploaderFile := uploader.File{
		File: file,
	}

	suite.T().Run("Upload is created successfully", func(t *testing.T) {
		uploadCreator := NewUploadCreator(suite.DB(), suite.logger, storer)
		upload, err := uploadCreator.CreateUpload(uploaderFile, paymentRequest.ID, *activeUser.UserID)

		suite.NoError(err)
		suite.Equal("../../uploader/testdata/test.pdf", upload.Filename)
		suite.Equal(int64(10596), upload.Bytes)
		suite.Equal("application/pdf", upload.ContentType)

	})
}

func (suite *PaymentRequestServiceSuite) TestCreateUploadFailure() {
	storer := &mocks.FileStorer{}
	storer.On("Store",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*mem.File"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*string"),
	).Return(&storage.StoreResult{}, nil).Once()

	//activeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{}) // temp user-- will need to be connected to prime
	//paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	//file, err := suite.openLocalFile("../../uploader/testdata/test.pdf")
	//suite.NoError(err)
	//
	//uploaderFile := uploader.File{
	//	File: file,
	//}

	suite.T().Run("invalid payment request ID", func(t *testing.T) {

	})

	suite.T().Run("invalid user ID", func(t *testing.T) {

	})

	suite.T().Run("", func(t *testing.T) {

	})
}
