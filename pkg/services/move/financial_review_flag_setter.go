package move

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type financialReviewFlagSetter struct {
}

func NewFinancialReviewFlagSetter() services.MoveFinancialReviewFlagSetter {
	return &financialReviewFlagSetter{}
}

// SetFinancialReviewFlag is used to edit fields on a Move object that support the Financial Review Flag process for office users
func (f financialReviewFlagSetter) SetFinancialReviewFlag(appCtx appcontext.AppContext, moveID uuid.UUID, eTag string, flagForReview bool, remarks *string) (*models.Move, error) {
	// Let's get the existing move
	move, err := f.fetchMove(appCtx, moveID)
	if err != nil {
		return nil, err
	}
	updatedMove, err := f.updateMove(appCtx, move, flagForReview, remarks, eTag)
	if err != nil {
		return nil, err
	}

	return updatedMove, nil
}

func (f *financialReviewFlagSetter) fetchMove(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Move, error) {
	move := &models.Move{}
	err := appCtx.DB().Find(move, moveID)
	if err != nil {
		return nil, apperror.NewNotFoundError(moveID, "while looking for move")
	}
	return move, nil
}

func (f *financialReviewFlagSetter) updateMove(appCtx appcontext.AppContext, move *models.Move,
	flagForReview bool, remarks *string, eTag string) (*models.Move, error) {

	// Using a new move object to track the proposed changes to the move.
	// This will get used in validation function that checks against rules in the move service.
	moveDelta := move

	existingETag := etag.GenerateEtag(move.UpdatedAt)
	if existingETag != eTag {
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	moveDelta.FinancialReviewFlag = flagForReview
	currentTime := time.Now()
	// If this is true then we need to save the remarks and timestamp fields
	if flagForReview {
		moveDelta.FinancialReviewFlagSetAt = &currentTime
		moveDelta.FinancialReviewRemarks = remarks
		// If it isn't true, it's false which means we need to nil out the remarks and timestamp
	} else {
		moveDelta.FinancialReviewRemarks = nil
		moveDelta.FinancialReviewFlagSetAt = nil
	}

	// validate the proposed changes against rules in validation.go
	if verrs := validateMove(appCtx, *move, moveDelta, remarksNeededForFinancialFlag(), checkFinancialFlagRemoval()); verrs != nil {
		return move, verrs
	}

	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := txnAppCtx.DB().ValidateAndUpdate(moveDelta)
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
