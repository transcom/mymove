package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	ppm "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	ppmshipment "github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func preLoadData(suite *HandlerSuite) {
	testdatagen.FetchOrMakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
		GHCDieselFuelPrice: models.GHCDieselFuelPrice{
			FuelPriceInMillicents: unit.Millicents(281400),
			PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
		},
	})

	originDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea:      "056",
			ServicesSchedule: 3,
			SITPDSchedule:    3,
		},
		ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
	})

	testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			Contract:             originDomesticServiceArea.Contract,
			ContractID:           originDomesticServiceArea.ContractID,
			StartDate:            time.Date(2019, time.June, 1, 0, 0, 0, 0, time.UTC),
			EndDate:              time.Date(2020, time.May, 31, 0, 0, 0, 0, time.UTC),
			Escalation:           1.0,
			EscalationCompounded: 1.0,
		},
	})

	testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            originDomesticServiceArea.Contract,
			ContractID:          originDomesticServiceArea.ContractID,
			DomesticServiceArea: originDomesticServiceArea,
			Zip3:                "902",
		},
	})

	destDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			Contract:    originDomesticServiceArea.Contract,
			ContractID:  originDomesticServiceArea.ContractID,
			ServiceArea: "208",
		},
	})

	testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            destDomesticServiceArea.Contract,
			ContractID:          destDomesticServiceArea.ContractID,
			DomesticServiceArea: destDomesticServiceArea,
			Zip3:                "308",
		},
	})

	testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			Contract:              originDomesticServiceArea.Contract,
			ContractID:            originDomesticServiceArea.ContractID,
			DomesticServiceArea:   originDomesticServiceArea,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			WeightLower:           unit.Pound(500),
			WeightUpper:           unit.Pound(4999),
			MilesLower:            2001,
			MilesUpper:            2500,
			PriceMillicents:       unit.Millicents(412400),
		},
	})

	testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			Contract:              originDomesticServiceArea.Contract,
			ContractID:            originDomesticServiceArea.ContractID,
			DomesticServiceArea:   originDomesticServiceArea,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			WeightLower:           unit.Pound(500),
			WeightUpper:           unit.Pound(4999),
			MilesLower:            2001,
			MilesUpper:            2500,
			IsPeakPeriod:          true,
			PriceMillicents:       unit.Millicents(437600),
		},
	})

	testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			Contract:              originDomesticServiceArea.Contract,
			ContractID:            originDomesticServiceArea.ContractID,
			DomesticServiceArea:   originDomesticServiceArea,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			WeightLower:           unit.Pound(5000),
			WeightUpper:           unit.Pound(9999),
			MilesLower:            2001,
			MilesUpper:            2500,
			PriceMillicents:       unit.Millicents(606800),
		},
	})

	dopService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            originDomesticServiceArea.ContractID,
			Contract:              originDomesticServiceArea.Contract,
			ServiceID:             dopService.ID,
			Service:               dopService,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			DomesticServiceArea:   originDomesticServiceArea,
			IsPeakPeriod:          false,
			PriceCents:            unit.Cents(404),
		},
	})

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            originDomesticServiceArea.ContractID,
			Contract:              originDomesticServiceArea.Contract,
			ServiceID:             dopService.ID,
			Service:               dopService,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			DomesticServiceArea:   originDomesticServiceArea,
			IsPeakPeriod:          true,
			PriceCents:            unit.Cents(465),
		},
	})

	ddpService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            destDomesticServiceArea.ContractID,
			Contract:              destDomesticServiceArea.Contract,
			ServiceID:             ddpService.ID,
			Service:               ddpService,
			DomesticServiceAreaID: destDomesticServiceArea.ID,
			DomesticServiceArea:   destDomesticServiceArea,
			IsPeakPeriod:          false,
			PriceCents:            unit.Cents(832),
		},
	})

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            destDomesticServiceArea.ContractID,
			Contract:              destDomesticServiceArea.Contract,
			ServiceID:             ddpService.ID,
			Service:               ddpService,
			DomesticServiceAreaID: destDomesticServiceArea.ID,
			DomesticServiceArea:   destDomesticServiceArea,
			IsPeakPeriod:          true,
			PriceCents:            unit.Cents(957),
		},
	})

	dpkService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)

	testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticOtherPrice: models.ReDomesticOtherPrice{
			ContractID:   originDomesticServiceArea.ContractID,
			Contract:     originDomesticServiceArea.Contract,
			ServiceID:    dpkService.ID,
			Service:      dpkService,
			IsPeakPeriod: false,
			Schedule:     3,
			PriceCents:   7395,
		},
	})

	testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticOtherPrice: models.ReDomesticOtherPrice{
			ContractID:   originDomesticServiceArea.ContractID,
			Contract:     originDomesticServiceArea.Contract,
			ServiceID:    dpkService.ID,
			Service:      dpkService,
			IsPeakPeriod: true,
			Schedule:     3,
			PriceCents:   8000,
		},
	})

	dupkService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)

	testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticOtherPrice: models.ReDomesticOtherPrice{
			ContractID:   destDomesticServiceArea.ContractID,
			Contract:     destDomesticServiceArea.Contract,
			ServiceID:    dupkService.ID,
			Service:      dupkService,
			IsPeakPeriod: false,
			Schedule:     2,
			PriceCents:   597,
		},
	})

	testdatagen.FetchOrMakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticOtherPrice: models.ReDomesticOtherPrice{
			ContractID:   destDomesticServiceArea.ContractID,
			Contract:     destDomesticServiceArea.Contract,
			ServiceID:    dupkService.ID,
			Service:      dupkService,
			IsPeakPeriod: true,
			Schedule:     2,
			PriceCents:   650,
		},
	})

	dofsitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            originDomesticServiceArea.ContractID,
			Contract:              originDomesticServiceArea.Contract,
			ServiceID:             dofsitService.ID,
			Service:               dofsitService,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			DomesticServiceArea:   originDomesticServiceArea,
			IsPeakPeriod:          false,
			PriceCents:            1153,
		},
	})

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            originDomesticServiceArea.ContractID,
			Contract:              originDomesticServiceArea.Contract,
			ServiceID:             dofsitService.ID,
			Service:               dofsitService,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			DomesticServiceArea:   originDomesticServiceArea,
			IsPeakPeriod:          true,
			PriceCents:            1326,
		},
	})

	doasitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOASIT)

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            originDomesticServiceArea.ContractID,
			Contract:              originDomesticServiceArea.Contract,
			ServiceID:             doasitService.ID,
			Service:               doasitService,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			DomesticServiceArea:   originDomesticServiceArea,
			IsPeakPeriod:          false,
			PriceCents:            46,
		},
	})

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            originDomesticServiceArea.ContractID,
			Contract:              originDomesticServiceArea.Contract,
			ServiceID:             doasitService.ID,
			Service:               doasitService,
			DomesticServiceAreaID: originDomesticServiceArea.ID,
			DomesticServiceArea:   originDomesticServiceArea,
			IsPeakPeriod:          true,
			PriceCents:            53,
		},
	})

	ddfsitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            destDomesticServiceArea.ContractID,
			Contract:              destDomesticServiceArea.Contract,
			ServiceID:             ddfsitService.ID,
			Service:               ddfsitService,
			DomesticServiceAreaID: destDomesticServiceArea.ID,
			DomesticServiceArea:   destDomesticServiceArea,
			IsPeakPeriod:          false,
			PriceCents:            1612,
		},
	})

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            destDomesticServiceArea.ContractID,
			Contract:              destDomesticServiceArea.Contract,
			ServiceID:             ddfsitService.ID,
			Service:               ddfsitService,
			DomesticServiceAreaID: destDomesticServiceArea.ID,
			DomesticServiceArea:   destDomesticServiceArea,
			IsPeakPeriod:          true,
			PriceCents:            1854,
		},
	})

	ddasitService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDASIT)

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            destDomesticServiceArea.ContractID,
			Contract:              destDomesticServiceArea.Contract,
			ServiceID:             ddasitService.ID,
			Service:               ddasitService,
			DomesticServiceAreaID: destDomesticServiceArea.ID,
			DomesticServiceArea:   destDomesticServiceArea,
			IsPeakPeriod:          false,
			PriceCents:            55,
		},
	})

	testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
			ContractID:            destDomesticServiceArea.ContractID,
			Contract:              destDomesticServiceArea.Contract,
			ServiceID:             ddasitService.ID,
			Service:               ddasitService,
			DomesticServiceAreaID: destDomesticServiceArea.ID,
			DomesticServiceArea:   destDomesticServiceArea,
			IsPeakPeriod:          true,
			PriceCents:            63,
		},
	})
}

