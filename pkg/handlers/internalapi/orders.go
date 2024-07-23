package internalapi

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

func payloadForUploadModelFromAmendedOrdersUpload(storer storage.FileStorer, upload models.Upload, url string) (*internalmessages.Upload, error) {
	uploadPayload := &internalmessages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
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
	orderPayload, err := payloads.PayloadForDocumentModel(storer, order.UploadedOrders)
	if err != nil {
		return nil, err
	}

	var amendedOrderPayload *internalmessages.Document
	if order.UploadedAmendedOrders != nil {
		amendedOrderPayload, err = payloads.PayloadForDocumentModel(storer, *order.UploadedAmendedOrders)
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
	var entitlement internalmessages.Entitlement
	if order.Entitlement != nil {
		dBAuthorizedWeight = models.Int64Pointer(int64(*order.Entitlement.AuthorizedWeight()))
		entitlement.ProGear = models.Int64Pointer(int64(order.Entitlement.ProGearWeight))
		entitlement.ProGearSpouse = models.Int64Pointer(int64(order.Entitlement.ProGearWeightSpouse))
	}
	var originDutyLocation models.DutyLocation
	originDutyLocation = models.DutyLocation{}
	if order.OriginDutyLocation != nil {
		originDutyLocation = *order.OriginDutyLocation
	}

	var grade internalmessages.OrderPayGrade
	if order.Grade != nil {
		grade = internalmessages.OrderPayGrade(*order.Grade)
	}

	ordersType := order.OrdersType
	payload := &internalmessages.Orders{
		ID:                         handlers.FmtUUID(order.ID),
		CreatedAt:                  handlers.FmtDateTime(order.CreatedAt),
		UpdatedAt:                  handlers.FmtDateTime(order.UpdatedAt),
		ServiceMemberID:            handlers.FmtUUID(order.ServiceMemberID),
		IssueDate:                  handlers.FmtDate(order.IssueDate),
		ReportByDate:               handlers.FmtDate(order.ReportByDate),
		OrdersType:                 &ordersType,
		OrdersTypeDetail:           order.OrdersTypeDetail,
		OriginDutyLocation:         payloadForDutyLocationModel(originDutyLocation),
		OriginDutyLocationGbloc:    handlers.FmtStringPtr(order.OriginDutyLocationGBLOC),
		Grade:                      &grade,
		NewDutyLocation:            payloadForDutyLocationModel(order.NewDutyLocation),
		HasDependents:              handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:           handlers.FmtBool(order.SpouseHasProGear),
		UploadedOrders:             orderPayload,
		UploadedAmendedOrders:      amendedOrderPayload,
		OrdersNumber:               order.OrdersNumber,
		Moves:                      moves,
		Tac:                        order.TAC,
		Sac:                        order.SAC,
		DepartmentIndicator:        (*internalmessages.DeptIndicator)(order.DepartmentIndicator),
		Status:                     internalmessages.OrdersStatus(order.Status),
		AuthorizedWeight:           dBAuthorizedWeight,
		Entitlement:                &entitlement,
		ProvidesServicesCounseling: originDutyLocation.ProvidesServicesCounseling,
	}

	return payload, nil
}

// CreateOrdersHandler creates new orders via POST /orders
type CreateOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle ... creates new Orders from a request payload
func (h CreateOrdersHandler) Handle(params ordersop.CreateOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.CreateOrders

			serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			originDutyLocationID, err := uuid.FromString(payload.OriginDutyLocationID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			originDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), originDutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newDutyLocationID, err := uuid.FromString(payload.NewDutyLocationID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			newDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), newDutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newDutyLocationGBLOC, err := models.FetchGBLOCForPostalCode(appCtx.DB(), newDutyLocation.Address.PostalCode)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			originDutyLocationGBLOC, err := models.FetchGBLOCForPostalCode(appCtx.DB(), originDutyLocation.Address.PostalCode)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return nil, apperror.NewNotFoundError(originDutyLocation.ID, "while looking for Duty Location PostalCodeToGBLOC")
				default:
					return nil, apperror.NewQueryError("PostalCodeToGBLOC", err, "")
				}
			}

			grade := payload.Grade
			weightAllotment := models.GetWeightAllotment(*grade)

			weight := weightAllotment.TotalWeightSelf
			if *payload.HasDependents {
				weight = weightAllotment.TotalWeightSelfPlusDependents
			}

			// Assign default SIT allowance based on customer type.
			// We only have service members right now, but once we introduce more, this logic will have to change.
			sitDaysAllowance := models.DefaultServiceMemberSITDaysAllowance

			entitlement := models.Entitlement{
				DependentsAuthorized: payload.HasDependents,
				DBAuthorizedWeight:   models.IntPointer(weight),
				StorageInTransit:     models.IntPointer(sitDaysAllowance),
				ProGearWeight:        weightAllotment.ProGearWeight,
				ProGearWeightSpouse:  weightAllotment.ProGearWeightSpouse,
			}

			/*
				IF you get that to work you'll still have to add conditionals for all the places the entitlement is used because it
				isn't inheritly clear if it's using the spouse weight or not. So you'll be creating new variables and conditionals
				in move_dats.go, move_weights, and move_submitted, etc
			*/

			if saveEntitlementErr := appCtx.DB().Save(&entitlement); saveEntitlementErr != nil {
				return handlers.ResponseForError(appCtx.Logger(), saveEntitlementErr), saveEntitlementErr
			}

			var deptIndicator *string
			if payload.DepartmentIndicator != nil {
				converted := string(*payload.DepartmentIndicator)
				deptIndicator = &converted
			}

			if payload.OrdersType == nil {
				errMsg := "missing required field: OrdersType"
				return handlers.ResponseForError(appCtx.Logger(), errors.New(errMsg)), apperror.NewBadDataError("missing required field: OrdersType")
			}

			contractor, err := models.FetchGHCPrimeContractor(appCtx.DB())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			packingAndShippingInstructions := models.InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + models.InstructionsAfterContractNumber
			newOrder, verrs, err := serviceMember.CreateOrder(
				appCtx,
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
				&originDutyLocationGBLOC.GBLOC,
				packingAndShippingInstructions,
				&newDutyLocationGBLOC.GBLOC,
			)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			moveOptions := models.MoveOptions{
				Show: models.BoolPointer(true),
			}
			newMove, verrs, err := newOrder.CreateNewMove(appCtx.DB(), moveOptions)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}
			newOrder.Moves = append(newOrder.Moves, *newMove)

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), newOrder)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewCreateOrdersCreated().WithPayload(orderPayload), nil
		})
}

