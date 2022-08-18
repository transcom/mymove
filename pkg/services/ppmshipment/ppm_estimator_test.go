package ppmshipment

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	prhelpermocks "github.com/transcom/mymove/pkg/payment_request/mocks"

	"github.com/transcom/mymove/pkg/route/mocks"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestPPMEstimator() {

	mockedPlanner := &mocks.Planner{}
	mockedPaymentRequestHelper := &prhelpermocks.Helper{}
	ppmEstimator := NewEstimatePPM(mockedPlanner, mockedPaymentRequestHelper)

	// To avoid creating all of the re_services and their corresponding params using factories, we can create this
	// mapping to help mock the response
	serviceParamKeys := map[models.ServiceItemParamName]models.ServiceItemParamKey{
		models.ServiceItemParamNameActualPickupDate:                 {Key: models.ServiceItemParamNameActualPickupDate, Type: models.ServiceItemParamTypeDate},
		models.ServiceItemParamNameContractCode:                     {Key: models.ServiceItemParamNameContractCode, Type: models.ServiceItemParamTypeString},
		models.ServiceItemParamNameDistanceZip:                      {Key: models.ServiceItemParamNameDistanceZip, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameEIAFuelPrice:                     {Key: models.ServiceItemParamNameEIAFuelPrice, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier: {Key: models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier, Type: models.ServiceItemParamTypeDecimal},
		models.ServiceItemParamNameReferenceDate:                    {Key: models.ServiceItemParamNameReferenceDate, Type: models.ServiceItemParamTypeDate},
		models.ServiceItemParamNameRequestedPickupDate:              {Key: models.ServiceItemParamNameRequestedPickupDate, Type: models.ServiceItemParamTypeDate},
		models.ServiceItemParamNameServiceAreaDest:                  {Key: models.ServiceItemParamNameServiceAreaDest, Type: models.ServiceItemParamTypeString},
		models.ServiceItemParamNameServiceAreaOrigin:                {Key: models.ServiceItemParamNameServiceAreaOrigin, Type: models.ServiceItemParamTypeString},
		models.ServiceItemParamNameServicesScheduleDest:             {Key: models.ServiceItemParamNameServicesScheduleDest, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameServicesScheduleOrigin:           {Key: models.ServiceItemParamNameServicesScheduleOrigin, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameWeightAdjusted:                   {Key: models.ServiceItemParamNameWeightAdjusted, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameWeightBilled:                     {Key: models.ServiceItemParamNameWeightBilled, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameWeightEstimated:                  {Key: models.ServiceItemParamNameWeightEstimated, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameWeightOriginal:                   {Key: models.ServiceItemParamNameWeightOriginal, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameWeightReweigh:                    {Key: models.ServiceItemParamNameWeightReweigh, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameZipDestAddress:                   {Key: models.ServiceItemParamNameZipDestAddress, Type: models.ServiceItemParamTypeString},
		models.ServiceItemParamNameZipPickupAddress:                 {Key: models.ServiceItemParamNameZipPickupAddress, Type: models.ServiceItemParamTypeString},
	}

	serviceParams := models.ServiceParams{}
	// Domestic Linehaul
	for _, serviceParamKey := range []models.ServiceItemParamName{
		models.ServiceItemParamNameActualPickupDate,
		models.ServiceItemParamNameContractCode,
		models.ServiceItemParamNameDistanceZip,
		models.ServiceItemParamNameReferenceDate,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameServiceAreaOrigin,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameWeightBilled,
		models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightOriginal,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameZipDestAddress,
		models.ServiceItemParamNameZipPickupAddress,
	} {
		serviceParams = append(serviceParams, models.ServiceParam{Service: models.ReService{Code: models.ReServiceCodeDLH}, ServiceItemParamKey: serviceParamKeys[serviceParamKey]})
	}

	// Fuel Surcharge
	for _, serviceParamKey := range []models.ServiceItemParamName{
		models.ServiceItemParamNameActualPickupDate,
		models.ServiceItemParamNameContractCode,
		models.ServiceItemParamNameDistanceZip,
		models.ServiceItemParamNameEIAFuelPrice,
		models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameWeightBilled,
		models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightOriginal,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameZipDestAddress,
		models.ServiceItemParamNameZipPickupAddress,
	} {
		serviceParams = append(serviceParams, models.ServiceParam{Service: models.ReService{Code: models.ReServiceCodeFSC}, ServiceItemParamKey: serviceParamKeys[serviceParamKey]})
	}

	// Domestic Origin Price
	for _, serviceParamKey := range []models.ServiceItemParamName{
		models.ServiceItemParamNameActualPickupDate,
		models.ServiceItemParamNameContractCode,
		models.ServiceItemParamNameReferenceDate,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameServiceAreaOrigin,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameWeightBilled,
		models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightOriginal,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameZipPickupAddress,
	} {
		serviceParams = append(serviceParams, models.ServiceParam{Service: models.ReService{Code: models.ReServiceCodeDOP}, ServiceItemParamKey: serviceParamKeys[serviceParamKey]})
	}

	// Domestic Destination Price
	for _, serviceParamKey := range []models.ServiceItemParamName{
		models.ServiceItemParamNameActualPickupDate,
		models.ServiceItemParamNameContractCode,
		models.ServiceItemParamNameReferenceDate,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameServiceAreaDest,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameWeightBilled,
		models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightOriginal,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameZipDestAddress,
	} {
		serviceParams = append(serviceParams, models.ServiceParam{Service: models.ReService{Code: models.ReServiceCodeDDP}, ServiceItemParamKey: serviceParamKeys[serviceParamKey]})
	}

	// Domestic Packing
	for _, serviceParamKey := range []models.ServiceItemParamName{
		models.ServiceItemParamNameActualPickupDate,
		models.ServiceItemParamNameContractCode,
		models.ServiceItemParamNameReferenceDate,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameServiceAreaOrigin,
		models.ServiceItemParamNameServicesScheduleOrigin,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameWeightBilled,
		models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightOriginal,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameZipPickupAddress,
	} {
		serviceParams = append(serviceParams, models.ServiceParam{Service: models.ReService{Code: models.ReServiceCodeDPK}, ServiceItemParamKey: serviceParamKeys[serviceParamKey]})
	}

	// Domestic Unpacking
	for _, serviceParamKey := range []models.ServiceItemParamName{
		models.ServiceItemParamNameActualPickupDate,
		models.ServiceItemParamNameContractCode,
		models.ServiceItemParamNameReferenceDate,
		models.ServiceItemParamNameRequestedPickupDate,
		models.ServiceItemParamNameServiceAreaDest,
		models.ServiceItemParamNameServicesScheduleDest,
		models.ServiceItemParamNameWeightAdjusted,
		models.ServiceItemParamNameWeightBilled,
		models.ServiceItemParamNameWeightEstimated,
		models.ServiceItemParamNameWeightOriginal,
		models.ServiceItemParamNameWeightReweigh,
		models.ServiceItemParamNameZipDestAddress,
	} {
		serviceParams = append(serviceParams, models.ServiceParam{Service: models.ReService{Code: models.ReServiceCodeDUPK}, ServiceItemParamKey: serviceParamKeys[serviceParamKey]})
	}

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

		dopService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOP,
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

		ddpService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDP,
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

		dpkService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDPK,
			},
		})

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

		dupkService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDUPK,
			},
		})

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

		dofsitService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
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

		doasitService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
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

		ddfsitService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
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

		ddasitService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDASIT,
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

	suite.Run("Estimated Incentive", func() {
		suite.Run("Estimated Incentive - Success", func() {
			oldPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

			setupPricerData()

			// shipment has locations and date but is now updating the estimated weight for the first time
			estimatedWeight := unit.Pound(5000)
			newPPM := oldPPMShipment
			newPPM.EstimatedWeight = &estimatedWeight

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30813").Return(2294, nil)

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
			suite.Equal(unit.Pound(5000), *newPPM.EstimatedWeight)
			suite.Equal(unit.Cents(70064364), *ppmEstimate)
		})

		suite.Run("Estimated Incentive - Success - clears advance and advance requested values", func() {
			oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					Status: models.PPMShipmentStatusDraft,
				},
			})

			setupPricerData()

			newPPM := oldPPMShipment

			// updating the departure date will re-calculate the estimate and clear the previously requested advance
			newPPM.ExpectedDepartureDate = time.Date(testdatagen.GHCTestYear, time.March, 30, 0, 0, 0, 0, time.UTC)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30813").Return(2294, nil).Once()

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Nil(newPPM.HasRequestedAdvance)
			suite.Nil(newPPM.AdvanceAmountRequested)
			suite.Equal(unit.Cents(38213948), *ppmEstimate)
		})

		suite.Run("Estimated Incentive - does not change when required fields are the same", func() {
			oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					Status:             models.PPMShipmentStatusDraft,
					EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
				},
			})

			newPPM := oldPPMShipment
			newPPM.HasProGear = models.BoolPointer(false)

			estimatedIncentive, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
			suite.Equal(*oldPPMShipment.EstimatedWeight, *newPPM.EstimatedWeight)
			suite.Equal(oldPPMShipment.DestinationPostalCode, newPPM.DestinationPostalCode)
			suite.True(oldPPMShipment.ExpectedDepartureDate.Equal(newPPM.ExpectedDepartureDate))
			suite.Equal(*oldPPMShipment.EstimatedIncentive, *estimatedIncentive)
			suite.Equal(models.BoolPointer(true), newPPM.HasRequestedAdvance)
			suite.Equal(unit.Cents(598700), *newPPM.AdvanceAmountRequested)
		})

		suite.Run("Estimated Incentive - does not change when status is not DRAFT", func() {
			oldPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
				},
			})

			newPPM := models.PPMShipment{
				ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
				ShipmentID:            oldPPMShipment.ShipmentID,
				Status:                models.PPMShipmentStatusPaymentApproved,
				ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
				PickupPostalCode:      oldPPMShipment.PickupPostalCode,
				DestinationPostalCode: "94040",
				EstimatedWeight:       oldPPMShipment.EstimatedWeight,
				SITExpected:           oldPPMShipment.SITExpected,
				EstimatedIncentive:    models.CentPointer(unit.Cents(600000)),
			}

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Equal(oldPPMShipment.EstimatedIncentive, ppmEstimate)
		})

		suite.Run("Estimated Incentive - Success - is skipped when Estimated Weight is missing", func() {
			oldPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

			newPPM := oldPPMShipment
			newPPM.DestinationPostalCode = "94040"

			_, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NoError(err)
			suite.Nil(newPPM.EstimatedIncentive)
		})
	})

	suite.Run("SIT Estimated Cost", func() {
		// For comparison should be priced the same as ORGSIT in devseed
		suite.Run("Success - Origin First Day and Additional Day SIT", func() {
			setupPricerData()

			originLocation := models.SITLocationTypeOrigin
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			})
			shipmentOriginSIT := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					Shipment:                  mtoShipment,
					ShipmentID:                mtoShipment.ID,
					DestinationPostalCode:     "30813",
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &originLocation,
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
				},
				Stub: true,
			})

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30813").Return(2294, nil)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), models.PPMShipment{}, &shipmentOriginSIT)

			suite.NoError(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(50660, estimatedSITCost.Int())
		})

		suite.Run("Success - Destination First Day and Additional Day SIT", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			})
			shipmentOriginSIT := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					Shipment:                  mtoShipment,
					ShipmentID:                mtoShipment.ID,
					DestinationPostalCode:     "30813",
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationLocation,
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
				},
				Stub: true,
			})

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30813").Return(2294, nil)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), models.PPMShipment{}, &shipmentOriginSIT)

			suite.NoError(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(65240, estimatedSITCost.Int())
		})

		suite.Run("Success - same entry and departure dates only prices first day SIT", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			})
			shipmentOriginSIT := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					Shipment:                  mtoShipment,
					ShipmentID:                mtoShipment.ID,
					DestinationPostalCode:     "30813",
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationLocation,
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
					SITEstimatedEntryDate:     &entryDate,
					SITEstimatedDepartureDate: &entryDate,
				},
				Stub: true,
			})

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30813").Return(2294, nil)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), models.PPMShipment{}, &shipmentOriginSIT)

			suite.NoError(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(32240, estimatedSITCost.Int())
		})

		suite.Run("SIT cost is not calculated when required fields are missing", func() {
			setupPricerData()

			destinationSITLocation := models.SITLocationTypeDestination

			// an MTO Shipment ID is required for the shipment query
			shipmentSITFieldsNotUpdated := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{})
			shipmentSITNotExpected := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: shipmentSITFieldsNotUpdated.Shipment,
				Stub:        true,
			})
			shipmentSITWeightMissing := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: shipmentSITFieldsNotUpdated.Shipment,
				PPMShipment: models.PPMShipment{
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationSITLocation,
					SITEstimatedEntryDate:     models.TimePointer(time.Now()),
					SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
				},
				Stub: true,
			})
			shipmentSITEntryDateMissing := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: shipmentSITFieldsNotUpdated.Shipment,
				PPMShipment: models.PPMShipment{
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationSITLocation,
					SITEstimatedDepartureDate: models.TimePointer(time.Now()),
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
				},
				Stub: true,
			})
			shipmentSITDepartureDateMissing := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				MTOShipment: shipmentSITFieldsNotUpdated.Shipment,
				PPMShipment: models.PPMShipment{
					SITExpected:           models.BoolPointer(true),
					SITLocation:           &destinationSITLocation,
					SITEstimatedEntryDate: models.TimePointer(time.Now()),
					SITEstimatedWeight:    models.PoundPointer(unit.Pound(2999)),
				},
				Stub: true,
			})

			shipmentTestCases := []struct {
				oldShipment models.PPMShipment
				newShipment models.PPMShipment
				name        string
			}{
				{
					models.PPMShipment{},
					shipmentSITNotExpected,
					"PPM Shipment with SITExpected set to false",
				},
				{
					models.PPMShipment{},
					shipmentSITWeightMissing,
					"PPM Shipment with SIT Estimated Weight missing",
				},
				{
					models.PPMShipment{},
					shipmentSITEntryDateMissing,
					"PPM Shipment with SIT Entry Date missing",
				},
				{
					models.PPMShipment{},
					shipmentSITDepartureDateMissing,
					"PPM Shipment with SIT Departure Date missing",
				},
				{
					models.PPMShipment{},
					shipmentSITDepartureDateMissing,
					"PPM Shipment with SIT Departure Date missing",
				},
				{
					shipmentSITFieldsNotUpdated,
					shipmentSITFieldsNotUpdated,
					"PPM Shipment fields were not updated",
				},
			}

			for _, testCase := range shipmentTestCases {
				_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), testCase.oldShipment, &testCase.newShipment)
				suite.NoError(err, fmt.Sprintf("unexpected error running test %q", testCase.name))
				suite.Nil(estimatedSITCost, fmt.Sprintf("SIT cost was calculated when it shouldnt't have been during test %q", testCase.name))
			}
		})

		suite.Run("SIT cost is not re-calculated when fields are unchanged", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			shipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationLocation,
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
					SITEstimatedEntryDate:     models.TimePointer(time.Now()),
					SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
					SITEstimatedCost:          models.CentPointer(unit.Cents(89900)),
				},
			})
			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), shipment, &shipment)
			suite.NoError(err)
			suite.Equal(*shipment.SITEstimatedCost, *estimatedSITCost)
		})

		suite.Run("SIT cost is re-calculated when any dependent field is changed", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			originalShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationLocation,
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
					SITEstimatedEntryDate:     models.TimePointer(time.Now()),
					SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
					SITEstimatedCost:          models.CentPointer(unit.Cents(89900)),
				},
			})

			// PPM base shipment field changes will affect SIT pricing
			shipmentDifferentPickup := originalShipment
			shipmentDifferentPickup.PickupPostalCode = "90211"

			shipmentDifferentDestination := originalShipment
			shipmentDifferentDestination.DestinationPostalCode = "30814"

			shipmentDifferentDeparture := originalShipment
			// original date was Mar 15th so adding 3 months should affect the date peak period pricing
			shipmentDifferentDeparture.ExpectedDepartureDate = originalShipment.ExpectedDepartureDate.Add(time.Hour * 24 * 70)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90211", "30813").Return(2294, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90210", "30814").Return(2290, nil)

			// SIT specific field changes will likely cause the price to change, although adjusting dates may not change
			// the total number of days in SIT.

			shipmentDifferentLocation := originalShipment
			originLocation := models.SITLocationTypeOrigin
			shipmentDifferentLocation.SITLocation = &originLocation

			shipmentDifferentSITWeight := originalShipment
			shipmentDifferentSITWeight.SITEstimatedWeight = models.PoundPointer(unit.Pound(4555))

			shipmentDifferentEntryDate := originalShipment
			previousDay := originalShipment.SITEstimatedEntryDate.Add(time.Hour * -24)
			shipmentDifferentEntryDate.SITEstimatedEntryDate = &previousDay

			shipmentDifferentSITDepartureDate := originalShipment
			nextDay := shipmentDifferentSITDepartureDate.SITEstimatedDepartureDate.Add(time.Hour * 24)
			shipmentDifferentSITDepartureDate.SITEstimatedDepartureDate = &nextDay

			for _, updatedShipment := range []models.PPMShipment{
				shipmentDifferentPickup,
				shipmentDifferentDestination,
				shipmentDifferentDeparture,
				shipmentDifferentLocation,
				shipmentDifferentSITWeight,
				shipmentDifferentEntryDate,
				shipmentDifferentSITDepartureDate,
			} {
				copyOfShipment := updatedShipment

				_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), originalShipment, &copyOfShipment)

				suite.NoError(err)
				suite.NotNil(originalShipment.SITEstimatedCost)
				suite.NotNil(estimatedSITCost)
				suite.NotEqual(*originalShipment.SITEstimatedCost, *estimatedSITCost)
			}
		})

		suite.Run("SIT cost is set to nil when storage is no longer expected", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			originalShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					SITExpected:               models.BoolPointer(true),
					SITLocation:               &destinationLocation,
					SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
					SITEstimatedEntryDate:     models.TimePointer(time.Now()),
					SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
					SITEstimatedCost:          models.CentPointer(unit.Cents(89900)),
				},
			})

			shipmentSITNotExpected := originalShipment
			shipmentSITNotExpected.SITExpected = models.BoolPointer(false)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), originalShipment, &shipmentSITNotExpected)
			suite.NoError(err)
			suite.Nil(shipmentSITNotExpected.SITEstimatedCost)
			suite.Nil(estimatedSITCost)
		})
	})
}
