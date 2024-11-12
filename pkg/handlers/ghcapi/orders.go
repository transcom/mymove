package ghcapi

import (
	"database/sql"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	orderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// GetOrdersHandler fetches the information of a specific order
type GetOrdersHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
}

// Handle getting the information of a specific order
func (h GetOrdersHandler) Handle(params orderop.GetOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			orderID, _ := uuid.FromString(params.OrderID.String())
			order, err := h.FetchOrder(appCtx, orderID)
			if err != nil {
				appCtx.Logger().Error("fetching order", zap.Error(err))
				switch err {
				case sql.ErrNoRows:
					return orderop.NewGetOrderNotFound(), err
				default:
					return orderop.NewGetOrderInternalServerError(), err
				}
			}
			orderPayload := payloads.Order(order)
			return orderop.NewGetOrderOK().WithPayload(orderPayload), nil
		})
}

// UpdateOrderHandler updates an order via PATCH /orders/{orderId}
type UpdateOrderHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
	moveUpdater  services.MoveTaskOrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateOrderHandler) Handle(params orderop.UpdateOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateOrderNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewUpdateOrderUnprocessableEntity().WithPayload(payload), err
				case apperror.ConflictError:
					return orderop.NewUpdateOrderConflict().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateOrderForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateOrderInternalServerError(), err
				}
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateOrderAsTOO(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdateOrderEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateOrderOK().WithPayload(orderPayload), nil
		})
}