// ShowOrdersHandler returns orders for a user and order ID
type ShowOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h ShowOrdersHandler) Handle(params ordersop.ShowOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			orderID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), orderID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewShowOrdersOK().WithPayload(orderPayload), nil
		})
}

// UpdateOrdersHandler updates an order via PUT /orders/{orderId}
type UpdateOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle ... updates an order from a request payload
func (h UpdateOrdersHandler) Handle(params ordersop.UpdateOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			orderID, err := uuid.FromString(params.OrdersID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), orderID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := params.UpdateOrders
			dutyLocationID, err := uuid.FromString(payload.NewDutyLocationID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			dutyLocation, err := models.FetchDutyLocation(appCtx.DB(), dutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newDutyLocationGBLOC, err := models.FetchGBLOCForPostalCode(appCtx.DB(), dutyLocation.Address.PostalCode)
			if err != nil {
				err = apperror.NewBadDataError("New duty location GBLOC cannot be verified")
				appCtx.Logger().Error(err.Error())
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			if payload.OriginDutyLocationID != "" {
				originDutyLocationID, errorOrigin := uuid.FromString(payload.OriginDutyLocationID.String())
				if errorOrigin != nil {
					return handlers.ResponseForError(appCtx.Logger(), errorOrigin), errorOrigin
				}
				originDutyLocation, errorOrigin := models.FetchDutyLocation(appCtx.DB(), originDutyLocationID)
				if errorOrigin != nil {
					return handlers.ResponseForError(appCtx.Logger(), errorOrigin), errorOrigin
				}
				order.OriginDutyLocation = &originDutyLocation
				order.OriginDutyLocationID = &originDutyLocationID
			}

			if payload.OrdersType == nil {
				errMsg := "missing required field: OrdersType"
				return handlers.ResponseForError(appCtx.Logger(), errors.New(errMsg)), apperror.NewBadDataError("missing required field: OrdersType")
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
			order.DestinationGBLOC = &newDutyLocationGBLOC.GBLOC
			order.TAC = payload.Tac
			order.SAC = payload.Sac

			// Check if the grade is receiving an update
			if order.Grade != payload.Grade {
				weightAllotment := models.GetWeightAllotment(*payload.Grade)
				weight := weightAllotment.TotalWeightSelf
				if *payload.HasDependents {
					weight = weightAllotment.TotalWeightSelfPlusDependents
				}

				// Assign default SIT allowance based on customer type.
				// We only have service members right now, but once we introduce more, this logic will have to change.
				sitDaysAllowance := models.DefaultServiceMemberSITDaysAllowance

				entitlement := models.Entitlement{
					DependentsAuthorized: payload.HasDependents,
					DBAuthorizedWeight:   models.IntPointer(weight),
					StorageInTransit:     models.IntPointer(sitDaysAllowance),
					ProGearWeight:        weightAllotment.ProGearWeight,
					ProGearWeightSpouse:  weightAllotment.ProGearWeightSpouse,
				}

				/*
					IF you get that to work you'll still have to add conditionals for all the places the entitlement is used because it
					isn't inheritly clear if it's using the spouse weight or not. So you'll be creating new variables and conditionals
					in move_dats.go, move_weights, and move_submitted, etc
				*/

				if saveEntitlementErr := appCtx.DB().Save(&entitlement); saveEntitlementErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), saveEntitlementErr), saveEntitlementErr
				}

				order.EntitlementID = &entitlement.ID
				order.Entitlement = &entitlement
			}
			order.Grade = payload.Grade

			if payload.DepartmentIndicator != nil {
				order.DepartmentIndicator = handlers.FmtString(string(*payload.DepartmentIndicator))
			}

			verrs, err := models.SaveOrder(appCtx.DB(), &order)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewUpdateOrdersOK().WithPayload(orderPayload), nil
		})
}

