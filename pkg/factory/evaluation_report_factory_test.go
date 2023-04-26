package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildEvaluationReport() {
	suite.Run("Successful creation of default evaluation report", func() {
		// Under test:      BuildEvaluationReport
		// Set up:          Create an evaluation report with no customizations or traits
		// Expected outcome: evaluation report should be created with default values

		// SETUP
		defaultOfficeUser := BuildOfficeUser(nil, nil, nil)
		defaultMove := BuildMove(nil, nil, nil)

		evaluationReport := BuildEvaluationReport(suite.DB(), nil, nil)

		suite.False(evaluationReport.OfficeUserID.IsNil())
		suite.Equal(defaultOfficeUser.FirstName, evaluationReport.OfficeUser.FirstName)
		suite.False(evaluationReport.MoveID.IsNil())
		suite.Equal(defaultMove.Status, evaluationReport.Move.Status)
		suite.Equal(models.EvaluationReportTypeCounseling, evaluationReport.Type)
	})

	suite.Run("Successful creation of custom EvaluationReportTypeShipment with custom shipment", func() {
		// Under test:      BuildEvaluationReport
		// Set up:          Create an evaluation report and pass custom fields
		// Expected outcome: evaluation report should be created with custom values

		// SETUP
		customOfficeuser := models.OfficeUser{
			FirstName: "Test",
		}
		customMove := models.Move{
			Locator: "AAAA",
		}
		inspectionType := models.EvaluationReportInspectionTypePhysical
		location := models.EvaluationReportLocationTypeOther
		evalStart := time.Now().AddDate(0, 0, -4)
		evalEnd := time.Now().AddDate(0, 0, -2)
		customEvaluationReport := models.EvaluationReport{
			Type:                               models.EvaluationReportTypeShipment,
			InspectionDate:                     models.TimePointer(time.Now()),
			InspectionType:                     &inspectionType,
			Location:                           &location,
			LocationDescription:                models.StringPointer("location description"),
			ObservedShipmentDeliveryDate:       &evalEnd,
			ObservedShipmentPhysicalPickupDate: &evalStart,
			TimeDepart:                         &evalEnd,
			EvalStart:                          &evalStart,
			EvalEnd:                            &evalEnd,
			ViolationsObserved:                 models.BoolPointer(true),
			Remarks:                            models.StringPointer("This is a remark."),
			SeriousIncident:                    models.BoolPointer(true),
			SeriousIncidentDesc:                models.StringPointer("very serious incident"),
			ObservedClaimsResponseDate:         models.TimePointer(time.Now()),
			ObservedPickupDate:                 &evalStart,
			ObservedPickupSpreadStartDate:      &evalStart,
			ObservedPickupSpreadEndDate:        &evalEnd,
			ObservedDeliveryDate:               &evalEnd,
		}
		customShipment := models.MTOShipment{
			ShipmentType: models.NTSRaw,
		}

		// CALL FUNCTION UNDER TEST
		evaluationReport := BuildEvaluationReport(suite.DB(), []Customization{
			{Model: customOfficeuser},
			{Model: customMove},
			{Model: customEvaluationReport},
			{Model: customShipment},
		}, nil)

		suite.NotNil(evaluationReport.OfficeUser)
		suite.False(evaluationReport.OfficeUserID.IsNil())
		suite.Equal(customOfficeuser.FirstName, evaluationReport.OfficeUser.FirstName)
		suite.False(evaluationReport.MoveID.IsNil())
		suite.Equal(customMove.Locator, evaluationReport.Move.Locator)
		suite.NotNil(evaluationReport.Shipment)
		suite.False(evaluationReport.ShipmentID.IsNil())
		suite.Equal(customShipment.ShipmentType, evaluationReport.Shipment.ShipmentType)
		suite.Equal(customEvaluationReport.Type, evaluationReport.Type)
		suite.Equal(customEvaluationReport.InspectionDate, evaluationReport.InspectionDate)
		suite.Equal(customEvaluationReport.InspectionType, evaluationReport.InspectionType)
		suite.Equal(customEvaluationReport.Location, evaluationReport.Location)
		suite.Equal(customEvaluationReport.LocationDescription, evaluationReport.LocationDescription)
		suite.Equal(customEvaluationReport.ObservedShipmentDeliveryDate, evaluationReport.ObservedShipmentDeliveryDate)
		suite.Equal(customEvaluationReport.ObservedShipmentPhysicalPickupDate, evaluationReport.ObservedShipmentPhysicalPickupDate)
		suite.Equal(customEvaluationReport.TimeDepart, evaluationReport.TimeDepart)
		suite.Equal(customEvaluationReport.EvalStart, evaluationReport.EvalStart)
		suite.Equal(customEvaluationReport.EvalEnd, evaluationReport.EvalEnd)
		suite.Equal(customEvaluationReport.ViolationsObserved, evaluationReport.ViolationsObserved)
		suite.Equal(customEvaluationReport.SeriousIncident, evaluationReport.SeriousIncident)
		suite.Equal(customEvaluationReport.SeriousIncidentDesc, evaluationReport.SeriousIncidentDesc)
		suite.Equal(customEvaluationReport.ObservedClaimsResponseDate, evaluationReport.ObservedClaimsResponseDate)
		suite.Equal(customEvaluationReport.ObservedPickupDate, evaluationReport.ObservedPickupDate)
		suite.Equal(customEvaluationReport.ObservedPickupSpreadStartDate, evaluationReport.ObservedPickupSpreadStartDate)
		suite.Equal(customEvaluationReport.ObservedPickupSpreadEndDate, evaluationReport.ObservedPickupSpreadEndDate)
		suite.Equal(customEvaluationReport.ObservedDeliveryDate, evaluationReport.ObservedDeliveryDate)
	})

	suite.Run("Successful creation of EvaluationReportTypeShipment without custom shipment", func() {
		// Under test:      BuildEvaluationReport
		// Set up:          Create an evaluation report and pass custom fields
		// Expected outcome: evaluation report should be created with custom values

		// SETUP
		customOfficeuser := models.OfficeUser{
			FirstName: "Test",
		}
		customMove := models.Move{
			Locator: "AAAA",
		}
		customEvaluationReport := models.EvaluationReport{
			Type: models.EvaluationReportTypeShipment,
		}

		// CALL FUNCTION UNDER TEST
		evaluationReport := BuildEvaluationReport(suite.DB(), []Customization{
			{Model: customOfficeuser},
			{Model: customMove},
			{Model: customEvaluationReport},
		}, nil)

		suite.False(evaluationReport.OfficeUserID.IsNil())
		suite.Equal(customOfficeuser.FirstName, evaluationReport.OfficeUser.FirstName)
		suite.False(evaluationReport.MoveID.IsNil())
		suite.Equal(customMove.Locator, evaluationReport.Move.Locator)
		suite.NotNil(evaluationReport.Shipment)
		suite.False(evaluationReport.ShipmentID.IsNil())
		suite.Equal(models.MTOShipmentTypeHHG, evaluationReport.Shipment.ShipmentType)
		suite.Equal(customEvaluationReport.Type, evaluationReport.Type)
	})

	suite.Run("Successful creation of stubbed evaluation report", func() {
		// Under test:      BuildEvaluationReport
		// Set up:          Create a stubbed evaluation report
		// Expected outcome:evaluation report should be created with stubbed move and office user
		precount, err := suite.DB().Count(&models.EvaluationReport{})
		suite.NoError(err)

		evaluationReport := BuildEvaluationReport(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(evaluationReport.OfficeUserID.IsNil())
		suite.True(evaluationReport.OfficeUser.ID.IsNil())
		suite.NotNil(evaluationReport.OfficeUser.FirstName)
		suite.True(evaluationReport.MoveID.IsNil())
		suite.True(evaluationReport.Move.ID.IsNil())
		suite.NotNil(evaluationReport.Move.Locator)
		suite.Equal(models.EvaluationReportTypeCounseling, evaluationReport.Type)

		// Count how many notification are in the DB, no new
		// notifications should have been created
		count, err := suite.DB().Count(&models.EvaluationReport{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful return of linkOnly evaluation report", func() {
		// Under test:       BuildEvaluationReport
		// Set up:           Pass in a linkOnly evaluation report
		// Expected outcome: No new evaluation report should be created.

		// Check num evaluation report records
		precount, err := suite.DB().Count(&models.EvaluationReport{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		evaluationReport := BuildEvaluationReport(suite.DB(), []Customization{
			{
				Model: models.EvaluationReport{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.EvaluationReport{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, evaluationReport.ID)
	})
}
