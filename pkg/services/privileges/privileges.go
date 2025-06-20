package privileges

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type privilegesFetcher struct {
}

func NewPrivilegesFetcher() services.PrivilegeFetcher {
	return privilegesFetcher{}
}

func (f privilegesFetcher) FetchPrivilegeTypes(appCtx appcontext.AppContext) ([]roles.PrivilegeType, error) {
	var privilegeTypes []roles.PrivilegeType
	err := appCtx.DB().RawQuery("SELECT DISTINCT privilege_type FROM privileges").All(&privilegeTypes)
	return privilegeTypes, err
}
