package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
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
		ID:                   handlers.FmtUUID(movingExpenseDocument.MoveDocument.ID),
		MoveID:               handlers.FmtUUID(movingExpenseDocument.MoveDocument.MoveID),
		Document:             documentPayload,
		Title:                &movingExpenseDocument.MoveDocument.Title,
		MoveDocumentType:     internalmessages.MoveDocumentType(movingExpenseDocument.MoveDocument.MoveDocumentType),
		Status:               internalmessages.MoveDocumentStatus(movingExpenseDocument.MoveDocument.Status),
		Notes:                movingExpenseDocument.MoveDocument.Notes,
		MovingExpenseType:    internalmessages.MovingExpenseType(movingExpenseDocument.MovingExpenseType),
		RequestedAmountCents: int64(movingExpenseDocument.RequestedAmountCents),
		PaymentMethod:        movingExpenseDocument.PaymentMethod,
	}

	return &movingExpenseDocumentPayload, nil
}

// CreateMovingExpenseDocumentHandler creates a MovingExpenseDocument
type CreateMovingExpenseDocumentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CreateMovingExpenseDocumentHandler) Handle(params movedocop.CreateMovingExpenseDocumentParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID, nil)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := params.CreateMovingExpenseDocumentPayload

	uploadIds := payload.UploadIds
	haveReceipt := !payload.ReceiptMissing
	// To maintain old behavior that required / assumed
	// that users always had receipts
	if len(uploadIds) == 0 && haveReceipt {
		return movedocop.NewCreateMovingExpenseDocumentBadRequest()
	}

	// Fetch uploads to confirm ownership
	userUploads := models.UserUploads{}
	for _, id := range uploadIds {
		convertedUploadID := uuid.Must(uuid.FromString(id.String()))
		userUpload, fetchUploadErr := models.FetchUserUploadFromUploadID(ctx, h.DB(), session, convertedUploadID)
		if fetchUploadErr != nil {
			return handlers.ResponseForError(logger, fetchUploadErr)
		}
		userUploads = append(userUploads, userUpload)
	}

	var ppmID *uuid.UUID
	if payload.PersonallyProcuredMoveID != nil {
		id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

		// Enforce that the ppm's move_id matches our move
		ppm, fetchPPMErr := models.FetchPersonallyProcuredMove(h.DB(), session, id)
		if fetchPPMErr != nil {
			return handlers.ResponseForError(logger, fetchPPMErr)
		}
		if ppm.MoveID != moveID {
			return movedocop.NewCreateMovingExpenseDocumentBadRequest()
		}

		ppmID = &id
	}

	var storageStartDate *time.Time
	if payload.StorageStartDate != nil {
		storageStartDate = (*time.Time)(payload.StorageStartDate)
	}
	var storageEndDate *time.Time
	if payload.StorageEndDate != nil {
		storageEndDate = (*time.Time)(payload.StorageEndDate)
	}
	movingExpenseDocument := models.MovingExpenseDocument{
		MovingExpenseType:    models.MovingExpenseType(payload.MovingExpenseType),
		RequestedAmountCents: unit.Cents(*payload.RequestedAmountCents),
		PaymentMethod:        *payload.PaymentMethod,
		ReceiptMissing:       payload.ReceiptMissing,
		StorageEndDate:       storageEndDate,
		StorageStartDate:     storageStartDate,
	}
	newMovingExpenseDocument, verrs, err := move.CreateMovingExpenseDocument(
		h.DB(),
		userUploads,
		ppmID,
		models.MoveDocumentType(payload.MoveDocumentType),
		*payload.Title,
		payload.Notes,
		movingExpenseDocument,
		*move.SelectedMoveType,
	)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	newPayload, err := payloadForMovingExpenseDocumentModel(h.FileStorer(), *newMovingExpenseDocument)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return movedocop.NewCreateMovingExpenseDocumentOK().WithPayload(newPayload)
}
