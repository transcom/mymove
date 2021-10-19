package move

import (
	"time"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type financialReviewFlagSetter struct {
}

func NewFinancialReviewFlagSetter() services.MoveFinancialReviewFlagSetter {
	return &financialReviewFlagSetter{}
}

func (f financialReviewFlagSetter) SetFinancialReviewFlag(appCtx appcontext.AppContext, moveID uuid.UUID, eTag string, remarks string) (*models.Move, error) {
	if remarks == "" {
		verrs := validate.NewErrors()
		verrs.Add("remarks", "must not be empty")
		return nil, apperror.NewInvalidInputError(moveID, nil, verrs, "")
	}

	move := &models.Move{}
	err := appCtx.DB().Find(move, moveID)
	if err != nil {
		return nil, apperror.NewNotFoundError(moveID, "while looking for move")
	}

	existingETag := etag.GenerateEtag(move.UpdatedAt)
	if existingETag != eTag {
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
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
			return apperror.NewInvalidInputError(
				move.ID, err, verrs, "Validation errors found while setting financial review flag on move")
		} else if err != nil {
			return apperror.NewQueryError("Move", err, "Failed to request financial review for move")
		}

		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return move, nil
}
