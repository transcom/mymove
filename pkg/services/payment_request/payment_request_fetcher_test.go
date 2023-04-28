package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequest() {
	suite.Run("If a payment request is fetched, it should be returned", func() {

		fetcher := NewPaymentRequestFetcher()

		pr := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), pr.ID)

		suite.NoError(err)
		suite.Equal(pr.ID, paymentRequest.ID)
	})

	suite.Run("returns payment request with proof of service docs", func() {

		fetcher := NewPaymentRequestFetcher()

		primeUpload := factory.BuildPrimeUpload(suite.DB(), nil, nil)
		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), primeUpload.ProofOfServiceDoc.PaymentRequestID)

		suite.NoError(err)
		suite.Equal(primeUpload.ProofOfServiceDoc.PaymentRequest.ID, paymentRequest.ID)

		suite.Len(paymentRequest.ProofOfServiceDocs, 1)
		suite.Equal(primeUpload.ProofOfServiceDoc.ID, paymentRequest.ProofOfServiceDocs[0].ID)

		suite.Len(paymentRequest.ProofOfServiceDocs[0].PrimeUploads, 1)
		suite.Equal(primeUpload.ID, paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].ID)

		suite.Equal(primeUpload.UploadID, paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].UploadID)
	})

	suite.Run("returns payment request without soft deleted proof of service docs", func() {

		fetcher := NewPaymentRequestFetcher()
		primeUpload := factory.BuildPrimeUpload(suite.DB(), nil, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), primeUpload.ProofOfServiceDoc.PaymentRequest.ID)

		suite.NoError(err)
		suite.Equal(primeUpload.ProofOfServiceDoc.PaymentRequest.ID, paymentRequest.ID)

		suite.Len(paymentRequest.ProofOfServiceDocs, 1)
		suite.Equal(primeUpload.ProofOfServiceDoc.ID, paymentRequest.ProofOfServiceDocs[0].ID)

		suite.Len(paymentRequest.ProofOfServiceDocs[0].PrimeUploads, 0)
	})

	suite.Run("if there is an error, we get it with zero payment request", func() {
		fetcher := NewPaymentRequestFetcher()

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), uuid.Nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(models.PaymentRequest{}, paymentRequest)
	})
}
