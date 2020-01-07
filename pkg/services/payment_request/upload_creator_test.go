package paymentrequest

import (
	"os"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/storage/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

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

	testFile, err := os.Open("../../uploader/testdata/test.pdf")
	suite.NoError(err)

	suite.T().Run("Upload is created successfully", func(t *testing.T) {
		uploadCreator := NewPaymentRequestUploadCreator(suite.DB(), suite.logger, storer)
		upload, err := uploadCreator.CreateUpload(testFile, paymentRequest.ID, *activeUser.UserID)

		suite.NoError(err)
		suite.Equal("/app/payment-request-uploads/tempfile", upload.Filename)
		suite.Equal(int64(10596), upload.Bytes)
		suite.Equal("application/pdf", upload.ContentType)
	})

	testFile.Close()
}

func (suite *PaymentRequestServiceSuite) TestCreateUploadFailure() {
	storer := &mocks.FileStorer{}
	storer.On("Store",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*mem.File"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*string"),
	).Return(&storage.StoreResult{}, nil).Once()
	activeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{}) // temp user-- will need to be connected to prime
	testdatagen.MakeDefaultPaymentRequest(suite.DB())

	testFile, err := os.Open("../../uploader/testdata/test.pdf")
	suite.NoError(err)

	suite.T().Run("invalid payment request ID", func(t *testing.T) {
		uploadCreator := NewPaymentRequestUploadCreator(suite.DB(), suite.logger, storer)
		_, err := uploadCreator.CreateUpload(testFile, uuid.FromStringOrNil("96b77644-4028-48c2-9ab8-754f33309db9"), *activeUser.UserID)
		suite.Error(err)
	})
	suite.T().Run("invalid user ID", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		uploadCreator := NewPaymentRequestUploadCreator(suite.DB(), suite.logger, storer)
		_, err := uploadCreator.CreateUpload(testFile, paymentRequest.ID, uuid.FromStringOrNil("806e2f96-f9f9-4cbb-9a3d-d2f488539a1f"))
		suite.Error(err)
	})

	testFile.Close()
}
