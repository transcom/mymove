package ppmshipment

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services/mocks"
	paperwork "github.com/transcom/mymove/pkg/services/paperwork"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestVerifyAOAPacketSuccess() {

	mockSSWPPMGenerator := &mocks.SSWPPMGenerator{}
	mockSSWPPMComputer := &mocks.SSWPPMComputer{}
	mockPrimeDownloadMoveUploadPDFGenerator := &mocks.PrimeDownloadMoveUploadPDFGenerator{}
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)

	// Create an instance of aoaPacketCreator with mock dependencies
	a := &aoaPacketCreator{
		SSWPPMGenerator:                     mockSSWPPMGenerator,
		SSWPPMComputer:                      mockSSWPPMComputer,
		PrimeDownloadMoveUploadPDFGenerator: mockPrimeDownloadMoveUploadPDFGenerator,
		UserUploader:                        *userUploader,
	}
	// Set up service member ID to verify
	serviceMemberID, err := uuid.NewV4()
	suite.FatalNoError(err)
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: serviceMemberID,
			},
		},
	}, nil)
	// Copy ID to session for passing test
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID,
	})
	ppmShipmentID := ppmShipment.ID
	suite.MustSave(&ppmShipment)

	err = a.VerifyAOAPacketInternal(appCtx, ppmShipmentID)
	suite.Nil(err)

}

func (suite *PPMShipmentSuite) TestVerifyAOAPacketFail() {

	mockSSWPPMGenerator := &mocks.SSWPPMGenerator{}
	mockSSWPPMComputer := &mocks.SSWPPMComputer{}
	mockPrimeDownloadMoveUploadPDFGenerator := &mocks.PrimeDownloadMoveUploadPDFGenerator{}
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)

	// Create an instance of aoaPacketCreator with mock dependencies
	a := &aoaPacketCreator{
		SSWPPMGenerator:                     mockSSWPPMGenerator,
		SSWPPMComputer:                      mockSSWPPMComputer,
		PrimeDownloadMoveUploadPDFGenerator: mockPrimeDownloadMoveUploadPDFGenerator,
		UserUploader:                        *userUploader,
	}
	// Set up service member ID to verify
	serviceMemberID, err := uuid.NewV4()
	suite.FatalNoError(err)
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: serviceMemberID,
			},
		},
	}, nil)

	// Ensure appcontext id is different than service member id
	differentID, err := uuid.NewV4()
	suite.FatalNoError(err)

	appCtx := suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: differentID,
	})
	ppmShipmentID := ppmShipment.ID
	suite.MustSave(&ppmShipment)

	err = a.VerifyAOAPacketInternal(appCtx, ppmShipmentID)
	suite.Error(err)

}

func (suite *PPMShipmentSuite) TestCreateAOAPacketNotFound() {
	mockSSWPPMGenerator := &mocks.SSWPPMGenerator{}
	mockSSWPPMComputer := &mocks.SSWPPMComputer{}
	mockPrimeDownloadMoveUploadPDFGenerator := &mocks.PrimeDownloadMoveUploadPDFGenerator{}
	// mockAOAPacketCreator := &mocks.AOAPacketCreator{}
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)

	suite.Run("returns an error if the FetchDataShipmentSummaryWorksheet returns an error", func() {

		appCtx := suite.AppContextForTest()

		ppmshipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(suite.DB(), userUploader)

		ppmShipmentID := ppmshipment.ID
		errMsgPrefix := "error creating AOA packet"

		// Create an instance of aoaPacketCreator with mock dependencies
		a := &aoaPacketCreator{
			SSWPPMGenerator:                     mockSSWPPMGenerator,
			SSWPPMComputer:                      mockSSWPPMComputer,
			PrimeDownloadMoveUploadPDFGenerator: mockPrimeDownloadMoveUploadPDFGenerator,
			UserUploader:                        *userUploader,
		}
		fakeErr := apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		fakeErrWithWrap := fmt.Errorf("%s: %w", errMsgPrefix, fakeErr)

		// Define mock behavior for FetchDataShipmentSummaryWorksheetFormData
		mockSSWPPMComputer.On("FetchDataShipmentSummaryWorksheetFormData", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("*auth.Session"), mock.AnythingOfType("uuid.UUID")).Return(nil, fakeErr)

		// Test case: returns an error if FetchDataShipmentSummaryWorksheetFormData returns an error
		packet, err := a.CreateAOAPacket(appCtx, ppmShipmentID, false)
		suite.Error(err, err)
		suite.Equal(fakeErrWithWrap, err)
		if packet != nil {
			println("packet exists")
		}
	})

}

