package ghcapi

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservice "github.com/transcom/mymove/pkg/services/move"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetMoveHandler() {
	swaggerTimeFormat := "2006-01-02T15:04:05.99Z07:00"
	availableToPrimeAt := time.Now()
	submittedAt := availableToPrimeAt.Add(-1 * time.Hour)

	ordersID := uuid.Must(uuid.NewV4())
	var move models.Move
	var requestUser models.User
	setupTestData := func() {
		move = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
				SubmittedAt:        &submittedAt,
				Orders:             models.Order{ID: ordersID},
			},
		})
		requestUser = factory.BuildUser(nil, nil, nil)
	}

	suite.Run("Successful move fetch", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&move, nil)

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)
		payload := response.(*moveops.GetMoveOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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
		suite.Nil(payload.CloseoutOffice)
	})

	suite.Run("Successful move with a saved transportation office", func() {
		transportationOffice := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				ProvidesCloseout: true,
			},
		})

		move = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:           models.MoveStatusSUBMITTED,
				SubmittedAt:      &submittedAt,
				CloseoutOffice:   &transportationOffice,
				CloseoutOfficeID: &transportationOffice.ID,
			},
		})
		moveFetcher := moveservice.NewMoveFetcher()
		requestOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor})

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateOfficeRequest(req, requestOfficeUser)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   moveFetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)
		payload := response.(*moveops.GetMoveOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(transportationOffice.ID.String(), payload.CloseoutOfficeID.String())
		suite.Equal(transportationOffice.ID.String(), payload.CloseoutOffice.ID.String())
		suite.Equal(transportationOffice.AddressID.String(), payload.CloseoutOffice.Address.ID.String())

	})

	suite.Run("Unsuccessful move fetch - empty string bad request", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
		}
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		// Validate incoming payload: no body to validate

		response := handler.Handle(moveops.GetMoveParams{HTTPRequest: req, Locator: ""})
		suite.IsType(&moveops.GetMoveBadRequest{}, response)
		payload := response.(*moveops.GetMoveBadRequest).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move fetch - locator not found", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.NotFoundError{})
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
		payload := response.(*moveops.GetMoveNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move fetch - internal server error", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.QueryError{})

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveInternalServerError{}, response)
		payload := response.(*moveops.GetMoveInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestSearchMovesHandler() {

	var requestUser models.User
	setupTestData := func() *http.Request {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		return req

	}

	suite.Run("Successful move search by locator", func() {
		req := setupTestData()
		move := testdatagen.MakeDefaultMove(suite.DB())
		moves := models.Moves{move}

		mockSearcher := mocks.MoveSearcher{}

		handler := SearchMovesHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveSearcher:  &mockSearcher,
		}
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.Locator == move.Locator
			}),
		).Return(moves, 1, nil)

		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: &move.Locator,
				DodID:   nil,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload := response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		payloadMove := *(*payload).SearchMoves[0]
		suite.Equal(move.ID.String(), payloadMove.ID.String())
		suite.Equal(*move.Orders.ServiceMember.Edipi, *payloadMove.DodID)
		suite.Equal(move.Orders.NewDutyLocation.Address.PostalCode, payloadMove.DestinationDutyLocationPostalCode)
		suite.Equal(move.Orders.OriginDutyLocation.Address.PostalCode, payloadMove.OriginDutyLocationPostalCode)
		suite.Equal(ghcmessages.MoveStatusDRAFT, payloadMove.Status)
		suite.Equal("ARMY", payloadMove.Branch)
		suite.Equal(int64(0), payloadMove.ShipmentsCount)
		suite.NotEmpty(payloadMove.FirstName)
		suite.NotEmpty(payloadMove.LastName)
	})

	suite.Run("Successful move search by DoD ID", func() {
		req := setupTestData()
		move := testdatagen.MakeDefaultMove(suite.DB())
		moves := models.Moves{move}

		mockSearcher := mocks.MoveSearcher{}

		handler := SearchMovesHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveSearcher:  &mockSearcher,
		}
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.DodID == *move.Orders.ServiceMember.Edipi &&
					params.Locator == nil &&
					params.CustomerName == nil
			}),
		).Return(moves, 1, nil)

		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: nil,
				DodID:   move.Orders.ServiceMember.Edipi,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload := response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(move.ID.String(), (*payload).SearchMoves[0].ID.String())
	})
}

func (suite *HandlerSuite) TestSetFinancialReviewFlagHandler() {
	var move models.Move
	var requestUser models.User
	setupTestData := func() (*http.Request, models.Move) {
		move = testdatagen.MakeDefaultMove(suite.DB())
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		return req, move
	}
	defaultRemarks := "destination address is on the moon"
	fakeEtag := ""

	suite.Run("Successful flag setting to true", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			move.ID,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&move, nil)

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: swag.Bool(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagOK{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful flag - missing remarks", func() {
		req, move := setupTestData()
		paramsNilRemarks := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       nil,
				FlagForReview: swag.Bool(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}

		// Validate incoming payload
		suite.NoError(paramsNilRemarks.Body.Validate(strfmt.Default))

		response := handler.Handle(paramsNilRemarks)
		suite.IsType(&moveops.SetFinancialReviewFlagUnprocessableEntity{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful flag - move not found", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.NotFoundError{})

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: swag.Bool(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagNotFound{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful flag - internal server error", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.QueryError{})

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: swag.Bool(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagInternalServerError{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful flag - bad etag", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.PreconditionFailedError{})

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: swag.Bool(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagPreconditionFailed{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagPreconditionFailed).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestUpdateMoveCloseoutOfficeHandler() {
	var move models.Move
	var requestUser models.OfficeUser
	var transportationOffice models.TransportationOffice

	closeoutOfficeUpdater := moveservice.NewCloseoutOfficeUpdater(moveservice.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())

	setupTestData := func() (*http.Request, models.Move, models.TransportationOffice) {
		move = testdatagen.MakeDefaultMove(suite.DB())
		requestUser = factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor})
		transportationOffice = testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				ProvidesCloseout: true,
			},
		})

		req := httptest.NewRequest("GET", "/move/#{move.locator}/closeout-office", nil)
		req = suite.AuthenticateOfficeRequest(req, requestUser)
		return req, move, transportationOffice
	}

	suite.Run("Successful update of closeout office", func() {
		req, move, transportationOffice := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: move.Locator,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficeOK{}, response)
		payload := response.(*moveops.UpdateCloseoutOfficeOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(closeoutOfficeID, *payload.CloseoutOfficeID)
		suite.Equal(closeoutOfficeID, *payload.CloseoutOffice.ID)
		suite.Equal(transportationOffice.AddressID.String(), payload.CloseoutOffice.Address.ID.String())
	})

	suite.Run("Unsuccessful move not found", func() {
		req, move, transportationOffice := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: "ABC123",
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficeNotFound{}, response)
	})

	suite.Run("Unsuccessful closeout office not found", func() {
		transportationOfficeNonCloseout := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				ProvidesCloseout: false,
			},
		})

		req, move, _ := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOfficeNonCloseout.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: move.Locator,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficeNotFound{}, response)
	})

	suite.Run("Unsuccessful eTag does not match", func() {
		req, move, transportationOffice := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(time.Now()),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: move.Locator,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficePreconditionFailed{}, response)
	})
}
