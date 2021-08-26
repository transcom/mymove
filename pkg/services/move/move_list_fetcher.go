package move

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type moveListFetcher struct {
	builder moveListQueryBuilder
}

// FetchMoveList uses the passed query builder to fetch a list of moves
func (o *moveListFetcher) FetchMoveList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.Moves, error) {
	var moves models.Moves
	error := o.builder.FetchMany(appCtx, &moves, filters, associations, pagination, ordering)
	return moves, error
}

// FetchMoveCount uses the passed query builder to count moves
func (o *moveListFetcher) FetchMoveCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var moves models.Moves
	count, error := o.builder.Count(appCtx, &moves, filters)
	return count, error
}

// NewMoveListFetcher returns an implementation of OfficeUserListFetcher
func NewMoveListFetcher(builder moveListQueryBuilder) services.MoveListFetcher {
	return &moveListFetcher{builder}
}
