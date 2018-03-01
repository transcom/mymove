package models_test

import (
	"testing"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func Test_TrafficDistributionList(t *testing.T) {
	tdl := &TrafficDistributionList{}

	var expErrors = map[string][]string{
		"source_rate_area":   []string{"SourceRateArea can not be blank."},
		"destination_region": []string{"DestinationRegion can not be blank."},
		"code_of_service":    []string{"CodeOfService can not be blank."},
	}

	verifyValidationErrors(tdl, expErrors, t)
}

func Test_FetchTDLsAwaitingBandAssignment(t *testing.T) {
	foundTDL, _ := testdatagen.MakeTDL(dbConnection, "california", "90210", "2")
	foundTSP, _ := testdatagen.MakeTSP(dbConnection, "Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(dbConnection, foundTSP, foundTDL, nil, mps+1, 0)

	notFoundTDL, _ := testdatagen.MakeTDL(dbConnection, "california", "90210", "2")
	notFoundTSP, _ := testdatagen.MakeTSP(dbConnection, "Test Shipper", "TEST")
	testdatagen.MakeTSPPerformance(dbConnection, notFoundTSP, notFoundTDL, intPtr(1), mps+1, 0)

	tdls, err := FetchTDLsAwaitingBandAssignment(dbConnection)
	if err != nil {
		t.Fatal(err)
	}

	if len(tdls) != 1 {
		t.Errorf("Got wrong number of TDLs; expected: 1, got: %d", len(tdls))
	}
}

func intPtr(i int) *int {
	return &i
}
