package ppmshipment

import (
	// 	"fmt"
	// 	"io"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/afero"
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
				EagerPreloadAssociationMovingExpenses,
			},
			[]string{
				PostLoadAssociationWeightTicketUploads,
				PostLoadAssociationProgearWeightTicketUploads,
				PostLoadAssociationMovingExpenseUploads,
			},
			returnValue...,
		)
	}

	file, err := suite.openLocalFile("../../paperwork/testdata/orders1.pdf", generator.FileSystem())
	suite.FatalNil(err)
	mockAoaPacketCreator.On("CreateAOAPacket", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID")).Return(file, nil)

	suite.Run("generate pdf - INTERNAL", func() {

		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
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

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 10  (9 + 1)
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
		// just manually set the type to OIL
		var movingExpenseReceiptTypeOil models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeOil
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeOil

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
		// just manually set the type to PACKING_MATERIALS
		var movingExpenseReceiptTypePackingMaterials models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypePackingMaterials
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypePackingMaterials

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
		// just manually set the type to RENTAL_EQUIPMENT
		var movingExpenseReceiptTypeRentalEquipment models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeRentalEquipment
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeRentalEquipment

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
		// just manually set the type to STORAGE
		var movingExpenseReceiptTypeStorage models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeStorage
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeStorage

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
		// just manually set the type to TOLLS
		var movingExpenseReceiptTypeTolls models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeTolls
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeTolls

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
		// just manually set the type to WEIGHING_FEE
		var movingExpenseReceiptTypeWeighingFee models.MovingExpenseReceiptType = models.MovingExpenseReceiptTypeWeighingFee
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeWeighingFee

		// append 1 expense
		// it will have:
		// 1 expense with 1 doc
		// updated: total = 17  (16 + 1)
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
		ppmShipment.MovingExpenses[len(ppmShipment.MovingExpenses)-1].MovingExpenseType = &movingExpenseReceiptTypeOther

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

		pdf, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)
		suite.FatalNil(err)

		pdfBookmarks := extractBookmarks(suite, *generator, pdf)
		suite.True(len(pdfBookmarks.Bookmarks) == 19)

		// this is a way to verify doc order via bookmark title
		expectedLabels := [19]string{"Shipment Summary Worksheet and Orders", "Weight Moved: Trip #1: Empty Weight Document #1", "Weight Moved: Trip #1: Empty Weight Document #2",
			"Weight Moved: Trip #1: Full Weight Document #1", "Weight Moved: Trip #1: Full Weight Document #2", "Weight Moved: Trip #2: Empty Weight Document #1",
			"Weight Moved: Trip #2: Full Weight Document #1", "Pro-gear: Set #1:(Me) Weight Ticket Document #1", "Weight Moved: Trip #1: POV/Registration Document #1",
			"Weight Moved: Trip #2: POV/Registration Document #1", "Expenses: Receipt #1: Weighing Fee Document #1", "Expenses: Receipt #2: Rental Equipment Document #1",
			"Expenses: Receipt #3: Contracted Expense Document #1", "Expenses: Receipt #4: Oil Document #1", "Expenses: Receipt #5: Packing Materials Document #1", "Expenses: Receipt #6: Tolls Document #1",
			"Expenses: Receipt #7: Storage Document #1", "Expenses: Receipt #8: Other Document #1", "Expenses: Receipt #9: Other Document #1"}

		for i := 0; i < len(pdfBookmarks.Bookmarks); i++ {
			suite.Equal(expectedLabels[i], pdfBookmarks.Bookmarks[i].Title)
		}
	})

	suite.Run("returns a NotFoundError if the ppmShipment is not found", func() {
		appCtx := suite.AppContextForTest()

		ppmShipmentID := uuid.Must(uuid.NewV4())

		fakeErr := apperror.NewNotFoundError(ppmShipmentID, "PPMShipment")

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipmentID, nil, fakeErr)

		_, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipmentID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)
		}
	})

	suite.Run("returns a ForbiddenError if the ppmShipment does not belong to user in INTERNAL context", func() {
		serviceMemberID := uuid.Must(uuid.NewV4())
		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: serviceMemberID,
			ApplicationName: auth.MilApp,
		})

		fakeErr := apperror.NewForbiddenError(fmt.Sprintf("PPMShipmentId: %s", ppmShipment.ID.String()))

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, nil, fakeErr)

		_, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)
		}
	})

	suite.Run("generation even if PPM is not current user's - NON INTERNAL CONTEXT (Office/Admin)", func() {
		var apps = []auth.Application{
			auth.OfficeApp,
			auth.AdminApp,
		}

		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
		suite.NotNil(ppmShipment)

		notOwnerServiceMemberID := uuid.Must(uuid.NewV4())
		for _, app := range apps {
			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ServiceMemberID: notOwnerServiceMemberID,
				ApplicationName: app,
			})

			setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

			// should still generate even if PPM belongs to different user in office/admin
			pdf, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)
			suite.FatalNil(err)

			mergedBytes, err := io.ReadAll(pdf)
			suite.FatalNil(err)
			suite.True(len(mergedBytes) > 0)
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
		pdf, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipment.ID)
		suite.FatalNil(err)

		mergedBytes, err := io.ReadAll(pdf)
		suite.FatalNil(err)
		suite.True(len(mergedBytes) > 0)
	})

	suite.Run("generate with disabled bookmark and watermark", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
		suite.NotNil(ppmShipment)
		notOwnerServiceMemberID := uuid.Must(uuid.NewV4())
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: notOwnerServiceMemberID,
			ApplicationName: auth.OfficeApp,
		})

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

		// disable bookmark and watermark
		// TODO -- figure out how to determine if watermark was generated
		pdf, err := paymentPacketCreator.Generate(appCtx, ppmShipment.ID, false, false)
		suite.FatalNil(err)

		mergedBytes, err := io.ReadAll(pdf)
		suite.FatalNil(err)
		suite.True(len(mergedBytes) > 0)
	})

	suite.Run("generate with enable bookmark, disable watermark", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
		suite.NotNil(ppmShipment)
		notOwnerServiceMemberID := uuid.Must(uuid.NewV4())
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: notOwnerServiceMemberID,
			ApplicationName: auth.OfficeApp,
		})

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

		// enable bookmark, disable watermark
		pdf, err := paymentPacketCreator.Generate(appCtx, ppmShipment.ID, true, false)
		suite.FatalNil(err)

		bookmarks := extractBookmarks(suite, *generator, pdf)
		suite.True(len(bookmarks.Bookmarks) > 0)
	})

	suite.Run("generate with disable bookmark, enable watermark", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
		suite.NotNil(ppmShipment)
		notOwnerServiceMemberID := uuid.Must(uuid.NewV4())
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: notOwnerServiceMemberID,
			ApplicationName: auth.OfficeApp,
		})

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipment.ID, &ppmShipment, nil)

		// disable bookmark, enable watermark
		// TODO -- figure out how to determine if watermark was generated
		pdf, err := paymentPacketCreator.Generate(appCtx, ppmShipment.ID, false, true)
		suite.FatalNil(err)

		bookmarks := extractBookmarks(suite, *generator, pdf)
		suite.True(bookmarks == nil)
	})
}

func extractBookmarks(suite *PPMShipmentSuite, generator paperworkgenerator.Generator, pdf io.ReadCloser) *pdfBookmarks {
	mergedBytes, err := io.ReadAll(pdf)
	suite.FatalNil(err)
	suite.True(len(mergedBytes) > 0)

	memorybasedFs := afero.NewMemMapFs()

	outFile, err := memorybasedFs.Create("test")
	suite.FatalNil(err)
	defer outFile.Close()

	buf := new(bytes.Buffer)
	buf.Write(mergedBytes)

	_, err = io.Copy(outFile, buf)
	suite.FatalNil(err)

	info, err := generator.GetPdfFileInfoForReadSeeker(outFile)
	suite.FatalNil(err)
	suite.True(info.PageCount > 0)

	buf = new(bytes.Buffer)
	err = api.ExportBookmarksJSON(outFile, buf, "", nil)
	if err != nil {
		// no bookmarks
		return nil
	}

	pb := pdfBookmarks{}

	err = json.Unmarshal(buf.Bytes(), &pb)
	suite.FatalNil(err)

	return &pb
}

type pdfBookmarks struct {
	Bookmarks []bookmarks
}
type bookmarks struct {
	Title string `json:"title"`
	Page  int    `json:"page"`
}
