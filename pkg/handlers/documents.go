package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	authctx "github.com/transcom/mymove/pkg/auth"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDocumentModel(document models.Document) internalmessages.DocumentPayload {
	documentPayload := internalmessages.DocumentPayload{
		ID:      fmtUUID(document.ID),
		Name:    swag.String(document.Name),
		Uploads: []*internalmessages.UploadPayload{},
	}
	return documentPayload
}

// CreateDocumentHandler creates a new document via POST /documents/
type CreateDocumentHandler HandlerContext

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Error("Missing User ID in context")
		return documentop.NewCreateDocumentBadRequest()
	}

	serviceMemberID, err := uuid.FromString(params.DocumentPayload.ServiceMemberID.String())
	if err != nil {
		h.logger.Info("Badly formed UUID for serviceMemberId", zap.String("service_member_id", params.DocumentPayload.ServiceMemberID.String()), zap.Error(err))
		return documentop.NewCreateDocumentBadRequest()
	}

	newDocument := models.Document{
		UploaderID:      userID,
		ServiceMemberID: serviceMemberID,
		Name:            params.DocumentPayload.Name,
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
	documentPayload := payloadForDocumentModel(newDocument)
	return documentop.NewCreateDocumentCreated().WithPayload(&documentPayload)
}

/* NOTE - The code above is for the INTERNAL API. The code below is for the public API. These will, obviously,
need to be reconciled. This will be done when the NotImplemented code below is Implemented
*/

// CreateDocumentUploadHandler creates a new document upload via POST /document/{document_uuid}/uploads
type CreateDocumentUploadHandler HandlerContext

// Handle creates a new DocumentUpload from a request payload
func (h CreateDocumentUploadHandler) Handle(params apioperations.CreateDocumentUploadParams) middleware.Responder {
	return middleware.NotImplemented("operation .createDocumentUpload has not yet been implemented")
}
