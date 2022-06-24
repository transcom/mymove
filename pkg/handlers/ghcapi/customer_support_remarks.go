package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"

	customersupportremarksop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
)

// ListCustomerSupportRemarksHandler is a struct that describes listing customer support remarks for a move
type ListCustomerSupportRemarksHandler struct {
	handlers.HandlerConfig
	services.CustomerSupportRemarksFetcher
}

type CreateCustomerSupportRemarksHandler struct {
	handlers.HandlerConfig
	services.CustomerSupportRemarksCreator
}

type UpdateCustomerSupportRemarkHandler struct {
	handlers.HandlerConfig
	services.CustomerSupportRemarkUpdater
}
type DeleteCustomerSupportRemarkHandler struct {
	handlers.HandlerConfig
	services.CustomerSupportRemarkDeleter
}

// Handle handles the handling for getting a list of customer support remarks for a move
func (h ListCustomerSupportRemarksHandler) Handle(params customersupportremarksop.GetCustomerSupportRemarksForMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			customerSupportRemarks, err := h.ListCustomerSupportRemarks(appCtx, params.Locator)
			if err != nil {
				if err == models.ErrFetchNotFound {
					appCtx.Logger().Error("Error fetching customer support remarks: ", zap.Error(err))
					return customersupportremarksop.NewGetCustomerSupportRemarksForMoveNotFound(), err
				}
				appCtx.Logger().Error("Error fetching customer support remarks: ", zap.Error(err))
				return customersupportremarksop.NewGetCustomerSupportRemarksForMoveInternalServerError(), err
			}

			returnPayload := payloads.CustomerSupportRemarks(*customerSupportRemarks)
			return customersupportremarksop.NewGetCustomerSupportRemarksForMoveOK().WithPayload(returnPayload), nil
		})
}

func (h CreateCustomerSupportRemarksHandler) Handle(params customersupportremarksop.CreateCustomerSupportRemarkForMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body

			remark := payloads.CustomerSupportRemarkModelFromCreate(payload)

			customerSupportRemark, err := h.CreateCustomerSupportRemark(appCtx, remark, params.Locator)
			if err != nil {
				appCtx.Logger().Error("Error creating customer support remark: ", zap.Error(err))
				return customersupportremarksop.NewCreateCustomerSupportRemarkForMoveInternalServerError(), err
			}

			returnPayload := payloads.CustomerSupportRemark(customerSupportRemark)

			return customersupportremarksop.NewCreateCustomerSupportRemarkForMoveOK().WithPayload(returnPayload), nil
		})
}

func (h UpdateCustomerSupportRemarkHandler) Handle(params customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body

			customerSupportRemark, err := h.UpdateCustomerSupportRemark(appCtx, *payload)
			if err != nil {
				appCtx.Logger().Error("Error updating customer support remark: ", zap.Error(err))
				return customersupportremarksop.NewUpdateCustomerSupportRemarkForMoveInternalServerError(), err
			}

			returnPayload := payloads.CustomerSupportRemark(customerSupportRemark)

			return customersupportremarksop.NewUpdateCustomerSupportRemarkForMoveOK().WithPayload(returnPayload), nil
		})
}
func (h DeleteCustomerSupportRemarkHandler) Handle(params customersupportremarksop.DeleteCustomerSupportRemarkParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// TODO what's the point of this conversion?
			remarkID := uuid.FromStringOrNil(params.CustomerSupportRemarkID.String())
			err := h.DeleteCustomerSupportRemark(appCtx, remarkID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteCustomerSupportRemarkHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return customersupportremarksop.NewDeleteCustomerSupportRemarkNotFound(), err
				case apperror.ConflictError:
					return customersupportremarksop.NewDeleteCustomerSupportRemarkConflict(), err
				case apperror.ForbiddenError:
					return customersupportremarksop.NewDeleteCustomerSupportRemarkForbidden(), err
				case apperror.UnprocessableEntityError:
					return customersupportremarksop.NewDeleteCustomerSupportRemarkUnprocessableEntity(), err
				default:
					return customersupportremarksop.NewDeleteCustomerSupportRemarkInternalServerError(), err
				}
			}

			// TODO do we need to trigger an event for this? it should be done here if so

			return customersupportremarksop.NewDeleteCustomerSupportRemarkNoContent(), nil
		})
}
