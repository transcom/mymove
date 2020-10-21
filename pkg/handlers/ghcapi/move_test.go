package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
)

func (suite *HandlerSuite) TestGetMoveHandler() {
	move := testdatagen.MakeDefaultMove(suite.DB())

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/move/#{move.locator}"), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)
	params := moveops.GetMoveParams{
		HTTPRequest: req,
		Locator:     move.Locator,
	}

	suite.T().Run("Successful move fetch", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        &mockFetcher,
			NewQueryFilter: query.NewQueryFilter,
		}

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil)

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)
	})

	suite.T().Run("Unsuccessful move fetch - locator not found", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        &mockFetcher,
			NewQueryFilter: query.NewQueryFilter,
		}

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(errors.New("error Resource not found: e"))

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
	})

	suite.T().Run("Unsuccessful move fetch - internal server error", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			Fetcher:        &mockFetcher,
			NewQueryFilter: query.NewQueryFilter,
		}

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(errors.New("error"))

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveInternalServerError{}, response)
	})

}
