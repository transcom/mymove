package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/rickar/cal"
	"github.com/transcom/mymove/pkg/unit"
)

// MoveDatesSummary contains the set of dates for a move
type MoveDatesSummary struct {
	PackDays             []time.Time
	EstimatedPackDays    int64
	PickupDays           []time.Time
	TransitDays          []time.Time
	EstimatedTransitDays int64
	DeliveryDays         []time.Time
	ReportDays           []time.Time
}

// CalculateMoveDates will calculate the MoveDatesSummary given a Move object
func CalculateMoveDates(db *pop.Connection, move *Move, moveDate time.Time, estimatedPackDays *int64, estimatedTransitDays *int64) (MoveDatesSummary, error) {
	var summary MoveDatesSummary

	if estimatedPackDays == nil || estimatedTransitDays == nil {
		// Calculate the expected transit days given the entitlement and distance
		entitlementWeight := unit.Pound(GetEntitlement(*move.Orders.ServiceMember.Rank, move.Orders.HasDependents,
			move.Orders.SpouseHasProGear))

		// Calculate the expected packing and transit days given the entitlement
		if estimatedPackDays == nil {
			packDays := PackDays(entitlementWeight)
			estimatedPackDays = &packDays
		}
		if estimatedTransitDays == nil && move.TransitDistance != nil {
			transitDays, err := TransitDays(entitlementWeight, *move.TransitDistance)
			if err != nil {
				return summary, err
			}
			estimatedTransitDays = &transitDays
		}
	}

	// Calculate dates in the calendar
	usCalendar := NewUSCalendar()

	// MoveDate is the RequestedPickupDate
	lastPossiblePackDay := moveDate.AddDate(0, 0, -1)

	// Packing Days
	summary.PackDays = createPastMoveDates(lastPossiblePackDay, *estimatedPackDays, false, usCalendar)
	summary.EstimatedPackDays = *estimatedPackDays

	// Pickup Days
	firstPossiblePickupDay := moveDate
	pickupDays := createFutureMoveDates(firstPossiblePickupDay, 1, false, usCalendar)
	summary.PickupDays = pickupDays

	// Transit Days
	firstPossibleTransitDay := time.Time(pickupDays[len(pickupDays)-1]).AddDate(0, 0, 1)
	transitDays := createFutureMoveDates(firstPossibleTransitDay, *estimatedTransitDays, true, usCalendar)
	summary.TransitDays = transitDays
	summary.EstimatedTransitDays = *estimatedTransitDays

	// Delivery Days
	firstPossibleDeliveryDay := time.Time(transitDays[len(transitDays)-1]).AddDate(0, 0, 1)
	summary.DeliveryDays = createFutureMoveDates(firstPossibleDeliveryDay, 1, false, usCalendar)

	// Report Days
	summary.ReportDays = []time.Time{move.Orders.ReportByDate.UTC()}

	return summary, nil
}

// CalculateMoveDatesFromShipment will calculate the MoveDatesSummary given a Shipment object
func CalculateMoveDatesFromShipment(db *pop.Connection, shipment *Shipment) (MoveDatesSummary, error) {
	// Error checking
	if shipment.RequestedPickupDate == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have a RequestedPickupDate")
	}

	moveDate := time.Time(*shipment.RequestedPickupDate)
	return CalculateMoveDates(db, &shipment.Move, moveDate, shipment.EstimatedPackDays, shipment.EstimatedTransitDays)
}

func createFutureMoveDates(startDate time.Time, numDays int64, includeWeekendsAndHolidays bool, calendar *cal.Calendar) []time.Time {
	dates := make([]time.Time, 0, numDays)

	daysAdded := int64(0)
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, 1) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(d) {
			dates = append(dates, d)
			daysAdded++
		}
	}

	return dates
}

func createPastMoveDates(startDate time.Time, numDays int64, includeWeekendsAndHolidays bool, calendar *cal.Calendar) []time.Time {
	dates := make([]time.Time, numDays)

	daysAdded := int64(0)
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, -1) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(d) {
			// Since we're working backwards, put dates at end of slice.
			dates[numDays-daysAdded-1] = d
			daysAdded++
		}
	}

	return dates
}
