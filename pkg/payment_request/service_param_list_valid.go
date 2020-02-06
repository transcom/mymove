package paymentrequest

import (
	"github.com/transcom/mymove/pkg/models"
)

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
		if found == false {
			hasError = true
			errorString = errorString + " Param Key <" + serviceParam.ServiceItemParamKey.Key + ">"
		}
	}

	if hasError {
		errorMessage := " MTO Service Item <" + mtoServiceItem.ID.String() + "> missing params needed for pricing: " + errorString
		return !hasError, errorMessage
	}

	return !hasError, ""
}