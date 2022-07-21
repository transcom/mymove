package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type EventServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestEventServiceSuite(t *testing.T) {
	ts := &EventServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
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
	setupTestData := func() models.PaymentRequest {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})

		return paymentRequest
	}

	// Test successful event passing with Support API
	suite.Run("Success with support api endpoint", func() {
		paymentRequest := setupTestData()
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           paymentRequest.MoveTaskOrderID,
			UpdatedObjectID: paymentRequest.ID,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         uuid.Must(uuid.NewV4()),
		})
		suite.Nil(err)
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count+1, newCount)

	})

	// Test successful event passing with GHC API
	suite.Run("Success with ghc api endpoint", func() {
		paymentRequest := setupTestData()
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           paymentRequest.MoveTaskOrderID,
			UpdatedObjectID: paymentRequest.ID,
			EndpointKey:     GhcUpdatePaymentRequestStatusEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         uuid.Must(uuid.NewV4()),
		})
		suite.Nil(err)
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count+1, newCount)

	})

	// This test verifies that if the object updated is not on an MTO that
	// is available to prime, no notification is created.
	suite.Run("Fail with no notification - unavailable mto", func() {
		unavailablePaymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		unavailablePRID := unavailablePaymentRequest.ID
		unavailableMTOID := unavailablePaymentRequest.MoveTaskOrderID

		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           unavailableMTOID,
			UpdatedObjectID: unavailablePRID,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         uuid.Must(uuid.NewV4()),
		})
		suite.Nil(err)

		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, newCount)

	})

	suite.Run("Fail with bad event key", func() {
		paymentRequest := setupTestData()
		// Pass a bad event key
		_, err := TriggerEvent(Event{
			EventKey:        "BadEventKey",
			MtoID:           paymentRequest.MoveTaskOrderID,
			UpdatedObjectID: paymentRequest.ID,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         uuid.Must(uuid.NewV4()),
		})
		// Check that at least one error was returned
		suite.NotNil(err)
	})
	suite.Run("Fail with bad endpoint key", func() {
		paymentRequest := setupTestData()
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		// Pass a bad endpoint key
		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           paymentRequest.MoveTaskOrderID,
			UpdatedObjectID: paymentRequest.ID,
			EndpointKey:     "Bad Endpoint Key That Doesn't Exist",
			AppContext:      suite.AppContextForTest(),
			TraceID:         uuid.Must(uuid.NewV4()),
		})
		// Check that at least one error was returned
		suite.NotNil(err)
		// Check that no notification was created
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, newCount)

	})
	suite.Run("Fail with bad object ID", func() {
		paymentRequest := setupTestData()
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		// Pass a bad payment request ID
		randomID := uuid.Must(uuid.NewV4())
		_, err := TriggerEvent(Event{
			EventKey:        PaymentRequestCreateEventKey,
			MtoID:           paymentRequest.MoveTaskOrderID,
			UpdatedObjectID: randomID,
			EndpointKey:     SupportUpdatePaymentRequestStatusEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         uuid.Must(uuid.NewV4()),
		})
		// Check that at least one error was returned
		suite.NotNil(err)
		// Check that no notification was created
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count, newCount)

	})
}

func (suite *EventServiceSuite) Test_MTOEventTrigger() {

	// Test successful event
	suite.Run("Success with GHC MoveTaskOrder endpoint", func() {
		now := time.Now()
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})
		mtoID := mto.ID

		traceID := uuid.Must(uuid.NewV4())

		_, err := TriggerEvent(Event{
			EventKey:        MoveTaskOrderCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoID,
			EndpointKey:     GhcUpdateMoveTaskOrderEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})
		suite.Nil(err)

		// Get the notification
		notification, err := suite.getNotification(mtoID, traceID)
		suite.FatalNoError(err)
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
		// nolint:errcheck
		json.Unmarshal([]byte(notification.Payload), &mtoInPayload)
		// Check some params
		suite.Equal(mto.PPMType, &mtoInPayload.PpmType)
		suite.Equal(handlers.FmtDateTimePtr(mto.AvailableToPrimeAt).String(), mtoInPayload.AvailableToPrimeAt.String())

	})
}

