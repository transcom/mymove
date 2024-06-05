package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	moveRouter := moveservices.NewMoveRouter()
	fetcher := fetch.NewFetcher(builder)
	addressCreator := address.NewAddressCreator()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
	ppmEstimator := mocks.PPMEstimator{}
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
	shipmentRouter := mtoshipment.NewShipmentRouter()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		builder,
		mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		moveRouter,
	)
	shipmentCreator := shipmentorchestrator.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, shipmentRouter, moveTaskOrderUpdater)
	mockCreator := mocks.ShipmentCreator{}

	var pickupAddress primemessages.Address
	var destinationAddress primemessages.Address

	setupTestData := func() (CreateMTOShipmentHandler, models.Move) {

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		handler := CreateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentCreator,
			mtoChecker,
		}

		// Make stubbed addresses just to collect address data for payload
		newAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		pickupAddress = primemessages.Address{
			City:           &newAddress.City,
			Country:        newAddress.Country,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		newAddress = factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
		destinationAddress = primemessages.Address{
			City:           &newAddress.City,
			Country:        newAddress.Country,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		return handler, move

	}

	suite.Run("Successful POST - Integration Test", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Successful submission, status should be SUBMITTED
		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createMTOShipmentPayload.Validate(strfmt.Default))

		// check that the mto shipment status is Submitted
		suite.Require().Equal(createMTOShipmentPayload.Status, primemessages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.PrimeEstimatedWeight, params.Body.PrimeEstimatedWeight)
	})

	suite.Run("Successful POST - Integration Test (PPM)", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:      Create a PPM shipment on an available move
		// Expected:   Successful submission, status should be SUBMITTED
		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		counselorRemarks := "Some counselor remarks"
		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		pickupPostalCode := "30907"
		secondaryPickupPostalCode := "30809"
		destinationPostalCode := "29212"
		secondaryDestinationPostalCode := "29201"
		sitExpected := true
		sitLocation := primemessages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1500)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3200)
		hasProGear := true
		proGearWeight := unit.Pound(400)
		spouseProGearWeight := unit.Pound(250)
		estimatedIncentive := 123456
		sitEstimatedCost := 67500

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:  handlers.FmtUUID(move.ID),
				ShipmentType:     primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypePPM),
				CounselorRemarks: &counselorRemarks,
				PpmShipment: &primemessages.CreatePPMShipment{
					ExpectedDepartureDate:          handlers.FmtDate(expectedDepartureDate),
					PickupPostalCode:               &pickupPostalCode,
					SecondaryPickupPostalCode:      &secondaryPickupPostalCode,
					DestinationPostalCode:          &destinationPostalCode,
					SecondaryDestinationPostalCode: &secondaryDestinationPostalCode,
					SitExpected:                    &sitExpected,
					SitLocation:                    &sitLocation,
					SitEstimatedWeight:             handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:          handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate:      handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:                handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                     &hasProGear,
					ProGearWeight:                  handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:            handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createdShipment := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createdShipment.Validate(strfmt.Default))

		createdPPM := createdShipment.PpmShipment

		suite.Equal(move.ID.String(), createdShipment.MoveTaskOrderID.String())
		suite.Equal(primemessages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(primemessages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, createdShipment.Status)
		suite.Equal(&counselorRemarks, createdShipment.CounselorRemarks)

		suite.Equal(createdShipment.ID.String(), createdPPM.ShipmentID.String())
		suite.Equal(primemessages.PPMShipmentStatusSUBMITTED, createdPPM.Status)
		suite.Equal(handlers.FmtDatePtr(&expectedDepartureDate), createdPPM.ExpectedDepartureDate)
		suite.Equal(&pickupPostalCode, createdPPM.PickupPostalCode)
		suite.Equal(&secondaryPickupPostalCode, createdPPM.SecondaryPickupPostalCode)
		suite.Equal(&destinationPostalCode, createdPPM.DestinationPostalCode)
		suite.Equal(&secondaryDestinationPostalCode, createdPPM.SecondaryDestinationPostalCode)
		suite.Equal(&sitExpected, createdPPM.SitExpected)
		suite.Equal(&sitLocation, createdPPM.SitLocation)
		suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), createdPPM.SitEstimatedWeight)
		suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), createdPPM.SitEstimatedEntryDate)
		suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), createdPPM.SitEstimatedDepartureDate)
		suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), createdPPM.EstimatedWeight)
		suite.Equal(handlers.FmtBool(hasProGear), createdPPM.HasProGear)
		suite.Equal(handlers.FmtPoundPtr(&proGearWeight), createdPPM.ProGearWeight)
		suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), createdPPM.SpouseProGearWeight)
		suite.Equal(int64(estimatedIncentive), *createdPPM.EstimatedIncentive)
		suite.Equal(int64(sitEstimatedCost), *createdPPM.SitEstimatedCost)
	})

	suite.Run("Successful POST with Shuttle service items without primeEstimatedWeight - Integration Test", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:   Create an mto shipment on an available move
		// Expected:   Successful submission, status should be SUBMITTED
		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Reason: models.StringPointer("not applicable"),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
		}, nil)
		serviceItem.ID = uuid.Nil

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:     handlers.FmtUUID(move.ID),
				Agents:              nil,
				CustomerRemarks:     nil,
				PointOfContact:      "John Doe",
				RequestedPickupDate: handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:        primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:       struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:  struct{ primemessages.Address }{destinationAddress},
			},
		}

		mtoServiceItems := models.MTOServiceItems{serviceItem}
		params.Body.SetMtoServiceItems(*payloads.MTOServiceItems(&mtoServiceItems))

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createMTOShipmentPayload.Validate(strfmt.Default))

		// check that the mto shipment status is Submitted
		suite.Require().Equal(createMTOShipmentPayload.Status, primemessages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, "MTO Shipment should have been submitted")
	})

	suite.Run("POST failure - 500", func() {
		// Under Test: CreateMTOShipmentHandler
		// Mocked:     CreateMTOShipment creator
		// Setup:   If underlying CreateMTOShipment returns error, handler should return 500 response
		// Expected:   500 Response returned
		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		// Create a handler with the mocked creator
		handler.ShipmentCreator = &mockCreator

		err := errors.New("ServerError")

		mockCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil, err)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.Run("POST failure - 422 -- Bad agent IDs set on shipment", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with an agent that doesn't really exist, handler should return unprocessable entity
		// Expected:   422 Unprocessable Entity Response returned

		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		badID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		agent := &primemessages.MTOAgent{
			ID:            badID,
			MtoShipmentID: badID,
			FirstName:     handlers.FmtString("Mary"),
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				Agents:               primemessages.MTOAgents{agent},
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.NotEmpty(typedResponse.Payload.InvalidFields)
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:id")
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:mtoShipmentID")
	})

	suite.Run("POST failure - 422 - invalid input, missing pickup address", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with missing pickup address, handler should return unprocessable entity
		// Expected:   422 Unprocessable Entity Response returned

		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
			},
		}
		params.Body.PickupAddress.Address.StreetAddress1 = nil

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		unprocessableEntity := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		// TODO: Can't validate the response because of the issue noted below. Figure out a way to
		//   either alter the service or relax the swagger requirements.
		// suite.NoError(unprocessableEntity.Payload.Validate(strfmt.Default))
		// CreateShipment is returning apperror.InvalidInputError without any validation errors
		// so InvalidFields won't be added to the payload.

		suite.Contains(*unprocessableEntity.Payload.Detail, "PickupAddress is required")
	})

	suite.Run("POST failure - 404 -- not found", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment on a non-existent move
		// Expected:   404 Not Found returned
		handler, _ := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		// Generate a unique id
		badID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      &badID,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.CreateMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a request with no data in the body
		// Expected:   422 Unprocessable Entity Response returned

		handler, _ := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		paramsNilBody := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}

		// Validate incoming payload: nil body (the point of this test)

		response := handler.Handle(paramsNilBody)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
		responsePayload := response.(*mtoshipmentops.CreateMTOShipmentBadRequest).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 404 -- MTO is not available to Prime", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment on an unavailable move, prime cannot update these
		// Expected:   404 Not found returned

		handler, _ := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		unavailableMove := factory.BuildMove(suite.DB(), nil, nil)
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(unavailableMove.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primemessages.Address }{pickupAddress},
				DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentNotFound)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*typedResponse.Payload.Detail, unavailableMove.ID.String())
	})

	suite.Run("POST failure - 422 - modelType() not supported", func() {
		// Under Test: CreateMTOShipmentHandler
		// Setup:      Create a shipment with service items that don't match the modeltype
		// Expected:   422 Unprocessable Entity returned

		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)
		handler.ShipmentCreator = &mockCreator

		err := apperror.NotFoundError{}

		mockCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		// Create a service item that doesn't match the modeltype
		mtoServiceItems := models.MTOServiceItems{
			models.MTOServiceItem{
				MoveTaskOrderID:  move.ID,
				MTOShipmentID:    &uuid.Nil,
				ReService:        models.ReService{Code: models.ReServiceCodeMS},
				Reason:           nil,
				PickupPostalCode: nil,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primemessages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primemessages.NewMTOShipmentType(primemessages.MTOShipmentTypeHHG),
			},
		}

		params.Body.SetMtoServiceItems(*payloads.MTOServiceItems(&mtoServiceItems))

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*typedResponse.Payload.Detail, "MTOServiceItem modelType() not allowed")
	})
}

