package adminuser

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type requestedOfficeUsersListQueryBuilder interface {
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type requestedOfficeUserListFetcher struct {
	builder requestedOfficeUsersListQueryBuilder
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *requestedOfficeUserListFetcher) FetchRequestedOfficeUsersList(appCtx appcontext.AppContext, filterFuncs []func(*pop.Query), pagination services.Pagination, ordering services.QueryOrder) (models.OfficeUsers, int, error) {
	var query *pop.Query
	var requestedUsers models.OfficeUsers

	query = appCtx.DB().Q().EagerPreload(
		"User.Roles",
		"TransportationOffice").
		Join("transportation_offices", "office_users.transportation_office_id = transportation_offices.id")

	for _, filterFunc := range filterFuncs {
		filterFunc(query)
	}

	query = query.Where("status = ?", models.OfficeUserStatusREQUESTED)

	err := query.Paginate(pagination.Page(), pagination.PerPage()).All(&requestedUsers)
	if err != nil {
		return nil, 0, err
	}

	count := query.Paginator.TotalEntriesSize

	return requestedUsers, count, nil
}

// FetchAdminUserList uses the passed query builder to fetch a list of office users
func (o *requestedOfficeUserListFetcher) FetchRequestedOfficeUsersCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var requestedUsers models.OfficeUsers
	count, err := o.builder.Count(appCtx, &requestedUsers, filters)
	return count, err
}

// NewAdminUserListFetcher returns an implementation of AdminUserListFetcher
func NewRequestedOfficeUsersListFetcher(builder requestedOfficeUsersListQueryBuilder) services.RequestedOfficeUserListFetcher {
	return &requestedOfficeUserListFetcher{builder}
}
