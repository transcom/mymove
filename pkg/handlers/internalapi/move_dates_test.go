package internalapi

import (
	"github.com/transcom/mymove/pkg/models"
	"time"
)

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipment() {
	// create a shipment
	pickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	packDate := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	deliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)
	var shipment = models.Shipment{
		RequestedPickupDate:  &pickupDate,
		OriginalPackDate:     &packDate,
		OriginalDeliveryDate: &deliveryDate,
	}

	summary, err := calculateMoveDatesFromShipment(&shipment)

	// check that there is no error
	suite.Nil(err)
	// compare expected output with actual output
	expectedPickupDays := []time.Time{pickupDate}
	suite.Equal(expectedPickupDays, summary.PickupDays)

	expectedDeliveryDays := []time.Time{deliveryDate}
	suite.Equal(expectedDeliveryDays, summary.DeliveryDays)

	expectedPackDays := []time.Time{
		time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedPackDays, summary.PackDays)

	expectedTransitDays := []time.Time{
		time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedTransitDays, summary.TransitDays)

}

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

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentOriginalDeliveryDate() {
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

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentUsePMSurveyDates() {
	// create a shipment
	requestedPickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	pmSurveyPickupDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)
	originalPackDate := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	pmSurveyPackDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	originalDeliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)
	pmSurveyDeliveryDate := time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC)
	var shipment = models.Shipment{
		RequestedPickupDate:         &requestedPickupDate,
		OriginalPackDate:            &originalPackDate,
		OriginalDeliveryDate:        &originalDeliveryDate,
		PmSurveyPlannedPickupDate:   &pmSurveyPickupDate,
		PmSurveyPlannedPackDate:     &pmSurveyPackDate,
		PmSurveyPlannedDeliveryDate: &pmSurveyDeliveryDate,
	}

	summary, err := calculateMoveDatesFromShipment(&shipment)

	// check that there is no error
	suite.Nil(err)
	// compare expected output with actual output
	expectedPickupDays := []time.Time{pmSurveyPickupDate}
	suite.Equal(expectedPickupDays, summary.PickupDays)

	expectedDeliveryDays := []time.Time{pmSurveyDeliveryDate}
	suite.Equal(expectedDeliveryDays, summary.DeliveryDays)

	expectedPackDays := []time.Time{
		time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedPackDays, summary.PackDays)

	expectedTransitDays := []time.Time{
		time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedTransitDays, summary.TransitDays)
}

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentUseActualDates() {
	// create a shipment
	requestedPickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)
	originalPackDate := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	actualPackDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	originalDeliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)
	actualDeliveryDate := time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC)
	var shipment = models.Shipment{
		RequestedPickupDate:  &requestedPickupDate,
		OriginalPackDate:     &originalPackDate,
		OriginalDeliveryDate: &originalDeliveryDate,
		ActualPickupDate:     &actualPickupDate,
		ActualPackDate:       &actualPackDate,
		ActualDeliveryDate:   &actualDeliveryDate,
	}

	summary, err := calculateMoveDatesFromShipment(&shipment)

	// check that there is no error
	suite.Nil(err)
	// compare expected output with actual output
	expectedPickupDays := []time.Time{actualPickupDate}
	suite.Equal(expectedPickupDays, summary.PickupDays)

	expectedDeliveryDays := []time.Time{actualDeliveryDate}
	suite.Equal(expectedDeliveryDays, summary.DeliveryDays)

	expectedPackDays := []time.Time{
		time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedPackDays, summary.PackDays)

	expectedTransitDays := []time.Time{
		time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedTransitDays, summary.TransitDays)
}

func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentUseMostCurrentDates() {
	// create a shipment
	requestedPickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	pmSurveyPickupDate := time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC)
	originalPackDate := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	actualPackDate := time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)
	originalDeliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)

	var shipment = models.Shipment{
		RequestedPickupDate:  &requestedPickupDate,
		OriginalPackDate:     &originalPackDate,
		OriginalDeliveryDate: &originalDeliveryDate,

		ActualPackDate:            &actualPackDate,
		PmSurveyPlannedPickupDate: &pmSurveyPickupDate,
	}

	summary, err := calculateMoveDatesFromShipment(&shipment)

	// check that there is no error
	suite.Nil(err)
	// compare expected output with actual output
	expectedPickupDays := []time.Time{pmSurveyPickupDate}
	suite.Equal(expectedPickupDays, summary.PickupDays)

	expectedDeliveryDays := []time.Time{originalDeliveryDate}
	suite.Equal(expectedDeliveryDays, summary.DeliveryDays)

	expectedPackDays := []time.Time{
		time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedPackDays, summary.PackDays)

	expectedTransitDays := []time.Time{
		time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
	}
	suite.Equal(expectedTransitDays, summary.TransitDays)
}

//func (suite *HandlerSuite) TestCalculateMoveDatesFromShipmentMostCurrentDateBeforePreviousDate() {
//	// create a shipment
//	requestedPickupDate := time.Date(2018, 12, 4, 0, 0, 0, 0, time.UTC)
//	actualPickupDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)
//	originalPackDate := time.Date(2018, 11, 29, 0, 0, 0, 0, time.UTC)
//	actualPackDate := time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC)
//	originalDeliveryDate := time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC)
//	pmSurveyDeliveryDate := time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC)
//	var shipment = models.Shipment{
//		RequestedPickupDate:  &requestedPickupDate,
//		OriginalPackDate: &originalPackDate,
//		ActualPickupDate:  &actualPickupDate,
//		ActualPackDate: &actualPackDate,
//		OriginalDeliveryDate:    &originalDeliveryDate,
//
//		PmSurveyPlannedDeliveryDate:    &pmSurveyDeliveryDate,
//	}
//
//	summary, err := calculateMoveDatesFromShipment(&shipment)
//
//	// check that there is no error
//	suite.Nil(err)
//	// compare expected output with actual output
//	expectedPickupDays := []time.Time{actualPickupDate}
//	suite.Equal(expectedPickupDays, summary.PickupDays)
//
//	expectedDeliveryDays := []time.Time{pmSurveyDeliveryDate}
//	suite.Equal(expectedDeliveryDays, summary.DeliveryDays)
//
//	expectedPackDays := []time.Time{
//		time.Date(2018, 11, 29, 0, 0, 0, 0, time.UTC),
//		time.Date(2018, 11, 30, 0, 0, 0, 0, time.UTC),
//		time.Date(2018, 12, 3, 0, 0, 0, 0, time.UTC),
//	}
//	suite.Equal(expectedPackDays, summary.PackDays)
//
//	expectedTransitDays := []time.Time{
//		time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
//		//time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
//		//time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
//		//time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
//		//time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
//	}
//	suite.Equal(expectedTransitDays, summary.TransitDays)
//}
