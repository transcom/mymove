package invoice

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ghcPaymentRequestJobRunner struct {
	db *pop.Connection
}

// NewGHCJobRunner returns an implementation of the GHCJobRunner interface
func NewGHCJobRunner(db *pop.Connection) services.GHCJobRunner {
	return &ghcPaymentRequestJobRunner{
		db: db,
	}
}

func (g ghcPaymentRequestJobRunner) ApprovedPaymentRequestFetcher() (models.PaymentRequests, error) {
	var reviewedPaymentRequests models.PaymentRequests
	err := g.db.Q().
		Where("status = ?", models.PaymentRequestStatusReviewed).
		All(&reviewedPaymentRequests)
	if err != nil {
		return reviewedPaymentRequests, services.NewQueryError("PaymentRequests", err, fmt.Sprintf("Could not find reviewed payment requests: %s", err))
	}
	return reviewedPaymentRequests, err
}
