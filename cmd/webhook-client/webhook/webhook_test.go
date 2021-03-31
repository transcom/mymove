//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate test data for use in the unit test
//RA: Creation of test data generation for unit test consumption does not present any unexpected states and conditions
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

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
	logger, _ := logging.Config(logging.WithEnvironment("development"), logging.WithLoggingLevel("debug"))

	ts := &WebhookClientTestingSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
		certPath:     "../../config/tls/devlocal-mtls.cer",
		keyPath:      "../../config/tls/devlocal-mtls.key",
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *WebhookClientTestingSuite) Test_SendStgNotification() {
	defer teardownEngineRun(suite)

	// Parse flags from environment
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create a client
	client, _, err := utils.CreateClient(v)
	suite.Nil(err)

	// Create the engine
	engine := Engine{
		DB:                  suite.DB(),
		Logger:              suite.logger,
		Client:              client,
		MaxImmediateRetries: 3,
	}
	// Create a notification
	notification := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			Status:  models.WebhookNotificationPending,
			Payload: "{\"message\":\"This is an updated notification #1\"}",
		},
	})
	// Create a subscription
	subscription := testdatagen.MakeWebhookSubscription(suite.DB(), testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "https://api.stg.move.mil/support/v1/webhook-notify",
		},
	})

	// TESTCASE SCENARIO
	// What is being tested: sendOneNotification function
	// Mocked: None
	// Behaviour: The function gets passed in 2 models, one for notification
	// and one for subscription.
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status with SENT or FAILED accordingly

	suite.T().Run("Successful post to staging", func(t *testing.T) {

		// Under test: sendOneNotification function
		// Set up:     We provide a PENDING webhook notification, and point the
		//             subscription at live Staging environment
		// Expected outcome:
		//             Notification would be updated as SENT

		// Call the engine function.
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was no error
		suite.Nil(err)
		// Check that notification Status was set to Sent in the model
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationSent, notif.Status)
		// Check that first attempted at date was set
		suite.NotNil(notif.FirstAttemptedAt)

	})

}

