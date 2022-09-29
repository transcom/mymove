package evaluationreport

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite EvaluationReportSuite) TestSubmitEvaluationReport() {
	updater := NewEvaluationReportUpdater()

	suite.Run("Successful submission", func() {
		// Create office user
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		// Create a report
		inspectionType := models.EvaluationReportInspectionTypeVirtual
		location := models.EvaluationReportLocationTypeOrigin
		evaluationReport := testdatagen.MakeEvaluationReport(suite.DB(),
			testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
				OfficeUserID:            officeUser.ID,
				InspectionDate:          swag.Time(time.Now()),
				InspectionType:          &inspectionType,
				Location:                &location,
				EvaluationLengthMinutes: swag.Int(160),
				ViolationsObserved:      swag.Bool(false),
				Remarks:                 swag.String("This is a remark."),
			}})
		// Generate an etag
		eTag := etag.GenerateEtag(evaluationReport.UpdatedAt)
		// Submit the report
		err := updater.SubmitEvaluationReport(suite.AppContextForTest(), evaluationReport.ID, officeUser.ID, eTag)

		suite.NoError(err)
	})

	suite.Run("Bad etag", func() {
		// Create office user
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		// Create a report
		inspectionType := models.EvaluationReportInspectionTypeVirtual
		location := models.EvaluationReportLocationTypeOrigin
		evaluationReport := testdatagen.MakeEvaluationReport(suite.DB(),
			testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
				OfficeUserID:            officeUser.ID,
				InspectionDate:          swag.Time(time.Now()),
				InspectionType:          &inspectionType,
				Location:                &location,
				EvaluationLengthMinutes: swag.Int(160),
				ViolationsObserved:      swag.Bool(false),
				Remarks:                 swag.String("This is a remark."),
			}})
		// Generate an etag
		eTag := ""
		// Submit the report
		err := updater.SubmitEvaluationReport(suite.AppContextForTest(), evaluationReport.ID, officeUser.ID, eTag)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Missing required field", func() {
		// Create office user
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		// Create a report
		inspectionType := models.EvaluationReportInspectionTypeVirtual
		location := models.EvaluationReportLocationTypeOrigin
		// Missing inspection date
		evaluationReport := testdatagen.MakeEvaluationReport(suite.DB(),
			testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
				OfficeUserID:            officeUser.ID,
				InspectionType:          &inspectionType,
				Location:                &location,
				EvaluationLengthMinutes: swag.Int(160),
				ViolationsObserved:      swag.Bool(false),
				Remarks:                 swag.String("This is a remark."),
			}})
		// Generate an etag
		eTag := etag.GenerateEtag(evaluationReport.UpdatedAt)
		// Submit the report
		err := updater.SubmitEvaluationReport(suite.AppContextForTest(), evaluationReport.ID, officeUser.ID, eTag)
		suite.Equal(models.ErrInvalidTransition, errors.Cause(err))
	})

	suite.Run("Missing location description for physical location", func() {
		// Create office user
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		// Create a report
		inspectionType := models.EvaluationReportInspectionTypeVirtual
		location := models.EvaluationReportLocationTypeOther
		// Missing location description
		evaluationReport := testdatagen.MakeEvaluationReport(suite.DB(),
			testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
				OfficeUserID:            officeUser.ID,
				InspectionDate:          swag.Time(time.Now()),
				InspectionType:          &inspectionType,
				Location:                &location,
				EvaluationLengthMinutes: swag.Int(160),
				ViolationsObserved:      swag.Bool(false),
				Remarks:                 swag.String("This is a remark."),
			}})
		// Generate an etag
		eTag := etag.GenerateEtag(evaluationReport.UpdatedAt)
		// Submit the report
		err := updater.SubmitEvaluationReport(suite.AppContextForTest(), evaluationReport.ID, officeUser.ID, eTag)
		suite.Equal(models.ErrInvalidTransition, errors.Cause(err))
	})

	suite.Run("Missing travel time on physical location", func() {
		// Create office user
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{})
		// Create a report
		inspectionType := models.EvaluationReportInspectionTypePhysical
		location := models.EvaluationReportLocationTypeOrigin
		// Missing travel time
		evaluationReport := testdatagen.MakeEvaluationReport(suite.DB(),
			testdatagen.Assertions{EvaluationReport: models.EvaluationReport{
				OfficeUserID:            officeUser.ID,
				InspectionDate:          swag.Time(time.Now()),
				InspectionType:          &inspectionType,
				Location:                &location,
				EvaluationLengthMinutes: swag.Int(160),
				ViolationsObserved:      swag.Bool(false),
				Remarks:                 swag.String("This is a remark."),
			}})
		// Generate an etag
		eTag := etag.GenerateEtag(evaluationReport.UpdatedAt)
		// Submit the report
		err := updater.SubmitEvaluationReport(suite.AppContextForTest(), evaluationReport.ID, officeUser.ID, eTag)
		suite.Equal(models.ErrInvalidTransition, errors.Cause(err))
	})
}

