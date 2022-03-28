package ppmservices

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMServiceSuite) TestCalculateEstimateSuccess() {
	moveID := uuid.FromStringOrNil("02856e5d-cdd1-4403-ad54-60e52e249d0d")

	appCtx := suite.AppContextForTest()

	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	originDutyLocationZip := "94540"
	destDutyLocationZip := "95632"
	move := suite.setupCalculateEstimateTest(moveID, originDutyLocationZip, destDutyLocationZip)
	weightEstimate := unit.Pound(7500)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:           moveID,
			Move:             move,
			PickupPostalCode: swag.String("94540"),
			OriginalMoveDate: swag.Time(time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)),
			WeightEstimate:   &weightEstimate,
			HasSit:           swag.Bool(true),
			DaysInStorage:    swag.Int64(int64(30)),
		},
	})

	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		"94540",
		"95632",
	).Return(3200, nil)
	calculator := NewEstimateCalculator(planner)
	sitCharge, _, err := calculator.CalculateEstimates(suite.AppContextForTest(), &ppm, moveID)
	suite.NoError(err)
	suite.Equal(int64(171401), sitCharge)
}

func (suite *PPMServiceSuite) TestCalculateEstimateNoSITSuccess() {
	moveID := uuid.FromStringOrNil("02856e5d-cdd1-4403-ad54-60e52e249d0d")

	appCtx := suite.AppContextForTest()

	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	originDutyLocationZip := "94540"
	destDutyLocationZip := "95632"
	move := suite.setupCalculateEstimateTest(moveID, originDutyLocationZip, destDutyLocationZip)
	weightEstimate := unit.Pound(7500)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:           moveID,
			Move:             move,
			PickupPostalCode: swag.String("94540"),
			OriginalMoveDate: swag.Time(time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)),
			WeightEstimate:   &weightEstimate,
			HasSit:           swag.Bool(false),
			DaysInStorage:    swag.Int64(int64(30)),
		},
	})

	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		"94540",
		"95632",
	).Return(3200, nil)
	calculator := NewEstimateCalculator(planner)
	sitCharge, _, err := calculator.CalculateEstimates(suite.AppContextForTest(), &ppm, moveID)
	suite.NoError(err)
	suite.Equal(int64(0), sitCharge)
}

func (suite *PPMServiceSuite) TestCalculateEstimateBadMoveIDFails() {
	weightEstimate := unit.Pound(7000)

	moveID := uuid.FromStringOrNil("02856e5d-cdd1-4403-ad54-60e52e249d0d")

	appCtx := suite.AppContextForTest()

	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	pickupZip := "94540"
	originDutyLocationZip := "94540"
	destDutyLocationZip := "95632"
	move := suite.setupCalculateEstimateTest(moveID, originDutyLocationZip, destDutyLocationZip)

	moveDate := time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:           moveID,
			Move:             move,
			PickupPostalCode: &pickupZip,
			OriginalMoveDate: &moveDate,
			WeightEstimate:   &weightEstimate,
		},
	})
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(3200, nil)
	calculator := NewEstimateCalculator(planner)
	nonExistentMoveID, err := uuid.FromString("2ef27bd2-97ae-4808-96cb-0cadd7f48972")
	if err != nil {
		suite.Logger().Fatal("failure to get uuid from string")
	}
	_, _, err = calculator.CalculateEstimates(suite.AppContextForTest(), &ppm, nonExistentMoveID)

	suite.Error(err)
}

// we are currently hardcoding results for PPMs and aren't checking for valid zip codes
// keep this test around for when we do address PPM work
// func (suite *PPMServiceSuite) TestCalculateEstimateBadPickupZipFails() {
// 	weightEstimate := unit.Pound(7000)

// 	moveID := uuid.FromStringOrNil("02856e5d-cdd1-4403-ad54-60e52e249d0d")
//  appCtx := suite.AppContextForTest()
// 	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
// 		suite.FailNow("failed to run scenario 2: %+v", err)
// 	}

// 	invalidPickupZip := "11111"
// 	originDutyLocationZip := "94540"
// 	destDutyLocationZip := "95632"
// 	move := suite.setupCalculateEstimateTest(moveID, originDutyLocationZip, destDutyLocationZip)

