package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForMoveModel(storer storage.FileStorer, order models.Order, move models.Move) (*internalmessages.MovePayload, error) {

	var mtoPayloads internalmessages.MTOShipments
	for _, shipments := range move.MTOShipments {
		shipmentCopy := shipments // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload := payloads.MTOShipment(storer, &shipmentCopy)
		mtoPayloads = append(mtoPayloads, payload)
	}

	var SubmittedAt time.Time
	if move.SubmittedAt != nil {
		SubmittedAt = *move.SubmittedAt
	}

	eTag := etag.GenerateEtag(move.UpdatedAt)

	movePayload := &internalmessages.MovePayload{
		CreatedAt:       handlers.FmtDateTime(move.CreatedAt),
		SubmittedAt:     handlers.FmtDateTime(SubmittedAt),
		Locator:         models.StringPointer(move.Locator),
		ID:              handlers.FmtUUID(move.ID),
		UpdatedAt:       handlers.FmtDateTime(move.UpdatedAt),
		MtoShipments:    mtoPayloads,
		OrdersID:        handlers.FmtUUID(order.ID),
		ServiceMemberID: *handlers.FmtUUID(order.ServiceMemberID),
		Status:          internalmessages.MoveStatus(move.Status),
		ETag:            &eTag,
	}

	if move.CloseoutOffice != nil {
		movePayload.CloseoutOffice = payloads.TransportationOffice(*move.CloseoutOffice)
	}
	if move.PrimeCounselingCompletedAt != nil {
		movePayload.PrimeCounselingCompletedAt = *handlers.FmtDateTime(*move.PrimeCounselingCompletedAt)
	}
	return movePayload, nil
}

func payloadForInternalMove(storer storage.FileStorer, list models.Moves) []*internalmessages.InternalMove {
	var convertedCurrentMovesList []*internalmessages.InternalMove = []*internalmessages.InternalMove{}

	if len(list) == 0 {
		return convertedCurrentMovesList
	}

	// Convert moveList to internalmessages.InternalMove
	for _, move := range list {

		eTag := etag.GenerateEtag(move.UpdatedAt)
		shipments := move.MTOShipments
		var filteredShipments models.MTOShipments
		for _, shipment := range shipments {
			// Check if the DeletedAt field is nil
			if shipment.DeletedAt == nil {
				// If not nil, add the shipment to the filtered array
				filteredShipments = append(filteredShipments, shipment)
			}
		}
		var payloadShipments *internalmessages.MTOShipments = payloads.MTOShipments(storer, &filteredShipments)
		orders, _ := payloadForOrdersModel(storer, move.Orders)
		moveID := *handlers.FmtUUID(move.ID)

		var closeOutOffice internalmessages.TransportationOffice
		if move.CloseoutOffice != nil {
			closeOutOffice = *payloads.TransportationOffice(*move.CloseoutOffice)
		}

		currentMove := &internalmessages.InternalMove{
			CreatedAt:      *handlers.FmtDateTime(move.CreatedAt),
			ETag:           eTag,
			ID:             moveID,
			Status:         string(move.Status),
			MtoShipments:   *payloadShipments,
			MoveCode:       move.Locator,
			Orders:         orders,
			CloseoutOffice: &closeOutOffice,
			SubmittedAt:    handlers.FmtDateTimePtr(move.SubmittedAt),
		}

		if move.PrimeCounselingCompletedAt != nil {
			currentMove.PrimeCounselingCompletedAt = *handlers.FmtDateTime(*move.PrimeCounselingCompletedAt)
		}

		convertedCurrentMovesList = append(convertedCurrentMovesList, currentMove)
	}
	return convertedCurrentMovesList
}

func payloadForMovesList(storer storage.FileStorer, previousMovesList models.Moves, currentMoveList models.Moves, movesList models.Moves) *internalmessages.MovesList {

	if len(movesList) == 0 {
		return &internalmessages.MovesList{
			CurrentMove:   []*internalmessages.InternalMove{},
			PreviousMoves: []*internalmessages.InternalMove{},
		}
	}

	return &internalmessages.MovesList{
		CurrentMove:   payloadForInternalMove(storer, currentMoveList),
		PreviousMoves: payloadForInternalMove(storer, previousMovesList),
	}
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID, _ := uuid.FromString(params.MoveID.String())

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			// Fetch orders for authorized user
			orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return moveop.NewShowMoveOK().WithPayload(movePayload), nil
		})
}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler struct {
	handlers.HandlerConfig
	services.MoveCloseoutOfficeUpdater
}

