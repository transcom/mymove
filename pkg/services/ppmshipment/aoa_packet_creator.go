package ppmshipment

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

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

// saveAOAPacket uploads the AOA packet to S3 and saves the document data to the database, associating it with the PPM
// shipment.
func saveAOAPacket(
	appCtx appcontext.AppContext,
	ppmShipment *models.PPMShipment,
	aoaPacket io.ReadCloser,
	ppmShipmentUpdater services.PPMShipmentUpdater,
	userUploader *uploader.UserUploader,
) error {
	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Need to think about whether we want any special handling if we already have an AOA packet. I don't think we will
		// create more than one, but if we do, we'll need to handle it both here and on retrieval.
		if ppmShipment.AOAPacketID == nil {
			document := models.Document{
				ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID,
			}

			verrs, err := appCtx.DB().ValidateAndCreate(&document)

			errMsgPrefix := "failed to create AOA packet document"

			if verrs.HasAny() {
				appCtx.Logger().Error(errMsgPrefix, zap.Error(verrs))

				return fmt.Errorf("%s: %w", errMsgPrefix, verrs)
			} else if err != nil {
				appCtx.Logger().Error(errMsgPrefix, zap.Error(err))

				return fmt.Errorf("%s: %w", errMsgPrefix, err)
			}

			ppmShipment.AOAPacketID = &document.ID
			ppmShipment.AOAPacket = &document

			ppmShipment, err = ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(txnAppCtx, ppmShipment, ppmShipment.Shipment.ID)

			if err != nil {
				errMsgPrefix = "failed to update PPMShipment with AOA packet document"

				appCtx.Logger().Error(errMsgPrefix, zap.Error(err))

				return fmt.Errorf("%s: %w", errMsgPrefix, err)
			}
		}

		fileToUpload, prepErr := userUploader.PrepareFileForUpload(txnAppCtx, aoaPacket, "aoa_packet.pdf")

		if prepErr != nil {
			errMsgPrefix := "failed to prepare AOA packet for upload"

			appCtx.Logger().Error(errMsgPrefix, zap.Error(prepErr))

			return fmt.Errorf("%s: %w", errMsgPrefix, prepErr)
		}

		newUpload, uploadVerrs, uploadErr := userUploader.CreateUserUploadForDocument(
			txnAppCtx,
			ppmShipment.AOAPacketID,
			// We're doing this on behalf of the service member, so we'll use their user ID to store the upload
			ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
			uploader.File{File: fileToUpload},
			uploader.AllowedTypesPPMDocuments,
		)

		errMsgPrefix := "failed to upload AOA packet"
		if uploadVerrs.HasAny() {
			appCtx.Logger().Error(errMsgPrefix, zap.Error(uploadVerrs))

			return fmt.Errorf("%s: %w", errMsgPrefix, uploadVerrs)
		} else if uploadErr != nil {
			appCtx.Logger().Error(errMsgPrefix, zap.Error(uploadErr))

			return fmt.Errorf("%s: %w", errMsgPrefix, uploadErr)
		}

		ppmShipment.AOAPacket.UserUploads = append(ppmShipment.AOAPacket.UserUploads, *newUpload)

		return nil
	})

	if txnErr != nil {
		return txnErr
	}

	return nil
}
