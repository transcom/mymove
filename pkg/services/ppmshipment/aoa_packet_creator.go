package ppmshipment

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// aoaPacketCreator is the concrete struct implementing the AOAPacketCreator interface
type aoaPacketCreator struct {
	services.SSWPPMGenerator
	services.SSWPPMComputer
	services.PrimeDownloadMoveUploadPDFGenerator
	uploader.UserUploader
	pdfGenerator *paperwork.Generator
}

// NewAOAPacketCreator creates a new AOAPacketCreator with all of its dependencies
func NewAOAPacketCreator(
	sswPPMGenerator services.SSWPPMGenerator,
	sswPPMComputer services.SSWPPMComputer,
	primeDownloadMoveUploadPDFGenerator services.PrimeDownloadMoveUploadPDFGenerator,
	userUploader *uploader.UserUploader,
	pdfGenerator *paperwork.Generator,
) services.AOAPacketCreator {
	return &aoaPacketCreator{
		sswPPMGenerator,
		sswPPMComputer,
		primeDownloadMoveUploadPDFGenerator,
		*userUploader,
		pdfGenerator,
	}
}

func (a *aoaPacketCreator) VerifyAOAPacketInternal(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) error {
	appCtx.Logger().Info("Retrieving ServiceMember to verify authorization")
	ppmShipment := models.PPMShipment{}
	dbQErr := appCtx.DB().Q().Eager(
		"Shipment.MoveTaskOrder.Orders.ServiceMember",
	).Find(&ppmShipment, ppmShipmentID)

	if dbQErr != nil {
		appCtx.Logger().Error("Could not retrieve query from PPMShipment to ServiceMember")
		if errors.Cause(dbQErr).Error() == models.RecordNotFoundErrorString {
			return models.ErrFetchNotFound
		}
		return dbQErr
	}

	if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID == appCtx.Session().ServiceMemberID {
		return nil
	}

	// This returns the same client-facing error to prevent fishing for UUIDs after logging unauthorized access.
	appCtx.Logger().Error("Unauthorized AOA access attempted, Context Member: " +
		appCtx.Session().ServiceMemberID.String() + " attempted to access AOA Packet records for " +
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID.String())
	return errors.New("Not the authorized service member")

}

// CreateAOAPacket creates an AOA packet for a PPM Shipment, containing the shipment summary worksheet (SSW) and
// uploaded orders.
func (a *aoaPacketCreator) CreateAOAPacket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, isPaymentPacket bool) (afero.File, error) {
	errMsgPrefix := "error creating AOA packet"

	// First we begin by fetching SSW Data, computing obligations, formatting, and filling the SSWPDF
	ssfd, err := a.SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(appCtx, appCtx.Session(), ppmShipmentID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	page1Data, page2Data, page3Data, err := a.SSWPPMComputer.FormatValuesShipmentSummaryWorksheet(*ssfd, isPaymentPacket)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	SSWPPMWorksheet, SSWPDFInfo, err := a.SSWPPMGenerator.FillSSWPDFForm(page1Data, page2Data, page3Data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}
	// Ensure SSW PDF is not corrupted
	if SSWPDFInfo.PageCount != 3 {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	// Now that SSW is retrieved, find, convert to pdf, and append all orders and amendments
	// Query move, orders by ppm shipment
	ppmShipment := models.PPMShipment{}
	dbQErr := appCtx.DB().Q().Eager(
		"Shipment.MoveTaskOrder",
		"Shipment.MoveTaskOrder.Orders.ID",
	).Find(&ppmShipment, ppmShipmentID)

	if dbQErr != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, dbQErr)
	}

	// Find move attached to PPM Shipment
	move := models.Move(ppmShipment.Shipment.MoveTaskOrder)
	// This function retrieves all orders and amendments, converts and merges them into one PDF with bookmarks
	ordersFile, err := a.PrimeDownloadMoveUploadPDFGenerator.GenerateDownloadMoveUserUploadPDF(appCtx, services.MoveOrderUploadAll, move, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}
	// Ensure SSW PDF is not corrupted
	ordersFileInfo, err := a.pdfGenerator.GetPdfFileInfoByContents(ordersFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}
	if !(ordersFileInfo.PageCount > 0) {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	// Calling the PDF merge function in Generator with these filepaths creates issues due to instancing of the memory filesystem
	// Instead, we use a readseeker to pass in file information to merge the files in Generator.
	var files []io.ReadSeeker

	files = append(files, SSWPPMWorksheet)
	files = append(files, ordersFile)
	// Take all of generated PDFs and merge into a single PDF.
	mergedPdf, err := a.pdfGenerator.MergePDFFilesByContents(appCtx, files)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	return mergedPdf, nil
}

// remove all of the packet files from the temp directory associated with creating the AOA packet
func (a *aoaPacketCreator) CleanupAOAPacketFiles(appCtx appcontext.AppContext) error {
	err := a.pdfGenerator.Cleanup(appCtx)

	return err
}