func (suite *WebhookClientTestingSuite) Test_SendOneNotification() {
	mockClient := mocks.WebhookRuntimeClient{}
	defer teardownEngineRun(suite)

	// Create the engine replacing the client with the mock client
	engine := Engine{
		DB:                  suite.DB(),
		Logger:              suite.logger,
		Client:              &mockClient,
		MaxImmediateRetries: 3,
	}
	// Create a notification
	notification := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			Status:  models.WebhookNotificationSent,
			Payload: "{\"message\":\"This is an updated notification #1\"}",
		},
	})
	// Create a subscription
	subscription := testdatagen.MakeWebhookSubscription(suite.DB(), testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "/my/callback/url",
		},
	})
	// Some return values
	var responseSuccess = http.Response{}
	responseSuccess.StatusCode = 200
	responseSuccess.Status = "200 Success"

	// Create 400 response
	var responseFail = http.Response{}
	responseFail.StatusCode = 400
	responseFail.Status = "400 Not Found Error"

	var body = "A nice message saying that the notification was received."
	var bodyBytes = []byte(body)

	// TESTCASE SCENARIO
	// What is being tested: sendOneNotification function
	// Mocked: Client
	// Behaviour: The function gets passed in 2 models, one for notification
	// and one for subscription.
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status with SENT or FAILED accordingly

	suite.T().Run("Successful post, updated notification", func(t *testing.T) {

		// Under test: sendOneNotification function
		// Mocked:     Client
		// Set up:     We provide a PENDING webhook notification, and make the client return 200, success
		// Expected outcome:
		//             Client.Post called once in total
		//             Notification would be updated as SENT

		// Set beginning status as pending
		suite.DB().Find(&notification, notification.ID)
		notification.Status = models.WebhookNotificationPending
		suite.DB().ValidateAndUpdate(&notification)

		// Expectation: When Post is called, verify it was called with correct url.
		// Then, make it return 200 success and a body. It should run once.
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&responseSuccess, bodyBytes, nil).Once()

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was no error
		suite.Nil(err)
		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 1)
		// Check that notification Status was set to Sent in the model
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationSent, notif.Status)
		// Check that first attempted at date was set
		suite.NotNil(notif.FirstAttemptedAt)

	})

	suite.T().Run("Failed post, updated notification", func(t *testing.T) {

		// Under test: sendOneNotification function
		// Mocked:     Client
		// Set up:     We provide a PENDING webhook notification, and make the client return 400
		// Expected outcome:
		//             Client.Post called 3 times in total
		//             Notification would be updated as FAILING
		//             Subscription should be updated as FAILING

		// Set original status as pending
		suite.DB().Find(&notification, notification.ID)
		notification.Status = models.WebhookNotificationPending
		suite.DB().ValidateAndUpdate(&notification)

		// Recreate mock to clear stats
		mockClient := mocks.WebhookRuntimeClient{}
		engine.Client = &mockClient

		// Set Expectation: When Post is called, verify it was called with the callback url from the subscription.
		// Then, return failure.
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&responseFail, bodyBytes, nil)

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was an error returned
		suite.NotNil(err)
		// Check that the set expectations were met (the mockClient.On call) and that it was called 3 times
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 3)
		// Check that notification Status was set to Failing
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationFailing, notif.Status)
		// Check that first attempted at date was set
		suite.False(notif.FirstAttemptedAt.IsZero())
	})

	suite.T().Run("Failed post due to send error, updated notification", func(t *testing.T) {

		// Under test: sendOneNotification function
		// Mocked:     Client
		// Set up:     We provide a PENDING webhook notification, and make the client fail every time
		// Expected outcome:
		//             Client.Post called 3 times in total
		//             Notification would be updated as FAILING

		// Set original status as pending
		suite.DB().Find(&notification, notification.ID)
		notification.Status = models.WebhookNotificationPending
		suite.DB().ValidateAndUpdate(&notification)

		// Recreate mock to clear stats
		mockClient := mocks.WebhookRuntimeClient{}
		engine.Client = &mockClient

		// Expectation: When Post is called, verify it was called with the callback url from the subscription.
		// Then make it return an error.
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&responseSuccess, bodyBytes, errors.New("Error due to server down"))

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was an error returned
		suite.NotNil(err)
		// Check that the set expectations were met (the mockClient.On call) and that it was called 3 times
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 3)
		// Check that notification Status was set to Failing
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationFailing, notif.Status)
		// Check that first attempted at date was set
		suite.False(notif.FirstAttemptedAt.IsZero())
	})

	suite.T().Run("Failed post twice, then success, updated notification", func(t *testing.T) {

		// Under test: sendOneNotification function
		// Mocked:     Client
		// Set up:     We provide a FAILING webhook notification, and make the client fail twice to send
		//             and then succeed
		// Expected outcome:
		//             Client.Post called 3 times in total
		//             Notification would be updated as SENT

		// Set original status as failing
		suite.DB().Find(&notification, notification.ID)
		notification.Status = models.WebhookNotificationFailing
		suite.DB().ValidateAndUpdate(&notification)

		// Recreate mock to clear stats
		mockClient := mocks.WebhookRuntimeClient{}
		engine.Client = &mockClient

		// Set Expectation: When Post is called, return failure twice
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&responseFail, bodyBytes, nil).Twice()

		// Then return success
		mockClient.On("Post", mock.Anything, subscription.CallbackURL).Return(&responseSuccess, bodyBytes, nil)

		// Call the engine function. Internally it should call the mocked client
		err := engine.sendOneNotification(&notification, &subscription)

		// Check that there was no error
		suite.Nil(err)
		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 3)
		// Check that notification Status was set to Sent
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationSent, notif.Status)
		// Check that first attempted at date was set
		suite.False(notif.FirstAttemptedAt.IsZero())
	})

}

