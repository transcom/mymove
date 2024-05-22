package usersroles

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type usersRolesValidator interface {
	Validate(appCtx appcontext.AppContext, newUsersRoles, originalUsersRoles *[]models.UsersRoles) (*validate.Errors, error)
}

type usersRolesValidatorFunc func(appcontext.AppContext, *[]models.UsersRoles, *[]models.UsersRoles) (*validate.Errors, error)

func (fn usersRolesValidatorFunc) Validate(appCtx appcontext.AppContext, newUsersRoles, originalUsersRoles *[]models.UsersRoles) (*validate.Errors, error) {
	return fn(appCtx, newUsersRoles, originalUsersRoles)
}

func validateUsersRoles(appCtx appcontext.AppContext, newUsersRoles, originalUsersRoles *[]models.UsersRoles, checks ...usersRolesValidator) (*validate.Errors, error) {
	for _, check := range checks {
		if verrs, err := check.Validate(appCtx, newUsersRoles, originalUsersRoles); err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return nil, nil
}
