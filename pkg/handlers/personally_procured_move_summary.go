package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/uploader"
)

var moveDocumentSummaryTypes = []models.MoveDocumentType{
	models.MoveDocumentTypeOTHER,
	models.MoveDocumentTypeWEIGHTTICKET,
	models.MoveDocumentTypeSTORAGEEXPENSE,
	models.MoveDocumentTypeEXPENSE,
}

// CreatePersonallyProcuredMoveSummaryHandler creates a PPM Summary
type CreatePersonallyProcuredMoveSummaryHandler HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveSummaryHandler) Handle(params ppmop.CreatePPMSummaryParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.db, session, ppmID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	err = h.db.Load(ppm, "Move.Orders.UploadedOrders.Uploads")
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Fetch move documents with matching types
	moveDocs, err := ppm.FetchMoveDocumentsForTypes(h.db, moveDocumentSummaryTypes)
	if err != nil {
		return ppmop.NewCreatePPMSummaryInternalServerError()
	}
	if len(moveDocs) == 0 {
		return ppmop.NewCreatePPMSummaryFailedDependency()
	}

	// Init our tools
	loader := uploader.NewUploader(h.db, h.logger, h.storage)
	generator, err := paperwork.NewGenerator(h.db, h.logger, loader)
	if err != nil {
		h.logger.Error("failed to initialize generator", zap.Error(err))
		return ppmop.NewCreatePPMSummaryInternalServerError()
	}

	// Start with uploaded orders info
	// uploads := ppm.Move.Orders.UploadedOrders.Uploads
	uploads := models.Uploads{}

	// Flatten out uploads into a slice
	for _, moveDoc := range moveDocs {
		uploads = append(uploads, moveDoc.Document.Uploads...)
	}
	if len(uploads) == 0 {
		return ppmop.NewCreatePPMSummaryFailedDependency()
	}

	// Convert to PDF and merge into single PDF
	mergedPdf, err := generator.CreateMergedPDFUpload(uploads)
	if err != nil {
		h.logger.Error("failed to merge PDF files", zap.Error(err))
		return ppmop.NewCreatePPMSummaryUnprocessableEntity()
	}

	// Upload merged PDF to S3 and return Upload object
	pdfUpload, verrs, err := loader.CreateUpload(nil, session.UserID, mergedPdf)
	if verrs.HasAny() || err != nil {
		return responseForVErrors(h.logger, verrs, err)
	}

	url, err := loader.PresignedURL(pdfUpload)
	if err != nil {
		h.logger.Error("failed to get presigned url", zap.Error(err))
		return ppmop.NewCreatePPMSummaryInternalServerError()
	}

	uploadPayload := payloadForUploadModel(*pdfUpload, url)
	return ppmop.NewCreatePPMSummaryOK().WithPayload(uploadPayload)
}