func (suite *WebhookClientTestingSuite) Test_EngineRunSuccessful() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client
	// Set up:     We provide a 3 PENDING webhook notifications,
	//             1 active subscription
	//             And make the client return 200, success
	// Expected outcome:
	//             All three Notifications would be updated as SENT

	// SETUP SCENARIO
	// Setup 3 pending notifications notifications, and subscriptions, and the engine
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	var err error
	response.StatusCode = 200
	response.Status = "200 Success"

	// SETUP MOCKED OBJECT EXPECTATIONS
	// Expectation: When Post is called, verify it was called with the callback url from the subscription.
	// Then, make it return 200 success and a body
	bodyBytes := []byte("notification0 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[0].ID)
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	bodyBytes = []byte("notification1 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[1].ID)
	}), subscriptions[1].CallbackURL).Return(&response, bodyBytes, nil)

	bodyBytes = []byte("notification2 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[2].ID)
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	// RUN TEST
	// Call the engine function. Internally it should call the mocked client
	err = engine.run()

	// VERIFY RESULTS
	// Check that there was no error
	suite.Nil(err)
	// Check that the set expectations were met (the mockClient.On call)
	mockClient.AssertExpectations(suite.T())

	// Check that notification Status was set to Sent on all three notifications
	updatedNotifs := []models.WebhookNotification{}
	suite.DB().All(&updatedNotifs)
	for _, notif := range updatedNotifs {
		suite.Equal(models.WebhookNotificationSent, notif.Status)
		suite.False(notif.FirstAttemptedAt.IsZero())
	}
}

