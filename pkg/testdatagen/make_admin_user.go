package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAdminUser creates a single admin user and associated TransportOffice
func MakeAdminUser(db *pop.Connection, assertions Assertions) models.AdminUser {
	user := assertions.AdminUser.User
	// There's a uniqueness constraint on admin user emails so add some randomness
	email := fmt.Sprintf("admin_%s@example.com", makeRandomString(5))

	if assertions.AdminUser.UserID == nil || isZeroUUID(*assertions.AdminUser.UserID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		user = MakeUser(db, assertions)
	}

	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	adminUser := models.AdminUser{
		UserID:    &user.ID,
		User:      user,
		FirstName: "Leo",
		LastName:  "Spaceman",
		Email:     email,
		Role:      "SYSTEM_ADMIN",
		Disabled:  false,
	}

	mergeModels(&adminUser, assertions.AdminUser)

	mustCreate(db, &adminUser)

	return adminUser
}

// MakeDefaultAdminUser makes an AdminUser with default values
func MakeDefaultAdminUser(db *pop.Connection) models.AdminUser {
	return MakeAdminUser(db, Assertions{})
}
