package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeShipment creates a single shipment record.
func MakeShipment(db *pop.Connection, requestedPickup time.Time,
	pickup time.Time, delivery time.Time,
	tdl models.TrafficDistributionList, sourceGBLOC string, market *string, move *models.Move) (models.Shipment, error) {

	var moveObj models.Move
	if move == nil {
		moveObj = MakeDefaultMove(db)
	} else {
		moveObj = *move
	}
	pickupAddress := MakeAddress(db, Assertions{})
	codeOfService := "D"
	destinationGBLOC := "OHAI"
	shipment := models.Shipment{
		TrafficDistributionListID: uuidPointer(tdl.ID),
		TrafficDistributionList:   &tdl,
		PickupDate:                timePointer(pickup),
		DeliveryDate:              timePointer(delivery),
		// CreatedAt
		// UpdatedAt
		SourceGBLOC:                  stringPointer(sourceGBLOC),
		DestinationGBLOC:             stringPointer(destinationGBLOC),
		Market:                       market,
		CodeOfService:                &codeOfService,
		BookDate:                     timePointer(DateInsidePerformancePeriod),
		RequestedPickupDate:          timePointer(requestedPickup),
		MoveID:                       moveObj.ID,
		Move:                         &moveObj,
		Status:                       "DEFAULT",
		EstimatedPackDays:            models.Int64Pointer(2),
		EstimatedTransitDays:         models.Int64Pointer(3),
		PickupAddressID:              &pickupAddress.ID,
		PickupAddress:                &pickupAddress,
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

	return shipment, err
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
	now := time.Now()
	thirtyMin, _ := time.ParseDuration("30m")
	oneWeek, _ := time.ParseDuration("7d")
	market := "dHHG"
	sourceGBLOC := "OHAI"

	MakeShipment(db, now, now, now.Add(thirtyMin), tdlList[0], sourceGBLOC, &market, nil)
	MakeShipment(db, now.Add(oneWeek), now.Add(oneWeek), now.Add(oneWeek).Add(thirtyMin), tdlList[1], sourceGBLOC, &market, nil)
	MakeShipment(db, now.Add(oneWeek*2), now.Add(oneWeek*2), now.Add(oneWeek*2).Add(thirtyMin), tdlList[2], sourceGBLOC, &market, nil)
}
