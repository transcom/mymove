package customer

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type fetchCustomer struct{}

// NewCustomerFetcher creates a new struct with the service dependencies
func NewCustomerFetcher() services.CustomerFetcher {
	return &fetchCustomer{}
}

// FetchCustomer retrieves a Customer for a given UUID
func (f fetchCustomer) FetchCustomer(appCtx appcontext.AppContext, customerID uuid.UUID) (*models.ServiceMember, error) {
	customer := &models.ServiceMember{}
	if err := appCtx.DB().EagerPreload("ResidentialAddress.Country", "BackupMailingAddress.Country", "BackupContacts").Find(customer, customerID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.ServiceMember{}, apperror.NewNotFoundError(customerID, "")
		default:
			return &models.ServiceMember{}, apperror.NewQueryError("ServiceMember", err, "")
		}
	}

	if customer.ResidentialAddress != nil {
		if customer.ResidentialAddress.IsOconus == nil {
			// Evaluate address and populate addresses isOconus value
			isOconus, err := models.IsAddressOconus(appCtx.DB(), *customer.ResidentialAddress)
			if err != nil {
				return nil, err
			}
			customer.ResidentialAddress.IsOconus = &isOconus
		}
	}

	if customer.BackupMailingAddress != nil {
		if customer.BackupMailingAddress.IsOconus == nil {
			// Evaluate address and populate addresses isOconus value
			isOconus, err := models.IsAddressOconus(appCtx.DB(), *customer.BackupMailingAddress)
			if err != nil {
				return nil, err
			}
			customer.BackupMailingAddress.IsOconus = &isOconus
		}
	}

	return customer, nil
}
