package adminapi

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	officeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/office"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOfficesHandler() {
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		to := testdatagen.MakeDefaultTransportationOffice(suite.DB())
		params := officeop.IndexOfficesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/offices"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficesHandler{
			HandlerContext:    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter:    query.NewQueryFilter,
			OfficeListFetcher: office.NewOfficeListFetcher(queryBuilder),
			NewPagination:     pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeop.IndexOfficesOK{}, response)
		okResponse := response.(*officeop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(to.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("successful response", func() {
		id, _ := uuid.FromString("d874d002-5582-4a91-97d3-786e8f66c763")
		office := models.TransportationOffice{ID: id}
		params := officeop.IndexOfficesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/offices"),
		}
		officeListFetcher := &mocks.OfficeListFetcher{}
		officeListFetcher.On("FetchOfficeList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.TransportationOffices{office}, nil).Once()
		officeListFetcher.On("FetchOfficeCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(1, nil).Once()
		handler := IndexOfficesHandler{
			HandlerContext:    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter:    newQueryFilter,
			OfficeListFetcher: officeListFetcher,
			NewPagination:     pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeop.IndexOfficesOK{}, response)
		okResponse := response.(*officeop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(id.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		params := officeop.IndexOfficesParams{
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
			HandlerContext:    handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
