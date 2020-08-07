package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ServicesScheduleOriginLookup does lookup on services schedule origin
type ServicesScheduleOriginLookup struct {
	MTOShipment models.MTOShipment
}

func (s ServicesScheduleOriginLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Make sure there's a pickup and destination address since those are nullable
	pickupAddressID := s.MTOShipment.PickupAddressID
	if pickupAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
	}

	var pickupAddress models.Address
	err := db.Find(&pickupAddress, s.MTOShipment.PickupAddressID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(*s.MTOShipment.PickupAddressID, "looking for PickupAddressID")
		default:
			return "", err
		}
	}

	// find the service area by querying for the service area associated with the zip3
	zip := pickupAddress.PostalCode
	zip3 := zip[0:3]

	var domesticServiceArea models.ReDomesticServiceArea
	err = db.Q().
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
