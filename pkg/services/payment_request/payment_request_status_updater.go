package paymentrequest

import (
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

func (p *paymentRequestStatusUpdater) UpdatePaymentRequestStatus(paymentRequest *models.PaymentRequest, eTag string) (*models.PaymentRequest, error) {
	id := paymentRequest.ID
	status := paymentRequest.Status

	// Prevent changing status to REVIEWED if any service items are not reviewed
	if status == models.PaymentRequestStatusReviewed {
		var paymentServiceItems models.PaymentServiceItems
		serviceItemFilter := []services.QueryFilter{
			query.NewQueryFilter("payment_request_id", "=", id),
			query.NewQueryFilter("status", "=", models.PaymentServiceItemStatusRequested),
		}
		associations := query.NewQueryAssociations([]services.QueryAssociation{})
		error := p.builder.FetchMany(&paymentServiceItems, serviceItemFilter, associations, nil, nil)

		if error != nil {
			return nil, error
		}

		if len(paymentServiceItems) > 0 {
			return nil, services.NewInvalidInputError(id, nil, nil, "All PaymentServiceItems must be approved or denied to review this PaymentRequest")
		}
	}

	verrs, err := p.builder.UpdateOne(paymentRequest, &eTag)

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
