package paymentrequest

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type paymentRequestStatusQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type paymentRequestStatusUpdater struct {
	builder paymentRequestStatusQueryBuilder
}

// NewPaymentRequestStatusUpdater returns a new payment request status updater
func NewPaymentRequestStatusUpdater(builder paymentRequestStatusQueryBuilder) services.PaymentRequestStatusUpdater {
	return &paymentRequestStatusUpdater{builder}
}

func (p *paymentRequestStatusUpdater) UpdatePaymentRequestStatus(appCtx appcontext.AppContext, paymentRequest *models.PaymentRequest, eTag string) (*models.PaymentRequest, error) {
	id := paymentRequest.ID
	status := paymentRequest.Status

	// Prevent changing status to REVIEWED if any service items are not reviewed
	if status == models.PaymentRequestStatusReviewed {
		var paymentServiceItems models.PaymentServiceItems
		serviceItemFilter := []services.QueryFilter{
			query.NewQueryFilter("payment_request_id", "=", id),
			query.NewQueryFilter("status", "=", models.PaymentServiceItemStatusRequested),
		}
		error := p.builder.FetchMany(appCtx, &paymentServiceItems, serviceItemFilter, nil, nil, nil)

		if error != nil {
			return nil, error
		}

		if len(paymentServiceItems) > 0 {
			return nil, apperror.NewConflictError(id, "All PaymentServiceItems must be approved or denied to review this PaymentRequest")
		}
	}

	var verrs *validate.Errors
	var err error
	if eTag == "" {
		verrs, err = p.builder.UpdateOne(appCtx, paymentRequest, nil)
	} else {
		verrs, err = p.builder.UpdateOne(appCtx, paymentRequest, &eTag)
	}

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(id, err, verrs, "")
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.PaymentRequest{}, apperror.NewPreconditionFailedError(id, err)
		}

		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(id, "")
		default:
			return nil, apperror.NewQueryError("PaymentRequest", err, "")
		}
	}

	return paymentRequest, err
}
