package internalapi

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/rickar/cal"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// MoveDatesSummary contains the set of dates for a move
type MoveDatesSummary struct {
	PackDays     []time.Time
	PickupDays   []time.Time
	TransitDays  []time.Time
	DeliveryDays []time.Time
	ReportDays   []time.Time
}

func calculateMoveDates(db *pop.Connection, planner route.Planner, moveID uuid.UUID, moveDate time.Time) (MoveDatesSummary, error) {
	var summary MoveDatesSummary

	// FetchMoveForMoveDates will get all the required associations used below.
	move, err := models.FetchMoveForMoveDates(db, moveID)
	if err != nil {
		return summary, err
	}

	if move.Orders.ServiceMember.DutyStation.Address == (models.Address{}) {
		return summary, errors.New("DutyStation must have an address")
	}
	if move.Orders.NewDutyStation.Address == (models.Address{}) {
		return summary, errors.New("NewDutyStation must have an address")
	}

	var source = move.Orders.ServiceMember.DutyStation.Address
	var destination = move.Orders.NewDutyStation.Address

	transitDistance, err := planner.TransitDistance(&source, &destination)
	if err != nil {
		return summary, err
	}

	entitlementWeight := unit.Pound(models.GetEntitlement(*move.Orders.ServiceMember.Rank, move.Orders.HasDependents,
		move.Orders.SpouseHasProGear))

	numTransitDays, err := models.TransitDays(entitlementWeight, transitDistance)
	if err != nil {
		return summary, err
	}

	numPackDays := models.PackDays(entitlementWeight)
	usCalendar := handlers.NewUSCalendar()

	lastPossiblePackDay := moveDate.AddDate(0, 0, -1)
	summary.PackDays = createPastMoveDates(lastPossiblePackDay, numPackDays, false, usCalendar)

	firstPossiblePickupDay := moveDate
	pickupDays := createFutureMoveDates(firstPossiblePickupDay, 1, false, usCalendar)
	summary.PickupDays = pickupDays

	firstPossibleTransitDay := time.Time(pickupDays[len(pickupDays)-1]).AddDate(0, 0, 1)
	transitDays := createFutureMoveDates(firstPossibleTransitDay, numTransitDays, true, usCalendar)
	summary.TransitDays = transitDays

	firstPossibleDeliveryDay := time.Time(transitDays[len(transitDays)-1]).AddDate(0, 0, 1)
	summary.DeliveryDays = createFutureMoveDates(firstPossibleDeliveryDay, 1, false, usCalendar)

	summary.ReportDays = []time.Time{move.Orders.ReportByDate.UTC()}

	return summary, nil
}

func calculateMoveDatesFromShipment(shipment *models.Shipment) (MoveDatesSummary, error) {
	usCalendar := handlers.NewUSCalendar()

	if shipment.RequestedPickupDate == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have a RequestedPickupDate")
	}
	lastPossiblePackDay := time.Time(*shipment.RequestedPickupDate).AddDate(0, 0, -1)

	if shipment.EstimatedPackDays == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have a EstimatedPackDays")
	}
	packDates := createPastMoveDates(lastPossiblePackDay, int(*shipment.EstimatedPackDays), false, usCalendar)

	pickupDates := createFutureMoveDates(*shipment.RequestedPickupDate, 1, false, usCalendar)

	firstPossibleTransitDay := time.Time(pickupDates[len(pickupDates)-1]).AddDate(0, 0, 1)
	if shipment.EstimatedTransitDays == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have EstimatedTransitDays")
	}
	transitDates := createFutureMoveDates(firstPossibleTransitDay, int(*shipment.EstimatedTransitDays), true, usCalendar)

	firstPossibleDeliveryDay := time.Time(transitDates[int(*shipment.EstimatedTransitDays)-1].AddDate(0, 0, 1))
	deliveryDates := createFutureMoveDates(firstPossibleDeliveryDay, 1, false, usCalendar)

	summary := MoveDatesSummary{
		PackDays:     packDates,
		PickupDays:   pickupDates,
		TransitDays:  transitDates,
		DeliveryDays: deliveryDates,
	}
	return summary, nil
}

func createFutureMoveDates(startDate time.Time, numDays int, includeWeekendsAndHolidays bool, calendar *cal.Calendar) []time.Time {
	dates := make([]time.Time, 0, numDays)

	daysAdded := 0
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, 1) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(d) {
			dates = append(dates, d)
			daysAdded++
		}
	}

	return dates
}

func createPastMoveDates(startDate time.Time, numDays int, includeWeekendsAndHolidays bool, calendar *cal.Calendar) []time.Time {
	dates := make([]time.Time, numDays)

	daysAdded := 0
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, -1) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(d) {
			// Since we're working backwards, put dates at end of slice.
			dates[numDays-daysAdded-1] = d
			daysAdded++
		}
	}

	return dates
}
