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

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", "/offices", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := officeop.IndexOfficesParams{
			HTTPRequest: req,
		}

		handler := IndexOfficesHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeop.IndexOfficesOK{}, response)
		okResponse := response.(*officeop.IndexOfficesOK)
		suite.Len(okResponse.Payload, 0)
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := officeop.IndexOfficesParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		officeListFetcher := &mocks.OfficeListFetcher{}
		officeListFetcher.On("FetchOfficeList",
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexOfficesHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter: newQueryFilter,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
