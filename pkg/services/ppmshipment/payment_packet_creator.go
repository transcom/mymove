package ppmshipment

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
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

func (p *paymentPacketCreator) Generate(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, addWatermarks bool) (mergedPdf afero.File, dirPath string, returnErr error) {
	err := verifyPPMShipment(appCtx, ppmShipmentID)
	if err != nil {
		return nil, "", err
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
			EagerPreloadAssociationGunSafeWeightTickets,
			EagerPreloadAssociationMovingExpenses,
		},
		[]string{
			PostLoadAssociationWeightTicketUploads,
			PostLoadAssociationProgearWeightTicketUploads,
			PostLoadAssociationGunSafeWeightTicketUploads,
			PostLoadAssociationMovingExpenseUploads,
		},
	)

	// something bad happened on data retrieval of everything for PPM
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to load PPMShipment")
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, "", err
	}

	var pdfFilesToMerge []io.ReadSeeker

	// use aoa creator to generated SSW and Orders PDF
	aoaPacketFile, dirPath, err := p.aoaPacketCreator.CreateAOAPacket(appCtx, ppmShipmentID, true)
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, fmt.Sprintf("failed to generate AOA packet for ppmShipmentID: %s", ppmShipmentID.String()))

		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, dirPath, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	dirName := dirPath[strings.LastIndex(dirPath, "/")+1:]

	// AOA packet will be appended at the beginning of the final pdf file
	pdfFilesToMerge = append(pdfFilesToMerge, aoaPacketFile)

	// Start building individual PDFs for each expense/receipt docs. These files will then be merged as one PDF.
	var pdfFileNamesToMerge []string
	var pdfFileNamesToMergePdf afero.File
	var perr error

	sortedPaymentPacketItemsMap := buildPaymentPacketItemsMap(ppmShipment)

	for i := 0; i < len(sortedPaymentPacketItemsMap); i++ {
		pdfFileName, perr := p.pdfGenerator.ConvertUploadToPDF(appCtx, sortedPaymentPacketItemsMap[i].Upload, dirName)
		if perr != nil {
			errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to generate pdf for upload")
			appCtx.Logger().Error(errMsgPrefix, zap.Error(perr))
			return nil, dirPath, fmt.Errorf("%s: %w", errMsgPrefix, perr)
		}
		pdfFileNamesToMerge = append(pdfFileNamesToMerge, pdfFileName)
	}

	if len(pdfFileNamesToMerge) > 0 {
		pdfFileNamesToMergePdf, perr = p.pdfGenerator.MergePDFFiles(appCtx, pdfFileNamesToMerge, dirName)
		if perr != nil {
			errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed pdfGenerator.MergePDFFiles")
			appCtx.Logger().Error(errMsgPrefix, zap.Error(perr))
			return nil, dirPath, fmt.Errorf("%s: %w", errMsgPrefix, perr)
		}
		pdfFilesToMerge = append(pdfFilesToMerge, pdfFileNamesToMergePdf)
	}

	// Do final merge of all PDFs into one.
	finalMergePdf, err := p.pdfGenerator.MergePDFFilesByContents(appCtx, pdfFilesToMerge, dirName)
	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to generated file merged pdf")
		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))
		return nil, dirPath, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	// See https://github.com/transcom/mymove/pull/14496 for removal of bookmarks and watermarks

	defer func() {
		// if a panic occurred we set an error message that we can use to check for a recover in the calling method
		if r := recover(); r != nil {
			appCtx.Logger().Error("payment packet files panic", zap.Error(err))
			returnErr = fmt.Errorf("%s: panic", errMsgPrefix)
		}
	}()

	return finalMergePdf, dirPath, nil
}

func (p *paymentPacketCreator) GenerateDefault(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (afero.File, string, error) {
	return p.Generate(appCtx, ppmShipmentID, true)
}

// remove all of the packet files from the temp directory associated with creating the payment packet
func (p *paymentPacketCreator) CleanupPaymentPacketFile(packetFile afero.File, closeFile bool) error {
	if closeFile {
		if err := packetFile.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			return err
		}
	}

	fs := p.pdfGenerator.FileSystem()
	exists, err := afero.Exists(fs, packetFile.Name())

	if err != nil {
		return err
	}

	if exists {
		err := fs.Remove(packetFile.Name())

		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
				// File does not exist treat it as non-error:
				return nil
			}

			// Return the error if it's not a "file not found" error
			return err
		}
	}

	return nil
}

// remove all of the packet files from the temp directory associated with creating the Payment Packet
func (p *paymentPacketCreator) CleanupPaymentPacketDir(dirPath string) error {
	// RemoveAll does not return an error if the directory doesn't exist it will just do nothing and return nil
	return p.pdfGenerator.FileSystem().RemoveAll(dirPath)
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

	// process gun safe
	ppmShipment.GunSafeWeightTickets = ppmShipment.GunSafeWeightTickets.FilterRejected()
	gunSafeWeightTicketCnt := 1
	for _, gunSafeWeightTicket := range ppmShipment.GunSafeWeightTickets {
		sectionLabel := fmt.Sprintf("Pro-gear: Set #%s:", fmt.Sprint(gunSafeWeightTicketCnt))
		gunSafeWeightTicketSetCnt := 1
		for _, uu := range gunSafeWeightTicket.Document.UserUploads {
			sortedPaymentPacketItems[sortedPaymentPacketItemsIndex] = newPaymentPacketItem(fmt.Sprintf("%s Weight Ticket Document #%s", sectionLabel, fmt.Sprint(gunSafeWeightTicketSetCnt)), uu.Upload)
			sortedPaymentPacketItemsIndex++
			gunSafeWeightTicketSetCnt++
		}
		gunSafeWeightTicketCnt++
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
