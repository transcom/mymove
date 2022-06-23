package customersupportremarks

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type customerSupportRemarkDeleter struct {
}

func NewCustomerSupportRemarkDeleter() services.CustomerSupportRemarkDeleter {
	return &customerSupportRemarkDeleter{}
}

func (o customerSupportRemarkDeleter) DeleteCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemarkID uuid.UUID) error {
	var remark models.CustomerSupportRemark
	err := appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope()).Find(&remark, customerSupportRemarkID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(customerSupportRemarkID, "attempting to delete CustomerSupportRemark")
		default:
			return apperror.NewQueryError("CustomerSupportRemark", err, "")
		}
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err := utilities.SoftDestroy(appCtx.DB(), &remark)
		// TODO is this error handling here necessary? can we just bubble the error up?
		if err != nil {
			switch err.Error() {
			case "error updating model":
				return apperror.NewUnprocessableEntityError("while updating model")
			case "this model does not have deleted_at field":
				return apperror.NewPreconditionFailedError(remark.ID, errors.New("model or sub table missing deleted_at field"))
			default:
				return apperror.NewInternalServerError("failed attempt to soft delete model")
			}
		}
		return nil
	})

	if transactionError != nil {
		return transactionError
	}
	return nil
}
