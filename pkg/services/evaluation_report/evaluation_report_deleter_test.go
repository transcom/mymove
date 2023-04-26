package evaluationreport

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *EvaluationReportSuite) TestEvaluationReportDeleter() {
	setupTestData := func() (services.EvaluationReportDeleter, models.EvaluationReport, appcontext.AppContext) {
		deleter := NewEvaluationReportDeleter()
		report := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model: models.EvaluationReport{
					Type: models.EvaluationReportTypeShipment,
				},
			},
		}, nil)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: report.OfficeUserID,
		})

		return deleter, report, appCtx
	}

	suite.Run("delete existing report", func() {
		deleter, report, appCtx := setupTestData()

		suite.NoError(deleter.DeleteEvaluationReport(appCtx, report.ID))

		// Should not be able to find the deleted report
		var dbReport models.EvaluationReport
		err := suite.DB().Find(&dbReport, report.ID)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err)
	})

	suite.Run("Returns an error when delete non-existent report", func() {
		deleter, _, appCtx := setupTestData()

		uuid := uuid.Must(uuid.NewV4())
		err := deleter.DeleteEvaluationReport(appCtx, uuid)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)

	})
}
