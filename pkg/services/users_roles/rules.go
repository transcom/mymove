package usersroles

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func checkTransportationOfficerPolicyViolation() usersRolesValidator {
	return usersRolesValidatorFunc(func(appCtx appcontext.AppContext, newUsersRoles, existingUserRoles *[]models.UsersRoles) (*validate.Errors, error) {
		verrs := validate.NewErrors()

		hasTOO := false
		hasTIO := false

		// Check the existing roles
		if existingUserRoles != nil {
			for _, existingRole := range *existingUserRoles {
				var role roles.Role
				err := appCtx.DB().Find(&role, existingRole.RoleID)
				if err != nil {
					return nil, err
				}
				// Save if the existing role found was a TOO or TIO
				if role.RoleType == roles.RoleTypeTOO {
					hasTOO = true
				}
				if role.RoleType == roles.RoleTypeTIO {
					hasTIO = true
				}
			}
		}

		// Check the new roles being added
		if newUsersRoles != nil {
			for _, newRole := range *newUsersRoles {
				var role roles.Role
				err := appCtx.DB().Find(&role, newRole.RoleID)
				if err != nil {
					return nil, err
				}
				if role.RoleType == roles.RoleTypeTOO {
					if hasTIO {
						verrs.Add("roles", "a user cannot have both the TOO and TIO roles")
						break
					}
					hasTOO = true
				}
				if role.RoleType == roles.RoleTypeTIO {
					if hasTOO {
						verrs.Add("roles", "a user cannot have both the TOO and TIO roles")
						break
					}
					hasTIO = true
				}
			}
		}

		if verrs.HasAny() {
			return verrs, nil
		}

		return nil, nil
	})
}
