package internalapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/honeycombio/beeline-go"

	auth "github.com/transcom/mymove/pkg/auth"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForDocumentModel(storer storage.FileStorer, document models.Document) (*internalmessages.DocumentPayload, error) {
	uploads := make([]*internalmessages.UploadPayload, len(document.Uploads))
	for i, upload := range document.Uploads {
		url, err := storer.PresignedURL(upload.StorageKey, upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := payloadForUploadModel(upload, url)
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
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	serviceMemberID, err := uuid.FromString(params.DocumentPayload.ServiceMemberID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Fetch to check auth
	serviceMember, err := models.FetchServiceMemberForUser(ctx, h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	newDocument := models.Document{
		ServiceMemberID: serviceMember.ID,
	}

	verrs, err := h.DB().ValidateAndCreate(&newDocument)
	if err != nil {
		h.Logger().Info("DB Insertion", zap.Error(err))
		return documentop.NewCreateDocumentInternalServerError()
	} else if verrs.HasAny() {
		h.Logger().Error("Could not save document", zap.String("errors", verrs.Error()))
		return documentop.NewCreateDocumentBadRequest()
	}

	h.Logger().Info("created a document with id: ", zap.Any("new_document_id", newDocument.ID))
	documentPayload, err := payloadForDocumentModel(h.FileStorer(), newDocument)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return documentop.NewCreateDocumentCreated().WithPayload(documentPayload)
}

// ShowDocumentHandler shows a document via GETT /documents/:document_id
type ShowDocumentHandler struct {
	handlers.HandlerContext
}

// Handle creates a new Document from a request payload
func (h ShowDocumentHandler) Handle(params documentop.ShowDocumentParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	document, err := models.FetchDocument(ctx, h.DB(), session, documentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	documentPayload, err := payloadForDocumentModel(h.FileStorer(), document)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	return documentop.NewShowDocumentOK().WithPayload(documentPayload)
}
