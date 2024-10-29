package ppmcloseout

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	prhelpermocks "github.com/transcom/mymove/pkg/payment_request/mocks"
	"github.com/transcom/mymove/pkg/route/mocks"
	servicemocks "github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmBuildType string

const (
	ppmBuildReadyForCloseout  = "readyForCloseout"
	ppmBuildWaitingOnCustomer = "waitingOnCustomer"
)

func (suite *PPMCloseoutSuite) TestPPMShipmentCreator() {

	// One-time test setup
	mockedPlanner := &mocks.Planner{}
	mockedPaymentRequestHelper := &prhelpermocks.Helper{}
	mockPpmEstimator := &servicemocks.PPMEstimator{}
	ppmCloseoutFetcher := NewPPMCloseoutFetcher(mockedPlanner, mockedPaymentRequestHelper, mockPpmEstimator)
	serviceParams := mockServiceParamsTables()

	suite.PreloadData(func() {
		// Generate all the data needed for a PPM Closeout object to be calculated
		testdatagen.FetchOrMakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2020, time.March, 17, 0, 0, 0, 0, time.UTC),
			},
		})

		originDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "056",
				ServicesSchedule: 3,
				SITPDSchedule:    3,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.AppContextForTest().DB(), testdatagen.Assertions{}),
		})

		testdatagen.FetchOrMakeReContractYear(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             originDomesticServiceArea.Contract,
				ContractID:           originDomesticServiceArea.ContractID,
				StartDate:            time.Now(),
				EndDate:              time.Now().Add(time.Hour * 8760),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originDomesticServiceArea.Contract,
				ContractID:          originDomesticServiceArea.ContractID,
				DomesticServiceArea: originDomesticServiceArea,
				Zip3:                "902",
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originDomesticServiceArea.Contract,
				ContractID:          originDomesticServiceArea.ContractID,
				DomesticServiceArea: originDomesticServiceArea,
				Zip3:                "503",
			},
		})

		testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
				IsPeakPeriod:          false,
			},
		})
		testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
				IsPeakPeriod:          true,
			},
		})

		dopService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDOP)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             dopService.ID,
				Service:               dopService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            unit.Cents(404),
			},
		})

		destDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.AppContextForTest().DB(), testdatagen.Assertions{
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

		ddpService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDDP)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddpService.ID,
				Service:               ddpService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            unit.Cents(832),
			},
		})

		dpkService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDPK)

		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				ContractID:   originDomesticServiceArea.ContractID,
				Contract:     originDomesticServiceArea.Contract,
				ServiceID:    dpkService.ID,
				Service:      dpkService,
				IsPeakPeriod: true,
				Schedule:     3,
				PriceCents:   7395,
			},
		})

		dupkService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDUPK)

		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticOtherPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticOtherPrice: models.ReDomesticOtherPrice{
				ContractID:   destDomesticServiceArea.ContractID,
				Contract:     destDomesticServiceArea.Contract,
				ServiceID:    dupkService.ID,
				Service:      dupkService,
				IsPeakPeriod: true,
				Schedule:     2,
				PriceCents:   597,
			},
		})

		dofsitService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDOFSIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             dofsitService.ID,
				Service:               dofsitService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            1153,
			},
		})

		doasitService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDOASIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            originDomesticServiceArea.ContractID,
				Contract:              originDomesticServiceArea.Contract,
				ServiceID:             doasitService.ID,
				Service:               doasitService,
				DomesticServiceAreaID: originDomesticServiceArea.ID,
				DomesticServiceArea:   originDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            46,
			},
		})

		ddfsitService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDDFSIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddfsitService.ID,
				Service:               ddfsitService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            1612,
			},
		})

		ddasitService := factory.FetchOrBuildReServiceByCode(suite.AppContextForTest().DB(), models.ReServiceCodeDDASIT)

		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
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
		testdatagen.FetchOrMakeReDomesticServiceAreaPrice(suite.AppContextForTest().DB(), testdatagen.Assertions{
			ReDomesticServiceAreaPrice: models.ReDomesticServiceAreaPrice{
				ContractID:            destDomesticServiceArea.ContractID,
				Contract:              destDomesticServiceArea.Contract,
				ServiceID:             ddasitService.ID,
				Service:               ddasitService,
				DomesticServiceAreaID: destDomesticServiceArea.ID,
				DomesticServiceArea:   destDomesticServiceArea,
				IsPeakPeriod:          true,
				PriceCents:            55,
			},
		})
	})

	suite.Run("Can successfully GET a PPMCloseout Object", func() {
		// Under test:	GetPPMCloseout
		// Set up:		Established ZIPs, ReServices, and all pricing data
		// Expected:	PPMCloseout Object successfully retrieved
		appCtx := suite.AppContextForTest()

		mockedPlanner.On("ZipTransitDistance", mock.AnythingOfType("*appcontext.appContext"),
			"50309", "30813").Return(2294, nil)

		mockedPaymentRequestHelper.On(
			"FetchServiceParamsForServiceItems",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("[]models.MTOServiceItem")).Return(serviceParams, nil)

		mockIncentiveValue := unit.Cents(100000)
		mockPpmEstimator.On(
			"FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).Return(&mockIncentiveValue, nil)

		ppmShipment := suite.mockPPMShipmentForCloseoutTest(ppmBuildReadyForCloseout)

		ppmCloseoutObj, err := ppmCloseoutFetcher.GetPPMCloseout(appCtx, ppmShipment.ID)
		if err != nil {
			appCtx.Logger().Error("Error getting PPM closeout object: ", zap.Error(err))
		}

		suite.Nil(err)
		suite.NotNil(ppmCloseoutObj)
		suite.NotEmpty(ppmCloseoutObj)
	})

	suite.Run("Returns a \"NotFoundError\" if the PPM Shipment was not found using the given ID", func() {
		appCtx := suite.AppContextForTest()
		missingPpmID := uuid.Must(uuid.NewV4())
		ppmCloseoutObj, err := ppmCloseoutFetcher.GetPPMCloseout(appCtx, missingPpmID)

		suite.NotNil(err)
		suite.Nil(ppmCloseoutObj)
		suite.IsType(err, apperror.NotFoundError{})
	})

	suite.Run("Returns a \"PPMNotReadyForCloseoutError\" if shipment is not marked as either \"NEEDS_CLOSEOUT\" or \"APPROVED\"", func() {
		appCtx := suite.AppContextForTest()
		ppmShipment := suite.mockPPMShipmentForCloseoutTest(ppmBuildWaitingOnCustomer)
		ppmShipment.Status = models.PPMShipmentStatusSubmitted

		ppmCloseoutObj, err := ppmCloseoutFetcher.GetPPMCloseout(appCtx, ppmShipment.ID)

		suite.NotNil(err)
		suite.Nil(ppmCloseoutObj)
		appCtx.Logger().Debug("+%v", zap.Error(err))
		suite.IsType(err, apperror.PPMNotReadyForCloseoutError{})
	})
}

