package handlers

import (
	/*
		#nosec - we use md5 because it's required by the S3 API for
		validating data integrity.
		https://aws.amazon.com/premiumsupport/knowledge-center/data-integrity-s3/
	*/
	"crypto/md5"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/auth"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForUploadModel(upload models.Upload, url string) *internalmessages.UploadPayload {
	return &internalmessages.UploadPayload{
		ID:          fmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         fmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   fmtDateTime(upload.CreatedAt),
		UpdatedAt:   fmtDateTime(upload.UpdatedAt),
	}
}

// CreateUploadHandler creates a new upload via POST /documents/{documentID}/uploads
type CreateUploadHandler HandlerContext

// Handle creates a new Upload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {

	file, ok := params.File.(*runtime.File)
	if !ok {
		h.logger.Error("This should always be a runtime.File, something has changed in go-swagger.")
		return uploadop.NewCreateUploadInternalServerError()
	}
	h.logger.Info("File name and size: ", zap.String("name", file.Header.Filename), zap.Int64("size", file.Header.Size))

	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	reqApp := app.GetAppFromContext(params.HTTPRequest)

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		h.logger.Info("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	//fetching document to ensure user has access to it
	_, docErr := models.FetchDocument(h.db, user, reqApp, documentID)
	if docErr != nil {
		return responseForError(h.logger, docErr)
	}

	/*
		#nosec - we use md5 because it's required by the S3 API for
		validating data integrity.
		https://aws.amazon.com/premiumsupport/knowledge-center/data-integrity-s3/
	*/
	hash := md5.New()
	if _, err := io.Copy(hash, file.Data); err != nil {
		h.logger.Error("failed to hash uploaded file", zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}
	_, err = file.Data.Seek(0, io.SeekStart) // seek back to beginning of file
	if err != nil {
		h.logger.Error("failed to seek to beginning of uploaded file", zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	if file.Header.Size == 0 {
		h.logger.Error("File has a length of 0, aborting.")
		return uploadop.NewCreateUploadBadRequest()
	}

	buffer := make([]byte, 512)
	_, err = file.Data.Read(buffer)
	if err != nil {
		h.logger.Error("unable to read first 512 bytes of file", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	contentType := http.DetectContentType(buffer)

	_, err = file.Data.Seek(0, io.SeekStart) // seek back to beginning of file
	if err != nil {
		h.logger.Error("failed to seek to beginning of uploaded file", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	checksum := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	id := uuid.Must(uuid.NewV4())

	newUpload := models.Upload{
		ID:          id,
		DocumentID:  documentID,
		UploaderID:  user.ID,
		Filename:    file.Header.Filename,
		Bytes:       int64(file.Header.Size),
		ContentType: contentType,
		Checksum:    checksum,
	}

	// validate upload before pushing file to S3
	verrs, err := newUpload.Validate(h.db)
	if err != nil {
		h.logger.Error("Failed to validate", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	} else if verrs.HasAny() {
		payload := createFailedValidationPayload(verrs)
		return uploadop.NewCreateUploadBadRequest().WithPayload(payload)
	}

	// Push file to S3
	key := h.storage.Key("documents", documentID.String(), "uploads", id.String())
	_, err = h.storage.Store(key, file.Data, checksum)
	if err != nil {
		h.logger.Error("failed to store", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	// Already validated upload, so just save
	err = h.db.Create(&newUpload)
	if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	h.logger.Info("created an upload with id and key ", zap.Any("new_upload_id", newUpload.ID), zap.String("key", key))

	url, err := h.storage.PresignedURL(key, contentType)
	if err != nil {
		h.logger.Error("failed to get presigned url", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}
	uploadPayload := payloadForUploadModel(newUpload, url)
	return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload)
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler HandlerContext

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	app := app.GetAppFromContext(params.HTTPRequest)

	uploadID, _ := uuid.FromString(params.UploadID.String())
	upload, err := models.FetchUpload(h.db, user, app, uploadID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	key := h.storage.Key("documents", upload.DocumentID.String(), "uploads", upload.ID.String())
	err = h.storage.Delete(key)
	if err != nil {
		return responseForError(h.logger, err)
	}

	err = models.DeleteUpload(h.db, &upload)
	if err != nil {
		return responseForError(h.logger, err)
	}

	return uploadop.NewDeleteUploadCreated()
}

// DeleteUploadsHandler deletes a collection of uploads
type DeleteUploadsHandler HandlerContext

// Handle deletes uploads
func (h DeleteUploadsHandler) Handle(params uploadop.DeleteUploadsParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	app := app.GetAppFromContext(params.HTTPRequest)

	for _, uploadID := range params.UploadIds {
		uuid, _ := uuid.FromString(uploadID.String())
		upload, err := models.FetchUpload(h.db, user, app, uuid)
		if err != nil {
			return responseForError(h.logger, err)
		}

		key := h.storage.Key("documents", upload.DocumentID.String(), "uploads", upload.ID.String())
		err = h.storage.Delete(key)
		if err != nil {
			return responseForError(h.logger, err)
		}

		err = models.DeleteUpload(h.db, &upload)
		if err != nil {
			return responseForError(h.logger, err)
		}
	}

	return uploadop.NewDeleteUploadsCreated()
}