// CounselingUpdateOrderHandler updates an order via PATCH /counseling/orders/{orderId}
type CounselingUpdateOrderHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order as requested by a services counselor
func (h CounselingUpdateOrderHandler) Handle(
	params orderop.CounselingUpdateOrderParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewCounselingUpdateOrderNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewCounselingUpdateOrderUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewCounselingUpdateOrderPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewCounselingUpdateOrderForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewCounselingUpdateOrderInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				return handleError(apperror.NewForbiddenError("is not a Services Counselor"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateOrderAsCounselor(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerCounselingUpdateOrderEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewCounselingUpdateOrderOK().WithPayload(orderPayload), nil
		})
}

// CounselingUpdateOrderHandler create an order via POST /orders
type CreateOrderHandler struct {
	handlers.HandlerConfig
}

// Handle ... creates an order as requested by a services counselor
func (h CreateOrderHandler) Handle(params orderop.CreateOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			const SAC_LIMIT = 80
			payload := params.CreateOrders

			serviceMemberID, err := uuid.FromString(payload.ServiceMemberID.String())
			if err != nil {
				err = apperror.NewBadDataError("Error processing Service Member ID")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				err = apperror.NewBadDataError("Service member cannot be verified")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}

			if payload.Sac != nil && len(*payload.Sac) > SAC_LIMIT {
				err = apperror.NewBadDataError("SAC cannot be more than 80 characters.")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}

			originDutyLocationID, err := uuid.FromString(payload.OriginDutyLocationID.String())
			if err != nil {
				err = apperror.NewBadDataError("Error processing origin duty location ID")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}
			originDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), originDutyLocationID)
			if err != nil {
				err = apperror.NewBadDataError("Origin duty location cannot be verified")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}

			newDutyLocationID, err := uuid.FromString(payload.NewDutyLocationID.String())
			if err != nil {
				err = apperror.NewBadDataError("Error processing new duty location ID")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}
			newDutyLocation, err := models.FetchDutyLocation(appCtx.DB(), newDutyLocationID)
			if err != nil {
				err = apperror.NewBadDataError("New duty location cannot be verified")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}

			newDutyLocationGBLOC, err := models.FetchGBLOCForPostalCode(appCtx.DB(), newDutyLocation.Address.PostalCode)
			if err != nil {
				err = apperror.NewBadDataError("New duty location GBLOC cannot be verified")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
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

			grade := (internalmessages.OrderPayGrade)(*payload.Grade)
			weightAllotment := models.GetWeightAllotment(grade)
			weight := weightAllotment.TotalWeightSelf
			if *payload.HasDependents {
				weight = weightAllotment.TotalWeightSelfPlusDependents
			}

			sitDaysAllowance := models.DefaultServiceMemberSITDaysAllowance

			var dependentsTwelveAndOver *int
			var dependentsUnderTwelve *int
			if payload.DependentsTwelveAndOver != nil {
				dependentsTwelveAndOver = models.IntPointer(int(*payload.DependentsTwelveAndOver))
			}
			if payload.DependentsUnderTwelve != nil {
				dependentsUnderTwelve = models.IntPointer(int(*payload.DependentsUnderTwelve))
			}
			// Calculate UB allowance for the order entitlement
			unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, originDutyLocation.Address.IsOconus, newDutyLocation.Address.IsOconus, serviceMember.Affiliation, &grade, (*internalmessages.OrdersType)(payload.OrdersType), payload.HasDependents, payload.AccompaniedTour, dependentsUnderTwelve, dependentsTwelveAndOver)
			if err == nil {
				weightAllotment.UnaccompaniedBaggageAllowance = unaccompaniedBaggageAllowance
			}

			entitlement := models.Entitlement{
				DependentsAuthorized:    payload.HasDependents,
				DBAuthorizedWeight:      models.IntPointer(weight),
				StorageInTransit:        models.IntPointer(sitDaysAllowance),
				ProGearWeight:           weightAllotment.ProGearWeight,
				ProGearWeightSpouse:     weightAllotment.ProGearWeightSpouse,
				AccompaniedTour:         payload.AccompaniedTour,
				DependentsUnderTwelve:   dependentsUnderTwelve,
				DependentsTwelveAndOver: dependentsTwelveAndOver,
				UBAllowance:             &weightAllotment.UnaccompaniedBaggageAllowance,
			}

			if saveEntitlementErr := appCtx.DB().Save(&entitlement); saveEntitlementErr != nil {
				err = apperror.NewBadDataError("Error saving entitlement")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
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
				err = apperror.NewBadDataError("Error fetching contractor")
				appCtx.Logger().Error(err.Error())
				return orderop.NewCreateOrderUnprocessableEntity(), err
			}

			packingAndShippingInstructions := models.InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + models.InstructionsAfterContractNumber
			newOrder, verrs, err := serviceMember.CreateOrder(
				appCtx,
				time.Time(*payload.IssueDate),
				time.Time(*payload.ReportByDate),
				(internalmessages.OrdersType)(*payload.OrdersType),
				*payload.HasDependents,
				*payload.SpouseHasProGear,
				newDutyLocation,
				payload.OrdersNumber,
				payload.Tac,
				payload.Sac,
				deptIndicator,
				&originDutyLocation,
				&grade,
				&entitlement,
				&originDutyLocationGBLOC.GBLOC,
				packingAndShippingInstructions,
				&newDutyLocationGBLOC.GBLOC,
			)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			var status models.MoveStatus

			if appCtx.Session().IsOfficeApp() {
				status = models.MoveStatusNeedsServiceCounseling
			} else {
				status = models.MoveStatusDRAFT
			}

			moveOptions := models.MoveOptions{
				Show:   models.BoolPointer(true),
				Status: &status,
			}
			if !appCtx.Session().OfficeUserID.IsNil() {
				officeUser, err := models.FetchOfficeUserByID(appCtx.DB(), appCtx.Session().OfficeUserID)
				if err != nil {
					err = apperror.NewBadDataError("Unable to fetch office user.")
					appCtx.Logger().Error(err.Error())
					return orderop.NewCreateOrderUnprocessableEntity(), err
				} else {
					moveOptions.CounselingOfficeID = &officeUser.TransportationOfficeID
				}
			}

			if newOrder.OrdersType == "SAFETY" {
				// if creating a Safety move, clear out the OktaID for the customer since they won't log into MilMove
				err = models.UpdateUserOktaID(appCtx.DB(), &newOrder.ServiceMember.User, "")
				if err != nil {
					appCtx.Logger().Error("Authorization error updating user", zap.Error(err))
					return orderop.NewUpdateOrderInternalServerError(), err
				}
			}

			newMove, verrs, err := newOrder.CreateNewMove(appCtx.DB(), moveOptions)
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}
			newOrder.Moves = append(newOrder.Moves, *newMove)

			order := (models.Order)(newOrder)

			orderPayload := payloads.Order(&order)

			return orderop.NewCreateOrderOK().WithPayload(orderPayload), nil
		})
}

