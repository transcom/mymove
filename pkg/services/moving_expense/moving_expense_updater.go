package movingexpense

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseUpdater struct {
	checks []movingExpenseValidator
}

func NewMovingExpenseUpdater() services.MovingExpenseUpdater {
	return &movingExpenseUpdater{
		checks: updateChecks(),
	}
}

func (f *movingExpenseUpdater) UpdateMovingExpense(appCtx appcontext.AppContext, movingExpense models.MovingExpense, eTag string) (*models.MovingExpense, error) {
	originalMovingExpense, err := FetchMovingExpenseByID(appCtx, movingExpense.ID)

	if err != nil {
		return nil, err
	}

	if etag.GenerateEtag(originalMovingExpense.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalMovingExpense.ID, nil)
	}

	mergedMovingExpense := mergeMovingExpense(movingExpense, *originalMovingExpense)

	err = validateMovingExpense(appCtx, &mergedMovingExpense, originalMovingExpense, f.checks...)

	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndUpdate(&mergedMovingExpense)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(originalMovingExpense.ID, err, verrs, "")
		} else if err != nil {
			return apperror.NewQueryError("Moving Expense", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &mergedMovingExpense, nil
}

func FetchMovingExpenseByID(appContext appcontext.AppContext, movingExpenseID uuid.UUID) (*models.MovingExpense, error) {
	var movingExpense models.MovingExpense

	err := appContext.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload("Document.UserUploads.Upload").Find(&movingExpense, movingExpenseID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(movingExpenseID, "while looking for MovingExpense")
		default:
			return nil, apperror.NewQueryError("MovingExpense fetch original", err, "")
		}
	}

	// Assuming the document itself will not be deleted because of the model not null requirements.
	// We could not return nothing in the case it is soft deleted but that behavior is a bit
	// undefined here.
	movingExpense.Document.UserUploads = movingExpense.Document.UserUploads.FilterDeleted()

	return &movingExpense, nil
}

func mergeMovingExpense(updatedMovingExpense models.MovingExpense, originalMovingExpense models.MovingExpense) models.MovingExpense {
	mergedMovingExpense := originalMovingExpense

	movingExpenseType := services.SetOptionalStringField((*string)(updatedMovingExpense.MovingExpenseType), (*string)(mergedMovingExpense.MovingExpenseType))

	if movingExpenseType != nil {
		movingExpenseReceiptType := models.MovingExpenseReceiptType(*movingExpenseType)
		mergedMovingExpense.MovingExpenseType = &movingExpenseReceiptType

		if movingExpenseReceiptType == models.MovingExpenseReceiptTypeStorage {
			mergedMovingExpense.SITStartDate = services.SetOptionalDateTimeField(updatedMovingExpense.SITStartDate, mergedMovingExpense.SITStartDate)
			mergedMovingExpense.SITEndDate = services.SetOptionalDateTimeField(updatedMovingExpense.SITEndDate, mergedMovingExpense.SITEndDate)
		} else if originalMovingExpense.MovingExpenseType != nil && *originalMovingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage {
			// The receipt type has been changed from storage to something else so we should clear
			// the start and end values
			mergedMovingExpense.SITStartDate = nil
			mergedMovingExpense.SITEndDate = nil
		}

	} else {
		mergedMovingExpense.MovingExpenseType = nil
	}

	movingExpenseStatus := services.SetOptionalStringField((*string)(updatedMovingExpense.Status), (*string)(mergedMovingExpense.Status))
	if movingExpenseStatus != nil {
		ppmDocumentStatus := models.PPMDocumentStatus(*movingExpenseStatus)
		mergedMovingExpense.Status = &ppmDocumentStatus

		if ppmDocumentStatus == models.PPMDocumentStatusExcluded || ppmDocumentStatus == models.PPMDocumentStatusRejected {
			mergedMovingExpense.Reason = services.SetOptionalStringField(updatedMovingExpense.Reason, mergedMovingExpense.Reason)
		} else {
			// if that status is changed back to approved then we should clear the reason value
			mergedMovingExpense.Reason = nil
		}
	} else {
		mergedMovingExpense.Status = nil
	}

	mergedMovingExpense.Description = services.SetOptionalStringField(updatedMovingExpense.Description, mergedMovingExpense.Description)
	mergedMovingExpense.Amount = services.SetNoNilOptionalCentField(updatedMovingExpense.Amount, mergedMovingExpense.Amount)
	mergedMovingExpense.PaidWithGTCC = services.SetNoNilOptionalBoolField(updatedMovingExpense.PaidWithGTCC, mergedMovingExpense.PaidWithGTCC)
	mergedMovingExpense.MissingReceipt = services.SetNoNilOptionalBoolField(updatedMovingExpense.MissingReceipt, mergedMovingExpense.MissingReceipt)

	// TBD may be able to use the updater service for soft deleting instead of adding a dedicated one
	mergedMovingExpense.DeletedAt = services.SetOptionalDateTimeField(updatedMovingExpense.DeletedAt, mergedMovingExpense.DeletedAt)

	return mergedMovingExpense
}
