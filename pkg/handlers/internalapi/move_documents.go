package internalapi

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForMoveDocument(storer storage.FileStorer, moveDoc models.MoveDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDoc.Document)
	if err != nil {
		return nil, err
	}

	moveDocumentType := internalmessages.MoveDocumentType(moveDoc.MoveDocumentType)
	status := internalmessages.MoveDocumentStatus(moveDoc.Status)

	payload := internalmessages.MoveDocumentPayload{
		ID:                       handlers.FmtUUID(moveDoc.ID),
		MoveID:                   handlers.FmtUUID(moveDoc.MoveID),
		PersonallyProcuredMoveID: handlers.FmtUUIDPtr(moveDoc.PersonallyProcuredMoveID),
		Document:                 documentPayload,
		Title:                    &moveDoc.Title,
		MoveDocumentType:         &moveDocumentType,
		Status:                   &status,
		Notes:                    moveDoc.Notes,
	}

	if moveDoc.MovingExpenseDocument != nil {
		payload.MovingExpenseType = internalmessages.MovingExpenseType(moveDoc.MovingExpenseDocument.MovingExpenseType)
		payload.RequestedAmountCents = int64(moveDoc.MovingExpenseDocument.RequestedAmountCents)
		payload.PaymentMethod = moveDoc.MovingExpenseDocument.PaymentMethod
	}

	if moveDoc.MovingExpenseDocument != nil && moveDoc.MovingExpenseDocument.MovingExpenseType == models.MovingExpenseTypeSTORAGE {
		if moveDoc.MovingExpenseDocument.StorageStartDate != nil {
			payload.StorageStartDate = handlers.FmtDate(*moveDoc.MovingExpenseDocument.StorageStartDate)
		}
		if moveDoc.MovingExpenseDocument.StorageEndDate != nil {
			payload.StorageEndDate = handlers.FmtDate(*moveDoc.MovingExpenseDocument.StorageEndDate)
		}
	}

	if moveDoc.WeightTicketSetDocument != nil {
		if moveDoc.WeightTicketSetDocument.EmptyWeight != nil {
			payload.EmptyWeight = handlers.FmtInt64(int64(*moveDoc.WeightTicketSetDocument.EmptyWeight))
		}
		if moveDoc.WeightTicketSetDocument.FullWeight != nil {
			payload.FullWeight = handlers.FmtInt64(int64(*moveDoc.WeightTicketSetDocument.FullWeight))
		}
		if moveDoc.WeightTicketSetDocument.VehicleNickname != nil {
			payload.VehicleNickname = moveDoc.WeightTicketSetDocument.VehicleNickname
		}
		if moveDoc.WeightTicketSetDocument.VehicleMake != nil {
			payload.VehicleMake = moveDoc.WeightTicketSetDocument.VehicleMake
		}
		if moveDoc.WeightTicketSetDocument.VehicleModel != nil {
			payload.VehicleModel = moveDoc.WeightTicketSetDocument.VehicleModel
		}
		weightTicketSetType := internalmessages.WeightTicketSetType(moveDoc.WeightTicketSetDocument.WeightTicketSetType)
		payload.WeightTicketSetType = &weightTicketSetType
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
	var paymentMethod string
	if docExtractor.PaymentMethod != nil {
		paymentMethod = string(*docExtractor.PaymentMethod)
	}
	var requestedAmt unit.Cents
	if docExtractor.RequestedAmountCents != nil {
		requestedAmt = *docExtractor.RequestedAmountCents
	}
	var emptyWeight *int64
	if docExtractor.EmptyWeight != nil {
		emptyWeight = handlers.FmtInt64(int64(*docExtractor.EmptyWeight))
	}
	var emptyWeightTicketMissing *bool
	if docExtractor.EmptyWeightTicketMissing != nil {
		emptyWeightTicketMissing = docExtractor.EmptyWeightTicketMissing
	}
	var fullWeight *int64
	if docExtractor.FullWeight != nil {
		fullWeight = handlers.FmtInt64(int64(*docExtractor.FullWeight))
	}
	var fullWeightTicketMissing *bool
	if docExtractor.FullWeightTicketMissing != nil {
		fullWeightTicketMissing = docExtractor.FullWeightTicketMissing
	}
	var vehicleNickname string
	if docExtractor.VehicleNickname != nil {
		vehicleNickname = *docExtractor.VehicleNickname
	}
	var vehicleMake string
	if docExtractor.VehicleMake != nil {
		vehicleMake = *docExtractor.VehicleMake
	}
	var vehicleModel string
	if docExtractor.VehicleModel != nil {
		vehicleModel = *docExtractor.VehicleModel
	}
	var weightTicketDate *strfmt.Date
	if docExtractor.WeightTicketDate != nil {
		weightTicketDate = handlers.FmtDate(*docExtractor.WeightTicketDate)
	}
	var trailerOwnershipMissing *bool
	if docExtractor.TrailerOwnershipMissing != nil {
		trailerOwnershipMissing = docExtractor.TrailerOwnershipMissing
	}
	var receiptMissing *bool
	if docExtractor.ReceiptMissing != nil {
		receiptMissing = docExtractor.ReceiptMissing
	}
	var storageStartDate *strfmt.Date
	if docExtractor.StorageStartDate != nil {
		storageStartDate = handlers.FmtDate(*docExtractor.StorageStartDate)
	}
	var storageEndDate *strfmt.Date
	if docExtractor.StorageEndDate != nil {
		storageEndDate = handlers.FmtDate(*docExtractor.StorageEndDate)
	}

	moveDocumentType := internalmessages.MoveDocumentType(docExtractor.MoveDocumentType)
	status := internalmessages.MoveDocumentStatus(docExtractor.Status)

	payload := internalmessages.MoveDocumentPayload{
		ID:                       handlers.FmtUUID(docExtractor.ID),
		MoveID:                   handlers.FmtUUID(docExtractor.MoveID),
		PersonallyProcuredMoveID: handlers.FmtUUIDPtr(docExtractor.PersonallyProcuredMoveID),
		Document:                 documentPayload,
		Title:                    &docExtractor.Title,
		MoveDocumentType:         &moveDocumentType,
		Status:                   &status,
		Notes:                    docExtractor.Notes,
		MovingExpenseType:        expenseType,
		RequestedAmountCents:     int64(requestedAmt),
		PaymentMethod:            paymentMethod,
		ReceiptMissing:           receiptMissing,
		VehicleNickname:          &vehicleNickname,
		VehicleMake:              &vehicleMake,
		VehicleModel:             &vehicleModel,
		EmptyWeight:              emptyWeight,
		EmptyWeightTicketMissing: emptyWeightTicketMissing,
		FullWeight:               fullWeight,
		FullWeightTicketMissing:  fullWeightTicketMissing,
		WeightTicketDate:         weightTicketDate,
		TrailerOwnershipMissing:  trailerOwnershipMissing,
		StorageStartDate:         storageStartDate,
		StorageEndDate:           storageEndDate,
	}

	if docExtractor.WeightTicketSetType != nil {
		weightTicketSetType := internalmessages.WeightTicketSetType(*docExtractor.WeightTicketSetType)
		payload.WeightTicketSetType = &weightTicketSetType
	}

	return &payload, nil
}

