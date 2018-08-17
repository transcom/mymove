package internal

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForGenericMoveDocumentModel(storer storage.FileStorer, moveDocument models.MoveDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDocument.Document)
	if err != nil {
		return nil, err
	}

	genericMoveDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:               utils.FmtUUID(moveDocument.ID),
		MoveID:           utils.FmtUUID(moveDocument.MoveID),
		Document:         documentPayload,
		Title:            &moveDocument.Title,
		MoveDocumentType: internalmessages.MoveDocumentType(moveDocument.MoveDocumentType),
		Status:           internalmessages.MoveDocumentStatus(moveDocument.Status),
		Notes:            moveDocument.Notes,
	}

	return &genericMoveDocumentPayload, nil
}

// CreateGenericMoveDocumentHandler creates a MoveDocument
type CreateGenericMoveDocumentHandler utils.HandlerContext

// Handle is the handler
func (h CreateGenericMoveDocumentHandler) Handle(params movedocop.CreateGenericMoveDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
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
		upload, err := models.FetchUpload(h.db, session, converted)
		if err != nil {
			return utils.ResponseForError(h.logger, err)
		}
		uploads = append(uploads, upload)
	}

	var ppmID *uuid.UUID
	if payload.PersonallyProcuredMoveID != nil {
		id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

		// Enforce that the ppm's move_id matches our move
		ppm, err := models.FetchPersonallyProcuredMove(h.db, session, id)
		if err != nil {
			return utils.ResponseForError(h.logger, err)
		}
		if !uuid.Equal(ppm.MoveID, moveID) {
			return movedocop.NewCreateGenericMoveDocumentBadRequest()
		}

		ppmID = &id
	}

	newMoveDocument, verrs, err := move.CreateMoveDocument(h.db,
		uploads,
		ppmID,
		models.MoveDocumentType(payload.MoveDocumentType),
		*payload.Title,
		payload.Notes)

	if err != nil || verrs.HasAny() {
		return utils.ResponseForVErrors(h.logger, verrs, err)
	}

	newPayload, err := payloadForGenericMoveDocumentModel(h.Storage, *newMoveDocument)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}
	return movedocop.NewCreateGenericMoveDocumentOK().WithPayload(newPayload)
}
