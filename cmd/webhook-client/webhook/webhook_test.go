package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
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

func (suite *WebhookClientTestingSuite) Test_SendOneNotification() {
	mockClient := mocks.WebhookRuntimeClient{}

	// Create the engine replacing the client with the mock client
	engine := Engine{
		DB:     suite.DB(),
		Logger: suite.logger,
		Client: &mockClient,
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
	// and one for subscription.
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status with SENT or FAILED accordingly

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
		// Check that notification Status was set to Sent in the model
		suite.Equal(models.WebhookNotificationSent, notification.Status)
		// Check in the db
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationSent, notif.Status)
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
		// Check that notification Status was set to Failed in the model
		suite.Equal(models.WebhookNotificationFailed, notification.Status)
		// Check in the db
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationFailed, notif.Status)
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
		// Check that notification Status was set to Failed in the model
		suite.Equal(models.WebhookNotificationFailed, notification.Status)
		// Check in the db
		notif := models.WebhookNotification{}
		suite.DB().Find(&notif, notification.ID)
		suite.Equal(models.WebhookNotificationFailed, notif.Status)
	})

}

func (suite *WebhookClientTestingSuite) Test_EngineRunSuccessful() {

	// TESTCASE SCENARIO
	// What is being tested: Engine.run() function
	// Behaviour: The function checks the db for notifications and subscriptions
	// 1. Only pending notifications should get sent
	// 2. Only if an active subscription exists
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status in the model with SENT or FAILED accordingly

	// EXPECTED BEHAVIOR IN THIS TEST
	// All three notifications should get sent, in order of creation

	// SETUP SCENARIO
	// Setup some default notifications, and subscriptions, and the engine
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	var err error
	response.StatusCode = 200
	response.Status = "200 Success"

	fmt.Printf("%v", suite.DB())

	// SETUP MOCKED OBJECT EXPECTATIONS
	// Expectation: When Post is called, verify it was called with the callback url from the subscription.
	// Then, make it return 200 success and a body
	bodyBytes := []byte("notification0 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == notifications[0].ID
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	bodyBytes = []byte("notification1 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == notifications[1].ID
	}), subscriptions[1].CallbackURL).Return(&response, bodyBytes, nil)

	bodyBytes = []byte("notification2 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == notifications[2].ID
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
		fmt.Println(notif.ID, ":", notif.EventKey)
		suite.Equal(models.WebhookNotificationSent, notif.Status)
	}
}

func (suite *WebhookClientTestingSuite) Test_EngineRunInactiveSub() {

	// TESTCASE SCENARIO
	// What is being tested: Engine.run() function
	// Behaviour: The function checks the db for notifications and subscriptions
	// 1. Only pending notifications should get sent
	// 2. Only if an active subscription exists
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status in the model with SENT or FAILED accordingly

	// EXPECTED BEHAVIOR IN THIS TEST
	// We're going to make 2 Payment.Update events and one
	// Payment.Create event. We will have 1 active subscription for Payment.Update
	// and 1 inactive subscription for Payment.Create.
	// Therefore we expect only the 2 Payment.Update notifications to get sent.

	// SETUP SCENARIO
	engine, notifications, subscriptions := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	response.StatusCode = 200
	response.Status = "200 Success"

	// Change the eventkey on notification[1] and save it
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
		return message.ID == notifications[0].ID
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	bodyBytes = []byte("notification3 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == notifications[2].ID
	}), subscriptions[0].CallbackURL).Return(&response, bodyBytes, nil)

	// RUN TEST
	// Call the engine function. Internally it should call the mocked client
	err = engine.run()

	// VERIFY RESULTS
	// Check that there was no error
	suite.Nil(err)
	// Check that the set expectations were met (the mockClient.On call)
	mockClient.AssertExpectations(suite.T())

	// Check that notification Status was still pending on Payment.Create
	// but set to Sent on the other notifications
	updatedNotifs := []models.WebhookNotification{}
	suite.DB().All(&updatedNotifs)
	for _, notif := range updatedNotifs {
		if notif.EventKey == "Payment.Create" {
			// MYTODO: Should be skipped not pending after migration
			suite.Equal(models.WebhookNotificationPending, notif.Status)
		} else {
			suite.Equal(models.WebhookNotificationSent, notif.Status)
		}
	}

}

