package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ServiceAreaLookup does lookup of service area based on address postal code
type ServiceAreaLookup struct {
	Address models.Address
}

func (r ServiceAreaLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	zip := r.Address.PostalCode
	zip3 := zip[0:3]

	domesticServiceArea, err := fetchDomesticServiceArea(appCtx, keyData.ContractCode, zip3)

	return domesticServiceArea.ServiceArea, err
}

func (r ServiceAreaLookup) ParamValue(appCtx appcontext.AppContext, contractCode string) (string, error) {
	return r.lookup(appCtx, &ServiceItemParamKeyData{ContractCode: contractCode})
}
