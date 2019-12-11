package ghcapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

// GetCustomerInfoHandler fetches the information of a specific customer
type GetCustomerHandler struct {
	handlers.HandlerContext
}

// Handle getting the information of a specific customer
func (h GetCustomerHandler) Handle(params customercodeop.GetCustomerParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	customerID, _ := uuid.FromString(params.CustomerID.String())
	customer := &models.Customer{}
	err := h.DB().Find(customer, customerID)
	if err != nil {
		logger.Error("Loading Customer Info", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return customercodeop.NewGetCustomerNotFound()
		default:
			return customercodeop.NewGetCustomerInternalServerError()
		}
	}
	customerInfoPayload := payloads.PayloadForCustomer(customer)
	return customercodeop.NewGetCustomerOK().WithPayload(customerInfoPayload)
}
