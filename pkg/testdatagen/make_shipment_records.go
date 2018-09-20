package testdatagen

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeShipment creates a single shipment record
func MakeShipment(db *pop.Connection, assertions Assertions) models.Shipment {
	tdl := assertions.Shipment.TrafficDistributionList
	if tdl == nil {
		newTDL := MakeDefaultTDL(db)
		tdl = &newTDL
	}

	move := assertions.Shipment.Move
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Shipment.MoveID) {
		newMove := MakeMove(db, assertions)
		move = newMove
	}

	serviceMember := assertions.Shipment.ServiceMember
	if serviceMember == nil {
		serviceMember = &move.Orders.ServiceMember
	}

	pickupAddress := assertions.Shipment.PickupAddress
	if pickupAddress == nil {
		newPickupAddress := MakeDefaultAddress(db)
		pickupAddress = &newPickupAddress
	}

	deliveryAddress := assertions.Shipment.DeliveryAddress
	if deliveryAddress == nil {
		newDeliveryAddress := MakeAddress2(db, Assertions{})
		deliveryAddress = &newDeliveryAddress
	}

	status := assertions.Shipment.Status
	if status == "" {
		status = models.ShipmentStatusDRAFT
	}

	shipment := models.Shipment{
		TrafficDistributionListID:    uuidPointer(tdl.ID),
		TrafficDistributionList:      tdl,
		ServiceMemberID:              serviceMember.ID,
		ServiceMember:                serviceMember,
		ActualPickupDate:             timePointer(DateInsidePerformancePeriod),
		DeliveryDate:                 timePointer(DateOutsidePerformancePeriod),
		SourceGBLOC:                  stringPointer(DefaultSrcGBLOC),
		DestinationGBLOC:             stringPointer(DefaultSrcGBLOC),
		Market:                       &DefaultMarket,
		BookDate:                     timePointer(DateInsidePerformancePeriod),
		RequestedPickupDate:          timePointer(PerformancePeriodStart),
		MoveID:                       move.ID,
		Move:                         move,
		Status:                       models.ShipmentStatusDRAFT,
		EstimatedPackDays:            models.Int64Pointer(2),
		EstimatedTransitDays:         models.Int64Pointer(3),
		PickupAddressID:              &pickupAddress.ID,
		PickupAddress:                pickupAddress,
		HasSecondaryPickupAddress:    false,
		SecondaryPickupAddressID:     nil,
		SecondaryPickupAddress:       nil,
		HasDeliveryAddress:           true,
		DeliveryAddressID:            &deliveryAddress.ID,
		DeliveryAddress:              deliveryAddress,
		HasPartialSITDeliveryAddress: false,
		PartialSITDeliveryAddressID:  nil,
		PartialSITDeliveryAddress:    nil,
		WeightEstimate:               poundPointer(2000),
		ProgearWeightEstimate:        poundPointer(225),
		SpouseProgearWeightEstimate:  poundPointer(312),
	}

	// Overwrite values with those from assertions
	mergeModels(&shipment, assertions.Shipment)

	mustCreate(db, &shipment)

	return shipment
}

// MakeDefaultShipment makes a Shipment with default values
func MakeDefaultShipment(db *pop.Connection) models.Shipment {
	return MakeShipment(db, Assertions{})
}

// MakeShipmentData creates three shipment records
func MakeShipmentData(db *pop.Connection) {
	// Grab three UUIDs for individual TDLs
	// TODO: should this query be made in main, between creation functions,
	// and then sourced from one central place?
	tdlList := []models.TrafficDistributionList{}
	err := db.All(&tdlList)
	if err != nil {
		fmt.Println("TDL ID import failed.")
	}

	// Add three shipment table records using UUIDs from TDLs
	oneWeek, _ := time.ParseDuration("7d")
	now := time.Now()
	nowPlusOne := now.Add(oneWeek)
	nowPlusTwo := now.Add(oneWeek * 2)
	market := "dHHG"
	sourceGBLOC := "OHAI"

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &now,
			ActualPickupDate:        &now,
			DeliveryDate:            &now,
			TrafficDistributionList: &tdlList[0],
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &nowPlusOne,
			ActualPickupDate:        &nowPlusOne,
			DeliveryDate:            &nowPlusOne,
			TrafficDistributionList: &tdlList[1],
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &nowPlusTwo,
			ActualPickupDate:        &nowPlusTwo,
			DeliveryDate:            &nowPlusTwo,
			TrafficDistributionList: &tdlList[2],
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})
}