func (suite *WebhookClientTestingSuite) Test_EngineRunInactiveSub() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client
	// Set up:     We provide a 3 PENDING webhook notifications,
	//             1 active subscription for PaymentUpdate
	//             No active subscription for PaymentCreate
	//             And make the client return 200, success
	// Expected outcome:
	//             PaymentUpdate notification would be updated as SENT
	//             PaymentCreate notification would be updated as SKIPPED

	// SETUP SCENARIO
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	response.StatusCode = 200
	response.Status = "200 Success"

	// Change the eventkey on 1 notification to PaymentCreate
	notifications[1].EventKey = "Payment.Create"
	verrs, err := suite.DB().ValidateAndUpdate(&notifications[1])
	if verrs != nil {
		suite.False(verrs.HasAny())
	}
	suite.Nil(err)

	// Deactivate the subscription for Payment.Create and save it
	subscriptions[1].EventKey = "Payment.Create"
	subscriptions[1].Status = models.WebhookSubscriptionStatusDisabled
	verrs, err = suite.DB().ValidateAndUpdate(&subscriptions[1])
	if verrs != nil {
		suite.False(verrs.HasAny())
	}
	suite.Nil(err)

	// SETUP MOCKED OBJECT EXPECTATIONS
	// Expectation: Post will be called twice for notification 1 and 3
	// Check the notification ID in the param
	// Then, make it return 200 success and a body
	bodyBytes := []byte("notification1 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[0].ID)
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	bodyBytes = []byte("notification3 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[2].ID)
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	// RUN TEST
	// Call the engine function. Internally it should call the mocked client
	err = engine.run()

	// VERIFY RESULTS
	// Check that there was no error
	suite.Nil(err)
	// Check that the set expectations were met (the mockClient.On call)
	mockClient.AssertExpectations(suite.T())

	// Check that notification Status was Skipped on Payment.Create
	// but set to Sent on the other notifications
	updatedNotifs := []models.WebhookNotification{}
	suite.DB().All(&updatedNotifs)
	for _, notif := range updatedNotifs {
		if notif.EventKey == "Payment.Create" {
			// if there's no subscription, we except status to be skipped
			suite.Equal(models.WebhookNotificationSkipped, notif.Status)
			// And we except firstAttemptedAt to be unset
			suite.Nil(notif.FirstAttemptedAt)
		} else {
			suite.Equal(models.WebhookNotificationSent, notif.Status)
			suite.False(notif.FirstAttemptedAt.IsZero())
		}
	}

}
func (suite *WebhookClientTestingSuite) Test_EngineRunFailingSub() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client
	// Set up:     We provide a 3 PENDING webhook notifications
	//             1 active subscription for PaymentUpdate, client returns success
	//             1 active subscription for PaymentCreate, client returns failure
	// Expected outcome:
	//             PaymentUpdate notification would be updated as SENT
	//             PaymentCreate notification would be updated as FAILING
	//             PaymentCreate subscription would be updated as FAILING

	// SETUP SCENARIO
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	response.StatusCode = 200
	response.Status = "200 Success"

	// Change the eventkey on 1 notification to PaymentCreate
	notifications[1].EventKey = "Payment.Create"
	verrs, err := suite.DB().ValidateAndUpdate(&notifications[1])
	if verrs != nil {
		suite.False(verrs.HasAny())
	}
	suite.Nil(err)

	// Some responses for client to return
	var responseSuccess = http.Response{
		Status:     "200 Success",
		StatusCode: 200,
	}
	var responseFail = http.Response{
		Status:     "400 Not Found Error",
		StatusCode: 400,
	}

	// SETUP MOCKED OBJECT EXPECTATIONS
	// Expectation: Post will be once for notification 1 and return success
	// It will be called 3 times for notification 2 and return failure
	bodyBytes := []byte("notification1 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[0].ID)
	}), subscriptions[0].CallbackURL).Return(&responseSuccess, bodyBytes, nil)

	bodyBytes = []byte("notification2 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == *handlers.FmtUUID(notifications[1].ID)
	}), subscriptions[1].CallbackURL).Return(&responseFail, bodyBytes, nil)

	// RUN TEST
	// Call the engine function. Internally it should call the mocked client
	err = engine.run()

	// VERIFY RESULTS
	// Check that there was no error
	suite.Nil(err)
	// Check that the set expectations were met (the mockClient.On call)
	mockClient.AssertExpectations(suite.T())
	mockClient.AssertNumberOfCalls(suite.T(), "Post", 4)

	// Check that notification Status was SENT on 1st notification
	updatedNotifs := []models.WebhookNotification{}
	suite.DB().Order("created_at asc").All(&updatedNotifs)
	// First notification should have sent
	suite.Equal(models.WebhookNotificationSent, updatedNotifs[0].Status)
	suite.False(updatedNotifs[0].FirstAttemptedAt.IsZero())

	// Second notification should be FAILING
	suite.Equal(models.WebhookNotificationFailing, updatedNotifs[1].Status)
	suite.False(updatedNotifs[1].FirstAttemptedAt.IsZero())
	// Subscription should be set to FAILING
	suite.DB().Find(&subscriptions[1], subscriptions[1].ID)
	suite.Equal(models.WebhookSubscriptionStatusFailing, subscriptions[1].Status)

	// Third notification should be PENDING
	suite.Equal(models.WebhookNotificationPending, updatedNotifs[2].Status)
	suite.Nil(updatedNotifs[2].FirstAttemptedAt)

}

