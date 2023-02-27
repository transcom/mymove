package adminapi

import (
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	transportation_officesop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/office"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestIndexOfficesHandler() {
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		to := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		params := transportation_officesop.IndexOfficesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/offices"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficesHandler{
			HandlerConfig:     suite.HandlerConfig(),
			NewQueryFilter:    query.NewQueryFilter,
			OfficeListFetcher: office.NewOfficeListFetcher(queryBuilder),
			NewPagination:     pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&transportation_officesop.IndexOfficesOK{}, response)
		okResponse := response.(*transportation_officesop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(to.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		params := transportation_officesop.IndexOfficesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/offices"),
		}
		expectedError := models.ErrFetchNotFound
		officeListFetcher := &mocks.OfficeListFetcher{}
		officeListFetcher.On("FetchOfficeList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		officeListFetcher.On("FetchOfficeCount",
			mock.AnythingOfType("*appcontext.appContext"),
		).Return(0, expectedError).Once()
		handler := IndexOfficesHandler{
			HandlerConfig:     suite.HandlerConfig(),
			NewQueryFilter:    newQueryFilter,
			OfficeListFetcher: officeListFetcher,
			NewPagination:     pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
