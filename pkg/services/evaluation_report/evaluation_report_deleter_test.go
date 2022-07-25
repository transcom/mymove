package evaluationreport

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EvaluationReportSuite) TestEvaluationReportDeleter() {
	setupTestData := func() (services.EvaluationReportDeleter, models.EvaluationReport, appcontext.AppContext) {
		deleter := NewEvaluationReportDeleter()
		move := testdatagen.MakeDefaultMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: move.ID,
			},
		})
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				ShipmentID: &shipment.ID,
				Type:       models.EvaluationReportTypeShipment, OfficeUserID: officeUser.ID,
			},
		})
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: officeUser.ID,
		})

		return deleter, report, appCtx
	}

	suite.Run("delete existing report", func() {
		deleter, report, appCtx := setupTestData()

		suite.NoError(deleter.DeleteEvaluationReport(appCtx, report.ID))
		var dbReport models.EvaluationReport
		err := suite.DB().Find(&dbReport, report.ID)
		suite.NoError(err)
		suite.NotNil(dbReport.DeletedAt)
	})

	suite.Run("Returns an error when delete non-existent report", func() {
		deleter, _, appCtx := setupTestData()

		uuid := uuid.Must(uuid.NewV4())
		err := deleter.DeleteEvaluationReport(appCtx, uuid)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)

	})

	suite.Run("Returns an error when attempting to delete an already deleted report", func() {
		deleter, report, appCtx := setupTestData()

		suite.NoError(deleter.DeleteEvaluationReport(appCtx, report.ID))
		var dbReport models.EvaluationReport
		err := suite.DB().Find(&dbReport, report.ID)
		suite.NoError(err)
		suite.NotNil(dbReport.DeletedAt)

		err = deleter.DeleteEvaluationReport(appCtx, report.ID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
