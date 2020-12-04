package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

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
		suite.False(notif.FirstAttemptedAt.IsZero())

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

	// Check that notification Status was Skipped on Payment.Create
	// but set to Sent on the other notifications
	updatedNotifs := []models.WebhookNotification{}
	suite.DB().All(&updatedNotifs)
	for _, notif := range updatedNotifs {
		if notif.EventKey == "Payment.Create" {
			// if there's no subscription, we except status to be skipped
			suite.Equal(models.WebhookNotificationSkipped, notif.Status)
			// And we except firstAttemptedAt to be unset
			suite.True(notif.FirstAttemptedAt.IsZero())
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
		return message.ID == notifications[0].ID
	}), subscriptions[0].CallbackURL).Return(&responseSuccess, bodyBytes, nil)

	bodyBytes = []byte("notification2 received")
	mockClient.On("Post", mock.MatchedBy(func(body []byte) bool {
		message := convertBodyToPayload(body)
		return message.ID == notifications[1].ID
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
	suite.True(updatedNotifs[2].FirstAttemptedAt.IsZero())

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
		// #nosec G601 TODO needs review
		db.Destroy(&notif)
	}
}

// truncateAllSubscriptions truncates the subscriptions table
func truncateAllSubscriptions(db *pop.Connection) {
	subscriptions := []models.WebhookSubscription{}
	db.All(&subscriptions)
	for _, sub := range subscriptions {
		// #nosec G601 TODO needs review
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

type myData struct {
	attempt       time.Duration
	expectedLevel int
}

func (suite *WebhookClientTestingSuite) Test_GetSeverity() {
	thresholds := []int{1800, 3600, 7200}
	engine, _, _ := setupEngineRun(suite)
	testData := []myData{
		{attempt: -10 * time.Second, expectedLevel: 4},
		{attempt: -3000 * time.Second, expectedLevel: 3},
		{attempt: -3601 * time.Second, expectedLevel: 2},
		{attempt: -7201 * time.Second, expectedLevel: 1},
	}
	for _, data := range testData {
		suite.T().Run(fmt.Sprintf("Returns severity level %d", data.expectedLevel), func(t *testing.T) {
			currentTime := time.Now()
			attempt := currentTime.Add(data.attempt)
			severity := engine.GetSeverity(currentTime, attempt, thresholds)
			suite.Equal(data.expectedLevel, severity)
		})
	}
}
