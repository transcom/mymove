package adminapi

import (
	"encoding/json"
	"fmt"

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
		ID:              handlers.FmtUUID(move.ID),
		OrdersID:        handlers.FmtUUID(move.OrdersID),
		ServiceMemberID: *handlers.FmtUUID(move.Orders.ServiceMemberID),
		Locator:         &move.Locator,
		Status:          adminmessages.MoveStatus(move.Status),
		Show:            &showMove,
		CreatedAt:       handlers.FmtDateTime(move.CreatedAt),
		UpdatedAt:       handlers.FmtDateTime(move.UpdatedAt),
	}
}

// Handle retrieves a list of moves/MTOs
func (h IndexMovesHandler) Handle(params moveop.IndexMovesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	pagination := h.NewPagination(params.Page, params.PerPage)
	queryFilters := h.generateQueryFilters(params.Filter, logger)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("Orders.ServiceMember"),
	}
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	associations := query.NewQueryAssociations(queryAssociations)
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

// generateQueryFilters is helper to convert filter params from a json string
// of the form `{"move_type": "PPM" "code": "XYZBCS"}` to an array of services.QueryFilter
func (h IndexMovesHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		Locator string `json:"locator"`
	}
	f := Filter{}
	var queryFilters []services.QueryFilter
	if filters == nil {
		return queryFilters
	}
	b := []byte(*filters)
	err := json.Unmarshal(b, &f)
	if err != nil {
		fs := fmt.Sprintf("%v", filters)
		logger.Warn("unable to decode param", zap.Error(err),
			zap.String("filters", fs))
	}
	if f.Locator != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("locator", "=", f.Locator))
	}

	return queryFilters
}

// UpdateMoveHandler updates a given move
type UpdateMoveHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderUpdater
}

// Handle updates a given move
func (h UpdateMoveHandler) Handle(params moveop.UpdateMoveParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("adminapi.UpdateMoveHandler error - Bad MoveID passed in: %s", params.MoveID), zap.Error(err))
		return moveop.NewUpdateMoveBadRequest()
	}

	updatedMove, err := h.MoveTaskOrderUpdater.ShowHide(moveID, &params.Move.Show)
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
	services.MoveFetcher
	services.NewQueryFilter
}

// Handle retrieves a given move
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	return moveop.NewGetMoveNotImplemented()
}
