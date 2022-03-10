package paymentrequest

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// FetchServiceParamList fetches the service param list.  Returns a slice of ServiceParam models, each with an
// eagerly fetched ServiceItemParamKey association.
func (p *RequestPaymentHelper) FetchServiceParamList(appCtx appcontext.AppContext, mtoServiceItem models.MTOServiceItem) (models.ServiceParams, error) {
	// Get all service item param keys that do not come from pricers
	var serviceParams models.ServiceParams
	err := appCtx.DB().Q().
		InnerJoin("service_item_param_keys sipk", "service_params.service_item_param_key_id = sipk.id").
		Where("service_id = ?", mtoServiceItem.ReServiceID).
		Where("sipk.origin <> ?", models.ServiceItemParamOriginPricer).
		EagerPreload("ServiceItemParamKey").
		All(&serviceParams)
	if err != nil {
		return nil, fmt.Errorf("failure fetching service params for MTO Service Item ID <%s> with RE Service Item ID <%s>: %w", mtoServiceItem.ID.String(), mtoServiceItem.ReServiceID.String(), err)
	}

	return serviceParams, err
}
