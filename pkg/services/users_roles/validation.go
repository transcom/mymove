package usersroles

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type usersRolesValidator interface {
	Validate(appCtx appcontext.AppContext, newUsersRoles, originalUsersRoles *[]models.UsersRoles) error
}

type usersRolesValidatorFunc func(appcontext.AppContext, *[]models.UsersRoles, *[]models.UsersRoles) error

func (fn usersRolesValidatorFunc) Validate(appCtx appcontext.AppContext, newUsersRoles, originalUsersRoles *[]models.UsersRoles) error {
	return fn(appCtx, newUsersRoles, originalUsersRoles)
}

func validateUsersRoles(appCtx appcontext.AppContext, newUsersRoles, originalUsersRoles *[]models.UsersRoles, checks ...usersRolesValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newUsersRoles, originalUsersRoles); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// Accumulate all validation errors
				verrs.Append(e)
			default:
				// Non-validation errors have priority and short-circuit doing any further checks
				return err
			}
		}
	}

	if verrs.HasAny() {
		if newUsersRoles != nil {
			// Looping over the users roles will not be viable with a return, return the first ID
			return apperror.NewInvalidInputError((*newUsersRoles)[0].ID, nil, verrs, "Invalid input found while validating the usersRoles.")
		}
		return apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "Invalid input found while validating the usersRoles.")
	}
	return nil
}
