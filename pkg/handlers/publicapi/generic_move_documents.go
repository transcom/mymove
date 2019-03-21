package publicapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	movedocop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForDocumentModel(storer storage.FileStorer, document models.Document) (*apimessages.DocumentPayload, error) {
	uploads := make([]*apimessages.UploadPayload, len(document.Uploads))
	for i, upload := range document.Uploads {
		url, err := storer.PresignedURL(upload.StorageKey, upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := &apimessages.UploadPayload{
			ID:          handlers.FmtUUID(upload.ID),
			Filename:    swag.String(upload.Filename),
			ContentType: swag.String(upload.ContentType),
			URL:         handlers.FmtURI(url),
			Bytes:       &upload.Bytes,
			CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
			UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
		}
		uploads[i] = uploadPayload
	}

	documentPayload := &apimessages.DocumentPayload{
		ID:      handlers.FmtUUID(document.ID),
		Uploads: uploads,
	}
	return documentPayload, nil
}

func payloadForGenericMoveDocumentModel(storer storage.FileStorer, moveDocument models.MoveDocument, shipmentID uuid.UUID) (*apimessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDocument.Document)
	if err != nil {
		return nil, err
	}

	genericMoveDocumentPayload := apimessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		ShipmentID:       handlers.FmtUUID(shipmentID),
		Document:         documentPayload,
		Title:            &moveDocument.Title,
		MoveDocumentType: apimessages.MoveDocumentType(moveDocument.MoveDocumentType),
		Status:           apimessages.MoveDocumentStatus(moveDocument.Status),
		Notes:            moveDocument.Notes,
	}

	return &genericMoveDocumentPayload, nil
}

// CreateGenericMoveDocumentHandler creates a MoveDocument
type CreateGenericMoveDocumentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CreateGenericMoveDocumentHandler) Handle(params movedocop.CreateGenericMoveDocumentParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the TSP user is authorized to update move doc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	_, shipment, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
	if err != nil {
		if err.Error() == "USER_UNAUTHORIZED" {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		if err.Error() == "FETCH_FORBIDDEN" {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Fetch move
	move, err := models.FetchMove(h.DB(), session, shipment.Move.ID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := params.CreateGenericMoveDocumentPayload

	// Fetch uploads
	uploadIds := payload.UploadIds
	if len(uploadIds) == 0 {
		return movedocop.NewCreateGenericMoveDocumentBadRequest()
	}

	uploads := models.Uploads{}
	for _, id := range uploadIds {
		converted := uuid.Must(uuid.FromString(id.String()))
		upload, err := models.FetchUpload(ctx, h.DB(), session, converted)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
		uploads = append(uploads, upload)
	}

	newMoveDocument, verrs, err := move.CreateMoveDocument(h.DB(),
		uploads,
		&shipmentID,
		models.MoveDocumentType(payload.MoveDocumentType),
		*payload.Title,
		payload.Notes,
		*shipment.Move.SelectedMoveType)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	newPayload, err := payloadForGenericMoveDocumentModel(h.FileStorer(), *newMoveDocument, shipmentID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return movedocop.NewCreateGenericMoveDocumentOK().WithPayload(newPayload)
}
