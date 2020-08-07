package serviceparamvaluelookups

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ServiceAreaDestLookup does lookup on destination address postal code
type ServiceAreaDestLookup struct {
	MTOShipment models.MTOShipment
}

func (r ServiceAreaDestLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {

	db := *keyData.db

	// Make sure there's a destination address since those are nullable
	destinationAddressID := r.MTOShipment.DestinationAddressID
	if destinationAddressID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
	}

	var destinationAddress models.Address
	err := db.Find(&destinationAddress, r.MTOShipment.DestinationAddressID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(*r.MTOShipment.DestinationAddressID, "looking for DestinationAddressID")
		default:
			return "", err
		}
	}

	zip := destinationAddress.PostalCode
	zip3 := zip[0:3]

	var domesticServiceArea models.ReDomesticServiceArea

	query := db.Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("zip3 = ?", zip3).
		Where("re_contracts.code = ?", keyData.ContractCode)

	err = query.First(&domesticServiceArea)

	return domesticServiceArea.ServiceArea, err
}
