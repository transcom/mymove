package movingexpense

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func checkID() movingExpenseValidator {
	return movingExpenseValidatorFunc(func(_ appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error {
		verrs := validate.NewErrors()

		if newMovingExpense == nil || originalMovingExpense == nil {
			return verrs
		}

		if newMovingExpense.ID != originalMovingExpense.ID {
			verrs.Add("ID", "new MovingExpense ID must match original MovingExpense ID")
		}

		return verrs
	})
}

func checkBaseRequiredFields() movingExpenseValidator {
	return movingExpenseValidatorFunc(func(_ appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error {
		verrs := validate.NewErrors()

		if newMovingExpense.PPMShipmentID.IsNil() {
			verrs.Add("PPMShipmentID", "PPMShipmentID must exist")
		}

		if newMovingExpense.Document.ServiceMemberID.IsNil() {
			verrs.Add("ServiceMemberID", "Document ServiceMemberID must exist")
		}

		return verrs
	})
}

func checkAdditionalRequiredFields() movingExpenseValidator {
	return movingExpenseValidatorFunc(func(_ appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error {
		verrs := validate.NewErrors()

		// the model Validate should test for allowed values
		if newMovingExpense.MovingExpenseType == nil || *newMovingExpense.MovingExpenseType == "" {
			verrs.Add("MovingExpenseType", "MovingExpenseType must exist")
		} else if *newMovingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage {
			if newMovingExpense.SITStartDate == nil || newMovingExpense.SITStartDate.IsZero() {
				verrs.Add("SITStartDate", "SITStartDate is required for storage expenses")
			}

			if newMovingExpense.SITEndDate == nil || newMovingExpense.SITEndDate.IsZero() {
				verrs.Add("SITEndDate", "SITEndDate is required for storage expenses")
			}

			if newMovingExpense.SITStartDate != nil && newMovingExpense.SITEndDate != nil {
				if newMovingExpense.SITEndDate.Before(*newMovingExpense.SITStartDate) {
					verrs.Add("SITStartDate", "SITStartDate must be before SITEndDate")
				}
			}
		}

		if newMovingExpense.Description == nil || *newMovingExpense.Description == "" {
			verrs.Add("Description", "Description must have a value of at least 0")
		}

		if newMovingExpense.PaidWithGTCC == nil {
			verrs.Add("PaidWithGTCC", "PaidWithGTCC is required")
		}

		if newMovingExpense.Amount == nil || *newMovingExpense.Amount < 1 {
			verrs.Add("Amount", "Amount must have a value of at least 1")
		}

		if newMovingExpense.MissingReceipt == nil {
			verrs.Add("MissingReceipt", "MissingReceipt is required")
		}

		if len(originalMovingExpense.Document.UserUploads) < 1 {
			verrs.Add("Document", "At least 1 receipt file is required")
		}

		return verrs
	})
}

// verifyReasonAndStatusAreConstant ensures that the reason and status fields do not change
func verifyReasonAndStatusAreConstant() movingExpenseValidator {
	return movingExpenseValidatorFunc(func(_ appcontext.AppContext, newMovingExpense *models.MovingExpense, originalMovingExpense *models.MovingExpense) error {
		verrs := validate.NewErrors()

		if (originalMovingExpense.Status == nil && newMovingExpense.Status != nil) ||
			(originalMovingExpense.Status != nil && newMovingExpense.Status == nil) ||
			(originalMovingExpense.Status != nil && newMovingExpense.Status != nil && *originalMovingExpense.Status != *newMovingExpense.Status) {
			verrs.Add("Status", "status cannot be modified")
		}

		if (originalMovingExpense.Reason == nil && newMovingExpense.Reason != nil) ||
			(originalMovingExpense.Reason != nil && newMovingExpense.Reason == nil) ||
			(originalMovingExpense.Reason != nil && newMovingExpense.Reason != nil && *originalMovingExpense.Reason != *newMovingExpense.Reason) {
			verrs.Add("Reason", "reason cannot be modified")
		}

		return verrs
	})
}

// verifyReasonAndStatusAreValid ensures that the reason and status are only changed in valid ways
func verifyReasonAndStatusAreValid() movingExpenseValidator {
	return movingExpenseValidatorFunc(func(_ appcontext.AppContext, newMovingExpense *models.MovingExpense, _ *models.MovingExpense) error {
		verrs := validate.NewErrors()

		if newMovingExpense.Status != nil {
			if *newMovingExpense.Status == models.PPMDocumentStatusApproved && newMovingExpense.Reason != nil {
				verrs.Add("Reason", "reason must not be set if the status is Approved")
			} else if (*newMovingExpense.Status == models.PPMDocumentStatusExcluded || *newMovingExpense.Status == models.PPMDocumentStatusRejected) &&
				(newMovingExpense.Reason == nil || *newMovingExpense.Reason == "") {
				verrs.Add("Reason", "reason is mandatory if the status is Excluded or Rejected")
			}
		} else if newMovingExpense.Reason != nil {
			verrs.Add("Reason", "reason should not be set if the status is not set")
		}

		return verrs
	})
}

func createChecks() []movingExpenseValidator {
	return []movingExpenseValidator{
		checkID(),
		checkBaseRequiredFields(),
	}
}

func customerUpdateChecks() []movingExpenseValidator {
	return []movingExpenseValidator{
		checkID(),
		checkBaseRequiredFields(),
		checkAdditionalRequiredFields(),
		verifyReasonAndStatusAreConstant(),
	}
}

func officeUpdateChecks() []movingExpenseValidator {
	return []movingExpenseValidator{
		checkID(),
		checkBaseRequiredFields(),
		checkAdditionalRequiredFields(),
		verifyReasonAndStatusAreValid(),
	}
}
