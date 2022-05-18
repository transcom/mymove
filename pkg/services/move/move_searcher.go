package move

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveSearcher struct {
}

func NewMoveSearcher() services.MoveSearcher {
	return &moveSearcher{}
}

func (s moveSearcher) SearchMoves(appCtx appcontext.AppContext, locator *string, dodID *string) (models.Moves, error) {
	if locator == nil && dodID == nil {
		return models.Moves{}, fmt.Errorf("need at least one search filter")
	}
	if locator != nil && dodID != nil {
		return models.Moves{}, fmt.Errorf("SearchMoves requires exactly one search filter")
	}

	query := appCtx.DB().EagerPreload(
		"MTOShipments",
		"Orders.ServiceMember",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
	).
		Join("orders", "orders.id = moves.orders_id").
		Join("service_members", "service_members.id = orders.service_member_id").
		Where("show = TRUE")

	if locator != nil {
		query = query.Where("locator = ?", locator)
	}

	if dodID != nil {
		query = query.Where("service_members.edipi = ?", *dodID)
	}

	var moves models.Moves
	err := query.All(&moves)

	if err != nil {
		return models.Moves{}, apperror.NewQueryError("Move", err, "")
	}
	return moves, nil
}
