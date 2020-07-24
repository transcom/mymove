package events

import (
	"go.uber.org/zap"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// PaymentRequestModelToPayload This is a sample first example, will need to move to payloads_to_model file
func paymentRequestModelToPayload(paymentRequest *models.PaymentRequest) *primemessages.PaymentRequest {
	if paymentRequest == nil {
		return nil
	}

	return &primemessages.PaymentRequest{
		ID:                   strfmt.UUID(paymentRequest.ID.String()),
		IsFinal:              &paymentRequest.IsFinal,
		MoveTaskOrderID:      strfmt.UUID(paymentRequest.MoveTaskOrderID.String()),
		PaymentRequestNumber: paymentRequest.PaymentRequestNumber,
		RejectionReason:      paymentRequest.RejectionReason,
		Status:               primemessages.PaymentRequestStatus(paymentRequest.Status),
		ETag:                 etag.GenerateEtag(paymentRequest.UpdatedAt),
	}
}

// EventNotificationsHandler receives notifications from the events package
func EventNotificationsHandler(event *Event, db *pop.Connection, logger handlers.Logger) error {

	// Currently it logs information about the event. Eventually it will create an entry
	// in the notification table in the database
	logger.Info("Event handler ran:",
		zap.String("endpoint", string(event.EndpointKey)),
		zap.String("mtoID", event.MtoID.String()),
		zap.String("objectID", event.UpdatedObjectID.String()))

	// Create a query builder to query DB for notification details
	notificationQueryBuilder := query.NewQueryBuilder(db)

	//Todo Log who made the call
	//Todo Log whether mto was available

	//Log the payload
	filter := []services.QueryFilter{query.NewQueryFilter("id", "=", event.UpdatedObjectID.String())}
	//Get the type of model which is stored in the eventType
	modelBeingUpdated := event.EventType.ModelInstance
	// Based on which model was updated, construct the proper payload
	switch modelBeingUpdated.(type) {
	case models.PaymentRequest:
		model := models.PaymentRequest{}
		err := notificationQueryBuilder.FetchOne(&model, filter)
		if err != nil {
			logger.Info("Notification error:")
			return err
		}
		payload := paymentRequestModelToPayload(&model)
		logger.Info("Notification payload:", zap.Any("payload", *payload))
	}

	return nil
}
