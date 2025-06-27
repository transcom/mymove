package ppmshipment

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services/mocks"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestCreatePaymentPacket() {
	//--------------------//--------------------//--------------------
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, _ := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)
	generator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.FatalNil(err)

	mockPPMShipmentFetcher := &mocks.PPMShipmentFetcher{}
	mockAoaPacketCreator := &mocks.AOAPacketCreator{}
	paymentPacketCreator := NewPaymentPacketCreator(
		mockPPMShipmentFetcher,
		generator,
		mockAoaPacketCreator,
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
				EagerPreloadAssociationGunSafeWeightTickets,
				EagerPreloadAssociationMovingExpenses,
			},
			[]string{
				PostLoadAssociationWeightTicketUploads,
				PostLoadAssociationProgearWeightTicketUploads,
				PostLoadAssociationGunSafeWeightTicketUploads,
				PostLoadAssociationMovingExpenseUploads,
			},
			returnValue...,
		)
	}

	file, err := suite.openLocalFile("../../paperwork/testdata/orders1.pdf", generator.FileSystem())
	suite.FatalNil(err)
	mockAoaPacketCreator.On("CreateAOAPacket", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("bool")).Return(file, "testDir", nil)

	suite.Run("generate pdf - INTERNAL", func() {

		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)
		// initial test data will have one trip(as trip #1) containing:
		// 1 empty weight with 1 doc
		// 1 full weight with 1 doc
		// 1 POV/Registration with 1 doc
		// total = 3

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID,
			ApplicationName: auth.MilApp,
		})

		suite.NotNil(ppmShipment)

		// append another empty weight document to trip #1
		// updated: total = 4  (3 + 1)
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

		// append another full weight document to trip #1
		// updated: total = 5  (4 + 1)
		ppmShipment.WeightTickets[0].FullDocument.UserUploads = append(
			ppmShipment.WeightTickets[0].FullDocument.UserUploads,
			factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.WeightTickets[0].FullDocument,
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

		// append a new weight ticket (trip #2)
		// it will have:
		// 1 empty weight with 1 doc
		// 1 full weight with 1 doc
		// 1 pov/registration with 1 doc
		// updated: total = 8  (5 + 3)
		ppmShipment.WeightTickets = append(
			ppmShipment.WeightTickets,
			factory.BuildWeightTicket(suite.DB(), []factory.Customization{
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

		// append 1 pro-gear set for ME
		// it will have:
		// 1 pro weight with 1 doc
		// updated: total = 9  (8 + 1)
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

		// append 1 gun safe set for ME
		// it will have:
		// 1 gun safe weight with 1 doc
		// updated: total = 10  (9 + 1)
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

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 11  (10 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to CONTRACTED_EXPENSE
		var movingExpenseReceiptTypeContractedExpense models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeContractedExpense
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeContractedExpense

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 12  (11 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to OIL
		var movingExpenseReceiptTypeOil models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeOil
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeOil

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 13  (12 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to PACKING_MATERIALS
		var movingExpenseReceiptTypePackingMaterials models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypePackingMaterials
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypePackingMaterials

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 14  (13 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to RENTAL_EQUIPMENT
		var movingExpenseReceiptTypeRentalEquipment models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeRentalEquipment
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeRentalEquipment

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 15  (14 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to STORAGE
		var movingExpenseReceiptTypeStorage models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeStorage
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeStorage

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 16  (15 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to TOLLS
		var movingExpenseReceiptTypeTolls models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeTolls
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeTolls

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 17  (6 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to WEIGHING_FEE
		var movingExpenseReceiptTypeWeighingFee models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeWeighingFee
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeWeighingFee

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 18  (17 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to OTHER
		var movingExpenseReceiptTypeOther models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeOther
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeOther

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 19  (18 + 1)
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
		// note: factory data is created as Packing Material type expense. we will
		// just manually set the type to OTHER
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeOther

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

		//nolint:staticcheck
		_, dirPath, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)
		suite.FatalNil(err)

		err = paymentPacketCreator.CleanupPaymentPacketDir(dirPath)
		suite.NoError(err)
	})

	suite.Run("returns a NotFoundError if the ppmShipment is not found", func() {
		appCtx := suite.AppContextForTest()

		ppmShipmentID := uuid.Must(uuid.NewV4())

		fakeErr := apperror.NewNotFoundError(ppmShipmentID, "PPMShipment")

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipmentID, nil, fakeErr)

		_, dirPath, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipmentID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)
		}

		err = paymentPacketCreator.CleanupPaymentPacketDir(dirPath)
		suite.NoError(err)
	})

	suite.Run("returns a ForbiddenError if the ppmShipment does not belong to user in INTERNAL context", func() {
		serviceMemberID := uuid.Must(uuid.NewV4())
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMemberID,
			ApplicationName: auth.MilApp,
		})

		fakeErr := apperror.NewForbiddenError(fmt.Sprintf("PPMShipmentId: %s", ppmShipment.ID.String()))

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, nil, fakeErr)

		_, dirPath, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)
		}

		err = paymentPacketCreator.CleanupPaymentPacketDir(dirPath)
		suite.NoError(err)
	})

	suite.Run("generation even if PPM is not current user's - NON INTERNAL CONTEXT (Office/Admin)", func() {
		var apps = []auth.Application{
			auth.OfficeApp,
			auth.AdminApp,
		}

		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)
		suite.NotNil(ppmShipment)

		notOwnerServiceMemberID := uuid.Must(uuid.NewV4())
		for _, app := range apps {
			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ServiceMemberID: notOwnerServiceMemberID,
				ApplicationName: app,
			})

			setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

			// should still generate even if PPM belongs to different user in office/admin
			pdf, dirPath, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)
			suite.FatalNil(err)

			mergedBytes, err := io.ReadAll(pdf)
			suite.FatalNil(err)
			suite.True(len(mergedBytes) > 0)
			err = paymentPacketCreator.CleanupPaymentPacketDir(dirPath)
			suite.NoError(err)
		}
	})

	suite.Run("should still generate PDF if PPM has no uploaded docs (orders, expenses/receipts)", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		suite.NotNil(ppmShipment)
		notOwnerServiceMemberID := uuid.Must(uuid.NewV4())
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: notOwnerServiceMemberID,
			ApplicationName: auth.OfficeApp,
		})

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

		// should still generate even if PPM belongs to different user in office/admin
		pdf, dirPath, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)
		suite.FatalNil(err)

		mergedBytes, err := io.ReadAll(pdf)
		suite.FatalNil(err)
		suite.True(len(mergedBytes) > 0)
		err = paymentPacketCreator.CleanupPaymentPacketDir(dirPath)
		suite.NoError(err)
	})
}