// Handle ... patches a Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("PatchMoveHandler error", zap.Error(err))
				switch errors.Cause(err) {
				case models.ErrFetchForbidden:
					return moveop.NewPatchMoveForbidden(), err
				case models.ErrFetchNotFound:
					return moveop.NewPatchMoveNotFound(), err
				default:
					switch err.(type) {
					case apperror.NotFoundError:
						return moveop.NewPatchMoveNotFound(), err
					case apperror.PreconditionFailedError:
						return moveop.NewPatchMovePreconditionFailed(), err
					default:
						return moveop.NewPatchMoveInternalServerError(), err
					}
				}
			}

			if !appCtx.Session().IsMilApp() || !appCtx.Session().IsServiceMember() {
				return moveop.NewPatchMoveUnauthorized(), nil
			}

			moveID := uuid.FromStringOrNil(params.MoveID.String())

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handleError(err)
			}

			// Fetch orders for authorized user
			orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
			if err != nil {
				return handleError(err)
			}

			closeoutOfficeID := uuid.FromStringOrNil(params.PatchMovePayload.CloseoutOfficeID.String())
			move, err = h.MoveCloseoutOfficeUpdater.UpdateCloseoutOffice(appCtx, move.Locator, closeoutOfficeID, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), orders, *move)
			if err != nil {
				return handleError(err)
			}

			return moveop.NewPatchMoveOK().WithPayload(movePayload), nil
		})
}

// SubmitMoveHandler approves a move via POST /moves/{moveId}/submit
type SubmitMoveHandler struct {
	handlers.HandlerConfig
	services.MoveRouter
}

// Handle ... submit a move to TOO for approval
func (h SubmitMoveHandler) Handle(params moveop.SubmitMoveForApprovalParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveID, _ := uuid.FromString(params.MoveID.String())

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			newSignedCertification := payloads.SignedCertificationFromSubmit(params.SubmitMoveForApprovalPayload, appCtx.Session().UserID, params.MoveID)
			err = h.MoveRouter.Submit(appCtx, move, newSignedCertification)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			/* Don't send Move Creation email if orders type is BLUEBARK */
			if move.Orders.OrdersType != "BLUEBARK" {
				err = h.NotificationSender().SendNotification(appCtx,
					notifications.NewMoveSubmitted(moveID),
				)
				if err != nil {
					logger.Error("problem sending email to user", zap.Error(err))
					return handlers.ResponseForError(logger, err), err
				}
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			return moveop.NewSubmitMoveForApprovalOK().WithPayload(movePayload), nil
		})
}

// SubmitAmendedOrdersHandler approves a move via POST /moves/{moveId}/submit
type SubmitAmendedOrdersHandler struct {
	handlers.HandlerConfig
	services.MoveRouter
}

// Handle ... submit a move to TOO for approval
func (h SubmitAmendedOrdersHandler) Handle(params moveop.SubmitAmendedOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveID, _ := uuid.FromString(params.MoveID.String())

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			logger := appCtx.Logger().With(zap.String("moveLocator", move.Locator))

			err = h.MoveRouter.RouteAfterAmendingOrders(appCtx, move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}

			responseVErrors := validate.NewErrors()
			var responseError error

			if verrs, saveErr := appCtx.DB().ValidateAndSave(move); verrs.HasAny() || saveErr != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(saveErr, "Error Saving Move")
			}

			if responseVErrors.HasAny() {
				return handlers.ResponseForVErrors(logger, responseVErrors, responseError), responseError
			}

			movePayload, err := payloadForMoveModel(h.FileStorer(), move.Orders, *move)
			if err != nil {
				return handlers.ResponseForError(logger, err), err
			}
			return moveop.NewSubmitAmendedOrdersOK().WithPayload(movePayload), nil
		})
}

type GetAllMovesHandler struct {
	handlers.HandlerConfig
}

// GetAllMovesHandler returns the current and all previous moves of a service member
func (h GetAllMovesHandler) Handle(params moveop.GetAllMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// Grab service member ID from params
			serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())

			// Grab the serviceMember by serviceMemberId
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			var movesList models.Moves
			var latestMove models.Move
			var previousMovesList models.Moves
			var currentMovesList models.Moves

			// Get All Moves for the ServiceMember
			for _, order := range serviceMember.Orders {
				moves, fetchErr := models.FetchMovesByOrderID(appCtx.DB(), order.ID)
				if fetchErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}

				movesList = append(movesList, moves...)
			}

			// Find the move with the latest CreatedAt Date. That one will be the current move
			var nilTime time.Time
			for _, move := range movesList {
				if latestMove.CreatedAt == nilTime {
					latestMove = move
					break
				}
				if move.CreatedAt.After(latestMove.CreatedAt) && move.CreatedAt != latestMove.CreatedAt {
					latestMove = move
				}
			}

			// Place latest move in currentMovesList array
			currentMovesList = append(currentMovesList, latestMove)

			// Populate previousMovesList
			for _, move := range movesList {
				if move.ID != latestMove.ID {
					previousMovesList = append(previousMovesList, move)
				}
			}

			return moveop.NewGetAllMovesOK().WithPayload(payloadForMovesList(h.FileStorer(), previousMovesList, currentMovesList, movesList)), nil
		})
}
