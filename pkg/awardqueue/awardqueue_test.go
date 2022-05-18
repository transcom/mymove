//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package awardqueue

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AwardQueueSuite) Test_GetTSPsPerBandWithRemainder() {
	t := suite.T()
	// Check bands should expect differing num of TSPs when not divisible by 4
	// Remaining TSPs should be divided among bands in descending order
	tspPerBandList := getTSPsPerBand(10)
	expectedBandList := []int{3, 3, 2, 2}
	if !equalSlice(tspPerBandList, expectedBandList) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func (suite *AwardQueueSuite) Test_GetTSPsPerBandNoRemainder() {
	t := suite.T()
	// Check bands should expect correct num of TSPs when num of TSPs is divisible by 4
	tspPerBandList := getTSPsPerBand(8)
	expectedBandList := []int{2, 2, 2, 2}
	if !equalSlice(tspPerBandList, expectedBandList) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func (suite *AwardQueueSuite) Test_AssignTSPsToBands() {
	t := suite.T()
	queue := NewAwardQueue()
	tspsToMake := 5

	tdl := testdatagen.MakeDefaultTDL(suite.DB())

	var lastTSPP models.TransportationServiceProviderPerformance
	for i := 0; i < tspsToMake; i++ {
		tsp := testdatagen.MakeDefaultTSP(suite.DB())
		score := float64(mps + i + 1)
		rate := unit.NewDiscountRateFromPercent(45.3)

		var err error
		lastTSPP, err = testdatagen.MakeTSPPerformance(suite.DB(), testdatagen.Assertions{
			TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
				TransportationServiceProvider:   tsp,
				TransportationServiceProviderID: tsp.ID,
				TrafficDistributionListID:       tdl.ID,
				BestValueScore:                  score,
				LinehaulRate:                    rate,
				SITRate:                         rate,
			},
		})

		if err != nil {
			t.Errorf("Failed to MakeTSPPerformance: %v", err)
		}
	}

	err := queue.assignPerformanceBands(suite.AppContextForTest())

	if err != nil {
		t.Errorf("Failed to assign to performance bands: %v", err)
	}

	perfGroup := models.TSPPerformanceGroup{
		TrafficDistributionListID: lastTSPP.TrafficDistributionListID,
		PerformancePeriodStart:    lastTSPP.PerformancePeriodStart,
		PerformancePeriodEnd:      lastTSPP.PerformancePeriodEnd,
		RateCycleStart:            lastTSPP.RateCycleStart,
		RateCycleEnd:              lastTSPP.RateCycleEnd,
	}

	perfs, err := models.FetchTSPPerformancesForQualityBandAssignment(suite.DB(), perfGroup, mps)
	if err != nil {
		t.Errorf("Failed to fetch TSPPerformances: %v", err)
	}

	expectedBands := []int{1, 1, 2, 3, 4}

	for i, perf := range perfs {
		band := expectedBands[i]

		if perf.QualityBand == nil {
			t.Errorf("No quality band assigned for Performance #%v, got nil", perf.ID)
		} else if (*perf.QualityBand) != band {
			t.Errorf("Wrong quality band: expected %v, got %v", band, *perf.QualityBand)
		}
	}
}

func equalSlice(a []int, b []int) bool {
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

type AwardQueueSuite struct {
	testingsuite.PopTestSuite
}

func TestAwardQueueSuite(t *testing.T) {
	hs := &AwardQueueSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
