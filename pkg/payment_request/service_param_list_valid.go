package paymentrequest

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// ValidServiceParamList validates service params
func (p *RequestPaymentHelper) ValidServiceParamList(mtoServiceItem models.MTOServiceItem, serviceParams models.ServiceParams, paymentServiceItemParams models.PaymentServiceItemParams) (bool, string) {
	var errorString string
	hasError := false
	for _, serviceParam := range serviceParams {
		found := false
		for _, paymentServiceItemParam := range paymentServiceItemParams {
			if serviceParam.ServiceItemParamKey.Key == paymentServiceItemParam.ServiceItemParamKey.Key &&
				serviceParam.ServiceID.String() == mtoServiceItem.ReServiceID.String() {
				found = true
			}
		}
		if !found {
			if serviceParam.ServiceItemParamKey.Key == models.ServiceItemParamNameWeightEstimated {
				// This field is optional for a shipment, if the lookups did not error out when setting
				// the weights, then finding this param in the paymentServiceItemParams is not critical
				found = true
			} else {
				hasError = true
				errorString = fmt.Sprintf("%s Param Key <%s>", errorString, serviceParam.ServiceItemParamKey.Key)
			}
		}
	}

	if hasError {
		errorMessage := " MTO Service Item <" + mtoServiceItem.ID.String() + "> missing params needed for pricing: " + errorString
		return !hasError, errorMessage
	}

	return !hasError, ""
}
