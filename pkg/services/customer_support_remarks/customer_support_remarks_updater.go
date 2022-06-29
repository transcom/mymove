package customersupportremarks

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	customersupportremarksop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type CustomerSupportRemarkUpdater struct {
}

func NewCustomerSupportRemarkUpdater() services.CustomerSupportRemarkUpdater {
	return &CustomerSupportRemarkUpdater{}
}

func (o CustomerSupportRemarkUpdater) UpdateCustomerSupportRemark(appCtx appcontext.AppContext, params customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams) (*models.CustomerSupportRemark, error) {
	var remark models.CustomerSupportRemark
	remarkID := params.CustomerSupportRemarkID
	err := appCtx.DB().Find(&remark, remarkID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			remarkID, _ := uuid.FromString(remarkID.String())
			return nil, apperror.NewNotFoundError(remarkID, "customer support remark not found")
		default:
			return nil, apperror.NewQueryError("CustomerSupportRemark", err, "")
		}
	}

	/*
		userID on the session must match original remark creator

		should this be session().officeUserID?
	*/
	sessionUserID := appCtx.Session().UserID

	if remark.OfficeUser.UserID != &sessionUserID {
		appCtx.Logger().Warn("Customer Support Remarks may only be edited by the user who created them.", zap.String("Customer Support RemarkID", remarkID.String()))
		return nil, apperror.NewForbiddenError("Action not allowed")
	}

	remark.Content = *params.Body.Content
	remark.UpdatedAt = time.Now()

	verrs, err := appCtx.DB().Q().Connection.ValidateAndUpdate(&remark)
	if verrs.Count() != 0 || err != nil {
		return nil, err
	}

	return &remark, nil
}
