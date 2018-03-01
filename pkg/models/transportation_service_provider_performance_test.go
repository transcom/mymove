package models_test

import (
	"testing"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var mps = 10

func Test_BestValueScoreValidations(t *testing.T) {
	tspPerformance := &TransportationServiceProviderPerformance{BestValueScore: 101}

	var expErrors = map[string][]string{
		"best_value_score": []string{"101 is not less than 101."},
	}

	verifyValidationErrors(tspPerformance, expErrors, t)

	tspPerformance = &TransportationServiceProviderPerformance{BestValueScore: -1}

	expErrors = map[string][]string{
		"best_value_score": []string{"-1 is not greater than -1."},
	}

	verifyValidationErrors(tspPerformance, expErrors, t)
}

func Test_BVSWithLowMPS(t *testing.T) {
	tspsToMake := 5

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(dbConnection, "california", "90210", "2")

	// Make 5 (not divisible by 4) TSPs in this TDL with BVSs above MPS threshold
	for i := 0; i < tspsToMake; i++ {
		tsp, _ := testdatagen.MakeTSP(dbConnection, "Test Shipper", "TEST")
		testdatagen.MakeTSPPerformance(dbConnection, tsp, tdl, nil, 15, 0)
	}
	// Make 1 TSP in this TDL with BVS below the MPS threshold
	mpsTSP, _ := testdatagen.MakeTSP(dbConnection, "Low BVS Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(dbConnection, mpsTSP, tdl, nil, mps-1, 0)

	// Fetch TSPs in TDL
	tspsbb, err := FetchTSPPerformanceForQualityBandAssignment(dbConnection, tdl.ID, mps)

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

// Test_FetchTSPPerformanceForAwardQueue ensures that TSPs are returned in the expected
// order for the Award Queue operation.
func Test_FetchTSPPerformanceForAwardQueue(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(dbConnection, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 3", "TSP2")
	// TSPs should be orderd by award_count first, then BVS.
	testdatagen.MakeTSPPerformance(dbConnection, tsp1, tdl, nil, mps+1, 0)
	testdatagen.MakeTSPPerformance(dbConnection, tsp2, tdl, nil, mps+3, 1)
	testdatagen.MakeTSPPerformance(dbConnection, tsp3, tdl, nil, mps+2, 1)

	tsps, err := FetchTSPPerformanceForAwardQueue(dbConnection, tdl.ID, mps)

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

// Test_FetchTSPPerformanceForQualityBandAssignment ensures that TSPs are returned in the expected
// order for the division into quality bands.
func Test_FetchTSPPerformanceForQualityBandAssignment(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(dbConnection, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 2", "TSP2")
	tsp3, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 3", "TSP2")
	// What matter is the BVS score order; award_count has no influence.
	testdatagen.MakeTSPPerformance(dbConnection, tsp1, tdl, nil, 90, 0)
	testdatagen.MakeTSPPerformance(dbConnection, tsp2, tdl, nil, 50, 1)
	testdatagen.MakeTSPPerformance(dbConnection, tsp3, tdl, nil, 15, 1)

	tsps, err := FetchTSPPerformanceForQualityBandAssignment(dbConnection, tdl.ID, mps)

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
func Test_MinimumPerformanceScore(t *testing.T) {
	tdl, _ := testdatagen.MakeTDL(dbConnection, "source", "dest", "cos")
	tsp1, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 1", "TSP1")
	tsp2, _ := testdatagen.MakeTSP(dbConnection, "Test TSP 2", "TSP2")
	// Make 2 TSPs, one with a BVS above the MPS and one below the MPS.
	testdatagen.MakeTSPPerformance(dbConnection, tsp1, tdl, nil, mps+1, 0)
	testdatagen.MakeTSPPerformance(dbConnection, tsp2, tdl, nil, mps-1, 1)

	tsps, err := FetchTSPPerformanceForQualityBandAssignment(dbConnection, tdl.ID, mps)

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
