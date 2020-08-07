package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// ServicesScheduleDestLookup does lookup on services schedule destination
type ServicesScheduleDestLookup struct {
	MTOShipment models.MTOShipment
}

func (s ServicesScheduleDestLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Make sure there's a pickup and destination address since those are nullable
	destinationAddressID := s.MTOShipment.DestinationAddressID
	if destinationAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
	}

	var destinationAddress models.Address
	err := db.Find(&destinationAddress, s.MTOShipment.DestinationAddressID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(*s.MTOShipment.DestinationAddressID, "looking for DestinationAddressID")
		default:
			return "", err
		}
	}

	// find the service area by querying for the service area associated with the zip3
	zip := destinationAddress.PostalCode
	zip3 := zip[0:3]

	var domesticServiceArea models.ReDomesticServiceArea
	err = db.Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("re_zip3s.zip3 = ?", zip3).
		Where("re_contracts.code = ?", ghcrateengine.DefaultContractCode).
		First(&domesticServiceArea)
	if err != nil {
		return "", fmt.Errorf("unable to find domestic service area for %s under contract code %s", zip3, ghcrateengine.DefaultContractCode)
	}

	return strconv.Itoa(domesticServiceArea.ServicesSchedule), nil
}
