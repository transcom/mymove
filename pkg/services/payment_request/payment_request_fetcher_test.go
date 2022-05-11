package paymentrequest

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequest() {
	suite.T().Run("If a payment request is fetched, it should be returned", func(t *testing.T) {

		fetcher := NewPaymentRequestFetcher()

		pr := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), pr.ID)

		suite.NoError(err)
		suite.Equal(pr.ID, paymentRequest.ID)
	})

	suite.T().Run("returns payment request with proof of service docs", func(t *testing.T) {

		fetcher := NewPaymentRequestFetcher()

		pr := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		posd := testdatagen.MakeProofOfServiceDoc(suite.DB(), testdatagen.Assertions{
			ProofOfServiceDoc: models.ProofOfServiceDoc{
				PaymentRequestID: pr.ID,
			},
		})
		u := testdatagen.MakeDefaultUpload(suite.DB())
		pu := testdatagen.MakePrimeUpload(suite.DB(), testdatagen.Assertions{
			PrimeUpload: models.PrimeUpload{
				ProofOfServiceDocID: posd.ID,
				UploadID:            u.ID,
			},
		})

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), pr.ID)

		suite.NoError(err)
		suite.Equal(pr.ID, paymentRequest.ID)

		suite.Len(paymentRequest.ProofOfServiceDocs, 1)
		suite.Equal(posd.ID, paymentRequest.ProofOfServiceDocs[0].ID)

		suite.Len(paymentRequest.ProofOfServiceDocs[0].PrimeUploads, 1)
		suite.Equal(pu.ID, paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].ID)

		suite.Equal(u.ID, paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].UploadID)
	})

	suite.T().Run("returns payment request without soft deleted proof of service docs", func(t *testing.T) {

		fetcher := NewPaymentRequestFetcher()

		pr := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		posd := testdatagen.MakeProofOfServiceDoc(suite.DB(), testdatagen.Assertions{
			ProofOfServiceDoc: models.ProofOfServiceDoc{
				PaymentRequestID: pr.ID,
			},
		})
		deletedAt := time.Now()
		u := testdatagen.MakeUpload(suite.DB(), testdatagen.Assertions{
			Upload: models.Upload{DeletedAt: &deletedAt},
		})
		testdatagen.MakePrimeUpload(suite.DB(), testdatagen.Assertions{
			PrimeUpload: models.PrimeUpload{
				ProofOfServiceDocID: posd.ID,
				UploadID:            u.ID,
				DeletedAt:           &deletedAt,
			},
		})

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), pr.ID)

		suite.NoError(err)
		suite.Equal(pr.ID, paymentRequest.ID)

		suite.Len(paymentRequest.ProofOfServiceDocs, 1)
		suite.Equal(posd.ID, paymentRequest.ProofOfServiceDocs[0].ID)

		suite.Len(paymentRequest.ProofOfServiceDocs[0].PrimeUploads, 0)
	})

	suite.T().Run("if there is an error, we get it with zero payment request", func(t *testing.T) {
		fetcher := NewPaymentRequestFetcher()

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), uuid.Nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(models.PaymentRequest{}, paymentRequest)
	})
}