func (suite EvaluationReportSuite) TestUpdateEvaluationReport() {
	updater := NewEvaluationReportUpdater()

	suite.Run("successful save", func() {
		// Create a report
		originalReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})

		// Copy report to new object
		report := originalReport

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

	suite.Run("saving report with pre-existing violations should delete them if violationsObserved is false", func() {
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{ReportViolation: models.ReportViolation{
			ReportID:    report.ID,
			Violation:   models.PWSViolation{},
			ViolationID: uuid.UUID{},
		}})
		report.ViolationsObserved = swag.Bool(false)

		// do the update
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))
		suite.NoError(err)

		var reportViolations models.ReportViolations
		err = suite.DB().Where("report_id = ?", report.ID).All(&reportViolations)
		suite.NoError(err)
		// we shouldn't find any report violations, which means this object should have a 0 length
		suite.Equal(0, len(reportViolations))
	})

	suite.Run("saving report does not overwrite readonly fields", func() {
		// Create a report and save it to the database
		move := testdatagen.MakeDefaultMove(suite.DB())
		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		originalReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			Move:        move,
			MTOShipment: shipment,
		})

		// Copy report to new object
		reportPayload := originalReport

		wrongUUID := uuid.Must(uuid.NewV4())
		reportPayload.Remarks = swag.String("spectacular packing job!!")
		reportPayload.MoveID = wrongUUID
		reportPayload.ShipmentID = &wrongUUID

		// Attempt to update the reportPayload
		err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &reportPayload, reportPayload.OfficeUserID, etag.GenerateEtag(reportPayload.UpdatedAt))
		suite.NoError(err)

		// Fetch the reportPayload from the database and make sure that it got updated
		var updatedReport models.EvaluationReport
		err = suite.DB().Find(&updatedReport, originalReport.ID)
		suite.NoError(err)

		suite.Equal(reportPayload.Remarks, updatedReport.Remarks)
		suite.Equal(originalReport.MoveID, updatedReport.MoveID)
		suite.Equal(*originalReport.ShipmentID, *updatedReport.ShipmentID)

		swaggerTimeFormat := "2006-01-02T15:04:05.99Z07:00"
		suite.Equal(originalReport.CreatedAt.Format(swaggerTimeFormat), time.Time(updatedReport.CreatedAt).Format(swaggerTimeFormat))
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
		suite.IsType(apperror.NotFoundError{}, err)
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

	physical := models.EvaluationReportInspectionTypePhysical
	virtual := models.EvaluationReportInspectionTypeVirtual
	dataReview := models.EvaluationReportInspectionTypeDataReview
	currentTime := time.Now()

	testCases := map[string]struct {
		inspectionType    *models.EvaluationReportInspectionType
		travelTimeMinutes *int
		observedDate      *time.Time
		expectedError     bool
	}{
		"travel time set for physical report type should succeed": {
			inspectionType:    &physical,
			travelTimeMinutes: swag.Int(30),
			expectedError:     false,
		},
		"travel time set for virtual report type should fail": {
			inspectionType:    &virtual,
			travelTimeMinutes: swag.Int(30),
			expectedError:     true,
		},
		"travel time set for data review report type should fail": {
			inspectionType:    &dataReview,
			travelTimeMinutes: swag.Int(30),
			expectedError:     true,
		},
		"observed date set for physical report type should succeed": {
			inspectionType: &physical,
			observedDate:   &currentTime,
			expectedError:  false,
		},
		"observed date set for virtual report type should fail": {
			inspectionType: &virtual,
			observedDate:   &currentTime,
			expectedError:  true,
		},
		"observed date set for data review report type should fail": {
			inspectionType: &dataReview,
			observedDate:   &currentTime,
			expectedError:  true,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
			report.InspectionType = tc.inspectionType
			report.TravelTimeMinutes = tc.travelTimeMinutes
			report.ObservedDate = tc.observedDate
			err := updater.UpdateEvaluationReport(suite.AppContextForTest(), &report, report.OfficeUserID, etag.GenerateEtag(report.UpdatedAt))
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
