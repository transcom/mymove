package models_test

import (
	"context"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

var mps = 10.0

func (suite *ModelSuite) Test_PerformancePeriodValidations() {
	now := time.Now()
	earlier := now.AddDate(0, 0, -1)
	later := now.AddDate(0, 0, 1)

	tspPerformance := &TransportationServiceProviderPerformance{
		PerformancePeriodStart: later,
		PerformancePeriodEnd:   earlier,
	}

	var expErrors = map[string][]string{
		"performance_period_start": []string{"PerformancePeriodStart must be before PerformancePeriodEnd."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)
}

func (suite *ModelSuite) Test_RateCycleValidations() {
	now := time.Now()
	earlier := now.AddDate(0, 0, -1)
	later := now.AddDate(0, 0, 1)

	tspPerformance := &TransportationServiceProviderPerformance{
		RateCycleStart: later,
		RateCycleEnd:   earlier,
	}

	var expErrors = map[string][]string{
		"rate_cycle_start": []string{"RateCycleStart must be before RateCycleEnd."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)
}

func (suite *ModelSuite) Test_BestValueScoreValidations() {
	tspPerformance := &TransportationServiceProviderPerformance{BestValueScore: 101}

	var expErrors = map[string][]string{
		"best_value_score": {"101 is not less than 101."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)

	tspPerformance = &TransportationServiceProviderPerformance{BestValueScore: -1}

	expErrors = map[string][]string{
		"best_value_score": {"-1 is not greater than -1."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)
}

func (suite *ModelSuite) Test_GetRateCycle() {
	t := suite.T()

	peakRateCycleStart, peakRateCycleEnd := GetRateCycle(testdatagen.TestYear, true)
	nonPeakRateCycleStart, nonPeakRateCycleEnd := GetRateCycle(testdatagen.TestYear, false)

	if peakRateCycleStart != testdatagen.PeakRateCycleStart {
		t.Errorf("PeakRateCycleStart not calculated correctly. Expected %s, got %s",
			testdatagen.PeakRateCycleStart, peakRateCycleStart)
	}
	if peakRateCycleEnd != testdatagen.PeakRateCycleEnd {
		t.Errorf("PeakRateCycleEnd not calculated correctly. Expected %s, got %s",
			testdatagen.PeakRateCycleEnd, peakRateCycleEnd)
	}
	if nonPeakRateCycleStart != testdatagen.NonPeakRateCycleStart {
		t.Errorf("NonPeakRateCycleStart not calculated correctly. Expected %s, got %s",
			testdatagen.NonPeakRateCycleStart, nonPeakRateCycleStart)
	}
	if nonPeakRateCycleEnd != testdatagen.NonPeakRateCycleEnd {
		t.Errorf("NonPeakRateCycleEnd not calculated correctly. Expected %s, got %s",
			testdatagen.NonPeakRateCycleEnd, nonPeakRateCycleStart)
	}
}

func (suite *ModelSuite) Test_IncrementTSPPerformanceOfferCount() {
	t := suite.T()

	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	tsp := testdatagen.MakeDefaultTSP(suite.DB())
	perf, _ := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, nil, mps, 0, .2, .1)

	// Increment offer count twice
	performance, err := IncrementTSPPerformanceOfferCount(suite.DB(), perf.ID)
	if err != nil {
		t.Fatalf("Could not increment offer_count: %v", err)
	}
	if performance.OfferCount != 1 {
		t.Errorf("Wrong OfferCount returned: expected %d, got %d", 1, performance.OfferCount)
	}

	performance, err = IncrementTSPPerformanceOfferCount(suite.DB(), perf.ID)
	if err != nil {
		t.Fatalf("Could not increment offer_count: %v", err)
	}
	if performance.OfferCount != 2 {
		t.Errorf("Wrong OfferCount returned: expected %d, got %d", 2, performance.OfferCount)
	}

	performance = TransportationServiceProviderPerformance{}
	if err := suite.DB().Find(&performance, perf.ID); err != nil {
		t.Fatalf("could not find perf: %v", err)
	}

	if performance.OfferCount != 2 {
		t.Errorf("Wrong OfferCount: expected %d, got %d", 2, performance.OfferCount)
	}
}

func (suite *ModelSuite) Test_AssignQualityBandToTSPPerformance() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			CodeOfService: "2",
		},
	})
	tsp := testdatagen.MakeDefaultTSP(suite.DB())
	perf, _ := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, nil, mps, 0, .2, .3)
	band := 1

	err := AssignQualityBandToTSPPerformance(context.Background(), suite.DB(), band, perf.ID)
	if err != nil {
		t.Fatalf("Did not update quality band: %v", err)
	}

	performance := TransportationServiceProviderPerformance{}
	if err := suite.DB().Find(&performance, perf.ID); err != nil {
		t.Fatalf("could not find perf: %v", err)
	}

	if performance.QualityBand == nil {
		t.Errorf("No value for QualityBand: expected %v, got %v", band, performance.QualityBand)
	} else if *performance.QualityBand != band {
		t.Errorf("Wrong value for QualityBand: expected %d, got %d", band, *performance.QualityBand)
	}
}

