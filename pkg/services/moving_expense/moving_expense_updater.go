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
	checks    []movingExpenseValidator
	estimator services.PPMEstimator
}

func NewCustomerMovingExpenseUpdater(estimator services.PPMEstimator) services.MovingExpenseUpdater {
	return &movingExpenseUpdater{
		estimator: estimator,
		checks:    customerUpdateChecks(),
	}
}

func NewOfficeMovingExpenseUpdater(estimator services.PPMEstimator) services.MovingExpenseUpdater {
	return &movingExpenseUpdater{
		estimator: estimator,
		checks:    officeUpdateChecks(),
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

	// we only update the submitted dates if the moving expense is updated by the Customer
	if appCtx.Session().IsMilApp() {
		if mergedMovingExpense.Amount != nil {
			mergedMovingExpense.SubmittedAmount = mergedMovingExpense.Amount
		}
		if mergedMovingExpense.MovingExpenseType != nil {
			mergedMovingExpense.SubmittedMovingExpenseType = mergedMovingExpense.MovingExpenseType
		}
		if mergedMovingExpense.Description != nil {
			mergedMovingExpense.SubmittedDescription = mergedMovingExpense.Description
		}
		if mergedMovingExpense.SITStartDate != nil {
			mergedMovingExpense.SubmittedSITStartDate = mergedMovingExpense.SITStartDate
		}
		if mergedMovingExpense.SITEndDate != nil {
			mergedMovingExpense.SubmittedSITEndDate = mergedMovingExpense.SITEndDate
		}
	}

	if *mergedMovingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage &&
		mergedMovingExpense.PPMShipment.Status == models.PPMShipmentStatusNeedsCloseout {

		// We set sitExpected to true because this is a storage moving expense therefore SIT has to be true
		// The case where this could be false at this point is when the Customer created the shipment they answered No to SIT Expected question,
		// but later decided they needed SIT and submitted a moving expense for storage or if the Service Counselor adds one.
		sitExpected := true
		mergedMovingExpense.PPMShipment.SITExpected = &sitExpected
		estimatedCost, err := f.estimator.CalculatePPMSITEstimatedCost(appCtx, &mergedMovingExpense.PPMShipment)

		if err != nil {
			return nil, err
		}

		mergedMovingExpense.SITEstimatedCost = estimatedCost
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
			mergedMovingExpense.PPMShipment = updatedMovingExpense.PPMShipment
			mergedMovingExpense.PPMShipment.SITEstimatedWeight = updatedMovingExpense.WeightStored
			mergedMovingExpense.PPMShipment.SITEstimatedEntryDate = services.SetOptionalDateTimeField(updatedMovingExpense.SITStartDate, mergedMovingExpense.PPMShipment.SITEstimatedEntryDate)
			mergedMovingExpense.PPMShipment.SITEstimatedDepartureDate = services.SetOptionalDateTimeField(updatedMovingExpense.SITEndDate, mergedMovingExpense.PPMShipment.SITEstimatedDepartureDate)
			mergedMovingExpense.PPMShipment.SITLocation = updatedMovingExpense.SITLocation
			mergedMovingExpense.SITStartDate = services.SetOptionalDateTimeField(updatedMovingExpense.SITStartDate, mergedMovingExpense.SITStartDate)
			mergedMovingExpense.SITEndDate = services.SetOptionalDateTimeField(updatedMovingExpense.SITEndDate, mergedMovingExpense.SITEndDate)
			mergedMovingExpense.SITEstimatedCost = updatedMovingExpense.SITEstimatedCost
			mergedMovingExpense.SITReimburseableAmount = updatedMovingExpense.SITReimburseableAmount
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
			mergedMovingExpense.SITEstimatedCost = nil
			mergedMovingExpense.SITReimburseableAmount = nil
		}

		if movingExpenseReceiptType == models.MovingExpenseReceiptTypeSmallPackage {
			mergedMovingExpense.TrackingNumber = services.SetOptionalStringField(updatedMovingExpense.TrackingNumber, mergedMovingExpense.TrackingNumber)
			mergedMovingExpense.IsProGear = services.SetNoNilOptionalBoolField(updatedMovingExpense.IsProGear, mergedMovingExpense.IsProGear)

			if updatedMovingExpense.ProGearBelongsToSelf != nil {
				mergedMovingExpense.ProGearBelongsToSelf = updatedMovingExpense.ProGearBelongsToSelf
			}
			if updatedMovingExpense.ProGearDescription != nil {
				mergedMovingExpense.ProGearDescription = updatedMovingExpense.ProGearDescription
			}
			if *updatedMovingExpense.WeightShipped != 0 {
				mergedMovingExpense.WeightShipped = services.SetOptionalPoundField(updatedMovingExpense.WeightShipped, mergedMovingExpense.WeightShipped)
			}
			// description is not provided for small package expenses
			mergedMovingExpense.Description = nil
			updatedMovingExpense.Description = nil
		} else if originalMovingExpense.MovingExpenseType != nil && *originalMovingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeSmallPackage {
			// clearing related SPR values if expense type is being changed
			mergedMovingExpense.TrackingNumber = nil
			mergedMovingExpense.IsProGear = nil
			mergedMovingExpense.ProGearBelongsToSelf = nil
			mergedMovingExpense.ProGearDescription = nil
			mergedMovingExpense.WeightShipped = nil
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
