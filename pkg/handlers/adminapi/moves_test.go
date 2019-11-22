package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/services/pagination"

	"github.com/transcom/mymove/pkg/services"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	moveop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexMovesHandler() {
	// test that everything is wired up correctly
	m := testdatagen.MakeDefaultMove(suite.DB())
	req := httptest.NewRequest("GET", "/moves", nil)

	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := moveop.IndexMovesParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexMovesHandler{
			HandlerContext:  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:  query.NewQueryFilter,
			MoveListFetcher: move.NewMoveListFetcher(queryBuilder),
			NewPagination:   pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&moveop.IndexMovesOK{}, response)
		okResponse := response.(*moveop.IndexMovesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(m.ID.String(), okResponse.Payload[0].ID.String())
	})
	suite.T().Run("test failed response", func(t *testing.T) {
		params := moveop.IndexMovesParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		moveListFetcher := &mocks.MoveListFetcher{}
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
		moveListFetcher.On("FetchMoveList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexMovesHandler{
			HandlerContext:  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:  newQueryFilter,
			MoveListFetcher: moveListFetcher,
			NewPagination:   pagination.NewPagination,
		}
		response := handler.Handle(params)

		suite.CheckErrorResponse(response, http.StatusNotFound, expectedError.Error())
	})
}

func (suite *HandlerSuite) TestIndexMovesHandlerHelpers() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := IndexMovesHandler{
		HandlerContext:  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewQueryFilter:  query.NewQueryFilter,
		MoveListFetcher: move.NewMoveListFetcher(queryBuilder),
		NewPagination:   pagination.NewPagination,
	}

	suite.T().Run("test filters present", func(t *testing.T) {

		s := `{"locator":"TEST123", "service_member_id":"TESTID}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("locator", "=", "TEST123"),
			query.NewQueryFilter("service_member_id", "=", "TESTID"),
		}
		suite.Equal(expectedFilters, qfs)
	})

	suite.T().Run("test locator filter present", func(t *testing.T) {

		s := `{"locator":"TEST123"}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("locator", "=", "TEST123"),
		}
		suite.Equal(expectedFilters, qfs)
	})

	suite.T().Run("test service_member_id filter present", func(t *testing.T) {

		s := `{"service_member_id":"TESTID"}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("service_member_id", "=", "TESTID"),
		}
		suite.Equal(expectedFilters, qfs)
	})
}
