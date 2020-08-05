package event

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
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

func (suite *EventServiceSuite) Test_EventRecord() {

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
	paymentRequestID := paymentRequest.ID
	mtoID := paymentRequest.MoveTaskOrderID

	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger := zap.NewNop()
	handler := handlers.NewHandlerContext(suite.DB(), logger)

	suite.T().Run("trigger event passing", func(t *testing.T) {
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
	})
	suite.T().Run("trigger event error - bad object ID", func(t *testing.T) {
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
	})
}
