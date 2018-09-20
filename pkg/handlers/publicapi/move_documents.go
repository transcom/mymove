package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	// "github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	movedocop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
	// "github.com/transcom/mymove/pkg/storage"
	// "github.com/transcom/mymove/pkg/unit"
)

// IndexMoveDocumentsHandler returns a list of all the Move Documents associated with this move.
type IndexMoveDocumentsHandler struct {
	handlers.HandlerContext
}

// Handle handles the request
func (h IndexMoveDocumentsHandler) Handle(params movedocop.IndexMoveDocumentsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the logged in TSP user exists
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return movedocop.NewCreateGenericMoveDocumentUnauthorized()
	}

	// Verify that TSP user is authorized to create movedoc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return movedocop.NewCreateGenericMoveDocumentForbidden()
	}
	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, shipment.Move.ID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	moveDocs, err := move.FetchAllMoveDocumentsForMove(h.DB())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	moveDocumentsPayload := make(apimessages.IndexMoveDocumentPayload, len(moveDocs))
	for i, doc := range moveDocs {
		documentPayload, err := payloadForDocumentModel(h.FileStorer(), doc.Document)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
		moveDocumentPayload := apimessages.MoveDocumentPayload{
			ID:               handlers.FmtUUID(doc.ID),
			ShipmentID:       handlers.FmtUUIDPtr(doc.ShipmentID),
			Document:         documentPayload,
			Title:            &doc.Title,
			MoveDocumentType: apimessages.MoveDocumentType(doc.MoveDocumentType),
			Status:           apimessages.MoveDocumentStatus(doc.Status),
			Notes:            doc.Notes,
		}
		// moveDocumentPayload, err := payloadForMoveDocumentExtractor(h.FileStorer(), doc)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
		moveDocumentsPayload[i] = &moveDocumentPayload
	}

	response := movedocop.NewIndexMoveDocumentsOK().WithPayload(moveDocumentsPayload)
	return response
}

// // UpdateMoveDocumentHandler updates a move document via PUT /moves/{moveId}/documents/{moveDocumentId}
// type UpdateMoveDocumentHandler struct {
// 	handlers.HandlerContext
// }

// // Handle ... updates a move document from a request payload
// func (h UpdateMoveDocumentHandler) Handle(params movedocop.UpdateMoveDocumentParams) middleware.Responder {
// 	session := auth.SessionFromRequestContext(params.HTTPRequest)

// 	moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())

// 	// Fetch move document from move id
// 	moveDoc, err := models.FetchMoveDocument(h.DB(), session, moveDocID)
// 	if err != nil {
// 		return handlers.ResponseForError(h.Logger(), err)
// 	}

// 	payload := params.UpdateMoveDocument
// 	if payload.PersonallyProcuredMoveID != nil {
// 		ppmID := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))
// 		moveDoc.PersonallyProcuredMoveID = &ppmID
// 	}
// 	newType := models.MoveDocumentType(payload.MoveDocumentType)
// 	moveDoc.Title = *payload.Title
// 	moveDoc.Notes = payload.Notes
// 	moveDoc.MoveDocumentType = newType

// 	newStatus := models.MoveDocumentStatus(payload.Status)

// 	// If this is a shipment summary and it has been approved, we process the ppm.
// 	if newStatus != moveDoc.Status {
// 		err = moveDoc.AttemptTransition(newStatus)
// 		if err != nil {
// 			return handlers.ResponseForError(h.Logger(), err)
// 		}

// 		if newStatus == models.MoveDocumentStatusOK && moveDoc.MoveDocumentType == models.MoveDocumentTypeSHIPMENTSUMMARY {
// 			if moveDoc.PersonallyProcuredMoveID == nil {
// 				return handlers.ResponseForError(h.Logger(), errors.New("No PPM loaded for Approved Move Doc"))
// 			}

// 			ppm := &moveDoc.PersonallyProcuredMove
// 			// If the status has already been completed
// 			// (because the document has been toggled between OK and HAS_ISSUE and back)
// 			// then don't complete it again.
// 			if ppm.Status != models.PPMStatusCOMPLETED {
// 				err := ppm.Complete()
// 				if err != nil {
// 					return handlers.ResponseForError(h.Logger(), err)
// 				}
// 			}
// 		}
// 	}

// 	var saveAction models.MoveDocumentSaveAction

// 	// If we are an expense type, we need to either delete, create, or update a MovingExpenseType
// 	// depending on which type of document already exists
// 	if models.IsExpenseModelDocumentType(newType) {
// 		// We should have a MovingExpenseDocument model
// 		requestedAmt := unit.Cents(payload.RequestedAmountCents)
// 		paymentMethod := payload.PaymentMethod
// 		if moveDoc.MovingExpenseDocument == nil {
// 			// But we don't have one, so create it to be saved later
// 			moveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
// 				MoveDocumentID:       moveDoc.ID,
// 				MoveDocument:         *moveDoc,
// 				MovingExpenseType:    models.MovingExpenseType(payload.MovingExpenseType),
// 				RequestedAmountCents: requestedAmt,
// 				PaymentMethod:        paymentMethod,
// 			}
// 		} else {
// 			// We have one already, so update the fields
// 			moveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(payload.MovingExpenseType)
// 			moveDoc.MovingExpenseDocument.RequestedAmountCents = requestedAmt
// 			moveDoc.MovingExpenseDocument.PaymentMethod = paymentMethod
// 		}
// 		saveAction = models.MoveDocumentSaveActionSAVEEXPENSEMODEL
// 	} else {
// 		if moveDoc.MovingExpenseDocument != nil {
// 			// We just care if a MovingExpenseType exists, as it needs to be deleted
// 			saveAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
// 		}
// 	}

// 	verrs, err := models.SaveMoveDocument(h.DB(), moveDoc, saveAction)
// 	if err != nil || verrs.HasAny() {
// 		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
// 	}

// 	moveDocPayload, err := payloadForMoveDocument(h.FileStorer(), *moveDoc)
// 	if err != nil {
// 		return handlers.ResponseForError(h.Logger(), err)
// 	}
// 	return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
// }
