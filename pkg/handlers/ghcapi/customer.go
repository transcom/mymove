package ghcapi

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetCustomerHandler fetches the information of a specific customer
type GetCustomerHandler struct {
	handlers.HandlerContext
	services.CustomerFetcher
}

// Handle getting the information of a specific customer
func (h GetCustomerHandler) Handle(params customercodeop.GetCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			customerID, _ := uuid.FromString(params.CustomerID.String())
			customer, err := h.FetchCustomer(appCtx, customerID)
			if err != nil {
				appCtx.Logger().Error("Loading Customer Info", zap.Error(err))
				switch err {
				case sql.ErrNoRows:
					return customercodeop.NewGetCustomerNotFound()
				default:
					return customercodeop.NewGetCustomerInternalServerError()
				}
			}
			customerInfoPayload := payloads.Customer(customer)
			return customercodeop.NewGetCustomerOK().WithPayload(customerInfoPayload)
		})
}

// UpdateCustomerHandler updates a customer via PATCH /customer/{customerId}
type UpdateCustomerHandler struct {
	handlers.HandlerContext
	customerUpdater services.CustomerUpdater
}

// Handle updates a customer from a request payload
func (h UpdateCustomerHandler) Handle(params customercodeop.UpdateCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				appCtx.Logger().Error("user is not authenticated with service counselor office role")
				return customercodeop.NewUpdateCustomerForbidden()
			}

			customerID, err := uuid.FromString(params.CustomerID.String())
			if err != nil {
				appCtx.Logger().Error("unable to parse customer id param to uuid", zap.Error(err))
				return customercodeop.NewUpdateCustomerBadRequest()
			}

			newCustomer := payloads.CustomerToServiceMember(*params.Body)
			newCustomer.ID = customerID

			updatedCustomer, err := h.customerUpdater.UpdateCustomer(appCtx, params.IfMatch, newCustomer)

			if err != nil {
				appCtx.Logger().Error("error updating customer", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return customercodeop.NewGetCustomerNotFound()
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewUpdateCustomerUnprocessableEntity().WithPayload(payload)
				case apperror.PreconditionFailedError:
					return customercodeop.NewUpdateCustomerPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
				default:
					return customercodeop.NewUpdateCustomerInternalServerError()
				}
			}

			customerPayload := payloads.Customer(updatedCustomer)

			return customercodeop.NewUpdateCustomerOK().WithPayload(customerPayload)
		})
}
