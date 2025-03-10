package paymentrequest

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ValidServiceParamList validates service params
func (p *RequestPaymentHelper) ValidServiceParamList(appCtx appcontext.AppContext, mtoServiceItem models.MTOServiceItem, serviceParams models.ServiceParams, paymentServiceItemParams models.PaymentServiceItemParams) (bool, string) {
	var errorString string
	hasError := false

	// Resolve the ReService for our lookup
	reServiceForLookup, err := resolveReServiceForLookup(appCtx, mtoServiceItem)
	if err != nil {
		return true, err.Error()
	}

	for _, serviceParam := range serviceParams {
		if serviceParam.IsOptional {
			// Some params are considered optional.  If this is one, then we can skip looking for it.
			continue
		}

		found := false
		for _, paymentServiceItemParam := range paymentServiceItemParams {
			if serviceParam.ServiceItemParamKey.Key == paymentServiceItemParam.ServiceItemParamKey.Key &&
				serviceParam.ServiceID.String() == reServiceForLookup.ID.String() {
				found = true
			}
		}
		if !found {
			hasError = true
			errorString = fmt.Sprintf("%s Param Key <%s>", errorString, serviceParam.ServiceItemParamKey.Key)
		}
	}

	if hasError {
		errorMessage := " MTO Service Item <" + reServiceForLookup.ID.String() + "> missing params needed for pricing: " + errorString
		return !hasError, errorMessage
	}

	return !hasError, ""
}
