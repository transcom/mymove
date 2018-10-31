package testdatagen

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/dates"
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
	if isZeroUUID(assertions.Shipment.ServiceMemberID) {
		serviceMember = move.Orders.ServiceMember
	}

	pickupAddress := assertions.Shipment.PickupAddress
	if pickupAddress == nil {
		newPickupAddress := MakeDefaultAddress(db)
		pickupAddress = &newPickupAddress
	}

	hasDeliveryAddress := assertions.Shipment.HasDeliveryAddress
	deliveryAddress := assertions.Shipment.DeliveryAddress
	if deliveryAddress == nil && hasDeliveryAddress {
		newDeliveryAddress := MakeAddress2(db, Assertions{})
		deliveryAddress = &newDeliveryAddress
	}

	status := assertions.Shipment.Status
	if status == "" {
		status = models.ShipmentStatusDRAFT
	}

	sourceGBLOC := assertions.Shipment.SourceGBLOC
	if sourceGBLOC == nil {
		sourceGBLOC = stringPointer(DefaultSrcGBLOC)
	}

	destinationGBLOC := assertions.Shipment.DestinationGBLOC
	if destinationGBLOC == nil {
		destinationGBLOC = stringPointer(DefaultDstGBLOC)
	}

	var summary dates.MoveDatesSummary
	summary.CalculateMoveDates(PerformancePeriodStart, 2, 3)

	shipment := models.Shipment{
		Status:           status,
		SourceGBLOC:      sourceGBLOC,
		DestinationGBLOC: destinationGBLOC,
		GBLNumber:        nil,
		Market:           &DefaultMarket,

		// associations
		TrafficDistributionListID: uuidPointer(tdl.ID),
		TrafficDistributionList:   tdl,
		ServiceMemberID:           serviceMember.ID,
		ServiceMember:             serviceMember,
		MoveID:                    move.ID,
		Move:                      move,

		// dates
		ActualPickupDate:     nil,
		ActualPackDate:       nil,
		ActualDeliveryDate:   nil,
		BookDate:             timePointer(DateInsidePerformancePeriod),
		OriginalPackDate:     timePointer(summary.PackDays[0]),
		RequestedPickupDate:  timePointer(summary.MoveDate),
		OriginalDeliveryDate: timePointer(summary.DeliveryDays[0]),

		// calculated durations
		EstimatedPackDays:    models.Int64Pointer(int64(summary.EstimatedPackDays)),
		EstimatedTransitDays: models.Int64Pointer(int64(summary.EstimatedTransitDays)),

		// addresses
		PickupAddressID:              &pickupAddress.ID,
		PickupAddress:                pickupAddress,
		HasSecondaryPickupAddress:    false,
		SecondaryPickupAddressID:     nil,
		SecondaryPickupAddress:       nil,
		HasDeliveryAddress:           hasDeliveryAddress,
		DeliveryAddressID:            nil,
		DeliveryAddress:              nil,
		HasPartialSITDeliveryAddress: false,
		PartialSITDeliveryAddressID:  nil,
		PartialSITDeliveryAddress:    nil,

		// weights
		WeightEstimate:              poundPointer(2000),
		ProgearWeightEstimate:       poundPointer(225),
		SpouseProgearWeightEstimate: poundPointer(312),
		NetWeight:                   nil,
		GrossWeight:                 nil,
		TareWeight:                  nil,

		// pre-move survey
		PmSurveyConductedDate:               nil,
		PmSurveyPlannedPackDate:             nil,
		PmSurveyPlannedPickupDate:           nil,
		PmSurveyPlannedDeliveryDate:         nil,
		PmSurveyWeightEstimate:              nil,
		PmSurveyProgearWeightEstimate:       nil,
		PmSurveySpouseProgearWeightEstimate: nil,
		PmSurveyNotes:                       nil,
		PmSurveyMethod:                      "",
	}

	if hasDeliveryAddress {
		shipment.DeliveryAddressID = &deliveryAddress.ID
		shipment.DeliveryAddress = deliveryAddress
	}

	// Overwrite values with those from assertions
	mergeModels(&shipment, assertions.Shipment)

	mustCreate(db, &shipment)

	shipment.Move.Shipments = append(shipment.Move.Shipments, shipment)

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
	oneWeek := time.Hour * 168
	now := time.Now()
	nowPlusOne := now.Add(oneWeek)
	nowPlusTwo := now.Add(oneWeek * 2)
	market := "dHHG"
	sourceGBLOC := "KKFA"
	destinationGBLOC := "HAFC"

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &now,
			ActualPickupDate:        &now,
			ActualDeliveryDate:      &now,
			TrafficDistributionList: &tdlList[0],
			SourceGBLOC:             &sourceGBLOC,
			DestinationGBLOC:        &destinationGBLOC,
			Market:                  &market,
		},
	})

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &nowPlusOne,
			ActualPickupDate:        &nowPlusOne,
			ActualDeliveryDate:      &nowPlusOne,
			TrafficDistributionList: &tdlList[1],
			SourceGBLOC:             &sourceGBLOC,
			DestinationGBLOC:        &destinationGBLOC,
			Market:                  &market,
		},
	})

	MakeShipment(db, Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate:     &nowPlusTwo,
			ActualPickupDate:        &nowPlusTwo,
			ActualDeliveryDate:      &nowPlusTwo,
			TrafficDistributionList: &tdlList[2],
			SourceGBLOC:             &sourceGBLOC,
			DestinationGBLOC:        &destinationGBLOC,
			Market:                  &market,
		},
	})
}
