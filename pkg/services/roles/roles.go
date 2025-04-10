package roles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type rolesFetcher struct {
}

func NewRolesFetcher() services.RoleAssociater {
	return rolesFetcher{}
}

// FetchRolesForUser associates a given user with a set of roles
func (f rolesFetcher) FetchRolesForUser(appCtx appcontext.AppContext, userID uuid.UUID) (roles.Roles, error) {
	var roles roles.Roles
	err := appCtx.DB().Q().Join("users_roles", "users_roles.role_id = roles.id").
		Where("users_roles.deleted_at IS NULL AND users_roles.user_id = ?", (userID)).
		All(&roles)
	return roles, err
}

func (f rolesFetcher) FetchRolesPrivileges(appCtx appcontext.AppContext) ([]models.RolePrivilege, error) {
	var rolesPrivileges []models.RolePrivilege
	err := appCtx.DB().Q().EagerPreload("Role", "Privilege").All(&rolesPrivileges)
	return rolesPrivileges, err
}

func (f rolesFetcher) FetchRoleTypes(appCtx appcontext.AppContext) ([]roles.RoleType, error) {
	var roleTypes []roles.RoleType
	err := appCtx.DB().RawQuery("SELECT DISTINCT role_type FROM roles").All(&roleTypes)
	return roleTypes, err
}

func (f rolesFetcher) VerifyRolesPrivelegesAllowed(appCtx appcontext.AppContext, roleType *string, privilegeType *string) (bool, error) {
	var results []models.RolePrivilege
	query := appCtx.DB().RawQuery("select * from roles_privileges join roles on roles_privileges.role_id = roles.id join privileges on roles_privileges.privilege_id = privileges.id where role_type = $1 and privilege_type = $2", roleType, privilegeType)

	err := query.All(&results)

	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return false, nil
	}
	return true, nil
}
