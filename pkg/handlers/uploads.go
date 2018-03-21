package handlers

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/markbates/pop"
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
	pop.Debug = true
	file := params.File

	fmt.Printf("%s has a length of %d bytes.\n", file.Header.Filename, file.Header.Size)

	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Fatal("No User ID, this should never happen.")
	}

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		h.logger.Fatal("Invalid DocumentID, this should never happen.")
	}
	// 	cwd, err := os.Getwd()
	// 	if err != nil {
	// 		h.logger.Error("Could not get cwd", zap.Error(err))
	// 	}

	// 	uploadsDir := filepath.Join(cwd, "uploads")
	// 	if err = os.Mkdir(uploadsDir, 0777); err != nil {
	// 		h.logger.Error("Could not make directory", zap.Error(err))
	// 	}

	// 	destinationPath := filepath.Join(uploadsDir, file.Header.Filename)
	// 	destination, err := os.Create(destinationPath)
	// 	defer destination.Close()

	// 	if err != nil {
	// 		h.logger.Error("Could on open file", zap.Error(err))
	// 	}

	bucket := os.Getenv("AWS_S3_BUCKET_NAME")
	if len(bucket) == 0 {
		h.logger.Error("AWS_S3_BUCKET_NAME not configured")
		return uploadop.NewCreateUploadInternalServerError()
	}

	uploadID := uuid.Must(uuid.NewV4())
	key := fmt.Sprintf("moves/%s/documents/%s/uploads/%s", moveID, documentID, uploadID)

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

	fmt.Printf("%s has a length of %d bytes.\n", file.Header.Filename, file.Header.Size)

	var response middleware.Responder

	newUpload := models.Upload{
		DocumentID:  documentID,
		UploaderID:  userID,
		Filename:    file.Header.Filename,
		Bytes:       int64(file.Header.Size),
		ContentType: "application/pdf",
		Checksum:    "abcdefg",
		S3ID:        uploadID,
	}

	verrs, err := h.db.ValidateAndCreate(&newUpload)
	if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error(verrs.Error())
		return uploadop.NewCreateUploadBadRequest()
	} else {
		fmt.Printf("created an upload with id %s\n", newUpload.ID)
		uploadPayload := payloadForUploadModel(newUpload)
		response = uploadop.NewCreateUploadCreated().WithPayload(&uploadPayload)
	}
	return response
}
