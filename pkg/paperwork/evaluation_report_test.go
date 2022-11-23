package paperwork

import (
	"strings"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaperworkSuite) TestFormatValuesInspectionInformation() {
	suite.Run("FormatValuesInspectionInformation other location", func() {
		testDate := time.Date(2022, 10, 4, 0, 0, 0, 0, time.UTC)
		inspectionType := models.EvaluationReportInspectionTypePhysical
		location := models.EvaluationReportLocationTypeOrigin
		inspectionTime := time.Now().AddDate(0, 0, -4)
		timeDepart := inspectionTime
		evalStart := inspectionTime
		evalEnd := inspectionTime

		report := models.EvaluationReport{
			InspectionDate:     &testDate,
			InspectionType:     &inspectionType,
			Location:           &location,
			TimeDepart:         &timeDepart,
			EvalStart:          &evalStart,
			EvalEnd:            &evalEnd,
			ViolationsObserved: swag.Bool(false),
			Remarks:            swag.String("remarks"),
			UpdatedAt:          time.Time{},
		}
		values := FormatValuesInspectionInformation(report)
		suite.Equal("04 October 2022", values.DateOfInspection)
		suite.Equal("Physical", values.EvaluationType)
		suite.Equal("Origin", values.EvaluationLocation)
		suite.Equal(timeDepart, values.TimeDepart)
		suite.Equal(evalStart, values.EvalStart)
		suite.Equal(evalEnd, values.EvalEnd)
		suite.Equal("remarks", values.QAERemarks)
		suite.Equal("No", values.ViolationsObserved)

		// serious incident values are empty if no violations
		suite.Equal("", values.SeriousIncident)
		suite.Equal("", values.SeriousIncidentDescription)
	})

	suite.Run("FormatValuesInspectionInformation violations", func() {
		testDate := time.Date(2022, 10, 4, 0, 0, 0, 0, time.UTC)
		inspectionType := models.EvaluationReportInspectionTypePhysical
		location := models.EvaluationReportLocationTypeOther
		locationDescription := "other location"
		report := models.EvaluationReport{
			InspectionDate:      &testDate,
			InspectionType:      &inspectionType,
			Location:            &location,
			LocationDescription: &locationDescription,
			ViolationsObserved:  swag.Bool(true),
			Remarks:             swag.String("remarks"),
			SeriousIncident:     swag.Bool(true),
			SeriousIncidentDesc: swag.String("serious incident"),
			UpdatedAt:           time.Time{},
		}
		values := FormatValuesInspectionInformation(report)
		suite.Equal("04 October 2022", values.DateOfInspection)
		suite.Equal("Physical", values.EvaluationType)
		suite.Equal("Other\nother location", values.EvaluationLocation)
		suite.Equal("1 hr 0 min", values.EvaluationLength)
		suite.Equal("remarks", values.QAERemarks)
		suite.Equal("Yes\nViolations are listed on a subsequent page", values.ViolationsObserved)

		// serious incident values are empty if no violations
		suite.Equal("Yes", values.SeriousIncident)
		suite.Equal("serious incident", values.SeriousIncidentDescription)
	})
}

func (suite *PaperworkSuite) TestFormatValuesShipment() {
	suite.Run("storage facility with phone and email", func() {
		storageFacility := testdatagen.MakeDefaultStorageFacility(suite.DB())
		shipment := testdatagen.MakeNTSShipment(suite.DB(), testdatagen.Assertions{
			StorageFacility: storageFacility,
			MTOShipment:     models.MTOShipment{StorageFacility: &storageFacility},
		})

		shipmentValues := FormatValuesShipment(shipment)
		expectedContactInfo := strings.Join([]string{*storageFacility.Phone, *storageFacility.Email}, "\n")
		suite.Equal(expectedContactInfo, shipmentValues.StorageFacility)
	})

	suite.Run("storage facility with no phone number should not panic", func() {
		storageFacility := testdatagen.MakeDefaultStorageFacility(suite.DB())
		storageFacility.Phone = nil
		suite.MustSave(&storageFacility)
		shipment := testdatagen.MakeNTSShipment(suite.DB(), testdatagen.Assertions{
			StorageFacility: storageFacility,
			MTOShipment:     models.MTOShipment{StorageFacility: &storageFacility},
		})

		shipmentValues := FormatValuesShipment(shipment)
		suite.Equal(*storageFacility.Email, shipmentValues.StorageFacility)
	})

	suite.Run("storage facility with no email should not panic", func() {
		storageFacility := testdatagen.MakeDefaultStorageFacility(suite.DB())
		storageFacility.Email = nil
		suite.MustSave(&storageFacility)
		shipment := testdatagen.MakeNTSShipment(suite.DB(), testdatagen.Assertions{
			StorageFacility: storageFacility,
			MTOShipment:     models.MTOShipment{StorageFacility: &storageFacility},
		})

		shipmentValues := FormatValuesShipment(shipment)
		suite.Equal(*storageFacility.Phone, shipmentValues.StorageFacility)
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
