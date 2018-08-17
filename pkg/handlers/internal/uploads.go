package internal

import (
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

func payloadForUploadModel(upload models.Upload, url string) *internalmessages.UploadPayload {
	return &internalmessages.UploadPayload{
		ID:          utils.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         utils.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   utils.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   utils.FmtDateTime(upload.UpdatedAt),
	}
}

// CreateUploadHandler creates a new upload via POST /documents/{documentID}/uploads
type CreateUploadHandler utils.HandlerContext

// Handle creates a new Upload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {

	file, ok := params.File.(*runtime.File)
	if !ok {
		h.logger.Error("This should always be a runtime.File, something has changed in go-swagger.")
		return uploadop.NewCreateUploadInternalServerError()
	}

	h.logger.Info("File name and size: ", zap.String("name", file.Header.Filename), zap.Int64("size", file.Header.Size))

	// User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	var docID *uuid.UUID
	if params.DocumentID != nil {
		documentID, err := uuid.FromString(params.DocumentID.String())
		if err != nil {
			h.logger.Info("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
			return uploadop.NewCreateUploadBadRequest()
		}

		// Fetch document to ensure user has access to it
		document, docErr := models.FetchDocument(h.db, session, documentID)
		if docErr != nil {
			return utils.ResponseForError(h.logger, docErr)
		}
		docID = &document.ID
	}

	// Read the incoming data into a new afero.File for consumption
	aFile, err := h.Storage.FileSystem().Create(file.Header.Filename)
	if err != nil {
		h.logger.Error("Error opening afero file.", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	_, err = io.Copy(aFile, file.Data)
	if err != nil {
		h.logger.Error("Error copying incoming data into afero file.", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	uploader := uploaderpkg.NewUploader(h.db, h.logger, h.Storage)
	newUpload, verrs, err := uploader.CreateUpload(docID, session.UserID, aFile)
	if err != nil {
		if cause := errors.Cause(err); cause == uploaderpkg.ErrZeroLengthFile {
			return uploadop.NewCreateUploadBadRequest()
		}
		h.logger.Error("Failed to create upload", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	} else if verrs.HasAny() {
		payload := utils.CreateFailedValidationPayload(verrs)
		return uploadop.NewCreateUploadBadRequest().WithPayload(payload)
	}

	url, err := uploader.PresignedURL(newUpload)
	if err != nil {
		h.logger.Error("failed to get presigned url", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}
	uploadPayload := payloadForUploadModel(*newUpload, url)
	return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload)
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler utils.HandlerContext

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	uploadID, _ := uuid.FromString(params.UploadID.String())
	upload, err := models.FetchUpload(h.db, session, uploadID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	uploader := uploaderpkg.NewUploader(h.db, h.logger, h.Storage)
	if err = uploader.DeleteUpload(&upload); err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	return uploadop.NewDeleteUploadNoContent()
}

// DeleteUploadsHandler deletes a collection of uploads
type DeleteUploadsHandler utils.HandlerContext

// Handle deletes uploads
func (h DeleteUploadsHandler) Handle(params uploadop.DeleteUploadsParams) middleware.Responder {
	// User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	uploader := uploaderpkg.NewUploader(h.db, h.logger, h.Storage)

	for _, uploadID := range params.UploadIds {
		uuid, _ := uuid.FromString(uploadID.String())
		upload, err := models.FetchUpload(h.db, session, uuid)
		if err != nil {
			return utils.ResponseForError(h.logger, err)
		}

		if err = uploader.DeleteUpload(&upload); err != nil {
			return utils.ResponseForError(h.logger, err)
		}
	}

	return uploadop.NewDeleteUploadsNoContent()
}
