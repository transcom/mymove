package internalapi

import (
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"testing"
	"time"
)

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentMissingPickupDate() {
	// create a shipment
	transitDays := int64(5)
	packDays := int64(3)
	var shipment = models.Shipment{
		EstimatedTransitDays: &transitDays,
		EstimatedPackDays:    &packDays,
	}

	_, err := calculateMoveDatesFromShipment(&shipment)

	suite.Error(err)
}

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentMissingOriginalPackDate() {
	// create a shipment
	pickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	deliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)
	var shipment = models.Shipment{
		RequestedPickupDate:  &pickupDate,
		OriginalDeliveryDate: &deliveryDate,
	}
	_, err := calculateMoveDatesFromShipment(&shipment)

	suite.Error(err)
}

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentMissingOriginalDeliveryDate() {
	// create a shipment
	pickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	packDate := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	var shipment = models.Shipment{
		RequestedPickupDate: &pickupDate,
		OriginalPackDate:    &packDate,
	}

	_, err := calculateMoveDatesFromShipment(&shipment)

	suite.Error(err)
}

type testCase struct {
	name     string
	shipment models.Shipment
	summary  MoveDatesSummary
}

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipment() {
	originalPackDate := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	requestedPickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	originalDeliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)

	pmSurveyPackDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	pmSurveyPickupDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)
	pmSurveyDeliveryDate := time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC)

	actualPackDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)
	actualDeliveryDate := time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC)

	mostCurrentActualPackDate := time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)
	mostCurrentPmSurveyPickupDate := time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC)

	earlierDeliveryDate := time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)
	equalToDeliveryDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)

	var cases = []testCase{
		{
			name: "all original dates",
			shipment: models.Shipment{
				RequestedPickupDate:  &requestedPickupDate,
				OriginalPackDate:     &originalPackDate,
				OriginalDeliveryDate: &originalDeliveryDate,
			},
			summary: MoveDatesSummary{
				[]time.Time{time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)},
				[]time.Time{requestedPickupDate},
				[]time.Time{time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC)},
				[]time.Time{originalDeliveryDate},
				[]time.Time{},
			},
		},
		{
			name: "all pm survey dates",
			shipment: models.Shipment{
				RequestedPickupDate:  &requestedPickupDate,
				OriginalPackDate:     &originalPackDate,
				OriginalDeliveryDate: &originalDeliveryDate,

				PmSurveyPlannedPickupDate:   &pmSurveyPickupDate,
				PmSurveyPlannedPackDate:     &pmSurveyPackDate,
				PmSurveyPlannedDeliveryDate: &pmSurveyDeliveryDate,
			},
			summary: MoveDatesSummary{
				[]time.Time{
					time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{pmSurveyPickupDate},
				[]time.Time{
					time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{pmSurveyDeliveryDate},
				[]time.Time{},
			},
		},
		{
			name: "all actual dates",
			shipment: models.Shipment{
				RequestedPickupDate:  &requestedPickupDate,
				OriginalPackDate:     &originalPackDate,
				OriginalDeliveryDate: &originalDeliveryDate,

				ActualPickupDate:   &actualPickupDate,
				ActualPackDate:     &actualPackDate,
				ActualDeliveryDate: &actualDeliveryDate,
			},
			summary: MoveDatesSummary{
				[]time.Time{
					time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{actualPickupDate},
				[]time.Time{
					time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{actualDeliveryDate},
				[]time.Time{},
			},
		},
		{
			name: "mixed most current dates",
			shipment: models.Shipment{
				RequestedPickupDate:  &requestedPickupDate,
				OriginalPackDate:     &originalPackDate,
				OriginalDeliveryDate: &originalDeliveryDate,

				ActualPackDate:            &mostCurrentActualPackDate,
				PmSurveyPlannedPickupDate: &mostCurrentPmSurveyPickupDate,
			},
			summary: MoveDatesSummary{
				[]time.Time{
					time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{mostCurrentPmSurveyPickupDate},
				[]time.Time{
					time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{originalDeliveryDate},
				[]time.Time{},
			},
		},
		{
			name: "delivery date is before pickup date",
			shipment: models.Shipment{
				RequestedPickupDate:  &requestedPickupDate,
				OriginalPackDate:     &originalPackDate,
				OriginalDeliveryDate: &originalDeliveryDate,

				ActualPackDate:              &actualPackDate,
				ActualPickupDate:            &actualPickupDate,
				PmSurveyPlannedDeliveryDate: &earlierDeliveryDate,
			},
			summary: MoveDatesSummary{
				[]time.Time{
					time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{actualPickupDate},
				[]time.Time(nil),
				[]time.Time{earlierDeliveryDate},
				[]time.Time{},
			},
		},
		{
			name: "delivery date is equal to pickup date",
			shipment: models.Shipment{
				RequestedPickupDate:  &requestedPickupDate,
				OriginalPackDate:     &originalPackDate,
				OriginalDeliveryDate: &originalDeliveryDate,

				ActualPackDate:              &actualPackDate,
				ActualPickupDate:            &actualPickupDate,
				PmSurveyPlannedDeliveryDate: &equalToDeliveryDate,
			},
			summary: MoveDatesSummary{
				[]time.Time{
					time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{actualPickupDate},
				[]time.Time(nil),
				[]time.Time{equalToDeliveryDate},
				[]time.Time{},
			},
		},
	}

	for _, testCase := range cases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			summary, err := calculateMoveDatesFromShipment(&testCase.shipment)
			suite.Nil(err)
			suite.Equal(testCase.summary.PackDays, summary.PackDays, "PackDays did not match, expected %v, got %v", testCase.summary.PackDays, summary.PackDays)
			suite.Equal(testCase.summary.PickupDays, summary.PickupDays, "PickupDays did not match, expected %v, got %v", testCase.summary.PickupDays, summary.PickupDays)
			suite.Equal(testCase.summary.TransitDays, summary.TransitDays, "TransitDays did not match, expected %v, got %v", testCase.summary.TransitDays, summary.TransitDays)
			suite.Equal(testCase.summary.DeliveryDays, summary.DeliveryDays, "DeliveryDays did not match, expected %v, got %v", testCase.summary.DeliveryDays, summary.DeliveryDays)
		})
	}
}

func (suite *HandlerSuite) TestCreateValidDatesBetweenTwoDatesEndDateMustBeLater() {
	startDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	usCalendar := handlers.NewUSCalendar()
	_, err := createValidDatesBetweenTwoDates(startDate, endDate, true, false, usCalendar)
	suite.Error(err)
}
