package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	mto := testdatagen.MakeAvailableMove(suite.DB())
	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: models.MTOShipment{},
	})

	mtoShipment.MoveTaskOrderID = mto.ID

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	moveRouter := moveservices.NewMoveRouter()

	req := httptest.NewRequest("POST", "/mto-shipments", nil)

	params := mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest: req,
		Body: &primemessages.CreateMTOShipment{
			MoveTaskOrderID:      handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:               nil,
			CustomerRemarks:      nil,
			PointOfContact:       "John Doe",
			PrimeEstimatedWeight: 1200,
			RequestedPickupDate:  handlers.FmtDatePtr(mtoShipment.RequestedPickupDate),
			ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
		},
	}
	params.Body.DestinationAddress.Address = primemessages.Address{
		City:           &destinationAddress.City,
		Country:        destinationAddress.Country,
		PostalCode:     &destinationAddress.PostalCode,
		State:          &destinationAddress.State,
		StreetAddress1: &destinationAddress.StreetAddress1,
		StreetAddress2: destinationAddress.StreetAddress2,
		StreetAddress3: destinationAddress.StreetAddress3,
	}
	params.Body.PickupAddress.Address = primemessages.Address{
		City:           &pickupAddress.City,
		Country:        pickupAddress.Country,
		PostalCode:     &pickupAddress.PostalCode,
		State:          &pickupAddress.State,
		StreetAddress1: &pickupAddress.StreetAddress1,
		StreetAddress2: pickupAddress.StreetAddress2,
		StreetAddress3: pickupAddress.StreetAddress3,
	}

	suite.T().Run("Successful POST - Integration Test", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}
		response := handler.Handle(params)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		// check that the mto shipment status is Submitted
		suite.Require().Equal(createMTOShipmentPayload.Status, primemessages.MTOShipmentStatusSUBMITTED, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.PrimeEstimatedWeight, params.Body.PrimeEstimatedWeight)
	})

	suite.T().Run("POST failure - 500", func(t *testing.T) {
		mockCreator := mocks.MTOShipmentCreator{}

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			&mockCreator,
			mtoChecker,
		}

		err := errors.New("ServerError")

		mockCreator.On("CreateMTOShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.T().Run("POST failure - 422 -- Bad agent IDs set on shipment", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		badID := params.Body.MoveTaskOrderID
		agent := &primemessages.MTOAgent{
			ID:            *badID,
			MtoShipmentID: *badID,
			FirstName:     handlers.FmtString("Mary"),
		}

		paramsBadIDs := params
		paramsBadIDs.Body.Agents = primemessages.MTOAgents{agent}

		response := handler.Handle(paramsBadIDs)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)
		suite.NotEmpty(typedResponse.Payload.InvalidFields)
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:id")
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:mtoShipmentID")
	})

	suite.T().Run("POST failure - 422 - invalid input, missing pickup address", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		badParams := params
		badParams.Body.PickupAddress.Address.StreetAddress1 = nil

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("POST failure - 404 -- not found", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("POST failure - 400 -- nil body", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		paramsNilBody := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}
		response := handler.Handle(paramsNilBody)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("POST failure - 404 -- MTO is not available to Prime", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		mtoNotAvailable := testdatagen.MakeDefaultMove(suite.DB())
		mtoIDNotAvailable := strfmt.UUID(mtoNotAvailable.ID.String())

		paramsNotAvailable := params
		paramsNotAvailable.Body.MoveTaskOrderID = &mtoIDNotAvailable

		response := handler.Handle(paramsNotAvailable)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)

		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentNotFound)
		suite.Contains(*typedResponse.Payload.Detail, mtoNotAvailable.ID.String())
	})

	suite.T().Run("POST failure - 422 - modelType() not supported", func(t *testing.T) {
		mockCreator := mocks.MTOShipmentCreator{}

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}
		err := apperror.NotFoundError{}

		mockCreator.On("CreateMTOShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		mtoServiceItems := models.MTOServiceItems{
			models.MTOServiceItem{
				MoveTaskOrderID:  mto.ID,
				MTOShipmentID:    &mtoShipment.ID,
				ReService:        models.ReService{Code: models.ReServiceCodeMS},
				Reason:           nil,
				PickupPostalCode: nil,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}

		badModelTypeParams := params
		badModelTypeParams.Body.SetMtoServiceItems(*payloads.MTOServiceItems(&mtoServiceItems))
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})
}

