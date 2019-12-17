package ghcapi

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetCustomerInfoHandler fetches the information of a specific customer
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
