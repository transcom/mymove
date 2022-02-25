package internalapi

import (
	"errors"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/uploader"

	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForUploadModelFromAmendedOrdersUpload(storer storage.FileStorer, upload models.Upload, url string) (*internalmessages.UploadPayload, error) {
	uploadPayload := &internalmessages.UploadPayload{
		ID:          handlers.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         handlers.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload, nil
}

func payloadForOrdersModel(storer storage.FileStorer, order models.Order) (*internalmessages.Orders, error) {
	orderPayload, err := payloadForDocumentModel(storer, order.UploadedOrders)
	if err != nil {
		return nil, err
	}

	var amendedOrderPayload *internalmessages.DocumentPayload
	if order.UploadedAmendedOrders != nil {
		amendedOrderPayload, err = payloadForDocumentModel(storer, *order.UploadedAmendedOrders)
		if err != nil {
			return nil, err
		}
	}

	var moves internalmessages.IndexMovesPayload
	for _, move := range order.Moves {
		payload, err := payloadForMoveModel(storer, order, move)
		if err != nil {
			return nil, err
		}
		moves = append(moves, payload)
	}

	var dBAuthorizedWeight *int64
	dBAuthorizedWeight = nil
	if order.Entitlement != nil {
		dBAuthorizedWeight = swag.Int64(int64(*order.Entitlement.AuthorizedWeight()))
	}
	var originDutyStation models.DutyLocation
	originDutyStation = models.DutyLocation{}
	if order.OriginDutyLocation != nil {
		originDutyStation = *order.OriginDutyLocation
	}

	ordersType := order.OrdersType
	payload := &internalmessages.Orders{
		ID:                    handlers.FmtUUID(order.ID),
		CreatedAt:             handlers.FmtDateTime(order.CreatedAt),
		UpdatedAt:             handlers.FmtDateTime(order.UpdatedAt),
		ServiceMemberID:       handlers.FmtUUID(order.ServiceMemberID),
		IssueDate:             handlers.FmtDate(order.IssueDate),
		ReportByDate:          handlers.FmtDate(order.ReportByDate),
		OrdersType:            &ordersType,
		OrdersTypeDetail:      order.OrdersTypeDetail,
		OriginDutyLocation:    payloadForDutyLocationModel(originDutyStation),
		Grade:                 order.Grade,
		NewDutyLocation:       payloadForDutyLocationModel(order.NewDutyLocation),
		HasDependents:         handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:      handlers.FmtBool(order.SpouseHasProGear),
		UploadedOrders:        orderPayload,
		UploadedAmendedOrders: amendedOrderPayload,
		OrdersNumber:          order.OrdersNumber,
		Moves:                 moves,
		Tac:                   order.TAC,
		Sac:                   order.SAC,
		DepartmentIndicator:   (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
		Status:                internalmessages.OrdersStatus(order.Status),
		AuthorizedWeight:      dBAuthorizedWeight,
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler struct {
	handlers.HandlerContext
}

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	payload := params.CreateOrders

	serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	stationID, err := uuid.FromString(payload.NewDutyLocationID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	newDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), stationID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	originDutyLocation := serviceMember.DutyLocation
	grade := (*string)(serviceMember.Rank)

	weight, entitlementErr := models.GetEntitlement(*serviceMember.Rank, *payload.HasDependents)
	if entitlementErr != nil {
		return handlers.ResponseForError(appCtx.Logger(), entitlementErr)
	}

	// Assign default SIT allowance based on customer type.
	// We only have service members right now, but once we introduce more, this logic will have to change.
	sitDaysAllowance := models.DefaultServiceMemberSITDaysAllowance

	entitlement := models.Entitlement{
		DependentsAuthorized: payload.HasDependents,
		DBAuthorizedWeight:   models.IntPointer(weight),
		StorageInTransit:     models.IntPointer(sitDaysAllowance),
	}

	if saveEntitlementErr := appCtx.DB().Save(&entitlement); saveEntitlementErr != nil {
		return handlers.ResponseForError(appCtx.Logger(), saveEntitlementErr)
	}

	var deptIndicator *string
	if payload.DepartmentIndicator != nil {
		converted := string(*payload.DepartmentIndicator)
		deptIndicator = &converted
	}

	if payload.OrdersType == nil {
		return handlers.ResponseForError(appCtx.Logger(), errors.New("missing required field: OrdersType"))
	}
	newOrder, verrs, err := serviceMember.CreateOrder(
		appCtx.DB(),
		time.Time(*payload.IssueDate),
		time.Time(*payload.ReportByDate),
		*payload.OrdersType,
		*payload.HasDependents,
		*payload.SpouseHasProGear,
		newDutyLocation,
		payload.OrdersNumber,
		payload.Tac,
		payload.Sac,
		deptIndicator,
		&originDutyLocation,
		grade,
		&entitlement,
	)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
	}

	moveOptions := models.MoveOptions{
		SelectedType: nil,
		Show:         swag.Bool(true),
	}
	newMove, verrs, err := newOrder.CreateNewMove(appCtx.DB(), moveOptions)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
	}
	newOrder.Moves = append(newOrder.Moves, *newMove)

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), newOrder)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload)
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler struct {
	handlers.HandlerContext
}

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), orderID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return ordersop.NewShowOrdersOK().WithPayload(orderPayload)
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler struct {
	handlers.HandlerContext
}

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), orderID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	payload := params.UpdateOrders
	stationID, err := uuid.FromString(payload.NewDutyLocationID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	dutyLocation, err := models.FetchDutyLocation(appCtx.DB(), stationID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	if payload.OrdersType == nil {
		return handlers.ResponseForError(appCtx.Logger(), errors.New("missing required field: OrdersType"))
	}

	order.OrdersNumber = payload.OrdersNumber
	order.IssueDate = time.Time(*payload.IssueDate)
	order.ReportByDate = time.Time(*payload.ReportByDate)
	order.OrdersType = *payload.OrdersType
	order.OrdersTypeDetail = payload.OrdersTypeDetail
	order.HasDependents = *payload.HasDependents
	order.SpouseHasProGear = *payload.SpouseHasProGear
	order.NewDutyLocationID = dutyLocation.ID
	order.NewDutyLocation = dutyLocation
	order.TAC = payload.Tac
	order.SAC = payload.Sac

	if payload.DepartmentIndicator != nil {
		order.DepartmentIndicator = handlers.FmtString(string(*payload.DepartmentIndicator))
	}

	verrs, err := models.SaveOrder(appCtx.DB(), &order)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload)
}

