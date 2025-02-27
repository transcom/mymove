package adminuser

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type rejectedOfficeUsersListQueryBuilder interface {
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type rejectedOfficeUserListFetcher struct {
	builder rejectedOfficeUsersListQueryBuilder
}

// FetchRejectedUserList uses the passed query builder to fetch a list of office users
func (o *rejectedOfficeUserListFetcher) FetchRejectedOfficeUsersList(appCtx appcontext.AppContext, filterFuncs []func(*pop.Query), pagination services.Pagination, ordering services.QueryOrder) (models.OfficeUsers, int, error) {
	var query *pop.Query
	var rejectedUsers models.OfficeUsers

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

	query = query.Where("status = ?", models.OfficeUserStatusREJECTED)
	query.GroupBy("office_users.id")

	var order = "desc"
	if ordering.SortOrder() != nil && *ordering.SortOrder() {
		order = "asc"
	}

	var orderTerm = "id"
	if ordering.Column() != nil {
		orderTerm = *ordering.Column()
	}

	if orderTerm == "role" {
		if order == "asc" {
			query = query.Order("MIN(roles.role_name) ASC")
		} else {
			query = query.Order("MIN(roles.role_name) DESC")
		}
	} else {
		query = query.Order(fmt.Sprintf("%s %s", orderTerm, order))
	}

	query.Select("office_users.*")

	err := query.Paginate(pagination.Page(), pagination.PerPage()).All(&rejectedUsers)
	if err != nil {
		return nil, 0, err
	}

	for i := range rejectedUsers {
		sort.Slice(rejectedUsers[i].User.Roles, func(a, b int) bool {
			return rejectedUsers[i].User.Roles[a].RoleName < rejectedUsers[i].User.Roles[b].RoleName
		})
	}

	if orderTerm == "transportation_office_id" {
		if order == "desc" {
			sort.Slice(rejectedUsers, func(i, j int) bool {
				return rejectedUsers[i].TransportationOffice.Name > rejectedUsers[j].TransportationOffice.Name
			})
		} else {
			sort.Slice(rejectedUsers, func(i, j int) bool {
				return rejectedUsers[i].TransportationOffice.Name < rejectedUsers[j].TransportationOffice.Name
			})
		}
	}

	count := query.Paginator.TotalEntriesSize
	return rejectedUsers, count, nil
}

// FetchRejectedUserList uses the passed query builder to fetch a list of office users
func (o *rejectedOfficeUserListFetcher) FetchRejectedOfficeUsersCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var rejectedUsers models.OfficeUsers
	count, err := o.builder.Count(appCtx, &rejectedUsers, filters)
	return count, err
}

// NewRejectedUserListFetcher returns an implementation of RejectedUserListFetcher
func NewRejectedOfficeUsersListFetcher(builder rejectedOfficeUsersListQueryBuilder) services.RejectedOfficeUserListFetcher {
	return &rejectedOfficeUserListFetcher{builder}
}
