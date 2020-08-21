package webhook

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
	"github.com/transcom/mymove/cmd/webhook-client/utils/mocks"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// WebhookClientTestingSuite is a suite for testing the webhook client
type WebhookClientTestingSuite struct {
	testingsuite.PopTestSuite
	logger   utils.Logger
	certPath string
	keyPath  string
}

func TestWebhookClientTestingSuite(t *testing.T) {
	logger, _ := logging.Config("development", true)

	ts := &WebhookClientTestingSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
		certPath:     "../../config/tls/devlocal-mtls.cer",
		keyPath:      "../../config/tls/devlocal-mtls.key",
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *WebhookClientTestingSuite) Test_WebhookEngine() {

	// Create notification
	notification := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			Status:  models.WebhookNotificationSent,
			Payload: swag.String("{\"message\":\"This is an updated notification #1\"}"),
		},
	})
	//subscription := testdatagen.MakeDefaultWebhookSubscription(suite.DB())

	// type Engine struct {
	// 	Connection      *pop.Connection
	// 	Logger          utils.Logger
	// 	Client          *utils.WebhookRuntime
	// 	Cmd             *cobra.Command
	// 	PeriodInSeconds int
	// }

	// engine = Engine{
	// 	Connection: suite.DB(), //todo rename connection to dbconnection
	// 	Logger:     suite.logger,
	// }
	suite.logger.Info("Notification created",
		zap.String("eventKey", notification.EventKey),
		zap.String("status", string(notification.Status)),
		zap.String("payload", *notification.Payload),
	)
}

func (suite *WebhookClientTestingSuite) Test_SendOneNotification() {
	mockClient := mocks.WebhookRuntimeClient{}

	// Create the engine replacing the client with the mock client
	engine := Engine{
		Connection: suite.DB(),
		Logger:     suite.logger,
		Client:     &mockClient,
	}
	// Create a notification
	notification := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			Status:  models.WebhookNotificationSent,
			Payload: swag.String("{\"message\":\"This is an updated notification #1\"}"),
		},
	})
	// Create a subscription
	subscription := testdatagen.MakeWebhookSubscription(suite.DB(), testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "/my/callback/url",
		},
	})
	// Some return values
	var response = http.Response{}
	var body = "A nice message saying that the notification was received."
	var bodyBytes = []byte(body)

	// TESTCASE SCENARIO
	// What is being tested: sendOneNotification function
	// Behaviour: The function gets passed in 2 models, one for notification
	// and one for subscription. It does not touch the db.
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status in the model with SENT or FAILED accordingly

	suite.T().Run("Successful post, updated notification", func(t *testing.T) {

		// Expected behaviour: Response is set to 200, so model should be set to SENT
		response.StatusCode = 200
		response.Status = "200 Success"

		// Expectation: When Post is called, verify it was called with the callback url from the subscription.
		// Then, make it return 200 success and a body
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&response, bodyBytes, nil)

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was no error
		suite.Nil(err)
		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		// Check that notification Status was set to Sent
		suite.Equal(models.WebhookNotificationSent, notification.Status)
	})

	suite.T().Run("Failed post, updated notification", func(t *testing.T) {

		// Expected behaviour: Response is set to 400, so model should be set to FAILED
		response.StatusCode = 400
		response.Status = "400 Not Found Error"

		// Expectation: When Post is called, verify it was called with the callback url from the subscription.
		// Then, make it return 200 success and a body
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&response, bodyBytes, nil)

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was no error
		suite.NotNil(err)
		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		// Check that notification Status was set to Failed
		suite.Equal(models.WebhookNotificationFailed, notification.Status)
	})

	suite.T().Run("Failed post due to send error, updated notification", func(t *testing.T) {

		// Expected behaviour: Error was detected on sending, so model should be set to FAILED

		// Expectation: When Post is called, verify it was called with the callback url from the subscription.
		// Then, make it return 200 success and a body
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&response, bodyBytes, errors.New("Error due to server down"))

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was no error
		suite.NotNil(err)
		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		// Check that notification Status was set to Failed
		suite.Equal(models.WebhookNotificationFailed, notification.Status)
	})

}
