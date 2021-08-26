package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ServicesScheduleLookup does lookup on services schedule origin
type ServicesScheduleLookup struct {
	Address models.Address
}

func (s ServicesScheduleLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	// find the service area by querying for the service area associated with the zip3
	zip := s.Address.PostalCode
	zip3 := zip[0:3]

	domesticServiceArea, err := fetchDomesticServiceArea(appCtx, keyData.ContractCode, zip3)

	return strconv.Itoa(domesticServiceArea.ServicesSchedule), err
}
