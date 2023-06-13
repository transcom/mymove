package movetaskorderfetcher

import (
	"github.com/transcom/mymove/pkg/services"
	movetaskorderfetcherv1 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_fetcher/move_task_order_fetcher_v1"
	movetaskorderfetcherv2 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_fetcher/move_task_order_fetcher_v2"
)

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(apiVersionFlag string) services.MoveTaskOrderFetcher {
	if apiVersionFlag == "v1" {
		return movetaskorderfetcherv1.NewMoveTaskOrderFetcher()
	}
	if apiVersionFlag == "v2" {
		return movetaskorderfetcherv2.NewMoveTaskOrderFetcher()
	}
	return movetaskorderfetcherv1.NewMoveTaskOrderFetcher()

}
