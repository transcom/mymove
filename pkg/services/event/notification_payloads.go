package event

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrder has a custom payload definition, because it differs from the one in primeapi
type MoveTaskOrder struct {

	// available to prime at
	// Read Only: true
	// Format: date-time
	AvailableToPrimeAt *strfmt.DateTime `json:"availableToPrimeAt,omitempty"`

	// created at
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// e tag
	// Read Only: true
	ETag string `json:"eTag,omitempty"`

	// id
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// is canceled
	IsCanceled *bool `json:"isCanceled,omitempty"`

	// order ID
	// Format: uuid
	OrderID strfmt.UUID `json:"orderID,omitempty"`

	// ppm estimated weight
	PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

	// ppm type
	// Enum: [FULL PARTIAL]
	PpmType string `json:"ppmType,omitempty"`

	// reference Id
	ReferenceID string `json:"referenceId,omitempty"`

	// updated at
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
}

// MoveTaskOrderModelToPayload converts the Move model into a MoveTaskOrder payload
// Ideally it would be great to have this definition in the yaml - OpenAPI 3.0 should have
// ability to put callback payloads in the yaml
func MoveTaskOrderModelToPayload(moveTaskOrder *models.Move) *MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	payload := &MoveTaskOrder{
		ID:                 strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:          strfmt.DateTime(moveTaskOrder.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(moveTaskOrder.AvailableToPrimeAt),
		IsCanceled:         moveTaskOrder.IsCanceled(),
		OrderID:            strfmt.UUID(moveTaskOrder.OrdersID.String()),
		ReferenceID:        *moveTaskOrder.ReferenceID,
		UpdatedAt:          strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:               etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}

	if moveTaskOrder.PPMEstimatedWeight != nil {
		payload.PpmEstimatedWeight = int64(*moveTaskOrder.PPMEstimatedWeight)
	}

	if moveTaskOrder.PPMType != nil {
		payload.PpmType = *moveTaskOrder.PPMType
	}

	return payload
}

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
		ReferenceID:      paymentServiceItem.ReferenceID,
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
		copyOfPaymentServiceItem := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItemModelToPayload(&copyOfPaymentServiceItem)
	}
	return &payload
}
