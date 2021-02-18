package adminapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	webhooksubscriptionop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/webhook_subscriptions"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	webhooksubscriptionservice "github.com/transcom/mymove/pkg/services/webhook_subscription"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexWebhookSubscriptionsHandler() {
	// test that everything is wired up correctly
	m := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	req := httptest.NewRequest("GET", "/webhook_subscriptions", nil)

	suite.T().Run("200 - OK response", func(t *testing.T) {
		// Setup: Provide a valid request to endpoint, when there is data in the db
		// Expected outcome:
		//   GET request returns 200 and a list of length 1 containing a subscription
		params := webhooksubscriptionop.IndexWebhookSubscriptionsParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexWebhookSubscriptionsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.IndexWebhookSubscriptionsOK{}, response)
		okResponse := response.(*webhooksubscriptionop.IndexWebhookSubscriptionsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(m.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.T().Run("404 - Move not found", func(t *testing.T) {
		// Mocked: Fetcher for handler
		// Setup: Provide a valid request to endpoint and mock fetcher
		// Expected outcome:
		//   GET request returns 404 and no records are returned
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := webhooksubscriptionop.IndexWebhookSubscriptionsParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		webhookSubscriptionListFetcher := &mocks.ListFetcher{}

		webhookSubscriptionListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()

		webhookSubscriptionListFetcher.On("FetchRecordCount",
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()

		handler := IndexWebhookSubscriptionsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    webhookSubscriptionListFetcher,
			NewPagination:  pagination.NewPagination,
		}
		response := handler.Handle(params)

		suite.CheckErrorResponse(response, http.StatusNotFound, expectedError.Error())
	})
}

func (suite *HandlerSuite) TestGetWebhookSubscriptionHandler() {
	// Setup: Create a default webhook subscription and request
	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID), nil)

	suite.T().Run("200 - OK, Successfuly get webhook subscription", func(t *testing.T) {
		// Under test: 			GetWebhookSubscriptionHandler, Fetcher
		// Set up: 				Provide a valid request with the id of a webhook_subscription
		// 		   					to the getWebhookSubscription endpoint.
		// Expected Outcome: 	The webhookSubscription is returned and we
		//					 		review a 200 OK.
		params := webhooksubscriptionop.GetWebhookSubscriptionParams{
			HTTPRequest:           req,
			WebhookSubscriptionID: strfmt.UUID(webhookSubscription.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := GetWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			webhooksubscriptionservice.NewWebhookSubscriptionFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.GetWebhookSubscriptionOK{}, response)
		okResponse := response.(*webhooksubscriptionop.GetWebhookSubscriptionOK)
		suite.Equal(webhookSubscription.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("404 - Not Found", func(t *testing.T) {
		// Under test: 			GetWebhookSubscriptionHandler
		// Mocks:				WebhookSubscriptionFetcher
		// Set up: 				Provide an invalid uuid to the getWebhookSubscription
		// 		   					endpoint and mock Fetcher to return an error.
		// Expected Outcome: 	The handler returns a 404.

		webhookSubscriptionFetcher := &mocks.WebhookSubscriptionFetcher{}
		fakeID, err := uuid.NewV4()
		suite.NoError(err)

		expectedError := models.ErrFetchNotFound
		params := webhooksubscriptionop.GetWebhookSubscriptionParams{
			HTTPRequest:           req,
			WebhookSubscriptionID: strfmt.UUID(fakeID.String()),
		}

		webhookSubscriptionFetcher.On("FetchWebhookSubscription",
			mock.Anything,
		).Return(models.WebhookSubscription{}, expectedError).Once()

		handler := GetWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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

func (suite *HandlerSuite) TestUpdateWebhookSubscriptionHandler() {
	// Setup: Create a default webhook subscription and request
	webhookSubscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/webhook_subscriptions/%s", webhookSubscription.ID), nil)

	suite.T().Run("200 - OK, Successfuly updated webhook subscription", func(t *testing.T) {
		// Testing: 			UdateWebhookSubscriptionHandler, Updater
		// Set up: 				Provide a valid request with the id of a webhook_subscription
		// 		   					to the updateWebhookSubscription endpoint.
		// Expected Outcome: 	The webhookSubscription is updated and we
		//					 		receive a 200 OK.
		params := webhooksubscriptionop.UpdateWebhookSubscriptionParams{
			HTTPRequest:           req,
			WebhookSubscriptionID: strfmt.UUID(webhookSubscription.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := UpdateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			webhooksubscriptionservice.NewWebhookSubscriptionUpdater(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&webhooksubscriptionop.UpdateWebhookSubscriptionOK{}, response)
		okResponse := response.(*webhooksubscriptionop.UpdateWebhookSubscriptionOK)
		suite.Equal(webhookSubscription.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("404 - Not Found", func(t *testing.T) {
		// Testing: 			UpdateWebhookSubscriptionHandler
		// Mocks:				WebhookSubscriptionUpdater
		// Set up: 				Provide an invalid uuid to the updateWebhookSubscription
		// 		   					endpoint and mock Updater to return an error.
		// Expected Outcome: 	The handler returns a 404.

		webhookSubscriptionUpdater := &mocks.WebhookSubscriptionUpdater{}
		fakeID, err := uuid.NewV4()
		suite.NoError(err)

		expectedError := models.ErrFetchNotFound
		params := webhooksubscriptionop.UpdateWebhookSubscriptionParams{
			HTTPRequest:           req,
			WebhookSubscriptionID: strfmt.UUID(fakeID.String()),
		}

		webhookSubscriptionUpdater.On("UpdateWebhookSubscription",
			mock.Anything,
		).Return(models.WebhookSubscription{}, expectedError).Once()

		handler := UpdateWebhookSubscriptionHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			webhookSubscriptionUpdater,
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
