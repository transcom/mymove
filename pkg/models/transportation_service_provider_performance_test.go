package models_test

import (
	"github.com/go-openapi/swag"
	"github.com/satori/go.uuid"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var mps = 10

func (suite *ModelSuite) Test_BestValueScoreValidations() {
	tspPerformance := &TransportationServiceProviderPerformance{BestValueScore: 101}

	var expErrors = map[string][]string{
		"best_value_score": []string{"101 is not less than 101."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)

	tspPerformance = &TransportationServiceProviderPerformance{BestValueScore: -1}

	expErrors = map[string][]string{
		"best_value_score": []string{"-1 is not greater than -1."},
	}

	suite.verifyValidationErrors(tspPerformance, expErrors)
}

func (suite *ModelSuite) Test_IncrementTSPPerformanceAwardCount() {
	t := suite.T()

	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", "TEST")
	perf, _ := testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, nil, mps, 0)

	err := IncrementTSPPerformanceAwardCount(suite.db, perf.ID)
	if err != nil {
		t.Fatalf("Could not increament award_count: %v", err)
	}

	performance := TransportationServiceProviderPerformance{}
	if err := suite.db.Find(&performance, perf.ID); err != nil {
		t.Fatalf("could not find perf: %v", err)
	}

	if performance.AwardCount != 1 {
		t.Errorf("Wrong AwardCount: expected %d, got %d", 1, performance.AwardCount)
	}
}

func (suite *ModelSuite) Test_AssignQualityBandToTSPPerformance() {
	t := suite.T()

	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", "TEST")
	perf, _ := testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, nil, mps, 0)
	band := 1

	err := AssignQualityBandToTSPPerformance(suite.db, band, perf.ID)
	if err != nil {
		t.Fatalf("Did not update quality band: %v", err)
	}

	performance := TransportationServiceProviderPerformance{}
	if err := suite.db.Find(&performance, perf.ID); err != nil {
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
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")

	// Make 5 (not divisible by 4) TSPs in this TDL with BVSs above MPS threshold
	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", "TEST")
		testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, nil, 15, 0)
	}
	// Make 1 TSP in this TDL with BVS below the MPS threshold
	mpsTSP, _ := testdatagen.MakeTSP(suite.db, "Low BVS Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(suite.db, mpsTSP, tdl, nil, mps-1, 0)

	// Fetch TSPs in TDL
	tspsbb, err := FetchTSPPerformanceForQualityBandAssignment(suite.db, tdl.ID, mps)

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
	tdl, _ := testdatagen.MakeTDL(suite.db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(suite.db, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(suite.db, "Test TSP 3", "TSP2")
	// TSPs should be orderd by award_count first, then BVS.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, swag.Int(1), mps+1, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, swag.Int(1), mps+3, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp3, tdl, swag.Int(1), mps+2, 0)

	tsp, err := NextTSPPerformanceInQualityBand(suite.db, tdl.ID, 1)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if tsp.TransportationServiceProviderID != tsp2.ID {
		t.Errorf("Incorrect TSP returned.\n"+
			"\tExpected: %s \nFound: %s",
			tsp2.ID,
			tsp.TransportationServiceProviderID)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceAllZeros() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(4)}

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
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 1, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(4)}

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

func (suite *ModelSuite) Test_SelectNextTSPPerformanceOneFullRound() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 2, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 1, QualityBand: swag.Int(4)}

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

func (suite *ModelSuite) Test_SelectNextTSPPerformanceTwoFullRounds() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 10, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 6, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 4, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 2, QualityBand: swag.Int(4)}

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

func (suite *ModelSuite) Test_SelectNextTSPPerformanceFirstBandFilled() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(4)}

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
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 2, QualityBand: swag.Int(3)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp1 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp1.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformanceHalfAwarded() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp3 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp3.QualityBand, *chosen.QualityBand)
	}
}

