package movetaskorder

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderChecker struct {
}

// NewMoveTaskOrderChecker creates a new struct with the service dependencies
func NewMoveTaskOrderChecker() services.MoveTaskOrderChecker {
	return &moveTaskOrderChecker{}
}

//MTOAvailableToPrime retrieves a Move for a given UUID and checks if it is visible and available to prime
func (f moveTaskOrderChecker) MTOAvailableToPrime(appCfg appconfig.AppConfig, moveTaskOrderID uuid.UUID) (bool, error) {
	mto := &models.Move{}
	err := appCfg.DB().RawQuery("SELECT * FROM moves WHERE id = $1 AND show = TRUE", moveTaskOrderID).First(mto)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, services.NewNotFoundError(moveTaskOrderID, "for moveTaskOrder")
		default:
			return false, err
		}
	}

	if mto.AvailableToPrimeAt == nil {
		return false, nil
	}

	return true, nil
}
