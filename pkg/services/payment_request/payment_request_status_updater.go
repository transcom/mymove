package paymentrequest

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	moveservice "github.com/transcom/mymove/pkg/services/move"
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
		err := p.builder.FetchMany(appCtx, &paymentServiceItems, serviceItemFilter, nil, nil, nil)

		if err != nil {
			return nil, err
		}

		if len(paymentServiceItems) > 0 {
			return nil, apperror.NewConflictError(id, "All PaymentServiceItems must be approved or denied to review this PaymentRequest")
		}
	}

	paymentRequests := models.PaymentRequests{}
	moveID := paymentRequest.MoveTaskOrderID

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

	Qerr := appCtx.DB().Q().InnerJoin("moves", "payment_requests.move_id = moves.id").Where("moves.id = ?", moveID).All(&paymentRequests)
	if Qerr != nil {
		return nil, Qerr
	}

	paymentRequestNeedingReview := false
	for _, request := range paymentRequests {
		if request.Status != models.PaymentRequestStatusReviewed &&
			request.Status != models.PaymentRequestStatusReviewedAllRejected {
			paymentRequestNeedingReview = true
			break
		}
	}

	if !paymentRequestNeedingReview {
		_, err := moveservice.AssignedOfficeUserUpdater.DeleteAssignedOfficeUser(moveservice.AssignedOfficeUserUpdater{}, appCtx, moveID, roles.RoleTypeTIO)
		if err != nil {
			return nil, err
		}
		paymentRequest.MoveTaskOrder.TIOAssignedID = nil
		paymentRequest.MoveTaskOrder.TIOAssignedUser = nil
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