func (suite *ModelSuite) Test_SelectNextTSPPerformancePartialRound() {
	t := suite.T()
	tspp1 := TransportationServiceProviderPerformance{AwardCount: 5, QualityBand: swag.Int(1)}
	tspp2 := TransportationServiceProviderPerformance{AwardCount: 3, QualityBand: swag.Int(2)}
	tspp3 := TransportationServiceProviderPerformance{AwardCount: 1, QualityBand: swag.Int(3)}
	tspp4 := TransportationServiceProviderPerformance{AwardCount: 0, QualityBand: swag.Int(4)}

	choices := map[int]TransportationServiceProviderPerformance{
		1: tspp1,
		2: tspp2,
		3: tspp3,
		4: tspp4}

	chosen := SelectNextTSPPerformance(choices)

	if chosen != tspp3 {
		t.Errorf("Wrong TSPPerformance selected: expected band %v, got %v", *tspp3.QualityBand, *chosen.QualityBand)
	}
}

// Test_GatherNextEligibleTSPPerformanceByBand ensures that TSPs are returned in the expected
// order for the Award Queue operation.
func (suite *ModelSuite) Test_GatherNextEligibleTSPPerformances() {
	t := suite.T()
	tdl, _ := testdatagen.MakeTDL(suite.db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(suite.db, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(suite.db, "Test TSP 3", "TSP3")
	tsp4, _ := testdatagen.MakeTSP(suite.db, "Test TSP 4", "TSP4")
	tsp5, _ := testdatagen.MakeTSP(suite.db, "Test TSP 5", "TSP5")
	// TSPs should be orderd by award_count first, then BVS.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, swag.Int(1), mps+5, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, swag.Int(1), mps+4, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp3, tdl, swag.Int(2), mps+3, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp4, tdl, swag.Int(3), mps+2, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp5, tdl, swag.Int(4), mps+1, 0)

	tsps, err := GatherNextEligibleTSPPerformances(suite.db, tdl.ID)
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

// Test_FetchTSPPerformanceForQualityBandAssignment ensures that TSPs are returned in the expected
// order for the division into quality bands.
func (suite *ModelSuite) Test_FetchTSPPerformanceForQualityBandAssignment() {
	t := suite.T()

	tdl, _ := testdatagen.MakeTDL(suite.db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(suite.db, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(suite.db, "Test TSP 3", "TSP2")
	// What matter is the BVS score order; award_count has no influence.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, nil, 90, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, nil, 50, 1)
	testdatagen.MakeTSPPerformance(suite.db, tsp3, tdl, nil, 15, 1)

	tsps, err := FetchTSPPerformanceForQualityBandAssignment(suite.db, tdl.ID, mps)

	if err != nil {
		t.Errorf("Failed to find TSP: %v", err)
	} else if len(tsps) != 3 {
		t.Errorf("Failed to find TSPs. Expected to find 3, found %d", len(tsps))
	} else if tsps[0].TransportationServiceProviderID != tsp1.ID &&
		tsps[1].TransportationServiceProviderID != tsp2.ID &&
		tsps[2].TransportationServiceProviderID != tsp3.ID {

		t.Errorf("TSPs returned out of expected order.\n"+
			"\tExpected: [%s, %s, %s]\nFound:    [%s, %s, %s]",
			tsp1.ID, tsp2.ID, tsp3.ID,
			tsps[0].TransportationServiceProviderID,
			tsps[1].TransportationServiceProviderID,
			tsps[2].TransportationServiceProviderID)
	}
}

// Test_MinimumPerformanceScore ensures that TSPs whose BVS is below the MPS
// do not enter the Award Queue process.
func (suite *ModelSuite) Test_MinimumPerformanceScore() {
	t := suite.T()

	tdl, _ := testdatagen.MakeTDL(suite.db, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(suite.db, "Test TSP 2", "TSP2")
	// Make 2 TSPs, one with a BVS above the MPS and one below the MPS.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, nil, mps+1, 0)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, nil, mps-1, 1)

	tsps, err := FetchTSPPerformanceForQualityBandAssignment(suite.db, tdl.ID, mps)

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
