package adminapi

import (
	"net/http"

	"github.com/stretchr/testify/mock"

	notificationsop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/notification"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexNotificationsHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		notification0 := testdatagen.MakeDefaultNotification(suite.DB())
		testdatagen.MakeDefaultNotification(suite.DB())
		params := notificationsop.IndexNotificationsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/notifications"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexNotificationsHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&notificationsop.IndexNotificationsOK{}, response)
		okResponse := response.(*notificationsop.IndexNotificationsOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(notification0.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
		params := notificationsop.IndexNotificationsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/notifications"),
		}
		expectedError := models.ErrFetchNotFound
		listFetcher := &mocks.ListFetcher{}
		listFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		listFetcher.On("FetchRecordCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexNotificationsHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    listFetcher,
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
