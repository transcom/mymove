package move

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveListQueryBuilder interface {
	FetchMany(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) (int, error)
}

type moveListFetcher struct {
	builder moveListQueryBuilder
}

// FetchMoveList uses the passed query builder to fetch a list of moves
func (o *moveListFetcher) FetchMoveList(appCfg appconfig.AppConfig, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.Moves, error) {
	var moves models.Moves
	error := o.builder.FetchMany(appCfg, &moves, filters, associations, pagination, ordering)
	return moves, error
}

// FetchMoveCount uses the passed query builder to count moves
func (o *moveListFetcher) FetchMoveCount(appCfg appconfig.AppConfig, filters []services.QueryFilter) (int, error) {
	var moves models.Moves
	count, error := o.builder.Count(appCfg, &moves, filters)
	return count, error
}

// NewMoveListFetcher returns an implementation of OfficeUserListFetcher
func NewMoveListFetcher(builder moveListQueryBuilder) services.MoveListFetcher {
	return &moveListFetcher{builder}
}
