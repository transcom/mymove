package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
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

	requestedPickupDate := assertions.Shipment.RequestedPickupDate
	if requestedPickupDate == nil {
		requestedPickupDate = &PerformancePeriodStart
	}
	var summary dates.MoveDatesSummary
	summary.CalculateMoveDates(*requestedPickupDate, 2, 3)

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

// MakeShipmentForPricing makes a Shipment with all other depdendent models that would be expected in a real scenario
func MakeShipmentForPricing(db *pop.Connection, assertions Assertions) (models.Shipment, error) {
	var shipment models.Shipment

	// Shipment must have a NetWeight
	if assertions.Shipment.NetWeight == nil {
		weight := unit.Pound(5000)
		assertions.Shipment.NetWeight = &weight
	}

	shipment = MakeShipment(db, assertions)

	// Zip3 records must exist for pickup and delivery addresses
	addresses := []models.Address{
		*shipment.PickupAddress,
		shipment.Move.Orders.NewDutyStation.Address,
	}
	for _, a := range addresses {
		newAssertions := assertions
		newAssertions.Tariff400ngZip3.Zip3 = zip5ToZip3(a.PostalCode)
		zip3 := FetchOrMakeTariff400ngZip3(db, newAssertions)

		// Service area values must match between Zip3 and ServiceArea
		newAssertions.Tariff400ngServiceArea.ServiceArea = zip3.ServiceArea
		MakeTariff400ngServiceArea(db, newAssertions)
	}

	return shipment, nil
}
