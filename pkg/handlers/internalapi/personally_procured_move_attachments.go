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

// Handle is the handler
func (h CreatePersonallyProcuredMoveAttachmentsHandler) Handle(params ppmop.CreatePPMAttachmentsParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	//  UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())
	logger.Info("got ppm id: ", zap.Any("id", ppmID))

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	err = h.DB().Load(ppm, "Move.Orders.UploadedOrders.UserUploads.Upload")
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	moveDocs, err := ppm.FetchMoveDocumentsForTypes(h.DB(), params.DocTypes)
	if err != nil {
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}

	// Init our tools
	loader, err := uploader.NewUserUploader(h.DB(), logger, h.FileStorer(), 100*uploader.MB)
	if err != nil {
		logger.Error("could not instantiate uploader", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}
	generator, err := paperwork.NewGenerator(h.DB(), logger, loader.Uploader())
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
	uploads, err := models.UploadsFromUserUploads(h.DB(), ppm.Move.Orders.UploadedOrders.UserUploads)
	if err != nil {
		logger.Error("failed to get uploads for Orders.UploadedOrders", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsFailedDependency()
	}

	// Flatten out uploads into a slice
	for _, moveDoc := range moveDocs {
		moveDocUploads, moveDocUploadsErr := models.UploadsFromUserUploadsNoDatabase(moveDoc.Document.UserUploads)
		if moveDocUploadsErr != nil {
			logger.Error("failed to get uploads for moveDoc.Document.UserUploads", zap.Error(moveDocUploadsErr))
			return ppmop.NewCreatePPMAttachmentsFailedDependency()
		}
		uploads = append(uploads, moveDocUploads...)
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

	// Add relevant av-.* tags for generated objects (for s3)
	generatedObjectTags := handlers.FmtString("av-status=CLEAN&av-notes=GENERATED")
	file := uploader.File{File: mergedPdf, Tags: generatedObjectTags}
	// UserUpload merged PDF to S3 and return UserUpload object
	pdfUpload, verrs, err := loader.CreateUserUpload(session.UserID, file, uploader.AllowedTypesPDF)
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

	uploadPayload := payloadForUploadModel(h.FileStorer(), pdfUpload.Upload, url)
	return ppmop.NewCreatePPMAttachmentsOK().WithPayload(uploadPayload)
}
