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

	suite.Run("generate pdf", func() {
		appCtx := suite.AppContextForTest()

		ppmShipment := factory.BuildPPMShipmentThatNeedsPaymentApproval(suite.DB(), userUploader, nil)
		// initial test data will have one trip(as trip #1) containing:
		// 1 empty weight with 1 doc
		// 1 full weight with 1 doc
		// 1 POV/Registration with 1 doc
		// total = 3

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

		info, err := generator.GetPdfFileInfoByReadSeeker(outFile)
		suite.FatalNil(err)
		suite.True(info.PageCount > 0)

		buf = new(bytes.Buffer)
		api.ExportBookmarksJSON(outFile, buf, "", nil)

		pdfBookmarks := pdfBookmarks{}

		err = json.Unmarshal(buf.Bytes(), &pdfBookmarks)
		if err != nil {
			fmt.Println(err)
			return
		}

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

	suite.Run("returns an error if the PPMShipmentFetcher returns an error", func() {
		appCtx := suite.AppContextForTest()

		ppmShipmentID := uuid.Must(uuid.NewV4())

		fakeErr := apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")

		setUpMockPPMShipmentFetcherForPayment(appCtx, ppmShipmentID, nil, fakeErr)

		_, err := paymentPacketCreator.GenerateDefault(appCtx, ppmShipmentID)

		if suite.Error(err) {
			suite.ErrorIs(err, fakeErr)

			suite.ErrorContains(err, "not found while looking for PPMShipment")
		}
	})
}

type pdfBookmarks struct {
	Bookmarks []bookmarks
}
type bookmarks struct {
	Title string `json:"title"`
	Page  int    `json:"page"`
}
