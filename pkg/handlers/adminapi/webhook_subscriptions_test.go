package adminapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	webhooksubscriptionop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/webhook_subscriptions"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	webhooksubscription "github.com/transcom/mymove/pkg/services/webhook_subscription"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexWebhookSubscriptionsHandler() {
	suite.Run("200 - OK response", func() {
		// Setup: Provide a valid request to endpoint, when there is data in the db
		// Expected outcome:
		//   GET request returns 200 and a list of length 1 containing a subscription
		ws := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		params := webhooksubscriptionop.IndexWebhookSubscriptionsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/webhook_subscriptions"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexWebhookSubscriptionsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.IndexWebhookSubscriptionsOK{}, response)
		okResponse := response.(*webhooksubscriptionop.IndexWebhookSubscriptionsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(ws.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("404 - Move not found", func() {
		// Mocked: Fetcher for handler
		// Setup: Provide a valid request to endpoint and mock fetcher
		// Expected outcome:
		//   GET request returns 404 and no records are returned
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := webhooksubscriptionop.IndexWebhookSubscriptionsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/webhook_subscriptions"),
		}
		expectedError := models.ErrFetchNotFound
		webhookSubscriptionListFetcher := &mocks.ListFetcher{}

		webhookSubscriptionListFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()

		webhookSubscriptionListFetcher.On("FetchRecordCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()

		handler := IndexWebhookSubscriptionsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    webhookSubscriptionListFetcher,
			NewPagination:  pagination.NewPagination,
		}
		response := handler.Handle(params)

		suite.CheckErrorResponse(response, http.StatusNotFound, expectedError.Error())
	})
}

func (suite *HandlerSuite) TestGetWebhookSubscriptionHandler() {
	suite.Run("200 - OK, Successfuly get webhook subscription", func() {
		// Under test: 			GetWebhookSubscriptionHandler, Fetcher
		// Set up: 				Provide a valid request with the id of a webhook_subscription
		// 		   					to the getWebhookSubscription endpoint.
		// Expected Outcome: 	The webhookSubscription is returned and we
		//					 		review a 200 OK.
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		params := webhooksubscriptionop.GetWebhookSubscriptionParams{
			HTTPRequest:           suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID)),
			WebhookSubscriptionID: strfmt.UUID(webhookSubscription.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhooksubscription.NewWebhookSubscriptionFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.GetWebhookSubscriptionOK{}, response)
		okResponse := response.(*webhooksubscriptionop.GetWebhookSubscriptionOK)
		suite.Equal(webhookSubscription.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("404 - Not Found", func() {
		// Under test: 			GetWebhookSubscriptionHandler
		// Mocks:				WebhookSubscriptionFetcher
		// Set up: 				Provide an invalid uuid to the getWebhookSubscription
		// 		   					endpoint and mock Fetcher to return an error.
		// Expected Outcome: 	The handler returns a 404.
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		webhookSubscriptionFetcher := &mocks.WebhookSubscriptionFetcher{}
		fakeID, err := uuid.NewV4()
		suite.NoError(err)

		expectedError := models.ErrFetchNotFound
		params := webhooksubscriptionop.GetWebhookSubscriptionParams{
			HTTPRequest:           suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID)),
			WebhookSubscriptionID: strfmt.UUID(fakeID.String()),
		}

		webhookSubscriptionFetcher.On("FetchWebhookSubscription",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.WebhookSubscription{}, expectedError).Once()

		handler := GetWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhookSubscriptionFetcher,
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestCreateWebhookSubscriptionHandler() {
	setupHandler := func() CreateWebhookSubscriptionHandler {
		return CreateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhooksubscription.NewWebhookSubscriptionCreator(query.NewQueryBuilder()),
			query.NewQueryFilter,
		}
	}

	// Base WebhookSubscription setup
	webhookSubscription := models.WebhookSubscription{
		Status:      models.WebhookSubscriptionStatusActive,
		EventKey:    "PaymentRequest.Update",
		CallbackURL: "/my/callback/url",
	}
	status := adminmessages.WebhookSubscriptionStatus(webhookSubscription.Status)

	// Actually test handler and creator on test database,
	suite.Run("201 - Successful create", func() {
		subscriber := testdatagen.MakeDefaultContractor(suite.DB())
		params := webhooksubscriptionop.CreateWebhookSubscriptionParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/webhook_subscriptions"),
			WebhookSubscription: &adminmessages.CreateWebhookSubscription{
				Status:       &status,
				EventKey:     &webhookSubscription.EventKey,
				SubscriberID: handlers.FmtUUID(subscriber.ID),
				CallbackURL:  &webhookSubscription.CallbackURL,
			},
		}

		response := setupHandler().Handle(params)
		suite.IsType(&webhooksubscriptionop.CreateWebhookSubscriptionCreated{}, response)

		subscriptionCreated := response.(*webhooksubscriptionop.CreateWebhookSubscriptionCreated)
		suite.NotEqual(subscriptionCreated.Payload.ID.String(), "00000000-0000-0000-0000-000000000000")
	})

	suite.Run("400 - Invalid Request", func() {
		// Create Fake subscriber
		fakeSubscriberID, _ := uuid.NewV4()
		params := webhooksubscriptionop.CreateWebhookSubscriptionParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/webhook_subscriptions"),
			WebhookSubscription: &adminmessages.CreateWebhookSubscription{
				Status:       &status,
				EventKey:     &webhookSubscription.EventKey,
				SubscriberID: handlers.FmtUUID(fakeSubscriberID),
				CallbackURL:  &webhookSubscription.CallbackURL,
			},
		}

		response := setupHandler().Handle(params)
		suite.IsType(webhooksubscriptionop.NewCreateWebhookSubscriptionBadRequest(), response)
	})
}

func (suite *HandlerSuite) TestUpdateWebhookSubscriptionHandler() {
	suite.Run("200 - OK, Successfully updated webhook subscription", func() {
		// Testing:             UpdateWebhookSubscriptionHandler, WebhookSubscriptionUpdater
		// Set up:              Provide a valid request with the id of an existing webhook_subscription
		//                      to the updateWebhookSubscription endpoint.
		// Expected Outcome:    The webhookSubscription is updated and we receive a 200 OK.
		//                      Fields are changed as expected.
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		status := adminmessages.WebhookSubscriptionStatusFAILING
		subscriberID := strfmt.UUID(webhookSubscription.SubscriberID.String())

		params := webhooksubscriptionop.UpdateWebhookSubscriptionParams{
			HTTPRequest:           suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID)),
			WebhookSubscriptionID: strfmt.UUID(webhookSubscription.ID.String()),
			WebhookSubscription: &adminmessages.WebhookSubscription{
				CallbackURL:  swag.String("something.com"),
				Status:       &status,
				EventKey:     swag.String("WebhookSubscription.Update"),
				SubscriberID: &subscriberID,
			},
			IfMatch: etag.GenerateEtag(webhookSubscription.UpdatedAt),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := UpdateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhooksubscription.NewWebhookSubscriptionUpdater(queryBuilder),
			query.NewQueryFilter,
		}

		suite.NoError(params.WebhookSubscription.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.UpdateWebhookSubscriptionOK{}, response)
		okResponse := response.(*webhooksubscriptionop.UpdateWebhookSubscriptionOK)
		suite.Equal(params.WebhookSubscription.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(params.WebhookSubscription.CallbackURL, okResponse.Payload.CallbackURL)
		suite.Equal(params.WebhookSubscription.EventKey, okResponse.Payload.EventKey)
		suite.Equal(params.WebhookSubscription.Status, okResponse.Payload.Status)
		suite.Equal(params.WebhookSubscription.SubscriberID, okResponse.Payload.SubscriberID)
	})

	suite.Run("200 - OK, Successfully partial updated webhook subscription", func() {
		// Testing:           UpdateWebhookSubscriptionHandler, WebhookSubscriptionUpdater
		// Set up:            Provide a valid request with the id of a webhook_subscription
		//                    to the updateWebhookSubscription endpoint and only some of the fields
		// Expected Outcome:  The webhookSubscription is updated and we
		//                    receive a 200 OK. Updated fields have changed, and others have not.
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())

		params := webhooksubscriptionop.UpdateWebhookSubscriptionParams{
			HTTPRequest:           suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID)),
			WebhookSubscriptionID: strfmt.UUID(webhookSubscription.ID.String()),
			WebhookSubscription: &adminmessages.WebhookSubscription{
				CallbackURL: swag.String("somethingelse.com"),
				EventKey:    swag.String("WebhookSubscription.Delete"),
			},
			IfMatch: etag.GenerateEtag(webhookSubscription.UpdatedAt),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := UpdateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhooksubscription.NewWebhookSubscriptionUpdater(queryBuilder),
			query.NewQueryFilter,
		}

		// Run swagger validations
		suite.NoError(params.WebhookSubscription.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.UpdateWebhookSubscriptionOK{}, response)
		okResponse := response.(*webhooksubscriptionop.UpdateWebhookSubscriptionOK)
		// These fields should have changed
		suite.Equal(params.WebhookSubscription.CallbackURL, okResponse.Payload.CallbackURL)
		suite.Equal(params.WebhookSubscription.EventKey, okResponse.Payload.EventKey)
		// These fields should have not changed
		suite.Equal(webhookSubscription.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(webhookSubscription.Status, models.WebhookSubscriptionStatus(*okResponse.Payload.Status))
		suite.Equal(webhookSubscription.SubscriberID, uuid.FromStringOrNil(okResponse.Payload.SubscriberID.String()))
	})

	suite.Run("404 - Not Found", func() {
		// Testing:           UpdateWebhookSubscriptionHandler, WebhookSubscriptionUpdater
		// Set up:            Provide a valid request with the wrong ID
		//                    to the updateWebhookSubscription endpoint
		// Expected Outcome:  We receive a 404 Not Found error.
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		fakeID, err := uuid.NewV4()
		suite.NoError(err)
		status := adminmessages.WebhookSubscriptionStatusFAILING
		subscriberID := strfmt.UUID(webhookSubscription.SubscriberID.String())
		params := webhooksubscriptionop.UpdateWebhookSubscriptionParams{
			HTTPRequest:           suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID)),
			WebhookSubscriptionID: strfmt.UUID(fakeID.String()),
			WebhookSubscription: &adminmessages.WebhookSubscription{
				CallbackURL:  swag.String("something.com"),
				Status:       &status,
				EventKey:     swag.String("WebhookSubscription.Update"),
				SubscriberID: &subscriberID,
			},
			IfMatch: etag.GenerateEtag(webhookSubscription.UpdatedAt),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := UpdateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhooksubscription.NewWebhookSubscriptionUpdater(queryBuilder),
			query.NewQueryFilter,
		}

		suite.NoError(params.WebhookSubscription.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&webhooksubscriptionop.UpdateWebhookSubscriptionNotFound{}, response)
	})

	suite.Run("412 - Precondition Failed", func() {
		// Testing:           UpdateWebhookSubscriptionHandler, WebhookSubscriptionUpdater
		// Set up:            Provide a valid request with a stale ETag
		//                    to the updateWebhookSubscription endpoint
		// Expected Outcome:  We receive a 412 Precondition Failed error.
		webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
		params := webhooksubscriptionop.UpdateWebhookSubscriptionParams{
			HTTPRequest:           suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID)),
			WebhookSubscriptionID: strfmt.UUID(webhookSubscription.ID.String()),
			WebhookSubscription: &adminmessages.WebhookSubscription{
				CallbackURL: swag.String("stale.etag.com"),
			},
			IfMatch: etag.GenerateEtag(time.Now()),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := UpdateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			webhooksubscription.NewWebhookSubscriptionUpdater(queryBuilder),
			query.NewQueryFilter,
		}

		suite.NoError(params.WebhookSubscription.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&webhooksubscriptionop.UpdateWebhookSubscriptionPreconditionFailed{}, response)
	})
}
