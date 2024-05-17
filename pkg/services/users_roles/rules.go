package usersroles

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func checkTransportationOfficerPolicyViolation() usersRolesValidator {
	return usersRolesValidatorFunc(func(appCtx appcontext.AppContext, newUsersRoles, _ *[]models.UsersRoles) error {
		verrs := validate.NewErrors()

		hasTOO := false
		hasTIO := false

		// Check the new roles being added
		if newUsersRoles != nil {
			for _, newRole := range *newUsersRoles {
				var role roles.Role
				err := appCtx.DB().Find(&role, newRole.RoleID)
				if err != nil {
					return err
				}
				if role.RoleType == roles.RoleTypeTOO {
					hasTOO = true
					if hasTIO {
						verrs.Add("role", "a user cannot have both the TOO and TIO roles")
						break
					}
				}
				if role.RoleType == roles.RoleTypeTIO {
					hasTIO = true
					if hasTOO {
						verrs.Add("role", "a user cannot have both the TOO and TIO roles")
						break
					}
				}
			}
		}

		if verrs.HasAny() {
			return verrs
		}

		return nil
	})
}
