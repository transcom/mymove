package paymentrequest

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// FetchServiceParamList fetches the service param list
func (p *RequestPaymentHelper) FetchServiceParamList(mtoServiceID uuid.UUID) (models.ServiceParams, error) {
	mtoServiceItem := models.MTOServiceItem{}
	serviceParams := models.ServiceParams{}

	err := p.DB.Where("id = ?", mtoServiceID).First(&mtoServiceItem)
	if err != nil {
		return nil, fmt.Errorf("failure fetching MTO Service Item: %w", err)
	}

	// Get all service item param keys that do not come from pricers
	err = p.DB.Q().
		InnerJoin("service_item_param_keys sipk", "service_params.service_item_param_key_id = sipk.id").
		Where("service_id = ? AND sipk.origin <> ?", mtoServiceItem.ReServiceID, models.ServiceItemParamOriginPricer).
		Eager().All(&serviceParams)
	if err != nil {
		return nil, fmt.Errorf("failure fetching service params for MTO Service Item ID <%s> with RE Service Item ID <%s>: %w", mtoServiceID.String(), mtoServiceItem.ReServiceID.String(), err)
	}

	return serviceParams, err
}