func (suite *HandlerSuite) TestUpdateShipmentDestinationAddressHandler() {
	req := httptest.NewRequest("POST", "/mto-shipments/{mtoShipmentID}/shipment-address-updates", nil)

	makeSubtestData := func() mtoshipmentops.UpdateShipmentDestinationAddressParams {
		contractorRemark := "This is a contractor remark"
		body := primemessages.UpdateShipmentDestinationAddress{
			ContractorRemarks: &contractorRemark,
			NewAddress: &primemessages.Address{
				City:           swag.String("Beverly Hills"),
				PostalCode:     swag.String("90210"),
				State:          swag.String("CA"),
				StreetAddress1: swag.String("1234 N. 1st Street"),
			},
		}

		params := mtoshipmentops.UpdateShipmentDestinationAddressParams{
			HTTPRequest: req,
			Body:        &body,
		}

		return params

	}
	suite.Run("POST failure - 422 Unprocessable Entity Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := UpdateShipmentDestinationAddressHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}
		// InvalidInputError should generate an UnprocessableEntity response error
		// Need verrs incorporated to satisfy swagger validation
		verrs := validate.NewErrors()
		verrs.Add("some key", "some value")
		err := apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "unable to create ShipmentAddressUpdate")

		mockCreator.On("RequestShipmentDeliveryAddressUpdate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.Address"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData)
		suite.IsType(&mtoshipmentops.UpdateShipmentDestinationAddressUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.UpdateShipmentDestinationAddressUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 409 Request conflict reponse Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := UpdateShipmentDestinationAddressHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}
		// NewConflictError should generate a RequestConflict response error
		err := apperror.NewConflictError(uuid.Nil, "unable to create ShipmentAddressUpdate")

		mockCreator.On("RequestShipmentDeliveryAddressUpdate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.Address"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData)
		suite.IsType(&mtoshipmentops.UpdateShipmentDestinationAddressConflict{}, response)
		errResponse := response.(*mtoshipmentops.UpdateShipmentDestinationAddressConflict)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 404 Not Found response error", func() {

		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := UpdateShipmentDestinationAddressHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}
		// NewNotFoundError should generate a RequestNotFound response error
		err := apperror.NewNotFoundError(uuid.Nil, "unable to create ShipmentAddressUpdate")

		mockCreator.On("RequestShipmentDeliveryAddressUpdate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.Address"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData)
		suite.IsType(&mtoshipmentops.UpdateShipmentDestinationAddressNotFound{}, response)
		errResponse := response.(*mtoshipmentops.UpdateShipmentDestinationAddressNotFound)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("500 server error", func() {

		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := UpdateShipmentDestinationAddressHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}
		// NewQueryError should generate an InternalServerError response error
		err := apperror.NewQueryError("", nil, "unable to reach database")

		mockCreator.On("RequestShipmentDeliveryAddressUpdate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.Address"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData)
		suite.IsType(&mtoshipmentops.UpdateShipmentDestinationAddressInternalServerError{}, response)
		errResponse := response.(*mtoshipmentops.UpdateShipmentDestinationAddressInternalServerError)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
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

	// Create some usable weights
	primeEstimatedWeight := unit.Pound(500)
	primeActualWeight := unit.Pound(600)

	// Create service objects needed for handler
	// ghcDomesticTime is used in the planner, the planner checks transit distance.
	// We mock the planner to return 400, so we need an entry that will return a
	// transit time of 12 days for a distance of 400.

	// Mock planner to always return a distance of 400 mi
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()

	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
	ppmEstimator := mocks.PPMEstimator{}
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
	setupTestData := func() (UpdateMTOShipmentHandler, models.MTOShipment) {
		// Add a 12 day transit time for a distance of 400
		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 500,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetHHGPlanner(planner)
		handler := UpdateMTOShipmentHandler{
			handlerConfig,
			shipmentUpdater,
		}

		// Create an available shipment in DB
		now := testdatagen.CurrentDateWithoutTime()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusApproved,
					ApprovedDate: now,
				},
			},
			{
				Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3}),
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4}),
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			},
		}, nil)
		return handler, shipment
	}

	suite.Run("PATCH failure 500 Unit Test", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     MTOShipmentUpdater, Planner
		// Set up:     We provide an update but make MTOShipmentUpdater return a server error
		// Expected:   Handler returns Internal Server Error Response. This ensures if there is an
		//             unexpected error in the service object, we return the proper HTTP response
		handler, shipment := setupTestData()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				Diversion: true,
			},
			IfMatch: eTag,
		}

		mockUpdater := mocks.ShipmentUpdater{}
		mockHandler := UpdateMTOShipmentHandler{
			handler.HandlerConfig,
			&mockUpdater,
		}
		internalServerErr := errors.New("ServerError")

		mockUpdater.On("MTOShipmentsMTOAvailableToPrime",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(true, nil)

		mockUpdater.On("UpdateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.Run("PATCH success 200 minimal update", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We use the normal (non-minimal) shipment we created earlier
		//             We provide an update with minimal changes
		// Expected:   Handler returns OK
		//             Minimal updates are completed, old values retained for rest of
		//             shipment. This tests that PATCH is not accidentally clearing any existing
		//             data.
		handler, shipment := setupTestData()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)

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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(okPayload.Validate(strfmt.Default))

		suite.Equal(shipment.ID.String(), okPayload.ID.String())

		// Confirm PATCH working as expected; non-updated values still exist
		suite.EqualDatePtr(shipment.ApprovedDate, okPayload.ApprovedDate)
		suite.EqualDatePtr(shipment.FirstAvailableDeliveryDate, okPayload.FirstAvailableDeliveryDate)
		suite.EqualDatePtr(shipment.RequestedPickupDate, okPayload.RequestedPickupDate)
		suite.EqualDatePtr(shipment.RequiredDeliveryDate, okPayload.RequiredDeliveryDate)
		suite.EqualDatePtr(shipment.ScheduledPickupDate, okPayload.ScheduledPickupDate)
		suite.EqualDatePtr(shipment.ActualDeliveryDate, okPayload.ActualDeliveryDate)
		suite.EqualDatePtr(shipment.ScheduledDeliveryDate, okPayload.ScheduledDeliveryDate)

		suite.EqualAddress(*shipment.PickupAddress, &okPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &okPayload.DestinationAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryDeliveryAddress, &okPayload.SecondaryDeliveryAddress.Address, true)
		suite.EqualAddress(*shipment.SecondaryPickupAddress, &okPayload.SecondaryPickupAddress.Address, true)

		// Confirm new values
		suite.Equal(params.Body.Diversion, okPayload.Diversion)
		suite.Equal(params.Body.ActualPickupDate.String(), okPayload.ActualPickupDate.String())
	})

	suite.Run("PATCH success 200 update destination type", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We use the normal (non-minimal) shipment we created earlier
		//             We provide an update with minimal changes
		// Expected:   Handler returns OK
		//             Minimal updates are completed, old values retained for rest of
		//             shipment. This tests that PATCH is not accidentally clearing any existing
		//             data.
		handler, shipment := setupTestData()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)

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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(okPayload.Validate(strfmt.Default))

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

	suite.Run("Successful PATCH - Integration Test (PPM)", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner, PPMEstimator
		// Set up:     We create a ppm shipment
		//             We provide an update
		// Expected:   Handler returns OK
		//             Updates are completed
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
		}

		hasProGear := true
		now := time.Now()
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
				},
			},
			{
				Model: models.PPMShipment{
					HasProGear: &hasProGear,
				},
			},
		}, nil)
		year, month, day := time.Now().Date()
		actualMoveDate := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
		expectedDepartureDate := actualMoveDate.Add(time.Hour * 24 * 2)
		pickupPostalCode := "30907"
		secondaryPickupPostalCode := "30809"
		destinationPostalCode := "36106"
		secondaryDestinationPostalCode := "36101"
		sitExpected := true
		sitLocation := primemessages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1700)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3000)
		proGearWeight := unit.Pound(300)
		spouseProGearWeight := unit.Pound(200)
		estimatedIncentive := 654321
		sitEstimatedCost := 67500

		req := httptest.NewRequest("PATCH", "/mto-shipments/{MtoShipmentID}", nil)
		eTag := etag.GenerateEtag(ppmShipment.Shipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(ppmShipment.ShipmentID),
			IfMatch:       eTag,
			Body: &primemessages.UpdateMTOShipment{
				CounselorRemarks: handlers.FmtString("Test remark"),
				PpmShipment: &primemessages.UpdatePPMShipment{
					ExpectedDepartureDate:          handlers.FmtDatePtr(&expectedDepartureDate),
					PickupPostalCode:               &pickupPostalCode,
					SecondaryPickupPostalCode:      &secondaryPickupPostalCode,
					DestinationPostalCode:          &destinationPostalCode,
					SecondaryDestinationPostalCode: &secondaryDestinationPostalCode,
					SitExpected:                    &sitExpected,
					SitEstimatedWeight:             handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:          handlers.FmtDatePtr(&sitEstimatedEntryDate),
					SitEstimatedDepartureDate:      handlers.FmtDatePtr(&sitEstimatedDepartureDate),
					SitLocation:                    &sitLocation,
					EstimatedWeight:                handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                     &hasProGear,
					ProGearWeight:                  handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:            handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		suite.Equal(ppmShipment.Shipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(handlers.FmtDatePtr(&expectedDepartureDate), updatedShipment.PpmShipment.ExpectedDepartureDate)
		suite.Equal(&pickupPostalCode, updatedShipment.PpmShipment.PickupPostalCode)
		suite.Equal(&secondaryPickupPostalCode, updatedShipment.PpmShipment.SecondaryPickupPostalCode)
		suite.Equal(&destinationPostalCode, updatedShipment.PpmShipment.DestinationPostalCode)
		suite.Equal(&secondaryDestinationPostalCode, updatedShipment.PpmShipment.SecondaryDestinationPostalCode)
		suite.Equal(sitExpected, *updatedShipment.PpmShipment.SitExpected)
		suite.Equal(&sitLocation, updatedShipment.PpmShipment.SitLocation)
		suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), updatedShipment.PpmShipment.SitEstimatedWeight)
		suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), updatedShipment.PpmShipment.SitEstimatedEntryDate)
		suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), updatedShipment.PpmShipment.SitEstimatedDepartureDate)
		suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), updatedShipment.PpmShipment.EstimatedWeight)
		suite.Equal(int64(estimatedIncentive), *updatedShipment.PpmShipment.EstimatedIncentive)
		suite.Equal(int64(sitEstimatedCost), *updatedShipment.PpmShipment.SitEstimatedCost)
		suite.Equal(handlers.FmtBool(hasProGear), updatedShipment.PpmShipment.HasProGear)
		suite.Equal(handlers.FmtPoundPtr(&proGearWeight), updatedShipment.PpmShipment.ProGearWeight)
		suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), updatedShipment.PpmShipment.SpouseProGearWeight)
		suite.Equal(params.Body.CounselorRemarks, updatedShipment.CounselorRemarks)
	})

	suite.Run("PATCH failure 404 not found because not available to prime", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We provide an update to a shipment whose associated move isn't available to prime
		// Expected:   Handler returns Not Found error
		handler, _ := setupTestData()

		// Create a shipment unavailable to Prime in DB
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		suite.Nil(shipment.MoveTaskOrder.AvailableToPrimeAt)

		// Create params
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				Diversion: true,
			},
			IfMatch: etag.GenerateEtag(shipment.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// Verify not found response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure 404 not found because attempting to update an external vendor shipment", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We provide an update to a shipment that is handled by an external vendor
		// Expected:   Handler returns Not Found error
		handler, ogShipment := setupTestData()

		// Create a shipment handled by an external vendor
		externalShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    ogShipment.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		// Check that they point to the same move and that it's available
		suite.Equal(ogShipment.MoveTaskOrderID, externalShipment.MoveTaskOrderID)
		suite.NotNil(ogShipment.MoveTaskOrder.AvailableToPrimeAt)

		// Create params
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", externalShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(externalShipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				Diversion: true,
			},
			IfMatch: etag.GenerateEtag(externalShipment.UpdatedAt),
		}

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// Verify not found response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("PATCH success 200 update of primeEstimatedWeight and primeActualWeight", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We provide an update with actual and estimated weights
		// Expected:   Handler returns OK
		//             Weights are updated, and prime estimated weight recorded date is updated.
		handler, ogShipment := setupTestData()
		// Create a minimal shipment on the previously created move
		minimalShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    ogShipment.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		minimalReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", minimalShipment.ID.String()), nil)

		// Check that PrimeEstimatedWeightRecordedDate was nil at first
		suite.Nil(minimalShipment.PrimeEstimatedWeightRecordedDate)

		// Update the primeEstimatedWeight
		eTag := etag.GenerateEtag(minimalShipment.UpdatedAt)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   minimalReq,
			MtoShipmentID: *handlers.FmtUUID(minimalShipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight), // New estimated weight
				PrimeActualWeight:    handlers.FmtPoundPtr(&primeActualWeight),    // New actual weight
			},
			IfMatch: eTag,
		}

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(okPayload.Validate(strfmt.Default))

		suite.Equal(minimalShipment.ID.String(), okPayload.ID.String())

		// Confirm changes to weights
		suite.Equal(int64(primeActualWeight), *okPayload.PrimeActualWeight)
		suite.Equal(int64(primeEstimatedWeight), *okPayload.PrimeEstimatedWeight)
		// Confirm primeEstimatedWeightRecordedDate was added
		suite.NotNil(okPayload.PrimeEstimatedWeightRecordedDate)
		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okPayload.RequestedPickupDate)
		suite.EqualDatePtr(minimalShipment.RequestedPickupDate, okPayload.RequestedPickupDate)

		// refresh shipment from DB for getting the updated eTag
		minimalShipment = suite.refreshFromDB(minimalShipment.ID)

	})

	suite.Run("PATCH failure 422 cannot update primeEstimatedWeight again", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Use previously created shipment with primeEstimatedWeight updated
		//             Attempt to update primeEstimatedWeight
		// Expected:   Handler returns Unprocessable Entity
		//             primeEstimatedWeight cannot be updated more than once.
		handler, ogShipment := setupTestData()
		// Create a minimal shipment on the previously created move
		minimalShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    ogShipment.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		minimalReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", minimalShipment.ID.String()), nil)

		// Set the primeEstimatedWeight
		// Update the primeEstimatedWeight
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   minimalReq,
			MtoShipmentID: *handlers.FmtUUID(minimalShipment.ID),
			Body: &primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight), // New estimated weight
				PrimeActualWeight:    handlers.FmtPoundPtr(&primeActualWeight),    // New actual weight
			},
			IfMatch: etag.GenerateEtag(minimalShipment.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))

		minimalShipment = suite.refreshFromDB(minimalShipment.ID)
		// Check that primeEstimatedWeight was already populated
		suite.NotNil(minimalShipment.PrimeEstimatedWeight)

		// Attempt to update again
		updatedEstimatedWeight := primeEstimatedWeight + 100
		params.Body.PrimeEstimatedWeight = handlers.FmtPoundPtr(&updatedEstimatedWeight)
		params.IfMatch = etag.GenerateEtag(minimalShipment.UpdatedAt)

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response = handler.Handle(params)

		// Check response contains an error about primeEstimatedWeight
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errPayload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(errPayload.Validate(strfmt.Default))

		suite.Contains(errPayload.InvalidFields, "primeEstimatedWeight")
	})

	suite.Run("PATCH failure 404 unknown shipment", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with fake uuid
		// Expected:   Handler returns Not Found error
		handler, shipment := setupTestData()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)

		// Create request with non existent ID
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(uuid.Must(uuid.NewV4()).String()), // generate a UUID
			Body:          &primemessages.UpdateMTOShipment{},
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}
		// Call handler

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure 412 precondition failed", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     Attempt to update a shipment with old eTag
		// Expected:   Handler returns Precondition Failed error
		handler, shipment := setupTestData()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)

		// Create an update with a wrong eTag
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: strfmt.UUID(shipment.ID.String()),
			Body:          &primemessages.UpdateMTOShipment{Diversion: true}, // update anything
			IfMatch:       string(etag.GenerateEtag(time.Now())),             // use the wrong time to generate etag
		}

		// Call handler

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("PATCH success 200 returns all nested objects", func() {
		// Under test: updateMTOShipmentHandler.Handle
		// Mocked:     Planner
		// Set up:     We add service items to the shipment in the DB
		//             We provide an almost empty update so as to check that the
		//             nested objects in the response are fully populated
		// Expected:   Handler returns OK, all service items, agents and addresses are
		//             populated.
		handler, shipment := setupTestData()

		// Add service items to our shipment
		// Create a service item in the db, associate with the shipment
		reService := factory.BuildDDFSITReService(suite.DB())
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: shipment.MoveTaskOrderID,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					MoveTaskOrderID: shipment.MoveTaskOrderID,
					ReServiceID:     reService.ID,
					MTOShipmentID:   &shipment.ID,
					SITEntryDate:    models.TimePointer(time.Now()),
					Reason:          models.StringPointer("lorem epsum"),
				},
			},
		}, nil)

		// Add agents associated to our shipment
		agent1 := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName:    models.StringPointer("Test1"),
					LastName:     models.StringPointer("Agent"),
					Email:        models.StringPointer("test@test.email.com"),
					MTOAgentType: models.MTOAgentReceiving,
				},
			},
		}, nil)
		agent2 := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName:    models.StringPointer("Test2"),
					LastName:     models.StringPointer("Agent"),
					Email:        models.StringPointer("test@test.email.com"),
					MTOAgentType: models.MTOAgentReleasing,
				},
			},
		}, nil)

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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// Check response
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(okPayload.Validate(strfmt.Default))

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

	// CREATE HANDLER OBJECT
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	addressUpdater := address.NewAddressUpdater()

	planner.On("ZipTransitDistance",
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
	addressCreator := address.NewAddressCreator()
	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
	ppmEstimator := mocks.PPMEstimator{}

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

	setupTestData := func() (UpdateMTOShipmentHandler, models.MTOShipment) {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetHHGPlanner(planner)
		handler := UpdateMTOShipmentHandler{
			handlerConfig,
			shipmentUpdater,
		}
		// Create a shipment in the DB that has no addresses populated:
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		return handler, shipment

	}

	suite.Run("PATCH success 200 create addresses", func() {
		// Under test: updateMTOShipmentHandler.Handle, addresses mechanism - we can create but not update
		// Mocked:     Planner
		// Set up:     We use a shipment with minimal info, no addresses
		//             Update with PickupAddress, DestinationAddress, SecondaryPickupAddress, SecondaryDeliveryAddress
		// Expected:   Handler should return OK, new addresses created
		handler, shipment := setupTestData()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", shipment.ID.String()), nil)

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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(okPayload.Validate(strfmt.Default))

		// Check that addresses match what was sent
		suite.EqualAddressPayload(&update.PickupAddress.Address, &okPayload.PickupAddress.Address, false)
		suite.EqualAddressPayload(&update.DestinationAddress.Address, &okPayload.DestinationAddress.Address, false)
		suite.EqualAddressPayload(&update.SecondaryPickupAddress.Address, &okPayload.SecondaryPickupAddress.Address, false)
		suite.EqualAddressPayload(&update.SecondaryDeliveryAddress.Address, &okPayload.SecondaryDeliveryAddress.Address, false)

	})

	suite.Run("PATCH failure 422 update addresses not allowed", func() {
		// Under test: updateMTOShipmentHandler.Handle, addresses mechanism - we cannot update addresses
		// Mocked:     Planner
		// Set up:     We create a shipment with Pickup and Destination address.
		//             Then we update with PickupAddress, DestinationAddress, SecondaryPickupAddress, SecondaryDeliveryAddress
		// Expected:   Handler should return unprocessable entity error for those addresses already created, but not the new ones.
		//             Addresses cannot be updated with this endpoint, only created,  and should be listed in errors
		handler, _ := setupTestData()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", shipment.ID.String()), nil)

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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errPayload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(errPayload.Validate(strfmt.Default))

		suite.Contains(errPayload.InvalidFields, "pickupAddress")
		suite.Contains(errPayload.InvalidFields, "destinationAddress")
		suite.NotContains(errPayload.InvalidFields, "secondaryPickupAddress")
		suite.NotContains(errPayload.InvalidFields, "secondaryDeliveryAddress")

	})

	suite.Run("PATCH success 200 nil doesn't clear addresses", func() {
		// Under test: updateMTOShipmentHandler.Handle, addresses mechanism - we can create but not update
		// Mocked:     Planner
		// Set up:     Create a shipment with addresses populated.
		//             Update with nil for the addresses.
		// Expected:   Handler should return OK, addresses should be unchanged.
		//             This endpoint was previously blanking out addresses which is why we have this test.
		handler, _ := setupTestData()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", shipment.ID.String()), nil)

		// CREATE REQUEST
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.UpdateMTOShipment{}, // Empty payload
			IfMatch:       string(etag.GenerateEtag(shipment.UpdatedAt)),
		}

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		newPayload := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(newPayload.Validate(strfmt.Default))

		// Check that addresses match what was returned previously in the last successful payload
		suite.EqualAddress(*shipment.PickupAddress, &newPayload.PickupAddress.Address, true)
		suite.EqualAddress(*shipment.DestinationAddress, &newPayload.DestinationAddress.Address, true)
	})
}