// UploadAmendedOrdersHandler uploads amended orders to an order via PATCH /orders/{orderId}/upload_amended_orders
type UploadAmendedOrdersHandler struct {
	handlers.HandlerContext
	services.OrderUpdater
}

// Handle updates an order to attach amended orders from a request payload
func (h UploadAmendedOrdersHandler) Handle(params ordersop.UploadAmendedOrdersParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	file, ok := params.File.(*runtime.File)
	if !ok {
		appCtx.Logger().Error("This should always be a runtime.File, something has changed in go-swagger.")
		return handlers.ResponseForError(appCtx.Logger(), nil)
	}

	appCtx.Logger().Info(
		"File uploader and size",
		zap.String("userID", appCtx.Session().UserID.String()),
		zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
		zap.String("officeUserID", appCtx.Session().OfficeUserID.String()),
		zap.String("AdminUserID", appCtx.Session().AdminUserID.String()),
		zap.Int64("size", file.Header.Size),
	)

	orderID, err := uuid.FromString(params.OrdersID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	upload, url, verrs, err := h.OrderUpdater.UploadAmendedOrdersAsCustomer(appCtx, appCtx.Session().UserID, orderID, file.Data, file.Header.Filename, h.FileStorer())

	if verrs.HasAny() || err != nil {
		switch err.(type) {
		case uploader.ErrTooLarge:
			return ordersop.NewUploadAmendedOrdersRequestEntityTooLarge()
		case uploader.ErrFile:
			return ordersop.NewUploadAmendedOrdersInternalServerError()
		case uploader.ErrFailedToInitUploader:
			return ordersop.NewUploadAmendedOrdersInternalServerError()
		case apperror.NotFoundError:
			return ordersop.NewUploadAmendedOrdersNotFound()
		default:
			return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
		}
	}

	uploadPayload, err := payloadForUploadModelFromAmendedOrdersUpload(h.FileStorer(), upload, url)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return ordersop.NewUploadAmendedOrdersCreated().WithPayload(uploadPayload)
}
