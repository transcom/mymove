package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequestBulkDownload() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)

	generator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.FatalNil(err)

	primeUpload := factory.BuildPrimeUpload(suite.DB(), nil, nil)
	suite.FatalNil(err)
	if generator != nil {
		suite.FatalNil(err)
	}

	paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ProofOfServiceDocs: models.ProofOfServiceDocs{
					primeUpload.ProofOfServiceDoc,
				},
			},
		},
	}, nil)

	suite.NotNil(paymentRequest)
}
