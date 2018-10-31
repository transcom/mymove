package internalapi

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

func calculateMoveDatesFromMove(db *pop.Connection, planner route.Planner, moveID uuid.UUID, moveDate time.Time) (dates.MoveDatesSummary, error) {
	var summary dates.MoveDatesSummary

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

	estimatedPackDays := models.PackDays(entitlementWeight)
	estimatedTransitDays, err := models.TransitDays(entitlementWeight, transitDistance)
	if err != nil {
		return summary, err
	}

	summary.CalculateMoveDates(moveDate, estimatedPackDays, estimatedTransitDays)
	summary.ReportDays = []time.Time{move.Orders.ReportByDate.UTC()}

	return summary, nil
}

func calculateMoveDatesFromShipment(shipment *models.Shipment) (dates.MoveDatesSummary, error) {
	usCalendar := dates.NewUSCalendar()

	if shipment.RequestedPickupDate == nil {
		return dates.MoveDatesSummary{}, errors.New("Shipment must have a RequestedPickupDate")
	}
	lastPossiblePackDay := time.Time(*shipment.RequestedPickupDate).AddDate(0, 0, -1)

	if shipment.EstimatedPackDays == nil {
		return dates.MoveDatesSummary{}, errors.New("Shipment must have a EstimatedPackDays")
	}
	packDates := dates.CreatePastMoveDates(lastPossiblePackDay, int(*shipment.EstimatedPackDays), false, usCalendar)

	pickupDates := dates.CreateFutureMoveDates(*shipment.RequestedPickupDate, 1, false, usCalendar)

	firstPossibleTransitDay := time.Time(pickupDates[len(pickupDates)-1]).AddDate(0, 0, 1)
	if shipment.EstimatedTransitDays == nil {
		return dates.MoveDatesSummary{}, errors.New("Shipment must have EstimatedTransitDays")
	}
	transitDates := dates.CreateFutureMoveDates(firstPossibleTransitDay, int(*shipment.EstimatedTransitDays), true, usCalendar)

	firstPossibleDeliveryDay := time.Time(transitDates[int(*shipment.EstimatedTransitDays)-1].AddDate(0, 0, 1))
	deliveryDates := dates.CreateFutureMoveDates(firstPossibleDeliveryDay, 1, false, usCalendar)

	summary := dates.MoveDatesSummary{
		PackDays:     packDates,
		PickupDays:   pickupDates,
		TransitDays:  transitDates,
		DeliveryDays: deliveryDates,
	}
	return summary, nil
}
