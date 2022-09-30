package reportviolation

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ReportViolationSuite) TestReportViolationCreator() {
	creator := NewReportViolationCreator()

	suite.Run("Can associate violations to reports successfully", func() {
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		violation := testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})

		rvID := uuid.Must(uuid.NewV4())
		reportViolation := &models.ReportViolation{ReportID: report.ID, ViolationID: violation.ID, ID: rvID}
		reportViolations := models.ReportViolations{*reportViolation}

		err := creator.AssociateReportViolations(suite.AppContextForTest(), &reportViolations, report.ID)

		suite.Nil(err)

		// Fetch the reportViolation from the database to verify it was created
		var updatedReportViolation models.ReportViolation
		err = suite.DB().Find(&updatedReportViolation, rvID)
		suite.NoError(err)

		suite.Equal(reportViolation.ID, updatedReportViolation.ID)
		suite.Equal(reportViolation.ViolationID, updatedReportViolation.ViolationID)
		suite.Equal(reportViolation.ReportID, updatedReportViolation.ReportID)
	})
	suite.Run("Association requires a valid violation", func() {
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		badID := uuid.Must(uuid.NewV4())
		reportViolation := &models.ReportViolation{ReportID: report.ID, ViolationID: badID}
		reportViolations := models.ReportViolations{*reportViolation}

		err := creator.AssociateReportViolations(suite.AppContextForTest(), &reportViolations, report.ID)

		suite.Error(err)
		suite.Equal("Could not complete query related to object of type: reportViolations.", err.Error())
	})

	suite.Run("Association requires a valid report", func() {
		violation := testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})
		badID := uuid.Must(uuid.NewV4())
		reportViolation := &models.ReportViolation{ReportID: badID, ViolationID: violation.ID}
		reportViolations := models.ReportViolations{*reportViolation}

		err := creator.AssociateReportViolations(suite.AppContextForTest(), &reportViolations, badID)

		suite.Error(err)
		suite.Equal("Could not complete query related to object of type: reportViolations.", err.Error())
	})
	suite.Run("Adding new associations to a report replaces existing associations for the report", func() {
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		violation := testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})
		violation2 := testdatagen.MakePWSViolation(suite.DB(),
			testdatagen.Assertions{Violation: models.PWSViolation{
				ID:                   uuid.Must(uuid.NewV4()),
				DisplayOrder:         2,
				ParagraphNumber:      "1.2.3.4",
				Title:                "Title 2",
				Category:             "Category 2",
				SubCategory:          "Customer Support",
				RequirementSummary:   "RequirementSummary 2",
				RequirementStatement: "RequirementStatement 2",
				IsKpi:                false,
				AdditionalDataElem:   "",
			}})
		rvID := uuid.Must(uuid.NewV4())
		rvID2 := uuid.Must(uuid.NewV4())
		reportViolation := &models.ReportViolation{ReportID: report.ID, ViolationID: violation.ID, ID: rvID}
		reportViolation2 := &models.ReportViolation{ReportID: report.ID, ViolationID: violation2.ID, ID: rvID2}
		reportViolations := models.ReportViolations{*reportViolation}
		reportViolations2 := models.ReportViolations{*reportViolation2}

		err := creator.AssociateReportViolations(suite.AppContextForTest(), &reportViolations, report.ID)
		suite.Nil(err)

		// Fetch the reportViolation from the database to the report and violation were associated
		var updatedReportViolation models.ReportViolation
		err = suite.DB().Find(&updatedReportViolation, rvID)
		suite.NoError(err)
		suite.Equal(reportViolation.ID, updatedReportViolation.ID)
		suite.Equal(reportViolation.ViolationID, updatedReportViolation.ViolationID)
		suite.Equal(reportViolation.ReportID, updatedReportViolation.ReportID)

		// Re-associate the report with the second violation
		err = creator.AssociateReportViolations(suite.AppContextForTest(), &reportViolations2, report.ID)
		suite.Nil(err)

		// Verify the new violation is associated with the report
		var updatedReportViolation2 models.ReportViolation
		err = suite.DB().Find(&updatedReportViolation2, rvID2)
		suite.NoError(err)
		suite.Equal(reportViolation2.ID, updatedReportViolation2.ID)

		// Verify the old association was replaced/removed
		var oldReportViolation models.ReportViolation
		err = suite.DB().Find(&oldReportViolation, rvID)
		suite.Error(err)
	})

}
