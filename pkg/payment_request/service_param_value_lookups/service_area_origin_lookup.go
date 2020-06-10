package serviceparamvaluelookups

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ServiceAreaOrigin does lookup on actual weight billed
type ServiceAreaOriginLookup struct {
}

func (r ServiceAreaOriginLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {

	db := *keyData.db

	// Get the MTOServiceItem and associated MTOShipment
	mtoServiceItemID := keyData.MTOServiceItemID
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService", "MTOShipment").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return "", err
		}
	}

	// Make sure there's an MTOShipment since that's nullable
	mtoShipmentID := mtoServiceItem.MTOShipment.PickupAddress.
	if mtoShipmentID == nil {
		return "", services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
	}

	// Make sure there's an actual weight since that's nullable
	zip3 := mtoServiceItem.MTOShipment.PickupAddress.PostalCode[0:3]

	query := `
	SELECT service_area from re_domestic_service_areas
	JOIN re_zip3s on re_zip3s.domestic_service_area_id = re_domestic_service_areas.id
	`

	q := db.RawQuery(query).Where("zip3 = ?", zip3)
	// Select zip3, service_area from re_zip3s
	// JOIN re_domestic_service_areas rdsa on re_zip3s.domestic_service_area_id = rdsa.id

	return q
}