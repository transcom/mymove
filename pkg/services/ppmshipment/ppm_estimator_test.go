package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	prhelpermocks "github.com/transcom/mymove/pkg/payment_request/mocks"
	"github.com/transcom/mymove/pkg/route/mocks"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestPPMEstimator() {
	mockedPlanner := &mocks.Planner{}
	mockedPaymentRequestHelper := &prhelpermocks.Helper{}
	ppmEstimator := NewEstimatePPM(mockedPlanner, mockedPaymentRequestHelper)
	validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
	invalidGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-05-14")

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
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2020, time.March, 16, 0, 0, 0, 0, time.UTC),
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
				Zip3:                "503",
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

	suite.Run("Price Breakdown - Incentive-based PPM", func() {
		ppmShipment := factory.BuildPPMShipmentWithApprovedDocuments(suite.DB())

		setupPricerData()

		mockedPaymentRequestHelper.On(
			"FetchServiceParamsForServiceItems",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

		// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"50309", "30813").Return(2294, nil)

		linehaul, fuel, origin, dest, packing, unpacking, _, err := ppmEstimator.PriceBreakdown(suite.AppContextForTest(), &ppmShipment)
		suite.NilOrNoVerrs(err)

		mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"50309", "30813")
		mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

		suite.Equal(unit.Pound(4000), *ppmShipment.EstimatedWeight)
		suite.Equal(unit.Cents(48155648), linehaul)
		suite.Equal(unit.Cents(-641), fuel)
		suite.Equal(unit.Cents(24160), origin)
		suite.Equal(unit.Cents(36960), dest)
		suite.Equal(unit.Cents(321920), packing)
		suite.Equal(unit.Cents(26520), unpacking)

		total := linehaul + fuel + origin + dest + packing + unpacking
		suite.Equal(unit.Cents(48564567), total)

		// testing multiplier functionality when multiplier is not nil
		gccMultiplier := models.GCCMultiplier{
			Multiplier: 1.3,
		}
		ppmShipment.GCCMultiplier = &gccMultiplier

		linehaul, fuel, origin, dest, packing, unpacking, _, err = ppmEstimator.PriceBreakdown(suite.AppContextForTest(), &ppmShipment)
		suite.NilOrNoVerrs(err)

		suite.Equal(unit.Pound(4000), *ppmShipment.EstimatedWeight)
		suite.Equal(unit.Cents(62602342), linehaul)
		suite.Equal(unit.Cents(-641), fuel)
		suite.Equal(unit.Cents(31408), origin)
		suite.Equal(unit.Cents(48048), dest)
		suite.Equal(unit.Cents(418496), packing)
		suite.Equal(unit.Cents(34476), unpacking)

		totalWithMultiplier := linehaul + fuel + origin + dest + packing + unpacking
		suite.Equal(unit.Cents(63134129), totalWithMultiplier)
	})

	suite.Run("Price Breakdown - Small package PPM", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)
		suite.FatalNoError(uploaderErr)

		// this factory has two moving expenses that total 4000 pounds
		// pricing should be the same as the above test
		ppmShipment := factory.BuildPPMSPRShipmentWithoutPaymentPacketTwoExpenses(suite.DB(), userUploader)

		setupPricerData()

		mockedPaymentRequestHelper.On(
			"FetchServiceParamsForServiceItems",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

		// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"50309", "30813").Return(2294, nil)

		linehaul, fuel, origin, dest, packing, unpacking, _, err := ppmEstimator.PriceBreakdown(suite.AppContextForTest(), &ppmShipment)
		suite.NilOrNoVerrs(err)

		mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"50309", "30813")
		mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

		suite.Equal(unit.Pound(4000), *ppmShipment.EstimatedWeight)
		suite.Equal(unit.Cents(48155648), linehaul)
		suite.Equal(unit.Cents(-641), fuel)
		suite.Equal(unit.Cents(24160), origin)
		suite.Equal(unit.Cents(36960), dest)
		suite.Equal(unit.Cents(321920), packing)
		suite.Equal(unit.Cents(26520), unpacking)

		total := linehaul + fuel + origin + dest + packing + unpacking
		suite.Equal(unit.Cents(48564567), total)
	})

	suite.Run("Estimated Incentive", func() {
		suite.Run("Estimated Incentive - Success using estimated weight and not db authorized weight", func() {
			// when the PPM shipment is in draft, we use the estimated weight and not the db authorized weight
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status: models.PPMShipmentStatusDraft,
					},
				},
			}, nil)
			setupPricerData()

			// shipment has locations and date but is now updating the estimated weight for the first time
			estimatedWeight := unit.Pound(5000)
			newPPM := oldPPMShipment
			newPPM.EstimatedWeight = &estimatedWeight

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil).Twice()
			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(oldPPMShipment.PickupAddress.PostalCode, newPPM.PickupAddress.PostalCode)
			suite.Equal(unit.Pound(5000), *newPPM.EstimatedWeight)
			suite.Equal(unit.Cents(89071179), *ppmEstimate)

			// appending this to test functionality of the GCC multiplier
			ppmWithMultiplier := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ExpectedDepartureDate: validGccMultiplierDate,
					},
				},
			}, nil)
			newPPMWithMultiplier := ppmWithMultiplier
			newPPMWithMultiplier.EstimatedWeight = &estimatedWeight // setting weight to 5000
			ppmEstimatedWithMultiplier, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), ppmWithMultiplier, &newPPMWithMultiplier)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(unit.Pound(5000), *newPPMWithMultiplier.EstimatedWeight)
			suite.NotEqual(unit.Cents(89071179), *ppmEstimatedWithMultiplier)
			suite.Equal(unit.Cents(120722169), *ppmEstimatedWithMultiplier)
		})

		suite.Run("Estimated Incentive - Success using db authorize weight and not estimated incentive", func() {
			// when the PPM shipment is NOT in draft, we use the db authorized weight and not the estimated weight
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status: models.PPMShipmentStatusNeedsCloseout,
					},
				},
			}, nil)
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
				"50309", "30813").Return(2294, nil)

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(oldPPMShipment.PickupAddress.PostalCode, newPPM.PickupAddress.PostalCode)
			suite.Equal(unit.Pound(5000), *newPPM.EstimatedWeight)
			suite.Equal(unit.Cents(1000000), *ppmEstimate)
		})

		suite.Run("Estimated Incentive - Success when old Estimated Incentive is zero", func() {
			oldPPMShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

			zeroIncentive := unit.Cents(0)
			oldPPMShipment.EstimatedIncentive = &zeroIncentive

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
				"50309", "30813").Return(2294, nil)

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(oldPPMShipment.PickupAddress.PostalCode, newPPM.PickupAddress.PostalCode)
			suite.Equal(unit.Pound(5000), *newPPM.EstimatedWeight)
			suite.Equal(unit.Cents(89071179), *ppmEstimate)
		})

		suite.Run("Estimated Incentive - Success when old Estimated Incentive is zero", func() {
			oldPPMShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

			zeroIncentive := unit.Cents(0)
			oldPPMShipment.EstimatedIncentive = &zeroIncentive

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
				"50309", "30813").Return(2294, nil)

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(oldPPMShipment.PickupAddress.PostalCode, newPPM.PickupAddress.PostalCode)
			suite.Equal(unit.Pound(5000), *newPPM.EstimatedWeight)
			suite.Equal(unit.Cents(89071179), *ppmEstimate)
		})

		suite.Run("Estimated Incentive - Success - clears advance and advance requested values", func() {
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status: models.PPMShipmentStatusDraft,
					},
				},
			}, nil)
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
				"50309", "30813").Return(2294, nil).Once()

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Nil(newPPM.HasRequestedAdvance)
			suite.Nil(newPPM.AdvanceAmountRequested)
			suite.Equal(unit.Cents(48564567), *ppmEstimate)
		})

		suite.Run("Estimated Incentive - does not change when required fields are the same", func() {
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status:             models.PPMShipmentStatusDraft,
						EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
			}, nil)
			setupPricerData()
			newPPM := oldPPMShipment
			newPPM.HasProGear = models.BoolPointer(false)

			estimatedIncentive, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Equal(oldPPMShipment.PickupAddress.PostalCode, newPPM.PickupAddress.PostalCode)
			suite.Equal(*oldPPMShipment.EstimatedWeight, *newPPM.EstimatedWeight)
			suite.Equal(oldPPMShipment.DestinationAddress.PostalCode, newPPM.DestinationAddress.PostalCode)
			suite.True(oldPPMShipment.ExpectedDepartureDate.Equal(newPPM.ExpectedDepartureDate))
			suite.Equal(*oldPPMShipment.EstimatedIncentive, *estimatedIncentive)
			suite.Equal(models.BoolPointer(true), newPPM.HasRequestedAdvance)
			suite.Equal(unit.Cents(598700), *newPPM.AdvanceAmountRequested)
		})

		suite.Run("Estimated Incentive - does not change when status is not DRAFT", func() {
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						EstimatedIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
			}, nil)

			pickupAddress := models.Address{PostalCode: oldPPMShipment.PickupAddress.PostalCode}
			destinationAddress := models.Address{PostalCode: "94040"}
			newPPM := models.PPMShipment{
				ID:                    uuid.FromStringOrNil("575c25aa-b4eb-4024-9597-43483003c773"),
				ShipmentID:            oldPPMShipment.ShipmentID,
				Status:                models.PPMShipmentStatusCloseoutComplete,
				ExpectedDepartureDate: oldPPMShipment.ExpectedDepartureDate,
				PickupAddress:         &pickupAddress,
				DestinationAddress:    &destinationAddress,
				EstimatedWeight:       oldPPMShipment.EstimatedWeight,
				SITExpected:           oldPPMShipment.SITExpected,
				EstimatedIncentive:    models.CentPointer(unit.Cents(600000)),
			}

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Equal(oldPPMShipment.EstimatedIncentive, ppmEstimate)
		})

		suite.Run("Estimated Incentive - Success - is skipped when Estimated Weight is missing", func() {
			oldPPMShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

			newPPM := oldPPMShipment
			newPPM.DestinationAddress.PostalCode = "94040"
			_, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NoError(err)
			suite.Nil(newPPM.EstimatedIncentive)
		})
	})

	suite.Run("Max Incentive", func() {
		suite.Run("Max Incentive - Success", func() {
			maxIncentive2 := unit.Cents(900000000)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						MaxIncentive: &maxIncentive2,
					},
				},
			}, nil)
			setupPricerData()

			estimatedWeight := unit.Pound(5000)
			newPPM := oldPPMShipment
			newPPM.EstimatedWeight = &estimatedWeight
			newDate, _ := time.Parse("2006-01-02", "2025-05-02")
			newPPM.ExpectedDepartureDate = newDate

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil).Twice()

			maxIncentive, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.Equal(unit.Cents(142513951), *maxIncentive)

			// appending this to test functionality of the GCC multiplier
			ppmWithMultiplier := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ExpectedDepartureDate: validGccMultiplierDate,
						MaxIncentive:          &maxIncentive2,
					},
				},
			}, nil)
			newPPMWithMultiplier := ppmWithMultiplier
			newPPMWithMultiplier.EstimatedWeight = &estimatedWeight
			newDateInMultiplier, _ := time.Parse("2006-01-02", "2025-05-30")
			newPPMWithMultiplier.ExpectedDepartureDate = newDateInMultiplier
			ppmMaxWithMultiplier, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), ppmWithMultiplier, &newPPMWithMultiplier)
			suite.NilOrNoVerrs(err)

			suite.NotEqual(unit.Cents(142513951), *ppmMaxWithMultiplier)
			suite.Equal(unit.Cents(193155535), *ppmMaxWithMultiplier)
		})

		suite.Run("Max Incentive - Success - is skipped when Estimated Weight is missing", func() {
			oldPPMShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

			newPPM := oldPPMShipment
			newPPM.DestinationAddress.PostalCode = "94040"
			_, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NoError(err)
			suite.Nil(newPPM.MaxIncentive)
		})

		suite.Run("Max Incentive - Skips recalculation when GCC multiplier and departure date unchanged", func() {
			existingMaxIncentive := unit.Cents(123456789)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ExpectedDepartureDate: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
						MaxIncentive:          &existingMaxIncentive,
					},
				},
			}, nil)

			newPPM := oldPPMShipment
			estimatedWeight := unit.Pound(5000)
			newPPM.EstimatedWeight = &estimatedWeight

			// make sure departure date and GCC multiplier are the same
			newPPM.ExpectedDepartureDate = oldPPMShipment.ExpectedDepartureDate
			newPPM.GCCMultiplierID = oldPPMShipment.GCCMultiplierID

			maxIncentive, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NoError(err)
			suite.Equal(*oldPPMShipment.MaxIncentive, *maxIncentive)
		})

		suite.Run("Max Incentive - recalculates when old multiplier is nil and new is valid", func() {
			validMultiplierID := uuid.Must(uuid.NewV4())
			maxIncentive2 := unit.Cents(123456789)

			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						GCCMultiplierID:       nil,
						MaxIncentive:          &maxIncentive2,
						ExpectedDepartureDate: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
					},
				},
			}, nil)

			newPPM := oldPPMShipment
			estimatedWeight := unit.Pound(5000)
			newPPM.EstimatedWeight = &estimatedWeight
			newPPM.GCCMultiplierID = &validMultiplierID // introduce a new multiplier

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil).Once()

			maxIncentive, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NoError(err)

			suite.NotEqual(*oldPPMShipment.MaxIncentive, *maxIncentive)
		})

		suite.Run("Max Incentive - recalculates when actual move date is being updated", func() {
			maxIncentive2 := unit.Cents(123456789)

			// expected departure date is set to a multiplier of 1.3x
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						MaxIncentive:          &maxIncentive2,
						ExpectedDepartureDate: validGccMultiplierDate,
					},
				},
			}, nil)

			// now we are updating the actual move date
			newPPM := oldPPMShipment
			// simulating an actual move date being updated with GCC multiplier of 1x
			newPPM.ActualMoveDate = &invalidGccMultiplierDate
			newPPM.GCCMultiplier = &models.GCCMultiplier{Multiplier: 1}

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil).Once()

			maxIncentive, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NoError(err)

			suite.NotEqual(*oldPPMShipment.MaxIncentive, *maxIncentive)
		})
	})

	suite.Run("Final Incentive", func() {
		actualMoveDate := time.Date(2020, time.March, 14, 0, 0, 0, 0, time.UTC)

		suite.Run("Final Incentive - Success", func() {
			setupPricerData()
			weightOverride := unit.Pound(19500)
			maxIncentive := unit.Cents(90000000)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(actualMoveDate),
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						MaxIncentive:   &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPM := oldPPMShipment
			updatedMoveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			newPPM.ActualMoveDate = models.TimePointer(updatedMoveDate)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil).Twice()

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(5000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)
			suite.Equal(unit.Cents(80249474), *ppmFinal)

			// appending this to test functionality of the GCC multiplier
			maxIncentive2 := unit.Cents(900000000)
			ppmWithMultiplier := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ExpectedDepartureDate: validGccMultiplierDate,
						Status:                models.PPMShipmentStatusWaitingOnCustomer,
						MaxIncentive:          &maxIncentive2,
					},
				},
			}, nil)
			ppmWithMultiplier.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPMWithMultiplier := ppmWithMultiplier
			newPPMWithMultiplier.ActualMoveDate = models.TimePointer(updatedMoveDate)
			ppmFinalWithMultiplier, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), ppmWithMultiplier, &newPPMWithMultiplier)
			suite.NilOrNoVerrs(err)

			suite.NotEqual(unit.Cents(80249474), *ppmFinalWithMultiplier)
			suite.Equal(unit.Cents(104324316), *ppmFinalWithMultiplier)
		})

		suite.Run("Final Incentive - Success when capped at max gcc", func() {
			setupPricerData()
			weightOverride := unit.Pound(19500)
			maxIncentive := unit.Cents(500)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(actualMoveDate),
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						MaxIncentive:   &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPM := oldPPMShipment
			updatedMoveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			newPPM.ActualMoveDate = models.TimePointer(updatedMoveDate)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(5000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)
			suite.Equal(maxIncentive, *ppmFinal)
		})

		suite.Run("Final Incentive - Success with allowable weight less than net weight", func() {
			setupPricerData()
			weightOverride := unit.Pound(19500)
			maxIncentive := unit.Cents(90000000)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(actualMoveDate),
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						MaxIncentive:   &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPM := oldPPMShipment
			updatedMoveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			newPPM.ActualMoveDate = models.TimePointer(updatedMoveDate)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(5000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)
			suite.Equal(unit.Cents(80249474), *ppmFinal)

			// Repeat the above shipment with an allowable weight less than the net weight
			weightOverride = unit.Pound(19500)
			allowableWeightOverride := unit.Pound(4000)
			oldPPMShipment = factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate:  models.TimePointer(actualMoveDate),
						Status:          models.PPMShipmentStatusWaitingOnCustomer,
						AllowableWeight: &allowableWeightOverride,
						MaxIncentive:    &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPM = oldPPMShipment
			updatedMoveDate = time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			newPPM.ActualMoveDate = models.TimePointer(updatedMoveDate)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinalIncentiveLimitedByAllowableWeight, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight = SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(5000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)

			// Confirm the incentive is less than if all of the weight was allowable
			suite.Less(*ppmFinalIncentiveLimitedByAllowableWeight, *ppmFinal)
		})

		suite.Run("Final Incentive - Success with capping weight at total entitlement", func() {
			// The first half of this test tests the entitlement cap. The second half uses the allowable to check the entitlement cap.
			// The max entitlement for this test data is 8000 lbs.
			setupPricerData()
			weightOverride := unit.Pound(24500)
			maxIncentive := unit.Cents(900000000)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(actualMoveDate),
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						MaxIncentive:   &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPM := oldPPMShipment
			updatedMoveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			newPPM.ActualMoveDate = models.TimePointer(updatedMoveDate)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			// These two weights are passed in to final incentive calculation, but their values aren't returned. Therefore,
			// we're checking to make sure the value is 10000, same as before they've been sent. The ppmFinal will be calculated with
			// the allowable 8000 lbs. The second half of this test ensures that the calculation value is correct for 8k lbs.
			suite.Equal(unit.Pound(10000), originalWeight)
			suite.Equal(unit.Pound(10000), newWeight)
			suite.Equal(unit.Cents(128398858), *ppmFinal)

			// Repeat the above shipment with an allowable weight equal to the entitlement. Since the allowable is covered
			// by the test above, we can safely know that it's functioning to cap the value correctly. If we set the cap to
			// 8k, the same as the entitlement, we can then confirm that the final ppm incentive prices are equal, ensuring
			// that the entitlement calculation is adding up correctly.
			weightOverride = unit.Pound(24500)
			allowableWeightOverride := unit.Pound(8000)
			oldPPMShipment = factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate:  models.TimePointer(actualMoveDate),
						Status:          models.PPMShipmentStatusWaitingOnCustomer,
						AllowableWeight: &allowableWeightOverride,
						MaxIncentive:    &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			newPPM = oldPPMShipment
			updatedMoveDate = time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			newPPM.ActualMoveDate = models.TimePointer(updatedMoveDate)

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinalIncentiveLimitedByAllowableWeight, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"), "50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight = SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(10000), originalWeight)
			suite.Equal(unit.Pound(10000), newWeight)

			// Confirm the incentives are equal with the same value for allowable and entitlement weight caps with equal distances.
			suite.Equal(*ppmFinalIncentiveLimitedByAllowableWeight, *ppmFinal)
		})

		suite.Run("Final Incentive - Success with updated weights", func() {
			setupPricerData()
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			maxIncentive := unit.Cents(90000000)
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						MaxIncentive:   &maxIncentive,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})

			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), nil, nil),
			}

			newPPM := oldPPMShipment
			weightOverride := unit.Pound(19500)
			newPPM.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight: &weightOverride,
						},
					},
				}, nil),
			}

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			suite.NotEqual(*oldPPMShipment.ActualMoveDate, newPPM.ActualMoveDate)
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(4000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)
			suite.Equal(unit.Cents(80249474), *ppmFinal)
		})

		suite.Run("Final Incentive - Success with disregarding rejected weight tickets", func() {
			setupPricerData()
			oldEmptyWeight := unit.Pound(6000)
			oldFullWeight := unit.Pound(10000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			newPPM := oldPPMShipment
			newWeightTicket := newPPM.WeightTickets[0]
			rejected := models.PPMDocumentStatusRejected
			newWeightTicket.Status = &rejected
			newPPM.WeightTickets = models.WeightTickets{newWeightTicket}
			// At this point the updated weight tickets on the newPPMShipment could be saved to the DB
			// the save is being omitted here to reduce DB calls in our test

			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(4000), originalWeight)
			suite.Equal(unit.Pound(0), newWeight)
			suite.Nil(ppmFinal)
		})

		suite.Run("Final Incentive - Success updating finalIncentive with rejected weight tickets", func() {
			setupPricerData()
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			maxIncentive := unit.Cents(90000000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
						MaxIncentive:   &maxIncentive,
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
			oldPPMShipment.WeightTickets = models.WeightTickets{
				oldPPMShipment.WeightTickets[0],
				factory.BuildWeightTicket(suite.DB(), nil, nil),
			}

			newPPM := oldPPMShipment
			rejected := models.PPMDocumentStatusRejected
			approved := models.PPMDocumentStatusApproved
			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.Status = &rejected
			newWeightTicket2 := newPPM.WeightTickets[1]
			newWeightTicket2.Status = &approved
			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1, newWeightTicket2}
			// At this point the updated weight tickets on the newPPMShipment could be saved to the DB
			// the save is being omitted here to reduce DB calls in our test
			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(8000), originalWeight)
			suite.Equal(unit.Pound(4000), newWeight)
			suite.Equal(unit.Cents(48564567), *ppmFinal)
			suite.NotEqual(oldPPMShipment.FinalIncentive, *ppmFinal)
		})

		suite.Run("Final Incentive - Success updating finalIncentive when adjusted net weight is taken into account", func() {
			setupPricerData()
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			maxIncentive := unit.Cents(90000000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
						MaxIncentive:   &maxIncentive,
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
			oldPPMShipment.WeightTickets = models.WeightTickets{
				oldPPMShipment.WeightTickets[0],
				factory.BuildWeightTicket(suite.DB(), nil, nil),
			}

			newPPM := oldPPMShipment
			rejected := models.PPMDocumentStatusRejected
			approved := models.PPMDocumentStatusApproved
			adjustedNetWeight := unit.Pound(3000)

			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.AdjustedNetWeight = &adjustedNetWeight
			newWeightTicket1.Status = &rejected

			newWeightTicket2 := newPPM.WeightTickets[1]
			newWeightTicket2.AdjustedNetWeight = &adjustedNetWeight
			newWeightTicket2.Status = &approved

			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1, newWeightTicket2}

			// At this point the updated weight tickets on the newPPMShipment could be saved to the DB
			// the save is being omitted here to reduce DB calls in our test
			mockedPaymentRequestHelper.On(
				"FetchServiceParamsForServiceItems",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

			// DTOD distance is going to be less than the HHG Rand McNally distance of 2361 miles
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			ppmFinal, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)

			mockedPlanner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813")
			mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(8000), originalWeight)
			suite.Equal(unit.Pound(3000), newWeight)
			suite.Equal(unit.Cents(36423265), *ppmFinal)
			suite.NotEqual(oldPPMShipment.FinalIncentive, *ppmFinal)
		})

		suite.Run("Sum Weights - sum weights for original shipment with standard weight ticket and new shipment with standard weight ticket", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
			newPPM := oldPPMShipment
			newFullWeight := unit.Pound(8000)
			newEmptyWeight := unit.Pound(3000)
			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.FullWeight = &newFullWeight
			newWeightTicket1.EmptyWeight = &newEmptyWeight
			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1}

			//Both PPM's have valid weight tickets so both should return properly calculated totals
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(4000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)
		})

		suite.Run("Sum Weights - sum weights for original shipment with standard weight ticket and new shipment with standard weight ticket & rejected ticket", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			newPPM := oldPPMShipment
			newFullWeight1 := unit.Pound(8000)
			newEmptyWeight1 := unit.Pound(3000)
			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.FullWeight = &newFullWeight1
			newWeightTicket1.EmptyWeight = &newEmptyWeight1

			newFullWeight2 := unit.Pound(12000)
			newEmptyWeight2 := unit.Pound(4000)
			rejected := models.PPMDocumentStatusRejected
			newWeightTicket2 := newPPM.WeightTickets[0]
			newWeightTicket2.FullWeight = &newFullWeight2
			newWeightTicket2.EmptyWeight = &newEmptyWeight2
			newWeightTicket2.Status = &rejected

			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1, newWeightTicket2}

			//Weight for rejected ticket should NOT be included in newWeight total
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(4000), originalWeight)
			suite.Equal(unit.Pound(5000), newWeight)
		})

		suite.Run("Sum Weights - sum weights for original shipment with rejected weight ticket and new shipment with standard weight tickets", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			rejected := models.PPMDocumentStatusRejected
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
						Status:      &rejected,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			approved := models.PPMDocumentStatusApproved
			newPPM := oldPPMShipment
			newFullWeight1 := unit.Pound(8000)
			newEmptyWeight1 := unit.Pound(3000)
			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.FullWeight = &newFullWeight1
			newWeightTicket1.EmptyWeight = &newEmptyWeight1
			newWeightTicket1.Status = &approved

			newFullWeight2 := unit.Pound(12000)
			newEmptyWeight2 := unit.Pound(4000)
			newWeightTicket2 := newPPM.WeightTickets[0]
			newWeightTicket2.FullWeight = &newFullWeight2
			newWeightTicket2.EmptyWeight = &newEmptyWeight2
			newWeightTicket2.Status = &approved

			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1, newWeightTicket2}

			//Weight for rejected ticket should NOT be included in oldWeight total
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(0), originalWeight)
			suite.Equal(unit.Pound(13000), newWeight)
		})

		suite.Run("Sum Weights - sum weights for original shipment and new shipment with adjusted weight", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			approved := models.PPMDocumentStatusApproved
			newPPM := oldPPMShipment
			newFullWeight1 := unit.Pound(8000)
			newEmptyWeight1 := unit.Pound(3000)
			adjustedNetWeight1 := unit.Pound(4000)
			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.FullWeight = &newFullWeight1
			newWeightTicket1.EmptyWeight = &newEmptyWeight1
			newWeightTicket1.AdjustedNetWeight = &adjustedNetWeight1
			newWeightTicket1.Status = &approved

			newFullWeight2 := unit.Pound(12000)
			newEmptyWeight2 := unit.Pound(4000)
			adjustedNetWeight2 := unit.Pound(5000)
			newWeightTicket2 := newPPM.WeightTickets[0]
			newWeightTicket2.FullWeight = &newFullWeight2
			newWeightTicket2.EmptyWeight = &newEmptyWeight2
			newWeightTicket2.AdjustedNetWeight = &adjustedNetWeight2
			newWeightTicket2.Status = &approved

			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1, newWeightTicket2}

			//Weight for rejected ticket should NOT be included in oldWeight total
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(4000), originalWeight)
			//13000 comes from the full & empty weights being summed which we do not want in this scenario
			suite.NotEqual(unit.Pound(13000), newWeight)
			suite.Equal(unit.Pound(9000), newWeight)
		})

		suite.Run("Sum Weights - sum weights for original shipment and new shipment with 2 adjusted weights - one of them having a rejected status", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			newPPM := oldPPMShipment
			approved := models.PPMDocumentStatusApproved
			newFullWeight1 := unit.Pound(8000)
			newEmptyWeight1 := unit.Pound(3000)
			adjustedNetWeight1 := unit.Pound(4000)
			newWeightTicket1 := newPPM.WeightTickets[0]
			newWeightTicket1.FullWeight = &newFullWeight1
			newWeightTicket1.EmptyWeight = &newEmptyWeight1
			newWeightTicket1.AdjustedNetWeight = &adjustedNetWeight1
			newWeightTicket1.Status = &approved

			rejected := models.PPMDocumentStatusRejected
			newFullWeight2 := unit.Pound(12000)
			newEmptyWeight2 := unit.Pound(4000)
			adjustedNetWeight2 := unit.Pound(5000)
			newWeightTicket2 := newPPM.WeightTickets[0]
			newWeightTicket2.FullWeight = &newFullWeight2
			newWeightTicket2.EmptyWeight = &newEmptyWeight2
			newWeightTicket2.AdjustedNetWeight = &adjustedNetWeight2
			newWeightTicket2.Status = &rejected

			newPPM.WeightTickets = models.WeightTickets{newWeightTicket1, newWeightTicket2}

			//Weight for rejected ticket should NOT be included in oldWeight total
			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			suite.Equal(unit.Pound(4000), originalWeight)
			//13000 comes from the full & empty weights being summed which we do not want in this scenario
			suite.NotEqual(unit.Pound(13000), newWeight)
			suite.Equal(unit.Pound(4000), newWeight)
		})

		suite.Run("Sum Weights - sum weights for original shipment and new shipment with 2 adjusted moving expense statuses - PPM-SPR", func() {
			trackingNumber := "TRK1234"
			isProGear := true
			proGearBelongsToSelf := true
			proGearDescription := "Pro gear updated description"
			weightShipped := 1000
			ppmSpr := models.PPMTypeSmallPackage
			spr := models.MovingExpenseReceiptTypeSmallPackage
			approvedStatus := models.PPMDocumentStatusApproved
			rejectedStatus := models.PPMDocumentStatusRejected
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						PPMType:        ppmSpr,
						ActualMoveDate: models.TimePointer(moveDate),
					},
				},
			})

			expense1 := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
				{
					Model:    oldPPMShipment,
					LinkOnly: true,
				},
				{
					Model: models.MovingExpense{
						MovingExpenseType:    &spr,
						Status:               &rejectedStatus,
						PaidWithGTCC:         models.BoolPointer(false),
						MissingReceipt:       models.BoolPointer(false),
						Amount:               models.CentPointer(unit.Cents(8675309)),
						TrackingNumber:       &trackingNumber,
						IsProGear:            &isProGear,
						ProGearBelongsToSelf: &proGearBelongsToSelf,
						ProGearDescription:   &proGearDescription,
						WeightShipped:        (*unit.Pound)(&weightShipped),
					},
				},
			}, nil)

			expense2 := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
				{
					Model:    oldPPMShipment,
					LinkOnly: true,
				},
				{
					Model: models.MovingExpense{
						MovingExpenseType:    &spr,
						Status:               &rejectedStatus,
						PaidWithGTCC:         models.BoolPointer(false),
						MissingReceipt:       models.BoolPointer(false),
						Amount:               models.CentPointer(unit.Cents(8675309)),
						TrackingNumber:       &trackingNumber,
						IsProGear:            &isProGear,
						ProGearBelongsToSelf: &proGearBelongsToSelf,
						ProGearDescription:   &proGearDescription,
						WeightShipped:        (*unit.Pound)(&weightShipped),
					},
				},
			}, nil)

			oldPPMShipment.MovingExpenses = models.MovingExpenses{expense1, expense2}
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
			newPPM := oldPPMShipment

			// changing moving expense statuses to be approved
			newMovingExpense1 := newPPM.MovingExpenses[0]
			newMovingExpense1.Status = &approvedStatus

			newMovingExpense2 := newPPM.MovingExpenses[1]
			newMovingExpense2.Status = &approvedStatus

			newPPM.MovingExpenses = models.MovingExpenses{newMovingExpense1, newMovingExpense2}

			originalWeight, newWeight := SumWeights(oldPPMShipment, newPPM)
			// should be 0 because both were rejected
			suite.Equal(unit.Pound(0), originalWeight)
			// should be 2000 because both are now accepted so we add them together
			suite.Equal(unit.Pound(2000), newWeight)
		})

		suite.Run("Should Skip Calculating Final Incentive - should return false when the move date is changed", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			newPPMShipment := oldPPMShipment
			updatedMoveDate := time.Date(2020, time.March, 25, 0, 0, 0, 0, time.UTC)
			newPPMShipment.ActualMoveDate = models.TimePointer(updatedMoveDate)

			originalTotalWeight, newTotalWeight := SumWeights(oldPPMShipment, newPPMShipment)
			skipCalculateFinalIncentive := shouldSkipCalculatingFinalIncentive(&newPPMShipment, &oldPPMShipment, originalTotalWeight, newTotalWeight)
			suite.Equal(false, skipCalculateFinalIncentive)
		})

		suite.Run("Should Skip Calculating Final Incentive - should return false when the destination or pickup postal code is changed", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldShipmentPickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
				{
					Model: models.Address{
						StreetAddress1: "123 Main St",
						City:           "Beverly Hills",
						State:          "CA",
						PostalCode:     "90210",
					},
				},
			}, nil)
			oldShipmentDestinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
				{
					Model: models.Address{
						StreetAddress1: "321 Turbo St",
						City:           "Augusta",
						State:          "GA",
						PostalCode:     "30813",
					},
				},
			}, nil)

			newShipment1DestinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
				{
					Model: models.Address{
						StreetAddress1: "5 Jayden St",
						City:           "Augusta",
						State:          "GA",
						PostalCode:     "20906",
					},
				},
			}, nil)

			newShipment2PickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
				{
					Model: models.Address{
						StreetAddress1: "8 Ovechkin Ave",
						City:           "Beverly Hills",
						State:          "CA",
						PostalCode:     "99011",
					},
				},
			}, nil)

			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model:    oldShipmentPickupAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.PickupAddress,
				},
				{
					Model:    oldShipmentDestinationAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})
			newPPMShipment1 := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model:    oldShipmentPickupAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.PickupAddress,
				},
				{
					Model:    newShipment1DestinationAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})
			newPPMShipment2 := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model:    newShipment2PickupAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.PickupAddress,
				},
				{
					Model:    newShipment1DestinationAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			// Assert false is returned when the actual destination address postal code is changed
			originalTotalWeight1, newTotalWeight1 := SumWeights(oldPPMShipment, newPPMShipment1)
			skipCalculateFinalIncentive1 := shouldSkipCalculatingFinalIncentive(&newPPMShipment1, &oldPPMShipment, originalTotalWeight1, newTotalWeight1)
			suite.Equal(false, skipCalculateFinalIncentive1)

			originalTotalWeight2, newTotalWeight2 := SumWeights(oldPPMShipment, newPPMShipment2)
			skipCalculateFinalIncentive2 := shouldSkipCalculatingFinalIncentive(&newPPMShipment2, &oldPPMShipment, originalTotalWeight2, newTotalWeight2)
			suite.Equal(false, skipCalculateFinalIncentive2)
		})

		suite.Run("Should Skip Calculating Final Incentive - should return false when adjustedNetWeight is taken into account", func() {
			oldFullWeight := unit.Pound(10000)
			oldEmptyWeight := unit.Pound(6000)
			moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			oldPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate: models.TimePointer(moveDate),
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
					},
				},
				{
					Model: models.WeightTicket{
						FullWeight:  &oldFullWeight,
						EmptyWeight: &oldEmptyWeight,
					},
				},
			})

			// tests pass even if status is Needs Payment Approval,
			// but preserve in case it matters
			oldPPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

			newPPMShipment := oldPPMShipment
			newFullWeight := unit.Pound(10000)
			newEmptyWeight := unit.Pound(3000)
			adjustedNetWeight := unit.Pound(6000)
			approved := models.PPMDocumentStatusApproved

			newWeightTicket := newPPMShipment.WeightTickets[0]
			newWeightTicket.FullWeight = &newFullWeight
			newWeightTicket.EmptyWeight = &newEmptyWeight
			newWeightTicket.AdjustedNetWeight = &adjustedNetWeight
			newWeightTicket.Status = &approved
			newPPMShipment.WeightTickets = models.WeightTickets{newWeightTicket}

			originalTotalWeight, newTotalWeight := SumWeights(oldPPMShipment, newPPMShipment)
			suite.Equal(unit.Pound(4000), originalTotalWeight)
			suite.Equal(unit.Pound(6000), newTotalWeight)

			//Func should notice one of the total weights are different, triggering the recalculation
			skipCalculateFinalIncentive := shouldSkipCalculatingFinalIncentive(&newPPMShipment, &oldPPMShipment, originalTotalWeight, newTotalWeight)
			suite.Equal(false, skipCalculateFinalIncentive)
		})

		suite.Run("Final Incentive - does not change when required fields are the same", func() {
			setupPricerData()
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
						ActualMoveDate: models.TimePointer(actualMoveDate),
					},
				},
			}, nil)
			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), nil, nil),
			}
			newPPM := oldPPMShipment
			address := factory.BuildAddress(suite.DB(), nil, nil)
			newPPM.W2Address = &address

			finalIncentive, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.True(oldPPMShipment.ActualMoveDate.Equal(*newPPM.ActualMoveDate))
			suite.Equal(*oldPPMShipment.FinalIncentive, *finalIncentive)
		})

		suite.Run("Final Incentive - does not change when status is not WAITINGONCUSTOMER or NEEDSPAYMENTAPPROVAL", func() {
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status:         models.PPMShipmentStatusNeedsAdvanceApproval,
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
						ActualMoveDate: models.TimePointer(actualMoveDate),
					},
				},
			}, nil)

			newPPM := oldPPMShipment
			newPPM.Status = models.PPMShipmentStatusCloseoutComplete

			finalIncentive, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.Equal(oldPPMShipment.FinalIncentive, finalIncentive)
		})

		suite.Run("Final Incentive - set to nil when missing info", func() {
			setupPricerData()
			oldPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						Status:         models.PPMShipmentStatusWaitingOnCustomer,
						FinalIncentive: models.CentPointer(unit.Cents(500000)),
						ActualMoveDate: models.TimePointer(actualMoveDate),
					},
				},
			}, nil)
			oldPPMShipment.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), nil, nil),
			}

			newPPM := oldPPMShipment
			newPPM.WeightTickets = nil

			finalIncentive, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.True(oldPPMShipment.ActualMoveDate.Equal(*newPPM.ActualMoveDate))
			suite.Nil(finalIncentive)
		})
	})

	suite.Run("SIT Estimated Cost", func() {
		// For comparison should be priced the same as ORGSIT in devseed
		suite.Run("Success - Origin First Day and Additional Day SIT", func() {
			setupPricerData()

			originLocation := models.SITLocationTypeOrigin
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypePPM,
					},
				},
			}, nil)

			shipmentOriginSIT := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &originLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					},
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
			}, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), models.PPMShipment{}, &shipmentOriginSIT)

			suite.NoError(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(62720, estimatedSITCost.Int())
		})

		suite.Run("Success - Destination First Day and Additional Day SIT", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypePPM,
					},
				},
			}, nil)
			shipmentDestinationSIT := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "987 Other Avenue",
						StreetAddress2: models.StringPointer("P.O. Box 1234"),
						StreetAddress3: models.StringPointer("c/o Another Person"),
						City:           "Des Moines",
						State:          "IA",
						PostalCode:     "50309",
						County:         models.StringPointer("POLK"),
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "987 Other Avenue",
						StreetAddress2: models.StringPointer("P.O. Box 12345"),
						StreetAddress3: models.StringPointer("c/o Another Person"),
						City:           "Fort Eisenhower",
						State:          "GA",
						PostalCode:     "30813",
						County:         models.StringPointer("COLUMBIA"),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
			}, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), models.PPMShipment{}, &shipmentDestinationSIT)

			suite.NoError(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(72380, estimatedSITCost.Int())
		})

		suite.Run("Success - same entry and departure dates only prices first day SIT", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						ShipmentType: models.MTOShipmentTypePPM,
					},
				},
			}, nil)

			shipmentOriginSIT := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: &entryDate,
					},
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
			}, nil)
			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30813").Return(2294, nil)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), models.PPMShipment{}, &shipmentOriginSIT)

			suite.NoError(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(35780, estimatedSITCost.Int())
		})

		suite.Run("SIT cost is not calculated when required fields are missing", func() {
			setupPricerData()

			destinationSITLocation := models.SITLocationTypeDestination

			// an MTO Shipment ID is required for the shipment query
			shipmentSITFieldsNotUpdated := factory.BuildPPMShipment(suite.DB(), nil, nil)
			shipmentSITNotExpected := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model:    shipmentSITFieldsNotUpdated.Shipment,
					LinkOnly: true,
				},
			}, nil)
			shipmentSITWeightMissing := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationSITLocation,
						SITEstimatedEntryDate:     models.TimePointer(time.Now()),
						SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
					},
				},
				{
					Model:    shipmentSITFieldsNotUpdated.Shipment,
					LinkOnly: true,
				},
			}, nil)
			shipmentSITEntryDateMissing := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationSITLocation,
						SITEstimatedDepartureDate: models.TimePointer(time.Now()),
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
					},
				},
				{
					Model:    shipmentSITFieldsNotUpdated.Shipment,
					LinkOnly: true,
				},
			}, nil)
			shipmentSITDepartureDateMissing := factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:           models.BoolPointer(true),
						SITLocation:           &destinationSITLocation,
						SITEstimatedEntryDate: models.TimePointer(time.Now()),
						SITEstimatedWeight:    models.PoundPointer(unit.Pound(2999)),
					},
				},
				{
					Model:    shipmentSITFieldsNotUpdated.Shipment,
					LinkOnly: true,
				},
			}, nil)
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
				_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), testCase.oldShipment, &testCase.newShipment) //#nosec G601
				suite.NoError(err, fmt.Sprintf("unexpected error running test %q", testCase.name))
				suite.Nil(estimatedSITCost, fmt.Sprintf("SIT cost was calculated when it shouldnt't have been during test %q", testCase.name))
			}
		})

		suite.Run("SIT cost is not re-calculated when fields are unchanged", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			shipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
						SITEstimatedEntryDate:     models.TimePointer(time.Now()),
						SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
						SITEstimatedCost:          models.CentPointer(unit.Cents(89900)),
					},
				},
			}, nil)
			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), shipment, &shipment)
			suite.NoError(err)
			suite.Equal(*shipment.SITEstimatedCost, *estimatedSITCost)
		})

		suite.Run("SIT cost is re-calculated when any dependent field is changed", func() {
			setupPricerData()

			destinationLocation := models.SITLocationTypeDestination
			move := factory.BuildMove(suite.DB(), []factory.Customization{
				{
					Model: models.Order{
						ID: uuid.Must(uuid.NewV4()),
					},
				},
				{
					Model: models.Entitlement{
						ID:                 uuid.Must(uuid.NewV4()),
						DBAuthorizedWeight: models.IntPointer(2000),
					},
				},
			}, nil)
			originalShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
						SITEstimatedEntryDate:     models.TimePointer(time.Now()),
						SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
						SITEstimatedCost:          models.CentPointer(unit.Cents(89900)),
					},
				},
			}, nil)
			// PPM base shipment field changes will affect SIT pricing
			shipmentDifferentPickup := originalShipment
			pickupAddress := models.Address{
				StreetAddress1: originalShipment.PickupAddress.StreetAddress1,
				StreetAddress2: originalShipment.PickupAddress.StreetAddress2,
				StreetAddress3: originalShipment.PickupAddress.StreetAddress3,
				City:           originalShipment.PickupAddress.City,
				State:          originalShipment.PickupAddress.State,
				PostalCode:     "90211",
			}
			shipmentDifferentPickup.PickupAddress = &pickupAddress

			shipmentDifferentDestination := originalShipment
			destinationAddress := models.Address{
				StreetAddress1: originalShipment.PickupAddress.StreetAddress1,
				StreetAddress2: originalShipment.PickupAddress.StreetAddress2,
				StreetAddress3: originalShipment.PickupAddress.StreetAddress3,
				City:           originalShipment.PickupAddress.City,
				State:          originalShipment.PickupAddress.State,
				PostalCode:     "30814",
			}
			shipmentDifferentDestination.DestinationAddress = &destinationAddress

			shipmentDifferentDeparture := originalShipment
			// original date was Mar 15th so adding 3 months should affect the date peak period pricing
			shipmentDifferentDeparture.ExpectedDepartureDate = originalShipment.ExpectedDepartureDate.Add(time.Hour * 24 * 70)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"90211", "30813").Return(2294, nil)

			mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "30814").Return(2290, nil)

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
			originalShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &destinationLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2999)),
						SITEstimatedEntryDate:     models.TimePointer(time.Now()),
						SITEstimatedDepartureDate: models.TimePointer(time.Now().Add(time.Hour * 24)),
						SITEstimatedCost:          models.CentPointer(unit.Cents(89900)),
					},
				},
			}, nil)
			shipmentSITNotExpected := originalShipment
			shipmentSITNotExpected.SITExpected = models.BoolPointer(false)

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), originalShipment, &shipmentSITNotExpected)
			suite.NoError(err)
			suite.Nil(shipmentSITNotExpected.SITEstimatedCost)
			suite.Nil(estimatedSITCost)
		})
	})
}

