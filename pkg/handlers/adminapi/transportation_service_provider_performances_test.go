package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	tsppop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_service_provider_performance"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
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

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
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
			TransportationServiceProviderPerformanceListFetcher: tsppop.NewTransportationServiceProviderPerformanceListFetcher(queryBuilder),
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
		params := tsppop.IndexTSPPsHandler{
			HTTPRequest: req,
		}
		ListFetcher := &mocks.TransportationServiceProviderPerformanceFetcher{}
		ListFetcher.On("FetchTransportationServiceProviderPerformanceList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.TransportationServiceProviderPerformances{tspp}, nil).Once()
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
		params := tsppop.IndexTSPPsOK{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		ListFetcher := &mocks.TransportationServiceProviderPerformanceFetcher{}
		ListFetcher.On("FetchTransportationServiceProviderPerformanceList",
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
