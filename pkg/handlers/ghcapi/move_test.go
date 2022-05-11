package ghcapi

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
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

	suite.T().Run("Successful move fetch", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveFetcher:    &mockFetcher,
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

	suite.T().Run("Unsuccessful move fetch - empty string bad request", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveFetcher:    &mockFetcher,
		}

		response := handler.Handle(moveops.GetMoveParams{HTTPRequest: req, Locator: ""})
		suite.IsType(&moveops.GetMoveBadRequest{}, response)
	})

	suite.T().Run("Unsuccessful move fetch - locator not found", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveFetcher:    &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
	})

	suite.T().Run("Unsuccessful move fetch - internal server error", func(t *testing.T) {
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveFetcher:    &mockFetcher,
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

	suite.T().Run("Successful flag setting to true", func(t *testing.T) {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	suite.T().Run("Unsuccessful flag - missing remarks", func(t *testing.T) {
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
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}

		response := handler.Handle(paramsNilRemarks)
		suite.IsType(&moveops.SetFinancialReviewFlagUnprocessableEntity{}, response)
	})
	suite.T().Run("Unsuccessful flag - move not found", func(t *testing.T) {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
	suite.T().Run("Unsuccessful flag - internal server error", func(t *testing.T) {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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

	suite.T().Run("Unsuccessful flag - bad etag", func(t *testing.T) {
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerContext:                handlers.NewHandlerContext(suite.DB(), suite.Logger()),
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
