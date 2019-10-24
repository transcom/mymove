package movetaskorder

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// claimAccessCode is a service object to validate an access code.
type fetchMoveTaskOrder struct {
	db *pop.Connection
}

// NewAccessCodeClaimer creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &fetchMoveTaskOrder{db}
}

func (f fetchMoveTaskOrder) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	panic("implement me")
}
