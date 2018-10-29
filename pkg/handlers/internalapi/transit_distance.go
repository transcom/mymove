package internalapi

import (
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

// TransitDistance gets the transit distance from a move ID and planner
func TransitDistance(move *models.Move, planner route.Planner) (*int, error) {

	// Error if addresses are missing
	if move.Orders.ServiceMember.DutyStation.Address == (models.Address{}) {
		return nil, errors.New("DutyStation must have an address")
	}
	if move.Orders.NewDutyStation.Address == (models.Address{}) {
		return nil, errors.New("NewDutyStation must have an address")
	}

	var source = move.Orders.ServiceMember.DutyStation.Address
	var destination = move.Orders.NewDutyStation.Address

	// Get the transit distance
	transitDistance, err := planner.TransitDistance(&source, &destination)
	if err != nil {
		return nil, err
	}
	return &transitDistance, nil
}
