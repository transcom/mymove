package internalapi

import (
	"errors"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForMovingExpenseDocumentModel(storer storage.FileStorer, movingExpenseDocument models.MovingExpenseDocument) (*internalmessages.MoveDocumentPayload, error) {

	moveDocumentType := internalmessages.MoveDocumentType(movingExpenseDocument.MoveDocument.MoveDocumentType)
	status := internalmessages.MoveDocumentStatus(movingExpenseDocument.MoveDocument.Status)

	documentPayload, err := payloadForDocumentModel(storer, movingExpenseDocument.MoveDocument.Document)
	if err != nil {
		return nil, err
	}
	movingExpenseDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(movingExpenseDocument.MoveDocument.ID),
		MoveID:               handlers.FmtUUID(movingExpenseDocument.MoveDocument.MoveID),
		Document:             documentPayload,
		Title:                &movingExpenseDocument.MoveDocument.Title,
		MoveDocumentType:     &moveDocumentType,
		Status:               &status,
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
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
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
				userUpload, fetchUploadErr := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), convertedUploadID)
				if fetchUploadErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), fetchUploadErr)
				}
				userUploads = append(userUploads, userUpload)
			}

			var ppmID *uuid.UUID
			if payload.PersonallyProcuredMoveID != nil {
				id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

				// Enforce that the ppm's move_id matches our move
				ppm, fetchPPMErr := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), id)
				if fetchPPMErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), fetchPPMErr)
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
				RequestedAmountCents: unit.Cents(*payload.RequestedAmountCents),
				PaymentMethod:        *payload.PaymentMethod,
				ReceiptMissing:       payload.ReceiptMissing,
				StorageEndDate:       storageEndDate,
				StorageStartDate:     storageStartDate,
			}
			if payload.MovingExpenseType != nil {
				movingExpenseDocument.MovingExpenseType = models.MovingExpenseType(*payload.MovingExpenseType)
			}
			if payload.MoveDocumentType == nil {
				return handlers.ResponseForError(appCtx.Logger(), errors.New("missing required field: MoveDocumentType"))
			}
			newMovingExpenseDocument, verrs, err := move.CreateMovingExpenseDocument(
				appCtx.DB(),
				userUploads,
				ppmID,
				models.MoveDocumentType(*payload.MoveDocumentType),
				*payload.Title,
				payload.Notes,
				movingExpenseDocument,
				*move.SelectedMoveType,
			)

			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
			}

			newPayload, err := payloadForMovingExpenseDocumentModel(h.FileStorer(), *newMovingExpenseDocument)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}
			return movedocop.NewCreateMovingExpenseDocumentOK().WithPayload(newPayload)
		})
}
