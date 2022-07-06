package services

import (
	"io"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MoveListFetcher is the exported interface for fetching multiple moves
//go:generate mockery --name MoveListFetcher --disable-version-string
type MoveListFetcher interface {
	FetchMoveList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Moves, error)
	FetchMoveCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// MoveFetcher is the exported interface for fetching a move by locator
//go:generate mockery --name MoveFetcher --disable-version-string
type MoveFetcher interface {
	FetchMove(appCtx appcontext.AppContext, locator string, searchParams *MoveFetcherParams) (*models.Move, error)
}

//go:generate mockery --name MoveSearcher --disable-version-string
type MoveSearcher interface {
	SearchMoves(appCtx appcontext.AppContext, params *SearchMovesParams) (models.Moves, int, error)
}

// MoveFetcherParams is  public struct that's used to pass filter arguments to
// MoveFetcher queries
type MoveFetcherParams struct {
	IncludeHidden bool // indicates if a hidden/disabled move can be returned
}

// MoveRouter is the exported interface for routing moves at different stages
//go:generate mockery --name MoveRouter --disable-version-string
type MoveRouter interface {
	Approve(appCtx appcontext.AppContext, move *models.Move) error
	ApproveOrRequestApproval(appCtx appcontext.AppContext, move models.Move) (*models.Move, error)
	Cancel(appCtx appcontext.AppContext, reason string, move *models.Move) error
	CompleteServiceCounseling(appCtx appcontext.AppContext, move *models.Move) error
	SendToOfficeUser(appCtx appcontext.AppContext, move *models.Move) error
	Submit(appCtx appcontext.AppContext, move *models.Move) error
}

// MoveWeights is the exported interface for flagging a move with an excess weight risk
//go:generate mockery --name MoveWeights --disable-version-string
type MoveWeights interface {
	CheckExcessWeight(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment models.MTOShipment) (*models.Move, *validate.Errors, error)
	CheckAutoReweigh(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment *models.MTOShipment) (models.MTOShipments, error)
}

// MoveExcessWeightUploader is the exported interface for uploading an excess weight document for a move
//go:generate mockery --name MoveExcessWeightUploader --disable-version-string
type MoveExcessWeightUploader interface {
	CreateExcessWeightUpload(
		appCtx appcontext.AppContext,
		moveID uuid.UUID,
		file io.ReadCloser,
		uploadFilename string,
		uploadType models.UploadType,
	) (*models.Move, error)
}

// MoveFinancialReviewFlagSetter is the exported interface for flagging a move for financial review
//go:generate mockery --name MoveFinancialReviewFlagSetter --disable-version-string
type MoveFinancialReviewFlagSetter interface {
	SetFinancialReviewFlag(appCtx appcontext.AppContext, moveID uuid.UUID, eTag string, flagForReview bool, remarks *string) (*models.Move, error)
}

type SearchMovesParams struct {
	Branch                *string
	Locator               *string
	DodID                 *string
	CustomerName          *string
	DestinationPostalCode *string
	OriginPostalCode      *string
	Status                []string
	ShipmentsCount        *int64
	Page                  int64
	PerPage               int64
	Sort                  *string
	Order                 *string
}
