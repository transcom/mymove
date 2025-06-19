package adminuser

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
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
		Join("users", "users.id = office_users.user_id").
		Join("users_roles", "users.id = users_roles.user_id").
		Join("roles", "users_roles.role_id = roles.id").
		Join("transportation_offices", "office_users.transportation_office_id = transportation_offices.id")

	for _, filterFunc := range filterFuncs {
		filterFunc(query)
	}

	query = query.Where("status = ?", models.OfficeUserStatusREQUESTED)
	query.GroupBy("office_users.id")

	var order = "desc"
	if ordering.SortOrder() != nil && *ordering.SortOrder() {
		order = "asc"
	}

	var orderTerm = "id"
	if ordering.Column() != nil {
		orderTerm = *ordering.Column()
	}

	query.Order(fmt.Sprintf("%s %s", orderTerm, order))
	query.Select("office_users.*")

	err := query.Paginate(pagination.Page(), pagination.PerPage()).All(&requestedUsers)
	if err != nil {
		return nil, 0, err
	}

	if orderTerm == "transportation_office_id" {
		if order == "desc" {
			sort.Slice(requestedUsers, func(i, j int) bool {
				return requestedUsers[i].TransportationOffice.Name > requestedUsers[j].TransportationOffice.Name
			})
		} else {
			sort.Slice(requestedUsers, func(i, j int) bool {
				return requestedUsers[i].TransportationOffice.Name < requestedUsers[j].TransportationOffice.Name
			})
		}
	}
	for i := range requestedUsers {
		var liveRoles []roles.Role
		err := appCtx.DB().Q().
			Join("users_roles", "users_roles.role_id = roles.id").
			Where("users_roles.user_id = ?", requestedUsers[i].User.ID).
			Where("users_roles.deleted_at IS NULL").
			All(&liveRoles)
		if err != nil {
			return nil, 0, err
		}
		requestedUsers[i].User.Roles = liveRoles
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
