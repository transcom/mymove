package ghcv2api

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcv2api/ghcv2operations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcv2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	"github.com/transcom/mymove/pkg/swagger/nullable"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) getUpdateShipmentParams(originalShipment models.MTOShipment) mtoshipmentops.UpdateMTOShipmentParams {
	servicesCounselor := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	servicesCounselor.User.Roles = append(servicesCounselor.User.Roles, roles.Role{
		RoleType: roles.RoleTypeServicesCounselor,
	})
	pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
	pickupAddress.StreetAddress1 = "123 Fake Test St NW"
	destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"
	customerRemarks := "help"
	counselorRemarks := "counselor approved"
	billableWeightCap := int64(8000)
	billableWeightJustification := "Unable to perform reweigh because shipment was already unloaded."
	mtoAgent := factory.BuildMTOAgent(suite.DB(), nil, nil)
	agents := ghcv2messages.MTOAgents{&ghcv2messages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: string(mtoAgent.MTOAgentType),
	}}

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", originalShipment.MoveTaskOrderID.String(), originalShipment.ID.String()), nil)
	req = suite.AuthenticateOfficeRequest(req, servicesCounselor)

	eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

	now := strfmt.Date(time.Now())
	payload := ghcv2messages.UpdateShipment{
		BillableWeightJustification: &billableWeightJustification,
		BillableWeightCap:           &billableWeightCap,
		RequestedPickupDate:         &now,
		RequestedDeliveryDate:       &now,
		ShipmentType:                ghcv2messages.MTOShipmentTypeHHG,
		CustomerRemarks:             &customerRemarks,
		CounselorRemarks:            &counselorRemarks,
		Agents:                      agents,
		TacType:                     nullable.NewString("NTS"),
		SacType:                     nullable.NewString(""),
	}
	payload.DestinationAddress.Address = ghcv2messages.Address{
		City:           &destinationAddress.City,
		Country:        destinationAddress.Country,
		PostalCode:     &destinationAddress.PostalCode,
		State:          &destinationAddress.State,
		StreetAddress1: &destinationAddress.StreetAddress1,
		StreetAddress2: destinationAddress.StreetAddress2,
		StreetAddress3: destinationAddress.StreetAddress3,
	}
	payload.PickupAddress.Address = ghcv2messages.Address{
		City:           &pickupAddress.City,
		Country:        pickupAddress.Country,
		PostalCode:     &pickupAddress.PostalCode,
		State:          &pickupAddress.State,
		StreetAddress1: &pickupAddress.StreetAddress1,
		StreetAddress2: pickupAddress.StreetAddress2,
		StreetAddress3: pickupAddress.StreetAddress3,
	}

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(originalShipment.ID),
		Body:        &payload,
		IfMatch:     eTag,
	}

	return params
}

func (suite *HandlerSuite) TestUpdateShipmentHandler() {
	planner := &routemocks.Planner{}
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

	suite.Run("Successful PATCH - Integration Test", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		hhgLOAType := models.LOATypeHHG
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:             models.MTOShipmentStatusSubmitted,
					UsesExternalVendor: true,
					TACType:            &hhgLOAType,
					Diversion:          true,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		suite.Equal(oldShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(params.Body.BillableWeightCap, updatedShipment.BillableWeightCap)
		suite.Equal(params.Body.BillableWeightJustification, updatedShipment.BillableWeightJustification)
		suite.Equal(params.Body.CounselorRemarks, updatedShipment.CounselorRemarks)
		suite.Equal(params.Body.PickupAddress.StreetAddress1, updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(params.Body.DestinationAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), updatedShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.Agents[0].FirstName, updatedShipment.MtoAgents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, updatedShipment.MtoAgents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, updatedShipment.MtoAgents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, updatedShipment.MtoAgents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, updatedShipment.MtoAgents[0].AgentType)
		suite.Equal(oldShipment.ID.String(), string(updatedShipment.MtoAgents[0].MtoShipmentID))
		suite.NotEmpty(updatedShipment.MtoAgents[0].ID)
		suite.Equal(params.Body.RequestedDeliveryDate.String(), updatedShipment.RequestedDeliveryDate.String())
		suite.Equal(*params.Body.TacType.Value, string(*updatedShipment.TacType))
		suite.Nil(updatedShipment.SacType)

		// don't update non-nullable booleans if they're not passed in
		suite.Equal(oldShipment.Diversion, updatedShipment.Diversion)
		suite.Equal(oldShipment.UsesExternalVendor, updatedShipment.UsesExternalVendor)
	})

	suite.Run("Successful PATCH - Integration Test (PPM)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		hasProGear := true
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
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
		sitLocation := ghcv2messages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1700)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3000)
		proGearWeight := unit.Pound(300)
		spouseProGearWeight := unit.Pound(200)
		estimatedIncentive := 654321
		sitEstimatedCost := 67500

		params := suite.getUpdateShipmentParams(ppmShipment.Shipment)
		params.Body.ShipmentType = ghcv2messages.MTOShipmentTypePPM
		params.Body.PpmShipment = &ghcv2messages.UpdatePPMShipment{
			ActualMoveDate:                 handlers.FmtDatePtr(&actualMoveDate),
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
		suite.Equal(handlers.FmtDatePtr(&actualMoveDate), updatedShipment.PpmShipment.ActualMoveDate)
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
		suite.Equal(int64(sitEstimatedCost), *updatedShipment.PpmShipment.SitEstimatedCost)
		suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), updatedShipment.PpmShipment.EstimatedWeight)
		suite.Equal(int64(estimatedIncentive), *updatedShipment.PpmShipment.EstimatedIncentive)
		suite.Equal(handlers.FmtBool(hasProGear), updatedShipment.PpmShipment.HasProGear)
		suite.Equal(handlers.FmtPoundPtr(&proGearWeight), updatedShipment.PpmShipment.ProGearWeight)
		suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), updatedShipment.PpmShipment.SpouseProGearWeight)
	})

	suite.Run("PATCH failure - 400 -- nil body", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)
		params.Body = nil

		// Validate incoming payload: nil body (the point of this test)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		uuidString := handlers.FmtUUID(uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763"))
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)
		params.ShipmentID = *uuidString

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)
		params.IfMatch = "intentionally-bad-if-match-header-value"

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 412 -- shipment shouldn't be updatable", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		params := suite.getUpdateShipmentParams(oldShipment)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentForbidden{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentForbidden).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ShipmentUpdater{}
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			&mockUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateShipmentV1",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
