package publicapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	tsppop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/transportation_service_provider_performance"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func newMockQueryFilterBuilder(filter *mocks.QueryFilter) services.NewQueryFilter {
	return func(column string, comparator string, value interface{}) services.QueryFilter {
		return filter
	}
}

func (suite *HandlerSuite) TestLogTransportationServiceProviderPerformanceHandler() {
	tsppID, _ := uuid.NewV4()
	path := fmt.Sprintf("/transportation_service_provider_performances/%s", tsppID)
	noAuthReq := httptest.NewRequest("GET", path, nil)

	officeUserID, _ := uuid.NewV4()
	userID, _ := uuid.NewV4()
	officeUser := models.OfficeUser{ID: officeUserID, UserID: &userID}
	officeAuthReq := suite.AuthenticateOfficeRequest(noAuthReq, officeUser)

	officeUserParams := tsppop.LogTransportationServiceProviderPerformanceParams{
		HTTPRequest: officeAuthReq,
		TransportationServiceProviderPerformanceID: *handlers.FmtUUID(tsppID),
	}

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	tspp := models.TransportationServiceProviderPerformance{ID: tsppID}
	tsppFetcher := &mocks.TransportationServiceProviderPerformanceFetcher{}

	handler := LogTransportationServiceProviderPerformanceHandler{
		HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewQueryFilter: newQueryFilter,
		TransportationServiceProviderPerformanceFetcher: tsppFetcher,
	}

	suite.T().Run("has TSPP, no content response", func(t *testing.T) {
		tsppFetcher.On("FetchTransportationServiceProviderPerformance",
			mock.Anything,
		).Return(tspp, nil).Once()

		response := handler.Handle(officeUserParams)

		suite.IsType(&tsppop.LogTransportationServiceProviderPerformanceNoContent{}, response)
		// Note: No payload since it's just being logged.
	})

	suite.T().Run("does not have TSPP, no content response", func(t *testing.T) {
		expectedError := errors.New("sql: no rows in result set")
		tsppFetcher.On("FetchTransportationServiceProviderPerformance",
			mock.Anything,
		).Return(models.TransportationServiceProviderPerformance{}, expectedError).Once()

		response := handler.Handle(officeUserParams)

		suite.IsType(&tsppop.LogTransportationServiceProviderPerformanceNoContent{}, response)
		// Note: No payload since it's just being logged.
	})

	suite.T().Run("some other error", func(t *testing.T) {
		expectedError := errors.New("test error")
		tsppFetcher.On("FetchTransportationServiceProviderPerformance",
			mock.Anything,
		).Return(models.TransportationServiceProviderPerformance{}, expectedError).Once()

		response := handler.Handle(officeUserParams)

		suite.IsType(&tsppop.LogTransportationServiceProviderPerformanceInternalServerError{}, response)
	})

	suite.T().Run("use tsp user, get unauthorized response", func(t *testing.T) {
		// Create some params where auth'd as TSP to make sure we get an unauthorized response
		tspUserID, _ := uuid.NewV4()
		userID, _ = uuid.NewV4()
		tspUser := models.TspUser{ID: tspUserID, UserID: &userID}
		tspAuthReq := suite.AuthenticateTspRequest(noAuthReq, tspUser)

		tspUserParams := tsppop.LogTransportationServiceProviderPerformanceParams{
			HTTPRequest: tspAuthReq,
			TransportationServiceProviderPerformanceID: *handlers.FmtUUID(tsppID),
		}

		response := handler.Handle(tspUserParams)

		suite.IsType(&tsppop.LogTransportationServiceProviderPerformanceForbidden{}, response)
	})
}
