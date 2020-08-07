package serviceparamvaluelookups

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ServiceAreaOriginLookup does lookup on pickup address postal code
type ServiceAreaOriginLookup struct {
	MTOShipment models.MTOShipment
}

func (r ServiceAreaOriginLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Make sure there's a pickup and destination address since those are nullable
	pickupAddressID := r.MTOShipment.PickupAddressID
	if pickupAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
	}

	var pickupAddress models.Address
	err := db.Find(&pickupAddress, r.MTOShipment.PickupAddressID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(*r.MTOShipment.PickupAddressID, "looking for PickupAddressID")
		default:
			return "", err
		}
	}

	// find the service area by querying for the service area associated with the zip3
	zip := pickupAddress.PostalCode
	zip3 := zip[0:3]

	var domesticServiceArea models.ReDomesticServiceArea

	query := db.Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("re_zip3s.zip3 = ?", zip3).
		Where("re_contracts.code = ?", keyData.ContractCode)

	err = query.First(&domesticServiceArea)

	return domesticServiceArea.ServiceArea, err

}
