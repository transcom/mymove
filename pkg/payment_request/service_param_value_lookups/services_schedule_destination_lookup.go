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
}

func (s ServicesScheduleDestLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService", "MTOShipment.DestinationAddress").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// Make sure there's an MTOShipment since that's nullable
	mtoShipmentID := mtoServiceItem.MTOShipmentID
	if mtoShipmentID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	// Make sure there's a destination address since those are nullable
	destinationAddressID := mtoServiceItem.MTOShipment.DestinationAddressID
	if destinationAddressID == nil || *destinationAddressID == uuid.Nil {
		return "", fmt.Errorf("could not find destination address for MTOShipment [%s]", mtoShipmentID)
	}

	// find the service area by querying for the service area associated with the zip3
	zip := mtoServiceItem.MTOShipment.DestinationAddress.PostalCode
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
