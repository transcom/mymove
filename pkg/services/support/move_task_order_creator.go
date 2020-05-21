package support

import (
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
)

type moveTaskOrderCreator struct {
	db *pop.Connection
}

func (m moveTaskOrderCreator) CreateMoveTaskOrder(moveTaskOrder supportmessages.MoveTaskOrder) (*models.MoveTaskOrder, error) {
	panic("implement me")
}

// NewInternalMoveTaskOrderCreator creates a new struct with the service dependencies
func NewInternalMoveTaskOrderCreator(db *pop.Connection) InternalMoveTaskOrderCreator {
	return &moveTaskOrderCreator{db}
}