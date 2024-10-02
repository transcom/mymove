package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func fetchDomesticServiceArea(appCtx appcontext.AppContext, contractCode string, zip3 string) (models.ReDomesticServiceArea, error) {
	var domesticServiceArea models.ReDomesticServiceArea
	err := appCtx.DB().Q().
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

// Currently looking up a rate area by zip5 is the only supported method
func fetchInternationalRateArea(appCtx appcontext.AppContext, contractCode string, zip5 string) (models.ReRateArea, error) {
	var internationalRateArea models.ReRateArea
	err := appCtx.DB().Q().
		Join("re_zip5_rate_areas", "re_zip5_rate_areas.rate_area_id = re_rate_areas.id").
		Join("re_contracts", "re_contracts.id = re_rate_areas.contract_id").
		Where("re_zip5_rate_areas.zip5 = ?", zip5).
		Where("re_contracts.code = ?", contractCode).
		First(&internationalRateArea)
	if err != nil {
		return internationalRateArea, fmt.Errorf("unable to find international rate area for %s under contract code %s", zip5, contractCode)
	}

	return internationalRateArea, nil
}
