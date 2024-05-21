package primeapiv3

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapiv3/payloads"
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
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	moveRouter := moveservices.NewMoveRouter()
	fetcher := fetch.NewFetcher(builder)
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV2(builder, fetcher, moveRouter, addressCreator)
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

	var pickupAddress primev3messages.Address
	var secondaryPickupAddress primev3messages.Address
	var destinationAddress primev3messages.Address
	var secondaryDestinationAddress primev3messages.Address

	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, suite.TestNotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

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
		pickupAddress = primev3messages.Address{
			City:           &newAddress.City,
			Country:        newAddress.Country,
			PostalCode:     &newAddress.PostalCode,
			State:          &newAddress.State,
			StreetAddress1: &newAddress.StreetAddress1,
			StreetAddress2: newAddress.StreetAddress2,
			StreetAddress3: newAddress.StreetAddress3,
		}
		newAddress = factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
		destinationAddress = primev3messages.Address{
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				Agents:               nil,
				CustomerRemarks:      nil,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
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
		suite.Require().Equal(createMTOShipmentPayload.Status, primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.PrimeEstimatedWeight, params.Body.PrimeEstimatedWeight)
	})

	suite.Run("Successful POST/PATCH - Integration Test (PPM)", func() {
		// Under Test: CreateMTOShipment handler code
		// Setup:      Create a PPM shipment on an available move
		// Expected:   Successful submission, status should be SUBMITTED
		handler, move := setupTestData()
		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		counselorRemarks := "Some counselor remarks"
		expectedDepartureDate := time.Now().AddDate(0, 0, 10)
		sitExpected := true
		sitLocation := primev3messages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1500)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3200)
		hasProGear := true
		proGearWeight := unit.Pound(400)
		spouseProGearWeight := unit.Pound(250)
		estimatedIncentive := 123456
		sitEstimatedCost := 67500

		address1 := models.Address{
			StreetAddress1: "some address",
			City:           "city",
			State:          "CA",
			PostalCode:     "90210",
		}
		address2 := models.Address{
			StreetAddress1: "some address",
			City:           "city",
			State:          "IL",
			PostalCode:     "62225",
		}

		expectedPickupAddress := address1
		pickupAddress = primev3messages.Address{
			City:           &expectedPickupAddress.City,
			Country:        expectedPickupAddress.Country,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedSecondaryPickupAddress := address2
		secondaryPickupAddress = primev3messages.Address{
			City:           &expectedSecondaryPickupAddress.City,
			Country:        expectedSecondaryPickupAddress.Country,
			PostalCode:     &expectedSecondaryPickupAddress.PostalCode,
			State:          &expectedSecondaryPickupAddress.State,
			StreetAddress1: &expectedSecondaryPickupAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryPickupAddress.StreetAddress2,
			StreetAddress3: expectedSecondaryPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := address1
		destinationAddress = primev3messages.Address{
			City:           &expectedDestinationAddress.City,
			Country:        expectedDestinationAddress.Country,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}

		expectedSecondaryDestinationAddress := address2
		secondaryDestinationAddress = primev3messages.Address{
			City:           &expectedSecondaryDestinationAddress.City,
			Country:        expectedSecondaryDestinationAddress.Country,
			PostalCode:     &expectedSecondaryDestinationAddress.PostalCode,
			State:          &expectedSecondaryDestinationAddress.State,
			StreetAddress1: &expectedSecondaryDestinationAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryDestinationAddress.StreetAddress2,
			StreetAddress3: expectedSecondaryDestinationAddress.StreetAddress3,
		}

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:  handlers.FmtUUID(move.ID),
				ShipmentType:     primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypePPM),
				CounselorRemarks: &counselorRemarks,
				PpmShipment: &primev3messages.CreatePPMShipment{
					ExpectedDepartureDate:       handlers.FmtDate(expectedDepartureDate),
					PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
					SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
					DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
					SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
					SitExpected:                 &sitExpected,
					SitLocation:                 &sitLocation,
					SitEstimatedWeight:          handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:       handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate:   handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:             handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                  &hasProGear,
					ProGearWeight:               handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:         handlers.FmtPoundPtr(&spouseProGearWeight),
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
		suite.Equal(primev3messages.MTOShipmentTypePPM, createdShipment.ShipmentType)
		suite.Equal(primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, createdShipment.Status)
		suite.Equal(&counselorRemarks, createdShipment.CounselorRemarks)

		suite.Equal(createdShipment.ID.String(), createdPPM.ShipmentID.String())
		suite.Equal(primev3messages.PPMShipmentStatusSUBMITTED, createdPPM.Status)
		suite.Equal(handlers.FmtDatePtr(&expectedDepartureDate), createdPPM.ExpectedDepartureDate)
		suite.Equal(address1.PostalCode, *createdPPM.PickupAddress.PostalCode)
		suite.Equal(address1.PostalCode, *createdPPM.DestinationAddress.PostalCode)
		suite.Equal(address2.PostalCode, *createdPPM.SecondaryPickupAddress.PostalCode)
		suite.Equal(address2.PostalCode, *createdPPM.SecondaryDestinationAddress.PostalCode)
		suite.Equal(&sitExpected, createdPPM.SitExpected)
		suite.Equal(&sitLocation, createdPPM.SitLocation)
		suite.True(*models.BoolPointer(*createdPPM.HasSecondaryPickupAddress))
		suite.True(*models.BoolPointer(*createdPPM.HasSecondaryDestinationAddress))
		suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), createdPPM.SitEstimatedWeight)
		suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), createdPPM.SitEstimatedEntryDate)
		suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), createdPPM.SitEstimatedDepartureDate)
		suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), createdPPM.EstimatedWeight)
		suite.Equal(handlers.FmtBool(hasProGear), createdPPM.HasProGear)
		suite.Equal(handlers.FmtPoundPtr(&proGearWeight), createdPPM.ProGearWeight)
		suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), createdPPM.SpouseProGearWeight)
		suite.Equal(int64(estimatedIncentive), *createdPPM.EstimatedIncentive)
		suite.Equal(int64(sitEstimatedCost), *createdPPM.SitEstimatedCost)

		// ************
		// PATCH TESTS
		// ************
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Times(2)

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		patchHandler := UpdateMTOShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
		}

		patchReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", createdPPM.ShipmentID.String()), nil)

		var mtoShipment models.MTOShipment
		err := suite.DB().Find(&mtoShipment, createdPPM.ShipmentID)
		suite.NoError(err)
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		patchParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   patchReq,
			MtoShipmentID: createdPPM.ShipmentID,
			IfMatch:       eTag,
		}
		patchParams.Body = &primev3messages.UpdateMTOShipment{
			ShipmentType: primev3messages.MTOShipmentTypePPM,
		}
		// *************************************************************************************
		// *************************************************************************************
		// Run it without any flags, no deletes should occur for secondary addresses
		// *************************************************************************************
		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			HasProGear: &hasProGear,
		}

		// Validate incoming payload
		suite.NoError(patchParams.Body.Validate(strfmt.Default))

		patchResponse := patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse := patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment := okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM := updatedShipment.PpmShipment
		suite.NotNil(updatedPPM.SecondaryPickupAddress)
		suite.NotNil(updatedPPM.SecondaryDestinationAddress)
		suite.True(*models.BoolPointer(*updatedPPM.HasSecondaryPickupAddress))
		suite.True(*models.BoolPointer(*updatedPPM.HasSecondaryDestinationAddress))

		// *************************************************************************************
		// *************************************************************************************
		// Run it second time, but really delete secondary addresses with has flags set to false
		// *************************************************************************************
		err = suite.DB().Find(&mtoShipment, createdPPM.ShipmentID)
		suite.NoError(err)
		eTag = etag.GenerateEtag(mtoShipment.UpdatedAt)
		patchParams.IfMatch = eTag
		patchParams.Body.PpmShipment = &primev3messages.UpdatePPMShipment{
			HasProGear:                     &hasProGear,
			HasSecondaryPickupAddress:      models.BoolPointer(false),
			HasSecondaryDestinationAddress: models.BoolPointer(false),
		}
		suite.NoError(patchParams.Body.Validate(strfmt.Default))
		patchResponse = patchHandler.Handle(patchParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, patchResponse)
		okPatchResponse = patchResponse.(*mtoshipmentops.UpdateMTOShipmentOK)
		updatedShipment = okPatchResponse.Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		updatedPPM = updatedShipment.PpmShipment
		// secondary should be all nils
		suite.Nil(updatedPPM.SecondaryPickupAddress)
		suite.Nil(updatedPPM.SecondaryDestinationAddress)
		suite.False(*models.BoolPointer(*updatedPPM.HasSecondaryPickupAddress))
		suite.False(*models.BoolPointer(*updatedPPM.HasSecondaryDestinationAddress))
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:     handlers.FmtUUID(move.ID),
				Agents:              nil,
				CustomerRemarks:     nil,
				PointOfContact:      "John Doe",
				RequestedPickupDate: handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:        primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:       struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:  struct{ primev3messages.Address }{destinationAddress},
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
		suite.Require().Equal(createMTOShipmentPayload.Status, primev3messages.MTOShipmentWithoutServiceItemsStatusSUBMITTED, "MTO Shipment should have been submitted")
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
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
		agent := &primev3messages.MTOAgent{
			ID:            badID,
			MtoShipmentID: badID,
			FirstName:     handlers.FmtString("Mary"),
		}
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				Agents:               primev3messages.MTOAgents{agent},
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      &badID,
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(unavailableMove.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
				PickupAddress:        struct{ primev3messages.Address }{pickupAddress},
				DestinationAddress:   struct{ primev3messages.Address }{destinationAddress},
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
			Body: &primev3messages.CreateMTOShipment{
				MoveTaskOrderID:      handlers.FmtUUID(move.ID),
				PointOfContact:       "John Doe",
				PrimeEstimatedWeight: handlers.FmtInt64(1200),
				RequestedPickupDate:  handlers.FmtDatePtr(models.TimePointer(time.Now())),
				ShipmentType:         primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeHHG),
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
