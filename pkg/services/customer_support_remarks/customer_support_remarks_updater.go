package customersupportremarks

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type CustomerSupportRemarkUpdater struct {
}

func NewCustomerSupportRemarkUpdater() services.CustomerSupportRemarkUpdater {
	return &CustomerSupportRemarkUpdater{}
}

func (o CustomerSupportRemarkUpdater) UpdateCustomerSupportRemark(appCtx appcontext.AppContext, payload ghcmessages.UpdateCustomerSupportRemarkPayload) (*models.CustomerSupportRemark, error) {
	var remark models.CustomerSupportRemark
	err := appCtx.DB().Find(&remark, payload.ID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			remarkID, _ := uuid.FromString(payload.ID.String())
			return nil, apperror.NewNotFoundError(remarkID, "customer support remark not found")
		default:
			return nil, apperror.NewQueryError("CustomerSupportRemark", err, "")
		}
	}

	remark.Content = *payload.Content
	remark.UpdatedAt = time.Now()

	verrs, err := appCtx.DB().Q().Connection.ValidateAndUpdate(&remark)
	if verrs.Count() != 0 || err != nil {
		return nil, err
	}

	return &remark, nil
}
