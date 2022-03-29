package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			appCtx.Logger().Info("got ppm id: ", zap.Any("id", ppmID))

			ppm, err := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			err = appCtx.DB().Load(ppm, "Move.Orders.UploadedOrders.UserUploads.Upload")
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			moveDocs, err := ppm.FetchMoveDocumentsForTypes(appCtx.DB(), params.DocTypes)
			if err != nil {
				return ppmop.NewCreatePPMAttachmentsInternalServerError(), err
			}

			// Init our tools
			loader, err := uploader.NewUserUploader(h.FileStorer(), uploader.MaxOfficeUploadFileSizeLimit)
			if err != nil {
				appCtx.Logger().Error("could not instantiate uploader", zap.Error(err))
				return ppmop.NewCreatePPMAttachmentsInternalServerError(), err
			}
			generator, err := paperwork.NewGenerator(loader.Uploader())
			if err != nil {
				appCtx.Logger().Error("failed to initialize generator", zap.Error(err))
				return ppmop.NewCreatePPMAttachmentsInternalServerError(), err
			}
			defer func() {
				if cleanupErr := generator.Cleanup(appCtx); cleanupErr != nil {
					appCtx.Logger().Error("failed to cleanup", zap.Error(cleanupErr))
				}
			}()

			// Start with uploaded orders info
			uploads, err := models.UploadsFromUserUploads(appCtx.DB(), ppm.Move.Orders.UploadedOrders.UserUploads)
			if err != nil {
				appCtx.Logger().Error("failed to get uploads for Orders.UploadedOrders", zap.Error(err))
				return ppmop.NewCreatePPMAttachmentsFailedDependency(), err
			}

			// Flatten out uploads into a slice
			for _, moveDoc := range moveDocs {
				moveDocUploads, moveDocUploadsErr := models.UploadsFromUserUploadsNoDatabase(moveDoc.Document.UserUploads)
				if moveDocUploadsErr != nil {
					appCtx.Logger().Error("failed to get uploads for moveDoc.Document.UserUploads", zap.Error(moveDocUploadsErr))
					return ppmop.NewCreatePPMAttachmentsFailedDependency(), moveDocUploadsErr
				}
				uploads = append(uploads, moveDocUploads...)
			}
			if len(uploads) == 0 {
				return ppmop.NewCreatePPMAttachmentsFailedDependency(), apperror.NewInternalServerError("no uploads found")
			}

			// Convert to PDF and merge into single PDF
			mergedPdf, err := generator.CreateMergedPDFUpload(appCtx, uploads)
			if err != nil {
				appCtx.Logger().Error("failed to merge PDF files", zap.Error(err))
				return ppmop.NewCreatePPMAttachmentsUnprocessableEntity(), err
			}

			// Add relevant av-.* tags for generated objects (for s3)
			generatedObjectTags := handlers.FmtString("av-status=CLEAN&av-notes=GENERATED")
			file := uploader.File{File: mergedPdf, Tags: generatedObjectTags}
			// UserUpload merged PDF to S3 and return UserUpload object
			pdfUpload, verrs, err := loader.CreateUserUpload(appCtx, appCtx.Session().UserID, file, uploader.AllowedTypesPDF)
			if verrs.HasAny() || err != nil {
				switch err.(type) {
				case uploader.ErrTooLarge:
					return ppmop.NewCreatePPMAttachmentsRequestEntityTooLarge(), err
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
				}
			}

			url, err := loader.PresignedURL(appCtx, pdfUpload)
			if err != nil {
				appCtx.Logger().Error("failed to get presigned url", zap.Error(err))
				return ppmop.NewCreatePPMAttachmentsInternalServerError(), err
			}

			uploadPayload := payloadForUploadModel(h.FileStorer(), pdfUpload.Upload, url)
			return ppmop.NewCreatePPMAttachmentsOK().WithPayload(uploadPayload), nil
		})
}
