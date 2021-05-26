package paymentrequest

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type paymentRequestStatusQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type paymentRequestStatusUpdater struct {
	builder paymentRequestStatusQueryBuilder
}

// NewPaymentRequestStatusUpdater returns a new payment request status updater
func NewPaymentRequestStatusUpdater(builder paymentRequestStatusQueryBuilder) services.PaymentRequestStatusUpdater {
	return &paymentRequestStatusUpdater{builder}
}

func (p *paymentRequestStatusUpdater) UpdateReviewedPaymentRequestStatus(paymentRequest *models.PaymentRequest, eTag string) (*models.PaymentRequest, error) {
	id := paymentRequest.ID
	status := paymentRequest.Status

	// Prevent changing status to REVIEWED if any service items are not reviewed
	if status == models.PaymentRequestStatusReviewed {
		var paymentServiceItems models.PaymentServiceItems
		serviceItemFilter := []services.QueryFilter{
			query.NewQueryFilter("payment_request_id", "=", id),
			query.NewQueryFilter("status", "=", models.PaymentServiceItemStatusRequested),
		}
		error := p.builder.FetchMany(&paymentServiceItems, serviceItemFilter, nil, nil, nil)

		if error != nil {
			return nil, error
		}

		if len(paymentServiceItems) > 0 {
			return nil, services.NewConflictError(id, "All PaymentServiceItems must be approved or denied to review this PaymentRequest")
		}
	}
	// Payment request being updated with the UpdateReviewedPaymentRequestStatus can only be:
	// REVIEWED or REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED
	if status != models.PaymentRequestStatusReviewed && status != models.PaymentRequestStatusReviewedAllRejected {
		return nil, services.NewInvalidInputError(id, nil, nil, "Payment Request status can only be updated to REVIEWED or REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED")
	}

	now := time.Now()
	paymentRequest.ReviewedAt = &now

	var verrs *validate.Errors
	var err error
	if eTag == "" {
		verrs, err = p.builder.UpdateOne(paymentRequest, nil)
	} else {
		verrs, err = p.builder.UpdateOne(paymentRequest, &eTag)
	}

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(id, err, verrs, "")
	}

	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(id, "")
		}

		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.PaymentRequest{}, services.NewPreconditionFailedError(id, err)
		}
	}

	return paymentRequest, err
}

func (p *paymentRequestStatusUpdater) UpdateProcessedPaymentRequestStatus(paymentRequest *models.PaymentRequest, eTag string) (*models.PaymentRequest, error) {
	id := paymentRequest.ID
	status := paymentRequest.Status

	// Payment request being updated with the UpdateProcessedPaymentRequestStatus can only be:
	// SENT_TO_GEX, RECEIVED_BY_GEX, EDI_ERROR and PAID
	if status != models.PaymentRequestStatusSentToGex && status != models.PaymentRequestStatusEDIError &&
		status != models.PaymentRequestStatusReceivedByGex && status != models.PaymentRequestStatusPaid {
		return nil, services.NewInvalidInputError(id, nil, nil, "Payment Request status can only be updated to SENT_TO_GEX, RECEIVED_BY_GEX, EDI_ERROR or PAID")
	}

	var recGexDate time.Time
	var sentGexDate time.Time
	var paidAtDate time.Time

	if paymentRequest.ReceivedByGexAt != nil {
		recGexDate = *paymentRequest.ReceivedByGexAt
	}
	if paymentRequest.SentToGexAt != nil {
		sentGexDate = *paymentRequest.SentToGexAt
	}
	if paymentRequest.PaidAt != nil {
		paidAtDate = *paymentRequest.PaidAt
	}

	switch status {
	case models.PaymentRequestStatusSentToGex:
		sentGexDate = time.Now()
		paymentRequest.SentToGexAt = &sentGexDate
	case models.PaymentRequestStatusReceivedByGex:
		recGexDate = time.Now()
		paymentRequest.ReceivedByGexAt = &recGexDate
	case models.PaymentRequestStatusPaid:
		paidAtDate = time.Now()
		paymentRequest.PaidAt = &paidAtDate
	}

	var verrs *validate.Errors
	var err error
	if eTag == "" {
		verrs, err = p.builder.UpdateOne(paymentRequest, nil)
	} else {
		verrs, err = p.builder.UpdateOne(paymentRequest, &eTag)
	}

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(id, err, verrs, "")
	}

	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(id, "")
		}

		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.PaymentRequest{}, services.NewPreconditionFailedError(id, err)
		}
	}

	return paymentRequest, err
}
