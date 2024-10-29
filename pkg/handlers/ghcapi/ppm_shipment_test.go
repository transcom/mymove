package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	ppmsitops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
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

func (suite *HandlerSuite) TestGetPPMSITEstimatedCostHandler() {
	var ppmShipment models.PPMShipment
	newFakeSITEstimatedCost := models.CentPointer(unit.Cents(25500))

	setupPricerData := func() {
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

		dopService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOP)

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

		ddpService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDDP)

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

		dpkService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDPK)

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

		dupkService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)

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

		dofsitService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

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

		doasitService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOASIT)

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

		ddfsitService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)

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

		ddasitService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDDASIT)

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

	suite.PreloadData(func() {
		setupPricerData()
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
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"90210", "30813").Return(2294, nil)
	})

	setUpGetCostRequestAndParams := func() ppmsitops.GetPPMSITEstimatedCostParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/sit_location/%s/sit-estimated-cost", ppmShipment.ID.String(), *ppmShipment.SITLocation)

		req := httptest.NewRequest("GET", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		return ppmsitops.GetPPMSITEstimatedCostParams{
			HTTPRequest:   req,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
			SitLocation:   string(*ppmShipment.SITLocation),
		}
	}

	setUpUpdateCostRequestAndParams := func(sitLocation *ghcmessages.SITLocationType) ppmsitops.UpdatePPMSITParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/ppm-sit", ppmShipment.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)
		eTag := etag.GenerateEtag(ppmShipment.UpdatedAt)

		return ppmsitops.UpdatePPMSITParams{
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

		if suite.IsType(&ppmsitops.GetPPMSITEstimatedCostOK{}, response) {
			payload := response.(*ppmsitops.GetPPMSITEstimatedCostOK).Payload

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

		suite.IsType(&ppmsitops.GetPPMSITEstimatedCostNotFound{}, response)
	})

	suite.Run("Update PPM SIT Estimated Cost - ORIGIN", func() {
		sitLocation := ghcmessages.SITLocationTypeORIGIN
		subtestData := setUpForUpdateCostTests(nil)

		handler := setUpUpdateCostHandler(subtestData.ppmShipmentUpdater, subtestData.ppmShipmentFetcher)
		params := setUpUpdateCostRequestAndParams(&sitLocation)
		response := handler.Handle(params)

		if suite.IsType(&ppmsitops.UpdatePPMSITOK{}, response) {
			payload := response.(*ppmsitops.UpdatePPMSITOK).Payload

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

		suite.IsType(&ppmsitops.UpdatePPMSITNotFound{}, response)
	})
}
