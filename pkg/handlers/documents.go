package handlers

import (
	"fmt"
	"os"

	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/markbates/pop"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	authctx "github.com/transcom/mymove/pkg/auth/context"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDocumentModel(document models.Document, upload models.Upload) internalmessages.DocumentPayload {
	uploadPayload := &internalmessages.UploadPayload{
		ID:       fmtUUID(upload.ID),
		Filename: swag.String(upload.Filename),
		URL:      swag.String("download URL"),
	}
	uploads := []*internalmessages.UploadPayload{uploadPayload}
	documentPayload := internalmessages.DocumentPayload{
		ID:      fmtUUID(document.ID),
		Uploads: uploads,
	}
	return documentPayload
}

var s3Client *s3.S3

func init() {
	session := awsSession.Must(awsSession.NewSession())
	s3Client = s3.New(session)
}

// CreateDocumentHandler creates a new document via POST /issue
type CreateDocumentHandler HandlerContext

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
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
		// TODO should be a server error
		return documentop.NewCreateDocumentBadRequest()
	}

	uploadID := uuid.Must(uuid.NewV4())
	key := fmt.Sprintf("uploads/moves/%s/%s", moveID, uploadID)

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file.Data,
	}
	_, err = s3Client.PutObject(input)
	if err != nil {
		h.logger.Error("AWS_S3_BUCKET_NAME not configured")
		// TODO should be a server error
		return documentop.NewCreateDocumentBadRequest()
	}

	fmt.Printf("%s has a length of %d bytes.\n", file.Header.Filename, file.Header.Size)

	var response middleware.Responder

	newDocument := models.Document{
		UploaderID: userID,
		MoveID:     moveID,
		Name:       "test document",
	}

	verrs, err := h.db.ValidateAndCreate(&newDocument)
	if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		return documentop.NewCreateDocumentInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error(verrs.Error())
		return documentop.NewCreateDocumentBadRequest()
	}

	fmt.Printf("created a document with id %s\n", newDocument.ID)

	newUpload := models.Upload{
		DocumentID:  newDocument.ID,
		UploaderID:  userID,
		Filename:    file.Header.Filename,
		Bytes:       int64(file.Header.Size),
		ContentType: "application/pdf",
		Checksum:    "abcdefg",
		S3ID:        uploadID,
	}

	verrs, err = h.db.ValidateAndCreate(&newUpload)
	if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		return documentop.NewCreateDocumentInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error(verrs.Error())
		return documentop.NewCreateDocumentBadRequest()
	} else {
		fmt.Printf("created an upload with id %s\n", newUpload.ID)
		documentPayload := payloadForDocumentModel(newDocument, newUpload)
		response = documentop.NewCreateDocumentCreated().WithPayload(&documentPayload)
	}
	return response
}
