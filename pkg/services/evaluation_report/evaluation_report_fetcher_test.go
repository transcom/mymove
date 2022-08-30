package evaluationreport

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EvaluationReportSuite) TestFetchEvaluationReportList() {
	suite.Run("fetch for move with no evaluation reports should return empty array", func() {
		fetcher := NewEvaluationReportFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("fetch for nonexistent move or office user should not error", func() {
		// Since we're just checking if IDs in the evaluation reports match the provided IDs, and not
		// touching the moves or office users, we should get an empty response instead of an error.
		fetcher := NewEvaluationReportFetcher()
		badMoveID := uuid.Must(uuid.NewV4())
		badOfficeUserID := uuid.Must(uuid.NewV4())
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, badMoveID, badOfficeUserID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("submitted and draft reports for current user should be included", func() {
		fetcher := NewEvaluationReportFetcher()
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
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 2)
	})
	suite.Run("reports submitted by other office users should be included", func() {
		fetcher := NewEvaluationReportFetcher()
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
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 1)
		suite.Equal(report.ID, reports[0].ID)
	})
	suite.Run("draft reports by other office users should not be included", func() {
		fetcher := NewEvaluationReportFetcher()
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
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("deleted reports should not be included", func() {
		fetcher := NewEvaluationReportFetcher()
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
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("fetch counseling reports should only return counseling reports", func() {
		fetcher := NewEvaluationReportFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		counselingReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{Move: move})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
			},
			MTOShipment: shipment,
		})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 1)
		suite.Equal(counselingReport.ID, reports[0].ID)
	})
	suite.Run("fetch shipment reports should only return shipment reports", func() {
		fetcher := NewEvaluationReportFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{Move: move})
		shipmentReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				MoveID:       move.ID,
				OfficeUserID: officeUser.ID,
				SubmittedAt:  swag.Time(time.Now()),
			},
			MTOShipment: shipment,
		})
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeShipment, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 1)
		suite.Equal(shipmentReport.ID, reports[0].ID)
	})
}

func (suite *EvaluationReportSuite) TestFetchEvaluationReportByID() {
	// successful fetch
	suite.Run("fetch for a submitted evaluation report that exists should be successful", func() {
		fetcher := NewEvaluationReportFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeQAECSROfficeUser(suite.DB(), testdatagen.Assertions{})
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
			OfficeUserID: officeUser.ID,
			MoveID:       move.ID,
			SubmittedAt:  swag.Time(time.Now()),
		}})

		fetchedReport, err := fetcher.FetchEvaluationReportByID(suite.AppContextForTest(), report.ID, officeUser.ID)
		suite.NoError(err)
		suite.Equal(report.ID, fetchedReport.ID)
		suite.NotNil(report.Move.ReferenceID)
	})
	// forbidden if they don't own the draft
	suite.Run("fetch for a draft evaluation report should return a forbidden if the requester isn't the owner", func() {
		fetcher := NewEvaluationReportFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeQAECSROfficeUser(suite.DB(), testdatagen.Assertions{})
		officeUserOwner := testdatagen.MakeQAECSROfficeUser(suite.DB(), testdatagen.Assertions{})
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
			OfficeUserID: officeUserOwner.ID,
			MoveID:       move.ID,
			SubmittedAt:  nil,
		}})
		fetchedReport, err := fetcher.FetchEvaluationReportByID(suite.AppContextForTest(), report.ID, officeUser.ID)
		suite.Nil(fetchedReport)
		suite.Error(err, apperror.NewForbiddenError("Draft evaluation reports are viewable only by their owner/creator."))
	})
	// not found error if the ID is wrong
	suite.Run("fetch should return a not found error if the reportID doesn't exist", func() {
		fetcher := NewEvaluationReportFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeQAECSROfficeUser(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
			OfficeUserID: officeUser.ID,
			MoveID:       move.ID,
			SubmittedAt:  swag.Time(time.Now()),
		}})
		wrongID, _ := uuid.NewV4()
		expectedError := apperror.NewNotFoundError(wrongID, "while looking for evaluation report")

		_, err := fetcher.FetchEvaluationReportByID(suite.AppContextForTest(), wrongID, officeUser.ID)
		suite.Equal(expectedError, err)
	})
}
