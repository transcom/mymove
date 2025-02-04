package officeuser

import (
	"database/sql"
	"fmt"

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
	DeleteOne(appCtx appcontext.AppContext, model interface{}) error
	DeleteMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
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

func (o *officeUserFetcherPop) FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := appCtx.DB().Eager("TransportationOffice", "TransportationOfficeAssignments", "TransportationOfficeAssignments.TransportationOffice").Find(&officeUser, id)
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

// Fetch office users of the same role within a gbloc, for assignment purposes
func (o *officeUserFetcherPop) FetchOfficeUsersByRoleAndOffice(appCtx appcontext.AppContext, role roles.RoleType, officeID uuid.UUID) ([]models.OfficeUser, error) {
	var officeUsers []models.OfficeUser

	err := appCtx.DB().EagerPreload(
		"User",
		"User.Roles",
		"User.Privileges",
		"TransportationOffice",
		"TransportationOffice.Gbloc",
	).
		Join("users", "users.id = office_users.user_id").
		Join("users_roles", "users.id = users_roles.user_id").
		Join("roles", "users_roles.role_id = roles.id").
		Where("transportation_office_id = ?", officeID).
		Where("role_type = ?", role).
		Where("users_roles.deleted_at IS NULL").
		Where("office_users.active = TRUE").
		Order("last_name asc").
		All(&officeUsers)

	if err != nil {
		return nil, err
	}

	return officeUsers, nil
}

func (o *officeUserFetcherPop) FetchSafetyMoveOfficeUsersByRoleAndOffice(appCtx appcontext.AppContext, role roles.RoleType, officeID uuid.UUID) ([]models.OfficeUser, error) {
	var officeUsers []models.OfficeUser

	err := appCtx.DB().EagerPreload(
		"User",
		"User.Roles",
		"User.Privileges",
		"TransportationOffice",
		"TransportationOffice.Gbloc",
	).
		Join("users", "users.id = office_users.user_id").
		Join("users_roles", "users.id = users_roles.user_id").
		Join("roles", "users_roles.role_id = roles.id").
		LeftJoin("users_privileges", "users.id = users_privileges.user_id").
		LeftJoin("privileges", "privileges.id = users_privileges.privilege_id").
		Where("transportation_office_id = ?", officeID).
		Where("role_type = ?", role).
		Where("users_roles.deleted_at IS NULL").
		Where("office_users.active = TRUE").
		Where("users_privileges.deleted_at IS NULL").
		Where("privileges.privilege_type = 'safety'").
		Order("last_name asc").
		All(&officeUsers)

	if err != nil {
		return nil, err
	}

	return officeUsers, nil
}

// Fetch office users of the same role within a gbloc, with their workload, for assignment purposes
func (o *officeUserFetcherPop) FetchOfficeUsersWithWorkloadByRoleAndOffice(appCtx appcontext.AppContext, role roles.RoleType, officeID uuid.UUID) ([]models.OfficeUserWithWorkload, error) {
	var officeUsers []models.OfficeUserWithWorkload

	query :=
		`SELECT office_users.id,
			office_users.first_name,
			office_users.last_name,
			COUNT(DISTINCT moves.id) AS workload
		FROM office_users
		JOIN users_roles ON office_users.user_id = users_roles.user_id
		JOIN roles ON users_roles.role_id = roles.id
		JOIN transportation_offices ON office_users.transportation_office_id = transportation_offices.id
		LEFT JOIN moves
			ON (
				(roles.role_type = 'services_counselor' AND moves.sc_assigned_id = office_users.id) OR
				(roles.role_type = 'task_ordering_officer' AND moves.too_assigned_id = office_users.id) OR
				(roles.role_type = 'task_invoicing_officer' and moves.tio_assigned_id = office_users.id)
			)
		WHERE roles.role_type = $1
			AND transportation_offices.id = $2
			AND office_users.active = TRUE
		GROUP BY office_users.id, office_users.first_name, office_users.last_name`

	err := appCtx.DB().RawQuery(query, role, officeID).All(&officeUsers)
	if err != nil {
		return nil, fmt.Errorf("error fetching moves for office: %s with error %w", officeID, err)
	}

	return officeUsers, nil
}

// NewOfficeUserFetcherPop return an implementation of the OfficeUserFetcherPop interface
func NewOfficeUserFetcherPop() services.OfficeUserFetcherPop {
	return &officeUserFetcherPop{}
}
