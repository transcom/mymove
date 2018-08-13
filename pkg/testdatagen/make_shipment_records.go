package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeShipment creates a single shipment record
func MakeShipment(db *pop.Connection, assertions Assertions) models.Shipment {

	requestedPickupDate := assertions.Shipment.RequestedPickupDate
	pickupDate := assertions.Shipment.PickupDate
	deliveryDate := assertions.Shipment.DeliveryDate
	tdl := assertions.Shipment.TrafficDistributionList
	if tdl == nil {
		// TODO: Fix TDL
		// tdl = MakeDefaultTDL(db)
	}

	defaultGBLOC := "OHAI"
	sourceGBLOC := assertions.Shipment.SourceGBLOC
	if sourceGBLOC == nil {
		sourceGBLOC = &defaultGBLOC
	}
	destinationGBLOC := assertions.Shipment.DestinationGBLOC
	if destinationGBLOC == nil {
		destinationGBLOC = &defaultGBLOC
	}

	market := assertions.Shipment.Market

	move := assertions.Shipment.Move
	if move == nil {
		newMove := MakeDefaultMove(db)
		move = &newMove
	}

	serviceMember := assertions.Shipment.ServiceMember
	if serviceMember == nil {
		newServiceMember := MakeDefaultServiceMember(db)
		serviceMember = &newServiceMember
	}

	codeOfService := assertions.Shipment.CodeOfService
	if codeOfService == nil {
		newCodeOfService := "D"
		codeOfService = &newCodeOfService
	}

	pickupAddress := assertions.Shipment.PickupAddress
	if pickupAddress == nil {
		newPickupAddress := MakeAddress(db, Assertions{})
		pickupAddress = &newPickupAddress
	}

	status := assertions.Shipment.Status
	if status == "" {
		status = "DEFAULT"
	}

	shipment := models.Shipment{
		TrafficDistributionListID:    uuidPointer(tdl.ID),
		TrafficDistributionList:      tdl,
		ServiceMemberID:              serviceMember.ID,
		ServiceMember:                serviceMember,
		PickupDate:                   timePointer(*pickupDate),
		DeliveryDate:                 timePointer(*deliveryDate),
		SourceGBLOC:                  stringPointer(*sourceGBLOC),
		DestinationGBLOC:             stringPointer(*destinationGBLOC),
		Market:                       market,
		CodeOfService:                codeOfService,
		BookDate:                     timePointer(DateInsidePerformancePeriod),
		RequestedPickupDate:          timePointer(*requestedPickupDate),
		MoveID:                       move.ID,
		Move:                         move,
		Status:                       status,
		EstimatedPackDays:            models.Int64Pointer(2),
		EstimatedTransitDays:         models.Int64Pointer(3),
		PickupAddressID:              &pickupAddress.ID,
		PickupAddress:                pickupAddress,
		HasSecondaryPickupAddress:    false,
		SecondaryPickupAddressID:     nil,
		SecondaryPickupAddress:       nil,
		HasDeliveryAddress:           false,
		DeliveryAddressID:            nil,
		DeliveryAddress:              nil,
		HasPartialSITDeliveryAddress: false,
		PartialSITDeliveryAddressID:  nil,
		PartialSITDeliveryAddress:    nil,
		WeightEstimate:               poundPointer(2000),
		ProgearWeightEstimate:        poundPointer(225),
		SpouseProgearWeightEstimate:  poundPointer(312),
	}

	verrs, err := db.ValidateAndSave(&shipment)
	if verrs.HasAny() {
		err = fmt.Errorf("shipment validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

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
			PickupDate:              &now,
			DeliveryDate:            &now,
			TrafficDistributionList: &tdlList[0],
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &nowPlusOne,
			PickupDate:              &nowPlusOne,
			DeliveryDate:            &nowPlusOne,
			TrafficDistributionList: &tdlList[1],
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &nowPlusTwo,
			PickupDate:              &nowPlusTwo,
			DeliveryDate:            &nowPlusTwo,
			TrafficDistributionList: &tdlList[2],
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})
}