func (suite *EventServiceSuite) Test_MTOShipmentEventTrigger() {
	// Test successful event passing with Support API
	suite.Run("Success with GHC MTOShipment endpoint", func() {
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: swag.Time(time.Now()),
			},
		})

		mtoShipmentID := mtoShipment.ID
		mtoID := mtoShipment.MoveTaskOrderID

		traceID := uuid.Must(uuid.NewV4())

		_, err := TriggerEvent(Event{
			EventKey:        MTOShipmentUpdateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoShipmentID,
			EndpointKey:     GhcApproveShipmentEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})
		suite.NoError(err)

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
		// nolint:errcheck
		json.Unmarshal([]byte(notification.Payload), &mtoShipmentInPayload)
		// Check some params
		suite.EqualValues(mtoShipment.ShipmentType, mtoShipmentInPayload.ShipmentType)
		suite.EqualValues(handlers.FmtDatePtr(mtoShipment.RequestedPickupDate).String(), mtoShipmentInPayload.RequestedPickupDate.String())
		suite.Nil(mtoShipment.NTSRecordedWeight)
	})

	suite.Run("No notification for GHC MTOShipment endpoint when shipment uses external vendor", func() {
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				UsesExternalVendor: true,
			},
			Move: models.Move{
				AvailableToPrimeAt: swag.Time(time.Now()),
			},
		})

		mtoShipmentID := mtoShipment.ID
		mtoID := mtoShipment.MoveTaskOrderID
		traceID := uuid.Must(uuid.NewV4())

		count, err := suite.DB().Count(&models.WebhookNotification{})
		suite.NoError(err)

		_, err = TriggerEvent(Event{
			EventKey:        MTOShipmentUpdateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoShipmentID,
			EndpointKey:     GhcApproveShipmentEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})
		suite.NoError(err)

		// Get the notification
		newCount, err := suite.DB().Count(&models.WebhookNotification{})
		suite.NoError(err)
		suite.Equal(count, newCount)
	})

	// Test successful event passing with Support API
	suite.Run("Success with GHC MTOShipment endpoint for NTS Shipment", func() {
		ntsRecordedWeight := unit.Pound(6989)
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: swag.Time(time.Now()),
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:      models.MTOShipmentTypeHHGIntoNTSDom,
				NTSRecordedWeight: &ntsRecordedWeight,
			},
		})

		mtoShipmentID := mtoShipment.ID
		mtoID := mtoShipment.MoveTaskOrderID

		traceID := uuid.Must(uuid.NewV4())

		_, err := TriggerEvent(Event{
			EventKey:        MTOShipmentUpdateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoShipmentID,
			EndpointKey:     GhcApproveShipmentEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})
		suite.Nil(err)

		// Get the notification
		notification, err := suite.getNotification(mtoShipmentID, traceID)
		suite.NoError(err)
		suite.Equal(&mtoShipmentID, notification.ObjectID)

		// Reinflate the json from the notification payload
		suite.NotEmpty(notification.Payload)
		var mtoShipmentInPayload primemessages.MTOShipment
		unmarshallErr := json.Unmarshal([]byte(notification.Payload), &mtoShipmentInPayload)
		suite.NoError(unmarshallErr)
		// Check some params
		suite.EqualValues(mtoShipment.ShipmentType, mtoShipmentInPayload.ShipmentType)
		suite.NotNil(mtoShipment.RequestedPickupDate)
		suite.NotNil(mtoShipmentInPayload.RequestedPickupDate)
		suite.Equal(ntsRecordedWeight, *mtoShipment.NTSRecordedWeight)
		suite.EqualValues(handlers.FmtDatePtr(mtoShipment.RequestedPickupDate).String(), mtoShipmentInPayload.RequestedPickupDate.String())
		storageFacility := *mtoShipment.StorageFacility
		suite.Equal(storageFacility.FacilityName, mtoShipmentInPayload.StorageFacility.FacilityName)
		suite.Equal(storageFacility.LotNumber, mtoShipmentInPayload.StorageFacility.LotNumber)
		suite.Equal(storageFacility.Address.StreetAddress1, *mtoShipmentInPayload.StorageFacility.Address.StreetAddress1)
		suite.Equal(storageFacility.Address.State, *mtoShipmentInPayload.StorageFacility.Address.State)
		suite.Equal(storageFacility.Address.City, *mtoShipmentInPayload.StorageFacility.Address.City)
		suite.Equal(storageFacility.Address.PostalCode, *mtoShipmentInPayload.StorageFacility.Address.PostalCode)
	})

	// Test successful no event passing with Support API when shipment is assigned to external vendor
	suite.Run("Error with GHC MTOShipment endpoint for NTS Shipment using external vendor", func() {

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: swag.Time(time.Now()),
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGIntoNTSDom,
				UsesExternalVendor: true,
			},
		})

		mtoShipmentID := mtoShipment.ID
		mtoID := mtoShipment.MoveTaskOrderID

		traceID := uuid.Must(uuid.NewV4())

		_, err := TriggerEvent(Event{
			EventKey:        MTOShipmentUpdateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoShipmentID,
			EndpointKey:     GhcApproveShipmentEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})
		suite.NoError(err)

		// Get the notification
		notification, err := suite.getNotification(mtoShipmentID, traceID)
		suite.Equal("sql: no rows in result set", err.Error())
		suite.Equal(uuid.NullUUID{}.UUID, *notification.ObjectID)

		// Reinflate the json from the notification payload
		suite.Empty(notification.Payload)
	})
}

