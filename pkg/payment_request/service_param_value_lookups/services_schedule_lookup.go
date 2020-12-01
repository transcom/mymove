package serviceparamvaluelookups

import (
	"fmt"

	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

// ServicesScheduleLookup does lookup on services schedule origin
type ServicesScheduleLookup struct {
	Address models.Address
}

func (s ServicesScheduleLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// find the service area by querying for the service area associated with the zip3
	zip := s.Address.PostalCode
	zip3 := zip[0:3]

	var domesticServiceArea models.ReDomesticServiceArea
	err := db.Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("re_zip3s.zip3 = ?", zip3).
		Where("re_contracts.code = ?", keyData.ContractCode).
		First(&domesticServiceArea)
	if err != nil {
		return "", fmt.Errorf("unable to find domestic service area for %s under contract code %s", zip3, keyData.ContractCode)
	}

	return strconv.Itoa(domesticServiceArea.ServicesSchedule), nil
}
