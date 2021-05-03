package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	webhooksubscriptionop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/webhook_subscriptions"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// IndexWebhookSubscriptionsHandler returns a list of webhook subscriptions via GET /webhook_subscriptions
type IndexWebhookSubscriptionsHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of webhook subscriptions
func (h IndexWebhookSubscriptionsHandler) Handle(params webhooksubscriptionop.IndexWebhookSubscriptionsParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	ordering := query.NewQueryOrder(params.Sort, params.Order)
	pagination := h.NewPagination(params.Page, params.PerPage)

	var webhookSubscriptions models.WebhookSubscriptions
	err := h.ListFetcher.FetchRecordList(&webhookSubscriptions, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalWebhookSubscriptionsCount, err := h.ListFetcher.FetchRecordCount(&webhookSubscriptions, queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedWebhookSubscriptionsCount := len(webhookSubscriptions)

	payload := make(adminmessages.WebhookSubscriptions, queriedWebhookSubscriptionsCount)

	for i, s := range webhookSubscriptions {
		payload[i] = payloads.WebhookSubscriptionPayload(s)
	}

	return webhooksubscriptionop.NewIndexWebhookSubscriptionsOK().WithContentRange(fmt.Sprintf("webhookSubscriptions %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedWebhookSubscriptionsCount, totalWebhookSubscriptionsCount)).WithPayload(payload)
}

// GetWebhookSubscriptionHandler returns one webhookSubscription via GET /webhook_subscriptions/:ID
type GetWebhookSubscriptionHandler struct {
	handlers.HandlerContext
	services.WebhookSubscriptionFetcher
	services.NewQueryFilter
}

// Handle retrieves a webhook subscription
func (h GetWebhookSubscriptionHandler) Handle(params webhooksubscriptionop.GetWebhookSubscriptionParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	webhookSubscriptionID := uuid.FromStringOrNil(params.WebhookSubscriptionID.String())
	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscriptionID)}

	webhookSubscription, err := h.WebhookSubscriptionFetcher.FetchWebhookSubscription(queryFilters)

	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := payloads.WebhookSubscriptionPayload(webhookSubscription)
	return webhooksubscriptionop.NewGetWebhookSubscriptionOK().WithPayload(payload)
}

// CreateWebhookSubscriptionHandler is the handler for creating users.
type CreateWebhookSubscriptionHandler struct {
	handlers.HandlerContext
	services.WebhookSubscriptionCreator
	services.NewQueryFilter
}

// Handle creates an admin user
func (h CreateWebhookSubscriptionHandler) Handle(params webhooksubscriptionop.CreateWebhookSubscriptionParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	subscription := payloads.WebhookSubscriptionModelFromCreate(params.WebhookSubscription)
	subscriberIDFilter := []services.QueryFilter{
		h.NewQueryFilter("id", "=", subscription.SubscriberID),
	}

	createdWebhookSubscription, verrs, err := h.WebhookSubscriptionCreator.CreateWebhookSubscription(subscription, subscriberIDFilter)

	if verrs != nil {
		logger.Error("Error saving webhook subscription", zap.Error(verrs))
		return webhooksubscriptionop.NewCreateWebhookSubscriptionInternalServerError()
	}

	if err != nil {
		logger.Error("Error saving webhook subscription", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return webhooksubscriptionop.NewCreateWebhookSubscriptionBadRequest()
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("adminapi.CreateWebhookSubscriptionHandler query error", zap.Error(e.Unwrap()))
			}
			return webhooksubscriptionop.NewCreateWebhookSubscriptionInternalServerError()
		default:
			return webhooksubscriptionop.NewCreateWebhookSubscriptionInternalServerError()
		}
	}

	returnPayload := payloads.WebhookSubscriptionPayload(*createdWebhookSubscription)
	return webhooksubscriptionop.NewCreateWebhookSubscriptionCreated().WithPayload(returnPayload)
}

// UpdateWebhookSubscriptionHandler returns an updated webhook subscription via PATCH
type UpdateWebhookSubscriptionHandler struct {
	handlers.HandlerContext
	services.WebhookSubscriptionUpdater
	services.NewQueryFilter
}

// Handle updates a webhook subscription
func (h UpdateWebhookSubscriptionHandler) Handle(params webhooksubscriptionop.UpdateWebhookSubscriptionParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	payload := params.WebhookSubscription

	// Checks that ID in body matches ID in query
	payloadID := uuid.FromStringOrNil(payload.ID.String())
	if payloadID != uuid.Nil && params.WebhookSubscriptionID != payload.ID {
		return webhooksubscriptionop.NewUpdateWebhookSubscriptionUnprocessableEntity()
	}

	// If no ID in body, use query ID
	payload.ID = params.WebhookSubscriptionID

	// Convert payload to model
	webhookSubscription := payloads.WebhookSubscriptionModel(payload)

	// Note we are not checking etag as adminapi does not seem to use this
	updatedWebhookSubscription, err := h.WebhookSubscriptionUpdater.UpdateWebhookSubscription(webhookSubscription, payload.Severity, &params.IfMatch)

	// Return error response if not successful
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			logger.Error("Error finding webhookSubscription to update")
			return webhooksubscriptionop.NewUpdateWebhookSubscriptionNotFound()
		}
		switch err.(type) {
		case services.PreconditionFailedError:
			logger.Error("Error updating webhookSubscription due to stale eTag")
			return webhooksubscriptionop.NewUpdateWebhookSubscriptionPreconditionFailed()
		default:
			logger.Error(fmt.Sprintf("Error updating webhookSubscription %s", params.WebhookSubscriptionID.String()), zap.Error(err))
			return handlers.ResponseForError(logger, err)
		}
	}

	// Convert model back to a payload and return to caller
	payload = payloads.WebhookSubscriptionPayload(*updatedWebhookSubscription)
	return webhooksubscriptionop.NewUpdateWebhookSubscriptionOK().WithPayload(payload)
}
