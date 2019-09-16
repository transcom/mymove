package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/uploader"
)

// CreatePersonallyProcuredMoveAttachmentsHandler creates a PPM Attachments PDF
type CreatePersonallyProcuredMoveAttachmentsHandler struct {
	handlers.HandlerContext
}

func newGeneratedObjectMetaData() map[string]*string {
	metaData := make(map[string]*string)
	avStatus := "CLEANED"
	metaData["av-status"] = &avStatus
	avNotes := "GENERATED"
	metaData["av-notes"] = &avNotes
	return metaData
}

// Handle is the handler
func (h CreatePersonallyProcuredMoveAttachmentsHandler) Handle(params ppmop.CreatePPMAttachmentsParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())
	logger.Info("got ppm id: ", zap.Any("id", ppmID))

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	err = h.DB().Load(ppm, "Move.Orders.UploadedOrders.Uploads")
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	moveDocs, err := ppm.FetchMoveDocumentsForTypes(h.DB(), params.DocTypes)
	if err != nil {
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}

	// Init our tools
	loader, err := uploader.NewUploader(h.DB(), logger, h.FileStorer(), 100*uploader.MB)
	if err != nil {
		logger.Error("could not instantiate uploader", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}
	generator, err := paperwork.NewGenerator(h.DB(), logger, loader)
	if err != nil {
		logger.Error("failed to initialize generator", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}
	defer func() {
		if cleanupErr := generator.Cleanup(); cleanupErr != nil {
			logger.Error("failed to cleanup", zap.Error(cleanupErr))
		}
	}()

	// Start with uploaded orders info
	uploads := ppm.Move.Orders.UploadedOrders.Uploads

	// Flatten out uploads into a slice
	for _, moveDoc := range moveDocs {
		uploads = append(uploads, moveDoc.Document.Uploads...)
	}
	if len(uploads) == 0 {
		return ppmop.NewCreatePPMAttachmentsFailedDependency()
	}

	// Convert to PDF and merge into single PDF
	mergedPdf, err := generator.CreateMergedPDFUpload(uploads)
	if err != nil {
		logger.Error("failed to merge PDF files", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsUnprocessableEntity()
	}

	// Upload merged PDF to S3 and return Upload object
	metaData := newGeneratedObjectMetaData()
	pdfUpload, verrs, err := loader.CreateUpload(session.UserID, &mergedPdf, uploader.AllowedTypesPDF, metaData)
	if verrs.HasAny() || err != nil {
		switch err.(type) {
		case uploader.ErrTooLarge:
			return ppmop.NewCreatePPMAttachmentsRequestEntityTooLarge()
		default:
			return handlers.ResponseForVErrors(logger, verrs, err)
		}
	}

	url, err := loader.PresignedURL(pdfUpload)
	if err != nil {
		logger.Error("failed to get presigned url", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}

	uploadPayload := payloadForUploadModel(*pdfUpload, url)
	return ppmop.NewCreatePPMAttachmentsOK().WithPayload(uploadPayload)
}
