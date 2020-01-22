package paymentrequest

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (p *RequestPaymentHelper) FetchServiceParamList(mtoServiceID uuid.UUID) (models.ServiceParams, error) {
	mtoServiceItem := models.MTOServiceItem{}
	serviceParams := models.ServiceParams{}

	err := p.DB.Where("id = ?", mtoServiceID).First(&mtoServiceItem)
	if err != nil {
		return nil, fmt.Errorf("failure fetching MTO Service Item: %w", err)
	}

	err = p.DB.Where("service_id = ?", mtoServiceItem.ReServiceID).All(&serviceParams)
	if err != nil {
		return nil, fmt.Errorf("failure fetching service params for MTO Service Item ID <%s> with RE Service Item ID <%s>: %w", mtoServiceID.String(), mtoServiceItem.ReServiceID.String(), err)
	}

	return serviceParams, err
}
