package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

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

	// Fetch uploads to confirm ownership
	uploadIds := payload.UploadIds
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
	return moveop.NewCreateMoveDocumentOK().WithPayload(newPayload)
}

// UpdateMoveDocumentHandler updates a move document via PUT /moves/{moveId}/documents/{moveDocumentId}
type UpdateMoveDocumentHandler HandlerContext

// Handle ... updates a move document from a request payload
func (h UpdateMoveDocumentHandler) Handle(params moveop.UpdateMoveDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	moveID, _ := uuid.FromString(params.MoveID.String())
	moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

	// Fetch move document from move id
	moveDoc, err := models.FetchMoveDocument(h.db, session, moveDocID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	if moveDoc.MoveID != moveID {
		h.logger.Info("Move ID for Move Document does not match requested Move Document ID", zap.String("requested move_id", moveID.String()), zap.String("actual move_id", moveDoc.MoveID.String()))
		return moveop.NewUpdateMoveDocumentBadRequest()
	}
	payload := params.UpdateMoveDocument
	moveDoc.Title = *payload.Title
	moveDoc.Status = models.MoveDocumentStatus(payload.Status)
	moveDoc.Notes = payload.Notes

	verrs, err := models.SaveMoveDocument(h.db, moveDoc)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	moveDocPayload, err := payloadForMoveDocumentModel(h.storage, *moveDoc)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return moveop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
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
