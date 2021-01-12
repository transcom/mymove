package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/models"
)

// ServiceAreaLookup does lookup of service area based on address postal code
type ServiceAreaLookup struct {
	Address models.Address
}

func (r ServiceAreaLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	zip := r.Address.PostalCode
	zip3 := zip[0:3]

	domesticServiceArea, err := fetchDomesticServiceArea(&db, keyData.ContractCode, zip3)

	return domesticServiceArea.ServiceArea, err
}
