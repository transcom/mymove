package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_TrafficDistributionList() {
	tdl := &TrafficDistributionList{}

	expErrors := map[string][]string{
		"source_rate_area":   {"SourceRateArea can not be blank.", "SourceRateArea does not match the expected format."},
		"destination_region": {"DestinationRegion can not be blank.", "DestinationRegion does not match the expected format."},
		"code_of_service":    {"CodeOfService can not be blank."},
	}

	suite.verifyValidationErrors(tdl, expErrors)
}

func (suite *ModelSuite) Test_FetchTDL() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US28",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})
	fetchedTDL, err := FetchTDL(suite.DB(), "US28", "4", "2")
	if err != nil {
		t.Errorf("Something went wrong fetching the test object: %s\n", err)
	}

	if fetchedTDL.ID != tdl.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", tdl.ID, fetchedTDL.ID)
	}
}

func (suite *ModelSuite) Test_FetchOrCreateTDL() {
	t := suite.T()

	foundTDL := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US28",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})
	foundTSP := testdatagen.MakeDefaultTSP(suite.DB())
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
	//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
	//RA: in which this would be considered a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	testdatagen.MakeTSPPerformance(suite.DB(), testdatagen.Assertions{
		TransportationServiceProviderPerformance: TransportationServiceProviderPerformance{
			TransportationServiceProvider:   foundTSP,
			TransportationServiceProviderID: foundTSP.ID,
			TrafficDistributionListID:       foundTDL.ID,
		},
	})

	fetchedTDL, err := FetchOrCreateTDL(suite.DB(), "US28", "4", "2")
	if err != nil {
		t.Errorf("Didn't return expected TDL: %v", fetchedTDL)
	}

	if fetchedTDL.ID != foundTDL.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", foundTDL.ID, fetchedTDL.ID)
	}

	_, err = FetchOrCreateTDL(suite.DB(), "US23", "1", "2")
	if err != nil {
		t.Errorf("Something went wrong creating the test objects: %s\n", err)
	}

	tdls := TrafficDistributionLists{}
	sql := `SELECT
			*
		FROM
			traffic_distribution_lists;`

	err = suite.DB().RawQuery(sql).All(&tdls)

	if err != nil {
		t.Errorf("Something went wrong fetching the test objects: %s\n", err)
	}
	if len(tdls) != 2 {
		t.Errorf("A new object was not created")
	}
}

func (suite *ModelSuite) Test_TDLUniqueChannelCOS() {
	t := suite.T()

	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: TrafficDistributionList{
			SourceRateArea:    "US28",
			DestinationRegion: "4",
			CodeOfService:     "2",
		},
	})
	fetchedTDL, err := FetchOrCreateTDL(suite.DB(), "US28", "4", "2")
	if err != nil {
		t.Errorf("Something went wrong fetching the test object: %s\n", err)
	}

	if fetchedTDL.ID != tdl.ID {
		t.Errorf("Got wrong TDL; expected: %s, got: %s", tdl.ID, fetchedTDL.ID)
	}
}