// ClearNonUpdateFields clears out the MTOShipment payload fields that CANNOT be sent in for a successful update
func ClearNonUpdateFields(mtoShipment *models.MTOShipment) *primemessages.MTOShipment {
	mtoShipment.MoveTaskOrderID = uuid.FromStringOrNil("")
	mtoShipment.CreatedAt = time.Time{}
	mtoShipment.UpdatedAt = time.Time{}
	mtoShipment.PrimeEstimatedWeightRecordedDate = &time.Time{}
	mtoShipment.RequiredDeliveryDate = &time.Time{}
	mtoShipment.ApprovedDate = &time.Time{}
	mtoShipment.Status = ""
	mtoShipment.RejectionReason = nil
	mtoShipment.CustomerRemarks = nil
	mtoShipment.MTOAgents = nil

	return payloads.MTOShipment(mtoShipment)
}

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {

	// Create an available shipment in DB
	now := testdatagen.CurrentDateWithoutTime()
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:       models.MTOShipmentStatusApproved,
			ApprovedDate: now,
		},
		SecondaryPickupAddress:   testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{}),
		SecondaryDeliveryAddress: testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{}),
	})
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)

	// Create a minimal shipment
	minimalShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:              models.MTOShipmentStatusDraft,
			ScheduledPickupDate: &time.Time{},
		},
	})
	minimalReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", minimalShipment.ID.String()), nil)

	// Create some usable weights
	primeEstimatedWeight := unit.Pound(500)
	primeActualWeight := unit.Pound(600)

	// CREATE HANDLER OBJECT

	// ghcDomesticTime is used in the planner, the planner checks transit distance.
	// We mock the planner to return 400, so we need an entry that will return a
	// transit time of 12 days for a distance of 400.

	// Mock planner to always return a distance of 400 mi
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	// Add a 12 day transit time for a distance of 400
	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 1,
		DistanceMilesUpper: 500,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator)
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	suite.T().Run("PATCH failure 500 Unit Test", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     MTOShipmentUpdater, Planner
		// Set up:     We provide an update but make MTOShipmentUpdater return a server error
		// Expected:   Handler returns Internal Server Error Response. This ensures if there is an
		//             unexpected error in the service object, we return the proper HTTP response

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				Diversion: true,
			},
			IfMatch: eTag,
		}

		mockUpdater := mocks.MTOShipmentUpdater{}
		mockHandler := UpdateMTOShipmentHandler{
			context,
			&mockUpdater,
		}
		internalServerErr := errors.New("ServerError")

		mockUpdater.On("MTOShipmentsMTOAvailableToPrime",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(true, nil)

		mockUpdater.On("UpdateMTOShipmentPrime",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")

	})

	suite.T().Run("PATCH success 200 minimal update", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We use the normal (non-minimal) shipment we created earlier
		//             We provide an update with minimal changes
		// Expected:   Handler returns OK
		//             Minimal updates are completed, old values retained for rest of
		//             shipment. This tests that PATCH is not accidentally clearing any existing
		//             data.

		// Create an update with just diversion and actualPickupDate
		now := testdatagen.CurrentDateWithoutTime()
		minimalUpdate := primemessages.UpdateMTOShipment{
			Diversion:        true,
			ActualPickupDate: handlers.FmtDatePtr(now),
		}

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &minimalUpdate,
			IfMatch:       eTag,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload
		suite.Equal(shipment.ID.String(), okPayload.ID.String())

		// Confirm PATCH working as expected; non-updated values still exist
		suite.EqualDatePtr(shipment.ApprovedDate, okPayload.ApprovedDate)
		suite.EqualDatePtr(shipment.FirstAvailableDeliveryDate, okPayload.FirstAvailableDeliveryDate)
		suite.EqualDatePtr(shipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		suite.EqualDatePtr(shipment.RequiredDeliveryDate, okPayload.RequiredDeliveryDate)
		suite.EqualDatePtr(shipment.ScheduledPickupDate, okPayload.ScheduledPickupDate)

		suite.EqualAddress(*shipment.PickupAddress, &okPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &okPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &okPayload.SecondaryDeliveryAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &okPayload.SecondaryPickupAddress.Address, true)

		// Confirm new values
		suite.Equal(params.Body.Diversion, okPayload.Diversion)
		suite.Equal(params.Body.ActualPickupDate.String(), okPayload.ActualPickupDate.String())

		// Refresh local copy of shipment from DB for etag regeneration in future tests
		shipment = suite.refreshFromDB(shipment.ID)

	})

	suite.T().Run("PATCH success 200 update destination type", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We use the normal (non-minimal) shipment we created earlier
		//             We provide an update with minimal changes
		// Expected:   Handler returns OK
		//             Minimal updates are completed, old values retained for rest of
		//             shipment. This tests that PATCH is not accidentally clearing any existing
		//             data.

		// Create an update with just destinationAddressType
		destinationType := primemessages.DestinationTypeHOMEOFRECORD
		minimalUpdate := primemessages.UpdateMTOShipment{
			DestinationType: &destinationType,
		}

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &minimalUpdate,
			IfMatch:       eTag,
		}
		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload
		suite.Equal(shipment.ID.String(), okPayload.ID.String())

		// Confirm PATCH working as expected; non-updated values still exist
		suite.EqualDatePtr(shipment.ApprovedDate, okPayload.ApprovedDate)
		suite.EqualDatePtr(shipment.FirstAvailableDeliveryDate, okPayload.FirstAvailableDeliveryDate)
		suite.EqualDatePtr(shipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		suite.EqualDatePtr(shipment.RequiredDeliveryDate, okPayload.RequiredDeliveryDate)
		suite.EqualDatePtr(shipment.ScheduledPickupDate, okPayload.ScheduledPickupDate)

		suite.EqualAddress(*shipment.PickupAddress, &okPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &okPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &okPayload.SecondaryDeliveryAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &okPayload.SecondaryPickupAddress.Address, true)

		// Confirm new value
		suite.Equal(params.Body.DestinationType, okPayload.DestinationType)

		// Refresh local copy of shipment from DB for etag regeneration in future tests
		shipment = suite.refreshFromDB(shipment.ID)

	})

	suite.T().Run("PATCH failure 404 not found because not available to prime", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We provide an update to a shipment whose associated move isn't available to prime
		// Expected:   Handler returns Not Found error

		// Create a shipment unavailable to Prime in DB
		shipmentNotAvailable := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		suite.Nil(shipmentNotAvailable.MoveTaskOrder.AvailableToPrimeAt)

		// Create params
		notAvReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipmentNotAvailable.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   notAvReq,
			MtoShipmentID: *handlers.FmtUUID(shipmentNotAvailable.ID),
			Body: &primemessages.UpdateMTOShipment{
				Diversion: true,
			},
			IfMatch: etag.GenerateEtag(shipmentNotAvailable.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Verify not found response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PATCH failure 404 not found because attempting to update an external vendor shipment", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We provide an update to a shipment that is handled by an external vendor
		// Expected:   Handler returns Not Found error

		// Create a shipment handled by an external vendor
		externalShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
			},
		})
		suite.NotNil(externalShipment.MoveTaskOrder.AvailableToPrimeAt)

		// Create params
		notAvReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", externalShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   notAvReq,
			MtoShipmentID: *handlers.FmtUUID(externalShipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				Diversion: true,
			},
			IfMatch: etag.GenerateEtag(externalShipment.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Verify not found response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PATCH success 200 update of primeEstimatedWeight and primeActualWeight", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We provide an update with actual and estimated weights
		// Expected:   Handler returns OK
		//             Weights are updated, and prime estimated weight recorded date is updated.

		// Check that PrimeEstimatedWeightRecordedDate was nil at first
		suite.Nil(minimalShipment.PrimeEstimatedWeightRecordedDate)

		// Update the primeEstimatedWeight
		eTag := etag.GenerateEtag(minimalShipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   minimalReq,
			MtoShipmentID: *handlers.FmtUUID(minimalShipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: int64(primeEstimatedWeight), // New estimated weight
				PrimeActualWeight:    int64(primeActualWeight),    // New actual weight
			},
			IfMatch: eTag,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload
		suite.Equal(minimalShipment.ID.String(), okPayload.ID.String())

		// Confirm changes to weights
		suite.Equal(int64(primeActualWeight), okPayload.PrimeActualWeight)
		suite.Equal(int64(primeEstimatedWeight), okPayload.PrimeEstimatedWeight)
		// Confirm primeEstimatedWeightRecordedDate was added
		suite.NotNil(okPayload.PrimeEstimatedWeightRecordedDate)
		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okPayload.RequestedPickupDate)
		suite.EqualDatePtr(minimalShipment.RequestedPickupDate, okPayload.RequestedPickupDate)

		// refresh shipment from DB for getting the updated eTag
		minimalShipment = suite.refreshFromDB(minimalShipment.ID)

	})
	suite.T().Run("PATCH failure 422 cannot update primeEstimatedWeight again", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Use previously created shipment with primeEstimatedWeight updated
		//             Attempt to update primeEstimatedWeight
		// Expected:   Handler returns Unprocessable Entity
		//             primeEstimatedWeight cannot be updated more than once.

		// Check that primeEstimatedWeight was already populated
		suite.NotNil(minimalShipment.PrimeEstimatedWeight)

		// Attempt to update again
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   minimalReq,
			MtoShipmentID: *handlers.FmtUUID(minimalShipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: int64(primeEstimatedWeight + 100), // New estimated weight
			},
			IfMatch: etag.GenerateEtag(minimalShipment.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response contains an error about primeEstimatedWeight
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errPayload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload
		suite.Contains(errPayload.InvalidFields, "primeEstimatedWeight")

	})

	suite.T().Run("PATCH failure 422 cannot update estimatedWeight without scheduledPickupDate", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Create a shipment with no scheduledPickupDate
		//             Attempt to update the primeEstimatedWeight
		// Expected:   Handler returns Unprocessable entity because the
		//             primeEstimatedWeight cannot be set if the scheduledPickupDate is not set.

		// Create a shipment with no scheduled pickup date
		noScheduledPickupShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusSubmitted,
				ScheduledPickupDate: nil,
			},
		})

		// Create an update with updated weights
		mtoShipment := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			PrimeActualWeight:    int64(primeActualWeight),
		}

		// Create request to UpdateMTOShipment
		noPickupReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", noScheduledPickupShipment.ID.String()), nil)
		noPickupParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   noPickupReq,
			MtoShipmentID: *handlers.FmtUUID(noScheduledPickupShipment.ID),
			Body:          &mtoShipment,
			IfMatch:       etag.GenerateEtag(noScheduledPickupShipment.UpdatedAt),
		}

		suite.NoError(noPickupParams.Body.Validate(strfmt.Default))
		response := handler.Handle(noPickupParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)

		// Expect validation error due to updated weight
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.NotEmpty(errResponse.Payload.InvalidFields)
		suite.Contains(errResponse.Payload.InvalidFields, "primeEstimatedWeight")
	})

	suite.T().Run("PATCH failure 404 unknown shipment", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with fake uuid
		// Expected:   Handler returns Not Found error

		// Create request with non existent ID
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(uuid.Must(uuid.NewV4()).String()), // generate a UUID
			Body:          &primemessages.UpdateMTOShipment{},
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}
		// Call handler
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PATCH failure 412 precondition failed", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with old eTag
		// Expected:   Handler returns Precondition Failed error

		// Create an update with an old eTag
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(shipment.ID.String()),
			Body:          &primemessages.UpdateMTOShipment{Diversion: true}, // update anything
			IfMatch:       string(etag.GenerateEtag(shipment.CreatedAt)),     // Use createdAt to generate eTag
		}

		// Call handler
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.T().Run("PATCH success 200 returns all nested objects", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We add service items to the shipment in the DB
		//             We provide an almost empty update so as to check that the
		//             nested objects in the response are fully populated
		// Expected:   Handler returns OK, all service items, agents and addresses are
		//             populated.

		// Add service items to our shipment
		// Create a service item in the db, associate with the shipment
		reService := testdatagen.MakeDDFSITReService(suite.DB())
		testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move:        move,
			MTOShipment: shipment,
			ReService:   reService,
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: move.ID,
				ReServiceID:     reService.ID,
				MTOShipmentID:   &shipment.ID,
			},
		})

		// Add agents associated to our shipment
		agent1 := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				FirstName:    swag.String("Test1"),
				LastName:     swag.String("Agent"),
				Email:        swag.String("test@test.email.com"),
				MTOAgentType: models.MTOAgentReceiving,
			},
			MTOShipment: shipment,
		})
		agent2 := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				FirstName:    swag.String("Test2"),
				LastName:     swag.String("Agent"),
				Email:        swag.String("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
			MTOShipment: shipment,
		})

		// Create an almost empty update
		// We only want to see the response payload to make sure it is populated correctly
		update := primemessages.UpdateMTOShipment{
			PointOfContact: "John McRand",
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &update,
			IfMatch:       etag.GenerateEtag(shipment.UpdatedAt),
		}

		// Call the handler
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Check that there's one service item of model type DestSIT in the payload
		suite.Equal(1, len(okPayload.MtoServiceItems()))
		serviceItem := okPayload.MtoServiceItems()[0]
		suite.Equal(primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT, serviceItem.ModelType())

		// Check the reServiceCode string
		serviceItemDestSIT := serviceItem.(*primemessages.MTOServiceItemDestSIT)
		suite.Equal(string(reService.Code), *serviceItemDestSIT.ReServiceCode)

		// Check that there's 2 agents, then check them against the ones we created
		suite.Equal(2, len(okPayload.Agents))
		for _, item := range okPayload.Agents {
			if item.AgentType == primemessages.MTOAgentType(agent1.MTOAgentType) {
				suite.Equal(agent1.FirstName, item.FirstName)
			}
			if item.AgentType == primemessages.MTOAgentType(agent2.MTOAgentType) {
				suite.Equal(agent2.FirstName, item.FirstName)
			}
		}

		// Check all dates and addresses in the payload
		suite.EqualDatePtr(shipment.ApprovedDate, okPayload.ApprovedDate)
		suite.EqualDatePtr(shipment.FirstAvailableDeliveryDate, okPayload.FirstAvailableDeliveryDate)
		suite.EqualDatePtr(shipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		suite.EqualDatePtr(shipment.RequiredDeliveryDate, okPayload.RequiredDeliveryDate)
		suite.EqualDatePtr(shipment.ScheduledPickupDate, okPayload.ScheduledPickupDate)

		suite.EqualAddress(*shipment.PickupAddress, &okPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &okPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &okPayload.SecondaryDeliveryAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &okPayload.SecondaryPickupAddress.Address, true)

	})
}

