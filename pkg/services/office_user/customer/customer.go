package customer

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type fetchCustomer struct{}

// NewCustomerFetcher creates a new struct with the service dependencies
func NewCustomerFetcher() services.CustomerFetcher {
	return &fetchCustomer{}
}

//FetchCustomer retrieves a Customer for a given UUID
func (f fetchCustomer) FetchCustomer(appCfg appconfig.AppConfig, customerID uuid.UUID) (*models.ServiceMember, error) {
	customer := &models.ServiceMember{}
	if err := appCfg.DB().Eager().Find(customer, customerID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.ServiceMember{}, services.NewNotFoundError(customerID, "")
		default:
			return &models.ServiceMember{}, err
		}
	}
	return customer, nil
}
