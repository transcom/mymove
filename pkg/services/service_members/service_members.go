package servicemembers

import (
	"github.com/gobuffalo/validate/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type serviceMemberCreator struct {
	services.UserRoleAssociator
}

// NewServiceMemberCreator creates a new struct with the service dependencies
func NewServiceMemberCreator() services.ServiceMemberAssociator {
	return serviceMemberCreator{}
}

func (s serviceMemberCreator) CreateServiceMember(appCtx appcontext.AppContext, newServiceMember models.ServiceMember) (*validate.Errors, error) {
	smVerrs, err := models.SaveServiceMember(appCtx, &newServiceMember)
	if smVerrs.HasAny() || err != nil {
		return smVerrs, err
	}
	// Create customer user role for new service member
	_, err = s.UpdateUserRoles(appCtx, newServiceMember.UserID, []roles.RoleType{roles.RoleTypeCustomer})
	if err != nil {
		appCtx.Logger().Error("Error updating user roles", zap.Error(err))
		return nil, err
	}
	return nil, nil
}