// TestUpdateMTOShipmentAddressLogic tests the create/update address logic
// This endpoint can create but not update addresses due to optimistic locking
func (suite *HandlerSuite) TestUpdateMTOShipmentAddressLogic() {

	// Create a shipment in the DB that has no addresses populated:
	now := time.Now()
	shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Status:             "APPROVED",
		}, // Prime-available move
	})
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", shipment.ID.String()), nil)

	// CREATE HANDLER OBJECT
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator)
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	suite.T().Run("PATCH success 200 create addresses", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle, addresses mechanism - we can create but not update
		// Mocked:     Planner
		// Set up:     We use a shipment with minimal info, no addresses
		//             Update with PickupAddress, DestinationAddress, SecondaryPickupAddress, SecondaryDeliveryAddress
		// Expected:   Handler should return OK, new addresses created

		// CREATE REQUEST
		// Create an update message with all addresses provided
		update := primemessages.UpdateMTOShipment{
			PickupAddress:            getFakeAddress(),
			DestinationAddress:       getFakeAddress(),
			SecondaryPickupAddress:   getFakeAddress(),
			SecondaryDeliveryAddress: getFakeAddress(),
		}
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &update,
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Check that addresses match what was sent
		suite.EqualAddressPayload(&update.PickupAddress.Address, &okPayload.PickupAddress.Address, false)
		suite.EqualAddressPayload(&update.DestinationAddress.Address, &okPayload.DestinationAddress.Address, false)
		suite.EqualAddressPayload(&update.SecondaryPickupAddress.Address, &okPayload.SecondaryPickupAddress.Address, false)
		suite.EqualAddressPayload(&update.SecondaryDeliveryAddress.Address, &okPayload.SecondaryDeliveryAddress.Address, false)

		shipment = suite.refreshFromDB(shipment.ID)

	})

	suite.T().Run("PATCH failure 422 update addresses not allowed", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle, addresses mechanism - we cannot update addresses
		// Mocked:     Planner
		// Set up:     We use the previous shipment with all addresses populated
		//             Update with PickupAddress, DestinationAddress, SecondaryPickupAddress, SecondaryDeliveryAddress
		// Expected:   Handler should return unprocessable entity error. All addresses cannot be updated and should
		//             be listed in errors

		// CREATE REQUEST
		// Create an update message with all new addresses provided
		update := primemessages.UpdateMTOShipment{
			PickupAddress:            getFakeAddress(),
			DestinationAddress:       getFakeAddress(),
			SecondaryPickupAddress:   getFakeAddress(),
			SecondaryDeliveryAddress: getFakeAddress(),
		}
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &update,
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errPayload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload
		suite.Contains(errPayload.InvalidFields, "pickupAddress")
		suite.Contains(errPayload.InvalidFields, "destinationAddress")
		suite.Contains(errPayload.InvalidFields, "secondaryPickupAddress")
		suite.Contains(errPayload.InvalidFields, "secondaryDeliveryAddress")

	})

	suite.T().Run("PATCH success 200 nil doesn't clear addresses", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle, addresses mechanism - we can create but not update
		// Mocked:     Planner
		// Set up:     We use the previous shipment with addresses populated.
		//             Update with nil for the addresses.
		// Expected:   Handler should return OK, addresses should be unchanged.
		//             This endpoint was previously blanking out addresses which is why we have this test.

		// CREATE REQUEST
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.UpdateMTOShipment{}, // Empty payload
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		newPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Check that addresses match what was returned previously in the last successful payload
		suite.EqualAddress(*shipment.PickupAddress, &newPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &newPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &newPayload.SecondaryPickupAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &newPayload.SecondaryDeliveryAddress.Address, true)
	})
}

