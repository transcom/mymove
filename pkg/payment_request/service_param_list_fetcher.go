package paymentrequest

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ResolveReServiceForLookup ensures that the correct ReService is used for parameter lookup
// This is because some service items don't have parameters that can be looked up because they inherit the logic from existing items.
// For example, INPK. INPK is for iHHG shipments going into non-temporary storage.
// This means we are packing an iHHG shipment, so we price by IHPK, but with a special
// pricer for INPK. INPK is iHHG -> iNTS. Prices by IHPK multiplied by NTS market factor
func ResolveReServiceForLookup(appCtx appcontext.AppContext, mtoServiceItem models.MTOServiceItem) (models.ReService, error) {
	var reService models.ReService

	if mtoServiceItem.ReService.Code != "" {
		reService = mtoServiceItem.ReService
	} else {
		reServicePtr, err := models.FetchReServiceByCode(appCtx.DB(), models.ReServiceCodeIHPK)
		if err != nil {
			return models.ReService{}, err
		}
		reService = *reServicePtr
	}

	// Handle special cases where we need to swap lookup services
	switch reService.Code {
	case models.ReServiceCodeINPK:
		// INPK is priced using IHPK parameters but with an NTS market factor multiplier.
		reServicePtr, err := models.FetchReServiceByCode(appCtx.DB(), models.ReServiceCodeIHPK)
		if err != nil {
			return models.ReService{}, fmt.Errorf("failed to fetch IHPK for INPK lookup: %w", err)
		}
		reService = *reServicePtr
	}

	return reService, nil
}

// FetchServiceParamList fetches the service param list.  Returns a slice of ServiceParam models, each with an
// eagerly fetched ServiceItemParamKey association.
func (p *RequestPaymentHelper) FetchServiceParamList(appCtx appcontext.AppContext, mtoServiceItem models.MTOServiceItem) (models.ServiceParams, error) {

	// Resolve the ReService for our lookup
	reServiceForLookup, err := ResolveReServiceForLookup(appCtx, mtoServiceItem)
	if err != nil {
		return nil, err
	}

	// Get all service item param keys that do not come from pricers
	// using the ReService identified above
	var serviceParams models.ServiceParams
	err = appCtx.DB().Q().
		InnerJoin("service_item_param_keys sipk", "service_params.service_item_param_key_id = sipk.id").
		Where("service_id = ?", reServiceForLookup.ID).
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
