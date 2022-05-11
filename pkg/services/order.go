package services

import (
	"io"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// OrderFetcher is the service object interface for FetchOrder
//go:generate mockery --name OrderFetcher --disable-version-string
type OrderFetcher interface {
	FetchOrder(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListOrders(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *ListOrderParams) ([]models.Move, int, error)
}

//OrderUpdater is the service object interface for updating fields of an Order
//go:generate mockery --name OrderUpdater --disable-version-string
type OrderUpdater interface {
	UploadAmendedOrdersAsCustomer(appCtx appcontext.AppContext, userID uuid.UUID, orderID uuid.UUID, file io.ReadCloser, filename string, storer storage.FileStorer) (models.Upload, string, *validate.Errors, error)
	UpdateOrderAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.UpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateOrderAsCounselor(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.CounselingUpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateAllowanceAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.UpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateAllowanceAsCounselor(appCtx appcontext.AppContext, orderID uuid.UUID, payload ghcmessages.CounselingUpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error)
}

//ExcessWeightRiskManager is the service object interface for updating the max billable weight for an Order's Entitlement
//go:generate mockery --name ExcessWeightRiskManager --disable-version-string
type ExcessWeightRiskManager interface {
	AcknowledgeExcessWeightRisk(appCtx appcontext.AppContext, moveID uuid.UUID, eTag string) (*models.Move, error)
	UpdateBillableWeightAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, weight *int, eTag string) (*models.Order, uuid.UUID, error)
	UpdateMaxBillableWeightAsTIO(appCtx appcontext.AppContext, orderID uuid.UUID, weight *int, remarks *string, eTag string) (*models.Order, uuid.UUID, error)
}

// ListOrderParams is a public struct that's used to pass filter arguments to the ListOrders
type ListOrderParams struct {
	Branch                  *string
	Locator                 *string
	DodID                   *string
	LastName                *string
	DestinationDutyLocation *string
	OriginDutyLocation      *string
	OriginGBLOC             *string
	SubmittedAt             *time.Time
	RequestedMoveDate       *string
	Status                  []string
	Page                    *int64
	PerPage                 *int64
	Sort                    *string
	Order                   *string
}
