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

// paymentPacketCreator is the concrete struct implementing the PaymentPacketCreator interface
type paymentPacketCreator struct {
	services.PPMShipmentFetcher
	// TODO: Add the SSW PDF generation service object here.
	services.UserUploadToPDFConverter
	services.PDFMerger
	services.PPMShipmentUpdater
	*uploader.UserUploader
}

// NewPaymentPacketCreator creates a new PaymentPacketCreator with all of its dependencies
func NewPaymentPacketCreator(
	ppmShipmentFetcher services.PPMShipmentFetcher,
	userUploadToPDFConverter services.UserUploadToPDFConverter,
	pdfMerger services.PDFMerger,
	ppmShipmentUpdater services.PPMShipmentUpdater,
	userUploader *uploader.UserUploader,
) services.PaymentPacketCreator {
	return &paymentPacketCreator{
		ppmShipmentFetcher,
		userUploadToPDFConverter,
		pdfMerger,
		ppmShipmentUpdater,
		userUploader,
	}
}

// CreatePaymentPacket creates a payment packet for a PPM Shipment containing the shipment summary worksheet (SSW),
// uploaded orders, and any accepted uploaded PPM documents (i.e. weight tickets, pro-gear weight tickets, and moving
// expenses).
func (p *paymentPacketCreator) CreatePaymentPacket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) error {
	errMsgPrefix := "error creating payment packet"

	ppmShipment, err := p.PPMShipmentFetcher.GetPPMShipment(
		appCtx,
		ppmShipmentID,
		// TODO: We'll need to update this with any associations we need to build the SSW and to get the correct orders,
		//  taking amended orders into account if needed.
		[]string{
			EagerPreloadAssociationServiceMember,
			EagerPreloadAssociationWeightTickets,
			EagerPreloadAssociationProgearWeightTickets,
			EagerPreloadAssociationMovingExpenses,
			EagerPreloadAssociationPaymentPacket,
		},
		[]string{
			PostLoadAssociationWeightTicketUploads,
			PostLoadAssociationProgearWeightTicketUploads,
			PostLoadAssociationMovingExpenseUploads,
			PostLoadAssociationUploadedOrders,
		},
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

	filesToMerge, filesToMergeErr := p.UserUploadToPDFConverter.ConvertUserUploadsToPDF(
		appCtx,
		gatherPPMShipmentUserUploads(appCtx, ppmShipment),
	)

	if filesToMergeErr != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to convert uploads to PDF")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(filesToMergeErr))

		return fmt.Errorf("%s: %w", errMsgPrefix, filesToMergeErr)
	}

	for _, fileToMerge := range filesToMerge {
		fileToMerge := fileToMerge

		pdfsToMerge = append(pdfsToMerge, fileToMerge.PDFStream)
	}

	paymentPacket, mergeErr := p.PDFMerger.MergePDFs(appCtx, pdfsToMerge)

	if mergeErr != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to merge PDFs")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(mergeErr))

		return fmt.Errorf("%s: %w", errMsgPrefix, mergeErr)
	}

	saveErr := savePaymentPacket(appCtx, ppmShipment, paymentPacket, p.PPMShipmentUpdater, p.UserUploader)

	if saveErr != nil {
		errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to save payment packet")

		appCtx.Logger().Error(errMsgPrefix, zap.Error(saveErr))

		return fmt.Errorf("%s: %w", errMsgPrefix, saveErr)
	}

	return nil
}

// gatherPPMShipmentUserUploads is a helper func that gathers all of the user uploads associated with a PPM shipment.
// Mainly exists to keep the main func logic a little more legible
func gatherPPMShipmentUserUploads(_ appcontext.AppContext, ppmShipment *models.PPMShipment) models.UserUploads {
	userUploads := models.UserUploads{
		ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0],
	}

	for _, weightTicket := range ppmShipment.WeightTickets {
		weightTicket := weightTicket

		userUploads = append(userUploads, weightTicket.EmptyDocument.UserUploads...)
		userUploads = append(userUploads, weightTicket.FullDocument.UserUploads...)
		userUploads = append(userUploads, weightTicket.ProofOfTrailerOwnershipDocument.UserUploads...)
	}

	for _, progearWeightTicket := range ppmShipment.ProgearWeightTickets {
		progearWeightTicket := progearWeightTicket

		userUploads = append(userUploads, progearWeightTicket.Document.UserUploads...)
	}

	for _, movingExpense := range ppmShipment.MovingExpenses {
		movingExpense := movingExpense

		userUploads = append(userUploads, movingExpense.Document.UserUploads...)
	}

	return userUploads
}

// savePaymentPacket uploads the payment packet to S3 and saves the document data to the database, associating it with
// the PPM shipment.
func savePaymentPacket(
	appCtx appcontext.AppContext,
	ppmShipment *models.PPMShipment,
	paymentPacket io.ReadCloser,
	ppmShipmentUpdater services.PPMShipmentUpdater,
	userUploader *uploader.UserUploader,
) error {
	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Need to think about whether we want any special handling if we already have an payment packet. I don't think we will
		// create more than one, but if we do, we'll need to handle it both here and on retrieval.
		if ppmShipment.PaymentPacketID == nil {
			document := models.Document{
				ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID,
			}

			verrs, err := appCtx.DB().ValidateAndCreate(&document)

			errMsgPrefix := "failed to create payment packet document"

			if verrs.HasAny() {
				appCtx.Logger().Error(errMsgPrefix, zap.Error(verrs))

				return fmt.Errorf("%s: %w", errMsgPrefix, verrs)
			} else if err != nil {
				appCtx.Logger().Error(errMsgPrefix, zap.Error(err))

				return fmt.Errorf("%s: %w", errMsgPrefix, err)
			}

			ppmShipment.PaymentPacketID = &document.ID
			ppmShipment.PaymentPacket = &document

			ppmShipment, err = ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(txnAppCtx, ppmShipment, ppmShipment.Shipment.ID)

			if err != nil {
				errMsgPrefix = "failed to update PPMShipment with payment packet document"

				appCtx.Logger().Error(errMsgPrefix, zap.Error(err))

				return fmt.Errorf("%s: %w", errMsgPrefix, err)
			}
		}

		fileToUpload, prepErr := userUploader.PrepareFileForUpload(txnAppCtx, paymentPacket, "payment_packet.pdf")

		if prepErr != nil {
			errMsgPrefix := "failed to prepare payment packet for upload"

			appCtx.Logger().Error(errMsgPrefix, zap.Error(prepErr))

			return fmt.Errorf("%s: %w", errMsgPrefix, prepErr)
		}

		newUpload, uploadVerrs, uploadErr := userUploader.CreateUserUploadForDocument(
			txnAppCtx,
			ppmShipment.PaymentPacketID,
			// We're doing this on behalf of the service member, so we'll use their user ID to store the upload
			ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
			uploader.File{File: fileToUpload},
			uploader.AllowedTypesPPMDocuments,
		)

		errMsgPrefix := "failed to upload payment packet"
		if uploadVerrs.HasAny() {
			appCtx.Logger().Error(errMsgPrefix, zap.Error(uploadVerrs))

			return fmt.Errorf("%s: %w", errMsgPrefix, uploadVerrs)
		} else if uploadErr != nil {
			appCtx.Logger().Error(errMsgPrefix, zap.Error(uploadErr))

			return fmt.Errorf("%s: %w", errMsgPrefix, uploadErr)
		}

		ppmShipment.PaymentPacket.UserUploads = append(ppmShipment.PaymentPacket.UserUploads, *newUpload)

		return nil
	})

	if txnErr != nil {
		return txnErr
	}
	return nil
}
