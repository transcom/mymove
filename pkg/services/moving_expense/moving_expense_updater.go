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

func NewCustomerMovingExpenseUpdater() services.MovingExpenseUpdater {
	return &movingExpenseUpdater{
		checks: customerUpdateChecks(),
	}
}

func NewOfficeMovingExpenseUpdater() services.MovingExpenseUpdater {
	return &movingExpenseUpdater{
		checks: officeUpdateChecks(),
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

	if appCtx.Session().IsMilApp() {
		mergedMovingExpense.SubmittedAmount = mergedMovingExpense.Amount
		if mergedMovingExpense.SITStartDate != nil {
			mergedMovingExpense.SubmittedSITStartDate = mergedMovingExpense.SITStartDate
		}
		if mergedMovingExpense.SITEndDate != nil {
			mergedMovingExpense.SubmittedSITEndDate = mergedMovingExpense.SITEndDate
		}
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

	query := appContext.DB().Scope(utilities.ExcludeDeletedScope(models.MovingExpense{}))

	if appContext.Session().IsMilApp() {
		serviceMemberID := appContext.Session().ServiceMemberID

		query = query.
			LeftJoin("ppm_shipments", "ppm_shipments.id = moving_expenses.ppm_shipment_id").
			LeftJoin("mto_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
			LeftJoin("moves", "moves.id = mto_shipments.move_id").
			LeftJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", serviceMemberID)
	}

	err := query.EagerPreload("Document.UserUploads.Upload").Find(&movingExpense, movingExpenseID)

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

			// if weightStored was omitted we check for the zero value that is passed in and don't update it since we don't want to null out
			// a previous value
			if *updatedMovingExpense.WeightStored != 0 {
				mergedMovingExpense.WeightStored = services.SetOptionalPoundField(updatedMovingExpense.WeightStored, mergedMovingExpense.WeightStored)
			}

			if updatedMovingExpense.SITLocation != nil {
				mergedMovingExpense.SITLocation = updatedMovingExpense.SITLocation
			}
		} else if originalMovingExpense.MovingExpenseType != nil && *originalMovingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage {
			// The receipt type has been changed from storage to something else so we should clear
			// the start and end values
			mergedMovingExpense.SITStartDate = nil
			mergedMovingExpense.SITEndDate = nil
			mergedMovingExpense.WeightStored = nil
			mergedMovingExpense.SITLocation = nil
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
