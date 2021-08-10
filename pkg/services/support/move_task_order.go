package support

import (
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
)

// InternalMoveTaskOrderCreator is the service object interface for InternalCreateMoveTaskOrder
//go:generate mockery --name InternalMoveTaskOrderCreator --disable-version-string
type InternalMoveTaskOrderCreator interface {
	InternalCreateMoveTaskOrder(appCfg appconfig.AppConfig, moveTaskOrder supportmessages.MoveTaskOrder, logger *zap.Logger) (*models.Move, error)
}
