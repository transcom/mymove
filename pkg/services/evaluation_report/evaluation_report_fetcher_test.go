package evaluationreport

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *EvaluationReportSuite) TestFetchEvaluationReportList() {
	suite.Run("fetch for move with no evaluation reports should return empty array", func() {
		fetcher := NewEvaluationReportFetcher()
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
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
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 2)
	})
	suite.Run("reports submitted by other office users should be included", func() {
		fetcher := NewEvaluationReportFetcher()
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{
			factory.GetTraitOfficeUserEmail,
		})
		otherOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{
			factory.GetTraitOfficeUserEmail,
		})
		report := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    otherOfficeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 1)
		suite.Equal(report.ID, reports[0].ID)
	})
	suite.Run("draft reports by other office users should not be included", func() {
		fetcher := NewEvaluationReportFetcher()
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{
			factory.GetTraitOfficeUserEmail,
		})
		otherOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{
			factory.GetTraitOfficeUserEmail,
		})
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    otherOfficeUser,
				LinkOnly: true,
			},
		}, nil)
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Empty(reports)
	})
	suite.Run("fetch counseling reports should only return counseling reports", func() {
		fetcher := NewEvaluationReportFetcher()
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		counselingReport := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)

		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type:        models.EvaluationReportTypeShipment,
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		reports, err := fetcher.FetchEvaluationReports(suite.AppContextForTest(), models.EvaluationReportTypeCounseling, move.ID, officeUser.ID)
		suite.NoError(err)
		suite.Len(reports, 1)
		suite.Equal(counselingReport.ID, reports[0].ID)
	})
	suite.Run("fetch shipment reports should only return shipment reports", func() {
		fetcher := NewEvaluationReportFetcher()
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)

		shipmentReport := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type:        models.EvaluationReportTypeShipment,
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
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
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		report := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		fetchedReport, err := fetcher.FetchEvaluationReportByID(suite.AppContextForTest(), report.ID, officeUser.ID)
		suite.NoError(err)
		suite.Equal(report.ID, fetchedReport.ID)
		suite.NotNil(report.Move.ReferenceID)
	})
	// forbidden if they don't own the draft
	suite.Run("fetch for a draft evaluation report should return a forbidden if the requester isn't the owner", func() {
		fetcher := NewEvaluationReportFetcher()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		officeUserOwner := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		report := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUserOwner,
				LinkOnly: true,
			},
		}, nil)
		fetchedReport, err := fetcher.FetchEvaluationReportByID(suite.AppContextForTest(), report.ID, officeUser.ID)
		suite.Nil(fetchedReport)
		suite.Error(err, apperror.NewForbiddenError("Draft evaluation reports are viewable only by their owner/creator."))
	})
	// not found error if the ID is wrong
	suite.Run("fetch should return a not found error if the reportID doesn't exist", func() {
		fetcher := NewEvaluationReportFetcher()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		wrongID, _ := uuid.NewV4()
		expectedError := apperror.NewNotFoundError(wrongID, "while looking for evaluation report")

		_, err := fetcher.FetchEvaluationReportByID(suite.AppContextForTest(), wrongID, officeUser.ID)
		suite.Equal(expectedError, err)
	})
}