func (suite *HandlerSuite) TestGetPPMSITEstimatedCostHandler() {
	var ppmShipment models.PPMShipment
	newFakeSITEstimatedCost := models.CentPointer(unit.Cents(25500))

	suite.PreloadData(func() {
		preLoadData(suite)
	})

	setupData := func() {
		sitLocationDestination := models.SITLocationTypeDestination
		entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "30813"
		destinationAddress := &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: &streetAddress2,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		ppmShipment = factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate:     entryDate.Add(time.Hour * 24 * 30),
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &sitLocationDestination,
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					SITEstimatedWeight:        models.PoundPointer(1000),
					SITEstimatedCost:          newFakeSITEstimatedCost,
				},
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		ppmShipment.DestinationAddress = destinationAddress
		mockedPlanner := &routemocks.Planner{}
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "90210", "30813").Return(2294, nil)
	}
	setupData()

	setUpGetCostRequestAndParams := func() ppm.GetPPMSITEstimatedCostParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/sit_location/%s/sit-estimated-cost", ppmShipment.ID.String(), *ppmShipment.SITLocation)

		req := httptest.NewRequest("GET", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		return ppm.GetPPMSITEstimatedCostParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
			SitLocation:   string(*ppmShipment.SITLocation),
		}
	}

	setUpUpdateCostRequestAndParams := func(sitLocation *ghcmessages.SITLocationType) ppm.UpdatePPMSITParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/ppm-sit", ppmShipment.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)
		eTag := etag.GenerateEtag(ppmShipment.UpdatedAt)

		return ppm.UpdatePPMSITParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
			IfMatch:       eTag,
			Body: &ghcmessages.PPMShipmentSIT{
				SitLocation: sitLocation,
			},
		}
	}

	type ppmShipmentSubtestData struct {
		ppmShipmentUpdater services.PPMShipmentUpdater
		ppmEstimator       services.PPMEstimator
		ppmShipmentFetcher services.PPMShipmentFetcher
	}

	setUpForGetCostTests := func(sitEstimatedCost *unit.Cents, sitEstimatedError error) (subtestData ppmShipmentSubtestData) {
		ppmEstimator := mocks.PPMEstimator{}
		ppmEstimatedCostInfo := &models.PPMSITEstimatedCostInfo{}
		ppmEstimatedCostInfo.EstimatedSITCost = sitEstimatedCost
		ppmEstimatedCostInfo.PriceFirstDaySIT = sitEstimatedCost
		ppmEstimatedCostInfo.PriceAdditionalDaySIT = sitEstimatedCost
		ppmEstimator.
			On(
				"CalculatePPMSITEstimatedCostBreakdown",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(ppmEstimatedCostInfo, sitEstimatedError)

		subtestData.ppmShipmentFetcher = ppmshipment.NewPPMShipmentFetcher()
		subtestData.ppmEstimator = &ppmEstimator
		return subtestData
	}

	setUpForUpdateCostTests := func(sitEstimatedError error) (subtestData ppmShipmentSubtestData) {
		ppmEstimator := mocks.PPMEstimator{}

		ppmEstimator.
			On(
				"CalculatePPMSITEstimatedCostBreakdown",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(ppmShipment, sitEstimatedError)

		ppmShipmentUpdater := mocks.PPMShipmentUpdater{}

		ppmShipmentUpdater.
			On(
				"UpdatePPMShipmentSITEstimatedCost",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(&ppmShipment, sitEstimatedError)

		subtestData.ppmShipmentFetcher = ppmshipment.NewPPMShipmentFetcher()
		subtestData.ppmEstimator = &ppmEstimator
		subtestData.ppmShipmentUpdater = &ppmShipmentUpdater
		return subtestData
	}

	setUpGetCostHandler := func(mockPPMEstimator services.PPMEstimator, mockPPMShipmentFetcher services.PPMShipmentFetcher) GetPPMSITEstimatedCostHandler {
		return GetPPMSITEstimatedCostHandler{
			suite.createS3HandlerConfig(),
			mockPPMEstimator,
			mockPPMShipmentFetcher,
		}
	}

	setUpUpdateCostHandler := func(mockPPMShipmentUpdater services.PPMShipmentUpdater, mockPPMShipmentFetcher services.PPMShipmentFetcher) UpdatePPMSITHandler {
		return UpdatePPMSITHandler{
			suite.createS3HandlerConfig(),
			mockPPMShipmentUpdater,
			mockPPMShipmentFetcher,
		}
	}

	suite.Run("Get PPM SIT Estimated Cost - DESTINATION", func() {
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForGetCostTests(newFakeSITEstimatedCost, nil)

		handler := setUpGetCostHandler(subtestData.ppmEstimator, subtestData.ppmShipmentFetcher)
		params := setUpGetCostRequestAndParams()
		response := handler.Handle(params)

		if suite.IsType(&ppm.GetPPMSITEstimatedCostOK{}, response) {
			payload := response.(*ppm.GetPPMSITEstimatedCostOK).Payload

			suite.NoError(payload.Validate(strfmt.Default))
			suite.NotEqual(payload.SitCost, ppmShipment.SITEstimatedCost)
		}
	})

	suite.Run("FAIL to get PPM Shipment - SIT Estimated Cost - DESTINATION", func() {
		newFakeSITEstimatedCost := models.CentPointer(unit.Cents(62500))

		subtestData := setUpForGetCostTests(newFakeSITEstimatedCost, nil)

		handler := setUpGetCostHandler(subtestData.ppmEstimator, subtestData.ppmShipmentFetcher)
		params := setUpGetCostRequestAndParams()
		params.PpmShipmentID = strfmt.UUID("")
		response := handler.Handle(params)

		suite.IsType(&ppm.GetPPMSITEstimatedCostNotFound{}, response)
	})

	suite.Run("Update PPM SIT Estimated Cost - ORIGIN", func() {
		sitLocation := ghcmessages.SITLocationTypeORIGIN
		subtestData := setUpForUpdateCostTests(nil)

		handler := setUpUpdateCostHandler(subtestData.ppmShipmentUpdater, subtestData.ppmShipmentFetcher)
		params := setUpUpdateCostRequestAndParams(&sitLocation)
		response := handler.Handle(params)

		if suite.IsType(&ppm.UpdatePPMSITOK{}, response) {
			payload := response.(*ppm.UpdatePPMSITOK).Payload

			suite.NoError(payload.Validate(strfmt.Default))
			suite.NotEqual(payload.SitEstimatedCost, newFakeSITEstimatedCost)
		}
	})

	suite.Run("FAIL to update PPM SIT Estimated Cost - ORIGIN", func() {
		sitLocation := ghcmessages.SITLocationTypeORIGIN
		subtestData := setUpForUpdateCostTests(nil)

		handler := setUpUpdateCostHandler(subtestData.ppmShipmentUpdater, subtestData.ppmShipmentFetcher)
		params := setUpUpdateCostRequestAndParams(&sitLocation)
		params.PpmShipmentID = strfmt.UUID("")
		response := handler.Handle(params)

		suite.IsType(&ppm.UpdatePPMSITNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestSendPPMToCustomerHandler() {
	var ppmShipment models.PPMShipment
	newFakeSITEstimatedCost := models.CentPointer(unit.Cents(25500))

	suite.PreloadData(func() {
		preLoadData(suite)
	})

	setupData := func() {
		sitLocationDestination := models.SITLocationTypeDestination
		entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		streetAddress1 := "10642 N Second Ave"
		streetAddress2 := "Apt. 308"
		city := "Atco"
		state := "NJ"
		postalCode := "30813"
		destinationAddress := &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: &streetAddress2,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		ppmShipment = factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate:     entryDate.Add(time.Hour * 24 * 30),
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &sitLocationDestination,
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					SITEstimatedWeight:        models.PoundPointer(1000),
					SITEstimatedCost:          newFakeSITEstimatedCost,
				},
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		ppmShipment.DestinationAddress = destinationAddress
		mockedPlanner := &routemocks.Planner{}
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "90210", "30813").Return(2294, nil)
	}
	setupData()

	setUpPatchSendToCustomerAndParams := func() ppm.SendPPMToCustomerParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/send-to-customer", ppmShipment.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		return ppm.SendPPMToCustomerParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
		}
	}

	type ppmShipmentSubtestData struct {
		ppmShipmentMTOUpdater services.MoveTaskOrderUpdater
		ppmShipmentFetcher    services.PPMShipmentFetcher
	}

	setUpForSendPPMToCustomer := func() (subtestData ppmShipmentSubtestData) {
		subtestData.ppmShipmentFetcher = ppmshipment.NewPPMShipmentFetcher()
		mockMoveTaskOrderUpdater := mocks.MoveTaskOrderUpdater{}

		mockMoveTaskOrderUpdater.
			On(
				"UpdateStatusServiceCounselingSendPPMToCustomer",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.PPMShipment"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("*models.Move"),
			).
			Return(&ppmShipment, nil)

		subtestData.ppmShipmentMTOUpdater = &mockMoveTaskOrderUpdater

		return subtestData
	}

	setUpSendPPMToCustomerHandler := func(mockMoveTaskOrderUpdater services.MoveTaskOrderUpdater, mockPPMShipmentFetcher services.PPMShipmentFetcher) SendPPMToCustomerHandler {
		return SendPPMToCustomerHandler{
			suite.createS3HandlerConfig(),
			mockPPMShipmentFetcher,
			mockMoveTaskOrderUpdater,
		}
	}

	suite.Run("Successfully send PPM to Customer", func() {
		subtestData := setUpForSendPPMToCustomer()

		handler := setUpSendPPMToCustomerHandler(subtestData.ppmShipmentMTOUpdater, subtestData.ppmShipmentFetcher)
		params := setUpPatchSendToCustomerAndParams()
		response := handler.Handle(params)

		if suite.IsType(&ppm.SendPPMToCustomerOK{}, response) {
			payload := response.(*ppm.SendPPMToCustomerOK).Payload

			suite.NoError(payload.Validate(strfmt.Default))
		}
	})

	suite.Run("FAIL null ppm shipment id - send PPM to Customer", func() {
		subtestData := setUpForSendPPMToCustomer()

		handler := setUpSendPPMToCustomerHandler(subtestData.ppmShipmentMTOUpdater, subtestData.ppmShipmentFetcher)
		params := setUpPatchSendToCustomerAndParams()
		params.PpmShipmentID = strfmt.UUID("")
		response := handler.Handle(params)

		suite.IsType(&ppm.SendPPMToCustomerBadRequest{}, response)
	})

	suite.Run("FAIL can't get shipment - send PPM to Customer", func() {
		subtestData := setUpForSendPPMToCustomer()

		handler := setUpSendPPMToCustomerHandler(subtestData.ppmShipmentMTOUpdater, subtestData.ppmShipmentFetcher)
		params := setUpPatchSendToCustomerAndParams()
		params.PpmShipmentID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		response := handler.Handle(params)

		suite.IsType(&ppm.SendPPMToCustomerNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestSubmitPPMShipmentDocumentationHandlerUnit() {
	setUpPPMShipment := func() models.PPMShipment {

		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)

		ppmShipment.ID = uuid.Must(uuid.NewV4())
		ppmShipment.CreatedAt = time.Now()
		ppmShipment.UpdatedAt = ppmShipment.CreatedAt.AddDate(0, 0, 5)
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Must(uuid.NewV4())

		return ppmShipment
	}

	setUpRequestAndParams := func(
		ppmShipmentID uuid.UUID,
	) ppm.SubmitPPMShipmentDocumentationParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipmentID.String())
		request := httptest.NewRequest("POST", endpoint, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		request = suite.AuthenticateOfficeRequest(request, officeUser)

		return ppm.SubmitPPMShipmentDocumentationParams{
			HTTPRequest:   request,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipmentID),
		}
	}

	setUpPPMShipmentNewSubmitter := func(returnValue ...interface{}) services.PPMShipmentNewSubmitter {
		mockSubmitter := &mocks.PPMShipmentNewSubmitter{}

		mockSubmitter.On(
			"SubmitNewCustomerCloseOut",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockSubmitter
	}

	setUpHandler := func(submitter services.PPMShipmentNewSubmitter) SubmitPPMShipmentDocumentationHandler {
		return SubmitPPMShipmentDocumentationHandler{
			suite.HandlerConfig(),
			submitter,
		}
	}

	suite.Run("Returns an error if the PPMShipment ID in the url is invalid", func() {
		params := setUpRequestAndParams(uuid.Nil)

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppm.SubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppm.SubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Message, "nil PPM shipment ID")
		}
	})

	suite.Run("Returns an error if the request doesn't come from the office app", func() {
		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))
		ppmShipment := setUpPPMShipment()
		endpoint := fmt.Sprintf("/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipment.ID.String())
		request := httptest.NewRequest("POST", endpoint, nil)
		params := ppm.SubmitPPMShipmentDocumentationParams{
			HTTPRequest:   request,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
		}
		params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)
		response := handler.Handle(params)

		err := apperror.NewSessionError("Request should come from the office app.")

		if suite.IsType(&ppm.SubmitPPMShipmentDocumentationForbidden{}, response) {
			errResponse := response.(*ppm.SubmitPPMShipmentDocumentationForbidden)

			suite.Contains(*errResponse.Payload.Message, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a NotFoundError", func() {
		ppmShipment := setUpPPMShipment()

		params := setUpRequestAndParams(ppmShipment.ID)

		err := apperror.NewNotFoundError(ppmShipment.ID, "Can't find PPM shipment")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))
		response := handler.Handle(params)

		if suite.IsType(&ppm.SubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppm.SubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Message, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a QueryError", func() {
		ppmShipment := setUpPPMShipment()

		params := setUpRequestAndParams(ppmShipment.ID)

		err := apperror.NewQueryError("PPMShipment", nil, "Error getting PPM shipment")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))
		response := handler.Handle(params)

		suite.IsType(&ppm.SubmitPPMShipmentDocumentationInternalServerError{}, response)
	})

	suite.Run("Returns an error if the submitter service returns a ConflictError", func() {
		ppmShipment := setUpPPMShipment()

		params := setUpRequestAndParams(ppmShipment.ID)

		err := apperror.NewConflictError(ppmShipment.ID, "Can't route PPM shipment")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))
		response := handler.Handle(params)

		if suite.IsType(&ppm.SubmitPPMShipmentDocumentationConflict{}, response) {
			errResponse := response.(*ppm.SubmitPPMShipmentDocumentationConflict)

			suite.Contains(*errResponse.Payload.Message, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns an unexpected error", func() {
		ppmShipment := setUpPPMShipment()

		params := setUpRequestAndParams(ppmShipment.ID)

		err := apperror.NewNotImplementedError("Not implemented")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))
		response := handler.Handle(params)

		suite.IsType(&ppm.SubmitPPMShipmentDocumentationInternalServerError{}, response)
	})

	suite.Run("Returns the PPM shipment if all goes well", func() {
		ppmShipment := setUpPPMShipment()

		params := setUpRequestAndParams(ppmShipment.ID)

		expectedPPMShipment := ppmShipment
		expectedPPMShipment.Status = models.PPMShipmentStatusNeedsCloseout
		expectedPPMShipment.SubmittedAt = models.TimePointer(time.Now())
		signedCertification := models.SignedCertification{}

		expectedPPMShipment.SignedCertification = &signedCertification

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(&expectedPPMShipment, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppm.SubmitPPMShipmentDocumentationOK{}, response) {
			okResponse := response.(*ppm.SubmitPPMShipmentDocumentationOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(expectedPPMShipment.ID, returnedPPMShipment.ID)
		}
	})
}
