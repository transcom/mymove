package ppmshipment

import (
	"fmt"
	"io"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
)

var sortedExpenseType = [8]string{"WEIGHING_FEE", "RENTAL_EQUIPMENT", "CONTRACTED_EXPENSE", "OIL", "PACKING_MATERIALS", "TOLLS", "STORAGE", "OTHER"}

type paymentPacketItem struct {
	Label    string
	Upload   models.Upload
	PageSize int
}

// NewFileInfo creates a new FileInfo struct.
func newPaymentPacketItem(label string, userUpload models.Upload) paymentPacketItem {
	return paymentPacketItem{
		Label:    label,
		Upload:   userUpload,
		PageSize: 0,
	}
}

// paymentPacketCreator is the concrete struct implementing the PaymentPacketCreator interface
type paymentPacketCreator struct {
	services.PPMShipmentFetcher
	pdfGenerator     paperwork.Generator
	aoaPacketCreator services.AOAPacketCreator
}

// NewPaymentPacketCreator creates a new PaymentPacketCreator with all of its dependencies
func NewPaymentPacketCreator(
	ppmShipmentFetcher services.PPMShipmentFetcher,
	pdfGenerator *paperwork.Generator,
	aoaPacketCreator services.AOAPacketCreator,
) services.PaymentPacketCreator {
	return &paymentPacketCreator{
		ppmShipmentFetcher,
		*pdfGenerator,
		aoaPacketCreator,
	}
}

func (p *paymentPacketCreator) Generate(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, addBookmarks bool, addWatermarks bool) (io.ReadCloser, error) {

	err := verifyPPMShipment(appCtx, ppmShipmentID)
	if err != nil {
		return nil, err
	}

	errMsgPrefix := "error creating payment packet"

	ppmShipment, err := p.PPMShipmentFetcher.GetPPMShipment(
		appCtx,
		ppmShipmentID,
		// Note: Orders will be generate via SSW creator service
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
	)

	// something bad happened on data retrieval of everything for PPM
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to load PPMShipment")
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, err
	}

	var pdfFilesToMerge []io.ReadSeeker

	// use aoa creator to generated SSW and Orders PDF
	aoaPacketFile, err := p.aoaPacketCreator.CreateAOAPacket(appCtx, ppmShipmentID, true)
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, fmt.Sprintf("failed to generate AOA packet for ppmShipmentID: %s", ppmShipmentID.String()))
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	// AOA packet will be appended at the beginning of the final pdf file
	pdfFilesToMerge = append(pdfFilesToMerge, aoaPacketFile)

	// Start building individual PDFs for each expense/receipt docs. These files will then be merged as one PDF.
	var pdfFileNamesToMerge []string
	sortedPaymentPacketItemsMap := buildPaymentPacketItemsMap(ppmShipment)

	for i := 0; i < len(sortedPaymentPacketItemsMap); i++ {
		pdfFileName, perr := p.pdfGenerator.ConvertUploadToPDF(appCtx, sortedPaymentPacketItemsMap[i].Upload)
		if perr != nil {
			errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to generate pdf for upload")
			appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
			return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
		}
		pdfFileNamesToMerge = append(pdfFileNamesToMerge, pdfFileName)
	}

	if len(pdfFileNamesToMerge) > 0 {
		pdfFileNamesToMergePdf, perr := p.pdfGenerator.MergePDFFiles(appCtx, pdfFileNamesToMerge)
		if perr != nil {
			errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed pdfGenerator.MergePDFFiles")
			appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
			return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
		}
		pdfFilesToMerge = append(pdfFilesToMerge, pdfFileNamesToMergePdf)
	}

	// Do final merge of all PDFs into one.
	finalMergePdf, err := p.pdfGenerator.MergePDFFilesByContents(appCtx, pdfFilesToMerge)
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to generated file merged pdf")
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	// Start building bookmarks and watermarks
	bookmarks, err := buildBookMarks(pdfFileNamesToMerge, sortedPaymentPacketItemsMap, aoaPacketFile, p.pdfGenerator)
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to generate bookmarks for PDF")
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	watermarks, err := buildWaterMarks(bookmarks, p.pdfGenerator)
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to generate watermarks for PDF")
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	// Apply bookmarks and watermarks based on flag
	if addWatermarks && len(watermarks) > 0 {
		pdfWithWatermarks, err := p.pdfGenerator.AddWatermarks(finalMergePdf, watermarks)
		if err != nil {
			errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to add watermarks to PDF")
			appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
			return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
		}
		if addBookmarks {
			return p.pdfGenerator.AddPdfBookmarks(pdfWithWatermarks, bookmarks)
		}
		return pdfWithWatermarks, nil
	}

	if addBookmarks {
		return p.pdfGenerator.AddPdfBookmarks(finalMergePdf, bookmarks)
	}

	// bookmark and watermark both disabled
	return finalMergePdf, nil
}

