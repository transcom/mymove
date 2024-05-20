package clientcert

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// updatePrimeRoleForUser updates the prime role for the user
// associated with the client cert. If the client cert is nil, the
// cert has been deleted. If the user has no other client certs
// associated, the prime role will be removed
func updatePrimeRoleForUser(appCtx appcontext.AppContext, userID uuid.UUID, clientCert *models.ClientCert, _ clientCertQueryBuilder, associator services.UserRoleAssociator) error {

	userRoles, err := roles.FetchRolesForUser(appCtx.DB(), userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	allowPrime := false

	if clientCert != nil && clientCert.AllowPrime {
		allowPrime = true
	}
	if clientCert == nil {
		var clientCertCount int
		err = appCtx.DB().RawQuery("SELECT count(*) FROM client_certs WHERE allow_prime AND user_id = ?",
			userID.String()).First(&clientCertCount)
		if err != nil {
			return err
		}
		if clientCertCount > 0 {
			// if the number of client certs associated with this user
			// that have prime access is greater than 0, the user
			// should have the prime role
			allowPrime = true
		} else {
			allowPrime = false
		}
	}

	if !allowPrime && userRoles.HasRole(roles.RoleTypePrime) {
		// remove prime role
		newRoles := []roles.RoleType{}
		for _, role := range userRoles {
			if role.RoleType != roles.RoleTypePrime {
				newRoles = append(newRoles, role.RoleType)
			}
			_, _, err = associator.UpdateUserRoles(appCtx, userID, newRoles)
			if err != nil {
				return err
			}
		}
	}

	if allowPrime && !userRoles.HasRole(roles.RoleTypePrime) {
		// add prime role
		newRoles := []roles.RoleType{}
		for _, role := range userRoles {
			newRoles = append(newRoles, role.RoleType)
		}
		newRoles = append(newRoles, roles.RoleTypePrime)
		_, _, err = associator.UpdateUserRoles(appCtx, userID, newRoles)
		if err != nil {
			return err
		}
	}

	return nil
}
