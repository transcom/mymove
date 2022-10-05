package paperwork

import (
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaperworkSuite) TestFormatValuesInspectionInformation() {
	suite.Run("FormatValuesInspectionInformation other location", func() {
		testDate := time.Date(2022, 10, 4, 0, 0, 0, 0, time.UTC)
		inspectionType := models.EvaluationReportInspectionTypePhysical
		testDurationMinutes := 60
		location := models.EvaluationReportLocationTypeOther
		locationDescription := "other location"
		report := models.EvaluationReport{
			InspectionDate:          &testDate,
			InspectionType:          &inspectionType,
			TravelTimeMinutes:       &testDurationMinutes,
			Location:                &location,
			LocationDescription:     &locationDescription,
			EvaluationLengthMinutes: &testDurationMinutes,
			ViolationsObserved:      swag.Bool(false),
			Remarks:                 swag.String("remarks"),
			UpdatedAt:               time.Time{},
		}
		values := FormatValuesInspectionInformation(report)
		suite.Equal("04 October 2022", values.DateOfInspection)
		suite.Equal("Physical", values.EvaluationType)
		suite.Equal("1 hr 0 min", values.TravelTimeToEvaluation)
		suite.Equal("Other\nother location", values.EvaluationLocation)
		//suite.Equal("04 October 2022", values.ObservedPickupDate)
		suite.Equal("1 hr 0 min", values.EvaluationLength)
		suite.Equal("remarks", values.QAERemarks)
		suite.Equal("No", values.ViolationsObserved)

		// serious incident values are empty if no violations
		suite.Equal("", values.SeriousIncident)
		suite.Equal("", values.SeriousIncidentDescription)
	})

	suite.Run("FormatValuesInspectionInformation violations", func() {
		testDate := time.Date(2022, 10, 4, 0, 0, 0, 0, time.UTC)
		inspectionType := models.EvaluationReportInspectionTypePhysical
		testDurationMinutes := 60
		location := models.EvaluationReportLocationTypeOther
		locationDescription := "other location"
		report := models.EvaluationReport{
			InspectionDate:          &testDate,
			InspectionType:          &inspectionType,
			TravelTimeMinutes:       &testDurationMinutes,
			Location:                &location,
			LocationDescription:     &locationDescription,
			EvaluationLengthMinutes: &testDurationMinutes,
			ViolationsObserved:      swag.Bool(true),
			Remarks:                 swag.String("remarks"),
			SeriousIncident:         swag.Bool(true),
			SeriousIncidentDesc:     swag.String("serious incident"),
			UpdatedAt:               time.Time{},
		}
		values := FormatValuesInspectionInformation(report)
		suite.Equal("04 October 2022", values.DateOfInspection)
		suite.Equal("Physical", values.EvaluationType)
		suite.Equal("1 hr 0 min", values.TravelTimeToEvaluation)
		suite.Equal("Other\nother location", values.EvaluationLocation)
		suite.Equal("1 hr 0 min", values.EvaluationLength)
		suite.Equal("remarks", values.QAERemarks)
		suite.Equal("Yes\nViolations are listed on a subsequent page", values.ViolationsObserved)

		// serious incident values are empty if no violations
		suite.Equal("Yes", values.SeriousIncident)
		suite.Equal("serious incident", values.SeriousIncidentDescription)
	})
}

func (suite *PaperworkSuite) TestPickShipmentCardLayout() {
	suite.Run("HHG", func() {
		suite.ElementsMatch(HHGShipmentCardLayout, PickShipmentCardLayout(models.MTOShipmentTypeHHG))
	})
	suite.Run("PPM", func() {
		suite.ElementsMatch(PPMShipmentCardLayout, PickShipmentCardLayout(models.MTOShipmentTypePPM))
	})
	suite.Run("NTS", func() {
		suite.ElementsMatch(NTSShipmentCardLayout, PickShipmentCardLayout(models.MTOShipmentTypeHHGIntoNTSDom))
	})
	suite.Run("NTS-R", func() {
		suite.ElementsMatch(NTSRShipmentCardLayout, PickShipmentCardLayout(models.MTOShipmentTypeHHGOutOfNTSDom))
	})
}
