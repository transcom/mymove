package models_test

import (
	"github.com/go-openapi/swag"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_TrafficDistributionList() {
	tdl := &TrafficDistributionList{}

	expErrors := map[string][]string{
		"source_rate_area":   []string{"SourceRateArea can not be blank."},
		"destination_region": []string{"DestinationRegion can not be blank."},
		"code_of_service":    []string{"CodeOfService can not be blank."},
	}

	suite.verifyValidationErrors(tdl, expErrors)
}

func (suite *ModelSuite) Test_FetchTDLsAwaitingBandAssignment() {
	t := suite.T()

	foundTDL, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	foundTSP, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", testdatagen.RandomSCAC())
	testdatagen.MakeTSPPerformance(suite.db, foundTSP, foundTDL, nil, mps+1, 0)

	notFoundTDL, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	notFoundTSP, _ := testdatagen.MakeTSP(suite.db, "Test Shipper", testdatagen.RandomSCAC())
	testdatagen.MakeTSPPerformance(suite.db, notFoundTSP, notFoundTDL, swag.Int(1), mps+1, 0)

	tdls, err := FetchTDLsAwaitingBandAssignment(suite.db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tdls) != 1 {
		t.Errorf("Got wrong number of TDLs; expected: 1, got: %d", len(tdls))
	}

	if tdls[0].ID != foundTDL.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", foundTDL.ID, tdls[0].ID)
	}
}
