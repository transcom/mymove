package ghcapi

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	uploadop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/upload"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

type CreateUploadHandler struct {
	handlers.HandlerConfig
}

func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := apperror.NewBadDataError("error creating upload")

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
				uploaderpkg.AllowedTypesServiceMember,
				docID,
				models.UploadTypeOFFICE,
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

			uploadPayload := payloads.PayloadForUploadModel(h.FileStorer(), newUserUpload.Upload, url)
			return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload), nil
		})
}

type UpdateUploadHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

func (h UpdateUploadHandler) Handle(params uploadop.UpdateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().IsOfficeApp() {
				forbiddenError := apperror.NewForbiddenError("User is not an Office User.")
				appCtx.Logger().Error(forbiddenError.Error())
				return uploadop.NewUpdateUploadForbidden(), forbiddenError
			}

			uploadID, _ := uuid.FromString(params.UploadID.String())
			updater := upload.NewUploadUpdater()
			newUpload, err := updater.UpdateUploadForRotation(appCtx, uploadID, params.Body.Rotation)
			if err != nil {
				return nil, apperror.NewBadDataError("unable to update upload")
			}

			url, err := h.FileStorer().PresignedURL(newUpload.StorageKey, newUpload.ContentType)
			if err != nil {
				return nil, err
			}

			uploadPayload := payloads.Upload(h.FileStorer(), *newUpload, url)

			return uploadop.NewUpdateUploadCreated().WithPayload(uploadPayload), nil
		})
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler struct {
	handlers.HandlerConfig
	services.UploadInformationFetcher
}

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().IsOfficeApp() {
				forbiddenError := apperror.NewForbiddenError("User is not an Office User.")
				appCtx.Logger().Error(forbiddenError.Error())
				return uploadop.NewDeleteUploadForbidden(), forbiddenError
			}

			uploadID, _ := uuid.FromString(params.UploadID.String())
			userUpload, err := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), uploadID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			userUploader, err := uploaderpkg.NewUserUploader(
				h.FileStorer(),
				uploaderpkg.MaxCustomerUserUploadFileSizeLimit,
			)
			if err != nil {
				appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			}
			if err = userUploader.DeleteUserUpload(appCtx, &userUpload); err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			return uploadop.NewDeleteUploadNoContent(), nil

		})
}
