package serviceparamvaluelookups

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// MTOAvailableToPrimeAtLookup does lookup on the MTOAvailableToPrime timestamp
type MTOAvailableToPrimeAtLookup struct {
}

func (m MTOAvailableToPrimeAtLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	// Get the MoveTaskOrder
	moveTaskOrderID := keyData.MoveTaskOrderID
	var moveTaskOrder models.MoveTaskOrder
	err := db.Find(&moveTaskOrder, moveTaskOrderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", services.NewNotFoundError(moveTaskOrderID, "looking for MoveTaskOrderID")
		default:
			return "", err
		}
	}

	availableToPrimeAt := moveTaskOrder.AvailableToPrimeAt
	if availableToPrimeAt == nil {
		return "", services.NewBadDataError("This move task order is not available to prime")
	}

	return (*availableToPrimeAt).Format(TimestampParamFormat), nil
}
