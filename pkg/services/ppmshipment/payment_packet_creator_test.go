package ppmshipment

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestCreatePaymentPacket() {
	mockPPMShipmentFetcher := &mocks.PPMShipmentFetcher{}
	mockUserUploadToPDFConverter := &mocks.UserUploadToPDFConverter{}
	mockPDFMerger := &mocks.PDFMerger{}
	mockPPMShipmentUpdater := &mocks.PPMShipmentUpdater{}

	fakeS3 := storageTest.NewFakeS3Storage(true)

	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)

	suite.FatalNoError(uploaderErr)

	paymentPacketCreator := NewPaymentPacketCreator(
		mockPPMShipmentFetcher,
		mockUserUploadToPDFConverter,
		mockPDFMerger,
		mockPPMShipmentUpdater,
		userUploader,
	)

	setUpMockPPMShipmentFetcherForPayment := func(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, returnValue ...interface{}) {
		setUpMockPPMShipmentFetcher(
			mockPPMShipmentFetcher,
			appCtx,
			ppmShipmentID,
			[]string{
				EagerPreloadAssociationServiceMember,
				EagerPreloadAssociationWeightTickets,
				EagerPreloadAssociationProgearWeightTickets,
				EagerPreloadAssociationMovingExpenses,
				EagerPreloadAssociationPaymentPacket,
			},
			[]string{
				PostLoadAssociationWeightTicketUploads,
				PostLoadAssociationProgearWeightTicketUploads,
				PostLoadAssociationMovingExpenseUploads,
				PostLoadAssociationUploadedOrders,
			},
			returnValue...,
		)
	}

	// prepMockInfo is a helper function to prep the data needed for mocks and also passes back a cleanup func that
	// should be deferred to run after the test is done.
	prepMockInfo := func(ppmShipment *models.PPMShipment) (
		models.UserUploads,
		[]*services.FileInfo,
		[]io.ReadCloser,
		func(),
	) {
		userUploads := models.UserUploads{}
		fileInfoSet := []*services.FileInfo{}

		uploadedOrdersUserUpload := &ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0]

		if uploadedOrdersUserUpload.ID.IsNil() {
			uploadedOrdersUserUpload.ID = uuid.Must(uuid.NewV4())
		}

		ordersFileInfo := services.NewFileInfo(uploadedOrdersUserUpload, factory.FixtureOpen("test.png"))
		ordersFileInfo.PDFStream = factory.FixtureOpen("filled-out-orders.pdf")

		userUploads = append(userUploads, *uploadedOrdersUserUpload)
		fileInfoSet = append(fileInfoSet, ordersFileInfo)

		for _, weightTicket := range ppmShipment.WeightTickets {
			weightTicket := weightTicket

			for _, userUpload := range weightTicket.EmptyDocument.UserUploads {
				userUpload := userUpload

				emptyDocumentFileInfo := services.NewFileInfo(&userUpload, factory.FixtureOpen("empty-weight-ticket.png"))
				emptyDocumentFileInfo.PDFStream = factory.FixtureOpen("empty-weight-ticket.pdf")

				userUploads = append(userUploads, userUpload)
				fileInfoSet = append(fileInfoSet, emptyDocumentFileInfo)
			}

			for _, userUpload := range weightTicket.FullDocument.UserUploads {
				userUpload := userUpload

				fullDocumentFileInfo := services.NewFileInfo(&userUpload, factory.FixtureOpen("full-weight-ticket.png"))
				fullDocumentFileInfo.PDFStream = factory.FixtureOpen("full-weight-ticket.pdf")

				userUploads = append(userUploads, userUpload)
				fileInfoSet = append(fileInfoSet, fullDocumentFileInfo)
			}

			for _, userUpload := range weightTicket.ProofOfTrailerOwnershipDocument.UserUploads {
				userUpload := userUpload

				proofOfTrailerOwnershipDocumentFileInfo := services.NewFileInfo(&userUpload, factory.FixtureOpen("test.png"))
				proofOfTrailerOwnershipDocumentFileInfo.PDFStream = factory.FixtureOpen("test.pdf")

				userUploads = append(userUploads, userUpload)
				fileInfoSet = append(fileInfoSet, proofOfTrailerOwnershipDocumentFileInfo)
			}
		}

		for _, progearWeightTicket := range ppmShipment.ProgearWeightTickets {
			progearWeightTicket := progearWeightTicket

			for _, userUpload := range progearWeightTicket.Document.UserUploads {
				userUpload := userUpload

				documentFileInfo := services.NewFileInfo(&userUpload, factory.FixtureOpen("empty-weight-ticket.png"))
				documentFileInfo.PDFStream = factory.FixtureOpen("empty-weight-ticket.pdf")

				userUploads = append(userUploads, userUpload)
				fileInfoSet = append(fileInfoSet, documentFileInfo)
			}
		}

		for _, movingExpense := range ppmShipment.MovingExpenses {
			movingExpense := movingExpense

			for _, userUpload := range movingExpense.Document.UserUploads {
				userUpload := userUpload

				documentFileInfo := services.NewFileInfo(&userUpload, factory.FixtureOpen("test.png"))
				documentFileInfo.PDFStream = factory.FixtureOpen("test.pdf")

				userUploads = append(userUploads, userUpload)
				fileInfoSet = append(fileInfoSet, documentFileInfo)
			}
		}

		cleanUpFunc := func() {
			for _, fileInfo := range fileInfoSet {
				fileInfo.OriginalUploadStream.Close()
				fileInfo.PDFStream.Close()
			}
		}

		pdfStreams := []io.ReadCloser{}

		for _, fileInfo := range fileInfoSet {
			pdfStreams = append(pdfStreams, fileInfo.PDFStream)
		}

		return userUploads, fileInfoSet, pdfStreams, cleanUpFunc
	}

	suite.Run("returns an error if the PPMShipmentFetcher returns an error", func() {
		appCtx := suite.AppContextForTest()

		ppmShipmentID := uuid.Must(uuid.NewV4())

		fakeErr := apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipmentID, nil, fakeErr)

		err := paymentPacketCreator.CreatePaymentPacket(appCtx, ppmShipmentID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)

			suite.ErrorContains(err, "error creating payment packet: failed to load PPMShipment")
		}
	})

	// TODO: add test case(s) for the SSW gen call

	suite.Run("returns an error if we get an error trying to convert uploaded files to PDFs", func() {
		ppmShipment := factory.BuildPPMShipmentWithAllDocTypesApprovedMissingPaymentPacket(
			nil,
			nil,
			[]factory.Customization{
				{
					Model: models.PPMShipment{
						ID: uuid.Must(uuid.NewV4()),
					},
				},
			},
		)

		userUploads, _, _, cleanUpFunc := prepMockInfo(&ppmShipment)

		defer cleanUpFunc()

		setUpMockPPMShipmentFetcherForPayment(suite.AppContextForTest(), ppmShipment.ID, &ppmShipment, nil)

		uploadedOrdersUserUpload := ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0]

		fakeErr := fmt.Errorf(
			"failed to convert file %s (UserUploadID: %s) to PDF",
			uploadedOrdersUserUpload.Upload.Filename,
			uploadedOrdersUserUpload.ID,
		)

		setUpMockUserUploadToPDFConverter(
			mockUserUploadToPDFConverter,
			suite.AppContextForTest(),
			userUploads,
			nil,
			fakeErr,
		)

		err := paymentPacketCreator.CreatePaymentPacket(suite.AppContextForTest(), ppmShipment.ID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)

			suite.ErrorContains(err, "error creating payment packet: failed to convert uploads to PDF")
		}
	})

	suite.Run("returns an error if we get an error trying to merge the PDFs", func() {
		ppmShipment := factory.BuildPPMShipmentWithAllDocTypesApprovedMissingPaymentPacket(
			nil,
			nil,
			[]factory.Customization{
				{
					Model: models.PPMShipment{
						ID: uuid.Must(uuid.NewV4()),
					},
				},
			},
		)

		userUploads, fileInfoSet, pdfStreams, cleanUpFunc := prepMockInfo(&ppmShipment)

		defer cleanUpFunc()

		setUpMockPPMShipmentFetcherForPayment(suite.AppContextForTest(), ppmShipment.ID, &ppmShipment, nil)

		setUpMockUserUploadToPDFConverter(
			mockUserUploadToPDFConverter,
			suite.AppContextForTest(),
			userUploads,
			fileInfoSet,
			nil,
		)

		fakeErr := fmt.Errorf("failed to merge PDFs")

		setUpMockPDFMerger(mockPDFMerger, suite.AppContextForTest(), pdfStreams, nil, fakeErr)

		err := paymentPacketCreator.CreatePaymentPacket(suite.AppContextForTest(), ppmShipment.ID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)

			suite.ErrorContains(err, "error creating payment packet: failed to merge PDFs")
		}
	})

	suite.Run("returns an error if we get an error trying to save the payment packet", func() {
		// These tests rely on failures being raised because of bad service member IDs, but the important part is us
		// getting the error back, so this could change to be anything that triggers an error when saving. Service
		// member ID is mainly chosen because it's one of the first things we can error on and is easy to set up.

		ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(
			nil,
			nil,
			[]factory.Customization{
				{
					Model: models.PPMShipment{
						ID: uuid.Must(uuid.NewV4()),
					},
				},
				{
					Model: models.ServiceMember{
						ID: uuid.Nil,
					},
				},
				{
					Model: models.UserUpload{
						ID: uuid.Must(uuid.NewV4()),
					},
				},
			},
		)

		userUploads, fileInfoSet, pdfStreams, cleanUpFunc := prepMockInfo(&ppmShipment)

		defer cleanUpFunc()

		setUpMockPPMShipmentFetcherForPayment(suite.AppContextForTest(), ppmShipment.ID, &ppmShipment, nil)

		setUpMockUserUploadToPDFConverter(
			mockUserUploadToPDFConverter,
			suite.AppContextForTest(),
			userUploads,
			fileInfoSet,
			nil,
		)

		mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

		defer mockMergedPDF.Close()

		setUpMockPDFMerger(mockPDFMerger, suite.AppContextForTest(), pdfStreams, mockMergedPDF, nil)

		err := paymentPacketCreator.CreatePaymentPacket(suite.AppContextForTest(), ppmShipment.ID)

		if suite.Error(err) {
			suite.ErrorContains(err, "error creating payment packet: failed to save payment packet")

			suite.ErrorContains(err, "failed to create payment packet document")

			suite.ErrorContains(err, "ServiceMemberID can not be blank")
		}
	})

	suite.Run("returns nil if all goes well", func() {
		ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(suite.DB(), nil, nil)

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since the
		// savePaymentPacket function will be using a transaction.
		suite.NoError(suite.AppContextForTest().NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			userUploads, fileInfoSet, pdfStreams, cleanUpFunc := prepMockInfo(&ppmShipment)

			defer cleanUpFunc()

			setUpMockPPMShipmentFetcherForPayment(txnAppCtx, ppmShipment.ID, &ppmShipment, nil)

			setUpMockUserUploadToPDFConverter(
				mockUserUploadToPDFConverter,
				txnAppCtx,
				userUploads,
				fileInfoSet,
				nil,
			)

			mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

			defer mockMergedPDF.Close()

			setUpMockPDFMerger(mockPDFMerger, txnAppCtx, pdfStreams, mockMergedPDF, nil)

			setUpMockPPMShipmentUpdater(
				mockPPMShipmentUpdater,
				txnAppCtx,
				&ppmShipment,
				// This function will get called instead of the regular update function, so it needs to have the same
				// signature.
				func(_ appcontext.AppContext, ppmShipment *models.PPMShipment, _ uuid.UUID) (*models.PPMShipment, error) {
					// We'll just pass it back. In reality, the updatedAt field would have been updated, but it's not
					// super relevant to what we're testing.
					return ppmShipment, nil
				},
			)

			err := paymentPacketCreator.CreatePaymentPacket(txnAppCtx, ppmShipment.ID)

			if suite.NoError(err) {
				suite.NotNil(ppmShipment.PaymentPacketID)
			}

			return nil
		}))
	})
}