func (suite *EventServiceSuite) Test_MTOServiceItemEventTrigger() {

	// Test successful event passing with Support API
	suite.Run("Success with GHC ServiceItem endpoint", func() {
		now := time.Now()
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})

		mtoServiceItemID := mtoServiceItem.ID
		mtoID := mtoServiceItem.MoveTaskOrderID

		traceID := uuid.Must(uuid.NewV4())
		count, _ := suite.DB().Count(&models.WebhookNotification{})

		_, err := TriggerEvent(Event{
			EventKey:        MTOServiceItemCreateEventKey,
			MtoID:           mtoID,
			UpdatedObjectID: mtoServiceItemID,
			EndpointKey:     GhcCreateMTOServiceItemEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})

		suite.Nil(err)
		newCount, _ := suite.DB().Count(&models.WebhookNotification{})
		suite.Equal(count+1, newCount)

	})
}

func (suite *EventServiceSuite) TestOrderEventTrigger() {

	// Test successful event passing with Support API
	suite.Run("Success with GHC ServiceItem endpoint", func() {

		move := testdatagen.MakeAvailableMove(suite.DB())
		traceID := uuid.Must(uuid.NewV4())
		_, err := TriggerEvent(Event{
			EventKey:        OrderUpdateEventKey,
			MtoID:           move.ID,
			UpdatedObjectID: move.OrdersID,
			EndpointKey:     InternalUpdateOrdersEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
		})
		suite.Nil(err)

		// Get the notification
		notification, err := suite.getNotification(move.OrdersID, traceID)
		suite.FatalNoError(err)
		suite.Equal(&move.OrdersID, notification.ObjectID)

		// Reinflate the json from the notification payload
		suite.NotEmpty(notification.Payload)
		var orderPayload primemessages.Order
		err = json.Unmarshal([]byte(notification.Payload), &orderPayload)
		suite.FatalNoError(err)

		// Check some params
		suite.Equal(move.Orders.ServiceMember.ID.String(), orderPayload.Customer.ID.String())
		suite.Equal(move.Orders.Entitlement.ID.String(), orderPayload.Entitlement.ID.String())
		suite.Equal(move.Orders.OriginDutyLocation.ID.String(), orderPayload.OriginDutyLocation.ID.String())
	})
}

func (suite *EventServiceSuite) TestNotificationEventHandler() {

	// Test a nil MTO ID is present and no notification stored
	suite.Run("No move and notification stored", func() {
		order := testdatagen.MakeDefaultOrder(suite.DB())
		traceID := uuid.Must(uuid.NewV4())
		count, _ := suite.DB().Count(&models.WebhookNotification{})
		event := Event{
			EventKey:        OrderUpdateEventKey,
			MtoID:           uuid.Nil,
			UpdatedObjectID: order.ID,
			EndpointKey:     InternalUpdateOrdersEndpointKey,
			AppContext:      suite.AppContextForTest(),
			TraceID:         traceID,
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
