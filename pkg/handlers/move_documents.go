package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForMoveDocument(storer storage.FileStorer, moveDoc models.MoveDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDoc.Document)
	if err != nil {
		return nil, err
	}

	payload := internalmessages.MoveDocumentPayload{
		ID:               fmtUUID(moveDoc.ID),
		MoveID:           fmtUUID(moveDoc.MoveID),
		Document:         documentPayload,
		Title:            &moveDoc.Title,
		MoveDocumentType: internalmessages.MoveDocumentType(moveDoc.MoveDocumentType),
		Status:           internalmessages.MoveDocumentStatus(moveDoc.Status),
		Notes:            moveDoc.Notes,
	}

	if moveDoc.MovingExpenseDocument != nil {
		payload.MovingExpenseType = internalmessages.MovingExpenseType(moveDoc.MovingExpenseDocument.MovingExpenseType)
		payload.Reimbursement = payloadForReimbursementModel(&moveDoc.MovingExpenseDocument.Reimbursement)
	}

	return &payload, nil
}

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

	newType := models.MoveDocumentType(payload.MoveDocumentType)
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes
	moveDoc.Status = models.MoveDocumentStatus(payload.Status)
	moveDoc.MoveDocumentType = newType

	var saveAction models.MoveDocumentSaveAction

	// If we are an expense type, we need to either delete, create, or update a MovingExpenseType
	// depending on which type of document already exists
	if models.IsExpenseModelDocumentType(newType) {
		// We should have a MovingExpenseDocument model
		reimbursementAmt := unit.Cents(*payload.Reimbursement.RequestedAmount)
		reimbursementMethod := models.MethodOfReceipt(*payload.Reimbursement.MethodOfReceipt)
		if moveDoc.MovingExpenseDocument == nil {
			// But we don't have one, so create it to be saved later
			reimbursement := models.BuildRequestedReimbursement(
				reimbursementAmt,
				reimbursementMethod)
			moveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
				MoveDocumentID:    moveDoc.ID,
				MoveDocument:      *moveDoc,
				MovingExpenseType: models.MovingExpenseType(payload.MovingExpenseType),
				Reimbursement:     reimbursement,
			}
		} else {
			// We have one already, so update the fields
			moveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(payload.MovingExpenseType)
			moveDoc.MovingExpenseDocument.Reimbursement.RequestedAmount = reimbursementAmt
			moveDoc.MovingExpenseDocument.Reimbursement.MethodOfReceipt = reimbursementMethod
		}
		saveAction = models.MoveDocumentSaveActionSAVE_EXPENSE_MODEL
	} else {
		if moveDoc.MovingExpenseDocument != nil {
			// We just care if a MovingExpenseType exists, as it needs to be deleted
			saveAction = models.MoveDocumentSaveActionDELETE_EXPENSE_MODEL
		}
	}

	verrs, err := models.SaveMoveDocument(h.db, moveDoc, saveAction)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	moveDocPayload, err := payloadForMoveDocument(h.storage, *moveDoc)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
}
