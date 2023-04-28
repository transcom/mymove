package ppmshipment

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// aoaPacketCreator is the concrete struct implementing the AOAPacketCreator interface
type aoaPacketCreator struct {
	services.PPMShipmentFetcher
	// TODO: Add the SSW PDF generation service object here.
	services.UserUploadToPDFConverter
	services.PDFMerger
	services.PPMShipmentUpdater
	*uploader.UserUploader
}

// NewAOAPacketCreator creates a new AOAPacketCreator with all of its dependencies
func NewAOAPacketCreator(
	ppmShipmentFetcher services.PPMShipmentFetcher,
	userUploadToPDFConverter services.UserUploadToPDFConverter,
	pdfMerger services.PDFMerger,
	ppmShipmentUpdater services.PPMShipmentUpdater,
	userUploader *uploader.UserUploader,
) services.AOAPacketCreator {
	return &aoaPacketCreator{
		ppmShipmentFetcher,
		userUploadToPDFConverter,
		pdfMerger,
		ppmShipmentUpdater,
		userUploader,
	}
}

// CreateAOAPacket creates an AOA packet for a PPM Shipment, containing the shipment summary worksheet (SSW) and
// uploaded orders.
func (a *aoaPacketCreator) CreateAOAPacket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) error {
	errMsgPrefix := "error creating AOA packet"

	ppmShipment, err := a.PPMShipmentFetcher.GetPPMShipment(
		appCtx,
		ppmShipmentID,
		// TODO: We'll need to update this with any associations we need to build the SSW and to get the correct orders,
		//  taking amended orders into account if needed.
		[]string{EagerPreloadAssociationServiceMember, EagerPreloadAssociationAOAPacket},
		[]string{PostLoadAssociationUploadedOrders},
	)

	if err != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to load PPMShipment")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(err))

		return fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	pdfsToMerge := []io.ReadCloser{}

	// TODO: Call yet-to-be-written SSW PDF generation service object and add it to the list of PDFs to merge.
	//  using fake names here for the service object and its receiver function as a sample of how this might look
	//ssw, sswErr := a.SSWPDFGenerator.GenerateSSWPDF(appCtx, ppmShipmentID)
	//
	//if sswErr != nil {
	//	errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "error generating SSW PDF")
	//
	//	appCtx.Logger().Error(errMsgPrefix, zap.Error(sswErr))
	//
	//	return fmt.Errorf("%s: %w", errMsgPrefix, sswErr)
	//}
	//
	//pdfsToMerge = append(pdfsToMerge, ssw)

	filesToMerge, filesToMergeErr := a.UserUploadToPDFConverter.ConvertUserUploadsToPDF(
		appCtx,
		models.UserUploads{ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0]},
	)

	if filesToMergeErr != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to convert orders to PDF")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(filesToMergeErr))

		return fmt.Errorf("%s: %w", errMsgPrefix, filesToMergeErr)
	}

	for _, fileToMerge := range filesToMerge {
		fileToMerge := fileToMerge

		pdfsToMerge = append(pdfsToMerge, fileToMerge.PDFStream)
	}

	aoaPacket, mergeErr := a.PDFMerger.MergePDFs(appCtx, pdfsToMerge)

	if mergeErr != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to merge SSW and orders PDFs")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(mergeErr))

		return fmt.Errorf("%s: %w", errMsgPrefix, mergeErr)
	}

	saveErr := saveAOAPacket(appCtx, ppmShipment, aoaPacket, a.PPMShipmentUpdater, a.UserUploader)

	if saveErr != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to save AOA packet")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(saveErr))

		return fmt.Errorf("%s: %w", errMsgPrefix, saveErr)
	}

	return nil
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
