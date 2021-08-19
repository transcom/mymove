package adminapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/services/audit"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	moveop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/move"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// IndexMovesHandler returns a list of moves/MTOs via GET /moves
type IndexMovesHandler struct {
	handlers.HandlerContext
	services.MoveListFetcher
	services.NewQueryFilter
	services.NewPagination
}

func payloadForMoveModel(move models.Move) *adminmessages.Move {
	showMove := true
	if move.Show != nil {
		showMove = *move.Show
	}

	return &adminmessages.Move{
		ID:       handlers.FmtUUID(move.ID),
		OrdersID: handlers.FmtUUID(move.OrdersID),
		Locator:  &move.Locator,
		ServiceMember: &adminmessages.ServiceMember{
			ID:         *handlers.FmtUUID(move.Orders.ServiceMember.ID),
			UserID:     *handlers.FmtUUID(move.Orders.ServiceMember.UserID),
			FirstName:  move.Orders.ServiceMember.FirstName,
			MiddleName: move.Orders.ServiceMember.MiddleName,
			LastName:   move.Orders.ServiceMember.LastName,
		},
		Status:    adminmessages.MoveStatus(move.Status),
		Show:      &showMove,
		CreatedAt: handlers.FmtDateTime(move.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(move.UpdatedAt),
	}
}

var locatorFilterConverters = map[string]func(string) []services.QueryFilter{
	"locator": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("locator", "=", content)}
	},
}

// Handle retrieves a list of moves/MTOs
func (h IndexMovesHandler) Handle(params moveop.IndexMovesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	pagination := h.NewPagination(params.Page, params.PerPage)
	queryFilters := generateQueryFilters(logger, params.Filter, locatorFilterConverters)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("Orders.ServiceMember"),
	}
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	associations := query.NewQueryAssociationsPreload(queryAssociations)
	moves, err := h.MoveListFetcher.FetchMoveList(queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	movesCount := len(moves)

	totalMoveCount, err := h.MoveListFetcher.FetchMoveCount(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := make(adminmessages.Moves, movesCount)
	for i, s := range moves {
		payload[i] = payloadForMoveModel(s)
	}

	return moveop.NewIndexMovesOK().WithContentRange(fmt.Sprintf("moves %d-%d/%d", pagination.Offset(), pagination.Offset()+movesCount, totalMoveCount)).WithPayload(payload)
}

// UpdateMoveHandler updates a given move
type UpdateMoveHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderUpdater
}

// Handle updates a given move
func (h UpdateMoveHandler) Handle(params moveop.UpdateMoveParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("adminapi.UpdateMoveHandler error - Bad MoveID passed in: %s", params.MoveID), zap.Error(err))
		return moveop.NewUpdateMoveBadRequest()
	}

	updatedMove, err := h.MoveTaskOrderUpdater.ShowHide(moveID, params.Move.Show)
	if err != nil {
		switch e := err.(type) {
		case services.NotFoundError:
			return moveop.NewUpdateMoveNotFound()
		case services.InvalidInputError:
			return moveop.NewUpdateMoveUnprocessableEntity() // todo payload
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("adminapi.UpdateMoveHandler query error", zap.Error(e.Unwrap()))
			}
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError()
		default:
			return moveop.NewUpdateMoveInternalServerError()
		}
	}

	_, err = audit.Capture(updatedMove, params.Move, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Error capturing audit record", zap.Error(err))
	}

	if updatedMove == nil {
		logger.Debug(fmt.Sprintf("adminapi.UpdateMoveHandler - No Move returned from ShowHide update, but no error returned either. ID: %s", moveID))
		return moveop.NewUpdateMoveInternalServerError()
	}

	movePayload := payloadForMoveModel(*updatedMove)

	return moveop.NewUpdateMoveOK().WithPayload(movePayload)
}

// GetMoveHandler retrieves the info for a given move
type GetMoveHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a given move by move id
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	move := models.Move{}
	// Returns move by id and associated order and the service memeber associated with the order
	err := h.DB().Eager("Orders", "Orders.ServiceMember").Find(&move, params.MoveID.String())

	if err != nil {
		switch e := err.(type) {
		case services.NotFoundError:
			return moveop.NewGetMoveNotFound()
		case services.InvalidInputError:
			return moveop.NewGetMoveBadRequest()
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("adminapi.GetMoveHandler query error", zap.Error(e.Unwrap()))
			}
			return moveop.NewGetMoveInternalServerError()
		default:
			return moveop.NewGetMoveInternalServerError()
		}
	}

	payload := payloadForMoveModel(move)
	return moveop.NewGetMoveOK().WithPayload(payload)
}
