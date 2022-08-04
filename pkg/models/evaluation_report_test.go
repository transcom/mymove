package models_test

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func newUUIDPtr() *uuid.UUID {
	u := uuid.Must(uuid.NewV4())
	return &u
}

func (suite *ModelSuite) TestReport() {
	suite.Run("Create and query a report successfully", func() {
		reports := models.EvaluationReports{}
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		err := suite.DB().All(&reports)
		suite.NoError(err)
	})

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())
	virtualInspection := models.EvaluationReportInspectionTypeVirtual
	dataReviewInspection := models.EvaluationReportInspectionTypeDataReview

	testCases := map[string]struct {
		report         models.EvaluationReport
		expectedErrors map[string][]string
	}{
		"ShipmentID set for shipment report": {
			report: models.EvaluationReport{
				ID:           uuid.Must(uuid.NewV4()),
				OfficeUser:   officeUser,
				OfficeUserID: officeUser.ID,
				Move:         move,
				MoveID:       move.ID,
				Shipment:     nil,
				ShipmentID:   newUUIDPtr(),
				Type:         models.EvaluationReportTypeShipment,
			},
			expectedErrors: map[string][]string{},
		},
		"ShipmentID not set for shipment report": {
			report: models.EvaluationReport{
				ID:           uuid.Must(uuid.NewV4()),
				OfficeUser:   officeUser,
				OfficeUserID: officeUser.ID,
				Move:         move,
				MoveID:       move.ID,
				Shipment:     nil,
				ShipmentID:   nil,
				Type:         models.EvaluationReportTypeShipment,
			},
			expectedErrors: map[string][]string{
				"shipment_id": {"If report type is SHIPMENT, ShipmentID must not be null"},
			},
		},
		"ShipmentID set for non-shipment report": {
			report: models.EvaluationReport{
				ID:           uuid.Must(uuid.NewV4()),
				OfficeUser:   officeUser,
				OfficeUserID: officeUser.ID,
				Move:         move,
				MoveID:       move.ID,
				Shipment:     nil,
				ShipmentID:   newUUIDPtr(),
				Type:         models.EvaluationReportTypeCounseling,
			},
			expectedErrors: map[string][]string{
				"type": {"COUNSELING does not equal SHIPMENT."},
			},
		},
		"Non-physical inspection cannot have non-nil travel time": {
			report: models.EvaluationReport{
				ID:                uuid.Must(uuid.NewV4()),
				OfficeUser:        officeUser,
				OfficeUserID:      officeUser.ID,
				Move:              move,
				MoveID:            move.ID,
				Type:              models.EvaluationReportTypeCounseling,
				InspectionType:    &virtualInspection,
				TravelTimeMinutes: swag.Int(10),
			},
			expectedErrors: map[string][]string{
				"inspection_type": {"VIRTUAL does not equal PHYSICAL."},
			},
		},
		"ObservedDate cannot be set for virtual inspections": {
			report: models.EvaluationReport{
				ID:             uuid.Must(uuid.NewV4()),
				OfficeUser:     officeUser,
				OfficeUserID:   officeUser.ID,
				Move:           move,
				MoveID:         move.ID,
				Type:           models.EvaluationReportTypeCounseling,
				InspectionType: &virtualInspection,
				ObservedDate:   swag.Time(time.Now()),
			},
			expectedErrors: map[string][]string{
				"inspection_type": {"VIRTUAL does not equal PHYSICAL."},
			},
		},
		"ObservedDate cannot be set for data review inspections": {
			report: models.EvaluationReport{
				ID:             uuid.Must(uuid.NewV4()),
				OfficeUser:     officeUser,
				OfficeUserID:   officeUser.ID,
				Move:           move,
				MoveID:         move.ID,
				Type:           models.EvaluationReportTypeCounseling,
				InspectionType: &dataReviewInspection,
				ObservedDate:   swag.Time(time.Now()),
			},
			expectedErrors: map[string][]string{
				"inspection_type": {"DATA_REVIEW does not equal PHYSICAL."},
			},
		},
	}
	for name, test := range testCases {
		suite.Run(name, func() {
			suite.verifyValidationErrors(&test.report, test.expectedErrors)
		})
	}

}
