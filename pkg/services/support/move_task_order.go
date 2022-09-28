package support

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
)

// InternalMoveTaskOrderCreator is the service object interface for InternalCreateMoveTaskOrder
//
//go:generate mockery --name InternalMoveTaskOrderCreator --disable-version-string
type InternalMoveTaskOrderCreator interface {
	InternalCreateMoveTaskOrder(appCtx appcontext.AppContext, moveTaskOrder supportmessages.MoveTaskOrder) (*models.Move, error)
}
