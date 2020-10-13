package customer

import (
	"database/sql"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type fetchCustomer struct {
	db *pop.Connection
}

// NewCustomerFetcher creates a new struct with the service dependencies
func NewCustomerFetcher(db *pop.Connection) services.CustomerFetcher {
	return &fetchCustomer{db}
}

//FetchCustomer retrieves a Customer for a given UUID
func (f fetchCustomer) FetchCustomer(customerID uuid.UUID) (*models.ServiceMember, error) {
	customer := &models.ServiceMember{}
	if err := f.db.Eager().Find(customer, customerID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.ServiceMember{}, services.NewNotFoundError(customerID, "")
		default:
			return &models.ServiceMember{}, err
		}
	}
	return customer, nil
}