// TestUpdateMTOShipmentDateLogic tests various restrictions related to timelines that
// Prime is required to abide by
// More details about these rules can be found in the Performance Work Statement for the
// Global Household Goods Contract HTC711-19-R-R004
func (suite *HandlerSuite) TestUpdateMTOShipmentDateLogic() {

	// ghcDomesticTime is used in the planner, the planner checks transit distance.
	// We mock the planner to return 400, so we need an entry that will return a
	// transit time of 12 days for a distance of 400.

	// Mock planner to always return a distance of 400 mi
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	addressUpdater := address.NewAddressUpdater()

	// Add a 12 day transit time for a distance of 400
	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 1,
		DistanceMilesUpper: 500,
	}

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
	addressCreator := address.NewAddressCreator()

	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
	ppmEstimator := mocks.PPMEstimator{}

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

	setupTestData := func() (UpdateMTOShipmentHandler, models.Move) {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetHHGPlanner(planner)
		handler := UpdateMTOShipmentHandler{
			handlerConfig,
			shipmentUpdater,
		}
		// Create an available move to be used for the shipments
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		// Add the transit time record
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		return handler, move
	}

	primeEstimatedWeight := unit.Pound(500)
	now := time.Now()

	// Date checks for revised functionality without approval dates
	suite.Run("Successful case if estimated weight added before scheduled pickup date", func() {
		handler, move := setupTestData()

		twoDaysFromNow := now.AddDate(0, 0, 2)
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ScheduledPickupDate: &twoDaysFromNow,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
		}
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(oldShipment.ID.String(), okResponse.Payload.ID.String())

		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okResponse.Payload.RequestedPickupDate)
		suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))
	})

	suite.Run("Successful case if estimated weight added on scheduled pickup date", func() {
		handler, move := setupTestData()

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ScheduledPickupDate: &now,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
		}
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(oldShipment.ID.String(), okResponse.Payload.ID.String())

		// Confirm PATCH working as expected; non-updated value still exists
		suite.NotNil(okResponse.Payload.RequestedPickupDate)
		suite.Equal(oldShipment.RequestedPickupDate.Format(time.ANSIC), time.Time(*okResponse.Payload.RequestedPickupDate).Format(time.ANSIC))
	})

	suite.Run("Failed case if estimated weight added after scheduled pickup date", func() {
		handler, move := setupTestData()

		yesterday := now.AddDate(0, 0, -1)
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusApproved,
					ScheduledPickupDate: &yesterday,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
		}
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", oldShipment.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Contains(errResponse.Payload.InvalidFields, "primeEstimatedWeight")
	})

	// Prime-specific validations tested below
	suite.Run("PATCH Success 200 RequiredDeliveryDate updated on scheduledPickupDate update", func() {
		handler, move := setupTestData()

		address := factory.BuildAddress(suite.DB(), nil, nil)
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

		hhgShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
					Status:       models.MTOShipmentStatusApproved,
					ApprovedDate: &now,
				},
			},
		}, nil)

		ntsShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
					Status:       models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
		}, nil)

		NTSRecordedWeight := unit.Pound(1400)
		ntsrShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:      models.MTOShipmentTypeHHGOutOfNTSDom,
					NTSRecordedWeight: &NTSRecordedWeight,
					Status:            models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		tenDaysFromNow := now.AddDate(0, 0, 11)
		schedDate := strfmt.Date(tenDaysFromNow)

		testCases := []struct {
			shipment models.MTOShipment
			payload  primemessages.UpdateMTOShipment
		}{
			{hhgShipment, primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
				ScheduledPickupDate:  &schedDate,
			}},
			{ntsShipment, primemessages.UpdateMTOShipment{
				PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
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
				Body:          &testCase.payload, //#nosec G601
				IfMatch:       eTag,
			}

			// Validate incoming payload
			suite.NoError(params.Body.Validate(strfmt.Default))

			response := handler.Handle(params)

			suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
			okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
			responsePayload := okResponse.Payload

			// Validate outgoing payload
			suite.NoError(responsePayload.Validate(strfmt.Default))

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

	suite.Run("PATCH sends back unprocessable response when dest address is updated for approved shipment", func() {
		handler, move := setupTestData()

		// Create shipment with populated estimated weight and scheduled date
		tenDaysFromNow := now.AddDate(0, 0, 11)
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		// setting shipment status to approved
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					ApprovedDate:         &now,
					PrimeEstimatedWeight: &primeEstimatedWeight,
					ScheduledPickupDate:  &tenDaysFromNow,
				},
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
		}, nil)

		// adding destination address to update to get back error
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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)

	})

	suite.Run("PATCH Success 200 RequiredDeliveryDate updated on destinationAddress creation", func() {
		// Under test: updateMTOShipmentHandler.Handle, RequiredDeliveryDate logic
		// Mocked:     Planner
		// Set up:     We use a shipment with primeEstimatedWeight and ScheduledPickupDate set
		//             Update with new destinationAddress
		// Expected:   Handler should return OK, new DestinationAddress should be saved
		//             requiredDeliveryDate should be set to 12 days from scheduledPickupDate
		handler, move := setupTestData()

		// Create shipment with populated estimated weight and scheduled date
		tenDaysFromNow := now.AddDate(0, 0, 11)
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		oldShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusSubmitted,
					ApprovedDate:         &now,
					PrimeEstimatedWeight: &primeEstimatedWeight,
					ScheduledPickupDate:  &tenDaysFromNow,
				},
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
		}, nil)

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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))

		// Confirm destination address in payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.EqualAddressPayload(&update.DestinationAddress.Address, &responsePayload.DestinationAddress.Address, false)

		// Confirm that auto-generated requiredDeliveryDate matches expected value
		expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 12)
		suite.EqualDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.Run("PATCH Success 200 RequiredDeliveryDate for Alaska", func() {
		// Under test: updateMTOShipmentHandler.Handle, RequiredDeliveryDate logic
		// Mocked:     Planner
		// Set up:     We use a shipment with an Alaska Address
		//             Update with new DestinationAddress
		// Expected:   Handler should return OK, new DestinationAddress should be saved
		//             requiredDeliveryDate should be set to 12 + 10 = 22 days from scheduledPickupDate
		//             which is a special rule for Alaska
		handler, move := setupTestData()

		// Create shipment with Alaska destination
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusApproved,
					ApprovedDate: &now,
				},
			},
			{
				Model: models.Address{
					PostalCode: "12345",
					State:      "AK",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// CREATE REQUEST
		// Update with scheduledPickupDate and PrimeEstimatedWeight
		tenDaysFromNow := now.AddDate(0, 0, 11)
		schedDate := strfmt.Date(tenDaysFromNow)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))

		// Check that updated fields are saved
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
		suite.NotNil(responsePayload.ScheduledPickupDate)

		// Check that RDD is set to 12 + 10 days after scheduled pickup date
		expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 22)
		suite.EqualDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})

	suite.Run("PATCH Success 200 RequiredDeliveryDate for Adak, Alaska", func() {
		// Under test: updateMTOShipmentHandler.Handle, RequiredDeliveryDate logic
		// Mocked:     Planner
		// Set up:     We use a shipment with an Alaska Address, specifically Adak
		//             Update with new DestinationAddress
		// Expected:   Handler should return OK, new DestinationAddress should be saved
		//             requiredDeliveryDate should be set to 12 + 20 = 32 days from scheduledPickupDate,
		//             which is a special rule for Adak (look at it on a map!)
		handler, move := setupTestData()

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusApproved,
					ApprovedDate: &now,
				},
			},
			{
				Model: models.Address{
					PostalCode: "99546",
					City:       "Adak",
					State:      "AK",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// CREATE REQUEST
		// Update with scheduledPickupDate and PrimeEstimatedWeight
		tenDaysFromNow := now.AddDate(0, 0, 11)
		schedDate := strfmt.Date(tenDaysFromNow)
		payload := primemessages.UpdateMTOShipment{
			PrimeEstimatedWeight: handlers.FmtPoundPtr(&primeEstimatedWeight),
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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		// CHECK RESPONSE
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))

		// Check that updated fields are saved
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
		suite.NotNil(responsePayload.ScheduledPickupDate)

		// Check that RDD is set to 12 + 20 days after scheduled pickup date
		expectedRDD := time.Time(*responsePayload.ScheduledPickupDate).AddDate(0, 0, 32)
		suite.EqualDatePtr(&expectedRDD, responsePayload.RequiredDeliveryDate)

	})
}

