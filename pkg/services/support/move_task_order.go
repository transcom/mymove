package support

import (
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
)

// InternalMoveTaskOrderCreator is the service object interface for SupportCreateMoveTaskOrder
//go:generate mockery -name SupportCreateMoveTaskOrder
type InternalMoveTaskOrderCreator interface {
	CreateMoveTaskOrder(moveTaskOrder supportmessages.MoveTaskOrder) (*models.MoveTaskOrder, error)
}