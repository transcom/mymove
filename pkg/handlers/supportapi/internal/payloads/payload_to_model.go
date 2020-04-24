package payloads

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerModel model
func CustomerModel(customer *supportmessages.Customer) *models.Customer {
	if customer == nil {
		return nil
	}
	return &models.Customer{
		ID:          uuid.FromStringOrNil(customer.ID.String()),
		Agency:      &customer.Agency,
		FirstName:   &customer.FirstName,
		LastName:    &customer.LastName,
		DODID:       &customer.DodID,
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
	}
}
