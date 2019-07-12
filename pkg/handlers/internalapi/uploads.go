package internalapi

import (
	"io"
	"reflect"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/honeycombio/beeline-go"
	"go.uber.org/zap"

	uploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
)

func payloadForUploadModel(upload models.Upload, url string) *internalmessages.UploadPayload {
	return &internalmessages.UploadPayload{
		ID:          handlers.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         handlers.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
	}
}

// CreateUploadHandler creates a new upload via POST /documents/{documentID}/uploads
type CreateUploadHandler struct {
	handlers.HandlerContext
}

// Handle creates a new Upload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromContext(ctx)

	ctx, span := beeline.StartSpan(ctx, reflect.TypeOf(h).Name())
	defer span.Send()

	file, ok := params.File.(*runtime.File)
	if !ok {
		logger.Error("This should always be a runtime.File, something has changed in go-swagger.")
		return uploadop.NewCreateUploadInternalServerError()
	}

	logger.Info(
		"File name and size",
		zap.String("name", file.Header.Filename),
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
		document, docErr := models.FetchDocument(ctx, h.DB(), session, documentID)
		if docErr != nil {
			return handlers.ResponseForError(logger, docErr)
		}
		docID = &document.ID
	}

	// Read the incoming data into a temporary afero.File for consumption
	aFile, err := h.FileStorer().TempFileSystem().Create(file.Header.Filename)
	if err != nil {
		logger.Error("Error opening afero file.", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	_, err = io.Copy(aFile, file.Data)
	if err != nil {
		logger.Error("Error copying incoming data into afero file.", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	uploader := uploaderpkg.NewUploader(h.DB(), logger, h.FileStorer())
	newUpload, verrs, err := uploader.CreateUploadForDocument(docID, session.UserID, aFile, uploaderpkg.AllowedTypesServiceMember)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	url, err := uploader.PresignedURL(newUpload)
	if err != nil {
		logger.Error("failed to get presigned url", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}
	uploadPayload := payloadForUploadModel(*newUpload, url)
	return uploadop.NewCreateUploadCreated().WithPayload(uploadPayload)
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler struct {
	handlers.HandlerContext
}

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	uploadID, _ := uuid.FromString(params.UploadID.String())
	upload, err := models.FetchUpload(ctx, h.DB(), session, uploadID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	uploader := uploaderpkg.NewUploader(h.DB(), logger, h.FileStorer())
	if err = uploader.DeleteUpload(&upload); err != nil {
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

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	// User should always be populated by middleware
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	uploader := uploaderpkg.NewUploader(h.DB(), logger, h.FileStorer())

	for _, uploadID := range params.UploadIds {
		uuid, _ := uuid.FromString(uploadID.String())
		upload, err := models.FetchUpload(ctx, h.DB(), session, uuid)
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}

		if err = uploader.DeleteUpload(&upload); err != nil {
			return handlers.ResponseForError(logger, err)
		}
	}

	return uploadop.NewDeleteUploadsNoContent()
}
