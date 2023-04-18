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
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EvaluationReportSuite) TestEvaluationReportDeleter() {
	setupTestData := func() (services.EvaluationReportDeleter, models.EvaluationReport, appcontext.AppContext) {
		deleter := NewEvaluationReportDeleter()
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			OfficeUser:  officeUser,
			MTOShipment: shipment,
			Move:        move,
		})
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: officeUser.ID,
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
