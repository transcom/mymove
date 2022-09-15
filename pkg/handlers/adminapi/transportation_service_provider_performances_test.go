//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package adminapi

import (
	"fmt"
	"net/http"

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
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		tspp, err := testdatagen.MakeDefaultTSPPerformance(suite.DB())
		suite.Require().NoError(err)

		params := tsppop.IndexTSPPsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/transportation_service_provider_performances"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexTSPPsHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: query.NewQueryFilter,
			TransportationServiceProviderPerformanceListFetcher: tsp.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
			NewPagination: pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&tsppop.IndexTSPPsOK{}, response)
		okResponse := response.(*tsppop.IndexTSPPsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(tspp.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		params := tsppop.IndexTSPPsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/transportation_service_provider_performances"),
		}
		expectedError := models.ErrFetchNotFound
		ListFetcher := &mocks.TransportationServiceProviderPerformanceListFetcher{}
		ListFetcher.On("FetchTransportationServiceProviderPerformanceList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexTSPPsHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: newMockQueryFilterBuilder(&mocks.QueryFilter{}),
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
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		tspp, err := testdatagen.MakeDefaultTSPPerformance(suite.DB())
		suite.Require().NoError(err)

		params := tsppop.GetTSPPParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/transportation_service_provider_performances/%s", tspp.ID)),
			TsppID:      *handlers.FmtUUID(tspp.ID),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := GetTSPPHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: query.NewQueryFilter,
			TransportationServiceProviderPerformanceFetcher: tsp.NewTransportationServiceProviderPerformanceFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&tsppop.GetTSPPOK{}, response)
		okResponse := response.(*tsppop.GetTSPPOK)
		suite.Equal(tspp.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		tspp, err := testdatagen.MakeDefaultTSPPerformance(suite.DB())
		suite.Require().NoError(err)

		params := tsppop.GetTSPPParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/transportation_service_provider_performances/%s", tspp.ID)),
			TsppID:      *handlers.FmtUUID(tspp.ID),
		}
		expectedError := models.ErrFetchNotFound
		Fetcher := &mocks.TransportationServiceProviderPerformanceFetcher{}
		Fetcher.On("FetchTransportationServiceProviderPerformance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.TransportationServiceProviderPerformance{}, expectedError).Once()
		handler := GetTSPPHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: newMockQueryFilterBuilder(&mocks.QueryFilter{}),
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
	suite.Run("test both filters present", func() {

		s := `{"traffic_distribution_list_id":"001a4a1b-8b04-4621-b9ec-711d828f67e3", "transportation_service_provider_id":"8f166861-b8c4-4a8f-a43e-77ed5e745086"}`
		qfs := generateQueryFilters(suite.Logger(), &s, tsppFilterConverters)
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("traffic_distribution_list_id", "=", "001a4a1b-8b04-4621-b9ec-711d828f67e3"),
			query.NewQueryFilter("transportation_service_provider_id", "=", "8f166861-b8c4-4a8f-a43e-77ed5e745086"),
		}
		suite.ElementsMatch(expectedFilters, qfs) // order not important
	})
	suite.Run("test only traffic_distribution_list_id present", func() {
		s := `{"traffic_distribution_list_id":"001a4a1b-8b04-4621-b9ec-711d828f67e3"}`
		qfs := generateQueryFilters(suite.Logger(), &s, tsppFilterConverters)
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("traffic_distribution_list_id", "=", "001a4a1b-8b04-4621-b9ec-711d828f67e3"),
		}
		suite.Equal(expectedFilters, qfs)
	})
	suite.Run("test only transportation_service_provider_id present", func() {
		s := `{"transportation_service_provider_id":"8f166861-b8c4-4a8f-a43e-77ed5e745086"}`
		qfs := generateQueryFilters(suite.Logger(), &s, tsppFilterConverters)
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("transportation_service_provider_id", "=", "8f166861-b8c4-4a8f-a43e-77ed5e745086"),
		}
		suite.Equal(expectedFilters, qfs)
	})
}
