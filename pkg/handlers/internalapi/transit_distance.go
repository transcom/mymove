package internalapi

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

// TransitDistance gets the transit distance from a move ID and planner
func TransitDistance(db *pop.Connection, moveID uuid.UUID, planner route.Planner) (*models.Move, *int, error) {
	// FetchMoveForMoveDates will get all the required associations used below.
	move, err := models.FetchMoveForMoveDates(db, moveID)
	if err != nil {
		return nil, nil, err
	}

	// Error if addresses are missing
	if move.Orders.ServiceMember.DutyStation.Address == (models.Address{}) {
		return nil, nil, errors.New("DutyStation must have an address")
	}
	if move.Orders.NewDutyStation.Address == (models.Address{}) {
		return nil, nil, errors.New("NewDutyStation must have an address")
	}

	var source = move.Orders.ServiceMember.DutyStation.Address
	var destination = move.Orders.NewDutyStation.Address

	// Get the transit distance
	transitDistance, err := planner.TransitDistance(&source, &destination)
	if err != nil {
		return &move, nil, err
	}
	return &move, &transitDistance, nil
}
