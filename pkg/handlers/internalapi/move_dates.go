package internalapi

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
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

// calculateMoveDates is used on the hhg wizard DatePicker page to show move dates summary
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

// calculateMoveDatesFromShipment takes stored values on the shipment to calculate the most up-to-date move date ranges
// this is used to display date ranges for the SM HHG review page and the status timeline on the post-hhg-submission landing page
func calculateMoveDatesFromShipment(shipment *models.Shipment) (MoveDatesSummary, error) {
	usCalendar := handlers.NewUSCalendar()

	if shipment.EstimatedPackDays == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have EstimatedPackDays")
	}

	if shipment.RequestedPickupDate == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have a RequestedPickupDate")
	}

	if shipment.EstimatedTransitDays == nil {
		return MoveDatesSummary{}, errors.New("Shipment must have EstimatedTransitDays")
	}

	var mostCurrentPackDate time.Time
	if shipment.ActualPackDate != nil {
		mostCurrentPackDate = *shipment.ActualPackDate
	} else if shipment.PmSurveyPlannedPackDate != nil {
		mostCurrentPackDate = *shipment.PmSurveyPlannedPackDate
	} else if shipment.OriginalPackDate != nil {
		mostCurrentPackDate = *shipment.OriginalPackDate
	} else {
		lastPossiblePackDay := shipment.RequestedPickupDate.AddDate(0, 0, -1)
		mostCurrentPackDate = createPastMoveDates(lastPossiblePackDay, int(*shipment.EstimatedPackDays), false, usCalendar)[0]
	}

	var mostCurrentPickupDate time.Time
	if shipment.ActualPickupDate != nil {
		mostCurrentPickupDate = *shipment.ActualPickupDate
	} else if shipment.PmSurveyPlannedPickupDate != nil {
		mostCurrentPickupDate = *shipment.PmSurveyPlannedPickupDate
	} else {
		mostCurrentPickupDate = *shipment.RequestedPickupDate
	}

	var mostCurrentDeliveryDate time.Time

	if shipment.ActualDeliveryDate != nil {
		mostCurrentDeliveryDate = *shipment.ActualDeliveryDate
	} else if shipment.PmSurveyPlannedDeliveryDate != nil {
		mostCurrentDeliveryDate = *shipment.PmSurveyPlannedDeliveryDate
	} else if shipment.OriginalDeliveryDate != nil {
		mostCurrentDeliveryDate = *shipment.OriginalDeliveryDate
	} else {
		// transit days can be on weekends and holidays and delivery cannot, so calculations must be separated out
		estimatedTransitDates := createFutureMoveDates(*shipment.RequestedPickupDate, int(*shipment.EstimatedTransitDays), true, usCalendar)
		lastEstimatedTransitDate := estimatedTransitDates[len(estimatedTransitDates)-1]
		mostCurrentDeliveryDate = createFutureMoveDates(lastEstimatedTransitDate.AddDate(0, 0, 1), 1, false, usCalendar)[0]
	}
	// assigns the pack dates
	packDates, err := createValidDatesBetweenTwoDates(mostCurrentPackDate, mostCurrentPickupDate, false, true, usCalendar)
	if err != nil {
		return MoveDatesSummary{}, err
	}
	pickupDates := createFutureMoveDates(mostCurrentPickupDate, 1, false, usCalendar)

	firstPossibleTransitDay := time.Time(pickupDates[len(pickupDates)-1]).AddDate(0, 0, 1)

	transitDates, err := createValidDatesBetweenTwoDates(firstPossibleTransitDay, mostCurrentDeliveryDate, true, true, usCalendar)
	if err != nil {
		return MoveDatesSummary{}, err
	}
	deliveryDates := createFutureMoveDates(mostCurrentDeliveryDate, 1, false, usCalendar)

	summary := MoveDatesSummary{
		PackDays:     packDates,
		PickupDays:   pickupDates,
		TransitDays:  transitDates,
		DeliveryDays: deliveryDates,
	}
	return summary, nil
}

// createFutureMoveDates is used when a startDate is known, an end date is not, number of days into the future is known.  Inclusive of startDate.
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

// createPastMoveDates is used when a startDate is known and number of days into the past is known.  Inclusive of startDate.
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

func createValidDatesBetweenTwoDates(startDate time.Time, endDate time.Time, includeWeekendsAndHolidays bool, allowEarlierOrSameEndDate bool, calendar *cal.Calendar) ([]time.Time, error) {
	// returns date range inclusive of startDate, exclusive of endDate (unless endDate is before or equal to startDate and allowEarlierEndDate)
	var dates []time.Time

	if !calendar.IsWorkday(endDate) && !includeWeekendsAndHolidays {
		return dates, errors.New("End date cannot be a weekend or holiday")
	}

	if startDate.After(endDate) || startDate == endDate {
		if allowEarlierOrSameEndDate == true {
			return dates, nil
		}
		return dates, errors.New("End date cannot be before or equal to start date")
	}

	dateToAdd := startDate

	for dateToAdd.Before(endDate) {
		if includeWeekendsAndHolidays || calendar.IsWorkday(dateToAdd) {
			dates = append(dates, dateToAdd)
		}
		dateToAdd = dateToAdd.AddDate(0, 0, 1)
	}
	return dates, nil
}
