//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used set up environment variables
//RA: Given the functions causing the lint errors are used to set environment variables for testing purposes, it does not present a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package paymentrequest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/services"

	"github.com/stretchr/testify/mock"

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

	var responseSuccess = http.Response{}
	responseSuccess.StatusCode = http.StatusOK
	responseSuccess.Status = "200 Success"

	var responseFailure = http.Response{}
	responseFailure.StatusCode = http.StatusInternalServerError
	responseFailure.Status = "500 Internal Server Error"

	suite.T().Run("process reviewed payment request successfully (0 Payments to review)", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		gexSender := services.GexSender(nil)
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
		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.T().Run("process reviewed payment request successfully (1 payment request reviewed all rejected excluded)", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		rejectionReason := "Voided"
		rejectedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: &rejectionReason,
			},
		})

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		gexSender := services.GexSender(nil)
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
		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		// Ensure that payment requst was not sent to gex
		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, _ := fetcher.FetchPaymentRequest(rejectedPaymentRequest.ID)
		suite.Nil(paymentRequest.SentToGexAt)
		suite.Equal(rejectedPaymentRequest.Status, models.PaymentRequestStatusReviewedAllRejected)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.T().Run("process reviewed payment request successfully (do not send file)", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		prs := suite.createPaymentRequest(4)

		_ = testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				ID: uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
			},
		})

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())

		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		gexSender := services.GexSender(nil)
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
		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		// Ensure that sent_to_gex_at timestamp has been added
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.NotNil(paymentRequest.SentToGexAt)
			suite.Equal(false, paymentRequest.SentToGexAt.IsZero())
		}

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(len(prs), ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.T().Run("process reviewed payment request, failed EDI generator", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		prs := suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		gexSender := services.GexSender(nil)
		sendToSyncada := false

		// ediinvoice.Invoice858C, error
		ediGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		ediGenerator.On("InitDB", mock.IsType(&pop.Connection{}))
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
		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "function ProcessReviewedPaymentRequest failed call")

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)

		// Check that an error was recorded in the EdiErrors table.
		var ediErrors models.EdiErrors
		err = suite.DB().Where("edi_type = ?", models.EDIType858).All(&ediErrors)
		suite.NoError(err)
		// ProcessReviewedPaymentRequest() stops processing requests after it hits an error, so
		// we only expect the first payment request with an error to be recorded.
		suite.Len(ediErrors, 1)
		suite.Contains(*(ediErrors[0].Description), "test error")
		suite.Equal(ediErrors[0].PaymentRequestID, prs[0].ID)

		// Make sure that PR status is updated
		var updatedPaymentRequest models.PaymentRequest
		err = suite.DB().Where("id = ?", prs[0].ID).First(&updatedPaymentRequest)
		suite.NoError(err)
		suite.Equal(updatedPaymentRequest.Status, models.PaymentRequestStatusEDIError)
	})

	suite.T().Run("process reviewed payment request, failed EDI generator (mock GEX HTTP)", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		SFTPSender := services.SyncadaSFTPSender(nil)
		sendToSyncada := false

		// Get list of PRs before processing them
		prs, err := reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest()
		suite.NoError(err)

		// Record PR statuses
		type prStatus struct {
			id     uuid.UUID
			status models.PaymentRequestStatus
		}
		type prStatuses []prStatus

		// Set up mock HTTP server and mock GEX
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		mockGexSender := invoice.NewGexSenderHTTP(mockServer.URL, "", false, nil, "", "")
		if mockGexSender == nil {
			suite.T().Fatal("Failed to create mockGexSender")
		}

		// ediinvoice.Invoice858C, error
		ediGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		ediGenerator.On("InitDB", mock.IsType(&pop.Connection{}))
		ediGenerator.
			On("Generate", mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, errors.New("test error"))

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			mockGexSender,
			SFTPSender)
		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "function ProcessReviewedPaymentRequest failed call")

		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		afterProcessingStatus := prStatuses{}
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			var paymentRequest models.PaymentRequest
			paymentRequest, err = fetcher.FetchPaymentRequest(pr.ID)
			suite.NoError(err)
			suite.Nil(paymentRequest.SentToGexAt)
			afterProcessingStatus = append(afterProcessingStatus, prStatus{id: paymentRequest.ID, status: paymentRequest.Status})
		}

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)

		// Check that an error was recorded in the EdiErrors table.
		var ediError models.EdiError
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("created_at desc").First(&ediError)
		suite.NoError(err)

		// ProcessReviewedPaymentRequest() stops processing requests after it hits an error, so
		// we only expect the first payment request with an error to be recorded.
		suite.Contains(*(ediError.Description), "test error")
		paymentRequest, _ := fetcher.FetchPaymentRequest(ediError.PaymentRequestID)
		suite.Equal(ediError.PaymentRequestID, paymentRequest.ID)

		countUpdated := 0
		foundUpdatedPR := false

		for _, pr := range prs {
			for _, uPR := range afterProcessingStatus {
				if pr.ID == uPR.id {
					if pr.Status != uPR.status {
						suite.Equal(ediError.PaymentRequestID, uPR.id)
						suite.Equal(models.PaymentRequestStatusEDIError, uPR.status)
						foundUpdatedPR = true
						countUpdated++
					}
				}
			}
		}
		suite.True(foundUpdatedPR, "Found expected PR with EDI_ERROR")
		suite.Equal(1, countUpdated, "Expected 1 update to PR status")

	})

	suite.T().Run("process reviewed payment request, failed payment request fetcher", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		prs := suite.createPaymentRequest(4)

		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
		suite.NoError(SFTPSessionError)
		gexSender := services.GexSender(nil)
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

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "function ProcessReviewedPaymentRequest failed call")

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.T().Run("process reviewed payment request, fail SFTP send", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		prs := suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		gexSender := services.GexSender(nil)
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

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "error sending the following EDI")

		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Nil(paymentRequest.SentToGexAt)
		}

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.T().Run("process reviewed payment request, successful SFTP send", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		numPrs := 4
		_ = suite.createPaymentRequest(numPrs)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		gexSender := services.GexSender(nil)
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

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		// This test creates 4 payment requests, and there are 9 PRs from previous tests that didn't get statuses changed
		suite.Equal(13, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)

	})

	suite.T().Run("process reviewed payment request, fail POST to GEX", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		pr := suite.createPaymentRequest(1)[0]

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		sftpSender := services.SyncadaSFTPSender(nil)
		sendToSyncada := true // Call SendToSyncadaViaSFTP but using mock here

		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", mock.Anything, mock.Anything).Return(nil, errors.New("test error"))

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			sftpSender)

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "error sending the following EDI")

		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
		suite.Nil(paymentRequest.SentToGexAt)
		// TODO: bug, when GEX is the reason for the failure or even SFTP we shouldn't
		// TODO: mark the EDI status as failed, it should be marked as REVIEWED so that it can be retried.
		// TODO: created bug to fix this https://dp3.atlassian.net/browse/MB-7736
		suite.Equal(models.PaymentRequestStatusEDIError, paymentRequest.Status)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.T().Run("process reviewed payment request, non-200 response from GEX", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		pr := suite.createPaymentRequest(1)[0]

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		sftpSender := services.SyncadaSFTPSender(nil)
		sendToSyncada := true

		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", mock.Anything, mock.Anything).Return(&responseFailure, nil)

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			sftpSender)

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "error sending the following EDI")

		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
		suite.Nil(paymentRequest.SentToGexAt)
		// TODO: bug, when GEX is the reason for the failure or even SFTP we shouldn't
		// TODO: mark the EDI status as failed, it should be marked as REVIEWED so that it can be retried.
		// TODO: created bug to fix this https://dp3.atlassian.net/browse/MB-7736
		suite.Equal(models.PaymentRequestStatusEDIError, paymentRequest.Status)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})
	suite.T().Run("process reviewed payment request, successful POST to GEX", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		numPrs := 4
		prs := suite.createPaymentRequest(numPrs)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		sftpSender := services.SyncadaSFTPSender(nil)
		sendToSyncada := true // Call SendToSyncadaViaSFTP but using mock here

		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", mock.Anything, mock.Anything).Return(&responseSuccess, nil)

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			gexSender,
			sftpSender)

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(4, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)

		// Ensure that status is updated to SENT_TO_GEX when PRs are sent successfully
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			suite.Equal(models.PaymentRequestStatusSentToGex, paymentRequest.Status)
		}
	})

	suite.T().Run("process reviewed payment request, failed due to both senders being nil", func(t *testing.T) {

		prs := suite.createPaymentRequest(4)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := services.GexSender(nil)
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
		_, err := InitNewPaymentRequestReviewedProcessor(suite.DB(), suite.logger, false, suite.icnSequencer, nil)
		suite.NoError(err)
	})
}

