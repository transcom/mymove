package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/go-openapi/swag"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	movedocop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"go.uber.org/zap"
)

func payloadForDocumentModel(storer storage.FileStorer, document models.Document) (*apimessages.DocumentPayload, error) {
	uploads := make([]*apimessages.UploadPayload, len(document.Uploads))
	for i, upload := range document.Uploads {
		url, err := storer.PresignedURL(upload.StorageKey, upload.ContentType)
		if err != nil {
			return nil, err
		}

		// uploadPayload := payloadForUploadModel(upload, url)

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
		ID: handlers.FmtUUID(document.ID),
		// ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads: uploads,
	}
	return documentPayload, nil
}

func payloadForGenericMoveDocumentModel(storer storage.FileStorer, moveDocument models.MoveDocument) (*apimessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDocument.Document)
	if err != nil {
		return nil, err
	}

	genericMoveDocumentPayload := apimessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the logged in TSP user exists
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return movedocop.NewCreateGenericMoveDocumentUnauthorized()
	}

	// Verify that TSP user is authorized to create movedoc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return movedocop.NewCreateGenericMoveDocumentForbidden()
	}

	// #nosec UUID is pattern matched by swagger and will be ok
	// moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, shipment.Move.ID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := params.CreateGenericMoveDocumentPayload

	// Fetch uploads to confirm ownership
	uploadIds := payload.UploadIds
	if len(uploadIds) == 0 {
		return movedocop.NewCreateGenericMoveDocumentBadRequest()
	}

	uploads := models.Uploads{}
	for _, id := range uploadIds {
		converted := uuid.Must(uuid.FromString(id.String()))
		upload, err := models.FetchUpload(h.DB(), session, converted)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
		uploads = append(uploads, upload)
	}

	// var ppmID *uuid.UUID
	// if payload.PersonallyProcuredMoveID != nil {
	// 	id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

	// 	// Enforce that the ppm's move_id matches our move
	// 	ppm, err := models.FetchPersonallyProcuredMove(h.DB(), session, id)
	// 	if err != nil {
	// 		return handlers.ResponseForError(h.Logger(), err)
	// 	}
	// 	if !uuid.Equal(ppm.MoveID, moveID) {
	// 		return movedocop.NewCreateGenericMoveDocumentBadRequest()
	// 	}

	// 	ppmID = &id
	// }

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

	newPayload, err := payloadForGenericMoveDocumentModel(h.FileStorer(), *newMoveDocument)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return movedocop.NewCreateGenericMoveDocumentOK().WithPayload(newPayload)
}
