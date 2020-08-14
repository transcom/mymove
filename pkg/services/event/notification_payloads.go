package event

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
)

// PaymentRequestModelToPayload creates a model we can use to populate the payload column
// Currently we are using the primemessages struct as the payload.
func PaymentRequestModelToPayload(paymentRequest *models.PaymentRequest) *primemessages.PaymentRequest {
	if paymentRequest == nil {
		return nil
	}

	paymentServiceItems := PaymentServiceItemsModelToPayload(&paymentRequest.PaymentServiceItems)
	return &primemessages.PaymentRequest{
		ID:                   strfmt.UUID(paymentRequest.ID.String()),
		IsFinal:              &paymentRequest.IsFinal,
		MoveTaskOrderID:      strfmt.UUID(paymentRequest.MoveTaskOrderID.String()),
		PaymentRequestNumber: paymentRequest.PaymentRequestNumber,
		RejectionReason:      paymentRequest.RejectionReason,
		Status:               primemessages.PaymentRequestStatus(paymentRequest.Status),
		PaymentServiceItems:  *paymentServiceItems,
		ETag:                 etag.GenerateEtag(paymentRequest.UpdatedAt),
	}
}

// PaymentServiceItemModelToPayload payload
func PaymentServiceItemModelToPayload(paymentServiceItem *models.PaymentServiceItem) *primemessages.PaymentServiceItem {
	if paymentServiceItem == nil {
		return nil
	}

	payload := &primemessages.PaymentServiceItem{
		ID:               strfmt.UUID(paymentServiceItem.ID.String()),
		PaymentRequestID: strfmt.UUID(paymentServiceItem.PaymentRequestID.String()),
		MtoServiceItemID: strfmt.UUID(paymentServiceItem.MTOServiceItemID.String()),
		Status:           primemessages.PaymentServiceItemStatus(paymentServiceItem.Status),
		RejectionReason:  paymentServiceItem.RejectionReason,
		ETag:             etag.GenerateEtag(paymentServiceItem.UpdatedAt),
	}

	if paymentServiceItem.PriceCents != nil {
		payload.PriceCents = swag.Int64(int64(*paymentServiceItem.PriceCents))
	}

	return payload
}

// PaymentServiceItemsModelToPayload payload
func PaymentServiceItemsModelToPayload(paymentServiceItems *models.PaymentServiceItems) *primemessages.PaymentServiceItems {
	if paymentServiceItems == nil {
		return nil
	}

	payload := make(primemessages.PaymentServiceItems, len(*paymentServiceItems))

	for i, p := range *paymentServiceItems {
		payload[i] = PaymentServiceItemModelToPayload(&p)
	}
	return &payload
}
