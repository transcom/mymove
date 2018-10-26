package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCalculateMoveDatesFromShipment() {
	// create a shipment
	pickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	transitDays := int64(5)
	packDays := int64(3)
	var shipment = Shipment{
		RequestedPickupDate:  &pickupDate,
		EstimatedTransitDays: &transitDays,
		EstimatedPackDays:    &packDays,
	}

	summary, err := CalculateMoveDatesFromShipment(&shipment)

	// check that there is no error
	suite.Nil(err)
	// compare expected output with actual output
	expectedPickupDays := []time.Time{pickupDate}
	suite.Equal(expectedPickupDays, summary.PickupDays)

	deliveryDate := time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC)
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

func (suite *ModelSuite) TestCalculateMoveDatesFromShipmentMissingPickupDate() {
	// create a shipment
	transitDays := int64(5)
	packDays := int64(3)
	var shipment = Shipment{
		EstimatedTransitDays: &transitDays,
		EstimatedPackDays:    &packDays,
	}

	_, err := CalculateMoveDatesFromShipment(&shipment)

	suite.Error(err)
}

func (suite *ModelSuite) TestCalculateMoveDatesFromShipmentMissingTransitDays() {
	// create a shipment
	pickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	packDays := int64(3)
	var shipment = Shipment{
		RequestedPickupDate: &pickupDate,
		EstimatedPackDays:   &packDays,
	}
	_, err := CalculateMoveDatesFromShipment(&shipment)

	suite.Error(err)
}

func (suite *ModelSuite) TestCalculateMoveDatesFromShipmentMissingPackDays() {
	// create a shipment
	pickupDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	transitDays := int64(5)
	var shipment = Shipment{
		RequestedPickupDate:  &pickupDate,
		EstimatedTransitDays: &transitDays,
	}

	_, err := CalculateMoveDatesFromShipment(&shipment)

	suite.Error(err)
}