func (suite *PaymentRequestServiceSuite) TestProcessReviewedPaymentRequestFailedGEXMock() {
	suite.T().Run("process reviewed payment request, failed mock GEX HTTP send", func(t *testing.T) {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		numPrs := 2
		prs := suite.createPaymentRequest(numPrs)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
		sendToSyncada := true // Call SendToSyncadaViaSFTP but using mock here

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}))
		mockGexSender := invoice.NewGexSenderHTTP(mockServer.URL, "", false, nil, "", "")
		if mockGexSender == nil {
			suite.T().Fatal("Failed to create mockGexSender")
		}

		sftpSender := services.SyncadaSFTPSender(nil)

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
			suite.DB(),
			suite.logger,
			reviewedPaymentRequestFetcher,
			ediGenerator,
			sendToSyncada,
			mockGexSender,
			sftpSender)

		err = paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "error sending the following EDI")

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Equal(countProcessingRecordsBefore+1, newCount)

		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, _ := fetcher.FetchPaymentRequest(prs[0].ID)
		// TODO: bug, when GEX is the reason for the failure or even SFTP we shouldn't
		// TODO: mark the EDI status as failed, it should be marked as REVIEWED so that it can be retried.
		// TODO: created bug to fix this https://dp3.atlassian.net/browse/MB-7736
		suite.Equal(models.PaymentRequestStatusEDIError, paymentRequest.Status)
	})
}

