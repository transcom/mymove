package services

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// OrderFetcher is the service object interface for FetchOrder
//go:generate mockery --name OrderFetcher
type OrderFetcher interface {
	FetchOrder(moveTaskOrderID uuid.UUID) (*models.Order, error)
	ListOrders(officeUserID uuid.UUID, params *ListOrderParams) ([]models.Move, int, error)
}

//OrderUpdater is the service object interface for updating fields of an Order
//go:generate mockery --name OrderUpdater
type OrderUpdater interface {
	UploadAmendedOrders(orderID uuid.UUID, payload *internalmessages.UserUploadPayload, eTag string) (*models.Order, error)
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
