package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (p *RequestPaymentHelper) ValidServiceParamList(serviceID uuid.UUID, serviceParams models.ServiceParams, paymentServiceItemParams models.PaymentServiceItemParams) (bool, *string) {

	//var errorMessage string

	// Use list of params needed (`serviceParams`) for service item
	// Use list of params saved for payment service item (`paymentServiceItemParams`)
	// Verify all params are present for payment service item

	return true, nil
}