// 	moveDate := time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)
// 	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
// 		PersonallyProcuredMove: models.PersonallyProcuredMove{
// 			MoveID:           moveID,
// 			Move:             move,
// 			PickupPostalCode: &invalidPickupZip,
// 			OriginalMoveDate: &moveDate,
// 			WeightEstimate:   &weightEstimate,
// 		},
// 	})
// 	planner := &mocks.Planner{}
// 	planner.On("Zip5TransitDistanceLineHaul",
// 		mock.Anything,
// 		mock.Anything,
// 	).Return(3200, nil)
// 	calculator := NewEstimateCalculator(suite.DB(), planner)
// 	_, _, err := calculator.CalculateEstimates(&ppm, moveID, suite.logger)

// 	suite.Error(err)
// }

// PPMs estimates are being hardcoded because we are not loading tariff400ng data
// bypass tests as we are returning hard coded values and not checking zips right now
//
// PPMs estimates are being hardcoded because we are not loading tariff400ng data
// bypass tests as we are returning hard coded values and not checking zips right now
// func (suite *PPMServiceSuite) TestCalculateEstimateNewDutyLocationZipFails() {
// 	moveID := uuid.FromStringOrNil("02856e5d-cdd1-4403-ad54-60e52e249d0d")
//  appCtx := suite.AppContextForTest()
// 	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
// 		suite.FailNow("failed to run scenario 2: %+v", err)
// 	}
// 	originDutyLocationZip := "94540"
// 	invalidDestDutyLocationZip := "00000"
// 	pickupZip := "94540"
// 	move := suite.setupCalculateEstimateTest(moveID, originDutyLocationZip, invalidDestDutyLocationZip)

// 	moveDate := time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)
// 	weightEstimate := unit.Pound(7500)
// 	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
// 		PersonallyProcuredMove: models.PersonallyProcuredMove{
// 			MoveID:           moveID,
// 			Move:             move,
// 			PickupPostalCode: &pickupZip,
// 			OriginalMoveDate: &moveDate,
// 			WeightEstimate:   &weightEstimate,
// 		},
// 	})

// 	planner := &mocks.Planner{}
// 	planner.On("Zip5TransitDistanceLineHaul",
// 		mock.Anything,
// 		mock.Anything,
// 	).Return(3200, nil)
// 	calculator := NewEstimateCalculator(suite.DB(), planner)
// 	_, _, err := calculator.CalculateEstimates(&ppm, moveID, suite.logger)
// 	suite.Error(err)
// }

// we are currently hardcoding results for PPMs and aren't checking for valid zip codes
// keep this test around for when we do address PPM work
// func (suite *PPMServiceSuite) TestCalculateEstimateInvalidWeightFails() {
// 	moveID := uuid.FromStringOrNil("02856e5d-cdd1-4403-ad54-60e52e249d0d")
//  appCtx := suite.AppContextForTest()
// 	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
// 		suite.FailNow("failed to run scenario 2: %+v", err)
// 	}

// 	originDutyLocationZip := "94540"
// 	destDutyLocationZip := "95632"
// 	pickupZip := "94540"
// 	move := suite.setupCalculateEstimateTest(moveID, originDutyLocationZip, destDutyLocationZip)

// 	moveDate := time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)
// 	weightEstimate := unit.Pound(0)
// 	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
// 		PersonallyProcuredMove: models.PersonallyProcuredMove{
// 			MoveID:           moveID,
// 			Move:             move,
// 			PickupPostalCode: &pickupZip,
// 			OriginalMoveDate: &moveDate,
// 			WeightEstimate:   &weightEstimate,
// 		},
// 	})

// 	planner := &mocks.Planner{}
// 	planner.On("Zip5TransitDistanceLineHaul",
// 		"94540",
// 		"95632",
// 	).Return(3200, nil)
// 	calculator := NewEstimateCalculator(suite.DB(), planner)
// 	_, _, err := calculator.CalculateEstimates(&ppm, moveID, suite.logger)
// 	suite.Error(err)
// }

