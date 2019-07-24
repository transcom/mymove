package internalapi

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/pkg/errors"

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

	payload := internalmessages.MoveDocumentPayload{
		ID:                       handlers.FmtUUID(moveDoc.ID),
		MoveID:                   handlers.FmtUUID(moveDoc.MoveID),
		PersonallyProcuredMoveID: handlers.FmtUUIDPtr(moveDoc.PersonallyProcuredMoveID),
		Document:                 documentPayload,
		Title:                    &moveDoc.Title,
		MoveDocumentType:         internalmessages.MoveDocumentType(moveDoc.MoveDocumentType),
		Status:                   internalmessages.MoveDocumentStatus(moveDoc.Status),
		Notes:                    moveDoc.Notes,
	}

	if moveDoc.MovingExpenseDocument != nil {
		payload.MovingExpenseType = internalmessages.MovingExpenseType(moveDoc.MovingExpenseDocument.MovingExpenseType)
		payload.RequestedAmountCents = int64(moveDoc.MovingExpenseDocument.RequestedAmountCents)
		payload.PaymentMethod = moveDoc.MovingExpenseDocument.PaymentMethod
	}

	if moveDoc.WeightTicketSetDocument != nil {
		if moveDoc.WeightTicketSetDocument.EmptyWeight != nil {
			payload.EmptyWeight = handlers.FmtInt64(int64(*moveDoc.WeightTicketSetDocument.EmptyWeight))
		}
		if moveDoc.WeightTicketSetDocument.FullWeight != nil {
			payload.FullWeight = handlers.FmtInt64(int64(*moveDoc.WeightTicketSetDocument.FullWeight))
		}
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
	var vehicleOptions string
	if docExtractor.VehicleOptions != nil {
		vehicleOptions = *docExtractor.VehicleOptions
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

	payload := internalmessages.MoveDocumentPayload{
		ID:                       handlers.FmtUUID(docExtractor.ID),
		MoveID:                   handlers.FmtUUID(docExtractor.MoveID),
		PersonallyProcuredMoveID: handlers.FmtUUIDPtr(docExtractor.PersonallyProcuredMoveID),
		Document:                 documentPayload,
		Title:                    &docExtractor.Title,
		MoveDocumentType:         internalmessages.MoveDocumentType(docExtractor.MoveDocumentType),
		Status:                   internalmessages.MoveDocumentStatus(docExtractor.Status),
		Notes:                    docExtractor.Notes,
		MovingExpenseType:        expenseType,
		RequestedAmountCents:     int64(requestedAmt),
		PaymentMethod:            paymentMethod,
		ReceiptMissing:           receiptMissing,
		VehicleOptions:           vehicleOptions,
		VehicleNickname:          vehicleNickname,
		EmptyWeight:              emptyWeight,
		EmptyWeightTicketMissing: emptyWeightTicketMissing,
		FullWeight:               fullWeight,
		FullWeightTicketMissing:  fullWeightTicketMissing,
		WeightTicketDate:         weightTicketDate,
		TrailerOwnershipMissing:  trailerOwnershipMissing,
		StorageStartDate:         storageStartDate,
		StorageEndDate:           storageEndDate,
	}

	return &payload, nil
}

// IndexMoveDocumentsHandler returns a list of all the Move Documents associated with this move.
type IndexMoveDocumentsHandler struct {
	handlers.HandlerContext
}

// Handle handles the request
func (h IndexMoveDocumentsHandler) Handle(params movedocop.IndexMoveDocumentsParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	moveDocs, err := move.FetchAllMoveDocumentsForMove(h.DB())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	moveDocumentsPayload := make(internalmessages.MoveDocuments, len(moveDocs))
	for i, doc := range moveDocs {
		moveDocumentPayload, err := payloadForMoveDocumentExtractor(h.FileStorer(), doc)
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}
		moveDocumentsPayload[i] = moveDocumentPayload
	}

	response := movedocop.NewIndexMoveDocumentsOK().WithPayload(moveDocumentsPayload)
	return response
}

// UpdateMoveDocumentHandler updates a move document via PUT /moves/{moveId}/documents/{moveDocumentId}
type UpdateMoveDocumentHandler struct {
	handlers.HandlerContext
}

