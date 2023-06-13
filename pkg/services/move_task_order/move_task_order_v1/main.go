package movetaskorderv1

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// In this option we don't make a change to the old version of the service.
// ListPrimeMoveTaskOrders performs an optimized fetch for moves specifically targeting the Prime API.
func ListPrimeMoveTaskOrders(f services.MoveTaskOrderFetcher, appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	var moveTaskOrders models.Moves
	var err error

	sql := `SELECT moves.*
            FROM moves INNER JOIN orders ON moves.orders_id = orders.id
            WHERE moves.available_to_prime_at IS NOT NULL AND moves.show = TRUE`

	if searchParams != nil && searchParams.Since != nil {
		sql = sql + ` AND (moves.updated_at >= $1 OR orders.updated_at >= $1 OR
                          (moves.id IN (SELECT mto_shipments.move_id
                                        FROM mto_shipments WHERE mto_shipments.updated_at >= $1
                                        UNION
                                        SELECT mto_service_items.move_id
			                            FROM mto_service_items
			                            WHERE mto_service_items.updated_at >= $1
			                            UNION
			                            SELECT payment_requests.move_id
			                            FROM payment_requests
			                            WHERE payment_requests.updated_at >= $1)));`
		err = appCtx.DB().RawQuery(sql, *searchParams.Since).All(&moveTaskOrders)
	} else {
		sql = sql + `;`
		err = appCtx.DB().RawQuery(sql).All(&moveTaskOrders)
	}

	if err != nil {
		return models.Moves{}, apperror.NewQueryError("MoveTaskOrder", err, "Unexpected error while querying db.")
	}

	return moveTaskOrders, nil
}