func mockServiceParamsTables() models.ServiceParams {
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

	return serviceParams
}

func (suite *PPMCloseoutSuite) mockPPMShipmentForCloseoutTest(buildType ppmBuildType) models.PPMShipment {
	ppmID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	estWeight := unit.Pound(2000)
	actualMoveDate := time.Now()
	expectedDepartureDate := &actualMoveDate
	miles := unit.Miles(200)
	emptyWeight1 := unit.Pound(1000)
	emptyWeight2 := unit.Pound(1200)
	fullWeight1 := unit.Pound(1500)
	fullWeight2 := unit.Pound(1500)
	pgBoolCustomer := true
	pgBoolSpouse := false
	weightCustomer := unit.Pound(100)
	weightSpouse := unit.Pound(120)
	finalIncentive := unit.Cents(20000)

	weightTickets := models.WeightTickets{
		models.WeightTicket{
			EmptyWeight: &emptyWeight1,
			FullWeight:  &fullWeight1,
		},
		models.WeightTicket{
			EmptyWeight: &emptyWeight2,
			FullWeight:  &fullWeight2,
		},
	}
	progearWeightTickets := models.ProgearWeightTickets{
		models.ProgearWeightTicket{
			BelongsToSelf: &pgBoolCustomer,
			Weight:        &weightCustomer,
		},
		models.ProgearWeightTicket{
			BelongsToSelf: &pgBoolSpouse,
			Weight:        &weightSpouse,
		},
	}

	sitDaysAllowance := 20
	sitLocation := models.SITLocationTypeOrigin
	date := time.Now()

	ppmShipmentCustomization := []factory.Customization{
		{
			Model: models.MTOShipment{
				Distance:         &miles,
				SITDaysAllowance: &sitDaysAllowance,
			},
		},
		{
			Model: models.PPMShipment{
				ID:                        ppmID,
				ExpectedDepartureDate:     *expectedDepartureDate,
				ActualMoveDate:            &actualMoveDate,
				EstimatedWeight:           &estWeight,
				WeightTickets:             weightTickets,
				ProgearWeightTickets:      progearWeightTickets,
				FinalIncentive:            &finalIncentive,
				SITLocation:               &sitLocation,
				SITEstimatedEntryDate:     &date,
				SITEstimatedDepartureDate: &date,
			},
		},
	}

	var ppmShipment models.PPMShipment
	if buildType == ppmBuildReadyForCloseout {
		ppmShipment = factory.BuildPPMShipmentThatNeedsCloseout(suite.AppContextForTest().DB(), nil, ppmShipmentCustomization)
	} else if buildType == ppmBuildWaitingOnCustomer {
		ppmShipment = factory.BuildPPMShipment(suite.AppContextForTest().DB(), ppmShipmentCustomization, nil)
	}

	return ppmShipment
}