func (p *paymentPacketCreator) GenerateDefault(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (io.ReadCloser, error) {
	return p.Generate(appCtx, ppmShipmentID, true, true)
}

func buildBookMarks(fileNamesToMerge []string, sortedPaymentPacketItems map[int]paymentPacketItem, aoaPacketFile io.ReadSeeker, pdfGenerator paperwork.Generator) ([]pdfcpu.Bookmark, error) {
	// go out and retrieve PDF file info for each file name
	for i := 0; i < len(fileNamesToMerge); i++ {
		pdfFileInfo, err := pdfGenerator.GetPdfFileInfo(fileNamesToMerge[i])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fmt.Sprintf("failed to retrieve PDF file info for: %s", fileNamesToMerge[i]), err)
		}
		item := sortedPaymentPacketItems[i]
		// we just want the pagesize. update sortedPaymentPacketItems
		item.PageSize = pdfFileInfo.PageCount
		sortedPaymentPacketItems[i] = item
	}

	var bookmarks []pdfcpu.Bookmark

	// retrieve file info for AOA packet file
	aoaPacketFileInfo, err := pdfGenerator.GetPdfFileInfoForReadSeeker(aoaPacketFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "failed to retrieve PDF file info for AOA packet file", err)
	}
	// add first bookmark for AOA content
	bookmarks = append(bookmarks, pdfcpu.Bookmark{PageFrom: 1, PageThru: aoaPacketFileInfo.PageCount, Title: "Shipment Summary Worksheet and Orders"})

	// build bookmarks for all file names
	var pageFrom int
	var pageThru int
	for i := 0; i < len(fileNamesToMerge); i++ {
		item := sortedPaymentPacketItems[i]
		pageFrom = bookmarks[i].PageThru + 1
		pageThru = bookmarks[i].PageThru + item.PageSize
		bookmarks = append(bookmarks, pdfcpu.Bookmark{PageFrom: pageFrom, PageThru: pageThru, Title: item.Label})
	}
	return bookmarks, nil
}

// generate watermarks which will serve as page footer labels
func buildWaterMarks(bookMarks []pdfcpu.Bookmark, pdfGenerator paperwork.Generator) (map[int][]*model.Watermark, error) {
	m := make(map[int][]*model.Watermark)

	opacity := 1.0
	onTop := true
	update := false
	unit := types.POINTS

	desc := fmt.Sprintf("font:Times-Italic, points:10, sc:1 abs, pos:bc, off:0 8, rot:0, op:%f", opacity)

	creationTimeStamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	totalPages := bookMarks[len(bookMarks)-1].PageThru
	currentPage := 1
	bookMarkIndex := 0
	for _, bm := range bookMarks {
		cnt := bm.PageThru - bm.PageFrom
		for j := 0; j <= cnt; j++ {
			// do not add watermark on SSW pages
			if currentPage < 4 {
				currentPage++
				continue
			}
			wmText := bm.Title
			// we really can't use the bookmark title for the SSW+Orders.
			// we will just label it as only Orders
			if currentPage > 3 && bookMarkIndex == 0 {
				wmText = "Orders"
			}
			wms := make([]*model.Watermark, 0)
			pagingInfo := fmt.Sprintf("Page %d of %d", currentPage, totalPages)
			text := fmt.Sprintf("%s - Payment Packet[%s] (Creation Date: %v)", pagingInfo, wmText, creationTimeStamp)

			wm, _ := pdfGenerator.CreateTextWatermark(text, desc, onTop, update, unit)
			wms = append(wms, wm)
			// note: use current page because map is 1 based
			m[currentPage] = wms
			currentPage++
		}
		bookMarkIndex++
	}

	return m, nil
}