func (suite *WebhookClientTestingSuite) Test_EngineRunFailedSubWithSeverity() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client object that sends the HTTP request
	// Set up:     We provide a PENDING webhook notification with an ACTIVE subscription.
	//             Client returns failure repeatedly
	//             We update the firstAttemptedAt time to mimic a notification that's been failing
	//             for a while to test the severity thresholds
	// Expected outcome:
	//             After first failure - notif marked as failing, subscription severity = 3
	//             After second failure one minute later - notif marked as failing, subscription severity = 3
	//             After first threshold - notif marked as failing, subscription severity = 2
	//             After final threshold - notif marked as failed, subscription severity = 1, subscription deactivated
	//

	// SETUP SCENARIO
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	// We only need 1st notification, delete the others
	suite.DB().Destroy(&notifications[1])
	suite.DB().Destroy(&notifications[2])

	// Create a fail response for the mocked client to return
	var responseFail = http.Response{
		Status:     "400 Not Found Error",
		StatusCode: 400,
	}
	numExpectedPosts := 0

	suite.T().Run("Severity 3 failure", func(t *testing.T) {
		// Set up:     We provide a PENDING webhook notification with an ACTIVE subscription.
		//             Client returns failure repeatedly
		// Expected outcome:
		//             After first failure - notif marked as failing, subscription severity = 3

		// SETUP MOCKED OBJECT EXPECTATIONS
		// Expectation: Client.Post will be called and will return failure
		mockClient.On("Post", mock.Anything, subscriptions[0].CallbackURL).Return(&responseFail, nil, errors.New("Mocked webhook client fails to send"))

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		numExpectedPosts += engine.MaxImmediateRetries
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", numExpectedPosts)

		// Check that notification is marked as FAILING
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(notifications[0].Status, models.WebhookNotificationFailing)

		// Check that subscription is marked as FAILING, with severity 3
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusFailing, subscriptions[0].Status)
		suite.Equal(3, subscriptions[0].Severity)

	})

	suite.T().Run("Severity 3 failure, no raised severity", func(t *testing.T) {

		// Set up:     Notification has failed once, marked as FAILING
		// Expected outcome:
		//             After second failure one minute later - notif still marked as FAILING, subscription severity = 3

		// Update firstAttemptedTime to be a minute ago
		timestamp := *(notifications[0].FirstAttemptedAt)
		timestamp = timestamp.Add(-60 * time.Second)
		notifications[0].FirstAttemptedAt = &timestamp
		suite.DB().ValidateAndUpdate(&notifications[0])

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		numExpectedPosts += engine.MaxImmediateRetries
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", numExpectedPosts)

		// Check that notification is marked as FAILING
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(notifications[0].Status, models.WebhookNotificationFailing)

		// Check that subscription is marked as FAILING, with severity 3
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusFailing, subscriptions[0].Status)
		suite.Equal(3, subscriptions[0].Severity)

	})

	suite.T().Run("Severity 2 failure", func(t *testing.T) {

		// Set up:     Notification has failed once, marked as FAILING
		//			   We update the firstAttemptedAt time to mimic a notification that's been failing
		//             longer than the first threshold
		// Expected outcome:
		//             After first threshold - notif marked as FAILING, subscription severity = 2

		// Update firstAttemptedTime to be more than one threshold ago
		durationOffset := time.Duration(engine.SeverityThresholds[0]) * time.Second
		timestamp := *(notifications[0].FirstAttemptedAt)
		timestamp = timestamp.Add(-durationOffset)
		notifications[0].FirstAttemptedAt = &timestamp
		suite.DB().ValidateAndUpdate(&notifications[0])

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		numExpectedPosts += engine.MaxImmediateRetries
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", numExpectedPosts)

		// Check that notification is marked as FAILING
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(models.WebhookNotificationFailing, notifications[0].Status)

		// Check that subscription is marked as FAILING, with severity 2
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusFailing, subscriptions[0].Status)
		suite.Equal(2, subscriptions[0].Severity)

	})

	suite.T().Run("Severity 1 failure - deactivation", func(t *testing.T) {

		// Set up:     Notification is FAILING already
		//			   We update the firstAttemptedAt time to mimic a notification that's been failing
		//             longer than the final threshold
		// Expected outcome:
		//             After final threshold - notif marked as FAILED, subscription severity = 1, subscription DISABLED

		// Update firstAttemptedTime to be more than one threshold ago
		durationOffset := time.Duration(engine.SeverityThresholds[1]) * time.Second
		timestamp := *(notifications[0].FirstAttemptedAt)
		timestamp = timestamp.Add(-durationOffset)
		notifications[0].FirstAttemptedAt = &timestamp
		suite.DB().ValidateAndUpdate(&notifications[0])

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		numExpectedPosts += engine.MaxImmediateRetries
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", numExpectedPosts)

		// Check that notification is marked as FAILED
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(models.WebhookNotificationFailed, notifications[0].Status)

		// Check that subscription is marked as DISABLED, with severity 1
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusDisabled, subscriptions[0].Status)
		suite.Equal(1, subscriptions[0].Severity)

	})

	suite.T().Run("Notification not tried again", func(t *testing.T) {

		// Set up:     Notification has FAILED, subscription has been DISABLED
		// Expected outcome:
		//             Engine no longer attempts to send this notification.

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		// numExpectedPosts should have no change from previous run, because client was not called.
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", numExpectedPosts)

		// Check that notification is marked as FAILING
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(models.WebhookNotificationFailed, notifications[0].Status)

		// Check that subscription is marked as DISABLED, with severity 1
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusDisabled, subscriptions[0].Status)
		suite.Equal(1, subscriptions[0].Severity)

	})

}