func (suite *ModelSuite) Test_BVSWithLowMPS() {
	t := suite.T()
	tspsToMake := 5

	// Make a TDL to contain our tests
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			CodeOfService: "2",
		},
	})

	// Make 5 (not divisible by 4) TSPs in this TDL with BVSs above MPS threshold
	for i := 0; i < tspsToMake; i++ {
		tsp := testdatagen.MakeDefaultTSP(suite.DB())
		testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, nil, 15, 0, .5, .6)
	}
	// Make 1 TSP in this TDL with BVS below the MPS threshold
	mpsTSP := testdatagen.MakeDefaultTSP(suite.DB())
	lastTSPP, _ := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), mpsTSP, tdl, nil, mps-1, 0, .2, .9)

	perfGroup := TSPPerformanceGroup{
		TrafficDistributionListID: lastTSPP.TrafficDistributionListID,
		PerformancePeriodStart:    lastTSPP.PerformancePeriodStart,
		PerformancePeriodEnd:      lastTSPP.PerformancePeriodEnd,
		RateCycleStart:            lastTSPP.RateCycleStart,
		RateCycleEnd:              lastTSPP.RateCycleEnd,
	}

	// Fetch TSPs in TDL
	tspsbb, err := FetchTSPPerformancesForQualityBandAssignment(suite.DB(), perfGroup, mps)

	// Then: Expect to find TSPs in TDL
	if err != nil {
		t.Errorf("Failed to find TSPs: %v", err)
	}
	// Then: Expect TSP with low BVS won't be in sorted TSP slice
	for _, tsp := range tspsbb {
		if tsp.ID == mpsTSP.ID {
			t.Errorf("TSP: %v with a BVS below MPS incorrectly included.", mpsTSP.ID)
		}
	}
}

