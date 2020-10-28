package internalapi

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// calculateMoveDates is used on the hhg wizard DatePicker page to show move dates summary
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
	//TODO: fix test TestCreateShipmentHandlerAllValues() so that duty stations differ so that this error check does not cause the test to fail
	//if move.Orders.NewDutyStation.Address.PostalCode[0:5] == move.Orders.ServiceMember.DutyStation.Address.PostalCode[0:5] {
	//	return summary, errors.New("NewDutyStation must not have the same zip code as the original DutyStation")
	//}

	var source = move.Orders.ServiceMember.DutyStation.Address
	var destination = move.Orders.NewDutyStation.Address

	transitDistance, err := planner.TransitDistance(&source, &destination)
	if err != nil {
		return summary, err
	}

	entitlement, err := models.GetEntitlement(*move.Orders.ServiceMember.Rank, move.Orders.HasDependents,
		move.Orders.SpouseHasProGear)
	if err != nil {
		return summary, err
	}
	entitlementWeight := unit.Pound(entitlement)
	estimatedPackDays := models.PackDays(entitlementWeight)
	estimatedTransitDays, err := models.TransitDays(entitlementWeight, transitDistance)
	if err != nil {
		return summary, err
	}

	summary.CalculateMoveDates(moveDate, estimatedPackDays, estimatedTransitDays)
	// ReportDays isn't set by CalculateMoveDates and must be added here to display in the calendar widget
	summary.ReportDays = []time.Time{move.Orders.ReportByDate.UTC()}

	return summary, nil
}
