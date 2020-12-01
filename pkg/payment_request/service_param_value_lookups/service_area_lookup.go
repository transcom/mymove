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

	var domesticServiceArea models.ReDomesticServiceArea

	query := db.Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("zip3 = ?", zip3).
		Where("re_contracts.code = ?", keyData.ContractCode)

	err := query.First(&domesticServiceArea)

	return domesticServiceArea.ServiceArea, err
}
