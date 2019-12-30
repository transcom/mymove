package paymentrequest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/storage/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaymentRequestServiceSuite) TestCreateUpload() {
	storer := &mocks.FileStorer{}
	storer.On("Store",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*os.File"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*string"),
	).Return(&storage.StoreResult{}, nil).Once()

	activeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{}) // temp user-- will need to be connected to prime
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	file, _ := os.Open("../../uploader/testdata/test.pdf")

	uploaderFile := uploader.File{
		File: file,
	}

	suite.T().Run("Upload is created successfully", func(t *testing.T) {
		uploadCreator := NewUploadCreator(suite.DB(), suite.logger, storer)
		upload, err := uploadCreator.CreateUpload(uploaderFile, paymentRequest.ID, *activeUser.UserID)

		suite.NoError(err)
		suite.Equal("test.pdf", upload.Filename)
	})
}
