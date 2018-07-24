package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
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
		Title:            &moveDocument.Title,
		MoveDocumentType: internalmessages.MoveDocumentType(moveDocument.MoveDocumentType),
		Status:           internalmessages.MoveDocumentStatus(moveDocument.Status),
		Notes:            moveDocument.Notes,
	}

	return &moveDocumentPayload, nil
}

// CreateMoveDocumentHandler creates a MoveDocument
type CreateMoveDocumentHandler HandlerContext

// Handle is the handler
func (h CreateMoveDocumentHandler) Handle(params movedocop.CreateMoveDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreateMoveDocumentPayload

	// Fetch uploads to confirm ownership
	uploadIds := payload.UploadIds
	if len(uploadIds) == 0 {
		return movedocop.NewCreateMoveDocumentBadRequest()
	}

	uploads := models.Uploads{}
	for _, id := range uploadIds {
		converted := uuid.Must(uuid.FromString(id.String()))
		upload, err := models.FetchUpload(h.db, session, converted)
		if err != nil {
			return responseForError(h.logger, err)
		}
		uploads = append(uploads, upload)
	}

	newMoveDocument, verrs, err := move.CreateMoveDocument(h.db,
		uploads,
		models.MoveDocumentType(payload.MoveDocumentType),
		*payload.Title,
		models.MoveDocumentStatus(payload.Status),
		payload.Notes)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	newPayload, err := payloadForMoveDocumentModel(h.storage, *newMoveDocument)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return movedocop.NewCreateMoveDocumentOK().WithPayload(newPayload)
}

// UpdateMoveDocumentHandler updates a move document via PUT /moves/{moveId}/documents/{moveDocumentId}
type UpdateMoveDocumentHandler HandlerContext

// Handle ... updates a move document from a request payload
func (h UpdateMoveDocumentHandler) Handle(params movedocop.UpdateMoveDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

	// Fetch move document from move id
	moveDoc, err := models.FetchMoveDocument(h.db, session, moveDocID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.UpdateMoveDocument
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes

	verrs, err := models.SaveMoveDocument(h.db, moveDoc)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	moveDocPayload, err := payloadForMoveDocumentModel(h.storage, *moveDoc)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
}

// IndexMoveDocumentsHandler returns a list of all the Move Documents associated with this move.
type IndexMoveDocumentsHandler HandlerContext

// Handle handles the request
func (h IndexMoveDocumentsHandler) Handle(params movedocop.IndexMoveDocumentsParams) middleware.Responder {
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
	response := movedocop.NewIndexMoveDocumentsOK().WithPayload(moveDocumentsPayload)
	return response
}