func (suite *PPMShipmentSuite) TestInternationalPPMEstimator() {
	planner := &mocks.Planner{}
	paymentRequestHelper := &prhelpermocks.Helper{}
	ppmEstimator := NewEstimatePPM(planner, paymentRequestHelper)

	setupPricerData := func() {
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(2020, time.December, 31, 12, 0, 0, 0, time.UTC)
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
	}

	suite.Run("Estimated Incentive", func() {
		suite.Run("Estimated Incentive - Success using estimated weight and not db authorized weight for CONUS -> OCONUS", func() {
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			setupPricerData()

			estimatedWeight := unit.Pound(5000)
			newPPM := ppm
			newPPM.EstimatedWeight = &estimatedWeight

			planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"74133", "98421").Return(3000, nil).Twice()

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(ppmEstimate)

			// it should've called from the pickup -> port and NOT pickup -> dest
			planner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"74133", "98421")
			suite.Equal(unit.Cents(504512), *ppmEstimate)

			// appending this to test functionality of the GCC multiplier
			validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
			ppmWithMultiplier := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ExpectedDepartureDate: validGccMultiplierDate,
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)
			newPPMWithMultiplier := ppmWithMultiplier
			newPPMWithMultiplier.EstimatedWeight = &estimatedWeight // setting weight to 5000
			ppmEstimateWithMultiplier, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), ppmWithMultiplier, &newPPMWithMultiplier)
			suite.NilOrNoVerrs(err)

			suite.Equal(unit.Pound(5000), *newPPMWithMultiplier.EstimatedWeight)
			suite.NotEqual(unit.Cents(504512), *ppmEstimateWithMultiplier)
			suite.Equal(unit.Cents(771427), *ppmEstimateWithMultiplier)
		})

		suite.Run("Estimated Incentive - Success using estimated weight and not db authorized weight for OCONUS -> CONUS", func() {
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.PickupAddress,
				},
			}, nil)

			setupPricerData()

			estimatedWeight := unit.Pound(5000)
			newPPM := ppm
			newPPM.EstimatedWeight = &estimatedWeight

			planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"98421", "74133").Return(3000, nil)

			ppmEstimate, _, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(ppmEstimate)

			// it should've called from the pickup -> port and NOT pickup -> dest
			planner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"98421", "74133")
			suite.Equal(unit.Cents(464562), *ppmEstimate)
		})
	})

	suite.Run("Max Incentive", func() {
		suite.Run("Max Incentive - Success using db authorized weight and not estimated for CONUS -> OCONUS", func() {
			oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
				},
			}, nil)
			destDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
				{
					Model: models.DutyLocation{
						Name:      "Test OCONUS Duty Location",
						AddressID: oconusAddress.ID,
					},
				},
			}, nil)
			order := factory.BuildOrder(suite.DB(), []factory.Customization{
				{
					Model: models.Order{
						NewDutyLocationID: destDutyLocation.ID,
					},
				},
			}, nil)
			// when the PPM shipment is in draft, we use the estimated weight and not the db authorized weight
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.Move{
						OrdersID: order.ID,
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			setupPricerData()

			estimatedWeight := unit.Pound(5000)
			newPPM := ppm
			newPPM.EstimatedWeight = &estimatedWeight

			// DTOD will be called to get the distance between the origin duty location & the Tacoma Port ZIP
			planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "98421").Return(3000, nil).Twice()

			ppmMaxIncentive, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(ppmMaxIncentive)

			// it should've called from the pickup -> port and NOT pickup -> dest
			planner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"50309", "98421")
			suite.Equal(unit.Cents(720983), *ppmMaxIncentive)

			// appending this to test functionality of the GCC multiplier
			validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
			ppmWithMultiplier := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ExpectedDepartureDate: validGccMultiplierDate,
					},
				},
				{
					Model: models.Move{
						OrdersID: order.ID,
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)
			newPPMWithMultiplier := ppmWithMultiplier
			newPPMWithMultiplier.EstimatedWeight = &estimatedWeight // setting weight to 5000
			ppmEstimateWithMultiplier, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), ppmWithMultiplier, &newPPMWithMultiplier)
			suite.NilOrNoVerrs(err)

			suite.Equal(unit.Pound(5000), *newPPMWithMultiplier.EstimatedWeight)
			suite.NotEqual(unit.Cents(504512), *ppmEstimateWithMultiplier)
			suite.Equal(unit.Cents(1103119), *ppmEstimateWithMultiplier)
		})

		suite.Run("Max Incentive - Success using db authorized weight and not estimated for OCONUS -> CONUS", func() {
			oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
				},
			}, nil)
			pickupDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
				{
					Model: models.DutyLocation{
						Name:      "Test OCONUS Duty Location",
						AddressID: oconusAddress.ID,
					},
				},
			}, nil)
			order := factory.BuildOrder(suite.DB(), []factory.Customization{
				{
					Model: models.Order{
						OriginDutyLocationID: &pickupDutyLocation.ID,
					},
				},
			}, nil)
			// when the PPM shipment is in draft, we use the estimated weight and not the db authorized weight
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.Move{
						OrdersID: order.ID,
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			setupPricerData()

			estimatedWeight := unit.Pound(5000)
			newPPM := ppm
			newPPM.EstimatedWeight = &estimatedWeight

			// DTOD will be called to get the distance between the origin duty location & the Tacoma Port ZIP
			planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"98421", "30813").Return(3000, nil)

			ppmMaxIncentive, err := ppmEstimator.MaxIncentive(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(ppmMaxIncentive)

			// it should've called from the pickup -> port and NOT pickup -> dest
			planner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"98421", "30813")
			suite.Equal(unit.Cents(743383), *ppmMaxIncentive)
		})
	})

	suite.Run("Final Incentive", func() {
		suite.Run("Final Incentive - Success using estimated weight for CONUS -> OCONUS", func() {
			updatedMoveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate:  models.TimePointer(updatedMoveDate),
						Status:          models.PPMShipmentStatusWaitingOnCustomer,
						EstimatedWeight: models.PoundPointer(4000),
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			newPPM := ppm
			newFullWeight := unit.Pound(8000)
			newEmptyWeight := unit.Pound(3000)
			newPPM.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight:  &newFullWeight,
							EmptyWeight: &newEmptyWeight,
						},
					},
				}, nil),
			}

			setupPricerData()

			planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"74133", "98421").Return(3000, nil).Twice()

			ppmFinalIncentive, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(ppmFinalIncentive)

			// it should've called from the pickup -> port and NOT pickup -> dest
			planner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"74133", "98421")
			suite.Equal(unit.Cents(459178), *ppmFinalIncentive)

			// appending this to test functionality of the GCC multiplier
			validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
			ppmWithMultiplier := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate:        models.TimePointer(updatedMoveDate),
						Status:                models.PPMShipmentStatusWaitingOnCustomer,
						EstimatedWeight:       models.PoundPointer(4000),
						ExpectedDepartureDate: validGccMultiplierDate,
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)
			newPPMWithMultiplier := ppmWithMultiplier
			newPPMWithMultiplier.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight:  &newFullWeight,
							EmptyWeight: &newEmptyWeight,
						},
					},
				}, nil),
			}
			ppmEstimateWithMultiplier, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), ppmWithMultiplier, &newPPMWithMultiplier)
			suite.NilOrNoVerrs(err)

			suite.Equal(unit.Pound(4000), *newPPMWithMultiplier.EstimatedWeight)
			suite.NotEqual(unit.Cents(459178), *ppmEstimateWithMultiplier)
			suite.Equal(unit.Cents(596931), *ppmEstimateWithMultiplier)
		})

		suite.Run("Final Incentive - Success using estimated weight for OCONUS -> CONUS", func() {
			updatedMoveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						ActualMoveDate:  models.TimePointer(updatedMoveDate),
						Status:          models.PPMShipmentStatusWaitingOnCustomer,
						EstimatedWeight: models.PoundPointer(4000),
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.PickupAddress,
				},
			}, nil)

			newPPM := ppm
			newFullWeight := unit.Pound(8000)
			newEmptyWeight := unit.Pound(3000)
			newPPM.WeightTickets = models.WeightTickets{
				factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							FullWeight:  &newFullWeight,
							EmptyWeight: &newEmptyWeight,
						},
					},
				}, nil),
			}

			setupPricerData()

			planner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"98421", "74133").Return(3000, nil)

			ppmFinalIncentive, err := ppmEstimator.FinalIncentiveWithDefaultChecks(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(ppmFinalIncentive)

			// it should've called from the pickup -> port and NOT pickup -> dest
			planner.AssertCalled(suite.T(), "ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
				"98421", "74133")
			suite.Equal(unit.Cents(423178), *ppmFinalIncentive)
		})
	})

	suite.Run("SIT Costs for OCONUS PPMs", func() {
		suite.Run("CalculateSITCost - Success using estimated weight for CONUS -> OCONUS", func() {
			originLocation := models.SITLocationTypeOrigin
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						EstimatedWeight:           models.PoundPointer(4000),
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &originLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			newPPM := ppm
			newEstimatedWeight := models.PoundPointer(5500)
			newPPM.SITEstimatedWeight = newEstimatedWeight
			setupPricerData()

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(unit.Cents(27040), *estimatedSITCost)
		})

		suite.Run("CalculateSITCost - Success using estimated weight for CONUS -> OCONUS", func() {
			originLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						EstimatedWeight:           models.PoundPointer(4000),
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &originLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			newPPM := ppm
			newEstimatedWeight := models.PoundPointer(5500)
			newPPM.SITEstimatedWeight = newEstimatedWeight
			setupPricerData()

			_, estimatedSITCost, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), ppm, &newPPM)
			suite.NilOrNoVerrs(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(unit.Cents(46160), *estimatedSITCost)
		})

		suite.Run("CalculatePPMSITEstimatedCost - Success for OCONUS PPM", func() {
			originLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						EstimatedWeight:           models.PoundPointer(4000),
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &originLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			setupPricerData()

			estimatedSITCost, err := ppmEstimator.CalculatePPMSITEstimatedCost(suite.AppContextForTest(), &ppm)
			suite.NilOrNoVerrs(err)
			suite.NotNil(estimatedSITCost)
			suite.Equal(unit.Cents(23080), *estimatedSITCost)
		})

		suite.Run("CalculatePPMSITEstimatedCostBreakdown - Success for OCONUS PPM", func() {
			originLocation := models.SITLocationTypeDestination
			entryDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
			ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
				{
					Model: models.PPMShipment{
						EstimatedWeight:           models.PoundPointer(4000),
						SITExpected:               models.BoolPointer(true),
						SITLocation:               &originLocation,
						SITEstimatedWeight:        models.PoundPointer(unit.Pound(2000)),
						SITEstimatedEntryDate:     &entryDate,
						SITEstimatedDepartureDate: models.TimePointer(entryDate.Add(time.Hour * 24 * 30)),
					},
				},
				{
					Model: models.MTOShipment{
						MarketCode: models.MarketCodeInternational,
					},
				},
				{
					Model: models.Address{
						StreetAddress1: "Tester Address",
						City:           "Tulsa",
						State:          "OK",
						PostalCode:     "74133",
					},
					Type: &factory.Addresses.PickupAddress,
				},
				{
					Model: models.Address{
						StreetAddress1: "JBER",
						City:           "JBER",
						State:          "AK",
						PostalCode:     "99505",
						IsOconus:       models.BoolPointer(true),
					},
					Type: &factory.Addresses.DeliveryAddress,
				},
			}, nil)

			setupPricerData()

			sitCosts, err := ppmEstimator.CalculatePPMSITEstimatedCostBreakdown(suite.AppContextForTest(), &ppm)
			suite.NilOrNoVerrs(err)
			suite.NotNil(sitCosts)
			suite.Equal(unit.Cents(23080), *sitCosts.EstimatedSITCost)
			suite.Equal(unit.Cents(13480), *sitCosts.PriceFirstDaySIT)
			suite.Equal(unit.Cents(9600), *sitCosts.PriceAdditionalDaySIT)
		})
	})
}
