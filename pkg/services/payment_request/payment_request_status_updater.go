package paymentrequest

import (
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type paymentRequestStatusQueryBuilder interface {
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
	verrs, err := p.builder.UpdateOne(paymentRequest, &eTag)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(id, err, verrs, "")
	}

	if err != nil {
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return nil, services.NewNotFoundError(id)
		}

		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.PaymentRequest{}, services.NewPreconditionFailedError(id, err)
		}
	}

	return paymentRequest, err
}
