package ghcapi

import (
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"

	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
)

func (suite *HandlerSuite) TestGetMoveHandler() {
	swaggerTimeFormat := "2006-01-02T15:04:05.99Z07:00"
	availableToPrimeAt := time.Now()
	submittedAt := availableToPrimeAt.Add(-1 * time.Hour)

	ordersID := uuid.Must(uuid.NewV4())
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: &availableToPrimeAt,
			SubmittedAt:        &submittedAt,
			Orders:             models.Order{ID: ordersID},
		},
	})

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)
	params := moveops.GetMoveParams{
		HTTPRequest: req,
		Locator:     move.Locator,
	}

	suite.Run("Successful move fetch", func() {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFetcher:   &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&move, nil)

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)

		payload := response.(*moveops.GetMoveOK).Payload

		suite.Equal(move.ID.String(), payload.ID.String())
		suite.Equal(move.AvailableToPrimeAt.Format(swaggerTimeFormat), time.Time(*payload.AvailableToPrimeAt).Format(swaggerTimeFormat))
		suite.Equal(move.ContractorID.String(), payload.ContractorID.String())
		suite.Equal(move.Locator, payload.Locator)
		suite.Equal(move.OrdersID.String(), payload.OrdersID.String())
		suite.Equal(move.ReferenceID, payload.ReferenceID)
		suite.Equal(string(move.Status), string(payload.Status))
		suite.Equal(move.CreatedAt.Format(swaggerTimeFormat), time.Time(payload.CreatedAt).Format(swaggerTimeFormat))
		suite.Equal(move.SubmittedAt.Format(swaggerTimeFormat), time.Time(*payload.SubmittedAt).Format(swaggerTimeFormat))
		suite.Equal(move.UpdatedAt.Format(swaggerTimeFormat), time.Time(payload.UpdatedAt).Format(swaggerTimeFormat))
		suite.Equal(ordersID, move.Orders.ID)
	})

	suite.Run("Unsuccessful move fetch - empty string bad request", func() {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFetcher:   &mockFetcher,
		}

		response := handler.Handle(moveops.GetMoveParams{HTTPRequest: req, Locator: ""})
		suite.IsType(&moveops.GetMoveBadRequest{}, response)
	})

	suite.Run("Unsuccessful move fetch - locator not found", func() {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFetcher:   &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
	})

	suite.Run("Unsuccessful move fetch - internal server error", func() {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFetcher:   &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.QueryError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveInternalServerError{}, response)
	})

}

func (suite *HandlerSuite) TestSearchMovesHandler() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	moves := make(models.Moves, 1)
	moves[0] = move
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	suite.T().Run("Successful move search by locator", func(t *testing.T) {
		mockSearcher := mocks.MoveSearcher{}

		handler := SearchMovesHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveSearcher:  &mockSearcher,
		}

		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			&move.Locator,
			mock.Anything,
		).Return(moves, nil)

		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: &move.Locator,
				DodID:   nil,
			},
		}

		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)

		payload := response.(*moveops.SearchMovesOK).Payload

		payloadMove := *(*payload).SearchMoves[0]
		suite.Equal(move.ID.String(), payloadMove.ID.String())
		suite.Equal(*move.Orders.ServiceMember.Edipi, payloadMove.Customer.DodID)
		suite.Equal(move.Orders.NewDutyLocation.Address.PostalCode, *payloadMove.DestinationDutyLocation.Address.PostalCode)
		suite.Equal(move.Orders.OriginDutyLocation.Address.PostalCode, *payloadMove.OriginDutyLocation.Address.PostalCode)
		suite.Equal(ghcmessages.MoveStatusDRAFT, payloadMove.Status)
		suite.Equal("ARMY", payloadMove.Customer.Agency)
		suite.Equal(int64(0), payloadMove.ShipmentsCount)
		suite.NotEmpty(payloadMove.Customer.FirstName)
		suite.NotEmpty(payloadMove.Customer.LastName)
	})

	suite.T().Run("Successful move search by DoD ID", func(t *testing.T) {
		mockSearcher := mocks.MoveSearcher{}

		handler := SearchMovesHandler{
			HandlerConfig: handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveSearcher:  &mockSearcher,
		}

		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			move.Orders.ServiceMember.Edipi,
		).Return(moves, nil)

		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: nil,
				DodID:   move.Orders.ServiceMember.Edipi,
			},
		}
		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)

		payload := response.(*moveops.SearchMovesOK).Payload

		suite.Equal(move.ID.String(), (*payload).SearchMoves[0].ID.String())
	})
}

func (suite *HandlerSuite) TestSetFinancialReviewFlagHandler() {
	move := testdatagen.MakeDefaultMove(suite.DB())

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)
	defaultRemarks := "destination address is on the moon"
	fakeEtag := ""
	params := moveops.SetFinancialReviewFlagParams{
		HTTPRequest: req,
		IfMatch:     &fakeEtag,
		Body: moveops.SetFinancialReviewFlagBody{
			Remarks:       &defaultRemarks,
			FlagForReview: swag.Bool(true),
		},
		MoveID: *handlers.FmtUUID(move.ID),
	}

	suite.Run("Successful flag setting to true", func() {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			move.ID,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&move, nil)

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagOK{}, response)
	})

	suite.Run("Unsuccessful flag - missing remarks", func() {
		paramsNilRemarks := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks: nil,
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}

		response := handler.Handle(paramsNilRemarks)
		suite.IsType(&moveops.SetFinancialReviewFlagUnprocessableEntity{}, response)
	})
	suite.Run("Unsuccessful flag - move not found", func() {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagNotFound{}, response)
	})
	suite.Run("Unsuccessful flag - internal server error", func() {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.QueryError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagInternalServerError{}, response)
	})

	suite.Run("Unsuccessful flag - bad etag", func() {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.PreconditionFailedError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagPreconditionFailed{}, response)
	})
}
