package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			ID: id,
		},
	}
	testdatagen.MakeTransportationOffice(suite.DB(), assertions)

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/offices", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := officeop.IndexOfficesParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexOfficesHandler{
			HandlerContext:    handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:    query.NewQueryFilter,
			OfficeListFetcher: office.NewOfficeListFetcher(queryBuilder),
			NewPagination:     pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeop.IndexOfficesOK{}, response)
		okResponse := response.(*officeop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		office := models.TransportationOffice{ID: id}
		params := officeop.IndexOfficesParams{
			HTTPRequest: req,
		}
		officeListFetcher := &mocks.OfficeListFetcher{}
		officeListFetcher.On("FetchOfficeList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.TransportationOffices{office}, nil).Once()
		officeListFetcher.On("FetchOfficeCount",
			mock.Anything,
		).Return(1, nil).Once()
		handler := IndexOfficesHandler{
			HandlerContext:    handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:    newQueryFilter,
			OfficeListFetcher: officeListFetcher,
			NewPagination:     pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeop.IndexOfficesOK{}, response)
		okResponse := response.(*officeop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := officeop.IndexOfficesParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		officeListFetcher := &mocks.OfficeListFetcher{}
		officeListFetcher.On("FetchOfficeList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		officeListFetcher.On("FetchOfficeCount",
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexOfficesHandler{
			HandlerContext:    handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
