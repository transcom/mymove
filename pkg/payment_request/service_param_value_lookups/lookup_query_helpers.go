package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

func fetchDomesticServiceArea(db *pop.Connection, contractCode string, zip3 string) (models.ReDomesticServiceArea, error) {
	var domesticServiceArea models.ReDomesticServiceArea
	err := db.Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("re_zip3s.zip3 = ?", zip3).
		Where("re_contracts.code = ?", contractCode).
		First(&domesticServiceArea)
	if err != nil {
		return domesticServiceArea, fmt.Errorf("unable to find domestic service area for %s under contract code %s", zip3, contractCode)
	}

	return domesticServiceArea, nil
}
