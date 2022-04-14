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
	movehistory "github.com/transcom/mymove/pkg/services/move_history"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func getMoveHistoryForTest() models.MoveHistory {
	localUUID := uuid.Must(uuid.NewV4())
	transactionID := int64(3281)
	clientQuery := "UPDATE \"orders\" AS orders SET \"amended_orders_acknowledged_at\" = $1, \"department_indicator\" = $2, \"entitlement_id\" = $3, \"grade\" = $4, \"has_dependents\" = $5, \"issue_date\" = $6, \"new_duty_location_id\" = $7, \"nts_sac\" = $8, \"nts_tac\" = $9, \"orders_number\" = $10, \"orders_type\" = $11, \"orders_type_detail\" = $12, \"origin_duty_location_id\" = $13, \"report_by_date\" = $14, \"sac\" = $15, \"service_member_id\" = $16, \"spouse_has_pro_gear\" = $17, \"status\" = $18, \"tac\" = $19, \"updated_at\" = $20, \"uploaded_amended_orders_id\" = $21, \"uploaded_orders_id\" = $22 WHERE orders.id = $23"
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
				ClientQuery:     &clientQuery,
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

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)
	params := moveops.GetMoveHistoryParams{
		HTTPRequest: req,
		Locator:     "ABCD1234",
		Page:        swag.Int64(1),
		PerPage:     swag.Int64(20),
	}

	suite.T().Run("Successful move history fetch", func(t *testing.T) {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}

		handler := GetMoveHistoryHandler{
			HandlerContext:     handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		mockHistoryFetcher.On("FetchMoveHistory",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*services.FetchMoveHistoryParams"),
		).Return(&moveHistory, int64(1), nil)

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveHistoryOK{}, response)

		payload := response.(*moveops.GetMoveHistoryOK).Payload

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
		suite.Equal(maudit.ClientQuery, paudit.ClientQuery)
		suite.Equal(maudit.Action, paudit.Action)
		suite.Equal(maudit.EventName, paudit.EventName)
		suite.Equal(maudit.StatementOnly, paudit.StatementOnly)

		swaggerTimeFormat := "2006-01-02T15:04:05.99Z07:00"
		suite.Equal(maudit.ActionTstampTx.Format(swaggerTimeFormat), time.Time(paudit.ActionTstampTx).Format(swaggerTimeFormat))
		suite.Equal(maudit.ActionTstampStm.Format(swaggerTimeFormat), time.Time(paudit.ActionTstampStm).Format(swaggerTimeFormat))
		suite.Equal(maudit.ActionTstampClk.Format(swaggerTimeFormat), time.Time(paudit.ActionTstampClk).Format(swaggerTimeFormat))

	})

	suite.T().Run("Unsuccessful move history fetch - empty string bad request", func(t *testing.T) {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}

		handler := GetMoveHistoryHandler{
			HandlerContext:     handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		badParams := moveops.GetMoveHistoryParams{
			HTTPRequest: req,
			Locator:     "",
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(20),
		}
		response := handler.Handle(badParams)
		suite.IsType(&moveops.GetMoveHistoryBadRequest{}, response)
	})

	suite.T().Run("Unsuccessful move history fetch - locator not found", func(t *testing.T) {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}

		handler := GetMoveHistoryHandler{
			HandlerContext:     handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		mockHistoryFetcher.On("FetchMoveHistory",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*services.FetchMoveHistoryParams"),
		).Return(&models.MoveHistory{}, int64(0), apperror.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveHistoryNotFound{}, response)
	})

	suite.T().Run("Unsuccessful move history fetch - internal server error", func(t *testing.T) {
		mockHistoryFetcher := mocks.MoveHistoryFetcher{}

		handler := GetMoveHistoryHandler{
			HandlerContext:     handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			MoveHistoryFetcher: &mockHistoryFetcher,
		}

		mockHistoryFetcher.On("FetchMoveHistory",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*services.FetchMoveHistoryParams"),
		).Return(&models.MoveHistory{}, int64(0), apperror.QueryError{})

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveHistoryInternalServerError{}, response)
	})

	suite.T().Run("Paginated move history fetch results", func(t *testing.T) {
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

		context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		handler := GetMoveHistoryHandler{
			context,
			movehistory.NewMoveHistoryFetcher(),
		}

		response := handler.Handle(pagedParams)
		suite.IsNotErrResponse(response)

		suite.IsType(&moveops.GetMoveHistoryOK{}, response)
		payload := response.(*moveops.GetMoveHistoryOK).Payload

		// Returned row count of 2 (since page size = 2)
		suite.Len(payload.HistoryRecords, 2)
	})

}
