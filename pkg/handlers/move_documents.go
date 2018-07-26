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

func payloadForMoveDocumentExtractor(storer storage.FileStorer, docExtractor models.MoveDocumentExtractor) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, docExtractor.Document)
	if err != nil {
		return nil, err
	}

	var expenseType internalmessages.MovingExpenseType
	if docExtractor.MovingExpenseType != nil {
		expenseType = internalmessages.MovingExpenseType(*docExtractor.MovingExpenseType)
	}

	payload := internalmessages.MoveDocumentPayload{
		ID:                fmtUUID(docExtractor.ID),
		MoveID:            fmtUUID(docExtractor.MoveID),
		Document:          documentPayload,
		Title:             &docExtractor.Title,
		MoveDocumentType:  internalmessages.MoveDocumentType(docExtractor.MoveDocumentType),
		Status:            internalmessages.MoveDocumentStatus(docExtractor.Status),
		Notes:             docExtractor.Notes,
		MovingExpenseType: expenseType,
		Reimbursement:     payloadForReimbursementModel(&docExtractor.Reimbursement),
	}

	return &payload, nil
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

	moveDocs, err := move.FetchAllMoveDocumentsForMove(h.db)
	if err != nil {
		return responseForError(h.logger, err)
	}

	moveDocumentsPayload := make(internalmessages.IndexMoveDocumentPayload, len(moveDocs))
	for i, doc := range moveDocs {
		moveDocumentPayload, err := payloadForMoveDocumentExtractor(h.storage, doc)
		if err != nil {
			return responseForError(h.logger, err)
		}
		moveDocumentsPayload[i] = moveDocumentPayload
	}

	response := movedocop.NewIndexMoveDocumentsOK().WithPayload(moveDocumentsPayload)
	return response
}
