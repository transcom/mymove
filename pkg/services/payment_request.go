package services

import (
	"io"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PaymentRequestCreator is the exported interface for creating a payment request
//go:generate mockery --name PaymentRequestCreator --disable-version-string
type PaymentRequestCreator interface {
	CreatePaymentRequest(appCtx appcontext.AppContext, paymentRequest *models.PaymentRequest) (*models.PaymentRequest, error)
}

// PaymentRequestRecalculator is the exported interface for recalculating a payment request
//go:generate mockery --name PaymentRequestRecalculator --disable-version-string
type PaymentRequestRecalculator interface {
	RecalculatePaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (*models.PaymentRequest, error)
}

// PaymentRequestShipmentRecalculator is the exported interface for recalculating payment requests for a shipment
//go:generate mockery --name PaymentRequestShipmentRecalculator --disable-version-string
type PaymentRequestShipmentRecalculator interface {
	ShipmentRecalculatePaymentRequest(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.PaymentRequests, error)
}

// PaymentRequestListFetcher is the exported interface for fetching a list of payment requests
//go:generate mockery --name PaymentRequestListFetcher --disable-version-string
type PaymentRequestListFetcher interface {
	FetchPaymentRequestList(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *FetchPaymentRequestListParams) (*models.PaymentRequests, int, error)
	FetchPaymentRequestListByMove(appCtx appcontext.AppContext, officeUserID uuid.UUID, locator string) (*models.PaymentRequests, error)
}

// PaymentRequestFetcher is the exported interface for fetching a payment request
//go:generate mockery --name PaymentRequestFetcher --disable-version-string
type PaymentRequestFetcher interface {
	FetchPaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (models.PaymentRequest, error)
}

// PaymentRequestReviewedFetcher is the exported interface for fetching all payment requests in 'reviewed' status
//go:generate mockery --name PaymentRequestReviewedFetcher --disable-version-string
type PaymentRequestReviewedFetcher interface {
	FetchReviewedPaymentRequest(appCtx appcontext.AppContext) (models.PaymentRequests, error)
}

// PaymentRequestStatusUpdater is the exported interface for updating the status of a payment request
//go:generate mockery --name PaymentRequestStatusUpdater --disable-version-string
type PaymentRequestStatusUpdater interface {
	UpdatePaymentRequestStatus(appCtx appcontext.AppContext, paymentRequest *models.PaymentRequest, eTag string) (*models.PaymentRequest, error)
}

// PaymentRequestUploadCreator is the exported interface for creating a payment request upload
//go:generate mockery --name PaymentRequestUploadCreator --disable-version-string
type PaymentRequestUploadCreator interface {
	CreateUpload(appCtx appcontext.AppContext, file io.ReadCloser, paymentRequestID uuid.UUID, userID uuid.UUID, filename string) (*models.Upload, error)
}

// PaymentRequestReviewedProcessor is the exported interface for processing reviewed payment requests
//go:generate mockery --name PaymentRequestReviewedProcessor --disable-version-string
type PaymentRequestReviewedProcessor interface {
	ProcessReviewedPaymentRequest(appCtx appcontext.AppContext)
	ProcessAndLockReviewedPR(appCtx appcontext.AppContext, pr models.PaymentRequest) error
}

// FetchPaymentRequestListParams is a public struct that's used to pass filter arguments to FetchPaymentRequestList
type FetchPaymentRequestListParams struct {
	Branch                  *string
	Locator                 *string
	DodID                   *string
	LastName                *string
	DestinationDutyLocation *string
	Status                  []string
	Page                    *int64
	PerPage                 *int64
	SubmittedAt             *time.Time
	Sort                    *string
	Order                   *string
	OriginDutyLocation      *string
}

// ShipmentPaymentSITBalance is a public struct that's used to return current SIT balances to the TIO for a payment
// request
type ShipmentPaymentSITBalance struct {
	ShipmentID              uuid.UUID
	TotalSITDaysAuthorized  int
	TotalSITDaysRemaining   int
	TotalSITEndDate         time.Time
	PendingSITDaysInvoiced  int
	PendingBilledEndDate    time.Time
	PreviouslyBilledDays    *int
	PreviouslyBilledEndDate *time.Time
}

// ShipmentsPaymentSITBalance is the exported interface for returning SIT balances for all shipments of a payment
// request
//go:generate mockery --name ShipmentsPaymentSITBalance --disable-version-string
type ShipmentsPaymentSITBalance interface {
	ListShipmentPaymentSITBalance(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) ([]ShipmentPaymentSITBalance, error)
}
