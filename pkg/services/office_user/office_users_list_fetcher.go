package officeuser

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUsersListQueryBuilder interface {
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type officeUserListFetcher struct {
	builder officeUsersListQueryBuilder
}

// FetchOfficeUserList uses the passed query builder to fetch a list of office users
func (o *officeUserListFetcher) FetchOfficeUsersList(appCtx appcontext.AppContext, filterFuncs []func(*pop.Query), pagination services.Pagination, ordering services.QueryOrder) (models.OfficeUsers, int, error) {
	var query *pop.Query
	var officeUsers models.OfficeUsers

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

	query = query.Where("status = ?", models.OfficeUserStatusAPPROVED)
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

	err := query.Paginate(pagination.Page(), pagination.PerPage()).All(&officeUsers)
	if err != nil {
		return nil, 0, err
	}

	if orderTerm == "transportation_office_id" {
		if order == "desc" {
			sort.Slice(officeUsers, func(i, j int) bool {
				return officeUsers[i].TransportationOffice.Name > officeUsers[j].TransportationOffice.Name
			})
		} else {
			sort.Slice(officeUsers, func(i, j int) bool {
				return officeUsers[i].TransportationOffice.Name < officeUsers[j].TransportationOffice.Name
			})
		}
	}

	count := query.Paginator.TotalEntriesSize
	return officeUsers, count, nil
}

// FetchOfficeUserList uses the passed query builder to fetch a list of office users
func (o *officeUserListFetcher) FetchOfficeUsersCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var officeUsers models.OfficeUsers
	count, err := o.builder.Count(appCtx, &officeUsers, filters)
	return count, err
}

// NewOfficecUserListFetcher returns an implementation of OfficeUserListFetcher
func NewOfficeUsersListFetcher(builder officeUsersListQueryBuilder) services.OfficeUserListFetcher {
	return &officeUserListFetcher{builder}
}
