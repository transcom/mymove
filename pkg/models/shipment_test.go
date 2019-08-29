package models_test

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/dates"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_ShipmentValidations() {
	packDays := int64(-2)
	transitDays := int64(0)
	var weightEstimate unit.Pound = -3
	var progearWeightEstimate unit.Pound = -12
	var spouseProgearWeightEstimate unit.Pound = -9
	calendar := dates.NewUSCalendar()
	weekendDate := dates.NextNonWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 25, 0, 0, 0, 0, time.UTC))

	shipment := &Shipment{
		EstimatedPackDays:           &packDays,
		EstimatedTransitDays:        &transitDays,
		WeightEstimate:              &weightEstimate,
		ProgearWeightEstimate:       &progearWeightEstimate,
		SpouseProgearWeightEstimate: &spouseProgearWeightEstimate,
		RequestedPickupDate:         &weekendDate,
		OriginalDeliveryDate:        &weekendDate,
		OriginalPackDate:            &weekendDate,
		PmSurveyPlannedPackDate:     &weekendDate,
		PmSurveyPlannedPickupDate:   &weekendDate,
		PmSurveyPlannedDeliveryDate: &weekendDate,
		ActualPackDate:              &weekendDate,
		ActualPickupDate:            &weekendDate,
		ActualDeliveryDate:          &weekendDate,
	}

	stringDate := weekendDate.Format("2006-01-02 15:04:05 -0700 UTC")
	expErrors := map[string][]string{
		"move_id":                         {"move_id can not be blank."},
		"status":                          {"status can not be blank."},
		"estimated_pack_days":             {"-2 is less than or equal to zero."},
		"estimated_transit_days":          {"0 is less than or equal to zero."},
		"weight_estimate":                 {"-3 is less than zero."},
		"progear_weight_estimate":         {"-12 is less than zero."},
		"spouse_progear_weight_estimate":  {"-9 is less than zero."},
		"requested_pickup_date":           {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"original_delivery_date":          {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"original_pack_date":              {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"pm_survey_planned_pack_date":     {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"pm_survey_planned_pickup_date":   {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"pm_survey_planned_delivery_date": {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"actual_pack_date":                {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"actual_pickup_date":              {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
		"actual_delivery_date":            {fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate)},
	}

	suite.verifyValidationErrors(shipment, expErrors)
}

func (suite *ModelSuite) Test_ShipmentValidationsSubmittedMove() {
	shipment := &Shipment{
		Status: ShipmentStatusSUBMITTED,
	}

	verrs, err := shipment.Validate(suite.DB())
	suite.NoError(err)

	pickupAddressErrors := verrs.Get("pickup_address_id")
	suite.Equal(1, len(pickupAddressErrors), "expected one error on pickup_address_id, but there were %d: %v", len(pickupAddressErrors), pickupAddressErrors)

	suite.Equal(pickupAddressErrors[0], "pickup_address_id can not be blank.", "expected pickup_address_id to be required, but it was not")
}