// Test_FetchNextQualityBandTSPPerformance ensures that the TSP with the highest BVS is returned in the expected band
func (suite *ModelSuite) Test_FetchNextQualityBandTSPPerformance() {
	t := suite.T()

	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	tsp1 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp2 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp3 := testdatagen.MakeDefaultTSP(suite.DB())

	// TSPs should be orderd by offer_count first, then BVS.
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp1, tdl, swag.Int(1), mps+1, 0, .4, .4)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp2, tdl, swag.Int(1), mps+3, 0, .4, .4)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp3, tdl, swag.Int(1), mps+2, 0, .4, .4)

	tspp, err := NextTSPPerformanceInQualityBand(suite.DB(), tdl.ID, 1, testdatagen.DateInsidePerformancePeriod,
		testdatagen.DateInsidePeakRateCycle)

	if err != nil {
		t.Errorf("Failed to find TSPPerformance: %v", err)
	} else if tspp.TransportationServiceProviderID != tsp2.ID {
		t.Errorf("TSPPerformance for wrong TSP returned: expected %s, got %s",
			tsp2.ID, tspp.TransportationServiceProviderID)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceAllZeros() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp1 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp1.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceOneAssigned() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 1, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp1.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceOneFullRound() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 2, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 1, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp1.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceTwoFullRounds() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 10, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 6, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 4, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 2, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp1.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceFirstBandFilled() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp2.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceThreeBands() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 10, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 2, QualityBand: swag.Int(3)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp2.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceHalfOffered() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp3.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformancePartialRound() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{OfferCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{OfferCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{OfferCount: 1, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{OfferCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp2 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp3.QualityBand, *chosen.QualityBand)
	}
}

// Test_GatherNextEligibleTSPPerformanceByBand ensures that TSPs are returned in the expected
// order for the Award Queue operation.
func (suite *ModelSuite) Test_GatherNextEligibleTSPPerformances() {
	t := suite.T()
	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	tsp1 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp2 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp3 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp4 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp5 := testdatagen.MakeDefaultTSP(suite.DB())
	// TSPs should be orderd by offer_count first, then BVS.
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp1, tdl, swag.Int(1), mps+5, 0, .4, .4)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp2, tdl, swag.Int(1), mps+4, 0, .3, .3)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp3, tdl, swag.Int(2), mps+3, 0, .2, .2)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp4, tdl, swag.Int(3), mps+2, 0, .1, .1)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp5, tdl, swag.Int(4), mps+1, 0, .1, .1)

	tsps, err := GatherNextEligibleTSPPerformances(suite.DB(), tdl.ID, testdatagen.DateInsidePerformancePeriod,
		testdatagen.DateInsidePeakRateCycle)
	expectedTSPorder := []uuid.UUID{tsp1.ID, tsp3.ID, tsp4.ID, tsp5.ID}

	actualTSPorder := []uuid.UUID{
		tsps[1].TransportationServiceProviderID,
		tsps[2].TransportationServiceProviderID,
		tsps[3].TransportationServiceProviderID,
		tsps[4].TransportationServiceProviderID}

	if err != nil {
		t.Errorf("Failed to find TSPPerformances: %v", err)
	} else if len(tsps) != 4 {
		t.Errorf("Found wrong number of TSPPerformances. Expected to find 4, found %d", len(tsps))
	} else if !equalUUIDSlice(expectedTSPorder, actualTSPorder) {
		t.Errorf("TSPs returned out of expected order.\n"+
			"\tExpected: %v \nFound: %v",
			expectedTSPorder,
			actualTSPorder)
	}
}

// Test_FetchTSPPerformancesForQualityBandAssignment ensures that TSPs are returned in the expected
// order for the division into quality bands.
func (suite *ModelSuite) Test_FetchTSPPerformancesForQualityBandAssignment() {
	t := suite.T()

	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	tsp1 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp2 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp3 := testdatagen.MakeDefaultTSP(suite.DB())
	// What matter is the BVS score order; offer count has no influence.
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp1, tdl, nil, 90, 0, .5, .5)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp2, tdl, nil, 50, 1, .3, .9)
	lastTSPP, _ := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp3, tdl, nil, 15, 1, .1, .3)

	perfGroup := TSPPerformanceGroup{
		TrafficDistributionListID: lastTSPP.TrafficDistributionListID,
		PerformancePeriodStart:    lastTSPP.PerformancePeriodStart,
		PerformancePeriodEnd:      lastTSPP.PerformancePeriodEnd,
		RateCycleStart:            lastTSPP.RateCycleStart,
		RateCycleEnd:              lastTSPP.RateCycleEnd,
	}

	tsps, err := FetchTSPPerformancesForQualityBandAssignment(suite.DB(), perfGroup, mps)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 3 {
		t.Errorf("Failed to find TSPs. Expected to find 3, found %d", len(tsps))
	} else if tsps[0].TransportationServiceProviderID != tsp1.ID ||
		tsps[1].TransportationServiceProviderID != tsp2.ID ||
		tsps[2].TransportationServiceProviderID != tsp3.ID {
		t.Errorf("\tExpected: [%s, %s, %s]\nFound: [%s, %s, %s]",
			tsp1.ID, tsp2.ID, tsp3.ID,
			tsps[0].TransportationServiceProviderID,
			tsps[1].TransportationServiceProviderID,
			tsps[2].TransportationServiceProviderID,
		)
	}
}

