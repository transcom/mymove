package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	auth "github.com/transcom/mymove/pkg/auth"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	publicdocumentsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/documents"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

/*
 * ------------------------------------------
 * The code below is for the INTERNAL REST API.
 * ------------------------------------------
 */

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
		ID:              fmtUUID(document.ID),
		ServiceMemberID: fmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

// CreateDocumentHandler creates a new document via POST /documents/
type CreateDocumentHandler HandlerContext

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	serviceMemberID, err := uuid.FromString(params.DocumentPayload.ServiceMemberID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Fetch to check auth
	serviceMember, err := models.FetchServiceMember(h.db, session, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	newDocument := models.Document{
		ServiceMemberID: serviceMember.ID,
	}

	verrs, err := h.db.ValidateAndCreate(&newDocument)
	if err != nil {
		h.logger.Info("DB Insertion", zap.Error(err))
		return documentop.NewCreateDocumentInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error("Could not save document", zap.String("errors", verrs.Error()))
		return documentop.NewCreateDocumentBadRequest()
	}

	h.logger.Info("created a document with id: ", zap.Any("new_document_id", newDocument.ID))
	documentPayload, err := payloadForDocumentModel(h.storage, newDocument)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return documentop.NewCreateDocumentCreated().WithPayload(documentPayload)
}

// ShowDocumentHandler shows a document via GETT /documents/:document_id
type ShowDocumentHandler HandlerContext

// Handle creates a new Document from a request payload
func (h ShowDocumentHandler) Handle(params documentop.ShowDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	documentID, err := uuid.FromString(params.DocumentID.String())
	if err != nil {
		return responseForError(h.logger, err)
	}

	document, err := models.FetchDocument(h.db, session, documentID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	documentPayload, err := payloadForDocumentModel(h.storage, document)
	if err != nil {
		return responseForError(h.logger, err)
	}

	return documentop.NewShowDocumentOK().WithPayload(documentPayload)
}

/* NOTE - The code above is for the INTERNAL API. The code below is for the public API. These will, obviously,
need to be reconciled. This will be done when the NotImplemented code below is Implemented
*/

// CreateDocumentUploadHandler creates a new document upload via POST /document/{document_uuid}/uploads
type CreateDocumentUploadHandler HandlerContext

// Handle creates a new DocumentUpload from a request payload
func (h CreateDocumentUploadHandler) Handle(params publicdocumentsop.CreateDocumentUploadParams) middleware.Responder {
	return middleware.NotImplemented("operation .createDocumentUpload has not yet been implemented")
}

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

// NO CODE YET!
