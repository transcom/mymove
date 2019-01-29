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

	"github.com/transcom/mymove/pkg/auth"
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

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	file, ok := params.File.(*runtime.File)
	if !ok {
		h.Logger().Error("This should always be a runtime.File, something has changed in go-swagger.")
		return uploadop.NewCreateUploadInternalServerError()
	}

	h.Logger().Info("File name and size: ", zap.String("name", file.Header.Filename), zap.Int64("size", file.Header.Size))

	// User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	var docID *uuid.UUID
	if params.DocumentID != nil {
		documentID, err := uuid.FromString(params.DocumentID.String())
		if err != nil {
			h.Logger().Info("Badly formed UUID for document", zap.String("document_id", params.DocumentID.String()), zap.Error(err))
			return uploadop.NewCreateUploadBadRequest()
		}

		// Fetch document to ensure user has access to it
		document, docErr := models.FetchDocument(ctx, h.DB(), session, documentID)
		if docErr != nil {
			return handlers.ResponseForError(h.Logger(), docErr)
		}
		docID = &document.ID
	}

	// Read the incoming data into a temporary afero.File for consumption
	aFile, err := h.FileStorer().TempFileSystem().Create(file.Header.Filename)
	if err != nil {
		h.Logger().Error("Error opening afero file.", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	_, err = io.Copy(aFile, file.Data)
	if err != nil {
		h.Logger().Error("Error copying incoming data into afero file.", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}

	uploader := uploaderpkg.NewUploader(h.DB(), h.Logger(), h.FileStorer())
	newUpload, verrs, err := uploader.CreateUpload(docID, session.UserID, aFile)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	url, err := uploader.PresignedURL(newUpload)
	if err != nil {
		h.Logger().Error("failed to get presigned url", zap.Error(err))
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

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	uploadID, _ := uuid.FromString(params.UploadID.String())
	upload, err := models.FetchUpload(ctx, h.DB(), session, uploadID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	uploader := uploaderpkg.NewUploader(h.DB(), h.Logger(), h.FileStorer())
	if err = uploader.DeleteUpload(&upload); err != nil {
		return handlers.ResponseForError(h.Logger(), err)
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	uploader := uploaderpkg.NewUploader(h.DB(), h.Logger(), h.FileStorer())

	for _, uploadID := range params.UploadIds {
		uuid, _ := uuid.FromString(uploadID.String())
		upload, err := models.FetchUpload(ctx, h.DB(), session, uuid)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}

		if err = uploader.DeleteUpload(&upload); err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
	}

	return uploadop.NewDeleteUploadsNoContent()
}
