package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
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

			customerSupportRemark, err := h.UpdateCustomerSupportRemark(appCtx, params)
			if err != nil {
				appCtx.Logger().Error("Error updating customer support remark: ", zap.Error(err))
				return customersupportremarksop.NewUpdateCustomerSupportRemarkForMoveInternalServerError(), err
			}

			returnPayload := payloads.CustomerSupportRemark(customerSupportRemark)

			return customersupportremarksop.NewUpdateCustomerSupportRemarkForMoveOK().WithPayload(returnPayload), nil
		})
}
