package route

import (
	"github.com/transcom/mymove/pkg/models"
)

// Planner is the interface needed by Handlers to be able to evaluate the distance to be used for move accounting
type Planner interface {
	TransitDistance(source *models.Address, destination *models.Address) (int, error)
}
