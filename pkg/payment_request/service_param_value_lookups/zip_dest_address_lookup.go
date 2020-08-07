package serviceparamvaluelookups

import (
	"fmt"

	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ZipDestAddressLookup does lookup on actual weight billed
type ZipDestAddressLookup struct {
	MTOShipment models.MTOShipment
}

func (r ZipDestAddressLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
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

	value := fmt.Sprintf("%+v", destinationAddress.PostalCode)
	return value, nil
}
