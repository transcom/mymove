package publicapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	movedocop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// IndexMoveDocumentsHandler returns a list of all the Move Documents associated with this move.
type IndexMoveDocumentsHandler struct {
	handlers.HandlerContext
}

// Handle handles the request
func (h IndexMoveDocumentsHandler) Handle(params movedocop.IndexMoveDocumentsParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the TSP user is authorized to update move doc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	_, shipment, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
	if err != nil {
		if err.Error() == "USER_UNAUTHORIZED" {
			h.Logger().Error("DB Query", zap.Error(err))
			return h.RespondAndTraceError(ctx, err, "error user unauthorized")
		}
		if err.Error() == "FETCH_FORBIDDEN" {
			h.Logger().Error("DB Query", zap.Error(err))
			return h.RespondAndTraceError(ctx, err, "error fetch forbiddewn")
		}
		return h.RespondAndTraceError(ctx, err, "error fetching shipment", zap.String("shipmnt_id", shipmentID.String()))
	}

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.DB(), session, shipment.Move.ID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error user unauthorized")
	}

	moveDocs, err := move.FetchAllMoveDocumentsForMove(h.DB())
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching move")
	}

	moveDocumentsPayload := make(apimessages.MoveDocuments, len(moveDocs))
	for i, doc := range moveDocs {
		documentPayload, err := payloadForDocumentModel(h.FileStorer(), doc.Document)
		if err != nil {
			return h.RespondAndTraceError(ctx, err, "error fetching payload for document")
		}
		moveDocumentPayload := apimessages.MoveDocumentPayload{
			ID:               handlers.FmtUUID(doc.ID),
			ShipmentID:       handlers.FmtUUIDPtr(doc.ShipmentID),
			Document:         documentPayload,
			Title:            handlers.FmtStringPtr(&doc.Title),
			MoveDocumentType: apimessages.MoveDocumentType(doc.MoveDocumentType),
			Status:           apimessages.MoveDocumentStatus(doc.Status),
			Notes:            handlers.FmtStringPtr(doc.Notes),
		}
		if err != nil {
			return h.RespondAndTraceError(ctx, err, "error fetching payload for document")
		}
		moveDocumentsPayload[i] = &moveDocumentPayload
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
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the TSP user is authorized to update move doc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	_, _, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
	if err != nil {
		if err.Error() == "USER_UNAUTHORIZED" {
			h.Logger().Error("DB Query", zap.Error(err))
			return h.RespondAndTraceError(ctx, err, "error user unauthorized")
		}
		if err.Error() == "FETCH_FORBIDDEN" {
			h.Logger().Error("DB Query", zap.Error(err))
			return h.RespondAndTraceError(ctx, err, "error fetch forbidden")
		}
		return h.RespondAndTraceError(ctx, err, "error fetching shipment", zap.String("shipment_id", shipmentID.String()))
	}

	// Fetch move document from move doc id
	moveDocID, _ := uuid.FromString(params.MoveDocumentID.String())
	moveDoc, err := models.FetchMoveDocument(h.DB(), session, moveDocID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching move document", zap.String("move_id_id", moveDocID.String()))
	}
	if moveDoc.ShipmentID == nil {
		h.Logger().Error("Move document is not associated to a shipment.")
		return movedocop.NewCreateGenericMoveDocumentForbidden()
	}
	if *moveDoc.ShipmentID != shipmentID {
		h.Logger().Error("Move doc shipment ID does not match shipment ID.")
		return movedocop.NewCreateGenericMoveDocumentForbidden()
	}

	// Set new values on move document
	payload := params.UpdateMoveDocument
	moveDoc.ShipmentID = &shipmentID
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes
	moveDoc.MoveDocumentType = newType
	newStatus := models.MoveDocumentStatus(payload.Status)

	// If this is a shipment summary and it has been approved, we process the shipment.
	if newStatus != moveDoc.Status {
		err = moveDoc.AttemptTransition(newStatus)
		if err != nil {
			return h.RespondAndTraceError(ctx, err, "error changing move doc status")
		}
	}

	var saveAction models.MoveDocumentSaveAction

	verrs, err := models.SaveMoveDocument(h.DB(), moveDoc, saveAction)
	if err != nil || verrs.HasAny() {
		return h.RespondAndTraceVErrors(ctx, verrs, err, "error saving move document")
	}

	moveDocPayload, err := payloadForGenericMoveDocumentModel(h.FileStorer(), *moveDoc, shipmentID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching payload for generic move document")
	}
	return movedocop.NewUpdateMoveDocumentOK().WithPayload(moveDocPayload)
}
