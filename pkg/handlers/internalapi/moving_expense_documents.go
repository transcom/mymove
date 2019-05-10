package internalapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/auth"
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
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
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
		upload, fetchUploadErr := models.FetchUpload(ctx, h.DB(), session, converted)
		if fetchUploadErr != nil {
			return handlers.ResponseForError(h.Logger(), fetchUploadErr)
		}
		uploads = append(uploads, upload)
	}

	var ppmID *uuid.UUID
	if payload.PersonallyProcuredMoveID != nil {
		id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

		// Enforce that the ppm's move_id matches our move
		ppm, fetchPPMErr := models.FetchPersonallyProcuredMove(h.DB(), session, id)
		if fetchPPMErr != nil {
			return handlers.ResponseForError(h.Logger(), fetchPPMErr)
		}
		if ppm.MoveID != moveID {
			return movedocop.NewCreateMovingExpenseDocumentBadRequest()
		}

		ppmID = &id
	}

	newMovingExpenseDocument, verrs, err := move.CreateMovingExpenseDocument(
		h.DB(),
		uploads,
		ppmID,
		models.MoveDocumentType(payload.MoveDocumentType),
		*payload.Title,
		payload.Notes,
		unit.Cents(*payload.RequestedAmountCents),
		*payload.PaymentMethod,
		models.MovingExpenseType(payload.MovingExpenseType),
		*move.SelectedMoveType,
	)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	newPayload, err := payloadForMovingExpenseDocumentModel(h.FileStorer(), *newMovingExpenseDocument)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return movedocop.NewCreateMovingExpenseDocumentOK().WithPayload(newPayload)
}
