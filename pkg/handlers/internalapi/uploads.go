package internalapi

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

func payloadForUploadModel(storer storage.FileStorer, upload models.Upload, url string) *internalmessages.UploadPayload {
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

	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromContext(ctx)

	file, ok := params.File.(*runtime.File)
	if !ok {
		logger.Error("This should always be a runtime.File, something has changed in go-swagger.")
		return uploadop.NewCreateUploadInternalServerError()
	}

	logger.Info(
		"File uploader and size",
		zap.String("userID", session.ServiceMemberID.String()),
		zap.String("serviceMemberID", session.ServiceMemberID.String()),
		zap.String("officeUserID", session.OfficeUserID.String()),
		zap.String("AdminUserID", session.AdminUserID.String()),
		zap.Int64("size", file.Header.Size),
	)

	var docID *uuid.UUID
	if params.DocumentID != nil {
		documentID, err := uuid.FromString(params.DocumentID.String())
		if err != nil {
			logger.Info("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
			return uploadop.NewCreateUploadBadRequest()
		}

		// Fetch document to ensure user has access to it
		document, docErr := models.FetchDocument(ctx, h.DB(), session, documentID, true)
		if docErr != nil {
			return handlers.ResponseForError(logger, docErr)
		}
		docID = &document.ID
	}

	userUploader, err := uploaderpkg.NewUserUploader(h.DB(), logger, h.FileStorer(), 25*uploaderpkg.MB)
	if err != nil {
		logger.Fatal("could not instantiate uploader", zap.Error(err))
	}

	aFile, err := userUploader.PrepareFileForUpload(file.Data, file.Header.Filename)
	if err != nil {
		logger.Fatal("could not prepare file for uploader", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	newUserUpload, verrs, err := userUploader.CreateUserUploadForDocument(docID, session.UserID, uploaderpkg.File{File: aFile}, uploaderpkg.AllowedTypesServiceMember)
	if verrs.HasAny() || err != nil {
		switch err.(type) {
		case uploaderpkg.ErrTooLarge:
			return uploadop.NewCreateUploadRequestEntityTooLarge()
		default:
			return handlers.ResponseForVErrors(logger, verrs, err)
		}
	}

	url, err := userUploader.PresignedURL(newUserUpload)
	if err != nil {
		logger.Error("failed to get presigned url", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}
	uploadPayload := payloadForUploadModel(h.FileStorer(), newUserUpload.Upload, url)
	return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload)
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler struct {
	handlers.HandlerContext
}

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	uploadID, _ := uuid.FromString(params.UploadID.String())
	userUpload, err := models.FetchUserUploadFromUploadID(ctx, h.DB(), session, uploadID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	userUploader, err := uploaderpkg.NewUserUploader(h.DB(), logger, h.FileStorer(), 25*uploaderpkg.MB)
	if err != nil {
		logger.Fatal("could not instantiate uploader", zap.Error(err))
	}
	if err = userUploader.DeleteUserUpload(&userUpload); err != nil {
		return handlers.ResponseForError(logger, err)
	}

	return uploadop.NewDeleteUploadNoContent()
}

// DeleteUploadsHandler deletes a collection of uploads
type DeleteUploadsHandler struct {
	handlers.HandlerContext
}

// Handle deletes uploads
func (h DeleteUploadsHandler) Handle(params uploadop.DeleteUploadsParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	// User should always be populated by middleware
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	userUploader, err := uploaderpkg.NewUserUploader(h.DB(), logger, h.FileStorer(), 25*uploaderpkg.MB)
	if err != nil {
		logger.Fatal("could not instantiate uploader", zap.Error(err))
	}

	for _, uploadID := range params.UploadIds {
		uploadUUID, _ := uuid.FromString(uploadID.String())
		userUpload, err := models.FetchUserUploadFromUploadID(ctx, h.DB(), session, uploadUUID)
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}

		if err = userUploader.DeleteUserUpload(&userUpload); err != nil {
			return handlers.ResponseForError(logger, err)
		}
	}

	return uploadop.NewDeleteUploadsNoContent()
}
