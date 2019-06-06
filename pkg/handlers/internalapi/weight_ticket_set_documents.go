package internalapi

import (
	"reflect"
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/storage"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForWeightTicketSetMoveDocumentModel(storer storage.FileStorer, weightTicketSet models.WeightTicketSetDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, weightTicketSet.MoveDocument.Document)
	if err != nil {
		return nil, err
	}

	//
	genericMoveDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(weightTicketSet.ID),
		MoveID:           handlers.FmtUUID(weightTicketSet.MoveDocument.MoveID),
		Document:         documentPayload,
		Title:            &weightTicketSet.MoveDocument.Title,
		MoveDocumentType: internalmessages.MoveDocumentType(weightTicketSet.MoveDocument.MoveDocumentType),
		Status:           internalmessages.MoveDocumentStatus(weightTicketSet.MoveDocument.Status),
		Notes:            weightTicketSet.MoveDocument.Notes,
	}

	return &genericMoveDocumentPayload, nil
}

// CreateMovingExpenseDocumentHandler creates a MovingExpenseDocument
type CreateWeightTicketSetDocumentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CreateWeightTicketSetDocumentHandler) Handle(params movedocop.CreateWeightTicketDocumentParams) middleware.Responder {
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

	payload := params.CreateWeightTicketDocument

	// Fetch uploads to confirm ownership
	uploadIds := payload.UploadIds
	if len(uploadIds) == 0 {
		return movedocop.NewCreateWeightTicketDocumentBadRequest()
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

	ppmID := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

	// Enforce that the ppm's move_id matches our move
	ppm, fetchPPMErr := models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
	if fetchPPMErr != nil {
		return handlers.ResponseForError(h.Logger(), fetchPPMErr)
	}
	if ppm.MoveID != moveID {
		return movedocop.NewCreateMovingExpenseDocumentBadRequest()
	}

	newWeightTicketSetDocument, verrs, err := move.CreateWeightTicketSetDocument(
		h.DB(),
		uploads,
		&ppmID,
		*payload.VehicleNickname,
		*payload.VehicleOptions,
		payload.FullWeight,
		payload.EmptyWeight,
		(*time.Time)(payload.WeightTicketDate),
		*move.SelectedMoveType,
	)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	newPayload, err := payloadForWeightTicketSetMoveDocumentModel(h.FileStorer(), *newWeightTicketSetDocument)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return movedocop.NewCreateMovingExpenseDocumentOK().WithPayload(newPayload)
}