// TestUpdateMTOShipmentDateLogic tests various restrictions related to timelines that
// Prime is required to abide by
// More details about these rules can be found in the Performance Work Statement for the
// Global Household Goods Contract HTC711-19-R-R004
func (suite *HandlerSuite) TestUpdateMTOShipmentDateLogic() {

	// Create an available move to be used for the shipments
	move := testdatagen.MakeAvailableMove(suite.DB())
	primeEstimatedWeight := unit.Pound(500)
	now := time.Now()

	// ghcDomesticTime is used in the planner, the planner checks transit distance.
	// We mock the planner to return 400, so we need an entry that will return a
	// transit time of 12 days for a distance of 400.

	// Mock planner to always return a distance of 400 mi
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	// Add a 12 day transit time for a distance of 400
	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 1,
		DistanceMilesUpper: 500,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	// Create a handler object to use in the tests
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator)

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	// Prime-specific validations tested below
	suite.T().Run("Failed case if not both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		eightDaysFromNow := now.AddDate(0, 0, 8)
		threeDaysBefore := now.AddDate(0, 0, -3)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &eightDaysFromNow,
				ApprovedDate:        &threeDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.Contains(errResponse.Payload.InvalidFields, "primeEstimatedWeight")

	})

	suite.T().Run("Successful case if both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &tenDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(oldShipment.ID.String(), okResponse.Payload.ID.String())

		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okResponse.Payload.RequestedPickupDate)
		suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))

	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate updated on scheduledPickupDate update", func(t *testing.T) {
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
		storageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{})

		hhgShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
				Status:       models.MTOShipmentStatusApproved,
				ApprovedDate: &now,
			},
		})

		ntsShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType:      models.MTOShipmentTypeHHGIntoNTSDom,
				Status:            models.MTOShipmentStatusApproved,
				StorageFacility:   &storageFacility,
				StorageFacilityID: &storageFacility.ID,
				PickupAddress:     &address,
				PickupAddressID:   &address.ID,
			},
		})

		NTSRecordedWeight := unit.Pound(1400)
		ntsrShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
				NTSRecordedWeight:    &NTSRecordedWeight,
				Status:               models.MTOShipmentStatusApproved,
				StorageFacility:      &storageFacility,
				StorageFacilityID:    &storageFacility.ID,
				DestinationAddress:   &address,
				DestinationAddressID: &address.ID,
			},
		})

		tenDaysFromNow := now.AddDate(0, 0, 11)
		schedDate := strfmt.Date(tenDaysFromNow)

		testCases := []struct {
			shipment models.MTOShipment
			payload  primemessages.UpdateMTOShipment
		}{
			{hhgShipment, primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: int64(primeEstimatedWeight),
				ScheduledPickupDate:  &schedDate,
			}},
			{ntsShipment, primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: int64(primeEstimatedWeight),
				ScheduledPickupDate:  &schedDate,
			}},
			{ntsrShipment, primemessages.UpdateMTOShipment{
				ScheduledPickupDate: &schedDate,
			}},
		}
		for _, testCase := range testCases {
			oldShipment := testCase.shipment
			eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

			params := mtoshipmentops.UpdateMTOShipmentParams{
				HTTPRequest:   req,
				MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
				Body:          &testCase.payload,
				IfMatch:       eTag,
			}

			response := handler.Handle(params)

			okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
			suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

			responsePayload := okResponse.Payload
			suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
			suite.NotNil(responsePayload.RequiredDeliveryDate)
			suite.NotNil(responsePayload.ScheduledPickupDate)

			// Let's double check our maths.
			expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 12)
			actualRDD := time.Time(*responsePayload.RequiredDeliveryDate)
			suite.Equal(expectedRDD.Year(), actualRDD.Year())
			suite.Equal(expectedRDD.Month(), actualRDD.Month())
			suite.Equal(expectedRDD.Day(), actualRDD.Day())

			// Confirm PATCH working as expected; non-updated value still exists
			if testCase.shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom { //ntsr doesn't have a RequestedPickupDate
				suite.NotNil(okResponse.Payload.RequestedPickupDate)
				suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))
			}
		}
	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate updated on destinationAddress creation", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle, RequiredDeliveryDate logic
		// Mocked:     Planner
		// Set up:     We use a shipment with primeEstimatedWeight and ScheduledPickupDate set
		//             Update with new destinationAddress
		// Expected:   Handler should return OK, new DestinationAddress should be saved
		//             requiredDeliveryDate should be set to 12 days from scheduledPickupDate

		// Create shipment with populated estimated weight and scheduled date
		tenDaysFromNow := now.AddDate(0, 0, 11)
		pickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
		oldShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:               "APPROVED",
				ApprovedDate:         &now,
				PrimeEstimatedWeight: &primeEstimatedWeight,
				ScheduledPickupDate:  &tenDaysFromNow,
				PickupAddress:        &pickupAddress,
			},
		})

		// CREATE REQUEST
		// Update destination address
		update := primemessages.UpdateMTOShipment{
			DestinationAddress: getFakeAddress(),
		}
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &update,
			IfMatch:       etag.GenerateEtag(oldShipment.UpdatedAt),
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		responsePayload := okResponse.Payload

		// Confirm destination address in payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.EqualAddressPayload(&update.DestinationAddress.Address, &responsePayload.DestinationAddress.Address, false)

		// Confirm that auto-generated requiredDeliveryDate matches expected value
		expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 12)
		suite.EqualDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate for Alaska", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle, RequiredDeliveryDate logic
		// Mocked:     Planner
		// Set up:     We use a shipment with an Alaska Address
		//             Update with new DestinationAddress
		// Expected:   Handler should return OK, new DestinationAddress should be saved
		//             requiredDeliveryDate should be set to 12 + 10 = 22 days from scheduledPickupDate
		//             which is a special rule for Alaska

		// Create shipment with Alaska destination
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:       "APPROVED",
				ApprovedDate: &now,
			},
			DestinationAddress: models.Address{
				PostalCode: "12345",
				State:      "AK",
			},
		})

		// CREATE REQUEST
		// Update with scheduledPickupDate and PrimeEstimatedWeight
		tenDaysFromNow := now.AddDate(0, 0, 11)
		schedDate := strfmt.Date(tenDaysFromNow)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			ScheduledPickupDate:  &schedDate,
		}

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload

		// Check that updated fields are saved
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
		suite.NotNil(responsePayload.ScheduledPickupDate)

		// Check that RDD is set to 12 + 10 days after scheduled pickup date
		expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 22)
		suite.EqualDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate for Adak, Alaska", func(t *testing.T) {
		// Under test: updateMTOShipmentHandler.Handle, RequiredDeliveryDate logic
		// Mocked:     Planner
		// Set up:     We use a shipment with an Alaska Address, specifically Adak
		//             Update with new DestinationAddress
		// Expected:   Handler should return OK, new DestinationAddress should be saved
		//             requiredDeliveryDate should be set to 12 + 20 = 32 days from scheduledPickupDate,
		//             which is a special rule for Adak (look at it on a map!)

		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:       "APPROVED",
				ApprovedDate: &now,
			},
			DestinationAddress: models.Address{
				PostalCode: "99546",
				City:       "Adak",
				State:      "AK",
			},
		})

		// CREATE REQUEST
		// Update with scheduledPickupDate and PrimeEstimatedWeight
		tenDaysFromNow := now.AddDate(0, 0, 11)
		schedDate := strfmt.Date(tenDaysFromNow)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			ScheduledPickupDate:  &schedDate,
		}

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload

		// Check that updated fields are saved
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
		suite.NotNil(responsePayload.ScheduledPickupDate)

		// Check that RDD is set to 12 + 20 days after scheduled pickup date
		expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 32)
		suite.EqualDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.T().Run("Failed case if approved date is 3-9 days from scheduled move date but estimated weight recorded date isn't at least 3 days prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.Contains(errResponse.Payload.InvalidFields, "primeEstimatedWeight")
	})

	suite.T().Run("Successful case if approved date is 3-9 days from scheduled move date and estimated weight recorded date is at least 3 days prior to scheduled move date", func(t *testing.T) {
		sixDaysFromNow := now.AddDate(0, 0, 6)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &sixDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okResponse.Payload.RequestedPickupDate)
		suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))
	})

	suite.T().Run("Failed case if approved date is less than 3 days from scheduled move date but estimated weight recorded date isn't at least 1 day prior to scheduled move date", func(t *testing.T) {
		oneDayFromNow := now.AddDate(0, 0, 1)
		oneDayBefore := now.AddDate(0, 0, -1)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &oneDayFromNow,
				ApprovedDate:        &oneDayBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.Contains(errResponse.Payload.InvalidFields, "primeEstimatedWeight")

	})

	suite.T().Run("Successful case if approved date is less than 3 days from scheduled move date and estimated weight recorded date is at least 1 day prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okResponse.Payload.RequestedPickupDate)
		suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))
	})

}

