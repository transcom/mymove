package roles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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