func (suite *PPMShipmentSuite) TestSavePaymentPacket() {
	mockPPMShipmentUpdater := &mocks.PPMShipmentUpdater{}

	fakeS3 := storageTest.NewFakeS3Storage(true)

	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)

	suite.FatalNoError(uploaderErr)

	suite.Run("returns an error if we fail to create the packet document", func() {
		badServiceMemberIDTestCases := map[string]struct {
			serviceMemberID uuid.UUID
			expectedErrMsg  string
		}{
			"empty UUID": {
				serviceMemberID: uuid.Nil,
				expectedErrMsg:  "ServiceMemberID can not be blank",
			},
			"bad UUID": {
				serviceMemberID: uuid.Must(uuid.NewV4()),
				expectedErrMsg:  "insert or update on table \"documents\" violates foreign key constraint \"documents_service_members_id_fk\"",
			},
		}

		for name, testCase := range badServiceMemberIDTestCases {
			name, testCase := name, testCase

			// These tests rely on failures being raised because of bad service member IDs, but the important part is us
			// getting the error back, so this could change to be anything that triggers an error when creating the
			// document. Service member ID is mainly chosen because it's easy to set up to trigger both validation and
			// saving errors.
			suite.Run(fmt.Sprintf("bad service member ID: %s", name), func() {
				ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(
					nil,
					nil,
					[]factory.Customization{
						{
							Model: models.PPMShipment{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
						{
							Model: models.ServiceMember{
								ID: testCase.serviceMemberID,
							},
						},
						{
							Model: models.UserUpload{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
				)

				mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

				defer mockMergedPDF.Close()

				err := savePaymentPacket(suite.AppContextForTest(), &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

				if suite.Error(err) {
					suite.ErrorContains(err, "failed to create payment packet document")

					suite.ErrorContains(err, testCase.expectedErrMsg)
				}
			})
		}
	})

	suite.Run("returns an error if we fail to update the PPM shipment", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(appCtx.DB(), nil, nil)

		suite.FatalNil(ppmShipment.PaymentPacketID)

		mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

		defer mockMergedPDF.Close()

		fakeError := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since the
		// savePaymentPacket function will be using a transaction.
		suite.NoError(appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			setUpMockPPMShipmentUpdater(
				mockPPMShipmentUpdater,
				txnAppCtx,
				&ppmShipment,
				nil,
				fakeError,
			)

			err := savePaymentPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

			if suite.Error(err) {
				suite.ErrorIs(err, fakeError)

				suite.ErrorContains(err, "failed to update PPMShipment with payment packet document")
			}

			return nil
		}))
	})

	suite.Run("returns an error if we fail to prepare the file for upload", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(appCtx.DB(), nil, nil)

		suite.FatalNil(ppmShipment.PaymentPacketID)

		mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

		// If the file is closed, it should trigger an error when we try to prepare it for upload.
		mockMergedPDF.Close()

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since the
		// savePaymentPacket function will be using a transaction.
		suite.NoError(appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			setUpMockPPMShipmentUpdater(
				mockPPMShipmentUpdater,
				txnAppCtx,
				&ppmShipment,
				&ppmShipment,
				nil,
			)

			err := savePaymentPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

			if suite.Error(err) {
				suite.ErrorContains(err, "failed to prepare payment packet for upload")
				suite.ErrorContains(err, "Error copying incoming data into afero file")
			}

			return nil
		}))
	})

	suite.Run("returns an error if we fail to upload the file", func() {
		badUserIDTestCases := map[string]struct {
			userID        uuid.UUID
			expectedError string
			txnErrCheck   func(error)
		}{
			"blank UUID": {
				userID:        uuid.Nil,
				expectedError: "UploaderID can not be blank.",
				txnErrCheck: func(err error) {
					suite.NoError(err)
				},
			},
			"bad UUID": {
				userID:        uuid.Must(uuid.NewV4()),
				expectedError: "insert or update on table \"user_uploads\" violates foreign key constraint \"user_uploads_uploader_id_fkey\"",
				txnErrCheck: func(err error) {
					// Since we're triggering a DB error, there is a transaction error that gets raised.
					suite.Error(err)
				},
			},
		}

		for name, testCase := range badUserIDTestCases {
			name, testCase := name, testCase

			suite.Run(fmt.Sprintf("UserID error: %s", name), func() {
				ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(suite.DB(), nil, nil)

				suite.FatalNil(ppmShipment.PaymentPacketID)

				mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

				defer mockMergedPDF.Close()

				// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since
				// the savePaymentPacket function will be using a transaction.
				txnErr := suite.AppContextForTest().NewTransaction(func(txnAppCtx appcontext.AppContext) error {
					setUpMockPPMShipmentUpdater(
						mockPPMShipmentUpdater,
						txnAppCtx,
						&ppmShipment,
						&ppmShipment,
						nil,
					)

					// Setting a bad ID on this should trigger the upload to fail
					ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = testCase.userID

					err := savePaymentPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

					if suite.Error(err) {
						suite.ErrorContains(err, "failed to upload payment packet")

						suite.ErrorContains(err, testCase.expectedError)
					}

					return nil
				})

				testCase.txnErrCheck(txnErr)
			})
		}
	})

	suite.Run("returns nil if all goes well", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := factory.BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(appCtx.DB(), nil, nil)

		suite.FatalNil(ppmShipment.PaymentPacketID)

		mockMergedPDF := factory.FixtureOpen("payment-packet.pdf")

		defer mockMergedPDF.Close()

		expectedBytes, readExpectedErr := io.ReadAll(mockMergedPDF)

		suite.FatalNoError(readExpectedErr)

		_, seekErr := mockMergedPDF.Seek(0, io.SeekStart)

		suite.FatalNoError(seekErr)

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since
		// the savePaymentPacket function will be using a transaction.
		suite.NoError(appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			setUpMockPPMShipmentUpdater(
				mockPPMShipmentUpdater,
				txnAppCtx,
				&ppmShipment,
				// This function will get called instead of the regular update function, so it needs to have the same
				// signature.
				func(_ appcontext.AppContext, ppmShipment *models.PPMShipment, _ uuid.UUID) (*models.PPMShipment, error) {
					// We'll just pass it back. In reality, the updatedAt field would have been updated, but it's not
					// super relevant to what we're testing.
					return ppmShipment, nil
				},
			)

			err := savePaymentPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

			suite.NoError(err)

			return nil
		}))

		// Now we'll double check everything to make sure it was saved correctly.
		if suite.NotNil(ppmShipment.PaymentPacketID) {
			download, downloadErr := userUploader.Download(appCtx, &ppmShipment.PaymentPacket.UserUploads[0])

			suite.FatalNoError(downloadErr)

			actualBytes, readActualErr := io.ReadAll(download)

			suite.FatalNoError(readActualErr)

			suite.Equal(expectedBytes, actualBytes)
		}
	})
}
