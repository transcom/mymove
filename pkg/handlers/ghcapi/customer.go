package ghcapi

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"

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

// UpdateCustomerHandler updates a customer via PATCH /customer/{customerId}
type UpdateCustomerHandler struct {
	handlers.HandlerContext
	customerUpdater services.CustomerUpdater
}

// Handle updates a customer from a request payload
func (h UpdateCustomerHandler) Handle(params customercodeop.UpdateCustomerParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

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

	updatedCustomer, err := h.customerUpdater.UpdateCustomer(params.IfMatch, newCustomer)

	if err != nil {
		logger.Error("error updating customer", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return customercodeop.NewGetCustomerBadRequest()
		case services.InvalidInputError:
			payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceID(), validate.NewErrors())
			return customercodeop.NewUpdateCustomerUnprocessableEntity().WithPayload(payload)
		case services.PreconditionFailedError:
			return customercodeop.NewUpdateCustomerPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return customercodeop.NewUpdateCustomerInternalServerError()
		}
	}

	var customer models.ServiceMember
	query := h.DB().Where("user_id = ?", updatedCustomer.ID)
	err = query.First(&customer)

	if err != nil {
		logger.Error("ghcapi.UpdateCustomerHandler could not find customer")
	}

	customerPayload := payloads.Customer(updatedCustomer)

	return customercodeop.NewUpdateCustomerOK().WithPayload(customerPayload)
}

// Customer transforms UpdateCustomerPayload to ServiceMember model
func Customer(payload ghcmessages.UpdateCustomerPayload) (models.ServiceMember, error) {
	// TODO: move this to internal/models payload_to_modal?

	var address = models.Address{
		ID:             uuid.FromStringOrNil(payload.CurrentAddress.ID.String()),
		StreetAddress1: *payload.CurrentAddress.StreetAddress1,
		StreetAddress2: payload.CurrentAddress.StreetAddress2,
		StreetAddress3: payload.CurrentAddress.StreetAddress3,
		City:           *payload.CurrentAddress.City,
		State:          *payload.CurrentAddress.State,
		PostalCode:     *payload.CurrentAddress.PostalCode,
		Country:        payload.CurrentAddress.Country,
	}

	var backupContact = models.BackupContact{
		Email: *payload.BackupContact.Email,
		Name:  *payload.BackupContact.Name,
		Phone: payload.BackupContact.Phone,
	}

	var backupContacts []models.BackupContact
	backupContacts = append(backupContacts, backupContact)

	return models.ServiceMember{
		ResidentialAddress: &address,
		BackupContacts:     backupContacts,
		FirstName:          &payload.FirstName,
		LastName:           &payload.LastName,
		PersonalEmail:      payload.Email,
		Telephone:          payload.Phone,
	}, nil
}