func (suite *WebhookClientTestingSuite) Test_EngineRunFailingRecovery() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client
	// Set up:     We provide 3 PENDING webhook notifications with active subscriptions.
	//             Client returns failure repeatedly on the first run, then recovers on the second run.
	// Expected outcome:
	//             After failure - notif marked as failing, sub severity = 2, subscription status = failing
	//             After success - notif marked as sent, subscription severity = 0, subscription status = active

	// SETUP SCENARIO
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var responseSuccess = http.Response{
		Status:     "200 Success",
		StatusCode: 200,
	}
	var responseFail = http.Response{
		Status:     "400 Not Found Error",
		StatusCode: 400,
	}

	suite.T().Run("Severity 3 failure", func(t *testing.T) {

		// Set up:     We provide 3 PENDING webhook notifications with active subscriptions.
		//             Client returns failure repeatedly
		// Expected outcome:
		//             After failure - notif marked as FAILING, sub severity = 3, subscription status FAILING

		// Make mockClient fail to send
		mockClient.On("Post", mock.Anything, subscriptions[0].CallbackURL).Return(&responseFail, nil, errors.New("Mocked webhook client fails to send"))

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 3)

		// Check that notification is marked as FAILING
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(models.WebhookNotificationFailing, notifications[0].Status)

		// Check that subscription is marked as FAILING, with severity 3
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusFailing, subscriptions[0].Status)
		suite.Equal(3, subscriptions[0].Severity)

	})

	suite.T().Run("Successful recovery", func(t *testing.T) {
		// Set up:     We provide 3 PENDING webhook notifications with active subscriptions.
		//             One notification and subscription is marked as FAILING
		//             Client succeeds this time
		// Expected outcome:
		//             After success - All 3 notifs marked as sent, subscription severity = 0, subscription status = active

		// Set up mock for success
		mockClient = &mocks.WebhookRuntimeClient{}
		engine.Client = mockClient
		bodyBytes := []byte("notification0 received")
		mockClient.On("Post", mock.Anything, mock.Anything).Return(&responseSuccess, bodyBytes, nil)

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 3)

		// Check that notifications are marked as SENT
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(models.WebhookNotificationSent, notifications[0].Status)
		suite.DB().Find(&notifications[1], notifications[1].ID)
		suite.Equal(models.WebhookNotificationSent, notifications[1].Status)
		suite.DB().Find(&notifications[2], notifications[2].ID)
		suite.Equal(models.WebhookNotificationSent, notifications[2].Status)

		// Check that subscriptions are marked as ACTIVE, with severity 0
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusActive, subscriptions[0].Status)
		suite.Equal(0, subscriptions[0].Severity)
		suite.DB().Find(&subscriptions[1], subscriptions[1].ID)
		suite.Equal(models.WebhookSubscriptionStatusActive, subscriptions[1].Status)
		suite.Equal(0, subscriptions[1].Severity)

	})
}

