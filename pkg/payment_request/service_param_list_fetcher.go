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

// FetchServiceParamsForServiceItems fetches the service param list.  Returns a slice of ServiceItemParamKey models, each with an
// eagerly fetched ServiceItemParamKey association.
func (p *RequestPaymentHelper) FetchServiceParamsForServiceItems(appCtx appcontext.AppContext, mtoServiceItems []models.MTOServiceItem) (models.ServiceParams, error) {
	serviceItemCodes := make([]models.ReServiceCode, len(mtoServiceItems))
	for i, mtoServiceItem := range mtoServiceItems {
		serviceItemCodes[i] = mtoServiceItem.ReService.Code
	}

	// Get all service item system param keys
	var serviceParams models.ServiceParams
	err := appCtx.DB().EagerPreload("ServiceItemParamKey", "Service").
		InnerJoin("service_item_param_keys", "service_params.service_item_param_key_id = service_item_param_keys.id").
		InnerJoin("re_services", "re_services.id = service_params.service_id").
		Where("origin IN (?)", models.ServiceItemParamOriginSystem, models.ServiceItemParamOriginPrime).
		Where("re_services.code IN (?)", serviceItemCodes).
		All(&serviceParams)
	if err != nil {
		return nil, fmt.Errorf("failure fetching service params for RE Service Item IDs <%s>: %w", serviceItemCodes, err)
	}

	return serviceParams, err
}