// Test_MinimumPerformanceScore ensures that TSPs whose BVS is below the MPS
// do not enter the Award Queue process.
func (suite *ModelSuite) Test_MinimumPerformanceScore() {
	t := suite.T()

	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	tsp1 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp2 := testdatagen.MakeDefaultTSP(suite.DB())
	// Make 2 TSPs, one with a BVS above the MPS and one below the MPS.
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp1, tdl, nil, mps+1, 0, .3, .4)
	lastTSPP, _ := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp2, tdl, nil, mps-1, 1, .9, .7)

	perfGroup := TSPPerformanceGroup{
		TrafficDistributionListID: lastTSPP.TrafficDistributionListID,
		PerformancePeriodStart:    lastTSPP.PerformancePeriodStart,
		PerformancePeriodEnd:      lastTSPP.PerformancePeriodEnd,
		RateCycleStart:            lastTSPP.RateCycleStart,
		RateCycleEnd:              lastTSPP.RateCycleEnd,
	}

	tsps, err := FetchTSPPerformancesForQualityBandAssignment(suite.DB(), perfGroup, mps)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 1 {
		t.Errorf("Failed to find TSPs. Expected to find 1, found %d", len(tsps))
	} else if tsps[0].TransportationServiceProviderID != tsp1.ID {
		t.Errorf("Incorrect TSP returned. Expected %s, received %s.",
			tsp1.ID,
			tsps[0].TransportationServiceProviderID)
	}
}

func (suite *ModelSuite) Test_FetchUnbandedTSPPerformanceGroups() {
	t := suite.T()

	foundTDL := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US14",
			DestinationRegion: "3",
			CodeOfService:     "2",
		},
	})
	foundTSP := testdatagen.MakeDefaultTSP(suite.DB())
	foundTSPP, err := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), foundTSP, foundTDL, nil, float64(mps+1), 0, .2, .3)
	if err != nil {
		t.Errorf("Failed to MakeTSPPerformance for found TSPP: %v", err)
	}
	foundPerfGroup := TransportationServiceProviderPerformance{
		TrafficDistributionListID: foundTSPP.TrafficDistributionListID,
		PerformancePeriodStart:    foundTSPP.PerformancePeriodStart,
		PerformancePeriodEnd:      foundTSPP.PerformancePeriodEnd,
		RateCycleStart:            foundTSPP.RateCycleStart,
		RateCycleEnd:              foundTSPP.RateCycleEnd,
	}

	notFoundTDL := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US14",
			DestinationRegion: "5",
			CodeOfService:     "2",
		},
	})
	notFoundTSP := testdatagen.MakeDefaultTSP(suite.DB())
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), notFoundTSP, notFoundTDL, swag.Int(1), float64(mps+1), 0, .4, .3)

	unenrolledTDL := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US14",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})

	unenrolledTSP := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: TransportationServiceProvider{
			Enrolled: false,
		},
	})
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), unenrolledTSP, unenrolledTDL, nil, float64(mps+1), 0, .4, .3)

	perfGroups, err := FetchUnbandedTSPPerformanceGroups(suite.DB())
	if err != nil {
		t.Fatal(err)
	}

	suite.Len(perfGroups, 1, "Got wrong number of TSPP groups; expected: 1, got: %d",
		len(perfGroups))

	suite.Equal(perfGroups[0].TrafficDistributionListID, foundPerfGroup.TrafficDistributionListID,
		"TrafficDistributionListID in TSPP group did not match")
	suite.True(perfGroups[0].PerformancePeriodStart.Equal(foundPerfGroup.PerformancePeriodStart),
		"PerformancePeriodStart in TSPP group did not match")
	suite.True(perfGroups[0].PerformancePeriodEnd.Equal(foundPerfGroup.PerformancePeriodEnd),
		"PerformancePeriodEnd in TSPP group did not match")
	suite.True(perfGroups[0].RateCycleStart.Equal(foundPerfGroup.RateCycleStart),
		"RateCycleStart in TSPP group did not match")
	suite.True(perfGroups[0].RateCycleEnd.Equal(foundPerfGroup.RateCycleEnd),
		"RateCycleEnd in TSPP group did not match")
}