func (suite *PaymentRequestServiceSuite) lockPR(prID uuid.UUID) {
	query := `
		BEGIN;
		SELECT * FROM payment_requests
		WHERE id = $1 FOR NO KEY UPDATE SKIP LOCKED;
		UPDATE payment_requests
		SET
			status = $2,
		WHERE id = $1;
	`
	suite.DB().RawQuery(query, prID, models.PaymentRequestStatusPaid).Exec()
	time.Sleep(1 * time.Second)
	suite.DB().RawQuery(`COMMIT;`).Exec()
}

func (suite *PaymentRequestServiceSuite) TestProcessLockedReviewedPaymentRequest() {
	os.Setenv("SYNCADA_SFTP_PORT", "1234")
	os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
	os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
	os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
	os.Setenv("SYNCADA_SFTP_INBOUND_DIRECTORY", "/Dropoff")
	// generated fake host key to pass parser used following command and only saved the pub key
	//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
	os.Setenv("SYNCADA_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")

	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
	generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock())
	SFTPSession, SFTPSessionError := invoice.InitNewSyncadaSFTPSession()
	suite.NoError(SFTPSessionError)
	gexSender := services.GexSender(nil)
	sendToSyncada := false

	paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(
		suite.DB(),
		suite.logger,
		reviewedPaymentRequestFetcher,
		generator,
		sendToSyncada,
		gexSender,
		SFTPSession)

	suite.T().Run("successfully process prs even when a locked row has a delay", func(t *testing.T) {
		reviewedPaymentRequests := suite.createPaymentRequest(2)

		go suite.lockPR(reviewedPaymentRequests[0].ID)

		for _, pr := range reviewedPaymentRequests {
			err := paymentRequestReviewedProcessor.ProcessAndLockReviewedPR(pr)
			suite.NoError(err)
		}

		fetcher := NewPaymentRequestFetcher(suite.DB())
		for i, pr := range reviewedPaymentRequests {
			paymentRequest, _ := fetcher.FetchPaymentRequest(pr.ID)
			if i == 0 {
				suite.Equal(models.PaymentRequestStatusSentToGex, paymentRequest.Status)
			} else {
				suite.Equal(models.PaymentRequestStatusSentToGex, paymentRequest.Status)
			}
		}
	})
}
