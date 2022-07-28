package evaluationreport

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite EvaluationReportSuite) TestUpdateEvaluationReport() {
	// Since the service object doesn't hold any state, we can declare it up here and
	// reuse it in all of our tests
	updater := NewEvaluationReportUpdater()

	suite.Run("successful save", func() {
		// Create a report
		originalReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})

		// Copy report to new object
		report := originalReport

		// swag.String and other swag.<Type> functions just give you a pointer to the thing
		// you pass in. One might expect to be able to use the '&' operator, but Go won't
		// let you use it on literal values.
		report.Remarks = swag.String("spectacular packing job!!")

		// Attempt to update the report
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))
		suite.NoError(err)

		// Fetch the report from the database and make sure that it got updated
		var updatedReport models.EvaluationReport
		err = suite.DB().Find(&updatedReport, originalReport.ID)
		suite.NoError(err)

		suite.Equal(report.Remarks, updatedReport.Remarks)
		suite.Nil(updatedReport.SubmittedAt)
	})

	suite.Run("saving evaluation report with bad report id should fail", func() {
		// Create a report
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})

		// Overwrite the report's ID with some nonsense
		report.ID = uuid.Must(uuid.NewV4())

		// Attempt to update the report
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))

		// Our bogus report ID should cause an error
		suite.Error(err)
		// TODO check that we're getting the correct error
	})
	suite.Run("saving evaluation report created by another office user should fail", func() {
		// Create a report
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})

		otherOfficeUserID := uuid.Must(uuid.NewV4())

		// Attempt to update the report
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, otherOfficeUserID, etag.GenerateEtag(report.UpdatedAt))

		// Our bogus office user ID should cause an error
		suite.Error(err)
		suite.IsType(apperror.ForbiddenError{}, err)
	})

	suite.Run("updating a non-draft report should fail", func() {
		// Create a report
		originalReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				SubmittedAt: swag.Time(time.Now()),
			},
		})

		report := originalReport
		report.Remarks = swag.String("spectacular packing job!!")

		// Attempt to update the report
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))
		suite.Error(err)
	})
	suite.Run("updating a deleted report should fail", func() {
		// Create a report
		originalReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				DeletedAt: swag.Time(time.Now()),
			},
		})

		report := originalReport
		report.Remarks = swag.String("spectacular packing job!!")

		// Attempt to update the report
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
	suite.Run("updating a report with a bad ETag should fail", func() {
		// Create a report
		originalReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				SubmittedAt: swag.Time(time.Now()),
			},
		})

		report := originalReport
		report.Remarks = swag.String("spectacular packing job!!")

		// Attempt to update the report
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, "not a real etag")
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	// I'm going to need pointers to all of these things later on, and so I'm copying
	// them to variables that will be easy to turn into pointers with the address of operator (&) later.
	// Go won't directly give you a pointer to a constant or the return value of a function.
	// I also can't use swag.* functions on the constants here because they are a custom type.
	physical := models.EvaluationReportInspectionTypePhysical
	virtual := models.EvaluationReportInspectionTypeVirtual
	dataReview := models.EvaluationReportInspectionTypeDataReview
	currentTime := time.Now()

	// This is a "table driven test"
	// which means that I'm setting up a data structure that holds stuff that will be used
	// for multiple similar tests
	// and then iterating over it to generate a test for each item.
	// here's a blog with that explains some more about it: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
	// I'm not super happy with this one and will probably end up refactoring it a lot.
	testCases := map[string]struct {
		assertions testdatagen.Assertions
		// TODO when we figure out what the errors are going to look like, we should
		// TODO specify the expected errors in the test case.
		expectedError bool
	}{
		"travel time set for physical report type should succeed": {
			assertions: testdatagen.Assertions{
				EvaluationReport: models.EvaluationReport{
					InspectionType:    &physical,
					TravelTimeMinutes: swag.Int(30),
				},
			},
			expectedError: false,
		},
		"travel time set for virtual report type should fail": {
			assertions: testdatagen.Assertions{
				EvaluationReport: models.EvaluationReport{
					InspectionType:    &virtual,
					TravelTimeMinutes: swag.Int(30),
				},
			},
			expectedError: true,
		},
		"travel time set for data review report type should fail": {
			assertions: testdatagen.Assertions{
				EvaluationReport: models.EvaluationReport{
					InspectionType:    &dataReview,
					TravelTimeMinutes: swag.Int(30),
				},
			},
			expectedError: true,
		},
		"observed date set for physical report type should succeed": {
			assertions: testdatagen.Assertions{
				EvaluationReport: models.EvaluationReport{
					InspectionType: &physical,
					ObservedDate:   &currentTime,
				},
			},
			expectedError: false,
		},
		"observed date set for virtual report type should fail": {
			assertions: testdatagen.Assertions{
				EvaluationReport: models.EvaluationReport{
					InspectionType: &virtual,
					ObservedDate:   &currentTime,
				},
			},
			expectedError: true,
		},
		"observed date set for data review report type should fail": {
			assertions: testdatagen.Assertions{
				EvaluationReport: models.EvaluationReport{
					InspectionType: &dataReview,
					ObservedDate:   &currentTime,
				},
			},
			expectedError: true,
		},
	}

	// iterate over test cases
	for name, tc := range testCases {
		name := name
		tc := tc

		// Create a new test for each test case
		suite.Run(name, func() {
			report := testdatagen.MakeEvaluationReport(suite.DB(), tc.assertions)
			err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}

}