// IndexMoveDocumentsHandler returns a list of all the Move Documents associated with this move.
type IndexMoveDocumentsHandler struct {
	handlers.HandlerConfig
}

// Handle handles the request
func (h IndexMoveDocumentsHandler) Handle(params movedocop.IndexMoveDocumentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID, err := uuid.FromString(params.MoveID.String())

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			moveDocs, err := move.FetchAllMoveDocumentsForMove(appCtx.DB(), false)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			moveDocumentsPayload := make(internalmessages.MoveDocuments, len(moveDocs))
			for i, doc := range moveDocs {
				moveDocumentPayload, err := payloadForMoveDocumentExtractor(h.FileStorer(), doc)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				moveDocumentsPayload[i] = moveDocumentPayload
			}

			return movedocop.NewIndexMoveDocumentsOK().WithPayload(moveDocumentsPayload), nil
		})
}

// UpdateMoveDocumentHandler updates a move document via PUT /moves/{moveId}/documents/{moveDocumentId}
type UpdateMoveDocumentHandler struct {
	handlers.HandlerConfig
	moveDocumentUpdater services.MoveDocumentUpdater
}

// Handle ... updates a move document from a request payload
func (h UpdateMoveDocumentHandler) Handle(params movedocop.UpdateMoveDocumentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

			moveDoc, verrs, err := h.moveDocumentUpdater.Update(appCtx, params.UpdateMoveDocument, moveDocID)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}
			moveDocPayload, err := payloadForMoveDocument(h.FileStorer(), *moveDoc)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload), nil
		})
}

// DeleteMoveDocumentHandler deletes a move document via DELETE /moves/{moveId}/documents/{moveDocumentId}
type DeleteMoveDocumentHandler struct {
	handlers.HandlerConfig
}

// Handle ... deletes a move document
func (h DeleteMoveDocumentHandler) Handle(params movedocop.DeleteMoveDocumentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

			// for now, only delete if weight ticket set or expense
			moveDoc, err := models.FetchMoveDocument(appCtx.DB(), appCtx.Session(), moveDocID, false)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			err = models.DeleteMoveDocument(appCtx.DB(), moveDoc)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return movedocop.NewDeleteMoveDocumentNoContent(), nil
		})
}