// UpdateAllowanceHandler updates an order and entitlements via PATCH /orders/{orderId}/allowances
type UpdateAllowanceHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h UpdateAllowanceHandler) Handle(params orderop.UpdateAllowanceParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order allowance", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateAllowanceNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewUpdateAllowanceUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateAllowancePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateAllowanceInternalServerError(), err
				}
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateAllowanceAsTOO(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdatedAllowanceEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateAllowanceOK().WithPayload(orderPayload), nil
		})
}

// CounselingUpdateAllowanceHandler updates an order and entitlements via PATCH /counseling/orders/{orderId}/allowances
type CounselingUpdateAllowanceHandler struct {
	handlers.HandlerConfig
	orderUpdater services.OrderUpdater
}

// Handle ... updates an order from a request payload
func (h CounselingUpdateAllowanceHandler) Handle(
	params orderop.CounselingUpdateAllowanceParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating order allowance", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return orderop.NewCounselingUpdateAllowanceNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return orderop.NewCounselingUpdateAllowanceUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewCounselingUpdateAllowancePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewCounselingUpdateAllowanceForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewCounselingUpdateAllowanceInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				return handleError(apperror.NewForbiddenError("is not a Services Counselor"))
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedOrder, moveID, err := h.orderUpdater.UpdateAllowanceAsCounselor(
				appCtx,
				orderID,
				*params.Body,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerCounselingUpdateAllowanceEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewCounselingUpdateAllowanceOK().WithPayload(orderPayload), nil
		})
}

// UpdateBillableWeightHandler updates the max billable weight on an order's entitlements via PATCH /orders/{orderId}/update-billable-weight
type UpdateBillableWeightHandler struct {
	handlers.HandlerConfig
	excessWeightRiskManager services.ExcessWeightRiskManager
}

// Handle ... updates the authorized weight
func (h UpdateBillableWeightHandler) Handle(
	params orderop.UpdateBillableWeightParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating max billable weight", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateBillableWeightNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return orderop.NewUpdateBillableWeightUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateBillableWeightPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateBillableWeightForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateBillableWeightInternalServerError(), err
				}
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			dbAuthorizedWeight := models.IntPointer(int(*params.Body.AuthorizedWeight))
			updatedOrder, moveID, err := h.excessWeightRiskManager.UpdateBillableWeightAsTOO(
				appCtx,
				orderID,
				dbAuthorizedWeight,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdatedBillableWeightEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateBillableWeightOK().WithPayload(orderPayload), nil
		})
}

// UpdateMaxBillableWeightAsTIOHandler updates the max billable weight on an order's entitlements via PATCH /orders/{orderId}/update-billable-weight/tio
type UpdateMaxBillableWeightAsTIOHandler struct {
	handlers.HandlerConfig
	excessWeightRiskManager services.ExcessWeightRiskManager
}

