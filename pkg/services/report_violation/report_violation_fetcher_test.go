package reportviolation

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ReportViolationSuite) TestFetchReportViolationsByReportID() {
	suite.Run("fetch for report with no violations should return empty array", func() {
		fetcher := NewReportViolationFetcher()

		badReportID := uuid.Must(uuid.NewV4())
		fetchedReportViolations, err := fetcher.FetchReportViolationsByReportID(suite.AppContextForTest(), badReportID)

		suite.NoError(err)
		suite.Empty(fetchedReportViolations)
	})
	suite.Run("fetch by reportId when there are report-violations for the report should be successful", func() {
		fetcher := NewReportViolationFetcher()
		usprc, err := models.FindByZipCodeAndCity(suite.DB(), "90210", "Beverly Hills")
		suite.NoError(err)
		reportViolation := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				UsPostRegionCityID: &usprc.ID,
			},
		})

		fetchedReportViolations, err := fetcher.FetchReportViolationsByReportID(suite.AppContextForTest(), reportViolation.ReportID)

		suite.NoError(err)
		suite.Equal(1, len(fetchedReportViolations))
		suite.Equal(reportViolation.ID, fetchedReportViolations[0].ID)
		suite.Equal(reportViolation.ViolationID, fetchedReportViolations[0].Violation.ID)
	})
}