// Handle ... updates a move document from a request payload
func (h UpdateMoveDocumentHandler) Handle(params movedocop.UpdateMoveDocumentParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

	// Fetch move document from move id
	moveDoc, err := models.FetchMoveDocument(h.DB(), session, moveDocID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := params.UpdateMoveDocument
	if payload.PersonallyProcuredMoveID != nil {
		ppmID := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))
		moveDoc.PersonallyProcuredMoveID = &ppmID
	}
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes
	moveDoc.MoveDocumentType = newType

	newStatus := models.MoveDocumentStatus(payload.Status)
	oldStatus := moveDoc.Status

	if newStatus != oldStatus {
		err = moveDoc.AttemptTransition(newStatus)
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}

		// If this is a shipment summary and it has been approved, we process the ppm.
		if newStatus == models.MoveDocumentStatusOK && moveDoc.MoveDocumentType == models.MoveDocumentTypeSHIPMENTSUMMARY {
			if moveDoc.PersonallyProcuredMoveID == nil {
				return handlers.ResponseForError(logger, errors.New("No PPM loaded for Approved Move Doc"))
			}

			ppm := &moveDoc.PersonallyProcuredMove
			// If the status has already been completed
			// (because the document has been toggled between OK and HAS_ISSUE and back)
			// then don't complete it again.
			if ppm.Status != models.PPMStatusCOMPLETED {
				completeErr := ppm.Complete()
				if completeErr != nil {
					return handlers.ResponseForError(logger, completeErr)
				}
			}
		}

		// If this is a storage expense, we need to make changes to the total sit amount for the ppm.
		if moveDoc.MoveDocumentType == models.MoveDocumentTypeEXPENSE && moveDoc.MovingExpenseDocument.MovingExpenseType == models.MovingExpenseTypeSTORAGE {
			if moveDoc.PersonallyProcuredMoveID == nil {
				return handlers.ResponseForError(logger, errors.New("No PPM loaded for Approved Move Doc"))
			}

			ppm := &moveDoc.PersonallyProcuredMove
			storageRequestedAmt := unit.Cents(payload.RequestedAmountCents)
			var newCost unit.Cents

			// add to SIT total amount if OK
			if newStatus == models.MoveDocumentStatusOK {
				if ppm.TotalSITCost == nil {
					newCost = storageRequestedAmt
				} else {
					newCost = *ppm.TotalSITCost + storageRequestedAmt
				}
			}

			// subtract from SIT total amount if changed from OK
			if oldStatus == models.MoveDocumentStatusOK && newStatus != models.MoveDocumentStatusOK {
				newCost = *ppm.TotalSITCost - storageRequestedAmt
			}

			ppm.TotalSITCost = &newCost
		}
	}

	var saveExpenseAction models.MoveExpenseDocumentSaveAction

	// If we are an expense type, we need to either delete, create, or update a MovingExpenseType
	// depending on which type of document already exists
	if models.IsExpenseModelDocumentType(newType) {
		// We should have a MovingExpenseDocument model
		requestedAmt := unit.Cents(payload.RequestedAmountCents)
		paymentMethod := payload.PaymentMethod
		if moveDoc.MovingExpenseDocument == nil {
			// But we don't have one, so create it to be saved later
			moveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
				MoveDocumentID:       moveDoc.ID,
				MoveDocument:         *moveDoc,
				MovingExpenseType:    models.MovingExpenseType(payload.MovingExpenseType),
				RequestedAmountCents: requestedAmt,
				PaymentMethod:        paymentMethod,
			}
		} else {
			// We have one already, so update the fields
			moveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(payload.MovingExpenseType)
			moveDoc.MovingExpenseDocument.RequestedAmountCents = requestedAmt
			moveDoc.MovingExpenseDocument.PaymentMethod = paymentMethod
		}
		saveExpenseAction = models.MoveDocumentSaveActionSAVEEXPENSEMODEL
	} else {
		if moveDoc.MovingExpenseDocument != nil {
			// We just care if a MovingExpenseType exists, as it needs to be deleted
			saveExpenseAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
		}
	}

	var saveWeightTicketSetAction models.MoveWeightTicketSetDocumentSaveAction

	// If we are a weight ticket set type, we need to either delete, create, or update a WeightTicketSetDocument
	// depending on which type of document already exists
	if newType == models.MoveDocumentTypeWEIGHTTICKETSET {
		var emptyWeight, fullWeight *unit.Pound
		if payload.EmptyWeight != nil {
			ew := unit.Pound(*payload.EmptyWeight)
			emptyWeight = &ew
		}
		if payload.FullWeight != nil {
			fw := unit.Pound(*payload.FullWeight)
			fullWeight = &fw
		}
		var weightTicketDate *time.Time
		if payload.WeightTicketDate != nil {
			weightTicketDate = (*time.Time)(payload.WeightTicketDate)
		}
		var trailerOwnershipMissing bool
		if payload.TrailerOwnershipMissing != nil {
			trailerOwnershipMissing = *payload.TrailerOwnershipMissing
		}

		if moveDoc.WeightTicketSetDocument == nil {
			// create new weight ticket set
			moveDoc.WeightTicketSetDocument = &models.WeightTicketSetDocument{
				MoveDocumentID:          moveDoc.ID,
				MoveDocument:            *moveDoc,
				EmptyWeight:             emptyWeight,
				FullWeight:              fullWeight,
				VehicleNickname:         payload.VehicleNickname,
				VehicleOptions:          payload.VehicleOptions,
				WeightTicketDate:        weightTicketDate,
				TrailerOwnershipMissing: trailerOwnershipMissing,
			}
		} else {
			// update existing weight ticket set
			moveDoc.WeightTicketSetDocument.EmptyWeight = emptyWeight
			moveDoc.WeightTicketSetDocument.FullWeight = fullWeight
		}
		saveWeightTicketSetAction = models.MoveDocumentSaveActionSAVEWEIGHTTICKETSETMODEL
	} else {
		// delete if document exists but the move document is being converted to something else
		if moveDoc.WeightTicketSetDocument != nil {
			saveWeightTicketSetAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
		}
	}

	verrs, err := models.SaveMoveDocument(h.DB(), moveDoc, saveExpenseAction, saveWeightTicketSetAction)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	moveDocPayload, err := payloadForMoveDocument(h.FileStorer(), *moveDoc)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
}
