package ghcapi

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
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
	handlers.HandlerConfig
	services.CustomerFetcher
}

// Handle getting the information of a specific customer
func (h GetCustomerHandler) Handle(params customercodeop.GetCustomerParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)
	customerID, _ := uuid.FromString(params.CustomerID.String())
	customer, err := h.FetchCustomer(appCtx, customerID)
	if err != nil {
		logger.Error("Loading Customer Info", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return customercodeop.NewGetCustomerNotFound()
		default:
			return customercodeop.NewGetCustomerInternalServerError()
		}
	}
	customerInfoPayload := payloads.Customer(customer)
	return customercodeop.NewGetCustomerOK().WithPayload(customerInfoPayload)
}

// UpdateCustomerHandler updates a customer via PATCH /customer/{customerId}
type UpdateCustomerHandler struct {
	handlers.HandlerConfig
	customerUpdater services.CustomerUpdater
}

// Handle updates a customer from a request payload
func (h UpdateCustomerHandler) Handle(params customercodeop.UpdateCustomerParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	if !session.IsOfficeUser() || !session.Roles.HasRole(roles.RoleTypeServicesCounselor) {
		logger.Error("user is not authenticated with service counselor office role")
		return customercodeop.NewUpdateCustomerForbidden()
	}

	customerID, err := uuid.FromString(params.CustomerID.String())
	if err != nil {
		logger.Error("unable to parse customer id param to uuid", zap.Error(err))
		return customercodeop.NewUpdateCustomerBadRequest()
	}

	newCustomer := payloads.CustomerToServiceMember(*params.Body)
	newCustomer.ID = customerID

	updatedCustomer, err := h.customerUpdater.UpdateCustomer(appCtx, params.IfMatch, newCustomer)

	if err != nil {
		logger.Error("error updating customer", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return customercodeop.NewGetCustomerNotFound()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return customercodeop.NewUpdateCustomerUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return customercodeop.NewUpdateCustomerPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return customercodeop.NewUpdateCustomerInternalServerError()
		}
	}

	customerPayload := payloads.Customer(updatedCustomer)

	return customercodeop.NewUpdateCustomerOK().WithPayload(customerPayload)
}
