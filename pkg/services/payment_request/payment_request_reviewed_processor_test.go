package paymentrequest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/sequence"
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
	err := os.Setenv("SYNCADA_SFTP_PORT", "1234")
	suite.FatalNoError(err)
	err = os.Setenv("SYNCADA_SFTP_USER_ID", "FAKE_USER_ID")
	suite.FatalNoError(err)
	err = os.Setenv("SYNCADA_SFTP_IP_ADDRESS", "127.0.0.1")
	suite.FatalNoError(err)
	err = os.Setenv("SYNCADA_SFTP_PASSWORD", "FAKE PASSWORD")
	suite.FatalNoError(err)
	err = os.Setenv("SYNCADA_SFTP_INBOUND_DIRECTORY", "/Dropoff")
	suite.FatalNoError(err)
	// generated fake host key to pass parser used following command and only saved the pub key
	//   ssh-keygen -q -N "" -t ecdsa -f /tmp/ssh_host_ecdsa_key
	err = os.Setenv("SYNCADA_SFTP_HOST_KEY", "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI+M4xIGU6D4On+Wxz9k/QT12TieNvaXA0lvosnW135MRQzwZp5VDThQ6Vx7yhp18shgjEIxFHFTLxpmUc6JdMc= fake@localhost")
	suite.FatalNoError(err)

	var responseSuccess = http.Response{}
	responseSuccess.StatusCode = http.StatusOK
	responseSuccess.Status = "200 Success"

	var responseFailure = http.Response{}
	responseFailure.StatusCode = http.StatusInternalServerError
	responseFailure.Status = "500 Internal Server Error"

	suite.Run("process reviewed payment request successfully (0 Payments to review)", func() {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())
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
		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.Run("process reviewed payment request successfully (1 payment request reviewed all rejected excluded)", func() {
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
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())
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
		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		// Ensure that payment requst was not sent to gex
		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, err := fetcher.FetchPaymentRequest(rejectedPaymentRequest.ID)
		suite.NoError(err)
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

	suite.Run("process reviewed payment request successfully (do not send file)", func() {
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
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())

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
		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		// Ensure that sent_to_gex_at timestamp has been added
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, fetchErr := fetcher.FetchPaymentRequest(pr.ID)
			suite.NoError(fetchErr)
			suite.NotNil(paymentRequest.SentToGexAt)
			suite.Equal(false, paymentRequest.SentToGexAt.IsZero())
			suite.Equal(models.PaymentRequestStatusSentToGex, paymentRequest.Status)
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

	suite.Run("process reviewed payment request, failed EDI generator", func() {
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
		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, fetchErr := fetcher.FetchPaymentRequest(pr.ID)
			suite.NoError(fetchErr)
			suite.Nil(paymentRequest.SentToGexAt)
			suite.Equal(models.PaymentRequestStatusEDIError, paymentRequest.Status)
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
		suite.Len(ediErrors, len(prs))

		prsIDs := []string{}
		ediErrorPRIDs := []string{}
		for idx := range ediErrors {
			suite.Contains(*(ediErrors[idx].Description), "test error")
			prsIDs = append(prsIDs, prs[idx].ID.String())
			ediErrorPRIDs = append(ediErrorPRIDs, ediErrors[idx].PaymentRequestID.String())
		}
		sort.Slice(prsIDs, func(i, j int) bool {
			return prsIDs[i] < prsIDs[j]
		})
		sort.Slice(ediErrorPRIDs, func(i, j int) bool {
			return ediErrorPRIDs[i] < ediErrorPRIDs[j]
		})
		suite.Equal(prsIDs, ediErrorPRIDs)

		// Make sure that PR status is updated
		var updatedPaymentRequest models.PaymentRequest
		err = suite.DB().Where("id = ?", prs[0].ID).First(&updatedPaymentRequest)
		suite.NoError(err)
		suite.Equal(updatedPaymentRequest.Status, models.PaymentRequestStatusEDIError)
	})

	suite.Run("process reviewed payment request, failed EDI generator (mock GEX HTTP)", func() {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		SFTPSender := services.SyncadaSFTPSender(nil)
		sendToSyncada := false

		// Get list of PRs before processing them
		prs := suite.createPaymentRequest(2)

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
		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			var paymentRequest models.PaymentRequest
			paymentRequest, err = fetcher.FetchPaymentRequest(pr.ID)
			suite.NoError(err)
			suite.Nil(paymentRequest.SentToGexAt)
			suite.Equal(models.PaymentRequestStatusEDIError, paymentRequest.Status)
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
		suite.Len(ediErrors, len(prs))
		for idx := range ediErrors {
			suite.Contains(*(ediErrors[idx].Description), "test error")
			suite.Equal(ediErrors[idx].PaymentRequestID, prs[idx].ID)
		}
	})

	suite.Run("process reviewed payment request, failed payment request fetcher", func() {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		prs := suite.createPaymentRequest(2)

		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())
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

		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		// Ensure that sent_to_gex_at is Nil on unsucessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, fetchErr := fetcher.FetchPaymentRequest(pr.ID)
			suite.NoError(fetchErr)
			suite.Nil(paymentRequest.SentToGexAt)
			suite.Equal(models.PaymentRequestStatusReviewed, paymentRequest.Status)
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

	suite.Run("process reviewed payment request, Failure SendToSyncada", func() {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		pr := suite.createPaymentRequest(1)[0]

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())
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

		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		// Ensure that sent_to_gex_at is Nil on unsuccessful call to processReviewedPaymentRequest service
		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, err := fetcher.FetchPaymentRequest(pr.ID)
		suite.NoError(err)
		suite.Nil(paymentRequest.SentToGexAt)
		suite.Equal(models.PaymentRequestStatusReviewed, paymentRequest.Status)

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)
	})

	suite.Run("process reviewed payment request, successful SendToSyncada", func() {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		numPrs := 4
		prs := suite.createPaymentRequest(numPrs)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())
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

		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		// There are 4 in this test
		suite.Equal(4, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Greater(newCount, countProcessingRecordsBefore)
		suite.Equal(countProcessingRecordsBefore+1, newCount)

		// Ensure that status is updated to SENT_TO_GEX when PRs are sent successfully
		fetcher := NewPaymentRequestFetcher(suite.DB())
		for _, pr := range prs {
			paymentRequest, err := fetcher.FetchPaymentRequest(pr.ID)
			suite.NoError(err)
			suite.Equal(models.PaymentRequestStatusSentToGex, paymentRequest.Status)
		}
	})

	suite.Run("process reviewed payment request, successfully test init function", func() {
		// Run init with no issues
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		_, err := InitNewPaymentRequestReviewedProcessor(suite.DB(), suite.logger, false, icnSequencer, nil)
		suite.NoError(err)
	})
}

func (suite *PaymentRequestServiceSuite) TestProcessReviewedPaymentRequestFailedGEXMock() {
	suite.Run("process reviewed payment request, failed mock GEX HTTP send", func() {
		var ediProcessingBefore models.EDIProcessing
		countProcessingRecordsBefore, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessingBefore)
		suite.NoError(err, "Get count of EDIProcessing")

		numPrs := 2
		prs := suite.createPaymentRequest(numPrs)

		reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())
		icnSequencer := sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName)
		ediGenerator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.NewMock())
		sendToSyncada := true // Call GEXSender but using mock here

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

		paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()

		var ediProcessing models.EDIProcessing
		err = suite.DB().Where("edi_type = ?", models.EDIType858).Order("process_ended_at desc").First(&ediProcessing)
		suite.NoError(err, "Get number of processed files")
		suite.Equal(0, ediProcessing.NumEDIsProcessed)

		newCount, err := suite.DB().Where("edi_type = ?", models.EDIType858).Count(&ediProcessing)
		suite.NoError(err, "Get count of EDIProcessing")
		suite.Equal(countProcessingRecordsBefore+1, newCount)

		fetcher := NewPaymentRequestFetcher(suite.DB())
		paymentRequest, err := fetcher.FetchPaymentRequest(prs[0].ID)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusReviewed, paymentRequest.Status)
	})
}
