package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	webhooksubscriptionop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/webhook_subscriptions"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForWebhookSubscriptionModel(subscription models.WebhookSubscription) *adminmessages.WebhookSubscription {
	severity := int64(subscription.Severity)

	return &adminmessages.WebhookSubscription{
		ID:           *handlers.FmtUUID(subscription.ID),
		SubscriberID: *handlers.FmtUUID(subscription.SubscriberID),
		CallbackURL:  subscription.CallbackURL,
		Severity:     &severity,
		EventKey:     subscription.EventKey,
		Status:       adminmessages.WebhookSubscriptionStatus(subscription.Status),
		CreatedAt:    *handlers.FmtDateTime(subscription.CreatedAt),
		UpdatedAt:    *handlers.FmtDateTime(subscription.UpdatedAt),
	}
}

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
		payload[i] = payloadForWebhookSubscriptionModel(s)
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

	payload := payloadForWebhookSubscriptionModel(webhookSubscription)
	return webhooksubscriptionop.NewGetWebhookSubscriptionOK().WithPayload(payload)
}
