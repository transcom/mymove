package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequestBulkDownload() {
	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)

	generator, err := paperworkgenerator.NewGenerator(primeUploader.Uploader())
	suite.FatalNil(err)

	paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
	primeUpload := factory.BuildPrimeUpload(suite.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)
	posd := factory.BuildProofOfServiceDoc(suite.DB(), []factory.Customization{
		{
			Model:    primeUpload,
			LinkOnly: true,
		},
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)

	paymentRequest.ProofOfServiceDocs = append(paymentRequest.ProofOfServiceDocs, posd)

	creator := &paymentRequestBulkDownloadCreator{
		pdfGenerator: generator,
	}

	bulkDownload, err := creator.CreatePaymentRequestBulkDownload(suite.AppContextForTest(), paymentRequest.ID)
	suite.NoError(err)
	suite.NotNil(bulkDownload)
}
