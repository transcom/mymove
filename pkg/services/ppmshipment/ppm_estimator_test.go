package ppmshipment

import (
	"time"

	"github.com/stretchr/testify/mock"

	prhelpermocks "github.com/transcom/mymove/pkg/payment_request/mocks"

	"github.com/transcom/mymove/pkg/route/mocks"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestEstimatedIncentive() {

	mockedPlanner := &mocks.Planner{}
	mockedPaymentRequestHelper := &prhelpermocks.Helper{}
	ppmEstimator := NewEstimatePPM(mockedPlanner, mockedPaymentRequestHelper)

	// To avoid creating all of the re_services and their corresponding params using factories, we can create this
	// mapping to help mock the response
	serviceParamKeys := map[models.ServiceItemParamName]models.ServiceItemParamKey{
		models.ServiceItemParamNameActualPickupDate:                 {Key: models.ServiceItemParamNameActualPickupDate, Type: models.ServiceItemParamTypeDate},
		models.ServiceItemParamNameContractCode:                     {Key: models.ServiceItemParamNameContractCode, Type: models.ServiceItemParamTypeString},
		models.ServiceItemParamNameDistanceZip3:                     {Key: models.ServiceItemParamNameDistanceZip3, Type: models.ServiceItemParamTypeInteger},
		models.ServiceItemParamNameDistanceZip5:                     {Key: models.ServiceItemParamNameDistanceZip5, Type: models.ServiceItemParamTypeInteger},
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
		models.ServiceItemParamNameDistanceZip3,
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
		models.ServiceItemParamNameDistanceZip3,
		models.ServiceItemParamNameDistanceZip5,
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
		testdatagen.MakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
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

		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
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

		testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
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

		testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
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

		testdatagen.MakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReContract:            originDomesticServiceArea.Contract,
			ReDomesticServiceArea: originDomesticServiceArea,
			ReService: models.ReService{
				Code: models.ReServiceCodeDOP,
			},
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				IsPeakPeriod: false,
				PriceCents:   unit.Cents(404),
			},
		})

		testdatagen.MakeReDomesticServiceAreaPrice(suite.DB(), testdatagen.Assertions{
			ReContract:            destDomesticServiceArea.Contract,
			ReDomesticServiceArea: destDomesticServiceArea,
			ReService: models.ReService{
				Code: models.ReServiceCodeDDP,
			},
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				IsPeakPeriod: false,
				PriceCents:   unit.Cents(832),
			},
		})

		testdatagen.MakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
			ReContract: originDomesticServiceArea.Contract,
			ReService: models.ReService{
				Code: models.ReServiceCodeDPK,
			},
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				IsPeakPeriod: false,
				Schedule:     3,
				PriceCents:   7395,
			},
		})

		testdatagen.MakeReDomesticOtherPrice(suite.DB(), testdatagen.Assertions{
			ReContract: destDomesticServiceArea.Contract,
			ReService: models.ReService{
				Code: models.ReServiceCodeDUPK,
			},
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				IsPeakPeriod: false,
				Schedule:     2,
				PriceCents:   597,
			},
		})
	}

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

		mockedPlanner.On("Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"90210", "30813").Return(2361, nil).Once()

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)

		mockedPlanner.AssertCalled(suite.T(), "Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"90210", "30813")
		mockedPaymentRequestHelper.AssertCalled(suite.T(), "FetchServiceParamsForServiceItems", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("[]models.MTOServiceItem"))

		suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.Equal(unit.Pound(5000), *newPPM.EstimatedWeight)
		suite.Equal(unit.Cents(72097231), *ppmEstimate)
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

		mockedPlanner.On("Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"90210", "30813").Return(2361, nil).Once()

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Nil(newPPM.Advance)
		suite.Nil(newPPM.AdvanceRequested)
		suite.Equal(unit.Cents(39319267), *ppmEstimate)
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

		estimatedIncentive, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Equal(oldPPMShipment.PickupPostalCode, newPPM.PickupPostalCode)
		suite.Equal(*oldPPMShipment.EstimatedWeight, *newPPM.EstimatedWeight)
		suite.Equal(oldPPMShipment.DestinationPostalCode, newPPM.DestinationPostalCode)
		suite.True(oldPPMShipment.ExpectedDepartureDate.Equal(newPPM.ExpectedDepartureDate))
		suite.Equal(*oldPPMShipment.EstimatedIncentive, *estimatedIncentive)
		suite.Equal(models.BoolPointer(true), newPPM.AdvanceRequested)
		suite.Equal(unit.Cents(598700), *newPPM.Advance)
	})
	suite.Run("Estimated Incentive - Failure - is not created when status is not DRAFT", func() {
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
			SitExpected:           oldPPMShipment.SitExpected,
			EstimatedIncentive:    models.CentPointer(unit.Cents(500000)),
		}

		ppmEstimate, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NilOrNoVerrs(err)
		suite.Nil(ppmEstimate)
		suite.Equal(models.CentPointer(unit.Cents(500000)), newPPM.EstimatedIncentive)
	})

	suite.Run("Estimated Incentive - Failure - is not created when Estimated Weight is missing", func() {
		oldPPMShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

		newPPM := oldPPMShipment
		newPPM.DestinationPostalCode = "94040"

		_, err := ppmEstimator.EstimateIncentiveWithDefaultChecks(suite.AppContextForTest(), oldPPMShipment, &newPPM)
		suite.NoError(err)
		suite.Nil(newPPM.EstimatedIncentive)
	})
}
