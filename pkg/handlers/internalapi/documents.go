package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForDocumentModel(storer storage.FileStorer, document models.Document) (*internalmessages.DocumentPayload, error) {
	uploads := make([]*internalmessages.UploadPayload, len(document.UserUploads))
	for i, userUpload := range document.UserUploads {
		if userUpload.Upload.ID == uuid.Nil {
			return nil, errors.New("No uploads for user")
		}
		url, err := storer.PresignedURL(userUpload.Upload.StorageKey, userUpload.Upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := payloadForUploadModel(storer, userUpload.Upload, url)
		uploads[i] = uploadPayload
	}

	documentPayload := &internalmessages.DocumentPayload{
		ID:              handlers.FmtUUID(document.ID),
		ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

// CreateDocumentHandler creates a new document via POST /documents/
type CreateDocumentHandler struct {
	handlers.HandlerContext
}

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	serviceMemberID, err := uuid.FromString(params.DocumentPayload.ServiceMemberID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	// Fetch to check auth
	serviceMember, err := models.FetchServiceMemberForUser(h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	newDocument := models.Document{
		ServiceMemberID: serviceMember.ID,
	}

	verrs, err := h.DB().ValidateAndCreate(&newDocument)
	if err != nil {
		logger.Info("DB Insertion", zap.Error(err))
		return documentop.NewCreateDocumentInternalServerError()
	} else if verrs.HasAny() {
		logger.Error("Could not save document", zap.String("errors", verrs.Error()))
		return documentop.NewCreateDocumentBadRequest()
	}

	logger.Info("created a document with id", zap.Any("new_document_id", newDocument.ID))
	documentPayload, err := payloadForDocumentModel(h.FileStorer(), newDocument)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return documentop.NewCreateDocumentCreated().WithPayload(documentPayload)
}

// ShowDocumentHandler shows a document via GETT /documents/:document_id
type ShowDocumentHandler struct {
	handlers.HandlerContext
}

// Handle creates a new Document from a request payload
func (h ShowDocumentHandler) Handle(params documentop.ShowDocumentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	document, err := models.FetchDocument(h.DB(), session, documentID, false)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	documentPayload, err := payloadForDocumentModel(h.FileStorer(), document)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	return documentop.NewShowDocumentOK().WithPayload(documentPayload)
}
