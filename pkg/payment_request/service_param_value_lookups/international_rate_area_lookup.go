package serviceparamvaluelookups

import (
	"fmt"

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
	if len(zip) < 5 {
		// Looking up without zip5 not supported yet
		return "", fmt.Errorf("looking up the international rate area for addresses without zip5 codes is not supported yet")
	}
	zip5 := zip[0:5]

	internationalRateArea, err := fetchInternationalRateArea(appCtx, keyData.ContractCode, zip5)

	return internationalRateArea.Name, err
}

func (r InternationalRateAreaLookup) ParamValue(appCtx appcontext.AppContext, contractCode string) (string, error) {
	return r.lookup(appCtx, &ServiceItemParamKeyData{ContractCode: contractCode})
}
