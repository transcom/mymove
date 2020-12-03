package movetaskorder

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"

	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderHider struct {
	db *pop.Connection
}

// NewMoveTaskOrderHider creates a new struct with the service dependencies
func NewMoveTaskOrderHider(db *pop.Connection) services.MoveTaskOrderHider {
	return &moveTaskOrderHider{db}
}

// Hide hides any MTO that isn't using valid fake data
func (o *moveTaskOrderHider) Hide() (models.Moves, error) {
	var mtos models.Moves
	err := o.db.Q().
		Where("show = ?", swag.Bool(true)).
		All(&mtos)
	if err != nil {
		return nil, services.NewQueryError("Moves", err, fmt.Sprintf("Could not find move task orders: %s", err))
	}

	var invalidFakeMoves models.Moves
	for _, mto := range mtos {
		isValid, _ := fakedata.IsValidFakeServiceMember(mto.Orders.ServiceMember)
		if !isValid {
			dontShow := false
			mto.Show = &dontShow
			invalidFakeMoves = append(invalidFakeMoves, mto)
		}
	}

	return invalidFakeMoves, nil
}
