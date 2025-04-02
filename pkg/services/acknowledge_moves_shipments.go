package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//go:generate mockery --name MoveAndShipmentAcknowledgementUpdater
type MoveAndShipmentAcknowledgementUpdater interface {
	AcknowledgeMovesAndShipments(appCtx appcontext.AppContext, moves *models.Moves) error
}
