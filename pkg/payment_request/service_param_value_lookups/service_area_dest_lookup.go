package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ServiceAreaDestLookup does lookup on destination address postal code
type ServiceAreaDestLookup struct {
}

func (r ServiceAreaDestLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {

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
		//check for string of all zeros
		return "", fmt.Errorf("could not find destination address for MTOShipment [%s]", mtoShipmentID)
	}

	zip := mtoServiceItem.MTOShipment.DestinationAddress.PostalCode
	zip3 := zip[0:3]

	var domesticServiceArea models.ReDomesticServiceArea

	query := db.Q().Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Where("zip3 = ?", zip3)

	err = query.First(&domesticServiceArea)

	return domesticServiceArea.ServiceArea, err
}
