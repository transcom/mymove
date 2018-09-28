package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
	"math"
	"time"
)

// ShowUnavailableMoveDatesHandler returns the unavailable move dates starting at a given date.
type ShowUnavailableMoveDatesHandler struct {
	handlers.HandlerContext
}

// Handle returns the unavailable move dates.
func (h ShowUnavailableMoveDatesHandler) Handle(params calendarop.ShowUnavailableMoveDatesParams) middleware.Responder {
	startDate := time.Time(params.StartDate)

	var datesPayload []strfmt.Date
	datesPayload = append(datesPayload, strfmt.Date(startDate)) // The start date is always unavailable.

	const daysToCheck = 90
	const shortFuseTotalDays = 5
	daysChecked := 0
	shortFuseDaysFound := 0

	// TODO: Handle holidays.
	firstPossibleDate := startDate.AddDate(0, 0, 1)
	for d := firstPossibleDate; daysChecked < daysToCheck; d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			datesPayload = append(datesPayload, strfmt.Date(d))
		} else if shortFuseDaysFound < shortFuseTotalDays {
			datesPayload = append(datesPayload, strfmt.Date(d))
			shortFuseDaysFound++
		}
		daysChecked++
	}

	return calendarop.NewShowUnavailableMoveDatesOK().WithPayload(datesPayload)
}

// ShowMoveDatesSummaryHandler returns a summary of the dates in the move process given a move date and move ID.
type ShowMoveDatesSummaryHandler struct {
	handlers.HandlerContext
}

// Handle returns a summary of the dates in the move process.
func (h ShowMoveDatesSummaryHandler) Handle(params calendarop.ShowMoveDatesSummaryParams) middleware.Responder {
	startDate := time.Time(params.MoveDate)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// FetchMoveForMoveDates will get all the required associations used below.
	move, err := models.FetchMoveForMoveDates(h.DB(), moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	transitDistance, err := h.Planner().TransitDistance(&move.Orders.ServiceMember.DutyStation.Address,
		&move.Orders.NewDutyStation.Address)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	entitlementWeight := models.GetEntitlement(*move.Orders.ServiceMember.Rank, move.Orders.HasDependents,
		move.Orders.SpouseHasProGear)

	numTransitDays, err := models.TransitDays(unit.Pound(entitlementWeight), transitDistance)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	const poundsPerDay = 5000
	numPackDays := int(math.Ceil(float64(entitlementWeight) / float64(poundsPerDay)))

	packDays := createMoveDates(startDate, numPackDays, false)
	firstPossiblePickupDay := time.Time(packDays[len(packDays)-1]).AddDate(0, 0, 1)
	pickupDays := createMoveDates(firstPossiblePickupDay, 1, false)
	firstPossibleTransitDay := time.Time(pickupDays[len(pickupDays)-1])
	transitDays := createMoveDates(firstPossibleTransitDay, numTransitDays, false)
	firstPossibleDeliveryDay := time.Time(transitDays[len(transitDays)-1]).AddDate(0, 0, 1)
	deliveryDays := createMoveDates(firstPossibleDeliveryDay, 1, false)
	reportDays := []strfmt.Date{strfmt.Date(move.Orders.ReportByDate.UTC())}

	moveDatesSummaryPayload := &internalmessages.MoveDatesSummaryPayload{
		Pack:     packDays,
		Pickup:   pickupDays,
		Transit:  transitDays,
		Delivery: deliveryDays,
		Report:   reportDays,
	}

	return calendarop.NewShowMoveDatesSummaryOK().WithPayload(moveDatesSummaryPayload)
}

func createMoveDates(startDate time.Time, numDays int, includeWeekendsAndHolidays bool) []strfmt.Date {
	var dates []strfmt.Date

	// TODO: Handle holidays.
	daysAdded := 0
	for d := startDate; daysAdded < numDays; d = d.AddDate(0, 0, 1) {
		if includeWeekendsAndHolidays || (d.Weekday() != time.Saturday && d.Weekday() != time.Sunday) {
			dates = append(dates, strfmt.Date(d))
			daysAdded++
		}
	}

	return dates
}
