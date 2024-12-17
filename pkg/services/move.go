package services

import (
	"io"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/storage"
)

// MoveListFetcher is the exported interface for fetching multiple moves
//
//go:generate mockery --name MoveListFetcher
type MoveListFetcher interface {
	FetchMoveList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.Moves, error)
	FetchMoveCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// MoveFetcher is the exported interface for fetching a move by locator
//
//go:generate mockery --name MoveFetcher
type MoveFetcher interface {
	FetchMove(appCtx appcontext.AppContext, locator string, searchParams *MoveFetcherParams) (*models.Move, error)
	FetchMovesForPPTASReports(appCtx appcontext.AppContext, params *MoveTaskOrderFetcherParams) (models.Moves, error)
}

//go:generate mockery --name MoveSearcher
type MoveSearcher interface {
	SearchMoves(appCtx appcontext.AppContext, params *SearchMovesParams) (models.Moves, int, error)
}

// MoveFetcherParams is  public struct that's used to pass filter arguments to
// MoveFetcher queries
type MoveFetcherParams struct {
	IncludeHidden bool // indicates if a hidden/disabled move can be returned
}

// MoveRouter is the exported interface for routing moves at different stages
//
//go:generate mockery --name MoveRouter
type MoveRouter interface {
	Approve(appCtx appcontext.AppContext, move *models.Move) error
	ApproveOrRequestApproval(appCtx appcontext.AppContext, move models.Move) (*models.Move, error)
	Cancel(appCtx appcontext.AppContext, move *models.Move) error
	CompleteServiceCounseling(appCtx appcontext.AppContext, move *models.Move) error
	RouteAfterAmendingOrders(appCtx appcontext.AppContext, move *models.Move) error
	SendToOfficeUser(appCtx appcontext.AppContext, move *models.Move) error
	Submit(appCtx appcontext.AppContext, move *models.Move, newSignedCertification *models.SignedCertification) error
}

// MoveWeights is the exported interface for flagging a move with an excess weight risk
//
//go:generate mockery --name MoveWeights
type MoveWeights interface {
	CheckExcessWeight(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment models.MTOShipment) (*models.Move, *validate.Errors, error)
	CheckAutoReweigh(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment *models.MTOShipment) (models.MTOShipments, error)
	MoveShouldAutoReweigh(appCtx appcontext.AppContext, moveID uuid.UUID) (*bool, error)
}

// MoveExcessWeightUploader is the exported interface for uploading an excess weight document for a move
//
//go:generate mockery --name MoveExcessWeightUploader
type MoveExcessWeightUploader interface {
	CreateExcessWeightUpload(
		appCtx appcontext.AppContext,
		moveID uuid.UUID,
		file io.ReadCloser,
		uploadFilename string,
		uploadType models.UploadType,
	) (*models.Move, error)
}

type MoveAdditionalDocumentsUploader interface {
	CreateAdditionalDocumentsUpload(
		appCtx appcontext.AppContext,
		userID uuid.UUID,
		moveID uuid.UUID,
		file io.ReadCloser,
		uploadFilename string,
		storer storage.FileStorer,
		uploadType models.UploadType,
	) (models.Upload, string, *validate.Errors, error)
}

// MoveFinancialReviewFlagSetter is the exported interface for flagging a move for financial review
//
//go:generate mockery --name MoveFinancialReviewFlagSetter
type MoveFinancialReviewFlagSetter interface {
	SetFinancialReviewFlag(appCtx appcontext.AppContext, moveID uuid.UUID, eTag string, flagForReview bool, remarks *string) (*models.Move, error)
}

type SearchMovesParams struct {
	Branch                *string
	Locator               *string
	DodID                 *string
	Emplid                *string
	CustomerName          *string
	PaymentRequestCode    *string
	DestinationPostalCode *string
	OriginPostalCode      *string
	Status                []string
	ShipmentsCount        *int64
	Page                  int64
	PerPage               int64
	Sort                  *string
	Order                 *string
	PickupDate            *time.Time
	DeliveryDate          *time.Time
	MoveCreatedDate       *time.Time
}

type MoveCloseoutOfficeUpdater interface {
	UpdateCloseoutOffice(appCtx appcontext.AppContext, moveLocator string, closeoutOfficeID uuid.UUID, eTag string) (*models.Move, error)
}

type MoveCanceler interface {
	CancelMove(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Move, error)
}

type MoveAssignedOfficeUserUpdater interface {
	UpdateAssignedOfficeUser(appCtx appcontext.AppContext, moveID uuid.UUID, officeUser *models.OfficeUser, role roles.RoleType) (*models.Move, error)
	DeleteAssignedOfficeUser(appCtx appcontext.AppContext, moveID uuid.UUID, role roles.RoleType) (*models.Move, error)
}

type CheckForLockedMovesAndUnlockHandler interface {
	CheckForLockedMovesAndUnlock(appCtx appcontext.AppContext, officeUserID uuid.UUID) error
}
