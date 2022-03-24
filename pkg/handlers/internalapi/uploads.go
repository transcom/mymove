package internalapi

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

func payloadForUploadModel(
	storer storage.FileStorer,
	upload models.Upload,
	url string,
) *internalmessages.UploadPayload {
	uploadPayload := &internalmessages.UploadPayload{
		ID:          handlers.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         handlers.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

// CreateUploadHandler creates a new upload via POST /documents/{documentID}/uploads
type CreateUploadHandler struct {
	handlers.HandlerContext
}

// Handle creates a new UserUpload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := fmt.Errorf("Error creating upload")

			file, ok := params.File.(*runtime.File)
			if !ok {
				appCtx.Logger().Error("This should always be a runtime.File, something has changed in go-swagger.")
				return uploadop.NewCreateUploadInternalServerError(), rollbackErr
			}

			appCtx.Logger().Info(
				"File uploader and size",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
				zap.String("officeUserID", appCtx.Session().OfficeUserID.String()),
				zap.String("AdminUserID", appCtx.Session().AdminUserID.String()),
				zap.Int64("size", file.Header.Size),
			)

			var docID *uuid.UUID
			if params.DocumentID != nil {
				documentID, err := uuid.FromString(params.DocumentID.String())
				if err != nil {
					appCtx.Logger().Info("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
					return uploadop.NewCreateUploadBadRequest(), rollbackErr
				}

				// Fetch document to ensure user has access to it
				document, docErr := models.FetchDocument(appCtx.DB(), appCtx.Session(), documentID, true)
				if docErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), docErr), rollbackErr
				}
				docID = &document.ID
			}

			newUserUpload, url, verrs, createErr := uploaderpkg.CreateUserUploadForDocumentWrapper(
				appCtx,
				appCtx.Session().UserID,
				h.FileStorer(),
				file,
				file.Header.Filename,
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
				docID,
			)

			if verrs.HasAny() || createErr != nil {
				appCtx.Logger().Error("failed to create new user upload", zap.Error(createErr), zap.String("verrs", verrs.Error()))
				switch createErr.(type) {
				case uploaderpkg.ErrTooLarge:
					return uploadop.NewCreateUploadRequestEntityTooLarge(), rollbackErr
				case uploaderpkg.ErrFile:
					return uploadop.NewCreateUploadInternalServerError(), rollbackErr
				case uploaderpkg.ErrFailedToInitUploader:
					return uploadop.NewCreateUploadInternalServerError(), rollbackErr
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, createErr), rollbackErr
				}
			}

			uploadPayload := payloadForUploadModel(h.FileStorer(), newUserUpload.Upload, url)
			return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload), nil
		})
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler struct {
	handlers.HandlerContext
}

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			uploadID, _ := uuid.FromString(params.UploadID.String())
			userUpload, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			userUploader, err := uploaderpkg.NewUserUploader(
				h.FileStorer(),
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
			)
			if err != nil {
				appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			}
			if err = userUploader.DeleteUserUpload(appCtx, &userUpload); err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			return uploadop.NewDeleteUploadNoContent()
		})
}

// DeleteUploadsHandler deletes a collection of uploads
type DeleteUploadsHandler struct {
	handlers.HandlerContext
}

// Handle deletes uploads
func (h DeleteUploadsHandler) Handle(params uploadop.DeleteUploadsParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			userUploader, err := uploaderpkg.NewUserUploader(
				h.FileStorer(),
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
			)
			if err != nil {
				appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			}

			for _, uploadID := range params.UploadIds {
				uploadUUID, _ := uuid.FromString(uploadID.String())
				userUpload, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadUUID)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err)
				}

				if err = userUploader.DeleteUserUpload(appCtx, &userUpload); err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err)
				}
			}

			return uploadop.NewDeleteUploadsNoContent()
		})
}
