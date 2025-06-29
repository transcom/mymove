package supportapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	webhookops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *HandlerSuite) TestCreateWebhookNotification() {

	suite.Run("Success createWebhookNotification 201 Created", func() {

		// TESTCASE SCENARIO
		// Under test: CreateWebhookNotificationHandler
		// Mocked:     None
		// Set up:     We create a webhook notification with a defined payload
		// Expected outcome:
		//             Success, a webhook notification with the fields we passed in

		// CREATE THE REQUEST
		request := httptest.NewRequest("POST", "/webhook-notifications/", nil)
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		request = request.WithContext(trace.NewContext(request.Context(), traceID))

		requestPayload := &supportmessages.WebhookNotification{
			EventKey: "Test.Create",
			Object:   models.StringPointer("{ \"message\": \"This is an example notification.\" } "),
		}
		params := webhookops.CreateWebhookNotificationParams{
			HTTPRequest: request,
			Body:        requestPayload,
		}

		handlerConfig := suite.NewHandlerConfig()
		handler := CreateWebhookNotificationHandler{handlerConfig}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsNotErrResponse(response)
		suite.IsType(webhookops.NewCreateWebhookNotificationCreated(), response)
		responsePayload := response.(*webhookops.CreateWebhookNotificationCreated).Payload

		suite.Equal(requestPayload.EventKey, responsePayload.EventKey)
		suite.Equal(supportmessages.WebhookNotificationStatusPENDING, responsePayload.Status)
		suite.Equal(requestPayload.Object, responsePayload.Object)
		suite.NotNil(responsePayload.ID)
		suite.NotNil(responsePayload.CreatedAt)
		suite.NotNil(responsePayload.UpdatedAt)
	})

	suite.Run("Success createWebhookNotification 201 Created from empty payload", func() {
		// TESTCASE SCENARIO
		// Under test: CreateWebhookNotificationHandler
		// Mocked:     None
		// Set up:     We create a webhook notification with an empty payload
		// Expected outcome:
		//             Success, A basic webhook notification with default fields is created

		// CREATE THE REQUEST
		request := httptest.NewRequest("POST", "/webhook-notifications/", nil)
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		request = request.WithContext(trace.NewContext(request.Context(), traceID))
		params := webhookops.CreateWebhookNotificationParams{
			HTTPRequest: request,
		}

		handlerConfig := suite.NewHandlerConfig()
		handler := CreateWebhookNotificationHandler{handlerConfig}

		// CALL FUNCTION UNDER TEST
		response := handler.Handle(params)

		// CHECK RESULTS
		suite.IsNotErrResponse(response)
		suite.IsType(webhookops.NewCreateWebhookNotificationCreated(), response)
		responsePayload := response.(*webhookops.CreateWebhookNotificationCreated).Payload

		suite.Equal(string(event.TestCreateEventKey), responsePayload.EventKey)
		suite.Equal(supportmessages.WebhookNotificationStatusPENDING, responsePayload.Status)
		suite.NotNil(responsePayload.Object)
		suite.NotNil(responsePayload.ID)
		suite.NotNil(responsePayload.CreatedAt)
		suite.NotNil(responsePayload.UpdatedAt)
	})

	suite.Run("Failed createWebhookNotification 500 Failed due to non-existent MTO id", func() {
		// TESTCASE SCENARIO
		// Under test: CreateWebhookNotificationHandler
		// Mocked:     None
		// Set up:     We create a webhook notification with a moveTaskOrderID that doesn't exist
		// Expected outcome:
		//             Fail, no notification is created.
		//             Note, returning 500 here because this is a support api.

		// CREATE THE REQUEST
		request := httptest.NewRequest("POST", "/webhook-notifications/", nil)
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		request = request.WithContext(trace.NewContext(request.Context(), traceID))

		moveTaskOrderID := uuid.Must(uuid.NewV4())
		requestPayload := &supportmessages.WebhookNotification{
			EventKey:        "Test.Create",
			Object:          models.StringPointer("{ \"message\": \"This is an example notification.\" } "),
			MoveTaskOrderID: handlers.FmtUUID(moveTaskOrderID),
		}
		params := webhookops.CreateWebhookNotificationParams{
			HTTPRequest: request,
			Body:        requestPayload,
		}

		handlerConfig := suite.NewHandlerConfig()
		handler := CreateWebhookNotificationHandler{handlerConfig}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESULTS
		suite.IsNotErrResponse(response)
		suite.IsType(webhookops.NewCreateWebhookNotificationInternalServerError(), response)

	})
}
