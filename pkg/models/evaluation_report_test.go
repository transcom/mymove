package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func newUUIDPtr() *uuid.UUID {
	u := uuid.Must(uuid.NewV4())
	return &u
}

func (suite *ModelSuite) TestReport() {
	suite.Run("Create and query a report successfully", func() {
		reports := models.EvaluationReports{}
		factory.BuildEvaluationReport(suite.DB(), nil, nil)
		err := suite.DB().All(&reports)
		suite.NoError(err)
	})

	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	move := factory.BuildMove(suite.DB(), nil, nil)
	virtualInspection := models.EvaluationReportInspectionTypeVirtual
	dataReviewInspection := models.EvaluationReportInspectionTypeDataReview
	physicalInspection := models.EvaluationReportInspectionTypePhysical
	location := models.EvaluationReportLocationTypeOrigin

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
		"Physical inspection with time departed": {
			report: models.EvaluationReport{
				ID:             uuid.Must(uuid.NewV4()),
				OfficeUser:     officeUser,
				OfficeUserID:   officeUser.ID,
				Move:           move,
				MoveID:         move.ID,
				Type:           models.EvaluationReportTypeCounseling,
				InspectionType: &physicalInspection,
				TimeDepart:     models.TimePointer(time.Now()),
				Location:       &location,
			},
			expectedErrors: map[string][]string{},
		},
		"ObservedShipmentDeliveryDate cannot be set for virtual inspections": {
			report: models.EvaluationReport{
				ID:                           uuid.Must(uuid.NewV4()),
				OfficeUser:                   officeUser,
				OfficeUserID:                 officeUser.ID,
				Move:                         move,
				MoveID:                       move.ID,
				Type:                         models.EvaluationReportTypeCounseling,
				InspectionType:               &virtualInspection,
				ObservedShipmentDeliveryDate: models.TimePointer(time.Now()),
			},
			expectedErrors: map[string][]string{
				"inspection_type": {"VIRTUAL does not equal PHYSICAL."},
			},
		},
		"ObservedShipmentPhysicalPickupDate cannot be set for data review inspections": {
			report: models.EvaluationReport{
				ID:                                 uuid.Must(uuid.NewV4()),
				OfficeUser:                         officeUser,
				OfficeUserID:                       officeUser.ID,
				Move:                               move,
				MoveID:                             move.ID,
				Type:                               models.EvaluationReportTypeCounseling,
				InspectionType:                     &dataReviewInspection,
				ObservedShipmentPhysicalPickupDate: models.TimePointer(time.Now()),
			},
			expectedErrors: map[string][]string{
				"inspection_type": {"DATA_REVIEW does not equal PHYSICAL."},
			},
		},
	}
	for name, test := range testCases {
		suite.Run(name, func() {
			//nolint:gosec //G601
			suite.verifyValidationErrors(&test.report, test.expectedErrors, nil)
		})
	}

}
