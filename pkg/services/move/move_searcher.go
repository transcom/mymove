package move

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

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

	query := appCtx.DB().EagerPreload("MTOShipments", "Orders.ServiceMember").
		Join("orders", "orders.id = moves.orders_id").
		Join("service_members", "service_members.id = orders.service_member_id").
		Where("show = TRUE")

	if locator != nil {
		var moves models.Moves
		err := query.Where("locator = $1", locator).All(&moves)

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				// Not found error expects an id but we're querying by locator
				return models.Moves{}, apperror.NewNotFoundError(uuid.Nil, "move locator "+*locator)
			default:
				return models.Moves{}, apperror.NewQueryError("Move", err, "")
			}
		}
		return moves, nil
	}

	if dodID != nil {
		return models.Moves{}, nil
	}

	return models.Moves{}, fmt.Errorf("not implemented")
}
