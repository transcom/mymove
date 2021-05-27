package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	"github.com/transcom/mymove/pkg/services"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/swag"

	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	_ "github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
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

	builder := query.NewQueryBuilder(suite.DB())
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

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
			ShipmentType:         primemessages.MTOShipmentTypeHHG,
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}

		err := errors.New("ServerError")

		mockCreator.On("CreateMTOShipment",
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}
		err := services.NotFoundError{}

		mockCreator.On("CreateMTOShipment",
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

	primeEstimatedWeight := unit.Pound(500)
	primeActualWeight := unit.Pound(600)

	// Create an available shipment in DB
	now := time.Now()
	move := testdatagen.MakeAvailableMove(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:       models.MTOShipmentStatusApproved,
			ApprovedDate: &now,
		},
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

	// CREATE HANDLER OBJECT
	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	suite.T().Run("PATCH failure 500 Unit Test", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function
		// Mocked:     MTOShipmentUpdater, Planner
		// Set up:     We provide an update but make MTOShipmentUpdater return
		//             internal server error
		// Expected outcome:
		//             Handler returns Internal Server Error Response

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body: &primemessages.MTOShipment{
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
			mock.Anything,
		).Return(true, nil)

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")

	})

	suite.T().Run("PATCH success 200 minimal update", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function
		// Mocked:     None
		// Set up:     We use the normal (non-minimal) shipment we created earlier
		//             We provide an update with minimal changes
		// Expected outcome:
		//             Handler returns OK
		//             Minimal updates are completed, old values retained for rest of
		//             shipment
		now := time.Now()
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		minimalUpdate := primemessages.MTOShipment{
			Diversion:        true,
			ActualPickupDate: handlers.FmtDatePtr(&now),
		}
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
		suite.equalDatePtr(shipment.ApprovedDate, okPayload.ApprovedDate)
		suite.equalDatePtr(shipment.FirstAvailableDeliveryDate, okPayload.FirstAvailableDeliveryDate)
		suite.equalDatePtr(shipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		suite.equalDatePtr(shipment.RequiredDeliveryDate, okPayload.RequiredDeliveryDate)
		suite.equalDatePtr(shipment.ScheduledPickupDate, okPayload.ScheduledPickupDate)

		suite.EqualAddress(*shipment.PickupAddress, &okPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &okPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &okPayload.SecondaryDeliveryAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &okPayload.SecondaryPickupAddress.Address, true)

		// Confirm new values
		suite.Equal(params.Body.Diversion, okPayload.Diversion)
		suite.Equal(params.Body.ActualPickupDate.String(), okPayload.ActualPickupDate.String())

		// Refresh local copy of shipment from DB
		shipment = suite.refreshFromDB(shipment.ID)

	})

	suite.T().Run("PATCH failure 404 not found because not available to prime", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     We provide an update with minimal changes to a shipment
		//             whose associated move isn't available to prime
		// Expected outcome:
		//             Handler returns Not Found error

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
			Body: &primemessages.MTOShipment{
				Diversion: true,
			},
			IfMatch: etag.GenerateEtag(shipmentNotAvailable.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PATCH success 200 update of primeEstimatedWeight and primeActualWeight", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     We provide an update with actual and estimated weights
		// Expected outcome:
		//             Handler returns OK
		//             Weights are updated, and prime estimated weight recorded date is updated.

		eTag := etag.GenerateEtag(minimalShipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   minimalReq,
			MtoShipmentID: *handlers.FmtUUID(minimalShipment.ID),
			Body: &primemessages.MTOShipment{
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
		// Confirm new date was added
		suite.NotNil(okPayload.PrimeEstimatedWeightRecordedDate)
		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okPayload.RequestedPickupDate)
		suite.equalDatePtr(minimalShipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		fmt.Println("nefore", minimalShipment.UpdatedAt)
		minimalShipment = suite.refreshFromDB(minimalShipment.ID)
		fmt.Println("afer", minimalShipment.UpdatedAt)

	})
	suite.T().Run("PATCH failure cannot update primeEstimatedWeight again", func(t *testing.T) {
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     Use previously created shipment with estimated weight
		//             Attempt to update primeEstimatedWeight
		// Expected outcome:
		//             Handler returns Unprocessable Entity
		//             primeEstimatedWeight cannot be updated more than once.

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   minimalReq,
			MtoShipmentID: *handlers.FmtUUID(minimalShipment.ID),
			Body: &primemessages.MTOShipment{
				PrimeEstimatedWeight: int64(primeEstimatedWeight), // New estimated weight
				PrimeActualWeight:    int64(primeActualWeight),    // New actual weight
			},
			IfMatch: etag.GenerateEtag(minimalShipment.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errPayload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload
		suite.Contains(errPayload.InvalidFields, "primeEstimatedWeight")

	})

	suite.T().Run("PATCH failure 422 cannot update estimatedWeight without scheduledPickupDate", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     Create a shipment with no scheduledPickupDate
		//             Attempt to update the primeEstimatedWeight
		// Expected outcome:
		//             Handler returns Unprocessable entity because the
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
		mtoShipment := primemessages.MTOShipment{
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
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with nil uuid
		// Expected outcome:
		//             Handler returns Not Found error

		// Create request with non existent ID
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(uuid.Nil.String()), // unknown shipment id
			Body:          &primemessages.MTOShipment{},
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}
		// Call handler
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PATCH failure 422 set readOnly field", func(t *testing.T) {
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with counselorRemarks (a readOnly
		//             field) set
		// Expected outcome:
		//             Handler returns Unprocessable Entity error
		remarks := fmt.Sprintf("test conflict %s", time.Now())
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(shipment.ID.String()),
			Body:          &primemessages.MTOShipment{CustomerRemarks: &remarks},
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}
		// Call handler
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("PATCH failure 412 precondition failed", func(t *testing.T) {
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with old eTag
		// Expected outcome:
		//             Handler returns Precondition Failed error

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(shipment.ID.String()),
			Body:          &primemessages.MTOShipment{Diversion: true},   // update anything
			IfMatch:       string(etag.GenerateEtag(shipment.CreatedAt)), // Use createdAt to generate eTag
		}

		// Call handler
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.T().Run("PATCH success returns all nested objects", func(t *testing.T) {
		// Under test: Handle function
		// Mocked:     Planner
		// Set up:     We add service items to the shipment in the DB
		//             We provide an almost empty update so as to check that the
		//             nested objects in the response are fully populated
		// Expected outcome:
		//             Handler returns OK, all service items, agents and addresses are
		//             populated.
		// MYTODO agents

		// Add service items to our shipment
		// Create a service item in the db, associate with the shipment
		reService := testdatagen.MakeDDFSITReService(suite.DB())
		mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move:        move,
			MTOShipment: shipment,
			ReService:   reService,
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: move.ID,
				ReServiceID:     reService.ID,
				MTOShipmentID:   &shipment.ID,
			},
		})
		// Associate locally as well
		shipment.MTOServiceItems = models.MTOServiceItems{mtoServiceItem1}

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
		update := primemessages.MTOShipment{
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
		suite.equalDatePtr(shipment.ApprovedDate, okPayload.ApprovedDate)
		suite.equalDatePtr(shipment.FirstAvailableDeliveryDate, okPayload.FirstAvailableDeliveryDate)
		suite.equalDatePtr(shipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		suite.equalDatePtr(shipment.RequiredDeliveryDate, okPayload.RequiredDeliveryDate)
		suite.equalDatePtr(shipment.ScheduledPickupDate, okPayload.ScheduledPickupDate)

		suite.EqualAddress(*shipment.PickupAddress, &okPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &okPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &okPayload.SecondaryDeliveryAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &okPayload.SecondaryPickupAddress.Address, true)

	})
}

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
	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	suite.T().Run("PATCH success 200 create addresses", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function, addresses mechanism - we can create but not update
		// Mocked:     Planner
		// Set up:     We use a shipment with minimal info, no addresses
		//             Update with PickupAddress, DestinationAddress, SecondaryPickupAddress, SecondaryDeliveryAddress
		//             Expected: Handler should return OK, new addresses created

		// CREATE REQUEST
		// Create an update message with all addresses provided
		update := primemessages.MTOShipment{
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

	})

	suite.T().Run("PATCH failure 422 update addresses not allowed", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function, addresses mechanism - we cannot update addresses
		// Mocked:     Planner
		// Set up:     We use a shipment with addresses
		//             Update with PickupAddress, DestinationAddress, SecondaryPickupAddress, SecondaryDeliveryAddress
		// Expected:   Handler should return unprocessable entity error. All addresses cannot be updated and should
		//             be listed in errors

		// CREATE REQUEST
		// Create an update message with all new addresses provided
		update := primemessages.MTOShipment{
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
		// TESTCASE SCENARIO
		// Under test: Handle function, addresses mechanism - we can create but not update
		// Mocked:     Planner
		// Set up:     We use a shipment with addresses
		//             Update with nil for the addresses
		// Expected:   Handler should return OK, addresses should be unchanged.
		//             This endpoint was previously blanking out addresses which is why we have this test.

		// CREATE REQUEST
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.MTOShipment{}, // Empty payload
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
func (suite *HandlerSuite) TestUpdateMTOShipmentDateLogic() {

	primeEstimatedWeight := unit.Pound(500)

	// Create an available shipment in DB
	move := testdatagen.MakeAvailableMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)

	now := time.Now()
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	// used for all tests except the 500 server error:
	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	req2 := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", mtoShipment2.ID.String()), nil)

	// Prime-specific validations tested below
	suite.T().Run("Failed case if not both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		eightDaysFromNow := now.AddDate(0, 0, 8)
		threeDaysBefore := now.AddDate(0, 0, -3)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &eightDaysFromNow,
				ApprovedDate:        &threeDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req2,
			MtoShipmentID: *handlers.FmtUUID(mtoShipment2.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
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
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
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
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:       "APPROVED",
				ApprovedDate: &now,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		schedDate := strfmt.Date(tenDaysFromNow)
		payload := primemessages.MTOShipment{
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			ScheduledPickupDate:  &schedDate,
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
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
		suite.NotNil(okResponse.Payload.RequestedPickupDate)
		suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))

	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate updated on destinationAddress creation", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function, RequiredDeliveryDate logic
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
		update := primemessages.MTOShipment{
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
		suite.equalDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate for Alaska", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function, RequiredDeliveryDate logic
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
		payload := primemessages.MTOShipment{
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
		suite.equalDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.T().Run("PATCH Success 200 RequiredDeliveryDate for Adak, Alaska", func(t *testing.T) {
		// TESTCASE SCENARIO
		// Under test: Handle function, RequiredDeliveryDate logic
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
		payload := primemessages.MTOShipment{
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
		suite.equalDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

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
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
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
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

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
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
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
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

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
	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)

	handler := UpdateMTOShipmentStatusHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner),
		mtoshipment.NewMTOShipmentStatusUpdater(suite.DB(), builder,
			mtoserviceitem.NewMTOServiceItemCreator(builder), planner),
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

// Compares the time.Time from the model with the strfmt.date from the payload
// If one is nil, both should be nil, else they should match in value
// This is to be strictly used for dates as it drops any time parameters in the comparison
func (suite *HandlerSuite) equalDatePtr(expected *time.Time, actual *strfmt.Date) {
	if expected == nil || actual == nil {
		suite.Nil(expected)
		suite.Nil(actual)
	} else {
		isoDate := "2006-01-02" // Create a date format
		suite.Equal(expected.Format(isoDate), time.Time(*actual).Format(isoDate))
	}
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