func (suite *WebhookClientTestingSuite) Test_EngineRunNoPending() {

	// TESTCASE SCENARIO
	// What is being tested: Engine.run() function
	// Behaviour: The function checks the db for notifications and subscriptions
	// 1. Only pending notifications should get sent
	// 2. Only if an active subscription exists
	// It should create a payload from the notification and send it to the url
	// listed in the subscription. On success or failure, it should update the
	// notification.Status in the model with SENT or FAILED accordingly

	// EXPECTED BEHAVIOR IN THIS TEST
	// We're going to make all notifications status SENT or FAILED, so none should
	// actually be sent.

	// SETUP SCENARIO
	engine, notifications, _ := setupEngineRun(suite)
	mockClient := engine.Client.(*mocks.WebhookRuntimeClient)
	defer teardownEngineRun(suite)

	var response = http.Response{}
	response.StatusCode = 200
	response.Status = "200 Success"

	// Change the status on the notifications to sent
	notifications[0].Status = models.WebhookNotificationSent
	verrs, err := suite.DB().ValidateAndUpdate(&notifications[0])
	if verrs != nil {
		suite.False(verrs.HasAny())
	}
	suite.Nil(err)

	notifications[1].Status = models.WebhookNotificationFailed
	verrs, err = suite.DB().ValidateAndUpdate(&notifications[1])
	if verrs != nil {
		suite.False(verrs.HasAny())
	}
	suite.Nil(err)

	notifications[2].Status = models.WebhookNotificationSent
	verrs, err = suite.DB().ValidateAndUpdate(&notifications[2])
	if verrs != nil {
		suite.False(verrs.HasAny())
	}
	suite.Nil(err)

	// SETUP MOCKED OBJECT EXPECTATIONS
	// Expectation: We set up a possible call here, but we will be checking that in fact
	// it was NOT called.
	mockClient.On("Post", mock.Anything, mock.Anything).Return(&response, []byte(""), nil)

	// RUN TEST
	// Call the engine function. Internally it should call the mocked client
	err = engine.run()

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
		DB:     suite.DB(),
		Logger: suite.logger,
		Client: &mockClient,
	}
	// Create 3 notifications
	// Pending notification for Payment.Update
	notification0 := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			EventKey: "Payment.Update",
			Payload:  swag.String("{\"message\":\"This is an updated notification #0\"}"),
		},
	})
	// Pending notification for Payment.Create
	notification1 := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			EventKey: "Payment.Create",
			Payload:  swag.String("{\"message\":\"This is an updated notification #1\"}"),
		},
	})
	// Pending notification for Payment.Update
	notification2 := testdatagen.MakeWebhookNotification(suite.DB(), testdatagen.Assertions{
		WebhookNotification: models.WebhookNotification{
			EventKey: "Payment.Update",
			Payload:  swag.String("{\"message\":\"This is an updated notification #2\"}"),
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
		db.Destroy(&notif)
	}
}

// truncateAllSubscriptions truncates the subscriptions table
func truncateAllSubscriptions(db *pop.Connection) {
	subscriptions := []models.WebhookSubscription{}
	db.All(&subscriptions)
	for _, sub := range subscriptions {
		db.Destroy(&sub)
	}
}

// teardownEngineRun truncates the notifications and subscriptions tables
func teardownEngineRun(suite *WebhookClientTestingSuite) {
	truncateAllNotifications(suite.DB())
	truncateAllSubscriptions(suite.DB())
}

// convertBodyToPayload is a helper function to convert []byte to a webhookMessage payload
func convertBodyToPayload(body []byte) Message {
	message := Message{}
	json.Unmarshal(body, &message)
	return message
}
