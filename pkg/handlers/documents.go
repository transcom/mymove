package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	authctx "github.com/transcom/mymove/pkg/auth/context"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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

// CreateDocumentHandler creates a new document via POST /issue
type CreateDocumentHandler HandlerContext

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Fatal("No User ID, this should never happen.")
	}
	// userID := uuid.Must(uuid.FromString(params.MoveID.String()))

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}

	newDocument := models.Document{
		UploaderID: userID,
		MoveID:     moveID,
		Name:       params.DocumentPayload.Name,
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
	documentPayload := payloadForDocumentModel(newDocument)
	return documentop.NewCreateDocumentCreated().WithPayload(&documentPayload)
}
