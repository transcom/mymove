package officeuser

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	QueryForAssociations(appCtx appcontext.AppContext, model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type officeUserFetcher struct {
	builder officeUserQueryBuilder
}

// FetchOfficeUser fetches an office user given a slice of filters
func (o *officeUserFetcher) FetchOfficeUser(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := o.builder.FetchOne(appCtx, &officeUser, filters)
	return officeUser, err
}

// NewOfficeUserFetcher return an implementation of the OfficeUserFetcher interface
func NewOfficeUserFetcher(builder officeUserQueryBuilder) services.OfficeUserFetcher {
	return &officeUserFetcher{builder}
}

// TODO - Eventually move away from the query builder and back to pop
type officeUserFetcherPop struct {
}

// FetchOfficeUserByID fetches an office user given a slice of filters
func (o *officeUserFetcherPop) FetchOfficeUserByID(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := appCtx.DB().Eager("TransportationOffice").Find(&officeUser, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.OfficeUser{}, apperror.NewNotFoundError(id, "looking for OfficeUser")
		default:
			return models.OfficeUser{}, apperror.NewQueryError("OfficeUser", err, "")
		}
	}

	return officeUser, err
}

func (o *officeUserFetcherPop) FetchOfficeUserByRoleAndGbloc(appCtx appcontext.AppContext, role roles.RoleType, gbloc string) ([]models.OfficeUser, error) {
	// init office users array
	var officeUsers []models.OfficeUser

	err := appCtx.DB().EagerPreload(
		"User",
		"User.Roles",
		"User.Privileges",
		// "OfficeUser",
		"TransportationOffice",
		"TransportationOffice.Gbloc",
	).
		Join("users", "users.id = office_users.user_id").
		Join("users_roles", "users.id = users_roles.user_id").
		Join("roles", "users_roles.role_id = roles.id").
		Join("transportation_offices", "office_users.transportation_office_id = transportation_offices.id").
		Where("gbloc = ?", gbloc).
		Where("role_type = ?", role).
		All(&officeUsers)
	// err := appCtx.DB().EagerPreload(
	// 	"User",
	// 	"User.Roles",
	// 	"User.Privileges",
	// 	// "OfficeUser",
	// 	"TransportationOffice",
	// 	"TransportationOffice.Gbloc",
	// ).
	// 	Join("users", "users.id = office_users.user_id").
	// 	Join("transportation_offices", "office_users.transportation_office_id = transportation_offices.id").
	// 	Where("gbloc = ?", gbloc).
	// 	All(&officeUsers)
	// ).Join("transportation_offices", "transportation_offices.id = office_users.transportation_office_id").All(&officeUsers)

	if err != nil {
		return nil, err
	}

	return officeUsers, nil

	//	select ou.id, to2.gbloc, r.role_name, p.privilege_name  from office_users ou
	//
	// join users u on u.id = ou.user_id
	// join users_roles ur on ur.user_id = u.id
	// join roles r on r.id = ur.role_id
	// join transportation_offices to2 on ou.transportation_office_id = to2.id
	// join users_privileges up on u.id = up.user_id
	// join "privileges" p on p.id =up.privilege_id
	// where ou.id = 'a1018959-9523-44a1-8505-312c669533f5'
	// run query
}

// NewOfficeUserFetcherPop return an implementation of the OfficeUserFetcherPop interface
func NewOfficeUserFetcherPop() services.OfficeUserFetcherPop {
	return &officeUserFetcherPop{}
}
