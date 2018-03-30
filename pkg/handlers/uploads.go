package handlers

import (
	"crypto/md5"
	"encoding/base64"
	"io"

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
	file := params.File
	h.logger.Infof("%s has a length of %d bytes.\n", file.Header.Filename, file.Header.Size)

	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Panic("No User ID, this should never happen.")
	}

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Panic("Invalid MoveID, this should never happen.")
	}

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		h.logger.Panic("Invalid DocumentID, this should never happen.")
	}

	docExists, moveExists, userOwns := models.ValidateDocumentOwnership(h.db, userID, moveID, documentID)
	if !docExists || !moveExists {
		return uploadop.NewCreateUploadNotFound()
	}
	if !userOwns {
		return uploadop.NewCreateUploadForbidden()
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file.Data); err != nil {
		h.logger.Panic("failed to hash uploaded file", zap.Error(err))
	}
	_, err = file.Data.Seek(0, io.SeekStart) // seek back to beginning of file
	if err != nil {
		h.logger.Panic("failed to seek to beginning of uploaded file", zap.Error(err))
	}

	checksum := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	newUpload := models.Upload{
		DocumentID: documentID,
		UploaderID: userID,
		Filename:   file.Header.Filename,
		Bytes:      int64(file.Header.Size),
		// TODO replace this with a real content type by examining file content.
		ContentType: "text/plain",
		Checksum:    checksum,
	}

	verrs, err := h.db.ValidateAndCreate(&newUpload)
	if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error(verrs.Error())
		return uploadop.NewCreateUploadBadRequest()
	} else {
		h.logger.Infof("created an upload with id %s, s3 id %s\n", newUpload.ID, newUpload.ID)

		key := h.storage.Key("moves", moveID.String(), "documents", documentID.String(), "uploads", newUpload.ID.String())

		_, err := h.storage.Store(key, file.Data, checksum)
		if err != nil {
			h.logger.Error("failed to store", zap.Error(err))
			return uploadop.NewCreateUploadInternalServerError()
		}

		url, err := h.storage.PresignedURL(key)
		if err != nil {
			h.logger.Error("failed to get presigned url", zap.Error(err))
			return uploadop.NewCreateUploadInternalServerError()
		}
		uploadPayload := payloadForUploadModel(newUpload, url)
		return uploadop.NewCreateUploadCreated().WithPayload(&uploadPayload)
	}
}
