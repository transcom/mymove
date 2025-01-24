package paymentrequest

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestFetcher struct {
}

// NewPaymentRequestFetcher returns a new payment request fetcher
func NewPaymentRequestFetcher() services.PaymentRequestFetcher {
	return &paymentRequestFetcher{}
}

// FetchPaymentRequest finds the payment request by id
func (p *paymentRequestFetcher) FetchPaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (models.PaymentRequest, error) {
	var paymentRequest models.PaymentRequest

	// fetch the payment request first with proof of service docs
	// will error if payment request not found
	err := appCtx.DB().Eager(
		"PaymentServiceItems.MTOServiceItem.MTOShipment",
		"PaymentServiceItems.MTOServiceItem.ReService",
		"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"ProofOfServiceDocs").
		Find(&paymentRequest, paymentRequestID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.PaymentRequest{}, apperror.NewNotFoundError(paymentRequestID, "looking for PaymentRequest")
		default:
			return models.PaymentRequest{}, apperror.NewQueryError("PaymentRequest", err, "")
		}
	}

	// then fetch the uploads separately to omit soft deleted items
	// empty records are expected
	for index, posd := range paymentRequest.ProofOfServiceDocs {
		var primeUploads models.PrimeUploads
		err = appCtx.DB().Q().
			Where("prime_uploads.proof_of_service_docs_id = ? AND prime_uploads.deleted_at IS NULL AND u.deleted_at IS NULL", posd.ID).
			Eager("Upload").
			Join("uploads as u", "u.id = prime_uploads.upload_id").
			All(&primeUploads)
		if err != nil {
			return models.PaymentRequest{}, err
		}

		paymentRequest.ProofOfServiceDocs[index].PrimeUploads = primeUploads
	}

	return paymentRequest, err
}

type paymentRequestFetcherBulkAssignment struct {
}

// NewPaymentRequestFetcherBulkAssignment creates a new paymentRequestFetcherBulkAssignment service
func NewPaymentRequestFetcherBulkAssignment() services.PaymentRequestFetcherBulkAssignment {
	return &paymentRequestFetcherBulkAssignment{}
}

func (f paymentRequestFetcherBulkAssignment) FetchPaymentRequestsForBulkAssignment(appCtx appcontext.AppContext, gbloc string) ([]models.PaymentRequestWithEarliestRequestedDate, error) {
	var payment_requests []models.PaymentRequestWithEarliestRequestedDate

	sqlQuery := `
		SELECT
			payment_requests.id,
			payment_requests.requested_at
		FROM payment_requests
		INNER JOIN moves on moves.id = payment_requests.move_id
		INNER JOIN orders ON orders.id = moves.orders_id
		INNER JOIN service_members ON orders.service_member_id = service_members.id
		LEFT JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
		WHERE payment_requests.status = 'PENDING'
		AND moves.show = $1
		AND (orders.orders_type NOT IN ($2, $3, $4))
		AND moves.tio_assigned_id IS NULL `
	if gbloc == "USMC" {
		sqlQuery += `
			AND service_members.affiliation ILIKE 'MARINES' `
	} else {
		sqlQuery += `
		AND service_members.affiliation != 'MARINES'
		AND move_to_gbloc.gbloc = '` + gbloc + `' `
	}
	sqlQuery += `
		GROUP BY payment_requests.id
        ORDER BY payment_requests.requested_at ASC`

	err := appCtx.DB().RawQuery(sqlQuery,
		models.BoolPointer(true),
		internalmessages.OrdersTypeBLUEBARK,
		internalmessages.OrdersTypeWOUNDEDWARRIOR,
		internalmessages.OrdersTypeSAFETY).
		All(&payment_requests)

	if err != nil {
		return nil, fmt.Errorf("error fetching payment requests for GBLOC: %s with error %w", gbloc, err)
	}

	if len(payment_requests) < 1 {
		return nil, nil
	}

	return payment_requests, nil
}