func buildPaymentPacketItemsMap(ppmShipment *models.PPMShipment) map[int]paymentPacketItem {
	// items are sorted based on key(int), key represents order index
	sortedPaymentPacketItems := make(map[int]paymentPacketItem)
	sortedPaymentPacketItemsIndex := 0

	var povRegistrationPaymentPacketItems []paymentPacketItem

	// process weight tickets
	weightTicketCnt := 1
	for _, weightTicket := range ppmShipment.WeightTickets {
		sectionLabel := fmt.Sprintf("Weight Moved: Trip #%s:", fmt.Sprint(weightTicketCnt))
		emptyWeightCnt := 1
		for _, uu := range weightTicket.EmptyDocument.UserUploads {
			sortedPaymentPacketItems[sortedPaymentPacketItemsIndex] = newPaymentPacketItem(fmt.Sprintf("%s Empty Weight Document #%s", sectionLabel, fmt.Sprint(emptyWeightCnt)), uu.Upload)
			sortedPaymentPacketItemsIndex++
			emptyWeightCnt++
		}
		fullWeightCnt := 1
		for _, uu := range weightTicket.FullDocument.UserUploads {
			sortedPaymentPacketItems[sortedPaymentPacketItemsIndex] = newPaymentPacketItem(fmt.Sprintf("%s Full Weight Document #%s", sectionLabel, fmt.Sprint(fullWeightCnt)), uu.Upload)
			sortedPaymentPacketItemsIndex++
			fullWeightCnt++
		}

		// povRegistrations will be order after the other weight tickets
		povCnt := 1
		for _, uu := range weightTicket.ProofOfTrailerOwnershipDocument.UserUploads {
			// stuff into array so we may render it later in the correct order
			povRegistrationPaymentPacketItems = append(povRegistrationPaymentPacketItems, newPaymentPacketItem(fmt.Sprintf("%s POV/Registration Document #%s", sectionLabel, fmt.Sprint(povCnt)), uu.Upload))
			povCnt++
		}

		weightTicketCnt++
	}

	// process pro-gear
	proGearWeightTicketCnt := 1
	for _, progearWeightTicket := range ppmShipment.ProgearWeightTickets {
		sectionLabel := fmt.Sprintf("Pro-gear: Set #%s:", fmt.Sprint(proGearWeightTicketCnt))
		proGearWeightTicketSetCnt := 1
		for _, uu := range progearWeightTicket.Document.UserUploads {
			var belongsToLabel = "Me"
			if !*progearWeightTicket.BelongsToSelf {
				belongsToLabel = "Spouse"
			}
			sortedPaymentPacketItems[sortedPaymentPacketItemsIndex] = newPaymentPacketItem(fmt.Sprintf("%s(%s) Weight Ticket Document #%s", sectionLabel, belongsToLabel, fmt.Sprint(proGearWeightTicketSetCnt)), uu.Upload)
			sortedPaymentPacketItemsIndex++
			proGearWeightTicketSetCnt++
		}
		proGearWeightTicketCnt++
	}

	// place povRegistration after the weight items
	for _, item := range povRegistrationPaymentPacketItems {
		sortedPaymentPacketItems[sortedPaymentPacketItemsIndex] = item
		sortedPaymentPacketItemsIndex++
	}

	// process expenses, group by type as array list
	expenseTypeMap := make(map[string][]models.MovingExpense)
	for _, movingExpense := range ppmShipment.MovingExpenses {
		if value, exists := expenseTypeMap[string(*movingExpense.MovingExpenseType)]; exists {
			// add to existing array
			value = append(value, movingExpense)
			expenseTypeMap[string(*movingExpense.MovingExpenseType)] = value
		} else {
			// create new array and add first item
			expenseTypeMap[string(*movingExpense.MovingExpenseType)] = make([]models.MovingExpense, 0)
			expenseTypeMap[string(*movingExpense.MovingExpenseType)] = append(expenseTypeMap[string(*movingExpense.MovingExpenseType)], movingExpense)
		}
	}

	expensesCnt := 1
	for _, expenseType := range sortedExpenseType {
		for _, item := range expenseTypeMap[expenseType] {
			expensesDocCnt := 1
			for _, uu := range item.Document.UserUploads {
				sectionLabel := fmt.Sprintf("Expenses: Receipt #%s:", fmt.Sprint(expensesCnt))
				sortedPaymentPacketItems[sortedPaymentPacketItemsIndex] = newPaymentPacketItem(fmt.Sprintf("%s %s Document #%s", sectionLabel, getExpenseTypeLabel(expenseType), fmt.Sprint(expensesDocCnt)), uu.Upload)
				sortedPaymentPacketItemsIndex++
				expensesDocCnt++
			}
			expensesCnt++
		}
	}

	return sortedPaymentPacketItems
}

// Helper method to verify if PPMShipment is accessible for current user. This is to prevent
// backdoor access for unauthorized users in INTERNAL.
func verifyPPMShipment(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) error {
	ppmShipment := models.PPMShipment{}
	dbQErr := appCtx.DB().Q().Eager(
		"Shipment.MoveTaskOrder.Orders.ServiceMember",
	).Find(&ppmShipment, ppmShipmentID)

	if dbQErr != nil {
		if errors.Cause(dbQErr).Error() == models.RecordNotFoundErrorString {
			return apperror.NewNotFoundError(ppmShipmentID, "PPMShipment")
		}
		return dbQErr
	}

	// if request is from INTERNAL verify if PPM belongs to user
	if appCtx.Session().IsMilApp() && ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID != appCtx.Session().ServiceMemberID {
		return apperror.NewForbiddenError(fmt.Sprintf("PPMShipmentId: %s", ppmShipmentID.String()))
	}

	return nil
}

func getExpenseTypeLabel(value string) string {
	switch value {
	case "CONTRACTED_EXPENSE":
		return "Contracted Expense"
	case "OIL":
		return "Oil"
	case "PACKING_MATERIALS":
		return "Packing Materials"
	case "RENTAL_EQUIPMENT":
		return "Rental Equipment"
	case "STORAGE":
		return "Storage"
	case "TOLLS":
		return "Tolls"
	case "WEIGHING_FEE":
		return "Weighing Fee"
	case "OTHER":
		return "Other"
	default:
		return "Unknown"
	}
}
