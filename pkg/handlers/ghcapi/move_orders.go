package ghcapi

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetMoveOrdersHandler fetches the information of a specific move order
type GetMoveOrdersHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// Handle getting the information of a specific move order
func (h GetMoveOrdersHandler) Handle(params moveorderop.GetMoveOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrderID, _ := uuid.FromString(params.MoveOrderID.String())
	moveOrder, err := h.FetchMoveOrder(moveOrderID)
	if err != nil {
		logger.Error("fetching move order", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewGetMoveOrderNotFound()
		default:
			return moveorderop.NewGetMoveOrderInternalServerError()
		}
	}
	moveOrderPayload := payloads.MoveOrder(moveOrder)
	return moveorderop.NewGetMoveOrderOK().WithPayload(moveOrderPayload)
}

// ListMoveOrdersHandler fetches all the move orders
type ListMoveOrdersHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// Handle getting the all move orders
func (h ListMoveOrdersHandler) Handle(params moveorderop.ListMoveOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrders, err := h.ListMoveOrders()
	if err != nil {
		logger.Error("fetching all move orders", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewListMoveOrdersNotFound()
		default:
			return moveorderop.NewListMoveOrdersInternalServerError()
		}
	}
	moveOrdersPayload := make(ghcmessages.MoveOrders, len(moveOrders))
	for i, moveOrder := range moveOrders {
		moveOrdersPayload[i] = payloads.MoveOrder(&moveOrder)
	}
	return moveorderop.NewListMoveOrdersOK().WithPayload(moveOrdersPayload)
}

// ListMoveTaskOrdersHandler fetches all the move orders
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
	services.MoveTaskOrderFetcher
}

// Handle getting the all move orders
func (h ListMoveTaskOrdersHandler) Handle(params moveorderop.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveOrderID, _ := uuid.FromString(params.MoveOrderID.String())
	moveTaskOrders, err := h.ListMoveTaskOrders(moveOrderID)
	if err != nil {
		logger.Error("fetching all move orders", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return moveorderop.NewListMoveTaskOrdersNotFound()
		default:
			return moveorderop.NewListMoveTaskOrdersInternalServerError()
		}
	}
	moveTaskOrdersPayload := make(ghcmessages.MoveTaskOrders, len(moveTaskOrders))
	for i, moveTaskOrder := range moveTaskOrders {
		moveTaskOrdersPayload[i] = payloads.MoveTaskOrder(&moveTaskOrder)
	}
	return moveorderop.NewListMoveTaskOrdersOK().WithPayload(moveTaskOrdersPayload)
}

// UpdateMoveOrderHandler updates an order via PATCH /move-orders/{moveOrderId}
type UpdateMoveOrderHandler struct {
	handlers.HandlerContext
	services.MoveOrderFetcher
}

// Handle ... updates an order from a request payload
func (h UpdateMoveOrderHandler) Handle(params moveorderop.UpdateMoveOrderParams) middleware.Responder {

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	orderID, err := uuid.FromString(params.MoveOrderID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	order, err := models.FetchOrderForUser(h.DB(), session, orderID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := params.Body
	originStationID, err := uuid.FromString(payload.OriginDutyStationID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	originDutyStation, err := models.FetchDutyStation(h.DB(), originStationID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	newStationID, err := uuid.FromString(payload.NewDutyStationID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	newDutyStation, err := models.FetchDutyStation(h.DB(), newStationID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	ordersTypeDetail := internalmessages.OrdersTypeDetail(payload.OrdersTypeDetail)

	order.OrdersNumber = payload.OrdersNumber
	order.IssueDate = time.Time(*payload.IssueDate)
	order.ReportByDate = time.Time(*payload.ReportByDate)
	order.OrdersType = internalmessages.OrdersType(payload.OrdersType)
	order.OrdersTypeDetail = &ordersTypeDetail
	order.HasDependents = *payload.HasDependents
	order.SpouseHasProGear = *payload.SpouseHasProGear
	order.OriginDutyStationID = &originDutyStation.ID
	order.OriginDutyStation = &originDutyStation
	order.NewDutyStationID = newDutyStation.ID
	order.NewDutyStation = newDutyStation
	order.TAC = payload.Tac
	order.SAC = payload.Sac
	order.DepartmentIndicator = handlers.FmtString(string(payload.DepartmentIndicator))

	verrs, err := models.SaveOrder(h.DB(), &order)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return moveorderop.NewUpdateMoveOrderOK().WithPayload(orderPayload)
}

func payloadForOrdersModel(storer storage.FileStorer, order models.Order) (*ghcmessages.MoveOrder, error) {
	payload := &ghcmessages.MoveOrder{
		ID:                     *handlers.FmtUUID(order.ID),
		DateIssued:             *handlers.FmtDate(order.IssueDate),
		ReportByDate:           *handlers.FmtDate(order.ReportByDate),
		OrderType:              ghcmessages.OrdersType(order.OrdersType),
		OrderTypeDetail:        ghcmessages.OrdersTypeDetail(*order.OrdersTypeDetail),
		OriginDutyStation:      payloadForDutyStationModel(order.OriginDutyStation),
		DestinationDutyStation: payloadForDutyStationModel(&order.NewDutyStation),
		HasDependents:          *handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:       *handlers.FmtBool(order.SpouseHasProGear),
		OrderNumber:            order.OrdersNumber,
		Tac:                    order.TAC,
		Sac:                    order.SAC,
		DepartmentIndicator:    ghcmessages.DeptIndicator(*order.DepartmentIndicator),
	}

	return payload, nil
}

func payloadForDutyStationModel(station *models.DutyStation) *ghcmessages.DutyStation {
	// If the station ID has no UUID then it isn't real data
	// Unlike other payloads the
	if station == nil || station.ID == uuid.Nil {
		return nil
	}
	payload := ghcmessages.DutyStation{
		ID:        *handlers.FmtUUID(station.ID),
		Name:      station.Name,
		AddressID: *handlers.FmtUUID(station.AddressID),
		Address:   payloadForAddressModel(&station.Address),
	}

	return &payload
}

func payloadForAddressModel(a *models.Address) *ghcmessages.Address {
	if a == nil {
		return nil
	}
	return &ghcmessages.Address{
		ID:             strfmt.UUID(a.ID.String()),
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}