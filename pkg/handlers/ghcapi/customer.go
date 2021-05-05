package ghcapi

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/models"
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
	logger := h.LoggerFromRequest(params.HTTPRequest)
	customerID, _ := uuid.FromString(params.CustomerID.String())
	customer, err := h.FetchCustomer(customerID)
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

type UpdateCustomerHandler struct {
	handlers.HandlerContext
	customerUpdater services.CustomerUpdater
}

// Handle updates a customer from a request payload
func (h UpdateCustomerHandler) Handle(params customercodeop.UpdateCustomerParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	// TODO: is CustomerID correct id to use to search on?
	customerID, err := uuid.FromString(params.CustomerID.String())
	if err != nil {
		logger.Error("unable to parse customer id param to uuid", zap.Error(err))
		return customercodeop.NewUpdateCustomerBadRequest()
	}

	newCustomer, err := Customer(*params.Body)
	if err != nil {
		logger.Error("error converting payload to service member model", zap.Error(err))
		return customercodeop.NewUpdateCustomerBadRequest()
	}
	newCustomer.ID = customerID
	// TODO: finish
}

// Customer transforms
func Customer(payload ghcmessages.UpdateCustomerPayload) (models.ServiceMember, error) {
	// TOOD: finish
}