func (suite *PPMServiceSuite) setupCalculateEstimateTest(moveID uuid.UUID, originDutyLocationZip string, newDutyLocationZip string) models.Move {
	originZip3 := models.Tariff400ngZip3{
		Zip3:          "503",
		BasepointCity: "Des Moines",
		State:         "IA",
		ServiceArea:   "296",
		RateArea:      "US53",
		Region:        "7",
	}
	suite.MustSave(&originZip3)
	destinationZip3 := models.Tariff400ngZip3{
		Zip3:          "956",
		BasepointCity: "Sacramento",
		State:         "CA",
		ServiceArea:   "68",
		RateArea:      "US87",
		Region:        "2",
	}
	suite.MustSave(&destinationZip3)

	destinationZip5 := models.Tariff400ngZip5RateArea{
		Zip5:     "95630",
		RateArea: "US87",
	}
	suite.MustSave(&destinationZip5)

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Des Moines, IA",
		ServiceArea:        "296",
		LinehaulFactor:     unit.Cents(263),
		ServiceChargeCents: unit.Cents(489),
		ServicesSchedule:   3,
		EffectiveDateLower: scenario.May15TestYear,
		EffectiveDateUpper: scenario.May14FollowingYear,
		SIT185ARateCents:   unit.Cents(1447),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea)
	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Sacramento, CA",
		ServiceArea:        "68",
		LinehaulFactor:     unit.Cents(78),
		ServiceChargeCents: unit.Cents(452),
		ServicesSchedule:   3,
		EffectiveDateLower: scenario.May15TestYear,
		EffectiveDateUpper: scenario.May14FollowingYear,
		SIT185ARateCents:   unit.Cents(1642),
		SIT185BRateCents:   unit.Cents(70),
		SITPDSchedule:      3,
	}
	suite.MustSave(&destinationServiceArea)

	mySpecificRate := unit.Cents(474747)
	distanceLower := 3101
	distanceUpper := 3300
	weightLbsLower := unit.Pound(5000)
	weightLbsUpper := unit.Pound(10000)

	newBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: distanceLower,
		DistanceMilesUpper: distanceUpper,
		WeightLbsLower:     weightLbsLower,
		WeightLbsUpper:     weightLbsUpper,
		RateCents:          mySpecificRate,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.NonPeakRateCycleStart,
		EffectiveDateUpper: testdatagen.NonPeakRateCycleEnd,
	}
	suite.MustSave(&newBaseLinehaul)

	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "US53",
		DestinationRegion: "2",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl1)

	tdl2 := models.TrafficDistributionList{
		SourceRateArea:    "US87",
		DestinationRegion: "2",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl2)

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	suite.MustSave(&tsp)
	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          scenario.Oct1TestYear,
		PerformancePeriodEnd:            scenario.Dec31TestYear,
		RateCycleStart:                  scenario.Oct1TestYear,
		RateCycleEnd:                    scenario.May14FollowingYear,
		TrafficDistributionListID:       tdl1.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance)

	tspPerformance2 := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          scenario.Oct1TestYear,
		PerformancePeriodEnd:            scenario.Dec31TestYear,
		RateCycleStart:                  scenario.Oct1TestYear,
		RateCycleEnd:                    scenario.May14FollowingYear,
		TrafficDistributionListID:       tdl2.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance2)

	address1 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     originDutyLocationZip,
	}
	suite.MustSave(&address1)

	affiliationAirforce := internalmessages.AffiliationAIRFORCE
	locationName := "Origin Duty Location"
	originLocation := models.DutyLocation{
		Name:        locationName,
		Affiliation: &affiliationAirforce,
		AddressID:   address1.ID,
		Address:     address1,
	}
	suite.MustSave(&originLocation)

	address2 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     newDutyLocationZip,
	}
	suite.MustSave(&address2)
	locationName = "New Duty Location"
	newLocation := models.DutyLocation{
		Name:        locationName,
		Affiliation: &affiliationAirforce,
		AddressID:   address2.ID,
		Address:     address2,
	}
	suite.MustSave(&newLocation)

	orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			NewDutyLocationID: newLocation.ID,
			NewDutyLocation:   newLocation,
			ServiceMember: models.ServiceMember{
				DutyLocation:   originLocation,
				DutyLocationID: &originLocation.ID,
			},
		},
	})

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       moveID,
			OrdersID: orders.ID,
		},
		Order: orders,
	})

	return move
}