func (suite *HandlerSuite) TestUpdateMTOShipmentStatusHandler() {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	moveRouter := moveservices.NewMoveRouter()
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto_shipments/%s/status", uuid.Nil.String()), nil)

	setupTestData := func() (UpdateMTOShipmentStatusHandler, models.MTOShipment) {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetHHGPlanner(planner)
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		handler := UpdateMTOShipmentStatusHandler{
			handlerConfig,
			mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator),
			mtoshipment.NewMTOShipmentStatusUpdater(builder,
				mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()), planner),
		}

		// Set up Prime-available move
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusCancellationRequested,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		return handler, shipment
	}

	suite.Run("200 SUCCESS - Updated CANCELLATION_REQUESTED to CANCELED", func() {
		handler, shipment := setupTestData()
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       etag.GenerateEtag(shipment.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(string(models.MTOShipmentStatusCanceled), okResponse.Payload.Status)
		suite.Equal(shipment.MoveTaskOrderID.String(), okResponse.Payload.MoveTaskOrderID.String())
		suite.NotZero(okResponse.Payload.ETag)
	})

	suite.Run("404 FAIL - Bad shipment ID", func() {
		handler, shipment := setupTestData()

		badUUID := uuid.Must(uuid.NewV4())
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(badUUID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       etag.GenerateEtag(shipment.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusNotFound{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusNotFound)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*errResponse.Payload.Detail, badUUID.String())
	})

	suite.Run("404 FAIL - Shipment was not Prime-available", func() {
		handler, _ := setupTestData()

		nonPrimeShipment := factory.BuildMTOShipment(suite.DB(), nil, nil) // default is non-Prime available
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(nonPrimeShipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       etag.GenerateEtag(nonPrimeShipment.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusNotFound{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusNotFound)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*errResponse.Payload.Detail, nonPrimeShipment.ID.String())
	})

	suite.Run("412 FAIL - Stale eTag", func() {
		handler, shipment := setupTestData()

		staleShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MoveTaskOrderID: shipment.MoveTaskOrderID,
					Status:          models.MTOShipmentStatusCancellationRequested,
				},
			},
		}, nil)
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(staleShipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       "eTag",
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusPreconditionFailed{}, response)
		responsePayload := response.(*mtoshipmentops.UpdateMTOShipmentStatusPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("409 FAIL - Current status was not CANCELLATION_REQUESTED", func() {
		// Under test:       UpdateMTOShipmentStatusHandler
		// Mocked:           Planner
		// Set up:           Create a shipment with Canceled status, attempt to update to Canceled status
		// Expected outcome: Error since you can only cancel a shipment with CancellationRequested.
		handler, shipment := setupTestData()

		// Create a shipment in Canceled Status
		staleShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MoveTaskOrderID: shipment.MoveTaskOrderID,
					Status:          models.MTOShipmentStatusCanceled,
				},
			},
		}, nil)

		// Attempt to cancel again
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(staleShipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusCanceled)},
			IfMatch:       etag.GenerateEtag(staleShipment.UpdatedAt),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusConflict{}, response)
		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusConflict)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*errResponse.Payload.Detail, string(models.MTOShipmentStatusCanceled))
	})

	suite.Run("422 FAIL - Tried to use a status other than CANCELED", func() {
		_, shipment := setupTestData()

		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &primemessages.UpdateMTOShipmentStatus{Status: string(models.MTOShipmentStatusApproved)},
			IfMatch:       etag.GenerateEtag(shipment.UpdatedAt),
		}
		// Run swagger validations - should fail
		suite.Error(params.Body.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestDeleteMTOShipmentHandler() {
	setupTestData := func() DeleteMTOShipmentHandler {
		builder := query.NewQueryBuilder()
		moveRouter := moveservices.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter,
		)
		deleter := mtoshipment.NewPrimeShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		handlerConfig := suite.HandlerConfig()
		handler := DeleteMTOShipmentHandler{
			handlerConfig,
			deleter,
		}
		return handler
	}
	request := httptest.NewRequest("DELETE", "/shipments/{MtoShipmentID}", nil)

	suite.Run("Returns 204 when all validations pass", func() {
		handler := setupTestData()
		now := time.Now()
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		move.AvailableToPrimeAt = &now
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := mtoshipmentops.DeleteMTOShipmentParams{
			HTTPRequest:   request,
			MtoShipmentID: *handlers.FmtUUID(ppmShipment.ShipmentID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.DeleteMTOShipmentNoContent{}, response)

		// Validate outgoing payload: no payload
	})

	suite.Run("Returns a 403 when deleting a non-PPM shipment", func() {
		handler := setupTestData()
		now := time.Now()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
				},
			},
		}, nil)

		deletionParams := mtoshipmentops.DeleteMTOShipmentParams{
			HTTPRequest:   request,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)

		suite.IsType(&mtoshipmentops.DeleteMTOShipmentForbidden{}, response)
		responsePayload := response.(*mtoshipmentops.DeleteMTOShipmentForbidden).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Returns 404 when deleting a move not available to prime", func() {
		handler := setupTestData()
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: nil,
				},
			},
		}, nil)
		deletionParams := mtoshipmentops.DeleteMTOShipmentParams{
			HTTPRequest:   request,
			MtoShipmentID: *handlers.FmtUUID(ppmShipment.ShipmentID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteMTOShipmentNotFound{}, response)
		responsePayload := response.(*mtoshipmentops.DeleteMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Returns 409 - Conflict error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		deleter := &mocks.ShipmentDeleter{}
		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ConflictError{})
		handlerConfig := suite.HandlerConfig()
		handler := DeleteMTOShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteMTOShipmentParams{
			HTTPRequest:   request,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteMTOShipmentConflict{}, response)
		responsePayload := response.(*mtoshipmentops.DeleteMTOShipmentConflict).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(responsePayload)
	})

	suite.Run("Returns 422 - Unprocessable Entity error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		deleter := &mocks.ShipmentDeleter{}
		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.UnprocessableEntityError{})
		handlerConfig := suite.HandlerConfig()
		handler := DeleteMTOShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteMTOShipmentParams{
			HTTPRequest:   request,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteMTOShipmentUnprocessableEntity{}, response)
		responsePayload := response.(*mtoshipmentops.DeleteMTOShipmentUnprocessableEntity).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(responsePayload)
	})

	suite.Run("Returns 500 - Server error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		deleter := &mocks.ShipmentDeleter{}
		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.EventError{})
		handlerConfig := suite.HandlerConfig()
		handler := DeleteMTOShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := mtoshipmentops.DeleteMTOShipmentParams{
			HTTPRequest:   request,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&mtoshipmentops.DeleteMTOShipmentInternalServerError{}, response)
		responsePayload := response.(*mtoshipmentops.DeleteMTOShipmentInternalServerError).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}

func getFakeAddress() struct{ primemessages.Address } {
	// Use UUID to generate truly random address string
	streetAddr := fmt.Sprintf("%s %s", uuid.Must(uuid.NewV4()).String(), fakedata.RandomStreetAddress())
	// Using same zip so not a good helper for tests testing zip calculations
	return struct{ primemessages.Address }{
		Address: primemessages.Address{
			City:           models.StringPointer("San Diego"),
			PostalCode:     models.StringPointer("92102"),
			State:          models.StringPointer("CA"),
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
		"TertiaryPickupAddress",
		"TertiaryDeliveryAddress",
		"MTOAgents").Find(&dbShipment, id)
	suite.Nil(err)
	return dbShipment
}
