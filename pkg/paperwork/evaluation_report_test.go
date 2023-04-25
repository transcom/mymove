package paperwork

import (
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaperworkSuite) TestFormatValuesInspectionInformation() {
	suite.Run("FormatValuesInspectionInformation other location", func() {
		testDate := time.Date(2022, 10, 4, 0, 0, 0, 0, time.UTC)
		inspectionType := models.EvaluationReportInspectionTypePhysical
		location := models.EvaluationReportLocationTypeOrigin
		timeDepart := time.Now().AddDate(0, 0, -4)
		evalStart := timeDepart.Add(time.Hour * 2)
		evalEnd := evalStart.Add(time.Minute * 45)

		report := models.EvaluationReport{
			InspectionDate:     &testDate,
			InspectionType:     &inspectionType,
			Location:           &location,
			TimeDepart:         &timeDepart,
			EvalStart:          &evalStart,
			EvalEnd:            &evalEnd,
			ViolationsObserved: models.BoolPointer(false),
			Remarks:            models.StringPointer("remarks"),
			UpdatedAt:          time.Time{},
		}

		values := FormatValuesInspectionInformation(report)

		suite.Equal("04 October 2022", values.DateOfInspection)
		suite.Equal("Physical", values.EvaluationType)
		suite.Equal("Origin", values.EvaluationLocation)
		suite.Equal(timeDepart.Format(timeFormat), values.TimeDepart)
		suite.Equal(evalStart.Format(timeFormat), values.EvalStart)
		suite.Equal(evalEnd.Format(timeFormat), values.EvalEnd)
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
			ViolationsObserved:  models.BoolPointer(true),
			Remarks:             models.StringPointer("remarks"),
			SeriousIncident:     models.BoolPointer(true),
			SeriousIncidentDesc: models.StringPointer("serious incident"),
			UpdatedAt:           time.Time{},
		}
		values := FormatValuesInspectionInformation(report)
		suite.Equal("04 October 2022", values.DateOfInspection)
		suite.Equal("Physical", values.EvaluationType)
		suite.Equal("Other\nother location", values.EvaluationLocation)
		suite.Equal("remarks", values.QAERemarks)
		suite.Equal("Yes\nViolations are listed on a subsequent page", values.ViolationsObserved)

		// serious incident values are empty if no violations
		suite.Equal("Yes", values.SeriousIncident)
		suite.Equal("serious incident", values.SeriousIncidentDescription)
	})
}

func (suite *PaperworkSuite) TestFormatValuesShipment() {
	suite.Run("storage facility with phone and email", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)
		shipment := factory.BuildNTSShipment(suite.DB(), []factory.Customization{
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)

		shipmentValues := FormatValuesShipment(shipment)
		expectedContactInfo := strings.Join([]string{*storageFacility.Phone, *storageFacility.Email}, "\n")
		suite.Equal(expectedContactInfo, shipmentValues.StorageFacility)
	})

	suite.Run("storage facility with no phone number should not panic", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)
		storageFacility.Phone = nil
		suite.MustSave(&storageFacility)
		shipment := factory.BuildNTSShipment(suite.DB(), []factory.Customization{
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)

		shipmentValues := FormatValuesShipment(shipment)
		suite.Equal(*storageFacility.Email, shipmentValues.StorageFacility)
	})

	suite.Run("storage facility with no email should not panic", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)
		storageFacility.Email = nil
		suite.MustSave(&storageFacility)
		shipment := factory.BuildNTSShipment(suite.DB(), []factory.Customization{
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)

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
