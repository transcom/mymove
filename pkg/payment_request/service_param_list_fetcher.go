package payment_request

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
)

func (p *PaymentRequestHelper) FetchServiceParamList(serviceID uuid.UUID) (models.ServiceParams, error) {
	serviceParams := models.ServiceParams{}

	err := p.DB.Where("service_id = ?", serviceID).All(&serviceParams)
	if err != nil {
		return nil, fmt.Errorf("failure fetching service params: %w", err)
	}

	return serviceParams, err
}

