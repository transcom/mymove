package services

import (
	"io"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

// OrderFetcher is the service object interface for FetchOrder
//go:generate mockery --name OrderFetcher --disable-version-string
type OrderFetcher interface {
	FetchOrder(moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListOrders(officeUserID uuid.UUID, params *ListOrderParams) ([]models.Move, int, error)
}

//OrderUpdater is the service object interface for updating fields of an Order
//go:generate mockery --name OrderUpdater --disable-version-string
type OrderUpdater interface {
	UploadAmendedOrdersAsCustomer(logger *zap.Logger, userID uuid.UUID, orderID uuid.UUID, file io.ReadCloser, filename string, storer storage.FileStorer) (models.Upload, string, *validate.Errors, error)
	UpdateOrderAsTOO(orderID uuid.UUID, payload ghcmessages.UpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateOrderAsCounselor(orderID uuid.UUID, payload ghcmessages.CounselingUpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateAllowanceAsTOO(orderID uuid.UUID, payload ghcmessages.UpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateAllowanceAsCounselor(orderID uuid.UUID, payload ghcmessages.CounselingUpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error)
}

// ListOrderParams is a public struct that's used to pass filter arguments to the ListOrders
type ListOrderParams struct {
	Branch                 *string
	Locator                *string
	DodID                  *string
	LastName               *string
	DestinationDutyStation *string
	OriginGBLOC            *string
	SubmittedAt            *time.Time
	RequestedMoveDate      *string
	Status                 []string
	Page                   *int64
	PerPage                *int64
	Sort                   *string
	Order                  *string
}
