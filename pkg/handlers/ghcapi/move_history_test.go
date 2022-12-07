package ghcapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	movehistory "github.com/transcom/mymove/pkg/services/move_history"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func getMoveHistoryForTest() models.MoveHistory {
	localUUID := uuid.Must(uuid.NewV4())
	transactionID := int64(3281)
	eventName := "apiEndpoint"
	oldData := `{\"updated_at\": \"2022-03-08T19:08:44.664709\", \"postal_code\": \"90213\"}`
	changedData := `{\"updated_at\": \"2022-03-08T19:08:44.664709\", \"postal_code\": \"90213\"}`

	moveHistory := models.MoveHistory{
		ID:          uuid.Must(uuid.NewV4()),
		Locator:     "BILWEI",
		ReferenceID: handlers.FmtString("7858-9363"),
		AuditHistories: models.AuditHistories{
			{
				ID:              uuid.Must(uuid.NewV4()),
				SchemaName:      "",
				TableName:       "orders",
				RelID:           16879,
				ObjectID:        &localUUID,
				SessionUserID:   &localUUID,
				TransactionID:   &transactionID,
				Action:          "U",
				EventName:       &eventName,
				OldData:         &oldData,
				ChangedData:     &changedData,
				StatementOnly:   false,
				ActionTstampTx:  time.Now(),
				ActionTstampStm: time.Now(),
				ActionTstampClk: time.Now(),
			},
		},
	}
	return moveHistory
}

func (suite *HandlerSuite) TestMockGetMoveHistoryHandler() {
	moveHistory := getMoveHistoryForTest()

	suite.Run("Successful move history fetch", func() {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}
		requestUser := factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := moveops.GetMoveHistoryParams{
			HTTPRequest: req,
			Locator:     "ABCD1234",
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(20),
		}

		handler := GetMoveHistoryHandler{
			HandlerConfig:      suite.HandlerConfig(),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		mockHistoryFetcher.On("FetchMoveHistory",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*services.FetchMoveHistoryParams"),
		).Return(&moveHistory, int64(1), nil)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveHistoryOK{}, response)
		payload := response.(*moveops.GetMoveHistoryOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(moveHistory.ID.String(), payload.ID.String())
		suite.Equal(moveHistory.Locator, payload.Locator)
		suite.Equal(moveHistory.ReferenceID, payload.ReferenceID)

		suite.Equal(len(moveHistory.AuditHistories), len(payload.HistoryRecords))
		suite.Equal(1, len(payload.HistoryRecords))
		suite.Equal(len(moveHistory.AuditHistories), len(payload.HistoryRecords))
		maudit := moveHistory.AuditHistories[0]
		paudit := payload.HistoryRecords[0]
		suite.Equal(maudit.ID.String(), paudit.ID.String())
		suite.Equal(maudit.ObjectID.String(), paudit.ObjectID.String())
		suite.Equal(maudit.SessionUserID.String(), paudit.SessionUserID.String())
		suite.Equal(maudit.SchemaName, paudit.SchemaName)
		suite.Equal(maudit.TableName, paudit.TableName)
		suite.Equal(maudit.RelID, paudit.RelID)
		suite.Equal(maudit.Action, paudit.Action)
		suite.Equal(maudit.EventName, paudit.EventName)
		suite.Equal(maudit.StatementOnly, paudit.StatementOnly)
		swaggerTimeFormat := "2006-01-02T15:04:05.99Z07:00"
		suite.Equal(maudit.ActionTstampTx.Format(swaggerTimeFormat), time.Time(paudit.ActionTstampTx).Format(swaggerTimeFormat))
		suite.Equal(maudit.ActionTstampStm.Format(swaggerTimeFormat), time.Time(paudit.ActionTstampStm).Format(swaggerTimeFormat))
		suite.Equal(maudit.ActionTstampClk.Format(swaggerTimeFormat), time.Time(paudit.ActionTstampClk).Format(swaggerTimeFormat))

	})

	suite.Run("Unsuccessful move history fetch - empty string bad request", func() {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}
		requestUser := factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		handler := GetMoveHistoryHandler{
			HandlerConfig:      suite.HandlerConfig(),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		badParams := moveops.GetMoveHistoryParams{
			HTTPRequest: req,
			Locator:     "",
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(20),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(badParams)
		suite.IsType(&moveops.GetMoveHistoryBadRequest{}, response)
		payload := response.(*moveops.GetMoveHistoryBadRequest).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move history fetch - locator not found", func() {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}
		requestUser := factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := moveops.GetMoveHistoryParams{
			HTTPRequest: req,
			Locator:     "ABCD1234",
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(20),
		}

		handler := GetMoveHistoryHandler{
			HandlerConfig:      suite.HandlerConfig(),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		mockHistoryFetcher.On("FetchMoveHistory",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*services.FetchMoveHistoryParams"),
		).Return(&models.MoveHistory{}, int64(0), apperror.NotFoundError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveHistoryNotFound{}, response)
		payload := response.(*moveops.GetMoveHistoryNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move history fetch - internal server error", func() {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}
		requestUser := factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		params := moveops.GetMoveHistoryParams{
			HTTPRequest: req,
			Locator:     "ABCD1234",
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(20),
		}

		handler := GetMoveHistoryHandler{
			HandlerConfig:      suite.HandlerConfig(),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		mockHistoryFetcher.On("FetchMoveHistory",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*services.FetchMoveHistoryParams"),
		).Return(&models.MoveHistory{}, int64(0), apperror.QueryError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveHistoryInternalServerError{}, response)
		payload := response.(*moveops.GetMoveHistoryInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Paginated move history fetch results", func() {
		// Create a move
		move := testdatagen.MakeDefaultMove(suite.DB())

		// Add shipment to the move, giving the move some "history"
		shipment := models.MTOShipment{Status: models.MTOShipmentStatusSubmitted}
		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move:        move,
			MTOShipment: shipment,
		})

		// Build history request for a TIO user
		officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})
		request := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		pagedParams := moveops.GetMoveHistoryParams{
			HTTPRequest: request,
			Locator:     move.Locator,
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(2), // This should limit results to only the first 2 records of the possible 4
		}

		handlerConfig := suite.HandlerConfig()
		handler := GetMoveHistoryHandler{
			handlerConfig,
			movehistory.NewMoveHistoryFetcher(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(pagedParams)
		suite.IsNotErrResponse(response)

		suite.IsType(&moveops.GetMoveHistoryOK{}, response)
		payload := response.(*moveops.GetMoveHistoryOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		// Returned row count of 2 (since page size = 2)
		suite.Len(payload.HistoryRecords, 2)
	})
}
