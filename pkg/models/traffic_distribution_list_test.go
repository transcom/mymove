package models_test

import (
	"github.com/go-openapi/swag"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_TrafficDistributionList() {
	tdl := &TrafficDistributionList{}

	expErrors := map[string][]string{
		"source_rate_area":   []string{"SourceRateArea can not be blank.", "SourceRateArea does not match the expected format."},
		"destination_region": []string{"DestinationRegion can not be blank.", "DestinationRegion does not match the expected format."},
		"code_of_service":    []string{"CodeOfService can not be blank."},
	}

	suite.verifyValidationErrors(tdl, expErrors)
}

func (suite *ModelSuite) Test_FetchTDLsAwaitingBandAssignment() {
	t := suite.T()

	foundTDL := testdatagen.MakeTDL(suite.db, testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US14",
			DestinationRegion: "3",
			CodeOfService:     "2",
		},
	})
	foundTSP := testdatagen.MakeDefaultTSP(suite.db)
	testdatagen.MakeTSPPerformanceDeprecated(suite.db, foundTSP, foundTDL, nil, float64(mps+1), 0, .2, .3)

	notFoundTDL := testdatagen.MakeTDL(suite.db, testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US14",
			DestinationRegion: "5",
			CodeOfService:     "2",
		},
	})
	notFoundTSP := testdatagen.MakeDefaultTSP(suite.db)
	testdatagen.MakeTSPPerformanceDeprecated(suite.db, notFoundTSP, notFoundTDL, swag.Int(1), float64(mps+1), 0, .4, .3)

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

func (suite *ModelSuite) Test_FetchTDL() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.db, testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US28",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})
	fetchedTDL, err := FetchTDL(suite.db, "US28", "4", "2")
	if err != nil {
		t.Errorf("Something went wrong fetching the test object: %s\n", err)
	}

	if fetchedTDL.ID != tdl.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", tdl.ID, fetchedTDL.ID)
	}
}

func (suite *ModelSuite) Test_FetchOrCreateTDL() {
	t := suite.T()

	foundTDL := testdatagen.MakeTDL(suite.db, testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US28",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})
	foundTSP := testdatagen.MakeDefaultTSP(suite.db)
	testdatagen.MakeTSPPerformanceDeprecated(suite.db, foundTSP, foundTDL, swag.Int(1), float64(mps+1), 0, .2, .3)

	fetchedTDL, err := FetchOrCreateTDL(suite.db, "US28", "4", "2")
	if err != nil {
		t.Errorf("Didn't return expected TDL: %v", fetchedTDL)
	}

	if fetchedTDL.ID != foundTDL.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", foundTDL.ID, fetchedTDL.ID)
	}

	_, err = FetchOrCreateTDL(suite.db, "US23", "1", "2")
	if err != nil {
		t.Errorf("Something went wrong creating the test objects: %s\n", err)
	}

	tdls := TrafficDistributionLists{}
	sql := `SELECT
			*
		FROM
			traffic_distribution_lists;`

	err = suite.db.RawQuery(sql).All(&tdls)

	if err != nil {
		t.Errorf("Something went wrong fetching the test objects: %s\n", err)
	}
	if len(tdls) != 2 {
		t.Errorf("A new object was not created")
	}
}

func (suite *ModelSuite) Test_TDLUniqueChannelCOS() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.db, testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US28",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})
	fetchedTDL, err := FetchOrCreateTDL(suite.db, "US28", "4", "2")
	if err != nil {
		t.Errorf("Something went wrong fetching the test object: %s\n", err)
	}

	if fetchedTDL.ID != tdl.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", tdl.ID, fetchedTDL.ID)
	}
}
