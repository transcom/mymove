package move

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveStatusRouter struct {
	db *pop.Connection
}

// NewMoveStatusRouter creates a new moveStatusRouter service
func NewMoveStatusRouter(db *pop.Connection) services.MoveStatusRouter {
	return &moveStatusRouter{db}
}

//FetchOrder retrieves a Move if it is visible for a given locator
func (f moveStatusRouter) RouteMove(move *models.Move) error {
	// TODO: In future, add logic based on the service member's origin duty station
	// to route to send services counseling, otherwise submitted
	// submitDate := time.Now()
	// err := move.Submit(submitDate)
	err := move.SendToServiceCounseling()
	if err != nil {
		return err
	}
	return nil
}
