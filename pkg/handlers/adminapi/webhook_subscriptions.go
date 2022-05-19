package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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
	handlers.HandlerConfig
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of webhook subscriptions
func (h IndexWebhookSubscriptionsHandler) Handle(params webhooksubscriptionop.IndexWebhookSubscriptionsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
			queryFilters := []services.QueryFilter{}

			ordering := query.NewQueryOrder(params.Sort, params.Order)
			pagination := h.NewPagination(params.Page, params.PerPage)

			var webhookSubscriptions models.WebhookSubscriptions
			err := h.ListFetcher.FetchRecordList(appCtx, &webhookSubscriptions, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalWebhookSubscriptionsCount, err := h.ListFetcher.FetchRecordCount(appCtx, &webhookSubscriptions, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedWebhookSubscriptionsCount := len(webhookSubscriptions)

			payload := make(adminmessages.WebhookSubscriptions, queriedWebhookSubscriptionsCount)

			for i, s := range webhookSubscriptions {
				payload[i] = payloads.WebhookSubscriptionPayload(s)
			}

			return webhooksubscriptionop.NewIndexWebhookSubscriptionsOK().WithContentRange(fmt.Sprintf("webhookSubscriptions %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedWebhookSubscriptionsCount, totalWebhookSubscriptionsCount)).WithPayload(payload), nil
		})
}

// GetWebhookSubscriptionHandler returns one webhookSubscription via GET /webhook_subscriptions/:ID
type GetWebhookSubscriptionHandler struct {
	handlers.HandlerConfig
	services.WebhookSubscriptionFetcher
	services.NewQueryFilter
}

// Handle retrieves a webhook subscription
func (h GetWebhookSubscriptionHandler) Handle(params webhooksubscriptionop.GetWebhookSubscriptionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			webhookSubscriptionID := uuid.FromStringOrNil(params.WebhookSubscriptionID.String())
			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", webhookSubscriptionID)}

			webhookSubscription, err := h.WebhookSubscriptionFetcher.FetchWebhookSubscription(appCtx, queryFilters)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := payloads.WebhookSubscriptionPayload(webhookSubscription)
			return webhooksubscriptionop.NewGetWebhookSubscriptionOK().WithPayload(payload), nil
		})
}

// CreateWebhookSubscriptionHandler is the handler for creating users.
type CreateWebhookSubscriptionHandler struct {
	handlers.HandlerConfig
	services.WebhookSubscriptionCreator
	services.NewQueryFilter
}

// Handle creates an admin user
func (h CreateWebhookSubscriptionHandler) Handle(params webhooksubscriptionop.CreateWebhookSubscriptionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			subscription := payloads.WebhookSubscriptionModelFromCreate(params.WebhookSubscription)

			createdWebhookSubscription, verrs, err := h.WebhookSubscriptionCreator.CreateWebhookSubscription(appCtx, subscription)

			if verrs != nil {
				appCtx.Logger().Error("Error saving webhook subscription", zap.Error(verrs))
				return webhooksubscriptionop.NewCreateWebhookSubscriptionInternalServerError(), verrs
			}

			if err != nil {
				appCtx.Logger().Error("Error saving webhook subscription", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return webhooksubscriptionop.NewCreateWebhookSubscriptionBadRequest(), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("adminapi.CreateWebhookSubscriptionHandler query error", zap.Error(e.Unwrap()))
					}
					return webhooksubscriptionop.NewCreateWebhookSubscriptionInternalServerError(), err
				default:
					return webhooksubscriptionop.NewCreateWebhookSubscriptionInternalServerError(), err
				}
			}

			returnPayload := payloads.WebhookSubscriptionPayload(*createdWebhookSubscription)
			return webhooksubscriptionop.NewCreateWebhookSubscriptionCreated().WithPayload(returnPayload), nil
		})
}

// UpdateWebhookSubscriptionHandler returns an updated webhook subscription via PATCH
type UpdateWebhookSubscriptionHandler struct {
	handlers.HandlerConfig
	services.WebhookSubscriptionUpdater
	services.NewQueryFilter
}

// Handle updates a webhook subscription
func (h UpdateWebhookSubscriptionHandler) Handle(params webhooksubscriptionop.UpdateWebhookSubscriptionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.WebhookSubscription

			// Checks that ID in body matches ID in query
			payloadID := uuid.FromStringOrNil(payload.ID.String())
			if payloadID != uuid.Nil && params.WebhookSubscriptionID != payload.ID {
				webhookErr := apperror.NewUnprocessableEntityError("Payload ID does not match query ID")
				return webhooksubscriptionop.NewUpdateWebhookSubscriptionUnprocessableEntity(), webhookErr
			}

			// If no ID in body, use query ID
			payload.ID = params.WebhookSubscriptionID

			// Convert payload to model
			webhookSubscription := payloads.WebhookSubscriptionModel(payload)

			// Note we are not checking etag as adminapi does not seem to use this
			updatedWebhookSubscription, err := h.WebhookSubscriptionUpdater.UpdateWebhookSubscription(appCtx, webhookSubscription, payload.Severity, &params.IfMatch)

			// Return error response if not successful
			if err != nil {
				if err.Error() == models.RecordNotFoundErrorString {
					appCtx.Logger().Error("Error finding webhookSubscription to update")
					return webhooksubscriptionop.NewUpdateWebhookSubscriptionNotFound(), err
				}
				switch err.(type) {
				case apperror.PreconditionFailedError:
					appCtx.Logger().Error("Error updating webhookSubscription due to stale eTag")
					return webhooksubscriptionop.NewUpdateWebhookSubscriptionPreconditionFailed(), err
				default:
					appCtx.Logger().Error(fmt.Sprintf("Error updating webhookSubscription %s", params.WebhookSubscriptionID.String()), zap.Error(err))
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
			}

			// Convert model back to a payload and return to caller
			payload = payloads.WebhookSubscriptionPayload(*updatedWebhookSubscription)
			return webhooksubscriptionop.NewUpdateWebhookSubscriptionOK().WithPayload(payload), nil
		})
}
