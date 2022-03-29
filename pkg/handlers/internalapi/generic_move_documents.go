package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForGenericMoveDocumentModel(storer storage.FileStorer, moveDocument models.MoveDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, moveDocument.Document)
	if err != nil {
		return nil, err
	}

	moveDocumentType := internalmessages.MoveDocumentType(moveDocument.MoveDocumentType)
	status := internalmessages.MoveDocumentStatus(moveDocument.Status)
	genericMoveDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveID:           handlers.FmtUUID(moveDocument.MoveID),
		Document:         documentPayload,
		Title:            &moveDocument.Title,
		MoveDocumentType: &moveDocumentType,
		Status:           &status,
		Notes:            moveDocument.Notes,
	}

	return &genericMoveDocumentPayload, nil
}

// CreateGenericMoveDocumentHandler creates a MoveDocument
type CreateGenericMoveDocumentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h CreateGenericMoveDocumentHandler) Handle(params movedocop.CreateGenericMoveDocumentParams) middleware.Responder {
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

			payload := params.CreateGenericMoveDocumentPayload

			// Fetch uploads to confirm ownership
			uploadIds := payload.UploadIds
			if len(uploadIds) == 0 {
				badRequestErr := apperror.NewUnprocessableEntityError("There are no upload IDs")
				return movedocop.NewCreateGenericMoveDocumentBadRequest(), badRequestErr
			}

			userUploads := models.UserUploads{}
			for _, id := range uploadIds {
				convertedUploadID := uuid.Must(uuid.FromString(id.String()))
				userUpload, fetchUploadErr := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), convertedUploadID)
				if fetchUploadErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), fetchUploadErr), fetchUploadErr
				}
				userUploads = append(userUploads, userUpload)
			}

			var ppmID *uuid.UUID
			if payload.PersonallyProcuredMoveID != nil {
				id := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

				// Enforce that the ppm's move_id matches our move
				ppm, fetchPPMErr := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), id)
				if fetchPPMErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), fetchPPMErr), fetchPPMErr
				}
				if ppm.MoveID != moveID {
					ppmMatchErr := apperror.NewUnprocessableEntityError("PPM Move ID doesn't match original Move ID")
					return movedocop.NewCreateGenericMoveDocumentBadRequest(), ppmMatchErr
				}

				ppmID = &id
			}

			if payload.MoveDocumentType == nil {
				moveDocErr := apperror.NewBadDataError("missing required field: MoveDocumentType")
				return handlers.ResponseForError(appCtx.Logger(), moveDocErr), moveDocErr
			}
			newMoveDocument, verrs, err := move.CreateMoveDocument(appCtx.DB(),
				userUploads,
				ppmID,
				models.MoveDocumentType(*payload.MoveDocumentType),
				*payload.Title,
				payload.Notes,
				*move.SelectedMoveType)

			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			newPayload, err := payloadForGenericMoveDocumentModel(h.FileStorer(), *newMoveDocument)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return movedocop.NewCreateGenericMoveDocumentOK().WithPayload(newPayload), nil
		})
}
