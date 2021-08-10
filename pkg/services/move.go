package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// MoveListFetcher is the exported interface for fetching multiple moves
//go:generate mockery --name MoveListFetcher --disable-version-string
type MoveListFetcher interface {
	FetchMoveList(appCfg appconfig.AppConfig, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Moves, error)
	FetchMoveCount(appCfg appconfig.AppConfig, filters []QueryFilter) (int, error)
}

// MoveFetcher is the exported interface for fetching a move by locator
//go:generate mockery --name MoveFetcher --disable-version-string
type MoveFetcher interface {
	FetchMove(appCfg appconfig.AppConfig, locator string, searchParams *MoveFetcherParams) (*models.Move, error)
}

// MoveFetcherParams is  public struct that's used to pass filter arguments to
// MoveFetcher queries
type MoveFetcherParams struct {
	IncludeHidden bool // indicates if a hidden/disabled move can be returned
}

// MoveRouter is the exported interface for routing moves at different stages
//go:generate mockery --name MoveRouter --disable-version-string
type MoveRouter interface {
	Approve(appCfg appconfig.AppConfig, move *models.Move) error
	ApproveAmendedOrders(appCfg appconfig.AppConfig, moveID uuid.UUID, orderID uuid.UUID) (models.Move, error)
	Cancel(appCfg appconfig.AppConfig, reason string, move *models.Move) error
	CompleteServiceCounseling(appCfg appconfig.AppConfig, move *models.Move) error
	SendToOfficeUser(appCfg appconfig.AppConfig, move *models.Move) error
	Submit(appCfg appconfig.AppConfig, move *models.Move) error
}

// MoveWeights is the exported interface for flagging a move with an excess weight risk
//go:generate mockery --name MoveWeights --disable-version-string
type MoveWeights interface {
	CheckExcessWeight(appCfg appconfig.AppConfig, moveID uuid.UUID, updatedShipment models.MTOShipment) (*models.Move, *validate.Errors, error)
}