// UploadAmendedOrdersHandler uploads amended orders to an order via PATCH /orders/{orderId}/upload_amended_orders
type UploadAmendedOrdersHandler struct {
	handlers.HandlerConfig
	services.OrderUpdater
}

// Handle updates an order to attach amended orders from a request payload
func (h UploadAmendedOrdersHandler) Handle(params ordersop.UploadAmendedOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			file, ok := params.File.(*runtime.File)
			if !ok {
				errMsg := "This should always be a runtime.File, something has changed in go-swagger."

				appCtx.Logger().Error(errMsg)

				return ordersop.NewUploadAmendedOrdersInternalServerError(), nil
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
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			upload, url, verrs, err := h.OrderUpdater.UploadAmendedOrdersAsCustomer(appCtx, appCtx.Session().UserID, orderID, file.Data, file.Header.Filename, h.FileStorer())

			if verrs.HasAny() || err != nil {
				switch err.(type) {
				case uploader.ErrTooLarge:
					return ordersop.NewUploadAmendedOrdersRequestEntityTooLarge(), err
				case uploader.ErrFile:
					return ordersop.NewUploadAmendedOrdersInternalServerError(), err
				case uploader.ErrFailedToInitUploader:
					return ordersop.NewUploadAmendedOrdersInternalServerError(), err
				case apperror.NotFoundError:
					return ordersop.NewUploadAmendedOrdersNotFound(), err
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
				}
			}

			uploadPayload, err := payloadForUploadModelFromAmendedOrdersUpload(h.FileStorer(), upload, url)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return ordersop.NewUploadAmendedOrdersCreated().WithPayload(uploadPayload), nil
		})
}
