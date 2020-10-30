package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	tsppop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_service_provider_performances"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/tsp"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexTSPPsHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
			ID: id,
		},
	}
	testdatagen.MakeTSPPerformance(suite.DB(), assertions)

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/transportation_service_provider_performances", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := tsppop.IndexTSPPsParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexTSPPsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
			TransportationServiceProviderPerformanceListFetcher: tsp.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
			NewPagination: pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&tsppop.IndexTSPPsOK{}, response)
		okResponse := response.(*tsppop.IndexTSPPsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		tspp := models.TransportationServiceProviderPerformance{ID: id}
		params := tsppop.IndexTSPPsParams{
			HTTPRequest: req,
		}
		ListFetcher := &mocks.TransportationServiceProviderPerformanceListFetcher{}
		ListFetcher.On("FetchTransportationServiceProviderPerformanceList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.TransportationServiceProviderPerformances{tspp}, nil).Once()
		ListFetcher.On("FetchTransportationServiceProviderPerformanceCount",
			mock.Anything,
		).Return(1, nil).Once()
		handler := IndexTSPPsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			TransportationServiceProviderPerformanceListFetcher: ListFetcher,
			NewPagination: pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&tsppop.IndexTSPPsOK{}, response)
		okResponse := response.(*tsppop.IndexTSPPsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := tsppop.IndexTSPPsParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		ListFetcher := &mocks.TransportationServiceProviderPerformanceListFetcher{}
		ListFetcher.On("FetchTransportationServiceProviderPerformanceList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexTSPPsHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			TransportationServiceProviderPerformanceListFetcher: ListFetcher,
			NewPagination: pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestGetTSPPHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
			ID: id,
		},
	}
	testdatagen.MakeTSPPerformance(suite.DB(), assertions)

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/transportation_service_provider_performances/"+uuidString, nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := tsppop.GetTSPPParams{
			HTTPRequest: req,
			TsppID:      *handlers.FmtUUID(id),
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := GetTSPPHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
			TransportationServiceProviderPerformanceFetcher: tsp.NewTransportationServiceProviderPerformanceFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&tsppop.GetTSPPOK{}, response)
		okResponse := response.(*tsppop.GetTSPPOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		tspp := models.TransportationServiceProviderPerformance{ID: id}
		params := tsppop.GetTSPPParams{
			HTTPRequest: req,
			TsppID:      *handlers.FmtUUID(id),
		}
		Fetcher := &mocks.TransportationServiceProviderPerformanceFetcher{}
		Fetcher.On("FetchTransportationServiceProviderPerformance",
			mock.Anything,
		).Return(tspp, nil).Once()
		handler := GetTSPPHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			TransportationServiceProviderPerformanceFetcher: Fetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&tsppop.GetTSPPOK{}, response)
		okResponse := response.(*tsppop.GetTSPPOK)
		suite.Equal(uuidString, okResponse.Payload.ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := tsppop.GetTSPPParams{
			HTTPRequest: req,
			TsppID:      *handlers.FmtUUID(id),
		}
		expectedError := models.ErrFetchNotFound
		Fetcher := &mocks.TransportationServiceProviderPerformanceFetcher{}
		Fetcher.On("FetchTransportationServiceProviderPerformance",
			mock.Anything,
		).Return(models.TransportationServiceProviderPerformance{}, expectedError).Once()
		handler := GetTSPPHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
			TransportationServiceProviderPerformanceFetcher: Fetcher,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestIndexTSPPsHandlerHelpers() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := IndexTSPPsHandler{
		HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewQueryFilter: query.NewQueryFilter,
		TransportationServiceProviderPerformanceListFetcher: tsp.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
		NewPagination: pagination.NewPagination,
	}

	suite.T().Run("test both filters present", func(t *testing.T) {

		s := `{"traffic_distribution_list_id":"001a4a1b-8b04-4621-b9ec-711d828f67e3", "transportation_service_provider_id":"8f166861-b8c4-4a8f-a43e-77ed5e745086"}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("traffic_distribution_list_id", "=", "001a4a1b-8b04-4621-b9ec-711d828f67e3"),
			query.NewQueryFilter("transportation_service_provider_id", "=", "8f166861-b8c4-4a8f-a43e-77ed5e745086"),
		}
		suite.Equal(expectedFilters, qfs)
	})
	suite.T().Run("test only traffic_distribution_list_id present", func(t *testing.T) {
		s := `{"traffic_distribution_list_id":"001a4a1b-8b04-4621-b9ec-711d828f67e3"}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("traffic_distribution_list_id", "=", "001a4a1b-8b04-4621-b9ec-711d828f67e3"),
		}
		suite.Equal(expectedFilters, qfs)
	})
	suite.T().Run("test only transportation_service_provider_id present", func(t *testing.T) {
		s := `{"transportation_service_provider_id":"8f166861-b8c4-4a8f-a43e-77ed5e745086"}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("transportation_service_provider_id", "=", "8f166861-b8c4-4a8f-a43e-77ed5e745086"),
		}
		suite.Equal(expectedFilters, qfs)
	})
}
