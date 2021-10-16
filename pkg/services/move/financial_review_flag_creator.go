package move

import (
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type financialReviewFlagCreator struct {
}

func NewFinancialReviewFlagCreator() services.MoveFinancialReviewFlagCreator {
	return &financialReviewFlagCreator{}
}

func (f financialReviewFlagCreator) CreateFinancialReviewFlag(appCtx appcontext.AppContext, moveID uuid.UUID, remarks string) (*models.Move, error) {
	if remarks == "" {
		verrs := validate.NewErrors()
		verrs.Add("remarks", "must not be empty")
		return nil, services.NewInvalidInputError(moveID, nil, verrs, "")
	}

	move := &models.Move{}
	err := appCtx.DB().Find(move, moveID)

	if err != nil {
		return nil, services.NewNotFoundError(moveID, "while looking for move")
	}

	if move.FinancialReviewRequested {
		// If the flag has already been set, we do not want to update it
		return move, nil
	}

	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		move.FinancialReviewRequested = true
		currentTime := time.Now()
		move.FinancialReviewRequestedAt = &currentTime
		move.FinancialReviewRemarks = &remarks

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(move)
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(
				move.ID, err, verrs, "Validation errors found while setting financial review flag on move")
		} else if err != nil {
			return services.NewQueryError("Move", err, "Failed to request financial review for move")
		}

		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return move, nil
}
