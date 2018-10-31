package dates

import (
	"time"
)

// MoveDatesSummary contains the set of dates for a move
type MoveDatesSummary struct {
	MoveDate             time.Time
	PackDays             []time.Time
	EstimatedPackDays    int
	PickupDays           []time.Time
	TransitDays          []time.Time
	EstimatedTransitDays int
	DeliveryDays         []time.Time
	ReportDays           []time.Time
}

// CalculateMoveDates returns a MoveDatesSummary based on the expected move date and estimated pack and transit days
func (summary *MoveDatesSummary) CalculateMoveDates(moveDate time.Time, estimatedPackDays int, estimatedTransitDays int) {
	usCalendar := NewUSCalendar()

	summary.MoveDate = moveDate
	summary.EstimatedPackDays = estimatedPackDays
	summary.EstimatedTransitDays = estimatedTransitDays

	lastPossiblePackDay := moveDate.AddDate(0, 0, -1)
	summary.PackDays = CreatePastMoveDates(lastPossiblePackDay, estimatedPackDays, false, usCalendar)

	firstPossiblePickupDay := moveDate
	pickupDays := CreateFutureMoveDates(firstPossiblePickupDay, 1, false, usCalendar)
	summary.PickupDays = pickupDays

	firstPossibleTransitDay := time.Time(pickupDays[len(pickupDays)-1]).AddDate(0, 0, 1)
	transitDays := CreateFutureMoveDates(firstPossibleTransitDay, estimatedTransitDays, true, usCalendar)
	summary.TransitDays = transitDays

	firstPossibleDeliveryDay := time.Time(transitDays[len(transitDays)-1]).AddDate(0, 0, 1)
	summary.DeliveryDays = CreateFutureMoveDates(firstPossibleDeliveryDay, 1, false, usCalendar)
}