// Test_FetchDiscountRates tests that the discount rate for the TSP with the best BVS
// for the specified channel and date is returned.
func (suite *ModelSuite) Test_FetchDiscountRatesBVS() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US68",
			DestinationRegion: "5",
			CodeOfService:     "2",
		},
	}) // Victoria, TX to Salina, KS
	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	suite.MustSave(&Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: "320", Region: "6"})
	suite.MustSave(&Tariff400ngZip3{Zip3: "674", Region: "5", BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: "320"})

	pp2Start := date(testdatagen.TestYear, time.August, 1)
	pp2End := date(testdatagen.TestYear, time.September, 30)

	moveDate := date(testdatagen.TestYear, time.May, 26)

	tspPerformance := TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50.0),
	}
	suite.MustSave(&tspPerformance)

	lowerTSPPerformance := TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  89,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(55.5),
		SITRate:                         unit.NewDiscountRateFromPercent(52.0),
	}
	suite.MustSave(&lowerTSPPerformance)

	otherPerformancePeriodTSPPerformance := TransportationServiceProviderPerformance{
		PerformancePeriodStart:          pp2Start,
		PerformancePeriodEnd:            pp2End,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  91,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(55.5),
		SITRate:                         unit.NewDiscountRateFromPercent(53.0),
	}
	suite.MustSave(&otherPerformancePeriodTSPPerformance)

	discountRate, sitRate, err := FetchDiscountRates(suite.DB(), "77901", "67401", "2", moveDate)
	if err != nil {
		t.Fatalf("Failed to find tsp performance: %s", err)
	}

	expectedLinehaul := unit.DiscountRate(.505)
	if discountRate != expectedLinehaul {
		t.Errorf("Wrong discount rate: expected %v, got %v", expectedLinehaul, discountRate)
	}

	expectedSIT := unit.DiscountRate(.5)
	if sitRate != expectedSIT {
		t.Errorf("Wrong discount rate: expected %v, got %v", expectedSIT, sitRate)
	}
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func (suite *ModelSuite) Test_FetchDiscountRatesPerformancePeriodBoundaries() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US68",
			DestinationRegion: "5",
			CodeOfService:     "2",
		},
	}) // Victoria, TX to Salina, KS
	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	suite.MustSave(&Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: "320", Region: "6"})
	suite.MustSave(&Tariff400ngZip3{Zip3: "674", Region: "5", BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: "320"})

	ppStart := testdatagen.PerformancePeriodStart
	ppEnd := testdatagen.PerformancePeriodEnd

	tspPerformance := TransportationServiceProviderPerformance{
		PerformancePeriodStart:          ppStart,
		PerformancePeriodEnd:            ppEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50.0),
	}
	suite.MustSave(&tspPerformance)

	if _, _, err := FetchDiscountRates(suite.DB(), "77901", "67401", "2", ppEnd); err != nil {
		t.Fatalf("Failed to find tsp performance for last day in performance period: %s", err)
	}

	if _, _, err := FetchDiscountRates(suite.DB(), "77901", "67401", "2", ppStart); err != nil {
		t.Fatalf("Failed to find tsp performance for first day in performance period: %s", err)
	}

	if _, _, err := FetchDiscountRates(suite.DB(), "77901", "67401", "2", ppStart.Add(time.Hour*-24)); err == nil {
		t.Fatalf("Should not have found a TSPP for the last day before the start of a performance period: %s", err)
	}

	if _, _, err := FetchDiscountRates(suite.DB(), "77901", "67401", "2", ppEnd.Add(time.Hour*24)); err == nil {
		t.Fatalf("Should not have found a TSPP for the first day following a performance period: %s", err)
	}
}

func equalUUIDSlice(a []uuid.UUID, b []uuid.UUID) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