// Handle ... updates the authorized weight
func (h UpdateMaxBillableWeightAsTIOHandler) Handle(
	params orderop.UpdateMaxBillableWeightAsTIOParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error updating max billable weight", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return orderop.NewUpdateMaxBillableWeightAsTIONotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return orderop.NewUpdateMaxBillableWeightAsTIOUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewUpdateMaxBillableWeightAsTIOPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewUpdateMaxBillableWeightAsTIOForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewUpdateMaxBillableWeightAsTIOInternalServerError(), err
				}
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			dbAuthorizedWeight := models.IntPointer(int(*params.Body.AuthorizedWeight))
			remarks := params.Body.TioRemarks
			updatedOrder, moveID, err := h.excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(
				appCtx,
				orderID,
				dbAuthorizedWeight,
				remarks,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerUpdatedMaxBillableWeightAsTIOEvent(appCtx, orderID, moveID, params)

			orderPayload := payloads.Order(updatedOrder)

			return orderop.NewUpdateMaxBillableWeightAsTIOOK().WithPayload(orderPayload), nil
		})
}

// AcknowledgeExcessWeightRiskHandler is called when a TOO dismissed the alert to acknowledge the excess weight risk via POST /orders/{orderId}/acknowledge-excess-weight-risk
type AcknowledgeExcessWeightRiskHandler struct {
	handlers.HandlerConfig
	excessWeightRiskManager services.ExcessWeightRiskManager
}

// Handle ... updates the authorized weight
func (h AcknowledgeExcessWeightRiskHandler) Handle(
	params orderop.AcknowledgeExcessWeightRiskParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// handleError is a reusable function to deal with multiple errors
			// when it comes to updating orders.
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("error acknowledging excess weight risk", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return orderop.NewAcknowledgeExcessWeightRiskNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(handlers.ValidationErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return orderop.NewAcknowledgeExcessWeightRiskUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return orderop.NewAcknowledgeExcessWeightRiskPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				case apperror.ForbiddenError:
					return orderop.NewAcknowledgeExcessWeightRiskForbidden().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return orderop.NewAcknowledgeExcessWeightRiskInternalServerError(), err
				}
			}

			orderID := uuid.FromStringOrNil(params.OrderID.String())
			updatedMove, err := h.excessWeightRiskManager.AcknowledgeExcessWeightRisk(
				appCtx,
				orderID,
				params.IfMatch,
			)
			if err != nil {
				return handleError(err)
			}

			h.triggerAcknowledgeExcessWeightRiskEvent(appCtx, updatedMove.ID, params)

			movePayload, err := payloads.Move(updatedMove, h.FileStorer())
			if err != nil {
				return orderop.NewAcknowledgeExcessWeightRiskInternalServerError(), err
			}

			return orderop.NewAcknowledgeExcessWeightRiskOK().WithPayload(movePayload), nil
		})
}

func (h UpdateOrderHandler) triggerUpdateOrderEvent(
	appCtx appcontext.AppContext,
	orderID uuid.UUID,
	moveID uuid.UUID,
	params orderop.UpdateOrderParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateOrderHandler could not generate the event")
	}
}

func (h CounselingUpdateOrderHandler) triggerCounselingUpdateOrderEvent(
	appCtx appcontext.AppContext,
	orderID uuid.UUID,
	moveID uuid.UUID,
	params orderop.CounselingUpdateOrderParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcCounselingUpdateOrderEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateAllowanceHandler could not generate the event")
	}
}

func (h UpdateAllowanceHandler) triggerUpdatedAllowanceEvent(
	appCtx appcontext.AppContext,
	orderID uuid.UUID,
	moveID uuid.UUID,
	params orderop.UpdateAllowanceParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateAllowanceEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateAllowanceHandler could not generate the event")
	}
}

func (h CounselingUpdateAllowanceHandler) triggerCounselingUpdateAllowanceEvent(
	appCtx appcontext.AppContext,
	orderID uuid.UUID,
	moveID uuid.UUID,
	params orderop.CounselingUpdateAllowanceParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcCounselingUpdateAllowanceEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().
			Error("ghcapi.CounselingUpdateAllowanceHandler could not generate the event")
	}
}

func (h UpdateBillableWeightHandler) triggerUpdatedBillableWeightEvent(
	appCtx appcontext.AppContext,
	orderID uuid.UUID,
	moveID uuid.UUID,
	params orderop.UpdateBillableWeightParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateBillableWeightEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateBillableWeightHandler could not generate the event")
	}
}

