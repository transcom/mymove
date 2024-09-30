package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// InternationalRateAreaLookup does lookup of rate area based on address postal code
type InternationalRateAreaLookup struct {
	Address models.Address
}

// Looks up international rate area based on given address postal code.
// Domestic use should refer to the "Service Area" lookup
func (r InternationalRateAreaLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	zip := r.Address.PostalCode
	zip5 := zip[0:5]

	internationalRateArea, err := fetchInternationalRateArea(appCtx, keyData.ContractCode, zip5)

	return internationalRateArea.Name, err
}

func (r InternationalRateAreaLookup) ParamValue(appCtx appcontext.AppContext, contractCode string) (string, error) {
	return r.lookup(appCtx, &ServiceItemParamKeyData{ContractCode: contractCode})
}
