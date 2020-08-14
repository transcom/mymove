package event

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type EventServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *EventServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestEventServiceSuite(t *testing.T) {
	ts := &EventServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *EventServiceSuite) Test_EventTrigger() {

	now := time.Now()

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})
	paymentRequestID := paymentRequest.ID
	mtoID := paymentRequest.MoveTaskOrderID

	unavailablePaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})

	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)

	// Test successful event passing with Support API
	suite.T().Run("trigger event passing with support api endpoint", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: paymentRequestID,
			Request:         &dummyRequest,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		suite.Nil(err)
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count+1, newCount)

	})

	// Test successful event passing with GHC API
	suite.T().Run("trigger event passing with ghc api endpoint", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: paymentRequestID,
			Request:         &dummyRequest,
			EndpointKey:     GhcUpdatePaymentRequestStatusEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		suite.Nil(err)
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count+1, newCount)

	})

	// This test verifies that if the object updated is not on an MTO that
	// is available to prime, no notification is created.
	suite.T().Run("trigger event no notification - unavailable mto", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		unavailablePRID := unavailablePaymentRequest.ID
		unavailableMTOID := unavailablePaymentRequest.MoveTaskOrderID

		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           unavailableMTOID,
			UpdatedObjectID: unavailablePRID,
			Request:         &dummyRequest,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		suite.Nil(err)

		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, newCount)

	})

	suite.T().Run("trigger event error - bad event key", func(t *testing.T) {
		// Pass a bad event key
		_, err := TriggerEvent(Event{
			EventKey:        "BadEventKey",
			MtoID:           mtoID,
			UpdatedObjectID: paymentRequestID,
			Request:         &dummyRequest,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		// Check that at least one error was returned
		suite.NotNil(err)
	})
	suite.T().Run("trigger event error - bad endpoint key", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		// Pass a bad endpoint key
		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: paymentRequestID,
			Request:         &dummyRequest,
			EndpointKey:     "Bad Endpoint Key That Doesn't Exist",
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		// Check that at least one error was returned
		suite.NotNil(err)
		// Check that no notification was created
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, newCount)

	})
	suite.T().Run("trigger event error - bad object ID", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		// Pass a bad payment request ID
		uuidString := "88c9922f-58c7-45cd-8c10-48f2a52bbabc"
		paymentRequestID, _ = uuid.FromString(uuidString)
		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: paymentRequestID,
			Request:         &dummyRequest,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		// Check that at least one error was returned
		suite.NotNil(err)
		// Check that no notification was created
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, newCount)

	})
}
