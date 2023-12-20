package primeapiv2

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
	"github.com/transcom/mymove/pkg/factory"
	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	moveRouter := moveservices.NewMoveRouter()
	fetcher := fetch.NewFetcher(builder)
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	ppmEstimator := mocks.PPMEstimator{}
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
	shipmentRouter := mtoshipment.NewShipmentRouter()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		builder,
		mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter),
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

func (suite *HandlerSuite) TestCreateNonSITAddressUpdateRequestHandler() {
	req := httptest.NewRequest("POST", "/mto-shipments/{mtoShipmentID}/shipment-address-updates", nil)

	makeSubtestData := func() mtoshipmentops.CreateNonSITAddressUpdateRequestParams {
		contractorRemark := "This is a contractor remark"
		body := primemessages.CreateNonSITAddressUpdateRequest{
			ContractorRemarks: &contractorRemark,
			NewAddress: &primemessages.Address{
				City:           swag.String("Beverly Hills"),
				PostalCode:     swag.String("90210"),
				State:          swag.String("CA"),
				StreetAddress1: swag.String("1234 N. 1st Street"),
			},
		}

		params := mtoshipmentops.CreateNonSITAddressUpdateRequestParams{
			HTTPRequest: req,
			Body:        &body,
		}

		return params

	}
	suite.Run("POST failure - 422 Unprocessable Entity Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := CreateNonSITAddressUpdateRequestHandler{
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
		suite.IsType(&mtoshipmentops.CreateNonSITAddressUpdateRequestUnprocessableEntity{}, response)
		errResponse := response.(*mtoshipmentops.CreateNonSITAddressUpdateRequestUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 409 Request conflict reponse Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := CreateNonSITAddressUpdateRequestHandler{
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
		suite.IsType(&mtoshipmentops.CreateNonSITAddressUpdateRequestConflict{}, response)
		errResponse := response.(*mtoshipmentops.CreateNonSITAddressUpdateRequestConflict)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 404 Not Found response error", func() {

		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := CreateNonSITAddressUpdateRequestHandler{
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
		suite.IsType(&mtoshipmentops.CreateNonSITAddressUpdateRequestNotFound{}, response)
		errResponse := response.(*mtoshipmentops.CreateNonSITAddressUpdateRequestNotFound)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("500 server error", func() {

		subtestData := makeSubtestData()
		mockCreator := mocks.ShipmentAddressUpdateRequester{}
		handler := CreateNonSITAddressUpdateRequestHandler{
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
		suite.IsType(&mtoshipmentops.CreateNonSITAddressUpdateRequestInternalServerError{}, response)
		errResponse := response.(*mtoshipmentops.CreateNonSITAddressUpdateRequestInternalServerError)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
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
		"MTOAgents").Find(&dbShipment, id)
	suite.Nil(err)
	return dbShipment
}
