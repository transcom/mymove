package evaluationreport

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EvaluationReportSuite) TestFetchEvaluationReportList() {
	suite.Run("fetch for move with no evaluation reports should return empty array", func() {
		fetcher := NewEvaluationReportListFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("fetch for nonexistent move or office user should not error", func() {
		// if we want to detect this error, it will take another query to the moves table
		// would that be worth it?
		fetcher := NewEvaluationReportListFetcher()
		badMoveID := uuid.Must(uuid.NewV4())
		badOfficeUserID := uuid.Must(uuid.NewV4())
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), badMoveID, badOfficeUserID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("submitted and draft reports for current user should be included", func() {
		fetcher := NewEvaluationReportListFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
			},
		})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
			},
		})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 2)
	})
	suite.Run("reports submitted by other office users should be included", func() {
		fetcher := NewEvaluationReportListFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		otherOfficeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: otherOfficeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
			},
		})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 1)
		suite.Equal(report.ID, reports[0].ID)
	})
	suite.Run("draft reports by other office users should not be included", func() {
		fetcher := NewEvaluationReportListFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		otherOfficeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: otherOfficeUser.ID,
				SubmittedAt:  nil,
			},
		})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("deleted reports should not be included", func() {
		fetcher := NewEvaluationReportListFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				DeletedAt:    swag.Time(time.Now()),
			},
		})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
				DeletedAt:    swag.Time(time.Now()),
			},
		})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
}
