package evaluationreport

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *EvaluationReportSuite) TestAddAppealToSeriousIncident() {
	appealAdder := NewEvaluationReportSeriousIncidentAddAppeal()
	suite.Run("Successfully adds an appeal to a serious incident", func() {
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToSeriousIncident(
			suite.AppContextForTest(),
			report.ID,
			officeUser.ID,
			remarks,
			appealStatus,
		)

		suite.NoError(err)
		suite.NotNil(appeal)
	})

	suite.Run("Returns error for nil reportID", func() {
		officeUserID := uuid.Must(uuid.NewV4())
		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToSeriousIncident(
			suite.AppContextForTest(),
			uuid.Nil,
			officeUserID,
			remarks,
			appealStatus,
		)

		suite.Error(err)
		suite.Equal(models.GsrAppeal{}, appeal)
		suite.Contains(err.Error(), "reportID must be provided")
	})

	suite.Run("Returns error for nil officeUserID", func() {
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToSeriousIncident(
			suite.AppContextForTest(),
			report.ID,
			uuid.Nil,
			remarks,
			appealStatus,
		)

		suite.Error(err)
		suite.Equal(models.GsrAppeal{}, appeal)
		suite.Contains(err.Error(), "officeUserID must be provided")
	})

	suite.Run("Returns error for invalid appeal status", func() {
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		officeUserID := uuid.Must(uuid.NewV4())
		remarks := "Test remarks"
		appealStatus := "invalid_status"

		appeal, err := appealAdder.AddAppealToSeriousIncident(
			suite.AppContextForTest(),
			report.ID,
			officeUserID,
			remarks,
			appealStatus,
		)

		suite.Error(err)
		suite.Equal(models.GsrAppeal{}, appeal)
	})
}
