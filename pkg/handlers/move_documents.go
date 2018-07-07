package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForMoveDocumentModel(storer storage.FileStorer, moveDocument models.MoveDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDocument.Document)
	if err != nil {
		return nil, err
	}

	moveDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:               fmtUUID(moveDocument.ID),
		MoveID:           fmtUUID(moveDocument.MoveID),
		Document:         documentPayload,
		MoveDocumentType: internalmessages.MoveDocumentType(moveDocument.MoveDocumentType),
		Status:           internalmessages.MoveDocumentStatus(moveDocument.Status),
		Notes:            moveDocument.Notes,
	}

	return &moveDocumentPayload, nil
}

// CreateMoveDocumentHandler creates a MoveDocument
type CreateMoveDocumentHandler HandlerContext

// Handle is the handler
func (h CreateMoveDocumentHandler) Handle(params moveop.CreateMoveDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreateMoveDocumentPayload

	// Also validates access to the document
	documentID := uuid.Must(uuid.FromString(payload.DocumentID.String()))
	document, err := models.FetchDocument(h.db, session, documentID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	newMoveDocument, verrs, err := move.CreateMoveDocument(h.db,
		document,
		models.MoveDocumentType(payload.MoveDocumentType),
		models.MoveDocumentStatus(payload.Status),
		payload.Notes)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	newPayload, err := payloadForMoveDocumentModel(h.storage, *newMoveDocument)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewCreateMoveDocumentOK().WithPayload(newPayload)
}

// IndexMoveDocumentsHandler returns a list of all the Move Documents associated with this move.
type IndexMoveDocumentsHandler HandlerContext

// Handle handles the request
func (h IndexMoveDocumentsHandler) Handle(params moveop.IndexMoveDocumentsParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Fetch move documents on move documents model
	moveDocuments := move.MoveDocuments
	moveDocumentsPayload := make(internalmessages.IndexMoveDocumentPayload, len(moveDocuments))
	for i, moveDocument := range moveDocuments {
		moveDocumentPayload, err := payloadForMoveDocumentModel(h.storage, moveDocument)
		if err != nil {
			return responseForError(h.logger, err)
		}
		moveDocumentsPayload[i] = moveDocumentPayload
	}
	response := moveop.NewIndexMoveDocumentsOK().WithPayload(moveDocumentsPayload)
	return response
}
