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

func payloadForMovingExpenseDocumentModel(storer storage.FileStorer, movingExpenseDocument models.MovingExpenseDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, movingExpenseDocument.MoveDocument.Document)
	if err != nil {
		return nil, err
	}

	movingExpenseDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:                fmtUUID(movingExpenseDocument.MoveDocument.ID),
		MoveID:            fmtUUID(movingExpenseDocument.MoveDocument.MoveID),
		Document:          documentPayload,
		Title:             &movingExpenseDocument.MoveDocument.Title,
		MoveDocumentType:  internalmessages.MoveDocumentType(movingExpenseDocument.MoveDocument.MoveDocumentType),
		Status:            internalmessages.MoveDocumentStatus(movingExpenseDocument.MoveDocument.Status),
		Notes:             movingExpenseDocument.MoveDocument.Notes,
		MovingExpenseType: internalmessages.MovingExpenseType(movingExpenseDocument.MovingExpenseType),
		Reimbursement:     payloadForReimbursementModel(&movingExpenseDocument.Reimbursement),
	}

	return &movingExpenseDocumentPayload, nil
}

// CreateMovingExpenseDocumentHandler creates a MovingExpenseDocument
type CreateMovingExpenseDocumentHandler HandlerContext

// Handle is the handler
func (h CreateMovingExpenseDocumentHandler) Handle(params movedocop.CreateMovingExpenseDocumentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreateMovingExpenseDocumentPayload

	// Fetch uploads to confirm ownership
	uploadIds := payload.UploadIds
	if len(uploadIds) == 0 {
		return movedocop.NewCreateMovingExpenseDocumentBadRequest()
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

	var ppmID *uuid.UUID
	if payload.PersonallyProcuredMoveID != nil {
		id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

		// Enforce that the ppm's move_id matches our move
		ppm, err := models.FetchPersonallyProcuredMove(h.db, session, id)
		if err != nil {
			return responseForError(h.logger, err)
		}
		if !uuid.Equal(ppm.MoveID, moveID) {
			return movedocop.NewCreateMovingExpenseDocumentBadRequest()
		}

		ppmID = &id
	}

	reimbursement := models.BuildRequestedReimbursement(unit.Cents(*payload.Reimbursement.RequestedAmount), models.MethodOfReceipt(*payload.Reimbursement.MethodOfReceipt))

	newMovingExpenseDocument, verrs, err := move.CreateMovingExpenseDocument(
		h.db,
		uploads,
		ppmID,
		models.MoveDocumentType(payload.MoveDocumentType),
		*payload.Title,
		payload.Notes,
		reimbursement,
		models.MovingExpenseType(payload.MovingExpenseType),
	)

	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	newPayload, err := payloadForMovingExpenseDocumentModel(h.storage, *newMovingExpenseDocument)
	if err != nil {
		return responseForError(h.logger, err)
	}
	return movedocop.NewCreateMovingExpenseDocumentOK().WithPayload(newPayload)
}