func (suite *WebhookClientTestingSuite) Test_EngineRunNoThresholds() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client
	// Set up:     We provide 3 PENDING webhook notifications with active subscriptions.
	//             No thresholds are set, empty array
	//             Client returns failure repeatedly on the first run, then recovers on the second run.
	// Expected outcome:
	//             After failure - notif marked as FAILED, sub severity = 1, subscription status = DISABLED
	//             No panics!

	// SETUP SCENARIO
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	engine.SeverityThresholds = []int{}
	defer teardownEngineRun(suite)

	var responseFail = http.Response{
		Status:     "400 Not Found Error",
		StatusCode: 400,
	}

	suite.T().Run("Any failure is Severity 1", func(t *testing.T) {

		// Set up:     We provide 3 PENDING webhook notifications with active subscriptions.
		//             No thresholds are set, empty array
		//             Client returns failure repeatedly on the first run, then recovers on the second run.
		// Expected outcome:
		//             After failure - notif marked as FAILED, sub severity = 1, subscription status = DISABLED
		//             No panics!

		// Make mockClient fail to send
		mockClient.On("Post", mock.Anything, subscriptions[0].CallbackURL).Return(&responseFail, nil, errors.New("Mocked webhook client fails to send"))

		// RUN TEST
		// Call the engine function. Internally it should call the mocked client
		err := engine.run()

		// VERIFY RESULTS
		// Check that there was no error
		suite.Nil(err)

		// Check that the set expectations were met (the mockClient.On call)
		mockClient.AssertExpectations(suite.T())
		mockClient.AssertNumberOfCalls(suite.T(), "Post", 3)

		// Check that notification is marked as FAILED
		suite.DB().Find(&notifications[0], notifications[0].ID)
		suite.Equal(models.WebhookNotificationFailed, notifications[0].Status)

		// Check that subscription is marked as DISABLED, with severity 1
		suite.DB().Find(&subscriptions[0], subscriptions[0].ID)
		suite.Equal(models.WebhookSubscriptionStatusDisabled, subscriptions[0].Status)
		suite.Equal(1, subscriptions[0].Severity)
	})
}

func (suite *WebhookClientTestingSuite) Test_EngineRunNoPending() {

	// TESTCASE SCENARIO
	// Under test: Engine.run() function
	// Mocked:     Client
	// Set up:     We provide a 3 SENT, FAILED, SKIPPED webhook notifications,
	//             1 active subscription
	//             And make the client return 200, success
	// Expected outcome:
	//             Since no notifications are PENDING or FAILING
	//   		   Expect no calls to Client.Post

	// SETUP SCENARIO
	engine, notifications, _ := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	response.StatusCode = 200
	response.Status = "200 Success"

	// Change the status on the notifications to sent
	notifications[0].Status = models.WebhookNotificationSent
	suite.DB().ValidateAndUpdate(&notifications[0])

	notifications[1].Status = models.WebhookNotificationFailed
	suite.DB().ValidateAndUpdate(&notifications[1])

	notifications[2].Status = models.WebhookNotificationSkipped
	suite.DB().ValidateAndUpdate(&notifications[2])

	// SETUP MOCKED OBJECT EXPECTATIONS
	// Expectation: We set up a possible call here, but we will be checking that in fact
	// it was NOT called.
	mockClient.On("Post", mock.Anything, mock.Anything).Return(&response, []byte(""), nil)

	// RUN TEST
	// Call the engine function. Internally it should call the mocked client
	err := engine.run()

	// VERIFY RESULTS
	// Check that there was no error
	suite.Nil(err)
	// Check that the Post function was not called
	mockClient.AssertNotCalled(suite.T(), "Post", mock.Anything, mock.Anything)

}

