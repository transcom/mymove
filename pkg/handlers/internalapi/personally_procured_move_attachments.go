package internalapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
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
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())

	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching personally procured move", zap.String("personally_procured_move_id", ppmID.String()))
	}

	err = h.DB().Load(ppm, "Move.Orders.UploadedOrders.Uploads")
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error loading uploads")
	}

	moveDocs, err := ppm.FetchMoveDocumentsForTypes(h.DB(), params.DocTypes)
	if err != nil {
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}
	if len(moveDocs) == 0 {
		return ppmop.NewCreatePPMAttachmentsFailedDependency()
	}

	// Init our tools
	loader := uploader.NewUploader(h.DB(), h.Logger(), h.FileStorer())
	generator, err := paperwork.NewGenerator(h.DB(), h.Logger(), loader)
	if err != nil {
		h.Logger().Error("failed to initialize generator", zap.Error(err))
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
		h.Logger().Error("failed to merge PDF files", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsUnprocessableEntity()
	}

	// Upload merged PDF to S3 and return Upload object
	pdfUpload, verrs, err := loader.CreateUpload(session.UserID, &mergedPdf, uploader.AllowedTypesPDF)
	if verrs.HasAny() || err != nil {
		return h.RespondAndTraceVErrors(ctx, verrs, err, "error creating upload")
	}

	url, err := loader.PresignedURL(pdfUpload)
	if err != nil {
		h.Logger().Error("failed to get presigned url", zap.Error(err))
		return ppmop.NewCreatePPMAttachmentsInternalServerError()
	}

	uploadPayload := payloadForUploadModel(*pdfUpload, url)
	return ppmop.NewCreatePPMAttachmentsOK().WithPayload(uploadPayload)
}
