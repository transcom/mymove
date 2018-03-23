package handlers

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	authctx "github.com/transcom/mymove/pkg/auth/context"
	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForUploadModel(upload models.Upload) internalmessages.UploadPayload {
	return internalmessages.UploadPayload{
		ID:       fmtUUID(upload.ID),
		Filename: swag.String(upload.Filename),
		URL:      fmtURI("https://domain.text/file.ext"),
	}
}

// CreateUploadHandler creates a new upload via POST /issue
type CreateUploadHandler S3HandlerContext

// Handle creates a new Upload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	file := params.File
	if params.File == nil {
		// TODO Can swagger handle this check?
		return uploadop.NewCreateUploadBadRequest()
	}

	fmt.Printf("%s has a length of %d bytes.\n", file.Header.Filename, file.Header.Size)

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

	bucket := os.Getenv("AWS_S3_BUCKET_NAME")
	if len(bucket) == 0 {
		h.logger.Error("AWS_S3_BUCKET_NAME not configured")
		return uploadop.NewCreateUploadInternalServerError()
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file.Data); err != nil {
		h.logger.Panic("failed to hash uploaded file", zap.Error(err))
	}
	_, err = file.Data.Seek(0, io.SeekStart) // seek back to beginning of file
	if err != nil {
		h.logger.Panic("failed to seek to beginning of uploaded file", zap.Error(err))
	}

	newUpload := models.Upload{
		DocumentID:  documentID,
		UploaderID:  userID,
		Filename:    file.Header.Filename,
		Bytes:       int64(file.Header.Size),
		ContentType: "application/pdf",
		Checksum:    base64.StdEncoding.EncodeToString(hash.Sum(nil)),
	}

	verrs, err := h.db.ValidateAndCreate(&newUpload)
	if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error(verrs.Error())
		return uploadop.NewCreateUploadBadRequest()
	} else {
		fmt.Printf("created an upload with id %s, s3 id %s\n", newUpload.ID, newUpload.ID)

		key := fmt.Sprintf("moves/%s/documents/%s/uploads/%s", moveID, documentID, newUpload.ID)

		input := &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &key,
			Body:   file.Data,
		}
		_, err = h.s3.PutObject(input)
		if err != nil {
			h.logger.Error("PutObject failed")
			return uploadop.NewCreateUploadInternalServerError()
		}

		// TODO verify checksum

		uploadPayload := payloadForUploadModel(newUpload)
		return uploadop.NewCreateUploadCreated().WithPayload(&uploadPayload)
	}
}
