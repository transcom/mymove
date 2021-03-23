package event

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type EventServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *EventServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestEventServiceSuite(t *testing.T) {
	ts := &EventServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *EventServiceSuite) getNotification(mtoID uuid.UUID, traceID uuid.UUID) (*models.WebhookNotification, error) {
	var notification = new(models.WebhookNotification)
	err := suite.DB().Where("object_id = $1 AND trace_id = $2", mtoID.String(), traceID.String()).First(notification)
	return notification, err
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

	unavailablePaymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)

	// Test successful event passing with Support API
	suite.T().Run("Success with support api endpoint", func(t *testing.T) {
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
	suite.T().Run("Success with ghc api endpoint", func(t *testing.T) {
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
	suite.T().Run("Fail with no notification - unavailable mto", func(t *testing.T) {
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

	suite.T().Run("Fail with bad event key", func(t *testing.T) {
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
	suite.T().Run("Fail with bad endpoint key", func(t *testing.T) {
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
	suite.T().Run("Fail with bad object ID", func(t *testing.T) {
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

func (suite *EventServiceSuite) Test_MTOEventTrigger() {

	now := time.Now()
	mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})
	mtoID := mto.ID

	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)
	traceID, _ := uuid.NewV4()
	handler.SetTraceID(traceID)

	// Test successful event
	suite.T().Run("Success with GHC MoveTaskOrder endpoint", func(t *testing.T) {

		_, err := TriggerEvent(Event{
			EventKey:        MoveTaskOrderCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoID,
			Request:         &dummyRequest,
			EndpointKey:     GhcUpdateMoveTaskOrderEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		suite.Nil(err)

		// Get the notification
		notification, err := suite.getNotification(mtoID, traceID)
		suite.NoError(err)
		suite.Equal(&mtoID, notification.ObjectID)

		// Reinflate the json from the notification payload
		suite.NotEmpty(notification.Payload)
		var mtoInPayload MoveTaskOrder
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
		//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
		//RA: in a unit test, then there is no risk
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: N/A
		json.Unmarshal([]byte(notification.Payload), &mtoInPayload) // nolint:errcheck
		// Check some params
		suite.Equal(mto.PPMType, &mtoInPayload.PpmType)
		suite.Equal(handlers.FmtDateTimePtr(mto.AvailableToPrimeAt).String(), mtoInPayload.AvailableToPrimeAt.String())

	})
}

func (suite *EventServiceSuite) Test_MTOShipmentEventTrigger() {

	now := time.Now()
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})

	mtoShipmentID := mtoShipment.ID
	mtoID := mtoShipment.MoveTaskOrderID

	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)
	traceID, _ := uuid.NewV4()
	handler.SetTraceID(traceID)

	// Test successful event passing with Support API
	suite.T().Run("Success with GHC MTOShipment endpoint", func(t *testing.T) {

		_, err := TriggerEvent(Event{
			EventKey:        MTOShipmentUpdateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoShipmentID,
			Request:         &dummyRequest,
			EndpointKey:     GhcPatchMTOShipmentStatusEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		suite.Nil(err)

		// Get the notification
		notification, err := suite.getNotification(mtoShipmentID, traceID)
		suite.NoError(err)
		suite.Equal(&mtoShipmentID, notification.ObjectID)

		// Reinflate the json from the notification payload
		suite.NotEmpty(notification.Payload)
		var mtoShipmentInPayload primemessages.MTOShipment
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
		//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
		//RA: in a unit test, then there is no risk
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: N/A
		json.Unmarshal([]byte(notification.Payload), &mtoShipmentInPayload) // nolint:errcheck
		// Check some params
		suite.EqualValues(mtoShipment.ShipmentType, mtoShipmentInPayload.ShipmentType)
		suite.EqualValues(handlers.FmtDatePtr(mtoShipment.RequestedPickupDate).String(), mtoShipmentInPayload.RequestedPickupDate.String())

	})
}

func (suite *EventServiceSuite) Test_MTOServiceItemEventTrigger() {

	now := time.Now()
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
	})

	mtoServiceItemID := mtoServiceItem.ID
	mtoID := mtoServiceItem.MoveTaskOrderID

	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)

	// Test successful event passing with Support API
	suite.T().Run("Success with GHC ServiceItem endpoint", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		_, err := TriggerEvent(Event{
			EventKey:        MTOServiceItemCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoServiceItemID,
			Request:         &dummyRequest,
			EndpointKey:     GhcCreateMTOServiceItemEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})

		suite.Nil(err)
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count+1, newCount)

	})
}

func (suite *EventServiceSuite) TestOrderEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)
	traceID, _ := uuid.NewV4()
	handler.SetTraceID(traceID)

	// Test successful event passing with Support API
	suite.T().Run("Success with GHC ServiceItem endpoint", func(t *testing.T) {
		_, err := TriggerEvent(Event{
			EventKey:        OrderUpdateEventKey,
			MtoID:           move.ID,
			UpdatedObjectID: move.OrdersID,
			Request:         &dummyRequest,
			EndpointKey:     InternalUpdateOrdersEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		})
		suite.Nil(err)

		// Get the notification
		notification, err := suite.getNotification(move.OrdersID, traceID)
		suite.NoError(err)
		suite.Equal(&move.OrdersID, notification.ObjectID)

		// Reinflate the json from the notification payload
		suite.NotEmpty(notification.Payload)
		var orderPayload primemessages.Order
		err = json.Unmarshal([]byte(notification.Payload), &orderPayload)
		suite.FatalNoError(err)

		// Check some params
		suite.Equal(move.Orders.ServiceMember.ID.String(), orderPayload.Customer.ID.String())
		suite.Equal(move.Orders.Entitlement.ID.String(), orderPayload.Entitlement.ID.String())
		suite.Equal(move.Orders.OriginDutyStation.ID.String(), orderPayload.OriginDutyStation.ID.String())
	})
}

func (suite *EventServiceSuite) TestNotificationEventHandler() {
	order := testdatagen.MakeDefaultOrder(suite.DB())
	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewHandlerContext(suite.DB(), logger)
	traceID, _ := uuid.NewV4()
	handler.SetTraceID(traceID)

	// Test a nil MTO ID is present and no notification stored
	suite.T().Run("No move and notification stored", func(t *testing.T) {
		count, _ := suite.DB().Count(&models.WebhookNotification{})
		event := Event{
			EventKey:        OrderUpdateEventKey,
			MtoID:           uuid.Nil,
			UpdatedObjectID: order.ID,
			Request:         &dummyRequest,
			EndpointKey:     InternalUpdateOrdersEndpointKey,
			HandlerContext:  handler,
			DBConnection:    suite.DB(),
		}
		_, err := TriggerEvent(event)
		suite.NoError(err)

		afterCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, afterCount)

		// No notification stored and nil error returned
		_, err = suite.getNotification(order.ID, traceID)
		suite.Error(err)
		suite.Equal("sql: no rows in result set", err.Error())
	})
}
