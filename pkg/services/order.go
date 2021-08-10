package services

import (
	"io"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/storage"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

// OrderFetcher is the service object interface for FetchOrder
//go:generate mockery --name OrderFetcher --disable-version-string
type OrderFetcher interface {
	FetchOrder(appCfg appconfig.AppConfig, moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListOrders(appCfg appconfig.AppConfig, officeUserID uuid.UUID, params *ListOrderParams) ([]models.Move, int, error)
}

//OrderUpdater is the service object interface for updating fields of an Order
//go:generate mockery --name OrderUpdater --disable-version-string
type OrderUpdater interface {
	UploadAmendedOrdersAsCustomer(appCfg appconfig.AppConfig, userID uuid.UUID, orderID uuid.UUID, file io.ReadCloser, filename string, storer storage.FileStorer) (models.Upload, string, *validate.Errors, error)
	UpdateOrderAsTOO(appCfg appconfig.AppConfig, orderID uuid.UUID, payload ghcmessages.UpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateOrderAsCounselor(appCfg appconfig.AppConfig, orderID uuid.UUID, payload ghcmessages.CounselingUpdateOrderPayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateAllowanceAsTOO(appCfg appconfig.AppConfig, orderID uuid.UUID, payload ghcmessages.UpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error)
	UpdateAllowanceAsCounselor(appCfg appconfig.AppConfig, orderID uuid.UUID, payload ghcmessages.CounselingUpdateAllowancePayload, eTag string) (*models.Order, uuid.UUID, error)
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
