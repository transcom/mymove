package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	moveDocs, err := move.FetchAllMoveDocumentsForMove(h.DB())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	moveDocumentsPayload := make(internalmessages.MoveDocuments, len(moveDocs))
	for i, doc := range moveDocs {
		moveDocumentPayload, err := payloadForMoveDocumentExtractor(h.FileStorer(), doc)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

	// Fetch move document from move id
	moveDoc, err := models.FetchMoveDocument(h.DB(), session, moveDocID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
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

	// If this is a shipment summary and it has been approved, we process the ppm.
	if newStatus != moveDoc.Status {
		err = moveDoc.AttemptTransition(newStatus)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}

		if newStatus == models.MoveDocumentStatusOK && moveDoc.MoveDocumentType == models.MoveDocumentTypeSHIPMENTSUMMARY {
			if moveDoc.PersonallyProcuredMoveID == nil {
				return handlers.ResponseForError(h.Logger(), errors.New("No PPM loaded for Approved Move Doc"))
			}

			ppm := &moveDoc.PersonallyProcuredMove
			// If the status has already been completed
			// (because the document has been toggled between OK and HAS_ISSUE and back)
			// then don't complete it again.
			if ppm.Status != models.PPMStatusCOMPLETED {
				err := ppm.Complete()
				if err != nil {
					return handlers.ResponseForError(h.Logger(), err)
				}
			}
		}
	}

	var saveAction models.MoveDocumentSaveAction

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
		saveAction = models.MoveDocumentSaveActionSAVEEXPENSEMODEL
	} else {
		if moveDoc.MovingExpenseDocument != nil {
			// We just care if a MovingExpenseType exists, as it needs to be deleted
			saveAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
		}
	}

	verrs, err := models.SaveMoveDocument(h.DB(), moveDoc, saveAction)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	moveDocPayload, err := payloadForMoveDocument(h.FileStorer(), *moveDoc)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
}
