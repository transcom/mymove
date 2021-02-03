package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	webhooksubscriptionop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/webhook_subscriptions"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexWebhookSubscriptionsHandler() {
	// test that everything is wired up correctly
	m := testdatagen.MakeDefaultWebhookSubscription(suite.DB())
	req := httptest.NewRequest("GET", "/webhook_subscriptions", nil)

	suite.T().Run("200 - OK response", func(t *testing.T) {
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
