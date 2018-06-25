package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForReimbursementModel(r models.Reimbursement) *internalmessages.Reimbursement {

	methodOfReceipt := internalmessages.MethodOfReceipt(r.MethodOfReceipt)
	status := internalmessages.ReimbursementStatus(r.Status)

	return &internalmessages.Reimbursement{
		ID:              strfmt.UUID(r.ID.String()),
		MethodOfReceipt: &methodOfReceipt,
		RequestedAmount: swag.Int64(int64(r.RequestedAmount)),
		RequestedDate:   (*strfmt.Date)(r.RequestedDate),
		Status:          &status,
	}
}
