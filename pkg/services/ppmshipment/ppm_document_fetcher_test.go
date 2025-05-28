package ppmshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestPPMDocumentFetcher() {
	// Set up a fake S3 storage that we'll use to manage uploads
	fakeS3 := storageTest.NewFakeS3Storage(true)

	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)

	suite.FatalNoError(uploaderErr)

	ppmDocumentFetcher := NewPPMDocumentFetcher()

	makePPMShipmentWithAllDocuments := func(appCtx appcontext.AppContext) *models.PPMShipment {
		// Set up PPM shipment that is at the correct stage of processing for when we would typically use this service
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)

		suite.NotNil(ppmShipment)

		// Add an extra upload to one of the weight ticket documents to verify we get all non-deleted uploads later
		ppmShipment.WeightTickets[0].EmptyDocument.UserUploads = append(
			ppmShipment.WeightTickets[0].EmptyDocument.UserUploads,
			factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.WeightTickets[0].EmptyDocument,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		// The PPM shipment generated above only has a weight ticket, but we want to ensure the service works with all
		// types of documents, so we'll add some more here.
		ppmShipment.MovingExpenses = append(
			ppmShipment.MovingExpenses,
			factory.BuildMovingExpense(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		// Add an extra upload to the moving expense document to verify we get all non-deleted uploads later
		ppmShipment.MovingExpenses[0].Document.UserUploads = append(
			ppmShipment.MovingExpenses[0].Document.UserUploads,
			factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.MovingExpenses[0].Document,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		ppmShipment.ProgearWeightTickets = append(
			ppmShipment.ProgearWeightTickets,
			factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		// Add an extra upload to the progear weight ticket document to verify we get all non-deleted uploads later
		ppmShipment.ProgearWeightTickets[0].Document.UserUploads = append(
			ppmShipment.ProgearWeightTickets[0].Document.UserUploads,
			factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.ProgearWeightTickets[0].Document,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		ppmShipment.GunSafeWeightTickets = append(
			ppmShipment.GunSafeWeightTickets,
			factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		// Add an extra upload to the gun safe weight ticket document to verify we get all non-deleted uploads later
		ppmShipment.GunSafeWeightTickets[0].Document.UserUploads = append(
			ppmShipment.GunSafeWeightTickets[0].Document.UserUploads,
			factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.GunSafeWeightTickets[0].Document,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)

		suite.FatalTrue(len(ppmShipment.WeightTickets) > 0)
		suite.FatalTrue(len(ppmShipment.WeightTickets[0].EmptyDocument.UserUploads) > 1)
		suite.FatalTrue(len(ppmShipment.MovingExpenses) > 0)
		suite.FatalTrue(len(ppmShipment.MovingExpenses[0].Document.UserUploads) > 1)
		suite.FatalTrue(len(ppmShipment.ProgearWeightTickets) > 0)
		suite.FatalTrue(len(ppmShipment.ProgearWeightTickets[0].Document.UserUploads) > 1)
		suite.FatalTrue(len(ppmShipment.GunSafeWeightTickets) > 0)
		suite.FatalTrue(len(ppmShipment.GunSafeWeightTickets[0].Document.UserUploads) > 1)

		return &ppmShipment
	}

	suite.Run("Can retrieve weight tickets for a PPM shipment", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := makePPMShipmentWithAllDocuments(appCtx)

		fetchedDocument, err := ppmDocumentFetcher.GetPPMDocuments(appCtx, ppmShipment.Shipment.ID)

		if suite.NoError(err) && suite.NotNil(fetchedDocument) {
			suite.Equal(len(ppmShipment.WeightTickets), len(fetchedDocument.WeightTickets))

			for i := range ppmShipment.WeightTickets {
				suite.Equal(ppmShipment.WeightTickets[i].ID, fetchedDocument.WeightTickets[i].ID)
				suite.Equal(len(ppmShipment.WeightTickets[i].EmptyDocument.UserUploads), len(fetchedDocument.WeightTickets[i].EmptyDocument.UserUploads))

				for j := range ppmShipment.WeightTickets[i].EmptyDocument.UserUploads {
					suite.Equal(ppmShipment.WeightTickets[i].EmptyDocument.UserUploads[j].ID, fetchedDocument.WeightTickets[i].EmptyDocument.UserUploads[j].ID)
				}
			}

			suite.Equal(len(ppmShipment.MovingExpenses), len(fetchedDocument.MovingExpenses))

			for i := range ppmShipment.MovingExpenses {
				suite.Equal(ppmShipment.MovingExpenses[i].ID, fetchedDocument.MovingExpenses[i].ID)
				suite.Equal(len(ppmShipment.MovingExpenses[i].Document.UserUploads), len(fetchedDocument.MovingExpenses[i].Document.UserUploads))

				for j := range ppmShipment.MovingExpenses[i].Document.UserUploads {
					suite.Equal(ppmShipment.MovingExpenses[i].Document.UserUploads[j].ID, fetchedDocument.MovingExpenses[i].Document.UserUploads[j].ID)
				}
			}

			suite.Equal(len(ppmShipment.ProgearWeightTickets), len(fetchedDocument.ProgearWeightTickets))

			for i := range ppmShipment.ProgearWeightTickets {
				suite.Equal(ppmShipment.ProgearWeightTickets[i].ID, fetchedDocument.ProgearWeightTickets[i].ID)
				suite.Equal(len(ppmShipment.ProgearWeightTickets[i].Document.UserUploads), len(fetchedDocument.ProgearWeightTickets[i].Document.UserUploads))

				for j := range ppmShipment.ProgearWeightTickets[i].Document.UserUploads {
					suite.Equal(ppmShipment.ProgearWeightTickets[i].Document.UserUploads[j].ID, fetchedDocument.ProgearWeightTickets[i].Document.UserUploads[j].ID)
				}
			}

			suite.Equal(len(ppmShipment.GunSafeWeightTickets), len(fetchedDocument.GunSafeWeightTickets))

			for i := range ppmShipment.GunSafeWeightTickets {
				suite.Equal(ppmShipment.GunSafeWeightTickets[i].ID, fetchedDocument.GunSafeWeightTickets[i].ID)
				suite.Equal(len(ppmShipment.GunSafeWeightTickets[i].Document.UserUploads), len(fetchedDocument.GunSafeWeightTickets[i].Document.UserUploads))

				for j := range ppmShipment.GunSafeWeightTickets[i].Document.UserUploads {
					suite.Equal(ppmShipment.GunSafeWeightTickets[i].Document.UserUploads[j].ID, fetchedDocument.GunSafeWeightTickets[i].Document.UserUploads[j].ID)
				}
			}
		}
	})

	suite.Run("Returns empty slices if the shipment has been deleted", func() {
		appCtx := suite.AppContextForTest()

		// Not the proper way to delete a shipment, but we can't delete a shipment that has actually progressed to this
		// stage using our `utilities.SoftDestroy` function because of the related SignedCertification missing the
		// DeletedAt field, so we'll just set the PPMShipment.DeletedAt field directly.
		now := time.Now()
		ppmShipmentToDelete := factory.BuildPPMShipmentThatNeedsCloseout(appCtx.DB(), userUploader, []factory.Customization{
			{
				Model: models.PPMShipment{
					DeletedAt: &now,
				},
			},
		})

		fetchedDocument, err := ppmDocumentFetcher.GetPPMDocuments(suite.AppContextForTest(), ppmShipmentToDelete.Shipment.ID)

		if suite.NoError(err) && suite.NotNil(fetchedDocument) {
			suite.Equal(models.PPMDocuments{}, *fetchedDocument)
		}
	})

	suite.Run("Returns empty slices if the shipment ID is not found", func() {
		fetchedDocument, err := ppmDocumentFetcher.GetPPMDocuments(suite.AppContextForTest(), uuid.Must(uuid.NewV4()))

		if suite.NoError(err) && suite.NotNil(fetchedDocument) {
			suite.Equal(models.PPMDocuments{}, *fetchedDocument)
		}
	})

	suite.Run("Excludes deleted uploads", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := makePPMShipmentWithAllDocuments(appCtx)

		// Create an upload for a weight ticket that we then delete
		originalWeightTicket := ppmShipment.WeightTickets[0]
		numValidEmptyUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
		suite.FatalTrue(numValidEmptyUploads > 0)

		deletedWeightTicketUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    originalWeightTicket.EmptyDocument,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		err := userUploader.DeleteUserUpload(appCtx, &deletedWeightTicketUpload)

		suite.FatalNoError(err)

		suite.FatalNotNil(deletedWeightTicketUpload.Upload.DeletedAt)
		suite.FatalNotNil(deletedWeightTicketUpload.DeletedAt)

		// Create an upload for a progear weight ticket that we then delete
		originalProgearWeightTicket := ppmShipment.ProgearWeightTickets[0]
		numValidProgearWeightTicketUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
		suite.FatalTrue(numValidProgearWeightTicketUploads > 0)

		deletedProgearWeightTicketUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    originalProgearWeightTicket.Document,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		err = userUploader.DeleteUserUpload(appCtx, &deletedProgearWeightTicketUpload)

		suite.FatalNoError(err)

		suite.FatalNotNil(deletedProgearWeightTicketUpload.Upload.DeletedAt)
		suite.FatalNotNil(deletedProgearWeightTicketUpload.DeletedAt)

		// Create an upload for a gun safe weight ticket that we then delete
		originalGunSafeWeightTicket := ppmShipment.GunSafeWeightTickets[0]
		numValidGunSafeWeightTicketUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
		suite.FatalTrue(numValidGunSafeWeightTicketUploads > 0)

		deletedGunSafeWeightTicketUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    originalGunSafeWeightTicket.Document,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		err = userUploader.DeleteUserUpload(appCtx, &deletedGunSafeWeightTicketUpload)

		suite.FatalNoError(err)

		suite.FatalNotNil(deletedGunSafeWeightTicketUpload.Upload.DeletedAt)
		suite.FatalNotNil(deletedGunSafeWeightTicketUpload.DeletedAt)

		// Create an upload for a moving expense that we then delete
		originalMovingExpense := ppmShipment.MovingExpenses[0]
		numValidMovingExpenseUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
		suite.FatalTrue(numValidMovingExpenseUploads > 0)

		deletedMovingExpenseUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    originalMovingExpense.Document,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		err = userUploader.DeleteUserUpload(appCtx, &deletedMovingExpenseUpload)

		suite.FatalNoError(err)

		suite.FatalNotNil(deletedMovingExpenseUpload.Upload.DeletedAt)
		suite.FatalNotNil(deletedMovingExpenseUpload.DeletedAt)

		// now we're ready to test the fetcher
		fetchedDocument, err := ppmDocumentFetcher.GetPPMDocuments(suite.AppContextForTest(), ppmShipment.Shipment.ID)

		if suite.NoError(err) && suite.NotNil(fetchedDocument) {
			suite.Equal(len(ppmShipment.WeightTickets), len(fetchedDocument.WeightTickets))

			suite.Equal(originalWeightTicket.ID, fetchedDocument.WeightTickets[0].ID)
			retrievedWeightTicket := fetchedDocument.WeightTickets[0]

			if suite.Equal(numValidEmptyUploads, len(retrievedWeightTicket.EmptyDocument.UserUploads)) {
				for _, upload := range retrievedWeightTicket.EmptyDocument.UserUploads {
					suite.NotEqual(deletedWeightTicketUpload.ID, upload.ID)
					suite.Nil(upload.DeletedAt)
				}
			}

			suite.Equal(len(ppmShipment.MovingExpenses), len(fetchedDocument.MovingExpenses))

			suite.Equal(originalMovingExpense.ID, fetchedDocument.MovingExpenses[0].ID)
			retrievedMovingExpense := fetchedDocument.MovingExpenses[0]

			if suite.Equal(numValidMovingExpenseUploads, len(retrievedMovingExpense.Document.UserUploads)) {
				for _, upload := range retrievedMovingExpense.Document.UserUploads {
					suite.NotEqual(deletedMovingExpenseUpload.ID, upload.ID)
					suite.Nil(upload.DeletedAt)
				}
			}

			suite.Equal(len(ppmShipment.ProgearWeightTickets), len(fetchedDocument.ProgearWeightTickets))

			suite.Equal(originalProgearWeightTicket.ID, fetchedDocument.ProgearWeightTickets[0].ID)
			retrievedProgearWeightTicket := fetchedDocument.ProgearWeightTickets[0]

			if suite.Equal(numValidProgearWeightTicketUploads, len(retrievedProgearWeightTicket.Document.UserUploads)) {
				for _, upload := range retrievedProgearWeightTicket.Document.UserUploads {
					suite.NotEqual(deletedProgearWeightTicketUpload.ID, upload.ID)
					suite.Nil(upload.DeletedAt)
				}
			}

			suite.Equal(len(ppmShipment.GunSafeWeightTickets), len(fetchedDocument.GunSafeWeightTickets))

			suite.Equal(originalGunSafeWeightTicket.ID, fetchedDocument.GunSafeWeightTickets[0].ID)
			retrievedGunSafeWeightTicket := fetchedDocument.GunSafeWeightTickets[0]

			if suite.Equal(numValidGunSafeWeightTicketUploads, len(retrievedGunSafeWeightTicket.Document.UserUploads)) {
				for _, upload := range retrievedGunSafeWeightTicket.Document.UserUploads {
					suite.NotEqual(deletedGunSafeWeightTicketUpload.ID, upload.ID)
					suite.Nil(upload.DeletedAt)
				}
			}
		}
	})
}
