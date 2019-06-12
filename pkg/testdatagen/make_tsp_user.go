package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTspUser creates a single TSP user and associated Transportation Service Provider
func MakeTspUser(db *pop.Connection, assertions Assertions) models.TspUser {
	user := assertions.TspUser.User
	email := "leo_spaceman_tsp@example.com"

	if assertions.TspUser.UserID == nil || isZeroUUID(*assertions.TspUser.UserID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		user = MakeUser(db, assertions)
	}

	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	var tspAssertions = assertions.TransportationServiceProvider
	tsp := MakeTSP(db, Assertions{
		TransportationServiceProvider: tspAssertions,
	})

	tspUser := models.TspUser{
		UserID:                          &user.ID,
		User:                            user,
		TransportationServiceProvider:   tsp,
		TransportationServiceProviderID: tsp.ID,
		FirstName:                       "Leo",
		LastName:                        "Spaceman",
		Email:                           email,
		Telephone:                       "415-555-1212",
	}

	mergeModels(&tspUser, assertions.TspUser)

	mustCreate(db, &tspUser)

	return tspUser
}

// MakeDefaultTspUser makes an TspUser with default values
func MakeDefaultTspUser(db *pop.Connection) models.TspUser {
	return MakeTspUser(db, Assertions{})
}
