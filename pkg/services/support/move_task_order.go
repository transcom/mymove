package support

import (
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// InternalMoveTaskOrderCreator is the service object interface for InternalCreateMoveTaskOrder
//go:generate mockery -name InternalMoveTaskOrderCreator
type InternalMoveTaskOrderCreator interface {
	InternalCreateMoveTaskOrder(moveTaskOrder supportmessages.MoveTaskOrder, logger handlers.Logger) (*models.Move, error)
}