func (suite *PPMShipmentSuite) TestCreateAOAPacketFull() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)

	generator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.FatalNil(err)

	document := factory.BuildDocument(suite.DB(), nil, nil)
	file, err := suite.openLocalFile("../../paperwork/testdata/orders1.pdf", generator.FileSystem())
	suite.FatalNil(err)
	if generator != nil {
		suite.FatalNil(err)
	}

	mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
	SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(mockPPMCloseoutFetcher)
	ppmGenerator, err := shipmentsummaryworksheet.NewSSWPPMGenerator(generator)
	suite.FatalNoError(err)

	downloadMoveUploadGenerator, err := paperwork.NewMoveUserUploadToPDFDownloader(generator)
	suite.FatalNoError(err)
	order := factory.BuildOrder(suite.DB(), nil, nil)

	_, _, err = userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	err = suite.DB().Load(&document, "UserUploads.Upload")
	suite.FatalNil(err)
	suite.Equal(1, len(document.UserUploads))

	order.UploadedOrders = document
	order.UploadedOrdersID = document.ID
	suite.MustSave(&order)

	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	grade := models.ServiceMemberGradeE9
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
				Grade:      &grade,
			},
		},
		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    yuma,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.SignedCertification{},
		},
	}, nil)
	ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrders = document
	ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrdersID = document.ID
	ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedAmendedOrders = &document
	ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedAmendedOrdersID = &document.ID

	ppmShipmentID := ppmShipment.ID
	suite.MustSave(&ppmShipment)
	appCtx := suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
	})

	// Create an instance of aoaPacketCreator with mock dependencies
	a := &aoaPacketCreator{
		SSWPPMGenerator:                     ppmGenerator,
		SSWPPMComputer:                      SSWPPMComputer,
		PrimeDownloadMoveUploadPDFGenerator: downloadMoveUploadGenerator,
		UserUploader:                        *userUploader,
		pdfGenerator:                        generator,
	}

	_, err = models.SaveMoveDependencies(suite.DB(), &ppmShipment.Shipment.MoveTaskOrder)
	suite.NoError(err)

	packet, err := a.CreateAOAPacket(appCtx, ppmShipmentID, false)
	suite.NoError(err)
	suite.NotNil(packet) // ensures was generated with temp filesystem
}

func (suite *PPMShipmentSuite) TestSaveAOAPacket() {
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
				appCtx := suite.AppContextForTest()

				ppmShipment := factory.BuildPPMShipment(nil, []factory.Customization{
					{
						Model: models.MTOShipment{
							ID: uuid.Must(uuid.NewV4()),
						},
					},
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
				}, []factory.Trait{factory.GetTraitApprovedPPMShipment})

				mockMergedPDF := factory.FixtureOpen("aoa-packet.pdf")

				defer mockMergedPDF.Close()

				err := saveAOAPacket(appCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

				if suite.Error(err) {
					suite.ErrorContains(err, "failed to create AOA packet document")

					suite.ErrorContains(err, testCase.expectedErrMsg)
				}
			})
		}
	})

	suite.Run("returns an error if we fail to update the PPM shipment", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})
		suite.FatalNil(ppmShipment.AOAPacketID)

		mockMergedPDF := factory.FixtureOpen("aoa-packet.pdf")

		defer mockMergedPDF.Close()

		fakeError := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since the
		// saveAOAPacket function will be using a transaction.
		suite.NoError(appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			setUpMockPPMShipmentUpdater(
				mockPPMShipmentUpdater,
				txnAppCtx,
				&ppmShipment,
				nil,
				fakeError,
			)

			err := saveAOAPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

			if suite.Error(err) {
				suite.ErrorIs(err, fakeError)

				suite.ErrorContains(err, "failed to update PPMShipment with AOA packet document")
			}

			return nil
		}))
	})

	suite.Run("returns an error if we fail to prepare the file for upload", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})

		suite.FatalNil(ppmShipment.AOAPacketID)

		mockMergedPDF := factory.FixtureOpen("aoa-packet.pdf")

		// If the file is closed, it should trigger an error when we try to prepare it for upload.
		mockMergedPDF.Close()

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since the
		// saveAOAPacket function will be using a transaction.
		suite.NoError(appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			setUpMockPPMShipmentUpdater(
				mockPPMShipmentUpdater,
				txnAppCtx,
				&ppmShipment,
				&ppmShipment,
				nil,
			)

			err := saveAOAPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

			if suite.Error(err) {
				suite.ErrorContains(err, "failed to prepare AOA packet for upload")
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
				appCtx := suite.AppContextForTest()

				ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})

				suite.FatalNil(ppmShipment.AOAPacketID)

				mockMergedPDF := factory.FixtureOpen("aoa-packet.pdf")

				defer mockMergedPDF.Close()

				// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since
				// the saveAOAPacket function will be using a transaction.
				txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
					setUpMockPPMShipmentUpdater(
						mockPPMShipmentUpdater,
						txnAppCtx,
						&ppmShipment,
						&ppmShipment,
						nil,
					)

					// Setting a bad ID on this should trigger the upload to fail
					ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = testCase.userID

					err := saveAOAPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

					if suite.Error(err) {
						suite.ErrorContains(err, "failed to upload AOA packet")

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

		ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})

		suite.FatalNil(ppmShipment.AOAPacketID)

		mockMergedPDF := factory.FixtureOpen("aoa-packet.pdf")

		defer mockMergedPDF.Close()

		expectedBytes, readExpectedErr := io.ReadAll(mockMergedPDF)

		suite.FatalNoError(readExpectedErr)

		_, seekErr := mockMergedPDF.Seek(0, io.SeekStart)

		suite.FatalNoError(seekErr)

		// need to start a transaction so that our mocks know what the appCtx will actually be pointing to since
		// the saveAOAPacket function will be using a transaction.
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

			err := saveAOAPacket(txnAppCtx, &ppmShipment, mockMergedPDF, mockPPMShipmentUpdater, userUploader)

			suite.NoError(err)

			return nil
		}))

		// Now we'll double check everything to make sure it was saved correctly.
		if suite.NotNil(ppmShipment.AOAPacketID) {
			download, downloadErr := userUploader.Download(appCtx, &ppmShipment.AOAPacket.UserUploads[0])

			suite.FatalNoError(downloadErr)

			actualBytes, readActualErr := io.ReadAll(download)

			suite.FatalNoError(readActualErr)

			suite.Equal(expectedBytes, actualBytes)
		}
	})
}

func (suite *PPMShipmentSuite) closeFile(file afero.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

func (suite *PPMShipmentSuite) openLocalFile(path string, fs *afero.Afero) (afero.File, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	outputFile, err := fs.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating afero file")
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		return nil, errors.Wrap(err, "error copying over file contents")
	}

	suite.closeFile(outputFile)

	return outputFile, nil
}
