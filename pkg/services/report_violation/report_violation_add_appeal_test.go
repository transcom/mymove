package reportviolation

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ReportViolationSuite) TestAddAppealToViolation() {
	appealAdder := NewReportViolationsAddAppeal()
	suite.Run("Successfully adds an appeal to a violation", func() {
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		violation := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{})
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToViolation(
			suite.AppContextForTest(),
			report.ID,
			violation.ID,
			officeUser.ID,
			remarks,
			appealStatus,
		)

		suite.NoError(err)
		suite.NotNil(appeal)
	})

	suite.Run("Returns error for nil reportID", func() {
		violation := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{})
		officeUserID := uuid.Must(uuid.NewV4())
		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToViolation(
			suite.AppContextForTest(),
			uuid.Nil,
			violation.ID,
			officeUserID,
			remarks,
			appealStatus,
		)

		suite.Error(err)
		suite.Equal(models.GsrAppeal{}, appeal)
		suite.Contains(err.Error(), "reportID must be provided")
	})

	suite.Run("Returns error for nil reportViolationID", func() {
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		officeUserID := uuid.Must(uuid.NewV4())
		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToViolation(
			suite.AppContextForTest(),
			report.ID,
			uuid.Nil,
			officeUserID,
			remarks,
			appealStatus,
		)

		suite.Error(err)
		suite.Equal(models.GsrAppeal{}, appeal)
		suite.Contains(err.Error(), "reportViolationID must be provided")
	})

	suite.Run("Returns error for nil officeUserID", func() {
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		violation := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{})
		remarks := "Test remarks"
		appealStatus := "sustained"

		appeal, err := appealAdder.AddAppealToViolation(
			suite.AppContextForTest(),
			report.ID,
			violation.ID,
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
		violation := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{})
		officeUserID := uuid.Must(uuid.NewV4())
		remarks := "Test remarks"
		appealStatus := "invalid_status"

		appeal, err := appealAdder.AddAppealToViolation(
			suite.AppContextForTest(),
			report.ID,
			violation.ID,
			officeUserID,
			remarks,
			appealStatus,
		)

		suite.Error(err)
		suite.Equal(models.GsrAppeal{}, appeal)
	})
}
