package privileges

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type privilegesFetcher struct {
}

func NewPrivilegesFetcher() services.PrivilegeFetcher {
	return privilegesFetcher{}
}

// FetchPrivilegesForUser associates a given user with a set of privileges
func (f privilegesFetcher) FetchPrivilegesForUser(appCtx appcontext.AppContext, userID uuid.UUID) (roles.Privileges, error) {
	var privileges roles.Privileges
	err := appCtx.DB().Q().Join("users_privileges", "users_privileges.privilege_id = privileges.id").
		Where("users_privileges.deleted_at IS NULL AND users_privileges.user_id = ?", (userID)).
		All(&privileges)
	return privileges, err
}

func (f privilegesFetcher) FetchPrivilegeTypes(appCtx appcontext.AppContext) ([]roles.PrivilegeType, error) {
	var privilegeTypes []roles.PrivilegeType
	err := appCtx.DB().RawQuery("SELECT DISTINCT privilege_type FROM privileges").All(&privilegeTypes)
	return privilegeTypes, err
}
