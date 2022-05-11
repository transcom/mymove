package ghcapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// GetCustomerHandler fetches the information of a specific customer
type GetCustomerHandler struct {
	handlers.HandlerContext
	services.CustomerFetcher
}

// Handle getting the information of a specific customer
func (h GetCustomerHandler) Handle(params customercodeop.GetCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			customerID, _ := uuid.FromString(params.CustomerID.String())
			customer, err := h.FetchCustomer(appCtx, customerID)
			if err != nil {
				appCtx.Logger().Error("Loading Customer Info", zap.Error(err))
				switch err {
				case sql.ErrNoRows:
					return customercodeop.NewGetCustomerNotFound(), err
				default:
					return customercodeop.NewGetCustomerInternalServerError(), err
				}
			}
			customerInfoPayload := payloads.Customer(customer)
			return customercodeop.NewGetCustomerOK().WithPayload(customerInfoPayload), nil
		})
}

// UpdateCustomerHandler updates a customer via PATCH /customer/{customerId}
type UpdateCustomerHandler struct {
	handlers.HandlerContext
	customerUpdater services.CustomerUpdater
}

// Handle updates a customer from a request payload
func (h UpdateCustomerHandler) Handle(params customercodeop.UpdateCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				forbiddenError := apperror.NewForbiddenError("user is not authenticated with service counselor office role")
				appCtx.Logger().Error(forbiddenError.Error())
				return customercodeop.NewUpdateCustomerForbidden(), forbiddenError
			}

			customerID, err := uuid.FromString(params.CustomerID.String())
			if err != nil {
				appCtx.Logger().Error("unable to parse customer id param to uuid", zap.Error(err))
				return customercodeop.NewUpdateCustomerBadRequest(), err
			}

			newCustomer := payloads.CustomerToServiceMember(*params.Body)
			newCustomer.ID = customerID

			updatedCustomer, err := h.customerUpdater.UpdateCustomer(appCtx, params.IfMatch, newCustomer)

			if err != nil {
				appCtx.Logger().Error("error updating customer", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return customercodeop.NewGetCustomerNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewUpdateCustomerUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return customercodeop.NewUpdateCustomerPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return customercodeop.NewUpdateCustomerInternalServerError(), err
				}
			}

			customerPayload := payloads.Customer(updatedCustomer)

			return customercodeop.NewUpdateCustomerOK().WithPayload(customerPayload), nil
		})
}
