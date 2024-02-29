package ppmshipment

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// paymentPacketCreator is the concrete struct implementing the PaymentPacketCreator interface
type paymentPacketCreator struct {
	services.PPMShipmentFetcher
	*uploader.UserUploader
	pdfGenerator     paperwork.Generator
	aoaPacketCreator services.AOAPacketCreator
}

// NewPaymentPacketCreator creates a new PaymentPacketCreator with all of its dependencies
func NewPaymentPacketCreator(
	ppmShipmentFetcher services.PPMShipmentFetcher,
	// userUploadToPDFConverter services.UserUploadToPDFConverter,
	// pdfMerger services.PDFMerger,
	// ppmShipmentUpdater services.PPMShipmentUpdater,
	userUploader *uploader.UserUploader,
	aoaPacketCreator services.AOAPacketCreator,
) services.PaymentPacketCreator {
	pdfGenerator, _ := paperwork.NewGenerator(userUploader.Uploader())
	return &paymentPacketCreator{
		ppmShipmentFetcher,
		userUploader,
		*pdfGenerator,
		aoaPacketCreator,
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

	//pdfsToMerge := []io.ReadCloser{}

	uploads := gatherPPMShipmentUserUploads(appCtx, ppmShipment) //models.UserUploads

	//p.pdfGenerator.ConvertUploadsToPDF(appCtx, uploads)
	p.pdfGenerator.CreateMergedPDFUpload(appCtx, uploads)
	// filesToMerge, filesToMergeErr := p.UserUploadToPDFConverter.ConvertUserUploadsToPDF(
	// 	appCtx,
	// 	gatherPPMShipmentUserUploads(appCtx, ppmShipment),
	// )

	// if filesToMergeErr != nil {
	// 	errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to convert uploads to PDF")

	// 	appCtx.Logger().Error(errMsgPrefix, zap.Error(filesToMergeErr))

	// 	return fmt.Errorf("%s: %w", errMsgPrefix, filesToMergeErr)
	// }

	// for _, fileToMerge := range filesToMerge {
	// 	fileToMerge := fileToMerge

	// 	pdfsToMerge = append(pdfsToMerge, fileToMerge.PDFStream)
	// }

	// paymentPacket, mergeErr := p.PDFMerger.MergePDFs(appCtx, pdfsToMerge)

	// if mergeErr != nil {
	// 	errMsgPrefix = fmt.Sprintf("%s: %s", errMsgPrefix, "failed to merge PDFs")

	// 	appCtx.Logger().Error(errMsgPrefix, zap.Error(mergeErr))

	// 	return fmt.Errorf("%s: %w", errMsgPrefix, mergeErr)
	// }

	// if paymentPacket != nil {
	// 	appCtx.Logger().Debug("hello")
	// }

	return nil
}

func (p *paymentPacketCreator) Generate(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (io.ReadCloser, error) {
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

		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	//pdfsToMerge := []io.ReadCloser{}

	uploads := gatherPPMShipmentUserUploads(appCtx, ppmShipment)

	var pdfFileNames []string

	//aoaPacketFile, _ := p.aoaPacketCreator.CreateAOAPacket(appCtx, ppmShipmentID)
	//pdfFileNames = append(pdfFileNames, aoaPacketFile.Name())

	fileNames, _ := p.pdfGenerator.ConvertUploadsToPDF(appCtx, uploads) //return fileNames

	pdfFileNames = append(pdfFileNames, fileNames...)

	return p.pdfGenerator.MergePDFFiles(appCtx, pdfFileNames)

	//return p.pdfGenerator.CreateMergedPDFUpload(appCtx, uploads)
}

func gatherPPMShipmentUserUploads(_ appcontext.AppContext, ppmShipment *models.PPMShipment) models.Uploads {
	uploads := models.Uploads{
		ppmShipment.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].Upload,
	}

	for _, weightTicket := range ppmShipment.WeightTickets {
		for _, uu := range weightTicket.EmptyDocument.UserUploads {
			uploads = append(uploads, uu.Upload)
		}
		for _, uu := range weightTicket.FullDocument.UserUploads {
			uploads = append(uploads, uu.Upload)
		}
		for _, uu := range weightTicket.ProofOfTrailerOwnershipDocument.UserUploads {
			uploads = append(uploads, uu.Upload)
		}
	}

	for _, progearWeightTicket := range ppmShipment.ProgearWeightTickets {
		for _, uu := range progearWeightTicket.Document.UserUploads {
			uploads = append(uploads, uu.Upload)
		}
	}

	for _, movingExpense := range ppmShipment.MovingExpenses {
		for _, uu := range movingExpense.Document.UserUploads {
			uploads = append(uploads, uu.Upload)
		}
	}

	return uploads
}
