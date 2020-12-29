package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"

	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
)

func (suite *HandlerSuite) TestGetMoveHandler() {
	swaggerTimeFormat := "2006-01-02T15:04:05.999Z07:00"
	availableToPrimeAt := time.Now()
	submittedAt := availableToPrimeAt.Add(-1 * time.Hour)

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: &availableToPrimeAt,
			SubmittedAt:        &submittedAt,
		},
	})

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/move/#{move.locator}"), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)
	params := moveops.GetMoveParams{
		HTTPRequest: req,
		Locator:     move.Locator,
	}

	suite.T().Run("Successful move fetch", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MoveFetcher:    &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			move.Locator,
		).Return(&move, nil)

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)

		payload := response.(*moveops.GetMoveOK).Payload

		suite.Equal(move.ID.String(), payload.ID.String())
		suite.Equal(move.AvailableToPrimeAt.Format(swaggerTimeFormat), payload.AvailableToPrimeAt.String())
		suite.Equal(move.ContractorID.String(), payload.ContractorID.String())
		suite.Equal(move.Locator, payload.Locator)
		suite.Equal(move.OrdersID.String(), payload.OrdersID.String())
		suite.Equal(move.ReferenceID, payload.ReferenceID)
		suite.Equal(string(move.Status), string(payload.Status))
		suite.Equal(move.CreatedAt.Format(swaggerTimeFormat), payload.CreatedAt.String())
		suite.Equal(move.SubmittedAt.Format(swaggerTimeFormat), payload.SubmittedAt.String())
		suite.Equal(move.UpdatedAt.Format(swaggerTimeFormat), payload.UpdatedAt.String())
	})

	suite.T().Run("Unsuccessful move fetch - empty string bad request", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MoveFetcher:    &mockFetcher,
		}

		response := handler.Handle(moveops.GetMoveParams{HTTPRequest: req, Locator: ""})
		suite.IsType(&moveops.GetMoveBadRequest{}, response)
	})

	suite.T().Run("Unsuccessful move fetch - locator not found", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MoveFetcher:    &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			move.Locator,
		).Return(&models.Move{}, services.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
	})

	suite.T().Run("Unsuccessful move fetch - internal server error", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MoveFetcher:    &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			move.Locator,
		).Return(&models.Move{}, services.QueryError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveInternalServerError{}, response)
	})

}
