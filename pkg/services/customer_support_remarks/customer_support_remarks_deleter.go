package customersupportremarks

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type customerSupportRemarkDeleter struct {
}

func NewCustomerSupportRemarkDeleter() services.CustomerSupportRemarkDeleter {
	return &customerSupportRemarkDeleter{}
}

func (o customerSupportRemarkDeleter) DeleteCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemarkID uuid.UUID) error {
	var remark models.CustomerSupportRemark
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Find(&remark, customerSupportRemarkID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(customerSupportRemarkID, "attempting to delete CustomerSupportRemark")
		default:
			return apperror.NewQueryError("CustomerSupportRemark", err, "")
		}
	}

	/*
		https://dp3.atlassian.net/browse/MB-12730
		MB-12730 udpdates to customer support remarks are restricted to the original remark creator
	*/
	sessionUserID := appCtx.Session().OfficeUserID

	if remark.OfficeUserID != sessionUserID {
		appCtx.Logger().Warn("Customer Support Remarks may only be edited by the user who created them.", zap.String("Customer Support RemarkID", customerSupportRemarkID.String()))

		return apperror.NewForbiddenError("Action not allowed")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err := utilities.SoftDestroy(appCtx.DB(), &remark)
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
