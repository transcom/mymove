package paymentrequest

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const testDateFormat = "060102"

func (suite *PaymentRequestServiceSuite) createPaymentRequest(num int) models.PaymentRequests {
	var paymentRequests models.PaymentRequests
	for i := 0; i < num; i++ {
		currentTime := time.Now()
		basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTime.Format(testDateFormat),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "4242",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2424",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip5,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "24245",
			},
		}

		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: mto,
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		})

		requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)
		scheduledPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 20, 0, 0, 0, 0, time.UTC)
		actualPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 22, 0, 0, 0, 0, time.UTC)

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
				ScheduledPickupDate: &scheduledPickupDate,
				ActualPickupDate:    &actualPickupDate,
			},
		})

		assertions := testdatagen.Assertions{
			Move:           mto,
			MTOShipment:    mtoShipment,
			PaymentRequest: paymentRequest,
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		}

		// dlh
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
			assertions,
		)
		// fsc
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeFSC,
			basicPaymentServiceItemParams,
			assertions,
		)
		// ms
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeMS,
			basicPaymentServiceItemParams,
			assertions,
		)
		// cs
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeCS,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dsh
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDSH,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dop
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOP,
			basicPaymentServiceItemParams,
			assertions,
		)
		// ddp
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDP,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dpk
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDPK,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dupk
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDUPK,
			basicPaymentServiceItemParams,
			assertions,
		)
		paymentRequests = append(paymentRequests, paymentRequest)
	}
	return paymentRequests
}

func (suite *PaymentRequestServiceSuite) TestProcessReviewedPaymentRequest() {

	os.Setenv("SYNCADA_SFTP_PORT", "1234")
	os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
	os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
	os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
	os.Setenv("SYNCADA_SFTP_INBOUND_DIRECTORY", "/Dropoff")
	// generated fake host key to pass parser used following command and only saved the pub key
	//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
	os.Setenv("SYNCADA_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")

	suite.T().Run("process reviewed payment request successfully (0 Payments to review)", func(t *testing.T) {
		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := false

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			generator,
			sendToSyncada,
			gexSender,
			SFTPSession)
		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)
	})

	suite.T().Run("process reviewed payment request successfully (1 payment request reviewed all rejected excluded)", func(t *testing.T) {
		rejectionReason := "Voided"
		rejectedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: &rejectionReason,
			},
		})

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := false

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			generator,
			sendToSyncada,
			gexSender,
			SFTPSession)
		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		// Ensure that payment requst was not sent to gex
		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, _ := fetcher.FetchPaymentRequest(rejectedPaymentRequest.ID)
		suite.Nil(paymentRequest.SentToGexAt)
		suite.Equal(rejectedPaymentRequest.Status, models.PaymentRequestStatusReviewedAllRejected)
	})

	suite.T().Run("process reviewed payment request successfully (do not send file)", func(t *testing.T) {

		prs := suite.createPaymentRequest(4)

		_ = testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				ID: uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
			},
		})

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())

		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := false

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			generator,
			sendToSyncada,
			gexSender,
			SFTPSession)
		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		// Ensure that sent_to_gex_at timestamp has been added
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.NotNil(paymentRequest.SentToGexAt)
			suite.Equal(false, paymentRequest.SentToGexAt.IsZero())
		}
	})

	suite.T().Run("process reviewed payment request, failed EDI generator", func(t *testing.T) {

		prs := suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := false

		// ediinvoice.Invoice858C, error
		ediGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		ediGenerator.
			On("Generate", mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, errors.New("test error"))

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			SFTPSession)
		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "function ProcessReviewedPaymentRequest failed call")

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}
	})

	suite.T().Run("process reviewed payment request, failed payment request fetcher", func(t *testing.T) {

		prs := suite.createPaymentRequest(4)

		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := false

		// models.PaymentRequests, error
		reviewedPaymentRequestFetcher := &mocks.PaymentRequestReviewedFetcher{}
		reviewedPaymentRequestFetcher.
			On("FetchReviewedPaymentRequest").Return(models.PaymentRequests{}, errors.New("test error"))

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			SFTPSession)

		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "function ProcessReviewedPaymentRequest failed call")

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}
	})

	suite.T().Run("process reviewed payment request, fail SFTP send", func(t *testing.T) {

		prs := suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := true // Call SendToSyncadaViaSFTP but using mock here

		bytesSent := int64(0)
		// int64, error
		sftpSender := &mocks.SyncadaSFTPSender{}
		sftpSender.
			On("SendToSyncadaViaSFTP", mock.Anything, mock.Anything).Return(bytesSent, errors.New("test error"))

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			sftpSender)

		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "error sending the following EDIs")
		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}
	})

	suite.T().Run("process reviewed payment request, successful SFTP send", func(t *testing.T) {

		_ = suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := true // Call SendToSyncadaViaSFTP but using mock here

		bytesSent := int64(614)
		// int64, error
		sftpSender := &mocks.SyncadaSFTPSender{}
		sftpSender.
			On("SendToSyncadaViaSFTP", mock.Anything, mock.Anything).Return(bytesSent, nil)

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			sftpSender)

		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

	})

	suite.T().Run("process reviewed payment request, failed due to both senders being nil", func(t *testing.T) {

		prs := suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB(), suite.icnSequencer, clock.NewMock())
		var sftpSender services.SyncadaSFTPSender
		sftpSender = nil
		var gexSender services.GexSender
		gexSender = nil
		sendToSyncada := true

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			sftpSender)

		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "senders are nil")

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}
	})

	suite.T().Run("process reviewed payment request, successfully test init function", func(t *testing.T) {
		// Run init with no issues
		_, err := InitNewPaymentRequestReviewedProcessor(suite.DB(), suite.logger, false, suite.icnSequencer)
		suite.NoError(err)
	})
}
