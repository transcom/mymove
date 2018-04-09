package handlers

import (
	"crypto/md5"
	"encoding/base64"
	// "fmt"
	"io"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	authctx "github.com/transcom/mymove/pkg/auth/context"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForUploadModel(upload models.Upload, url string) internalmessages.UploadPayload {
	return internalmessages.UploadPayload{
		ID:       fmtUUID(upload.ID),
		Filename: swag.String(upload.Filename),
		URL:      fmtURI(url),
	}
}

// CreateUploadHandler creates a new upload via POST /moves/{moveID}/documents/{documentID}/uploads
type CreateUploadHandler FileHandlerContext

// Handle creates a new Upload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {

	file, ok := params.File.(*runtime.File)
	if !ok {
		h.logger.Error("This should always be a runtime.File, something has changed in go-swagger.")
		return uploadop.NewCreateUploadInternalServerError()
	}
	h.logger.Info("File name and size: ", zap.String("name", file.Header.Filename), zap.Int64("size", file.Header.Size))

	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Error("Missing User ID in context")
		return uploadop.NewCreateUploadBadRequest()
	}

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Error("Badly formed UUID for moveId", zap.String("move_id", params.MoveID.String()), zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		h.logger.Error("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	// Validate that the document and move exists in the db, and that they belong to user
	exists, userOwns := models.ValidateDocumentOwnership(h.db, userID, moveID, documentID)
	if !exists {
		h.logger.Error("document or move does not exist", zap.String("document_id", params.DocumentID.String()), zap.String("move_id", params.MoveID.String()), zap.Error(err))
		return uploadop.NewCreateUploadNotFound()
	}
	if !userOwns {
		h.logger.Error("user does not own document or move", zap.String("document_id", params.DocumentID.String()), zap.String("move_id", params.MoveID.String()), zap.Error(err))
		return uploadop.NewCreateUploadForbidden()
	}

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
		UploaderID:  userID,
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
		h.logger.Error(verrs.Error())
		return uploadop.NewCreateUploadBadRequest()
	}

	// Push file to S3
	key := h.storage.Key("moves", moveID.String(), "documents", documentID.String(), "uploads", id.String())
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
	return uploadop.NewCreateUploadCreated().WithPayload(&uploadPayload)
}