func (h UpdateMaxBillableWeightAsTIOHandler) triggerUpdatedMaxBillableWeightAsTIOEvent(
	appCtx appcontext.AppContext,
	orderID uuid.UUID,
	moveID uuid.UUID,
	params orderop.UpdateMaxBillableWeightAsTIOParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcUpdateMaxBillableWeightAsTIOEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.OrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: orderID,                   // ID of the updated logical object
		MtoID:           moveID,                    // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().
			Error("ghcapi.UpdateMaxBillableWeightAsTIOHandler could not generate the event")
	}
}

func (h AcknowledgeExcessWeightRiskHandler) triggerAcknowledgeExcessWeightRiskEvent(
	appCtx appcontext.AppContext,
	moveID uuid.UUID,
	params orderop.AcknowledgeExcessWeightRiskParams,
) {
	_, err := event.TriggerEvent(event.Event{
		EndpointKey: event.GhcAcknowledgeExcessWeightRiskEndpointKey,
		// Endpoint that is being handled
		EventKey:        event.MoveTaskOrderUpdateEventKey, // Event that you want to trigger
		UpdatedObjectID: moveID,                            // ID of the updated logical object
		MtoID:           moveID,                            // ID of the associated Move
		AppContext:      appCtx,
		TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
	})
	// If the event trigger fails, just log the error.
	if err != nil {
		appCtx.Logger().Error("ghcapi.UpdateBillableWeightHandler could not generate the event")
	}
}

func PayloadForOrdersModel(order models.Order) (*ghcmessages.OrderBody, error) {
	payload := &ghcmessages.OrderBody{
		ID: *handlers.FmtUUID(order.ID),
	}

	return payload, nil
}

// UploadAmendedOrdersHandler uploads amended orders to an order via POST /orders/{orderId}/upload_amended_orders
type UploadAmendedOrdersHandler struct {
	handlers.HandlerConfig
	services.OrderUpdater
}

// Handle updates an order to attach amended orders from a request payload
func (h UploadAmendedOrdersHandler) Handle(params orderop.UploadAmendedOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			file, ok := params.File.(*runtime.File)
			if !ok {
				errMsg := "This should always be a runtime.File, something has changed in go-swagger."

				appCtx.Logger().Error(errMsg)

				return orderop.NewUploadAmendedOrdersInternalServerError(), nil
			}

			appCtx.Logger().Info(
				"File uploader and size",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
				zap.String("officeUserID", appCtx.Session().OfficeUserID.String()),
				zap.String("AdminUserID", appCtx.Session().AdminUserID.String()),
				zap.Int64("size", file.Header.Size),
			)

			orderID, err := uuid.FromString(params.OrderID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			upload, url, verrs, err := h.OrderUpdater.UploadAmendedOrdersAsOffice(appCtx, appCtx.Session().UserID, orderID, file.Data, file.Header.Filename, h.FileStorer())

			if verrs.HasAny() || err != nil {
				switch err.(type) {
				case uploader.ErrTooLarge:
					return orderop.NewUploadAmendedOrdersRequestEntityTooLarge(), err
				case uploader.ErrFile:
					return orderop.NewUploadAmendedOrdersInternalServerError(), err
				case uploader.ErrFailedToInitUploader:
					return orderop.NewUploadAmendedOrdersInternalServerError(), err
				case apperror.NotFoundError:
					return orderop.NewUploadAmendedOrdersNotFound(), err
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
				}
			}

			uploadPayload, err := payloadForUploadModelFromAmendedOrdersUpload(h.FileStorer(), upload, url)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return orderop.NewUploadAmendedOrdersCreated().WithPayload(uploadPayload), nil
		})
}

func payloadForUploadModelFromAmendedOrdersUpload(storer storage.FileStorer, upload models.Upload, url string) (*ghcmessages.Upload, error) {
	uploadPayload := &ghcmessages.Upload{
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
