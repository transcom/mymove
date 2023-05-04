package internalapi

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForReimbursementModel(r *models.Reimbursement) *internalmessages.Reimbursement {
	if r == nil {
		return nil
	}
	methodOfReceipt := internalmessages.MethodOfReceipt(r.MethodOfReceipt)
	status := internalmessages.ReimbursementStatus(r.Status)

	return &internalmessages.Reimbursement{
		ID:              strfmt.UUID(r.ID.String()),
		MethodOfReceipt: &methodOfReceipt,
		RequestedAmount: models.Int64Pointer(int64(r.RequestedAmount)),
		RequestedDate:   (*strfmt.Date)(r.RequestedDate),
		Status:          &status,
	}
}