func (suite *HandlerSuite) TestUpdateMTOShipmentStatusHandler() {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetPlanner(planner)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	handler := UpdateMTOShipmentStatusHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator),
		mtoshipment.NewMTOShipmentStatusUpdater(builder,
			mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter), planner),
	}
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s/status", uuid.Nil.String()), nil)

	// Set up Prime-available move
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
			Status:          models.MTOShipmentStatusCancellationRequested,
		},
	})
	eTag := etag.GenerateEtag(shipment.UpdatedAt)

	suite.T().Run("200 SUCCESS - Updated CANCELLATION_REQUESTED to CANCELED", func(t *testing.T) {
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusOK)
		suite.Equal(string(models.MTOShipmentStatusCanceled), okResponse.Payload.Status)
		suite.Equal(move.ID.String(), okResponse.Payload.MoveTaskOrderID.String())
		suite.NotZero(okResponse.Payload.ETag)

		eTag = okResponse.Payload.ETag // updated for following tests
	})

	suite.T().Run("404 FAIL - Bad shipment ID", func(t *testing.T) {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000010")
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(badUUID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusNotFound{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusNotFound)
		suite.Contains(*errResponse.Payload.Detail, badUUID.String())
	})

	suite.T().Run("404 FAIL - Shipment was not Prime-available", func(t *testing.T) {
		nonPrimeShipment := testdatagen.MakeDefaultMTOShipment(suite.DB()) // default is non-Prime available
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(nonPrimeShipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       etag.GenerateEtag(nonPrimeShipment.UpdatedAt),
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusNotFound{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusNotFound)
		suite.Contains(*errResponse.Payload.Detail, nonPrimeShipment.ID.String())
	})

	suite.T().Run("412 FAIL - Stale eTag", func(t *testing.T) {
		staleShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrder:   move,
				MoveTaskOrderID: move.ID,
				Status:          models.MTOShipmentStatusCancellationRequested,
			},
		})
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(staleShipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       "eTag",
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusPreconditionFailed{}, response)
	})

	suite.T().Run("409 FAIL - Current status was not CANCELLATION_REQUESTED", func(t *testing.T) {
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID), // This shipment currently has CANCELED status already
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       eTag,
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusConflict{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusConflict)
		suite.Contains(*errResponse.Payload.Detail, string(models.MTOShipmentStatusCanceled))
	})

	suite.T().Run("422 FAIL - Tried to use a status other than CANCELED", func(t *testing.T) {
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusApproved)},
			IfMatch:       eTag,
		}
		// Run swagger validations - should fail
		suite.Error(params.Body.Validate(strfmt.Default))
	})
}

func getFakeAddress() struct{ primemessages.Address } {
	// Use UUID to generate truly random address string
	streetAddr := fmt.Sprintf("%s %s", uuid.Must(uuid.NewV4()).String(), fakedata.RandomStreetAddress())
	// Using same zip so not a good helper for tests testing zip calculations
	return struct{ primemessages.Address }{
		Address: primemessages.Address{
			City:           swag.String("San Diego"),
			PostalCode:     swag.String("92102"),
			State:          swag.String("CA"),
			StreetAddress1: &streetAddr,
		},
	}
}

func (suite *HandlerSuite) refreshFromDB(id uuid.UUID) models.MTOShipment {
	var dbShipment models.MTOShipment
	err := suite.DB().EagerPreload("PickupAddress",
		"DestinationAddress",
		"SecondaryPickupAddress",
		"SecondaryDeliveryAddress",
		"MTOAgents").Find(&dbShipment, id)
	suite.Nil(err)
	return dbShipment
}
