package internal

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/uploader"
)

var moveDocumentAttachmentsTypes = []models.MoveDocumentType{
	models.MoveDocumentTypeOTHER,
	models.MoveDocumentTypeWEIGHTTICKET,
	models.MoveDocumentTypeSTORAGEEXPENSE,
	models.MoveDocumentTypeEXPENSE,
}

// CreatePersonallyProcuredMoveAttachmentsHandler creates a PPM Attachments PDF
type CreatePersonallyProcuredMoveAttachmentsHandler utils.HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveAttachmentsHandler) Handle(params ppmop.CreatePPMAttachmentsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.Db, session, ppmID)
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}

	err = h.Db.Load(ppm, "Move.Orders.UploadedOrders.Uploads")
	if err != nil {
		return utils.ResponseForError(h.Logger, err)
	}

	// Fetch move documents with matching types
	moveDocs, err := ppm.FetchMoveDocumentsForTypes(h.Db, moveDocumentAttachmentsTypes)
	if err != nil {
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}
	if len(moveDocs) == 0 {
		return ppmop.NewCreatePPMAttachmentsFailedDependency()
	}

	// Init our tools
	loader := uploader.NewUploader(h.Db, h.Logger, h.Storage)
	generator, err := paperwork.NewGenerator(h.Db, h.Logger, loader)
	if err != nil {
		h.Logger.Error("failed to initialize generator", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}

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
		h.Logger.Error("failed to merge PDF files", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsUnprocessableEntity()
	}

	// Upload merged PDF to S3 and return Upload object
	pdfUpload, verrs, err := loader.CreateUpload(nil, session.UserID, mergedPdf)
	if verrs.HasAny() || err != nil {
		return utils.ResponseForVErrors(h.Logger, verrs, err)
	}

	url, err := loader.PresignedURL(pdfUpload)
	if err != nil {
		h.Logger.Error("failed to get presigned url", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}

	uploadPayload := payloadForUploadModel(*pdfUpload, url)
	return ppmop.NewCreatePPMAttachmentsOK().WithPayload(uploadPayload)
}