func setupEngineRun(suite *WebhookClientTestingSuite) (*Engine, []models.WebhookNotification, []models.WebhookSubscription) {
	mockClient := mocks.WebhookRuntimeClient{}

	// Create the engine replacing the client with the mock client
	engine := Engine{
		DB:                  suite.DB(),
		Logger:              suite.logger,
		Client:              &mockClient,
		MaxImmediateRetries: 3,
		SeverityThresholds:  []int{1800, 14400},
	}
	// Create 3 notifications
	// Pending notification for Payment.Update
	notification0 := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			EventKey: "Payment.Update",
			Payload:  "{\"message\":\"This is an updated notification #0\"}",
		},
	})
	// Pending notification for Payment.Create
	notification1 := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			EventKey: "Payment.Create",
			Payload:  "{\"message\":\"This is an updated notification #1\"}",
		},
	})
	// Pending notification for Payment.Update
	notification2 := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			EventKey: "Payment.Update",
			Payload:  "{\"message\":\"This is an updated notification #2\"}",
		},
	})

	// Active subscription for the Payment.Update event
	subscription0 := testdatagen.MakeWebhookSubscription(suite.DB(), testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			EventKey:    "Payment.Update",
			CallbackURL: "/my/callback/url/0",
		},
	})
	// Active subscription for the Payment.Create event
	subscription1 := testdatagen.MakeWebhookSubscription(suite.DB(), testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			EventKey:    "Payment.Create",
			CallbackURL: "/my/callback/url/1",
		},
	})

	// return an array of the object created, plus the engine
	notifications := []models.WebhookNotification{notification0, notification1, notification2}
	subscriptions := []models.WebhookSubscription{subscription0, subscription1}
	return &engine, notifications, subscriptions
}

// truncateAllNotifications truncates the notifications table
func truncateAllNotifications(db *pop.Connection) {
	notifications := []models.WebhookNotification{}
	db.All(&notifications)
	for _, notif := range notifications {
		copyOfNotify := notif // Make copy to avoid implicit memory aliasing of items from a range statement.
		db.Destroy(&copyOfNotify)
	}
}

// truncateAllSubscriptions truncates the subscriptions table
func truncateAllSubscriptions(db *pop.Connection) {
	subscriptions := []models.WebhookSubscription{}
	db.All(&subscriptions)
	for _, sub := range subscriptions {
		copyOfSub := sub // Make copy to avoid implicit memory aliasing of items from a range statement.
		db.Destroy(&copyOfSub)
	}
}

// teardownEngineRun truncates the notifications and subscriptions tables
func teardownEngineRun(suite *WebhookClientTestingSuite) {
	truncateAllNotifications(suite.DB())
	truncateAllSubscriptions(suite.DB())
}

// convertBodyToPayload is a helper function to convert []byte to a webhookMessage payload
func convertBodyToPayload(body []byte) supportmessages.WebhookNotification {
	message := supportmessages.WebhookNotification{}
	json.Unmarshal(body, &message)
	return message
}

type severityTestData struct {
	attempt       time.Duration
	expectedLevel int
}

func (suite *WebhookClientTestingSuite) Test_GetSeverity() {
	thresholds := []int{1800, 3600, 7200}
	engine, _, _ := setupEngineRun(suite)
	engine.SeverityThresholds = thresholds
	testData := []severityTestData{
		{attempt: -10 * time.Second, expectedLevel: 4},
		{attempt: -3000 * time.Second, expectedLevel: 3},
		{attempt: -3601 * time.Second, expectedLevel: 2},
		{attempt: -7201 * time.Second, expectedLevel: 1},
	}
	for _, data := range testData {
		suite.T().Run(fmt.Sprintf("Returns severity level %d", data.expectedLevel), func(t *testing.T) {
			currentTime := time.Now()
			attempt := currentTime.Add(data.attempt)
			severity := engine.GetSeverity(currentTime, attempt)
			suite.Equal(data.expectedLevel, severity)
		})
	}
}
